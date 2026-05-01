// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Defer - mechanics & order
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Scheduling cleanup tasks with the `defer` keyword.
//   - Last-In-First-Out (LIFO) execution order of deferred calls.
//   - When deferred arguments are evaluated.
//
// WHY THIS MATTERS:
//   - Real programs open files, connections, and locks. If a function returns
//     early due to an error, these resources might remain open forever (a "leak").
//     `defer` guarantees that cleanup code runs regardless of which path the
//     function takes to exit.
//
// RUN:
//   go run ./02-language-basics/03-control-flow/5-defer-basics
//
// KEY TAKEAWAY:
//   - `defer` decouples where you declare cleanup from where it executes.
// ============================================================================

package main

import (
	"fmt"
)

func main() {
	fmt.Println("CF.5: Defer Mechanics")
	fmt.Println("--------------------------------")

	// Defer 1 is scheduled first, so it runs last (LIFO order).
	defer fmt.Println("  (Defer 1) This runs LAST")

	// Defer 2 is scheduled second, so it runs before Defer 1.
	defer fmt.Println("  (Defer 2) This runs SECOND")

	fmt.Println("Performing main work...")
	fmt.Println("Work complete. Returning now.")

	fmt.Println()
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CF.6 -> 02-language-basics/03-control-flow/6-defer-use-cases")
	fmt.Println("Run    : go run ./02-language-basics/03-control-flow/6-defer-use-cases")
	fmt.Println("Current: CF.5 (defer-basics)")
	fmt.Println("---------------------------------------------------")
}
