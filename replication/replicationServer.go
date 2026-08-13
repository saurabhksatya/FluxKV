package replication

import (
	"fluxKV/internal"
	"fluxKV/replication/proto"
	"fluxKV/utils"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedReplicateServiceServer

	db       *internal.DataStore
	mu       sync.Mutex
	replicas map[string]proto.ReplicateService_ReplicationConnectServer

	grpc *grpc.Server
}

func NewReplicationServer(db *internal.DataStore) *Server {
	return &Server{
		db:       db,
		replicas: make(map[string]proto.ReplicateService_ReplicationConnectServer),
	}
}

func (s *Server) ReplicationConnect(stream proto.ReplicateService_ReplicationConnectServer) error {
	data, err := stream.Recv()
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.replicas[data.Id] = stream
	s.mu.Unlock()
	utils.Logger.Info("Replica connected: %s with Offset %d", data.Id, data.Offset)

	defer func() {
		s.mu.Lock()
		delete(s.replicas, data.Id)
		s.mu.Unlock()
		utils.Logger.Info("Replica disconnected: %s", data.Id)
	}()

	s.db.OffsetManager.Mutex.RLock()
	oldestOffset := s.db.OffsetManager.Offsets.Front()
	if oldestOffset != nil && oldestOffset.Value.(internal.OperationStructure).Offset < data.Offset {
		s.db.OffsetManager.Mutex.RUnlock()
		// Snapshot and return here
		return status.Error(codes.OutOfRange, "offset out of range")
	}

	for element := s.db.OffsetManager.Offsets.Front(); element != nil; element = element.Next() {
		value := element.Value.(internal.OperationStructure)
		if value.Offset > data.Offset {
			cmd := &proto.ReplicationStreamConnectionResponse{
				Op:       value.Operation,
				Key:      value.Key,
				Value:    value.Value,
				Offset:   value.Offset,
				Snapshot: false,
			}
			if err := stream.Send(cmd); err != nil {
				s.db.OffsetManager.Mutex.RUnlock()
				return err
			}
		}
	}
	s.db.OffsetManager.Mutex.RUnlock()

	<-stream.Context().Done()
	return stream.Context().Err()
}

func (s *Server) Broadcast(operation *internal.OperationStructure) {
	cmd := &proto.ReplicationStreamConnectionResponse{
		Op:       operation.Operation,
		Key:      operation.Key,
		Value:    operation.Value,
		Offset:   operation.Offset,
		Snapshot: false,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for id, stream := range s.replicas {
		if err := stream.Send(cmd); err != nil {
			utils.Logger.Error("Broadcast error: %s for replica id: %s", err, id)
			delete(s.replicas, id)
		}
	}
}

func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.grpc = grpc.NewServer()
	proto.RegisterReplicateServiceServer(s.grpc, s)
	utils.Logger.Info("Replication service listening on %s", addr)
	return s.grpc.Serve(lis)
}

func (s *Server) Stop() {
	if s.grpc != nil {
		s.grpc.GracefulStop()
	}
}
