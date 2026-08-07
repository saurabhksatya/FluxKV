package internal

import "testing"

func TestDataStore_SetGet(t *testing.T) {
	db := NewDataStore()

	db.Set("key1", "value1")
	val, ok := db.Get("key1")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if val != "value1" {
		t.Errorf("expected 'value1', got '%s'", val)
	}
}

func TestDataStore_GetMissing(t *testing.T) {
	db := NewDataStore()
	_, ok := db.Get("missing")
	if ok {
		t.Fatal("expected key to not exist")
	}
}

func TestDataStore_Overwrite(t *testing.T) {
	db := NewDataStore()
	db.Set("key", "v1")
	db.Set("key", "v2")
	val, ok := db.Get("key")
	if !ok || val != "v2" {
		t.Errorf("expected 'v2', got '%s' (ok=%v)", val, ok)
	}
}

func TestDataStore_Delete(t *testing.T) {
	db := NewDataStore()
	db.Set("key", "val")
	db.Delete("key")
	_, ok := db.Get("key")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestDataStore_DeleteMissing(t *testing.T) {
	db := NewDataStore()
	db.Delete("missing") // should not panic
}