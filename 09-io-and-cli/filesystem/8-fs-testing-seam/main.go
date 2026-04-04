// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// ============================================================================
// Section 10 Supplement: io/fs as a Testing Seam
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why accepting fs.FS instead of string paths makes code testable
//   - os.DirFS: wrapping a real directory as an fs.FS
//   - Using fstest.MapFS as an in-memory filesystem in tests (zero disk I/O)
//   - The "accept interfaces, return structs" principle at the filesystem level
//
// THE PROBLEM:
//   Functions that hardcode os.Open() or os.ReadDir() cannot be tested
//   without real files on disk. This means:
//     - Tests need a temp directory (slow, fragile, platform-dependent)
//     - Tests must clean up after themselves
//     - Parallel tests may conflict on the same temp path
//
// THE SOLUTION:
//   Accept fs.FS. In production, pass os.DirFS("/path/to/logs").
//   In tests, pass fstest.MapFS — an in-memory filesystem that lives in RAM.
//
//   PRODUCTION:
//     results, err := SearchLogs(os.DirFS("/var/log/app"), "ERROR")
//
//   TEST (zero disk I/O):
//     fakeFS := fstest.MapFS{
//         "app.log": &fstest.MapFile{Data: []byte("ERROR: db failed\n")},
//     }
//     results, err := SearchLogs(fakeFS, "ERROR")
//
// ENGINEERING DEPTH:
//   fs.FS (added in Go 1.16) is a minimal interface:
//     type FS interface { Open(name string) (File, error) }
//   The standard library provides implementations: os.DirFS, embed.FS, fstest.MapFS.
//   Any code that accepts fs.FS works with ALL of them — real disk, embedded binary,
//   or in-memory test fixture. This is the io.Reader pattern applied to filesystems.
//
// RUN: go run ./10-filesystem/8-fs-testing-seam
// ============================================================================

// SearchResult holds one matching line from the log search.
type SearchResult struct {
	File       string
	LineNumber int
	Text       string
}

// SearchLogs searches all .log and .txt files under the given filesystem
// for lines containing the keyword (case-insensitive).
//
// KEY DESIGN DECISION: The parameter is fs.FS, NOT a string directory path.
// This single change makes the function testable with fstest.MapFS.
func SearchLogs(fsys fs.FS, keyword string) ([]SearchResult, error) {
	keyword = strings.ToLower(keyword)
	var results []SearchResult

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".log" && ext != ".txt" {
			return nil
		}

		f, err := fsys.Open(path)
		if err != nil {
			slog.Warn("could not open file", "path", path, "error", err)
			return nil
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if strings.Contains(strings.ToLower(line), keyword) {
				results = append(results, SearchResult{
					File:       path,
					LineNumber: lineNum,
					Text:       line,
				})
			}
		}
		return scanner.Err()
	})

	return results, err
}

// ============================================================================
// ConfigLoader — another example of fs.FS as a testing seam
// ============================================================================

// LoadConfig reads a JSON config file from the given filesystem.
// In production: LoadConfig(os.DirFS("/etc/myapp"), "config.json")
// In tests:      LoadConfig(fstest.MapFS{"config.json": &fstest.MapFile{...}})
func LoadConfig(fsys fs.FS, path string) (map[string]string, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config %s: %w", path, err)
	}
	defer f.Close()

	config := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return config, scanner.Err()
}

func main() {
	// Production usage: real files from disk
	realFS := os.DirFS(".")

	results, err := SearchLogs(realFS, "error")
	if err != nil {
		slog.Error("search failed", "error", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Println("No matches found in current directory")
	}

	for _, r := range results {
		fmt.Printf("%s:%d  %s\n", r.File, r.LineNumber, r.Text)
	}

	// Testing with fstest.MapFS (shown below — run the test file to see it work):
	// go test -v ./10-filesystem/8-fs-testing-seam
	fmt.Print(`
To see the fs.FS testing seam in action, run:
  go test -v ./10-filesystem/8-fs-testing-seam

The test uses fstest.MapFS — zero disk I/O, fully in-memory.
No temp directories. No cleanup. Runs in microseconds.

KEY TAKEAWAY:
  - Accept fs.FS instead of string paths for testable filesystem code
  - os.DirFS("/path") wraps a real directory as fs.FS
  - fstest.MapFS is the in-memory test double — no temp files needed
  - fs.WalkDir, fs.ReadFile, fs.Glob all work with any fs.FS implementation
  - embed.FS from Section 10/5 also implements fs.FS — reuse the same code
`)
}
