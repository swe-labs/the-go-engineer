// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Break / Continue
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how to change a loop's behavior after the loop has already started.
//
// WHY THIS MATTERS:
//   - Loop control gives you two important tools: - `continue` skips the rest of the current iteration - `break` stops the loop completely That lets one ...
//
// RUN:
//   go run ./02-language-basics/03-control-flow/3-break-continue
//
// KEY TAKEAWAY:
//   - 'continue' skips the rest of the current loop iteration and moves to the next.
//   - 'break' immediately stops the entire loop.
//   - Used together, they allow fine-grained control over loop execution.
// ============================================================================

package main

import "fmt"

func main() {
	fmt.Println("Odd numbers until the stop point:")

	for i := 1; i <= 10; i++ {
		// If the number is even (i % 2 == 0), skip the rest of this iteration.
		// The loop does not stop, it just moves to the next 'i'.
		if i%2 == 0 {
			continue
		}

		// If we hit 7, we want to stop processing entirely.
		// 'break' exits the loop immediately, so 7 and any subsequent numbers
		// will not be printed or evaluated.
		if i == 7 {
			break
		}

		fmt.Println(i)
	}

	// Forward reference:
	// We've seen how 'if' conditions can get nested and complex. Next, we will
	// use 'switch' to evaluate multiple discrete branches cleanly.
	// See: ../4-switch/README.md
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CF.4 -> 02-language-basics/03-control-flow/4-switch")
	fmt.Println("Current: CF.3 (break / continue)")
	fmt.Println("Previous: CF.2 (for basics)")
	fmt.Println("---------------------------------------------------")
}
