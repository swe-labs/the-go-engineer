// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 10: Filesystem — Log Search Tool (Exercise Starter)
// Level: Intermediate → Advanced
// ============================================================================
//
// EXERCISE: Build a Directory Traversal Log Search Tool
//
// REQUIREMENTS:
//  1. [ ] Define a `SearchResult` struct with FilePath, LineNumber, LineText
//  2. [ ] Implement `searchFile(path, keyword string) ([]SearchResult, error)`
//         using bufio.Scanner for line-by-line reading
//  3. [ ] Implement `searchDirectory(rootDir, keyword string) ([]SearchResult, error)`
//         using filepath.WalkDir for recursive traversal
//  4. [ ] Filter by extension — only search .log and .txt files
//  5. [ ] Make the search case-insensitive using strings.ToLower
//  6. [ ] In main(), create sample log files and search them
//
// HINTS:
//   - bufio.NewScanner(file) reads line-by-line (memory efficient)
//   - filepath.WalkDir(root, func(path, d, err)) walks recursively
//   - d.IsDir() checks if the current entry is a directory
//   - filepath.Ext(path) returns the file extension including the dot
//
// RUN: go run ./10-filesystem/7-log-search/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Define SearchResult struct

// TODO: Implement searchFile

// TODO: Implement searchDirectory

func main() {
	fmt.Println("=== Log Search Tool Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your log search tool!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
