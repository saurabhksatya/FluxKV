package internal

type DataStore struct {
	data          map[string]string
	offsetManager *OffsetManager
}

func NewDataStore() *DataStore {
	return &DataStore{
		data:          make(map[string]string),
		offsetManager: NewOffsetManager(),
	}
}

func (db *DataStore) Set(key, value string) {
	db.data[key] = value
	db.offsetManager.OffsetQueue <- OperationStructure{
		Operation: "SET",
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
		db.offsetManager.OffsetQueue <- OperationStructure{
			Operation: "DEL",
			Key:       key,
			Value:     "",
		}
	}
	return exists
}

func (db *DataStore) Size() int {
	return len(db.data)
}
