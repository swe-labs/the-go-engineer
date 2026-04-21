package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTE10Golden(t *testing.T) {
	dir := t.TempDir()
	got := []byte("golden output\n")
	golden := filepath.Join(dir, "out.golden")
	if err := os.WriteFile(golden, got, 0o644); err != nil {
		t.Fatalf("write golden: %v", err)
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if string(got) != string(want) {
		t.Fatal("golden mismatch")
	}
}
