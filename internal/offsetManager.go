package internal

import (
	"container/list"
	"fluxKV/utils"
	"sync"
)

const MAX_OFFSET_LEN = 100

type OperationStructure struct {
	Operation string
	Key       string
	Value     string
}

type OffsetManager struct {
	Mutex       sync.RWMutex
	Offset      uint64
	Offsets     *list.List
	OffsetQueue chan OperationStructure
}

func NewOffsetManager() *OffsetManager {
	return &OffsetManager{
		Mutex:       sync.RWMutex{},
		Offset:      0,
		OffsetQueue: make(chan OperationStructure, MAX_OFFSET_LEN),
		Offsets:     list.New(),
	}
}

func (manager *OffsetManager) AddOffset(operation OperationStructure) {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	if operation.Operation == "ADD" || operation.Operation == "DEL" {
		if manager.Offsets.Len() > MAX_OFFSET_LEN {
			manager.Offsets.Remove(manager.Offsets.Front())
		}
		manager.Offset++
	} else {
		utils.Logger.Fatal("%s", "Invalid operation type.")
	}
}

func (om *OffsetManager) Manage() {
	for task := range om.OffsetQueue {
		om.AddOffset(task)
	}
}

//
//func (db *DataStore) Change(op string, key string, value string) {
//	switch op {
//	case "SET":
//		newOffset := map[string]string{
//			"op":    op,
//			"key":   key,
//			"value": value,
//		}
//		db.offsets.PushBack(newOffset)
//		if db.offsets.Len() > OFFSET_MAX {
//			db.offsets.Remove(db.offsets.Front())
//		}
//		db.offset++
//	case "DEL":
//		newOffset := map[string]string{
//			"op":    op,
//			"key":   key,
//			"value": value,
//		}
//		db.offsets.PushBack(newOffset)
//		if db.offsets.Len() > OFFSET_MAX {
//			db.offsets.Remove(db.offsets.Front())
//		}
//		db.offset++
//	}
//}
