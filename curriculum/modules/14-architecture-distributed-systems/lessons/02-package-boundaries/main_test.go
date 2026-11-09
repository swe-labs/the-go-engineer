package main

import (
	"testing"
)

func TestInMemoryStore_SaveAndFind(t *testing.T) {
	store := newInMemoryStore()
	err := store.Save("test-key")
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	val, err := store.Find("test-key")
	if err != nil {
		t.Fatalf("Find failed: %v", err)
	}
	if val != "test-key" {
		t.Errorf("expected test-key, got %s", val)
	}
}

func TestInMemoryStore_NotFound(t *testing.T) {
	store := newInMemoryStore()
	_, err := store.Find("missing")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestExportedService_UsesRepository(t *testing.T) {
	store := newInMemoryStore()
	svc := NewExportedService(store)
	if svc.repo == nil {
		t.Fatal("repository should not be nil")
	}
}

func TestRepositoryInterface_Satisfied(t *testing.T) {
	var _ Repository = (*inMemoryStore)(nil)
}
