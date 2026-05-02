// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: Functions Basics
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Defining logical boundaries using the `func` keyword.
//   - Improving maintainability by partitioning code into named blocks.
//   - The mechanics of program execution jumping between functions.
//
// WHY THIS MATTERS:
//   - Functions are the primary tool for complexity management. They allow
//     engineers to break large problems into small, manageable, and
//     testable units of work, facilitating parallel development and reuse.
//
// RUN:
//   go run ./03-functions-errors/1-functions-basics
//
// KEY TAKEAWAY:
//   - Functions encapsulate discrete responsibilities into named units.
// ============================================================================

package main

import "fmt"

// Section 03: Functions & Errors - Functions Basics
// Functions provide a mechanism for logic reuse and documentation by naming
// executable blocks.

// printBanner outputs the lesson header to stdout.
// printBanner (Function): outputs the lesson header to stdout.
func printBanner() {
	fmt.Println("=== Functions Basics: Named Logic Units ===")
}

// printGoal summarizes the technical objective of functional decomposition.
// printGoal (Function): summarizes the technical objective of functional decomposition.
func printGoal() {
	fmt.Println("Functions provide encapsulation and nameable boundaries for logic.")
}

// printChecklist enumerates the architectural benefits of functions.
// printChecklist (Function): enumerates the architectural benefits of functions.
func printChecklist() {
	fmt.Println("- Logic isolation for testing")
	fmt.Println("- Reduced cognitive load through naming")
	fmt.Println("- Reusable execution paths")
}

func main() {
	// 1. Procedural Execution.
	// The main function serves as the entry point and orchestrates calls
	// to specialized sub-routines.
	printBanner()
	printGoal()
	printChecklist()

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: FE.2 -> 03-functions-errors/2-parameters-and-returns")
	fmt.Println("Run    : go run ./03-functions-errors/2-parameters-and-returns")
	fmt.Println("Current: FE.1 (functions-basics)")
	fmt.Println("---------------------------------------------------")
}
