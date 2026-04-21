// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// ============================================================================
// Stage 05: Filesystem — File I/O
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - os.WriteFile / os.ReadFile — simple one-shot file operations
//   - os.Create / os.Open — file handles for streaming I/O
//   - bufio.Scanner — reading files line by line (memory efficient)
//   - os.OpenFile — fine-grained control with flags (append, create, etc.)
//   - File permissions (0644, 0755) — the Unix permission model
//   - defer file.Close() — the essential cleanup pattern
//
// ENGINEERING DEPTH:
//   When you call `os.Open()`, the Go runtime makes a "syscall" (System Call)
//   down to the Linux/Windows Kernel. The OS allocates memory for a file tracking
//   struct and returns a lightweight integer called a "File Descriptor" (FD).
//   Linux environments have a hard limit on how many FDs a single process can
//   hold open simultaneously (often just 1,024). If you forget to `defer file.Close()`,
//   your server will eventually "leak" descriptors and entirely crash with a
//   "too many open files" panic.
//
// ANALOGY:
//   Think of file operations like using a library:
//   - os.ReadFile = "Give me a photocopy of the entire book" (loads all into memory)
//   - bufio.Scanner = "Let me read one page at a time" (memory efficient)
//   - os.OpenFile = "I need to highlight and add notes" (read/write/append)
//
// FILE PERMISSIONS (Unix):
//   0644 = owner: read+write (6), group: read (4), others: read (4)
//   0755 = owner: all (7), group: read+execute (5), others: read+execute (5)
//   The first 0 means it's an octal number (base 8).
//
// RUN: go run ./05-packages-io/02-io-and-cli/filesystem/1-files
// ============================================================================

