// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Defer - mechanics & order
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how to schedule work to happen at the very end of a function, ensuring cleanup always runs.
//
// WHY THIS MATTERS:
//   - Think of `defer` like a "sticky note" you put on the exit door of a room. No matter what you do in the room or which way you leave, you must perfor...
//
// RUN:
//   go run ./02-language-basics/03-control-flow/5-defer-basics
//
// KEY TAKEAWAY:
//   - 'defer' schedules a function call to run exactly when the surrounding
//     function returns. Multiple defers execute in Last-In-First-Out (LIFO) order.
// ============================================================================

package main

import (
	"fmt"
)

func main() {
	fmt.Println("CF.5: Defer Mechanics")
	fmt.Println("--------------------------------")

	// Backward reference:
	// Unlike 'switch' and 'if' logic from previous lessons that execute exactly
	// when they are evaluated, 'defer' intercepts control flow at the exit point.
	// See: ../4-switch/README.md

	// Defer 1 is scheduled first, so it runs last (LIFO order).
	defer fmt.Println("  (Defer 1) This runs LAST")

	// Defer 2 is scheduled second, so it runs before Defer 1.
	defer fmt.Println("  (Defer 2) This runs SECOND")

	fmt.Println("Performing main work...")
	fmt.Println("Work complete. Returning now.")

	fmt.Println()
	// Forward reference:
	// Now that you understand the LIFO execution order of 'defer', we will apply
	// it to real-world cleanup scenarios like closing files in the next lesson.
	// See: ../6-defer-use-cases/README.md
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CF.6 defer-use-cases")
	fmt.Println("Current: CF.5 (defer-basics)")
	fmt.Println("Previous: CF.4 (switch)")
	fmt.Println("---------------------------------------------------")
}
