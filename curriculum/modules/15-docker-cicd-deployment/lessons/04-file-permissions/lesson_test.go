package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWriteSecureFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secret.key")

	err := writeSecureFile(path, []byte("my-secret"))
	if err != nil {
		t.Fatalf("writeSecureFile failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}

	mode := info.Mode().Perm()
	if runtime.GOOS != "windows" && mode != 0600 {
		t.Errorf("got mode %#o, want %#o", mode, 0600)
	}
}

func TestReadPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	os.WriteFile(path, []byte("hello"), 0644)

	mode, err := readPermissions(path)
	if err != nil {
		t.Fatalf("readPermissions failed: %v", err)
	}

	if runtime.GOOS != "windows" && mode != 0644 {
		t.Errorf("got mode %#o, want %#o", mode, 0644)
	}
}

func TestReadPermissionsNonexistent(t *testing.T) {
	_, err := readPermissions("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}
