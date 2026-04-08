// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 09: CLI Tools - File Organizer (Exercise Starter)
// Level: Intermediate
// ============================================================================
//
// EXERCISE: Build a CLI File Organizer
//
// REQUIREMENTS:
//  1. [ ] Use flag.String for --dir and flag.Bool for --dry-run
//  2. [ ] Read directory contents with os.ReadDir
//  3. [ ] Group files by extension into subdirectories
//  4. [ ] In dry-run mode, print what WOULD happen without moving files
//  5. [ ] In normal mode, use os.MkdirAll and os.Rename to move files
//
// RUN: go run ./09-io-and-cli/cli-tools/4-file-organizer/_starter --dir=./my-folder
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Define CLI flags with flag package

// TODO: Implement the file organization logic

func main() {
	fmt.Println("=== File Organizer Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your file organizer!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
