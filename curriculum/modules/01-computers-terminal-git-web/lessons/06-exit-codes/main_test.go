package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckFileStatus(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Existing file should return code 0
	filePath := filepath.Join(tempDir, "sample.txt")
	if err := os.WriteFile(filePath, []byte("data"), 0644); err != nil {
		t.Fatalf("failed to create sample file: %v", err)
	}

	code := CheckFileStatus(filePath)
	if code != 0 {
		t.Errorf("expected exit code 0 for existing file, got %d", code)
	}

	// 2. Missing file should return code 1
	missingPath := filepath.Join(tempDir, "does-not-exist.txt")
	code = CheckFileStatus(missingPath)
	if code != 1 {
		t.Errorf("expected exit code 1 for missing file, got %d", code)
	}

	// 3. Directory path should return code 1
	subDir := filepath.Join(tempDir, "sub-directory")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	code = CheckFileStatus(subDir)
	if code != 1 {
		t.Errorf("expected exit code 1 for directory, got %d", code)
	}
}
