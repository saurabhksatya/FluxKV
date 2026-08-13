package internal

import (
	"container/list"
	"fluxKV/replication/proto"
	"fluxKV/utils"
	"sync"
)

const MAX_OFFSET_LEN = 100

type OperationStructure struct {
	Operation proto.Command
	Key       string
	Value     string
	Offset    uint64
}

type OffsetManager struct {
	Mutex       sync.RWMutex
	Offset      uint64
	Offsets     *list.List
	OffsetQueue chan OperationStructure
	Enabled     bool
}

func NewOffsetManager(enabled bool) *OffsetManager {
	return &OffsetManager{
		Mutex:       sync.RWMutex{},
		Offset:      0,
		OffsetQueue: make(chan OperationStructure, MAX_OFFSET_LEN),
		Offsets:     list.New(),
		Enabled:     enabled,
	}
}

func (manager *OffsetManager) AddOffset(operation OperationStructure) {
	if !manager.Enabled {
		return
	}
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	switch operation.Operation {
	case proto.Command_SET, proto.Command_DEL:
		if manager.Offsets.Len() > MAX_OFFSET_LEN {
			manager.Offsets.Remove(manager.Offsets.Front())
		}
		operation.Offset = manager.Offset + 1
		manager.Offsets.PushBack(operation)
		manager.Offset++
	default:
		utils.Logger.Fatal("%s", "Invalid operation type.")
	}
	return
}
