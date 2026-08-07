package internal

type DataStore struct {
	data map[string]string
}

func NewDataStore() *DataStore {
	return &DataStore{
		data: make(map[string]string),
	}
}

func (db *DataStore) Set(key, value string) {
	db.data[key] = value
}

func (db *DataStore) Get(key string) (string, bool) {
	value, ok := db.data[key]
	return value, ok
}

func (db *DataStore) Delete(key string) {
	delete(db.data, key)
}
