// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: For Basics
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How Go uses a single keyword (`for`) for all looping logic.
//   - Three-part counted loops (init; condition; post).
//   - Condition-only loops (Go's version of `while`).
//   - A preview of `range` for collection iteration.
//
// WHY THIS MATTERS:
//   - Automation is the core of software engineering. Loops allow us to process
//     thousands of items (requests, logs, users) using the same logic without
//     duplicating code.
//
// RUN:
//   go run ./02-language-basics/03-control-flow/2-for-basics
//
// KEY TAKEAWAY:
//   - Go has exactly one looping keyword: `for`. Its simplicity reduces
//     syntactic overhead and forces consistency across the entire ecosystem.
// ============================================================================

package main

import "fmt"

func main() {
	fmt.Println("Counted loop:")
	// The classic three-part loop: init statement; condition; post statement.
	// Notice there are no parentheses around the loop signature.
	for i := 1; i <= 5; i++ {
		fmt.Printf("step %d\n", i)
	}

	fmt.Println()
	fmt.Println("Condition-only loop:")
	// Go doesn't have a 'while' keyword. A 'for' loop with only a condition
	// acts exactly like a traditional while loop.
	countdown := 3
	for countdown > 0 {
		fmt.Printf("countdown: %d\n", countdown)
		countdown--
	}

	fmt.Println()
	fmt.Println("Range preview:")

	words := []string{"go", "learn", "repeat"}
	for _, word := range words {
		fmt.Printf("word = %s\n", word)
	}

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CF.3 -> 02-language-basics/03-control-flow/3-break-continue")
	fmt.Println("Run    : go run ./02-language-basics/03-control-flow/3-break-continue")
	fmt.Println("Current: CF.2 (for basics)")
	fmt.Println("---------------------------------------------------")
}
