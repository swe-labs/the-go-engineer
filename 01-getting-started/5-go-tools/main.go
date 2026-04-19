// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 01: Getting Started — Lesson GT.5: Go Tools
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use go fmt to standardize code style
//   - How to use go vet to find common mistakes
//   - How to use go doc to read documentation
//
// WHY THIS MATTERS:
//   Standardized tools ensure that all Go codebases look and feel the same.
//
// RUN: go run ./01-getting-started/5-go-tools
// ============================================================================

package main

import (
	"fmt"
)

func main() {
	fmt.Println("GT.5: Mastering Go Tools")
	fmt.Println("--------------------------------")

	fmt.Println("1. go fmt: Formats your code.")
	fmt.Println("2. go vet: Examines source code and reports suspicious constructs.")
	fmt.Println("3. go doc: Prints documentation for symbols.")

	// KEY TAKEAWAY:
	// - Use go fmt to keep code clean and readable.
	// - Use go vet to catch bugs before they run.
	// - Use go doc to understand the tools you use.

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: GT.6 reading-compiler-errors")
	fmt.Println("Current: GT.5 (go-tools)")
	fmt.Println("---------------------------------------------------")
}
