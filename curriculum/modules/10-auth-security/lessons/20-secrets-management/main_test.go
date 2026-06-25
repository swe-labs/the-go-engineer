package main

import (
	"crypto/rand"
	"encoding/json"
	"os"
	"testing"
)

func TestSecretManagerSetGet(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	sm := NewSecretManager(key)

	tests := []struct {
		name  string
		set   string
		value string
	}{
		{name: "database password", set: "DB_PASSWORD", value: "s3cret"},
		{name: "api key", set: "API_KEY", value: "sk-abc123"},
		{name: "empty value", set: "EMPTY", value: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sm.Set(tc.set, tc.value)
			got, ok := sm.Get(tc.set)
			if !ok {
				t.Errorf("Get(%q) returned ok=false", tc.set)
			}
			if got != tc.value {
				t.Errorf("Get(%q) = %q, want %q", tc.set, got, tc.value)
			}
		})
	}
}

func TestSecretManagerGetMissing(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	sm := NewSecretManager(key)

	_, ok := sm.Get("NONEXISTENT")
	if ok {
		t.Error("expected ok=false for missing secret")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	tests := []struct {
		name string
		data []byte
	}{
		{name: "simple string", data: []byte("hello world")},
		{name: "json data", data: []byte(`{"key":"value"}`)},
		{name: "empty", data: []byte{}},
		{name: "binary", data: []byte{0x00, 0x01, 0xFF}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			enc, err := encrypt(tc.data, key)
			if err != nil {
				t.Fatalf("encrypt error: %v", err)
			}
			dec, err := decrypt(enc, key)
			if err != nil {
				t.Fatalf("decrypt error: %v", err)
			}
			if string(dec) != string(tc.data) {
				t.Errorf("decrypt = %q, want %q", dec, tc.data)
			}
		})
	}
}

func TestEncryptWrongKey(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	rand.Read(key1)
	rand.Read(key2)

	data := []byte("secret data")
	enc, err := encrypt(data, key1)
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}

	_, err = decrypt(enc, key2)
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("TEST_DB_URL", "postgres://db:5432")
	os.Setenv("TEST_API_KEY", "sk-test")
	os.Setenv("OTHER_VAR", "should-not-appear")

	secrets := loadFromEnv("TEST_")

	if secrets["TEST_DB_URL"] != "postgres://db:5432" {
		t.Errorf("TEST_DB_URL = %q, want %q", secrets["TEST_DB_URL"], "postgres://db:5432")
	}
	if secrets["TEST_API_KEY"] != "sk-test" {
		t.Errorf("TEST_API_KEY = %q, want %q", secrets["TEST_API_KEY"], "sk-test")
	}
	if _, ok := secrets["OTHER_VAR"]; ok {
		t.Error("OTHER_VAR should not be loaded")
	}
}

func TestEncryptedJSONRoundtrip(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	original := map[string]string{
		"DB_PASSWORD": "s3cret",
		"API_KEY":     "sk-abc",
	}
	data, _ := json.Marshal(original)

	enc, err := encrypt(data, key)
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}

	dec, err := decrypt(enc, key)
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}

	var restored map[string]string
	json.Unmarshal(dec, &restored)

	if restored["DB_PASSWORD"] != "s3cret" {
		t.Errorf("DB_PASSWORD = %q, want %q", restored["DB_PASSWORD"], "s3cret")
	}
	if restored["API_KEY"] != "sk-abc" {
		t.Errorf("API_KEY = %q, want %q", restored["API_KEY"], "sk-abc")
	}
}
