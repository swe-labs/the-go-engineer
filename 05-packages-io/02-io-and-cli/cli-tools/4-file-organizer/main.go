// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./05-packages-io/02-io-and-cli/cli-tools/4-file-organizer
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ============================================================================
// Stage 05: CLI Tools — File Organizer (Exercise)
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Using the `flag` package for command-line options
//   - Reading directory contents with os.ReadDir
//   - Organizing files by extension into subdirectories
//   - The --dry-run pattern for safe testing
//
// ENGINEERING DEPTH:
//   The `flag` package uses a registration pattern — you call `flag.String()`
//   to register each flag, which returns a *string pointer. The actual parsing
//   doesn't happen until you call `flag.Parse()`. This two-step approach allows
//   you to define flags anywhere in your code (even in init() functions across
//   packages) before a single centralized Parse() call processes them all.
//
// USAGE: go run ./05-packages-io/02-io-and-cli/cli-tools/4-file-organizer --dir=./my-messy-folder
// FLAGS: --dir (required), --dry-run (optional, default false)
// ============================================================================

func main() {
	// Define CLI flags using the flag package.
	// flag.String returns a *string — a pointer to the flag's value.
	dir := flag.String("dir", "", "Directory to organize (required)")
	dryRun := flag.Bool("dry-run", false, "Preview changes without moving files")

	// flag.Parse() processes os.Args and assigns values to the registered flags.
	flag.Parse()

	// Validate required flag
	if *dir == "" {
		fmt.Println("=== File Organizer ===")
		fmt.Println()
		fmt.Println("Organizes files in a directory by extension.")
		fmt.Println()
		flag.Usage()
		fmt.Println()
		fmt.Println("Running demo mode...")
		fmt.Println()

		// Create a demo directory for self-contained execution
		tmpDir, _ := os.MkdirTemp("", "organizer-demo-*")
		defer os.RemoveAll(tmpDir)

		// Create sample files with various extensions
		sampleFiles := []string{
			"report.pdf", "photo.jpg", "notes.txt", "main.go",
			"styles.css", "readme.md", "data.json", "script.py",
		}
		for _, name := range sampleFiles {
			os.WriteFile(filepath.Join(tmpDir, name), []byte("sample"), 0644)
		}

		*dir = tmpDir
		*dryRun = true // Demo always runs in dry-run mode
	}

	fmt.Printf("📂 Organizing: %s\n", *dir)
	if *dryRun {
		fmt.Println("🔍 DRY RUN — no files will actually be moved")
	}
	fmt.Println()

	// Read directory contents
	entries, err := os.ReadDir(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
		os.Exit(1)
	}

	moved := 0
	for _, entry := range entries {
		// Skip directories — we only organize files.
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		ext := strings.TrimPrefix(filepath.Ext(filename), ".") // "go", "txt", etc.

		// Skip files with no extension
		if ext == "" {
			fmt.Printf("  ⏭ Skipping (no extension): %s\n", filename)
			continue
		}

		// Target subdirectory named after the extension
		targetDir := filepath.Join(*dir, ext)
		targetPath := filepath.Join(targetDir, filename)
		sourcePath := filepath.Join(*dir, filename)

		if *dryRun {
			fmt.Printf("  📄 %s → %s/\n", filename, ext)
		} else {
			// Create the target subdirectory if it doesn't exist.
			// os.MkdirAll creates all parent directories as needed.
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "  ❌ Error creating dir %s: %v\n", targetDir, err)
				continue
			}

			// Move the file by renaming it to the new path.
			if err := os.Rename(sourcePath, targetPath); err != nil {
				fmt.Fprintf(os.Stderr, "  ❌ Error moving %s: %v\n", filename, err)
				continue
			}
			fmt.Printf("  ✅ %s → %s/\n", filename, ext)
		}
		moved++
	}

	fmt.Println()
	action := "would be moved"
	if !*dryRun {
		action = "moved"
	}
	fmt.Printf("Done! %d files %s.\n", moved, action)

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - flag.String/Bool/Int register flags; flag.Parse() processes them")
	fmt.Println("  - os.ReadDir reads directory contents efficiently")
	fmt.Println("  - filepath.Ext returns the extension (including the dot)")
	fmt.Println("  - --dry-run is essential for any CLI that modifies the filesystem")
}
