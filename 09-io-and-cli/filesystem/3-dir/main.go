// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// ============================================================================
// Section 09: Filesystem — Directory Operations
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - os.Mkdir / os.MkdirAll: creating directories
//   - os.ReadDir: listing directory contents
//   - filepath.WalkDir: recursively walking a directory tree
//   - os.RemoveAll: deleting directories
//   - os.Stat: checking if a file/directory exists
//
// ANALOGY:
//   os.Mkdir is like creating a single folder.
//   os.MkdirAll is like creating an entire folder path (mkdir -p).
//   filepath.WalkDir is like the `find` command — it visits every file
//   in a directory tree recursively.
//
// ENGINEERING DEPTH:
//   In Posix filesystems (Linux/macOS), a "Directory" is secretly just a specialized
//   file! Instead of containing text, a directory file contains a list of "inodes"
//   (Index Nodes) pointing to the physical disk sectors of the files inside it.
//   Because directories are just strings mapping to inodes, moving a 100GB file
//   to a different folder on the same disk is instantaneous—the OS just rewrites
//   a 16-byte inode pointer; it doesn't move any actual data!
//
// RUN: go run ./09-io-and-cli/filesystem/3-dir
// ============================================================================

func main() {
	fmt.Println("=== Directory Operations ===")
	fmt.Println()

	// Create a temp directory for our demo
	tmpDir, err := os.MkdirTemp("", "the-go-engineer-dirs-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // Cleanup everything when done
	fmt.Printf("  Working in: %s\n\n", tmpDir)

	// =====================================================================
	// 1. os.MkdirAll — Create nested directory structure
	// =====================================================================
	// os.MkdirAll creates the directory AND all parent directories.
	// It's like `mkdir -p` in the shell. If the directory already exists, it does nothing.
	// Permission 0755 = owner: rwx, group: r-x, others: r-x
	dirs := []string{
		filepath.Join(tmpDir, "project", "src", "handlers"),
		filepath.Join(tmpDir, "project", "src", "models"),
		filepath.Join(tmpDir, "project", "config"),
		filepath.Join(tmpDir, "project", "docs"),
	}

	fmt.Println("  1️⃣  Creating directory structure:")
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal("MkdirAll failed:", err)
		}
		// Print relative path for readability
		rel, _ := filepath.Rel(tmpDir, dir)
		fmt.Printf("     📁 %s\n", rel)
	}
	fmt.Println()

	// Create some dummy files for the walk demo
	dummyFiles := []string{
		filepath.Join(tmpDir, "project", "go.mod"),
		filepath.Join(tmpDir, "project", "main.go"),
		filepath.Join(tmpDir, "project", "src", "handlers", "user.go"),
		filepath.Join(tmpDir, "project", "src", "models", "user.go"),
		filepath.Join(tmpDir, "project", "config", "app.yaml"),
		filepath.Join(tmpDir, "project", "docs", "README.md"),
	}
	for _, f := range dummyFiles {
		os.WriteFile(f, []byte("// placeholder"), 0644)
	}

	// =====================================================================
	// 2. os.ReadDir — List directory contents (non-recursive)
	// =====================================================================
	// os.ReadDir returns a sorted list of directory entries.
	// Each entry has Name(), IsDir(), and Type() info.
	fmt.Println("  2️⃣  os.ReadDir (list one level):")
	projectDir := filepath.Join(tmpDir, "project")
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		log.Fatal("ReadDir failed:", err)
	}

	for _, entry := range entries {
		icon := "📄"
		if entry.IsDir() {
			icon = "📁"
		}
		fmt.Printf("     %s %s\n", icon, entry.Name())
	}
	fmt.Println()

	// =====================================================================
	// 3. filepath.WalkDir — Recursive directory traversal
	// =====================================================================
	// WalkDir visits EVERY file and directory in the tree.
	// The callback receives: path, directory entry info, and any error.
	// This is Go's equivalent of the Unix `find` command.
	fmt.Println("  3️⃣  filepath.WalkDir (recursive tree):")
	filepath.WalkDir(projectDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // Skip entries with errors
		}

		// Calculate indent based on depth
		rel, _ := filepath.Rel(projectDir, path)
		depth := len(filepath.SplitList(rel))
		if rel == "." {
			return nil // Skip root
		}

		icon := "📄"
		if d.IsDir() {
			icon = "📁"
		}
		indent := ""
		for i := 0; i < depth; i++ {
			indent += "  "
		}
		fmt.Printf("     %s%s %s\n", indent, icon, d.Name())
		return nil
	})
	fmt.Println()

	// =====================================================================
	// 4. os.Stat — Check if file/directory exists
	// =====================================================================
	// os.Stat returns file info. If the file doesn't exist, err != nil.
	// Use os.IsNotExist(err) to check specifically for "not found".
	fmt.Println("  4️⃣  Check existence:")
	checkPaths := []string{
		filepath.Join(tmpDir, "project", "main.go"),
		filepath.Join(tmpDir, "project", "missing.go"),
	}
	for _, p := range checkPaths {
		rel, _ := filepath.Rel(tmpDir, p)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			fmt.Printf("     ❌ %s — does NOT exist\n", rel)
		} else {
			fmt.Printf("     ✅ %s — exists\n", rel)
		}
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - os.MkdirAll creates nested dirs (like mkdir -p)")
	fmt.Println("  - os.ReadDir lists one level, filepath.WalkDir walks recursively")
	fmt.Println("  - os.Stat + os.IsNotExist checks if path exists")
	fmt.Println("  - os.RemoveAll recursively deletes directory and all contents")
	fmt.Println("  - Always use filepath.Join for cross-platform paths")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: FS.4 temp files")
	fmt.Println("   Current: FS.3 (directories)")
	fmt.Println("---------------------------------------------------")
}
