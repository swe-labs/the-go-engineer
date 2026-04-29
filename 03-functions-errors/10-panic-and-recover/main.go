// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: panic and recover
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn when panic is appropriate, when it is not, and how recover turns a crash into an explicit boundary decision.
//
// WHY THIS MATTERS:
//   - Errors describe expected failure. Panic describes broken assumptions. Recover belongs at process or request boundaries, not in ordinary business flow.
//
// RUN:
//   go run ./03-functions-errors/10-panic-and-recover
//
// KEY TAKEAWAY:
//   - Learn when panic is appropriate, when it is not, and how recover turns a crash into an explicit boundary decision.
// ============================================================================

package main

import "fmt"

//
// Mental model:
// Panic is for "impossible" errors. Recover is the safety net at the boundary.
//

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

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.1 structs")
	fmt.Println("Current: FE.10 (panic and recover)")
	fmt.Println("Previous: FE.7 (order-summary)")
	fmt.Println("---------------------------------------------------")
}
