package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAndReadFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.txt")
	data := []byte("hello file")

	err := os.WriteFile(path, data, 0644)
	if err != nil {
		t.Fatal(err)
	}

	read, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(read) != "hello file" {
		t.Errorf("expected 'hello file', got %q", read)
	}
}

func TestOpenAndStat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stat.txt")
	os.WriteFile(path, []byte("stats"), 0644)

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if info.Size() != 5 {
		t.Errorf("expected size 5, got %d", info.Size())
	}
	if info.Name() != "stat.txt" {
		t.Errorf("expected name 'stat.txt', got %q", info.Name())
	}
}

func TestFileNotExists(t *testing.T) {
	_, err := os.ReadFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestCreateFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "created.txt")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.Write([]byte("created"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() != 7 {
		t.Errorf("expected size 7, got %d", info.Size())
	}
}

func TestFilePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "perm.txt")
	os.WriteFile(path, []byte("perm"), 0600)

	info, _ := os.Stat(path)
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Logf("permissions: %o (note: umask may affect)", perm)
	}
}

func TestWriteFileOverwrites(t *testing.T) {
	path := filepath.Join(t.TempDir(), "overwrite.txt")
	os.WriteFile(path, []byte("original"), 0644)
	os.WriteFile(path, []byte("new"), 0644)

	data, _ := os.ReadFile(path)
	if string(data) != "new" {
		t.Errorf("expected 'new', got %q", data)
	}
}

func TestAppendToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "append.txt")
	os.WriteFile(path, []byte("hello "), 0644)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("world"))
	f.Close()

	data, _ := os.ReadFile(path)
	if string(data) != "hello world" {
		t.Errorf("expected 'hello world', got %q", data)
	}
}
