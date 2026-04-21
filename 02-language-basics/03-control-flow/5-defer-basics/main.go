// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics — Lesson CF.5: Defer — Mechanics & Order
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use the defer keyword
//   - The execution order of multiple defers (LIFO)
//   - When deferred functions are evaluated
//
// WHY THIS MATTERS:
//   Defer is Go's primary mechanism for resource cleanup and safety.
//
// RUN: go run ./02-language-basics/03-control-flow/5-defer-basics
// ============================================================================

package main

import (
	"fmt"
)

func main() {
	fmt.Println("CF.5: Defer Mechanics")
	fmt.Println("--------------------------------")

	defer fmt.Println("  (Defer 1) This runs LAST")
	defer fmt.Println("  (Defer 2) This runs SECOND")

	fmt.Println("Performing main work...")
	fmt.Println("Work complete. Returning now.")

	// KEY TAKEAWAY:
	// - defer schedules a function call to run when the surrounding function returns.
	// - Multiple defers run in Last-In-First-Out (LIFO) order.

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CF.6 defer-use-cases")
	fmt.Println("Current: CF.5 (defer-basics)")
	fmt.Println("---------------------------------------------------")
}
