// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"testing"
	"testing/fstest"
)

// ============================================================================
// Section 10 Supplement: fs.FS Testing — Test File
// ============================================================================
//
// These tests use fstest.MapFS — an in-memory filesystem.
// Zero disk I/O. No temp directories. No cleanup. Runs in microseconds.
//
// This is the definitive demonstration of why accepting fs.FS is superior
// to accepting a string path.
//
// RUN: go test -v ./10-filesystem/8-fs-testing-seam
// ============================================================================

func TestSearchLogs_FindsMatches(t *testing.T) {
	// fstest.MapFS is an in-memory filesystem — no real files needed.
	fakeFS := fstest.MapFS{
		"app.log":    {Data: []byte("INFO: server started\nERROR: connection failed\nINFO: retry ok\n")},
		"access.log": {Data: []byte("GET /api/v1 200\nPOST /api/v1 500 error\n")},
		"notes.txt":  {Data: []byte("TODO: fix the error handling\n")},
		"main.go":    {Data: []byte("// error handling code")}, // Should be skipped (.go)
	}

	results, err := SearchLogs(fakeFS, "error")
	if err != nil {
		t.Fatalf("SearchLogs returned error: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("expected 3 matches, got %d: %+v", len(results), results)
	}
}

func TestSearchLogs_CaseInsensitive(t *testing.T) {
	fakeFS := fstest.MapFS{
		"app.log": {Data: []byte("ERROR uppercase\nerror lowercase\nError mixed\n")},
	}

	results, err := SearchLogs(fakeFS, "ERROR")
	if err != nil {
		t.Fatalf("SearchLogs error: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("expected 3 case-insensitive matches, got %d", len(results))
	}
}

func TestSearchLogs_EmptyFS(t *testing.T) {
	results, err := SearchLogs(fstest.MapFS{}, "error")
	if err != nil {
		t.Fatalf("empty FS should not error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty FS, got %d", len(results))
	}
}

func TestSearchLogs_LineNumbers(t *testing.T) {
	fakeFS := fstest.MapFS{
		"app.log": {Data: []byte("line one\nERROR: on line two\nline three\nERROR: on line four\n")},
	}

	results, err := SearchLogs(fakeFS, "error")
	if err != nil {
		t.Fatalf("SearchLogs error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].LineNumber != 2 {
		t.Errorf("first match line = %d, want 2", results[0].LineNumber)
	}
	if results[1].LineNumber != 4 {
		t.Errorf("second match line = %d, want 4", results[1].LineNumber)
	}
}

func TestLoadConfig_ParsesKeyValues(t *testing.T) {
	fakeFS := fstest.MapFS{
		"config.env": {Data: []byte("# comment\nHOST = localhost\nPORT = 8080\nDEBUG = true\n")},
	}

	config, err := LoadConfig(fakeFS, "config.env")
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	tests := map[string]string{
		"HOST":  "localhost",
		"PORT":  "8080",
		"DEBUG": "true",
	}
	for key, want := range tests {
		if got := config[key]; got != want {
			t.Errorf("config[%q] = %q, want %q", key, got, want)
		}
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig(fstest.MapFS{}, "nonexistent.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
