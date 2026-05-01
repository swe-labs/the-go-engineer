// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: For Basics
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how Go repeats work with its single loop keyword: `for`.
//
// WHY THIS MATTERS:
//   - A loop says, "keep doing this work while the rule allows it." Go uses one keyword for several loop shapes: - counted loops - condition-only loops -...
//
// RUN:
//   go run ./02-language-basics/03-control-flow/2-for-basics
//
// KEY TAKEAWAY:
//   - Go has exactly one looping keyword: 'for'. It handles counted loops,
//     condition-only loops (while-style), infinite loops, and collection iteration.
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

	// Backward reference:
	// We touched upon short declarations (:=) in: ../../1-variables/README.md
	// Here we combine it with 'range' to cleanly iterate over a slice collection.
	words := []string{"go", "learn", "repeat"}
	for _, word := range words {
		fmt.Printf("word = %s\n", word)
	}

	// Forward reference:
	// Loops often need to stop early or skip an iteration. We will learn how to
	// control loop execution precisely in: ../3-break-continue/README.md
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CF.3 -> 02-language-basics/03-control-flow/3-break-continue")
	fmt.Println("Current: CF.2 (for basics)")
	fmt.Println("Previous: CF.1 (if / else)")
	fmt.Println("---------------------------------------------------")
}
