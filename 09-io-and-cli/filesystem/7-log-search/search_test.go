// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"os"
	"path/filepath"
	"testing"
)

// ============================================================================
// Tests for: Log Search Tool (Filesystem Exercise)
// ============================================================================
//
// These tests verify searchFile and searchDirectory work correctly
// with real temporary files on disk.
//
// RUN: go test -v ./10-filesystem/7-log-search
// ============================================================================

func TestSearchFileFindsMatches(t *testing.T) {
	// Create a temp file with known content
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	content := "INFO: Server started\nERROR: Connection failed\nINFO: Retry successful\nERROR: Timeout\n"
	os.WriteFile(logPath, []byte(content), 0644)

	results, err := searchFile(logPath, "error")
	if err != nil {
		t.Fatalf("searchFile returned error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(results))
	}

	// Verify line numbers are correct (1-indexed)
	if results[0].LineNumber != 2 {
		t.Errorf("first match line = %d, want 2", results[0].LineNumber)
	}
	if results[1].LineNumber != 4 {
		t.Errorf("second match line = %d, want 4", results[1].LineNumber)
	}
}

func TestSearchFileIsCaseInsensitive(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	content := "error lowercase\nERROR uppercase\nError mixed\n"
	os.WriteFile(logPath, []byte(content), 0644)

	results, err := searchFile(logPath, "ERROR")
	if err != nil {
		t.Fatalf("searchFile returned error: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("expected 3 case-insensitive matches, got %d", len(results))
	}
}

func TestSearchFileNoMatches(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	os.WriteFile(logPath, []byte("nothing relevant here\n"), 0644)

	results, err := searchFile(logPath, "FATAL")
	if err != nil {
		t.Fatalf("searchFile returned error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 matches, got %d", len(results))
	}
}

func TestSearchDirectoryFiltersExtensions(t *testing.T) {
	tmpDir := t.TempDir()

	// Write files with different extensions
	os.WriteFile(filepath.Join(tmpDir, "app.log"), []byte("ERROR: test\n"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "notes.txt"), []byte("error in code\n"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("// error handling\n"), 0644) // Should be skipped

	results, err := searchDirectory(tmpDir, "error")
	if err != nil {
		t.Fatalf("searchDirectory returned error: %v", err)
	}

	// Only .log and .txt files should be searched, not .go
	if len(results) != 2 {
		t.Errorf("expected 2 matches (log + txt only), got %d", len(results))
	}
}
