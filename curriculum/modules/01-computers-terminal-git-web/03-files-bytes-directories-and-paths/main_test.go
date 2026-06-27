package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateAndReadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	content := "hello world"

	info, err := CreateFile(path, content)
	if err != nil {
		t.Fatalf("CreateFile(%q, %q) error: %v", path, content, err)
	}

	if info.Content != content {
		t.Errorf("got content %q, want %q", info.Content, content)
	}
	if info.Size != int64(len(content)) {
		t.Errorf("got size %d, want %d", info.Size, len(content))
	}
	if info.IsDir {
		t.Errorf("expected IsDir=false, got true")
	}
	if info.Path != path {
		t.Errorf("got path %q, want %q", info.Path, path)
	}
}

func TestCreateFileCreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c", "nested.txt")
	content := "nested"

	info, err := CreateFile(path, content)
	if err != nil {
		t.Fatalf("CreateFile error: %v", err)
	}
	if info.Content != content {
		t.Errorf("got %q, want %q", info.Content, content)
	}
}

func TestReadFileNotFound(t *testing.T) {
	_, err := ReadFile(filepath.Join(t.TempDir(), "nonexistent.txt"))
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}

func TestReadFileEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")
	if err := os.WriteFile(path, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	info, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if info.Content != "" {
		t.Errorf("got %q, want empty", info.Content)
	}
	if info.Size != 0 {
		t.Errorf("got size %d, want 0", info.Size)
	}
}
