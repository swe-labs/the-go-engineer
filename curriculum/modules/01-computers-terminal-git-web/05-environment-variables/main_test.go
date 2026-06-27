package main

import (
	"os"
	"testing"
)

func TestGetEnvSet(t *testing.T) {
	os.Setenv("TEST_GO_ENGINEER_KEY", "test_value")

	val, ok := GetEnv("TEST_GO_ENGINEER_KEY")
	if !ok {
		t.Error("expected ok=true")
	}
	if val != "test_value" {
		t.Errorf("got %q, want %q", val, "test_value")
	}

	os.Unsetenv("TEST_GO_ENGINEER_KEY")
}

func TestGetEnvNotSet(t *testing.T) {
	os.Unsetenv("TEST_GO_ENGINEER_NONEXISTENT")

	val, ok := GetEnv("TEST_GO_ENGINEER_NONEXISTENT")
	if ok {
		t.Error("expected ok=false for unset variable")
	}
	if val != "" {
		t.Errorf("expected empty value, got %q", val)
	}
}

func TestGetEnvEmptyValue(t *testing.T) {
	os.Setenv("TEST_GO_ENGINEER_EMPTY", "")

	val, ok := GetEnv("TEST_GO_ENGINEER_EMPTY")
	if !ok {
		t.Error("expected ok=true for empty but set variable")
	}
	if val != "" {
		t.Errorf("expected empty value, got %q", val)
	}

	os.Unsetenv("TEST_GO_ENGINEER_EMPTY")
}

func TestGetEnvPathExists(t *testing.T) {
	val, ok := GetEnv("PATH")
	if !ok {
		t.Fatal("expected PATH to be set")
	}
	if val == "" {
		t.Error("expected PATH to have a non-empty value")
	}
}
