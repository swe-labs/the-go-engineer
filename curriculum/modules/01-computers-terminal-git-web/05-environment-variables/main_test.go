package main

import (
	"os"
	"testing"
)

func TestGetEnvSet(t *testing.T) {
	os.Setenv("TEST_GO_ENGINEER_KEY", "test_value")
	defer os.Unsetenv("TEST_GO_ENGINEER_KEY")

	ev := GetEnv("TEST_GO_ENGINEER_KEY")
	if !ev.Set {
		t.Error("expected Set=true")
	}
	if ev.Value != "test_value" {
		t.Errorf("got %q, want %q", ev.Value, "test_value")
	}
	if ev.Key != "TEST_GO_ENGINEER_KEY" {
		t.Errorf("got key %q, want %q", ev.Key, "TEST_GO_ENGINEER_KEY")
	}
}

func TestGetEnvNotSet(t *testing.T) {
	os.Unsetenv("TEST_GO_ENGINEER_NONEXISTENT")

	ev := GetEnv("TEST_GO_ENGINEER_NONEXISTENT")
	if ev.Set {
		t.Error("expected Set=false for unset variable")
	}
	if ev.Value != "" {
		t.Errorf("expected empty value, got %q", ev.Value)
	}
}

func TestGetEnvEmptyValue(t *testing.T) {
	os.Setenv("TEST_GO_ENGINEER_EMPTY", "")
	defer os.Unsetenv("TEST_GO_ENGINEER_EMPTY")

	ev := GetEnv("TEST_GO_ENGINEER_EMPTY")
	if !ev.Set {
		t.Error("expected Set=true for empty but set variable")
	}
	if ev.Value != "" {
		t.Errorf("expected empty value, got %q", ev.Value)
	}
}

func TestGetEnvPathExists(t *testing.T) {
	ev := GetEnv("PATH")
	if !ev.Set {
		t.Fatal("expected PATH to be set")
	}
	if ev.Value == "" {
		t.Error("expected PATH to have a non-empty value")
	}
}

func TestGetEnvKeyPreserved(t *testing.T) {
	key := "TEST_CASE_KEY_UPPER"
	os.Setenv(key, "val")
	defer os.Unsetenv(key)

	ev := GetEnv(key)
	if ev.Key != key {
		t.Errorf("got key %q, want %q", ev.Key, key)
	}
}
