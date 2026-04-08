// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// RUN: go run ./09-io-and-cli/filesystem/8-fs-testing-seam
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
// Section 09 Stretch: io/fs as a Testing Seam
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why accepting fs.FS instead of string paths makes code testable
//   - os.DirFS for wrapping a real directory as an fs.FS
//   - fstest.MapFS as an in-memory filesystem in tests
//   - The "accept interfaces, return structs" principle at the filesystem level
//
// ENGINEERING DEPTH:
//   fs.FS applies the same decoupling idea as io.Reader and io.Writer.
//   Code that accepts fs.FS can work with real disk, embedded assets, or
//   in-memory test fixtures without changing its core logic.
// ============================================================================

// SearchResult holds one matching line from the log search.
type SearchResult struct {
	File       string
	LineNumber int
	Text       string
}

// SearchLogs searches all .log and .txt files under the given filesystem
// for lines containing the keyword (case-insensitive).
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

// LoadConfig reads a key=value config file from the given filesystem.
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

	fmt.Print(`
To see the fs.FS testing seam in action, run:
  go test -v ./09-io-and-cli/filesystem/8-fs-testing-seam

The test uses fstest.MapFS — zero disk I/O, fully in-memory.
No temp directories. No cleanup. Runs in microseconds.

KEY TAKEAWAY:
  - Accept fs.FS instead of string paths for testable filesystem code
  - os.DirFS("/path") wraps a real directory as fs.FS
  - fstest.MapFS is the in-memory test double — no temp files needed
  - fs.WalkDir, fs.ReadFile, fs.Glob all work with any fs.FS implementation
  - embed.FS from Section 09/FS.5 also implements fs.FS — reuse the same code
`)
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT STEP: Continue to Section 10 web and database")
	fmt.Println("Current: FS.8 (fs.FS testing seam)")
	fmt.Println("---------------------------------------------------")
}
