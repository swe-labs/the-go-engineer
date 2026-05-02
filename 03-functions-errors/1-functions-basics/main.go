// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: Functions Basics
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Defining function boundaries using the `func` keyword.
//   - Naming blocks of code to improve readability.
//   - How the program execution jumps between `main` and other functions.
//
// WHY THIS MATTERS:
//   - Real-world programs are too complex to stay inside a single `main()`
//     function. Functions are the primary tool Go engineers use to break
//     large problems into small, manageable, and testable pieces.
//
// RUN:
//   go run ./03-functions-errors/1-functions-basics
//
// KEY TAKEAWAY:
//   - A function encapsulates one recognizable responsibility.
// ============================================================================

package main

import "fmt"

// Section 03: Functions & Errors - Functions Basics
//
// Mental model:
// A function gives one piece of work a name so main() can stay readable.
//

func printBanner() {
	fmt.Println("=== Functions Basics ===")
}

func printGoal() {
	fmt.Println("A function gives a piece of work a name.")
}

func printChecklist() {
	fmt.Println("- main() can call other functions")
	fmt.Println("- each function can do one small job")
	fmt.Println("- named steps are easier to read than one long inline block")
}

func main() {
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
