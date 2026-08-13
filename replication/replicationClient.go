package replication

import (
	"context"
	"errors"
	"fluxKV/internal"
	"fluxKV/replication/proto"
	"fluxKV/utils"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	leaderAddr string
	id         string
	db         *internal.DataStore
	conn       *grpc.ClientConn
}

func NewReplicationClient(leaderADdr, id string, db *internal.DataStore) *Client {
	return &Client{
		leaderAddr: leaderADdr,
		id:         id,
		db:         db,
	}
}

func (rc *Client) run(ctx context.Context) error {
	conn, err := grpc.NewClient(rc.leaderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	stream, err := proto.NewReplicateServiceClient(conn).ReplicationConnect(ctx)
	if err != nil {
		return err
	}
	rc.db.OffsetManager.Mutex.RLock()
	msg := &proto.ReplicationStreamConnectionRequest{
		Id:     rc.id,
		Offset: rc.db.OffsetManager.Offset,
	}

	if err := stream.Send(msg); err != nil {
		return err
	}
	rc.db.OffsetManager.Mutex.RUnlock()
	utils.Logger.Info("Connected to main server at %s", rc.leaderAddr)

	for {
		cmd, err := stream.Recv()
		if err != nil {
			return err
		}

		switch cmd.Op {
		case proto.Command_SET:
			rc.db.Set(cmd.Key, cmd.Value)
		case proto.Command_DEL:
			rc.db.Delete(cmd.Key)
		default:
			utils.Logger.Error("Unknown command %s", cmd.Op)
		}

		if cmd.Snapshot {
			rc.db.OffsetManager.Mutex.Lock()
			rc.db.OffsetManager.Offset = cmd.Offset
			rc.db.OffsetManager.Mutex.Unlock()
		}
	}
}

func (rc *Client) Start(ctx context.Context) {
	for {
		if err := rc.run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			utils.Logger.Error("ReplicationClient failed to start: %s", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
	}
}
