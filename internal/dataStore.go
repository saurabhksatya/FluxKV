package internal

import (
	"context"
	"fluxKV/replication/proto"
	"fluxKV/utils"
)

type ReplicationServer interface {
	Broadcast(operation *OperationStructure)
	Start(addr string) error
	Stop()
}

type ReplicationClient interface {
	Start(ctx context.Context)
}

type DataStore struct {
	data              map[string]string
	OffsetManager     *OffsetManager
	replicationServer ReplicationServer
	replicationClient ReplicationClient
}

func NewDataStore(offsetEnabled bool) *DataStore {
	db := &DataStore{
		data:          make(map[string]string),
		OffsetManager: NewOffsetManager(offsetEnabled),
	}

	go db.Manage()

	return db
}

func (db *DataStore) SetReplicationServer(server ReplicationServer) {
	db.replicationServer = server
}

func (db *DataStore) SetReplicationClient(client ReplicationClient) {
	db.replicationClient = client
}

func (db *DataStore) Manage() {
	for task := range db.OffsetManager.OffsetQueue {
		db.OffsetManager.AddOffset(task)
		if db.replicationServer != nil {
			db.replicationServer.Broadcast(&task)
		}
	}
}

func (db *DataStore) StartReplication(grpcAddr string) {
	if db.replicationServer != nil {
		go func() {
			if err := db.replicationServer.Start(grpcAddr); err != nil {
				utils.Logger.Fatal("%v", err)
			}
		}()
	}

	if db.replicationClient != nil {
		go db.replicationClient.Start(context.Background())
	}
}

func (db *DataStore) Set(key, value string) {
	db.data[key] = value
	db.OffsetManager.OffsetQueue <- OperationStructure{
		Operation: proto.Command_SET,
		Key:       key,
		Value:     value,
	}
}

func (db *DataStore) Get(key string) (string, bool) {
	value, ok := db.data[key]
	return value, ok
}

func (db *DataStore) Delete(key string) bool {
	_, exists := db.data[key]
	if exists {
		delete(db.data, key)
		db.OffsetManager.OffsetQueue <- OperationStructure{
			Operation: proto.Command_DEL,
			Key:       key,
			Value:     "",
		}
	}
	return exists
}

func (db *DataStore) Size() int {
	return len(db.data)
}
