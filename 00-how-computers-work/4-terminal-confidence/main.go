// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work — Terminal Confidence
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to output to stdout vs stderr
//   - How programs interact with the terminal
//
// WHY THIS MATTERS:
//   Production logs often split standard output and errors. Understanding how to
//   write to them is critical for observability.
//
// RUN: go run ./00-how-computers-work/4-terminal-confidence
// ============================================================================

package main

import (
	"fmt"
	"os"
)

func main() {
	// Standard Output (stdout)
	fmt.Println("This goes to standard output (stdout)")

	// Standard Error (stderr)
	fmt.Fprintln(os.Stderr, "This goes to standard error (stderr)")

	// KEY TAKEAWAY:
	// - Programs can write to multiple output streams
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.5 5-os-processes")
	fmt.Println("Run    : go run ./00-how-computers-work/5-os-processes")
	fmt.Println("Current: HC.4 (4-terminal-confidence)")
	fmt.Println("---------------------------------------------------")
}