func main() {
	fmt.Println("=== File I/O: Reading and Writing Files ===")
	fmt.Println()

	// Create a temp directory for our demo files.
	// os.MkdirTemp creates a temporary directory that won't conflict with
	// other programs. The second argument is a prefix for the directory name.
	tmpDir, err := os.MkdirTemp("", "go-bible-files-*")
	if err != nil {
		log.Fatal("Failed to create temp dir:", err)
	}
	// Clean up ALL demo files when main() exits.
	// os.RemoveAll recursively deletes the directory and everything in it.
	defer os.RemoveAll(tmpDir)
	fmt.Printf("Working in: %s\n\n", tmpDir)

	// =====================================================================
	// 1. os.WriteFile — Write an entire file in one call
	// =====================================================================
	// os.WriteFile(path, data, permission)
	//   - path: where to write
	//   - data: []byte content (strings must be converted with []byte(...))
	//   - permission: Unix file permission (0644 is standard for regular files)
	//
	// This is the SIMPLEST way to write a file. It creates the file if it
	// doesn't exist, or OVERWRITES it if it does.
	filePath := filepath.Join(tmpDir, "servers.txt")
	content := "web-01  10.0.1.10  us-east-1\ndb-01   10.0.2.20  us-west-2\napi-01  10.0.3.30  eu-west-1\n"

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Fatal("WriteFile failed:", err)
	}
	fmt.Println("1️⃣  Wrote servers.txt with os.WriteFile")

	// =====================================================================
	// 2. os.ReadFile — Read an entire file into memory
	// =====================================================================
	// os.ReadFile(path) returns ([]byte, error).
	// The entire file is loaded into RAM as a byte slice.
	//
	// CAUTION: Don't use this for large files (> 100MB).
	// A 1GB log file would consume 1GB of RAM. Use bufio.Scanner instead.
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("ReadFile failed:", err)
	}
	// Convert []byte to string for printing.
	// string(data) creates a new string from the byte slice.
	fmt.Printf("2️⃣  Read servers.txt:\n%s\n", string(data))

	// =====================================================================
	// 3. bufio.Scanner — Read a file LINE BY LINE
	// =====================================================================
	// For large files, reading line-by-line is memory efficient.
	// bufio.Scanner buffers reads internally (default 64KB buffer).
	//
	// Pattern:
	//   file := os.Open(path)
	//   defer file.Close()       ← Essential cleanup
	//   scanner := bufio.NewScanner(file)
	//   for scanner.Scan() {     ← Reads one line per iteration
	//       line := scanner.Text()
	//   }
	fmt.Println("3️⃣  Reading line-by-line with bufio.Scanner:")
	readLineByLine(filePath)
	fmt.Println()

	// =====================================================================
	// 4. os.OpenFile — Append to a file
	// =====================================================================
	// os.OpenFile gives you fine-grained control with FLAGS:
	//   os.O_APPEND — Write at the END of the file (don't overwrite)
	//   os.O_CREATE — Create the file if it doesn't exist
	//   os.O_WRONLY — Open for writing only
	//   os.O_RDONLY — Open for reading only
	//   os.O_RDWR   — Open for reading AND writing
	//   os.O_TRUNC  — Truncate (empty) the file when opening
	//
	// Combine flags with | (bitwise OR):
	//   os.O_APPEND | os.O_CREATE | os.O_WRONLY
	appendPath := filepath.Join(tmpDir, "deploy-log.txt")

	// First write — creates the file
	appendToFile(appendPath, "Deployed web-01 to us-east-1\n")
	appendToFile(appendPath, "Deployed db-01 to us-west-2\n")
	appendToFile(appendPath, "Deployed api-01 to eu-west-1\n")
	fmt.Println("4️⃣  Appended 3 entries to deploy-log.txt")

	// Read back the appended content
	logData, _ := os.ReadFile(appendPath)
	fmt.Printf("    Contents:\n%s\n", string(logData))

	// =====================================================================
	// 5. os.Create — Create a new file (or truncate existing)
	// =====================================================================
	// os.Create is shorthand for os.OpenFile(path, O_RDWR|O_CREATE|O_TRUNC, 0666).
	// WARNING: If the file exists, it EMPTIES it first (truncates to 0 bytes).
	createPath := filepath.Join(tmpDir, "config.ini")
	file, err := os.Create(createPath)
	if err != nil {
		log.Fatal("Create failed:", err)
	}
	// ALWAYS defer Close() immediately after a successful open/create.
	// If you forget, the OS file handle leaks (limited to ~1024 per process).
	defer file.Close()

	// Write using the file handle's WriteString method
	file.WriteString("[server]\n")
	file.WriteString("host = 0.0.0.0\n")
	file.WriteString("port = 8080\n")
	file.WriteString("workers = 4\n")
	fmt.Println("5️⃣  Created config.ini with os.Create")
	fmt.Println()

	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - os.WriteFile / os.ReadFile: one-shot for small files")
	fmt.Println("  - bufio.Scanner: line-by-line for large files (memory safe)")
	fmt.Println("  - os.OpenFile with flags: append, create, read/write control")
	fmt.Println("  - os.Create: new file or truncate existing (CAUTION: erases!)")
	fmt.Println("  - ALWAYS defer file.Close() after opening a file")
	fmt.Println("  - File permissions: 0644 (read-write owner, read others)")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: FS.2 paths")
	fmt.Println("   Current: FS.1 (files)")
	fmt.Println("---------------------------------------------------")
}

// readLineByLine demonstrates the bufio.Scanner pattern.
// This is the standard way to process large files in Go.
func readLineByLine(path string) {
	// os.Open opens for reading only (O_RDONLY).
	// It returns a *os.File handle and an error.
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Open failed:", err)
	}
	defer file.Close() // Cleanup: close when this function returns

	// bufio.NewScanner wraps the file reader with a buffer.
	// By default, it splits input by lines (scanner.Split(bufio.ScanLines)).
	// You can also split by words, bytes, or a custom split function.
	scanner := bufio.NewScanner(file)
	lineNum := 1

	// scanner.Scan() reads one line and returns true if successful.
	// When there are no more lines (EOF), it returns false.
	for scanner.Scan() {
		// scanner.Text() returns the current line as a string (without the newline)
		fmt.Printf("    Line %d: %s\n", lineNum, scanner.Text())
		lineNum++
	}

	// Always check for scanner errors after the loop.
	// Scan() can fail mid-file (disk error, permission change, etc.)
	if err := scanner.Err(); err != nil {
		log.Fatal("Scanner error:", err)
	}
}

// appendToFile adds text to the end of a file (creates it if needed).
func appendToFile(path, text string) {
	// Open with append + create + write-only flags
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("OpenFile failed:", err)
	}
	defer f.Close()

	// WriteString writes the string and returns (bytesWritten, error)
	if _, err := f.WriteString(text); err != nil {
		log.Fatal("WriteString failed:", err)
	}
}
