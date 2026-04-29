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
//   - Think of these tools as your "automated senior engineer": 1. `go fmt`: Fixes your style. 2. `go vet`: Catches suspicious logic. 3. `go doc`: Explai...
//
// RUN:
//   go run ./01-getting-started/5-go-tools
//
// KEY TAKEAWAY:
//   - Master the three essential tools that keep Go code clean, safe, and documented: `fmt`, `vet`, and `doc`.
// ============================================================================

//   Standardized tools ensure that all Go codebases look and feel the same.

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

	// - Use go fmt to keep code clean and readable.
	// - Use go vet to catch bugs before they run.
	// - Use go doc to understand the tools you use.

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: GT.6 reading-compiler-errors")
	fmt.Println("Current: GT.5 (go-tools)")
	fmt.Println("---------------------------------------------------")
}
