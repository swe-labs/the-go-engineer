// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ============================================================================
// Section 10: Filesystem — Log Search Tool (Exercise)
// Level: Intermediate → Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Recursive directory traversal with filepath.WalkDir
//   - Line-by-line file reading with bufio.Scanner
//   - Filtering files by extension
//   - Building a real CLI search tool from scratch
//
// ENGINEERING DEPTH:
//   `filepath.WalkDir` (introduced in Go 1.16) is the replacement for the older
//   `filepath.Walk`. The critical difference is that `WalkDir` uses `fs.DirEntry`
//   which is a lazy interface — it does NOT call `os.Stat()` on every single file.
//   The old `Walk` function called `os.Stat()` on every file, which triggers a
//   system call to the kernel for each file. On a directory with 100,000 files,
//   `WalkDir` can be 10x faster because it avoids these redundant stat() syscalls.
//
// RUN: go run ./10-filesystem/7-log-search
// ============================================================================

// SearchResult holds the location and content of a matching line.
type SearchResult struct {
	FilePath   string // Full path to the file containing the match
	LineNumber int    // 1-indexed line number
	LineText   string // The actual content of the matching line
}

// searchFile opens a single file and scans it line-by-line for the keyword.
// It returns a slice of all matching results found in that file.
//
// `bufio.Scanner` reads the file lazily — it never loads the entire file into
// memory at once. This is critical for searching large log files (1GB+).
func searchFile(filePath string, keyword string) ([]SearchResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []SearchResult
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Case-insensitive search: convert both to lowercase before comparing.
		if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
			results = append(results, SearchResult{
				FilePath:   filePath,
				LineNumber: lineNum,
				LineText:   line,
			})
		}
	}

	return results, scanner.Err()
}

// searchDirectory walks a directory tree recursively and searches all
// .log and .txt files for the given keyword.
//
// `filepath.WalkDir` calls our function for every file and directory
// it encounters. We filter by extension to only search relevant files.
func searchDirectory(rootDir string, keyword string) ([]SearchResult, error) {
	var allResults []SearchResult

	// Allowed file extensions for searching.
	allowedExts := map[string]bool{
		".log": true,
		".txt": true,
	}

	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// If we can't access a file/dir, skip it rather than crashing.
			fmt.Fprintf(os.Stderr, "  ⚠ Cannot access %s: %v\n", path, err)
			return nil
		}

		// Skip directories — we only want files.
		if d.IsDir() {
			return nil
		}

		// Check file extension.
		ext := strings.ToLower(filepath.Ext(path))
		if !allowedExts[ext] {
			return nil
		}

		// Search this file for the keyword.
		results, err := searchFile(path, keyword)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Error reading %s: %v\n", path, err)
			return nil
		}

		allResults = append(allResults, results...)
		return nil
	})

	return allResults, err
}

func main() {
	fmt.Println("=== Log Search Tool ===")
	fmt.Println()

	// For this exercise, we create sample log files to search through.
	// In a real tool, you'd accept the directory and keyword from os.Args.
	tmpDir, err := os.MkdirTemp("", "log-search-*")
	if err != nil {
		fmt.Println("Error creating temp dir:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	// Create sample log files
	sampleLogs := map[string]string{
		"app.log": `[2024-01-15 10:00:01] INFO: Server started on port 8080
[2024-01-15 10:00:05] INFO: Database connection established
[2024-01-15 10:05:12] WARNING: Slow query detected (2.3s)
[2024-01-15 10:10:33] ERROR: Connection timeout to payment gateway
[2024-01-15 10:10:34] ERROR: Failed to process order #1234
[2024-01-15 10:15:00] INFO: Retry successful for payment gateway`,
		"access.log": `192.168.1.10 - GET /api/users 200 12ms
192.168.1.15 - POST /api/orders 201 45ms
192.168.1.10 - GET /api/orders/1234 500 ERROR: internal server failure
10.0.0.1 - GET /health 200 1ms`,
		"notes.txt": `TODO: Fix the connection timeout error in payment module
TODO: Add retry logic for database connections
DONE: Deploy v2.1 to staging`,
	}

	// Write the sample files
	for name, content := range sampleLogs {
		path := filepath.Join(tmpDir, name)
		os.WriteFile(path, []byte(content), 0644)
	}

	// Search for "error" across all files
	keyword := "error"
	fmt.Printf("Searching for \"%s\" in %s...\n\n", keyword, tmpDir)

	results, err := searchDirectory(tmpDir, keyword)
	if err != nil {
		fmt.Println("Search failed:", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Println("No matches found.")
		return
	}

	fmt.Printf("Found %d matches:\n\n", len(results))
	for _, r := range results {
		// Show only the filename (not the full temp path) for clean output.
		filename := filepath.Base(r.FilePath)
		fmt.Printf("  %s:%d  %s\n", filename, r.LineNumber, r.LineText)
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - filepath.WalkDir recursively traverses directories efficiently")
	fmt.Println("  - bufio.Scanner reads files line-by-line (constant memory)")
	fmt.Println("  - Filter by extension to avoid searching binary files")
	fmt.Println("  - Always handle errors gracefully — skip files you can't read")
}
