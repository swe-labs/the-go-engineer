package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMkdir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "newdir")
	err := os.Mkdir(path, 0755)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Errorf("expected %q to be a directory", path)
	}
}

func TestMkdirAll(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "a", "b", "c")
	err := os.MkdirAll(path, 0755)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Errorf("expected %q to be a directory", path)
	}
}

func TestMkdirAllExisting(t *testing.T) {
	path := t.TempDir()
	err := os.MkdirAll(path, 0755)
	if err != nil {
		t.Errorf("MkdirAll on existing dir should succeed: %v", err)
	}
}

func TestReadDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644)
	os.Mkdir(filepath.Join(dir, "sub"), 0755)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestRemove(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "empty")
	os.Mkdir(sub, 0755)

	err := os.Remove(sub)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sub); !os.IsNotExist(err) {
		t.Errorf("expected dir to be removed")
	}
}

func TestRemoveNonEmpty(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0755)
	os.WriteFile(filepath.Join(sub, "file.txt"), []byte("x"), 0644)

	err := os.Remove(sub)
	if err == nil {
		t.Error("expected error when removing non-empty directory")
	}
}

func TestRemoveAll(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0755)
	os.WriteFile(filepath.Join(sub, "file.txt"), []byte("x"), 0644)

	err := os.RemoveAll(sub)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sub); !os.IsNotExist(err) {
		t.Errorf("expected dir tree to be removed")
	}
}

func TestDirEntryInfo(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.txt"), []byte("content"), 0644)

	entries, _ := os.ReadDir(dir)
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}

	e := entries[0]
	info, err := e.Info()
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() != 7 {
		t.Errorf("expected size 7, got %d", info.Size())
	}
	if e.Name() != "test.txt" {
		t.Errorf("expected name 'test.txt', got %q", e.Name())
	}
}

func TestMkdirAlreadyExists(t *testing.T) {
	dir := t.TempDir()
	err := os.Mkdir(dir, 0755)
	if err == nil {
		t.Error("expected error for Mkdir on existing dir")
	}
	if !os.IsExist(err) {
		t.Errorf("expected IsExist error, got %v", err)
	}
}
