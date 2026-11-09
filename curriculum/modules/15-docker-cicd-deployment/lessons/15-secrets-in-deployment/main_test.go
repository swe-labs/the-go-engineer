package main

import (
	"os"
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	store := NewInMemoryStore()
	err := store.Set("KEY", "value")
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}
	val, err := store.Get("KEY")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if val != "value" {
		t.Errorf("expected value, got %s", val)
	}
}

func TestInMemoryStoreMissing(t *testing.T) {
	store := NewInMemoryStore()
	_, err := store.Get("NONEXISTENT")
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestEnvSecretStore(t *testing.T) {
	os.Setenv("TEST_SECRET", "testvalue")
	store := &EnvSecretStore{}
	val, err := store.Get("TEST_SECRET")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if val != "testvalue" {
		t.Errorf("expected testvalue, got %s", val)
	}
}

func TestEnvSecretStoreMissing(t *testing.T) {
	store := &EnvSecretStore{}
	_, err := store.Get("NONEXISTENT_SECRET_KEY")
	if err == nil {
		t.Error("expected error for missing env var")
	}
}

func TestSecretsManagerGetSet(t *testing.T) {
	sm := NewSecretsManager(NewInMemoryStore())
	sm.Set("KEY", "val")
	v, _ := sm.Get("KEY")
	if v != "val" {
		t.Errorf("expected val, got %s", v)
	}
}

func TestResolveTemplate(t *testing.T) {
	sm := NewSecretsManager(NewInMemoryStore())
	sm.Set("PASSWORD", "hunter2")
	result, err := sm.ResolveTemplate("postgres://user:${PASSWORD}@localhost/db")
	if err != nil {
		t.Fatalf("resolve failed: %v", err)
	}
	if result != "postgres://user:hunter2@localhost/db" {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestResolveTemplateMissing(t *testing.T) {
	sm := NewSecretsManager(NewInMemoryStore())
	_, err := sm.ResolveTemplate("prefix${MISSING}suffix")
	if err == nil {
		t.Error("expected error for missing secret")
	}
}

func TestParseDotEnv(t *testing.T) {
	content := "KEY1=value1\nKEY2=value2\n"
	entries, err := ParseDotEnv(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "KEY1" || entries[0].Value != "value1" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
}

func TestParseDotEnvWithComments(t *testing.T) {
	content := "# This is a comment\nKEY=val\n"
	entries, err := ParseDotEnv(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestParseDotEnvQuoted(t *testing.T) {
	content := `KEY="quoted value"` + "\n"
	entries, err := ParseDotEnv(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "quoted value" {
		t.Errorf("expected 'quoted value', got %s", entries[0].Value)
	}
}

func TestParseDotEnvInvalid(t *testing.T) {
	_, err := ParseDotEnv("INVALID_LINE")
	if err == nil {
		t.Error("expected error for invalid line")
	}
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"abc", "***"},
		{"abcdefgh", "********"},
		{"abcdefghij", "ab******ij"},
		{"sk-abc123def456ghi789", "sk*****************89"},
	}
	for _, tc := range tests {
		got := Sanitize(tc.input)
		if got != tc.want {
			t.Errorf("Sanitize(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestSecretValidator(t *testing.T) {
	v := &SecretValidator{}
	if err := v.ValidateNotEmpty("KEY", ""); err == nil {
		t.Error("expected error for empty secret")
	}
	if err := v.ValidateNotEmpty("KEY", "val"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := v.ValidateLength("KEY", "short", 10); err == nil {
		t.Error("expected error for short secret")
	}
	if err := v.ValidateLength("KEY", "longenoughsecret", 10); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
