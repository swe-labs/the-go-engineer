// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// ============================================================================
// Section 09: Filesystem — File Paths
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - filepath.Join: building cross-platform paths (never use "/" or "\" manually)
//   - filepath.Base: extracting the filename from a path
//   - filepath.Dir: extracting the directory from a path
//   - filepath.Ext: extracting the file extension
//   - filepath.Abs: converting relative path to absolute
//   - filepath.Clean: normalizing messy paths (removing ../ and ./)
//   - filepath.Glob: finding files matching a pattern
//
// ENGINEERING DEPTH:
//   Go actually has two path packages: `path` and `path/filepath`. This is a
//   crucial distinction. The standard `path` package ONLY uses forward slashes(`/`)
//   and is specifically meant for manipulating URL paths (like `http://yoursite.com/a/b`).
//   The `path/filepath` package is OS-Aware. If you compile your code on Windows,
//   it automatically switches to backslashes (`\`). Always use `path/filepath`
//   when manipulating local file directories to guarantee cross-OS compatibility.
//
// WHY USE filepath.Join?
//   On Linux/macOS, paths use "/" (forward slash)
//   On Windows, paths use "\" (backslash)
//   filepath.Join automatically uses the correct separator for the OS.
//   NEVER hardcode "/" or "\" in path strings — your code will break on other OSes.
//
// RUN: go run ./09-io-and-cli/filesystem/2-paths
// ============================================================================

func main() {
	fmt.Println("=== File Paths: Cross-Platform Path Handling ===")
	fmt.Println()

	// --- filepath.Join: Build paths safely ---
	// Join combines path segments with the OS-correct separator.
	// It also cleans up redundant separators and dots.
	configPath := filepath.Join("home", "rasel", "projects", "the-go-engineer", "config.yaml")
	fmt.Printf("  Join:  %s\n", configPath)

	// Join handles extra slashes gracefully:
	messyPath := filepath.Join("/usr/", "/local/", "//bin//", "go")
	fmt.Printf("  Clean join: %s\n", messyPath) // Cleaned up: /usr/local/bin/go
	fmt.Println()

	// --- filepath.Base: Extract filename ---
	// Base returns the LAST element of a path (the filename).
	fullPath := "/var/log/app/server.log"
	fmt.Printf("  Base(%q) = %q\n", fullPath, filepath.Base(fullPath)) // "server.log"

	// --- filepath.Dir: Extract directory ---
	// Dir returns everything BEFORE the last element.
	fmt.Printf("  Dir(%q)  = %q\n", fullPath, filepath.Dir(fullPath)) // "/var/log/app"

	// --- filepath.Ext: Extract extension ---
	// Ext returns the file extension including the dot.
	fmt.Printf("  Ext(%q)  = %q\n", fullPath, filepath.Ext(fullPath)) // ".log"

	// Useful pattern: get filename WITHOUT extension
	base := filepath.Base(fullPath)
	nameOnly := base[:len(base)-len(filepath.Ext(base))]
	fmt.Printf("  Name (no ext): %q\n", nameOnly) // "server"
	fmt.Println()

	// --- filepath.Abs: Convert to absolute path ---
	// Abs resolves a relative path to an absolute path based on the current working dir.
	relPath := "./config/app.yaml"
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		fmt.Printf("  Abs error: %v\n", err)
	} else {
		fmt.Printf("  Abs(%q) = %s\n", relPath, absPath)
	}

	// Current working directory for reference
	cwd, _ := os.Getwd()
	fmt.Printf("  Working dir: %s\n", cwd)
	fmt.Println()

	// --- filepath.Clean: Normalize messy paths ---
	// Clean resolves "." (current dir), ".." (parent dir), and double slashes.
	dirtyPaths := []string{
		"./users/./docs/../downloads/./file.txt",
		"/home/rasel/../rasel/./projects//the-go-engineer",
		"a/b/../c/./d",
	}

	fmt.Println("  --- Clean (normalize paths) ---")
	for _, dp := range dirtyPaths {
		fmt.Printf("    %s\n", dp)
		fmt.Printf("    → %s\n\n", filepath.Clean(dp))
	}

	// --- filepath.Glob: Find files matching a pattern ---
	// Glob returns all files matching a wildcard pattern.
	// Patterns: * (any chars), ? (one char), [abc] (one of these)
	goFiles, err := filepath.Glob("*.go")
	if err != nil {
		fmt.Printf("  Glob error: %v\n", err)
	}
	fmt.Printf("  Go files in current dir: %v\n", goFiles)

	// Match all Go files recursively would need filepath.WalkDir (FS.3)
	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - ALWAYS use filepath.Join — never hardcode / or \\")
	fmt.Println("  - filepath.Base extracts filename, Dir extracts directory, Ext extracts extension")
	fmt.Println("  - filepath.Clean normalizes messy paths (resolves . and ..)")
	fmt.Println("  - filepath.Abs converts relative to absolute path")
	fmt.Println("  - filepath is cross-platform — works on Linux, macOS, AND Windows")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: FS.3 directories")
	fmt.Println("   Current: FS.2 (paths)")
	fmt.Println("---------------------------------------------------")
}
