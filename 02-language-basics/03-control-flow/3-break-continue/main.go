// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Break / Continue
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to skip the remainder of an iteration using `continue`.
//   - How to terminate a loop immediately using `break`.
//   - Why the order of control statements inside a loop matters.
//
// WHY THIS MATTERS:
//   - Real-world loops often encounter "special cases"—invalid data, a
//     successful search result, or an error. `break` and `continue` allow you
//     to handle these cases efficiently without over-complicating the
//     main loop condition.
//
// RUN:
//   go run ./02-language-basics/03-control-flow/3-break-continue
//
// KEY TAKEAWAY:
//   - `continue` skips the current step; `break` stops the entire process.
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

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CF.4 -> 02-language-basics/03-control-flow/4-switch")
	fmt.Println("Run    : go run ./02-language-basics/03-control-flow/4-switch")
	fmt.Println("Current: CF.3 (break / continue)")
	fmt.Println("---------------------------------------------------")
}
