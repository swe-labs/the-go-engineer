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

	p, size, isDir, c, err := CreateFile(path, content)
	if err != nil {
		t.Fatalf("CreateFile(%q, %q) error: %v", path, content, err)
	}

	if c != content {
		t.Errorf("got content %q, want %q", c, content)
	}
	if size != int64(len(content)) {
		t.Errorf("got size %d, want %d", size, len(content))
	}
	if isDir {
		t.Errorf("expected IsDir=false, got true")
	}
	if p != path {
		t.Errorf("got path %q, want %q", p, path)
	}
}

func TestCreateFileCreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c", "nested.txt")
	content := "nested"

	_, _, _, c, err := CreateFile(path, content)
	if err != nil {
		t.Fatalf("CreateFile error: %v", err)
	}
	if c != content {
		t.Errorf("got %q, want %q", c, content)
	}
}

func TestReadFileNotFound(t *testing.T) {
	_, _, _, _, err := ReadFile(filepath.Join(t.TempDir(), "nonexistent.txt"))
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
	_, size, _, c, err := ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if c != "" {
		t.Errorf("got %q, want empty", c)
	}
	if size != 0 {
		t.Errorf("got size %d, want 0", size)
	}
}
