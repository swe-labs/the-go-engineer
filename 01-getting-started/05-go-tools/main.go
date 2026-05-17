// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 01: Getting Started
// Title: go fmt, go vet, go doc
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Master the three essential tools that keep Go code clean, safe, and documented: `fmt`, `vet`, and `doc`.
//
// WHY THIS MATTERS:
//   - Think of these tools as your "automated senior engineer": 1. `go fmt`: Fixes your style.
//     2. `go vet`: Catches suspicious logic. 3. `go doc`: Explains how tools work.
//   - Standardized tools ensure that all Go codebases look and feel the same.
//
// RUN:
//   go run ./01-getting-started/05-go-tools
//
// KEY TAKEAWAY:
//   - Use the machine to enforce quality so you can focus on solving engineering problems.
// ============================================================================

package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

// runTool (Function): runs a go tool command and returns its output for display.
func runTool(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("(error: %v)", err)
	}
	return string(out)
}

func main() {
	fmt.Println("GT.5: Mastering Go Tools")
	fmt.Println("--------------------------------")
	fmt.Println()

	// 1. go fmt checks formatting without modifying (go fmt -d shows diffs)
	fmt.Println("1. go fmt -d (check formatting, show diffs):")
	fmt.Println(string(rune(9)), "Running: go fmt -d ./01-getting-started/05-go-tools/")
	out := runTool("go", "fmt", "-d", "./01-getting-started/05-go-tools/")
	if out == "" {
		fmt.Println(string(rune(9)), "No formatting issues found.")
	} else {
		fmt.Println(out)
	}
	fmt.Println()

	// 2. go vet checks for suspicious constructs
	fmt.Println("2. go vet (static analysis):")
	fmt.Println(string(rune(9)), "Running: go vet ./01-getting-started/05-go-tools/")
	vetOut := runTool("go", "vet", "./01-getting-started/05-go-tools/")
	if vetOut == "" {
		fmt.Println(string(rune(9)), "No issues found by go vet.")
	} else {
		fmt.Println(vetOut)
	}
	fmt.Println()

	// 3. go doc shows documentation
	fmt.Println("3. go doc (symbol documentation):")
	fmt.Println(string(rune(9)), "Documentation for fmt.Println:")
	docOut := runTool("go", "doc", "fmt.Println")
	if len(docOut) > 300 {
		docOut = docOut[:300] + "..."
	}
	fmt.Print(docOut)
	fmt.Println()

	// Show which Go version is active
	fmt.Println()
	fmt.Println("Go version:", runtime.Version())

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: GT.6 -> 01-getting-started/06-reading-compiler-errors")
	fmt.Println("Run    : go run ./01-getting-started/06-reading-compiler-errors")
	fmt.Println("Current: GT.5 (go-tools)")
	fmt.Println("---------------------------------------------------")
}
