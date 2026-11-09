package main

import (
	"strings"
	"testing"
	"time"
)

func TestCreateAPIKey_Valid(t *testing.T) {
	fullKey, hash, err := CreateAPIKey("api_key", []string{"read"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(fullKey, "api_key_") {
		t.Errorf("expected api_key_ prefix, got %s", fullKey[:8])
	}
	if len(fullKey) < 40 {
		t.Errorf("expected key length >= 40, got %d", len(fullKey))
	}
	if hash == "" {
		t.Error("expected non-empty hash")
	}
}

func TestHashKey(t *testing.T) {
	key := "api_key_test-key-123"
	hash := HashKey(key)
	if len(hash) != 64 {
		t.Errorf("expected 64-char hex hash, got %d", len(hash))
	}
	// Same key should produce same hash.
	if HashKey(key) != hash {
		t.Error("expected deterministic hashing")
	}
}

func TestValidateKey_Valid(t *testing.T) {
	store := NewKeyStore()
	fullKey, record, _ := GenerateAPIKey("test", []string{"read", "write"}, 1*time.Hour)
	store.Add(record)

	_, err := ValidateKey(store, fullKey, "read")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateKey_InvalidKey(t *testing.T) {
	store := NewKeyStore()
	_, err := ValidateKey(store, "invalid-key", "read")
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
}

func TestValidateKey_WrongScope(t *testing.T) {
	store := NewKeyStore()
	fullKey, record, _ := GenerateAPIKey("test", []string{"read"}, 1*time.Hour)
	store.Add(record)

	_, err := ValidateKey(store, fullKey, "write")
	if err == nil {
		t.Fatal("expected error for wrong scope")
	}
	if !strings.Contains(err.Error(), "insufficient scope") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateKey_Expired(t *testing.T) {
	store := NewKeyStore()
	fullKey, record, _ := GenerateAPIKey("test", []string{"read"}, -1*time.Hour)
	store.Add(record)

	_, err := ValidateKey(store, fullKey, "read")
	if err == nil {
		t.Fatal("expected error for expired key")
	}
}

func TestHasScope(t *testing.T) {
	tests := []struct {
		scopes   []string
		required string
		want     bool
	}{
		{[]string{"read"}, "read", true},
		{[]string{"read"}, "write", false},
		{[]string{"admin"}, "read", true},
		{[]string{"doc:*"}, "doc:read", true},
		{[]string{"doc:*"}, "user:read", false},
	}
	for _, tt := range tests {
		got := hasScope(tt.scopes, tt.required)
		if got != tt.want {
			t.Errorf("hasScope(%v, %s) = %v, want %v", tt.scopes, tt.required, got, tt.want)
		}
	}
}
