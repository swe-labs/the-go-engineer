package main

import (
	"os"
	"testing"
)

func TestWriteReadConfig(t *testing.T) {
	dir := t.TempDir()

	path, err := WriteConfig(dir, "hello world")
	if err != nil {
		t.Fatalf("WriteConfig failed: %v", err)
	}

	content, err := ReadConfig(path)
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}
	if content != "hello world" {
		t.Errorf("content = %q; want %q", content, "hello world")
	}
}

func TestCleanupOrder(t *testing.T) {
	var order []string

	t.Cleanup(func() {
		order = append(order, "A")
		if len(order) != 3 {
			t.Errorf("A: expected len 3, got %d (order=%v)", len(order), order)
		}
	})
	t.Cleanup(func() {
		order = append(order, "B")
		if len(order) != 2 {
			t.Errorf("B: expected len 2, got %d (order=%v)", len(order), order)
		}
	})
	t.Cleanup(func() {
		order = append(order, "C")
		if len(order) != 1 {
			t.Errorf("C: expected len 1, got %d (order=%v)", len(order), order)
		}
	})
}

func TestCleanupEnvVar(t *testing.T) {
	origVal, origExists := os.LookupEnv("MY_VAR")

	cleanup := SetupEnv("MY_VAR", "test-value")
	got := os.Getenv("MY_VAR")
	if got != "test-value" {
		t.Fatalf("MY_VAR = %q; want %q", got, "test-value")
	}

	t.Cleanup(func() {
		val, exists := os.LookupEnv("MY_VAR")
		if origExists {
			if val != origVal {
				t.Errorf("after cleanup MY_VAR = %q; want original %q", val, origVal)
			}
		} else {
			if exists {
				t.Errorf("after cleanup MY_VAR still set to %q; should be unset", val)
			}
		}
	})
	t.Cleanup(cleanup)
}
