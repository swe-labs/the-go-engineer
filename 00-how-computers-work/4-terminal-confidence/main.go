// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work
// Title: Terminal confidence
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How the shell launches programs and connects their output streams.
//   - The difference between stdout and stderr.
//
// WHY THIS MATTERS:
//   - The terminal is a professional interface for process control. Understanding
//     output streams is critical for logging, debugging, and piping data
//     between tools.
//
// RUN:
//   go run ./00-how-computers-work/4-terminal-confidence
//
// KEY TAKEAWAY:
//   - A program's visibility to the outside world is mediated through its
//     file descriptors (stdout, stderr).
// ============================================================================

package main

import (
	"fmt"
	"os"
)

// main (Function): entry point for the program. It writes to both stdout and stderr
// to illustrate how programs expose multiple output streams that the shell can
// redirect independently.
func main() {
	// Standard Output (stdout)
	fmt.Println("This goes to standard output (stdout)")

	// Standard Error (stderr)
	fmt.Fprintln(os.Stderr, "This goes to standard error (stderr)")

	// - Programs can write to multiple output streams
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.5 -> 00-how-computers-work/5-os-processes")
	fmt.Println("Run    : go run ./00-how-computers-work/5-os-processes")
	fmt.Println("Current: HC.4 (terminal-confidence)")
	fmt.Println("---------------------------------------------------")
}
