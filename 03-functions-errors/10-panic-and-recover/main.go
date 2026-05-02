// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: panic and recover
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Using panic for unrecoverable "Impossible" states.
//   - Catching panics with recover inside a deferred function.
//   - Understanding stack unwinding during a panic event.
//
// WHY THIS MATTERS:
//   - Errors are for expected failures; Panic is for broken assumptions.
//     Mastering the recovery boundary allows you to build resilient
//     systems (like web servers) that can survive a crash in a single
//     request handler without bringing down the entire service.
//
// RUN:
//   go run ./03-functions-errors/10-panic-and-recover
//
// KEY TAKEAWAY:
//   - Use panic only for system-level or programmer errors.
// ============================================================================

package main

import "fmt"

// Section 03: Functions & Errors - Panic and Recover
//
// Mental model:
// Panic is for "impossible" errors. Recover is the safety net at the boundary.
//

// accessDatabase (Function): runs the access database step and keeps its inputs, outputs, or errors visible.
func accessDatabase(connected bool) {
	// defer recover must be at the very top of the function you want to protect
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
			fmt.Println("Attempting to close connections safely...")
		}
	}()

	fmt.Println("Accessing database...")

	if !connected {
		// Panic is used when the program reaches a state it cannot handle
		panic("database connection lost during operation")
	}

	fmt.Println("Database operation successful.")
}

func main() {
	fmt.Println("--- Normal Operation ---")
	accessDatabase(true)

	fmt.Println("\n--- Crisis Operation ---")
	// This call will panic, but accessDatabase will recover
	accessDatabase(false)

	fmt.Println("\nProgram continued running after the recovered panic.")

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.1 -> 04-types-design/1-struct")
	fmt.Println("Run    : go run ./04-types-design/1-struct")
	fmt.Println("Current: FE.10 (panic-and-recover)")
	fmt.Println("---------------------------------------------------")
}
