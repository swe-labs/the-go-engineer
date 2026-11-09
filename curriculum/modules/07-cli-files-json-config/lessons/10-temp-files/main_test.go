package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateTempFile(t *testing.T) {
	f, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	if f.Name() == "" {
		t.Error("expected non-empty temp file name")
	}
	if !strings.HasPrefix(filepath.Base(f.Name()), "test-") {
		t.Errorf("expected prefix 'test-', got %q", f.Name())
	}
	if !strings.HasSuffix(f.Name(), ".txt") {
		t.Errorf("expected suffix '.txt', got %q", f.Name())
	}
}

func TestCreateTempFileDefaultDir(t *testing.T) {
	f, err := os.CreateTemp("", "nodir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	// Should be in system temp dir
	tmpDir := os.TempDir()
	if !strings.HasPrefix(f.Name(), tmpDir) {
		t.Logf("temp file %q is not in system temp dir %q (may vary by OS)", f.Name(), tmpDir)
	}
}

func TestCreateTempFileWritable(t *testing.T) {
	f, err := os.CreateTemp("", "write-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	_, err = f.Write([]byte("hello temp"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateTempFileRemoved(t *testing.T) {
	f, err := os.CreateTemp("", "remove-*")
	if err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	f.Close()
	os.Remove(name)

	if _, err := os.Stat(name); !os.IsNotExist(err) {
		t.Errorf("expected file to be removed")
	}
}

func TestMkdirTemp(t *testing.T) {
	dir, err := os.MkdirTemp("", "testdir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Errorf("expected %q to be a directory", dir)
	}
	if !strings.HasPrefix(filepath.Base(dir), "testdir-") {
		t.Errorf("expected prefix 'testdir-', got %q", filepath.Base(dir))
	}
}

func TestMkdirTempCustomDir(t *testing.T) {
	parent := t.TempDir()
	dir, err := os.MkdirTemp(parent, "custom-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if !strings.HasPrefix(dir, parent) {
		t.Errorf("expected temp dir %q to be under parent %q", dir, parent)
	}
}

func TestTempFilePattern(t *testing.T) {
	tests := []struct {
		pattern string
		valid   bool
	}{
		{"*.txt", true},
		{"prefix-*", true},
		{"no-star", true},
		{"", true},
	}

	for _, tt := range tests {
		f, err := os.CreateTemp("", tt.pattern)
		if err != nil && tt.valid {
			t.Errorf("CreateTemp with pattern %q failed: %v", tt.pattern, err)
		}
		if err == nil {
			os.Remove(f.Name())
			f.Close()
		}
	}
}

func TestTempDirWithFiles(t *testing.T) {
	dir, err := os.MkdirTemp("", "withfiles-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for i := 0; i < 3; i++ {
		f, err := os.CreateTemp(dir, "sub-*")
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) != 3 {
		t.Errorf("expected 3 entries in temp dir, got %d", len(entries))
	}
}
