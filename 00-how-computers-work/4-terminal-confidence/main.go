// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work
// Title: Terminal confidence
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Build confidence with the terminal as the text-based environment that launches programs, passes arguments, and handles output streams.
//
// WHY THIS MATTERS:
//   - The terminal is not "where scary commands live." It is a process-launching and output-reading interface for the operating system.
//
// RUN:
//   go run ./00-how-computers-work/4-terminal-confidence
//
// KEY TAKEAWAY:
//   - Build confidence with the terminal as the text-based environment that launches programs, passes arguments, and handles output streams.
// ============================================================================

//   Production logs often split standard output and errors. Understanding how to
//   write to them is critical for observability.

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

	// - Programs can write to multiple output streams
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.5 -> 00-how-computers-work/5-os-processes")
	fmt.Println("Run    : go run ./00-how-computers-work/5-os-processes")
	fmt.Println("Current: HC.4 (4-terminal-confidence)")
	fmt.Println("---------------------------------------------------")
}
