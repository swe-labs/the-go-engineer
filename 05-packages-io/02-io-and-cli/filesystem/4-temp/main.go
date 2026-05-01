// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 05: Packages and I/O
// Title: Temp Files
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to create secure, unique, and cross-platform temporary files and directories.
//   - The critical role of 'defer' in managing temporary resource life cycles.
//
// WHY THIS MATTERS:
//   - Temporary files are essential for testing, large data processing, and
//     intermediate storage. Using the standard library ensures your app
//     doesn't accidentally leak or overwrite system data.
//
// RUN:
//   go run ./05-packages-io/02-io-and-cli/filesystem/4-temp
//
// KEY TAKEAWAY:
//   - Use os.CreateTemp() and os.MkdirTemp() to avoid race conditions and
//     guarantee unique names.
// ============================================================================

// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Stage 05: Filesystem - Temporary Files & Directories
//
//   - os.MkdirTemp: creating secure temporary directories
//   - os.CreateTemp: creating secure temporary files
//   - The importance of defer for cleanup
//   - Why we use temp files (testing, intermediate processing, downloads)
//
// ANALOGY:
//   Temp files are like draft paper on a desk. You use it to work out
//   a complex math problem, and when you're done, you throw it in the
//   recycle bin. The OS guarantees you a clean piece of paper that
//   won't interfere with anyone else's desk.
//
// ENGINEERING DEPTH:
//   Why use these functions instead of hardcoding "/tmp/myfile"?
//   1. Cross-platform: Uses the correct temp dir depending on OS (e.g. C:\Temp or /tmp).
//   2. Race-free: The OS guarantees the generated filename is unique.
//      If two instances of your app run simultaneously, they won't overwrite each other.
//   3. Security: Temp files are created with restrictive permissions (0600),
//      preventing other users on the system from reading your temporary data.

func main() {
	fmt.Println("=== Temporary Files and Directories ===")
	fmt.Println()

	// 1. Temporary Directories
	// os.MkdirTemp(dir, pattern)
	//   dir: Where to create it ("" means the OS default temp directory)
	//   pattern: Prefix and/or suffix for the dir name ("logs-*" generates "logs-123456")

	fmt.Println("1. Temporary Directory:")
	tempDir, err := os.MkdirTemp("", "the-go-engineer-logs-*")
	if err != nil {
		log.Fatal("Failed to create temp dir:", err)
	}

	// ALWAYS defer cleanup immediately.
	// os.RemoveAll is necessary because a directory must be empty to use os.Remove.
	defer func() {
		fmt.Printf("   [CLEANUP] Cleaning up directory: %s\n", tempDir)
		_ = os.RemoveAll(tempDir)
	}()

	fmt.Printf("   [DIR] Created: %s\n", tempDir)
	fmt.Println()

	// 2. Temporary Files
	// os.CreateTemp(dir, pattern) creates and OPENs the file for reading/writing.

	fmt.Println("2. Temporary File:")
	// We create the file INSIDE our temporary directory
	tempFile, err := os.CreateTemp(tempDir, "data-*.csv")
	if err != nil {
		log.Fatal("Failed to create temp file:", err)
	}

	fmt.Printf("   [FILE] Created: %s\n", tempFile.Name())

	// Write some data to the temp file
	data := []byte("id,name,role\n1,admin,sysadmin\n2,user,guest\n")
	if _, err := tempFile.Write(data); err != nil {
		log.Fatal("Failed to write to temp file:", err)
	}

	fmt.Printf("   [WRITE] Wrote %d bytes to temp file\n", len(data))

	// Remember to close the file handle!
	if err := tempFile.Close(); err != nil {
		log.Fatal("Failed to close temp file:", err)
	}

	// 3. Verifying the temporary data
	fmt.Println()
	fmt.Println("3. Reading back temporary data:")

	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		log.Fatal("Failed to read temp file:", err)
	}

	fmt.Printf("   Contents of %s:\n", filepath.Base(tempFile.Name()))
	fmt.Printf("   ---\n   %s   ---\n", string(content))

	// Note: You will see the "[CLEANUP] Cleaning up directory..." print at the very end
	// when the main function returns, executing the deferred os.RemoveAll.

	fmt.Println()
	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  - os.MkdirTemp / os.CreateTemp generate unique, race-free names")
	fmt.Println("  - The first argument \"\" uses the OS default temp directory")
	fmt.Println("  - The second argument \"prefix-*\" determines the naming pattern")
	fmt.Println("  - ALWAYS defer os.RemoveAll(dir) or os.Remove(file) immediately after creation")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: FS.5 -> 05-packages-io/02-io-and-cli/filesystem/5-embed")
	fmt.Println("Current: FS.4 (temp)")
	fmt.Println("Previous: FS.3 (dir)")
	fmt.Println("---------------------------------------------------")
}
