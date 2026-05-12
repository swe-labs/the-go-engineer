// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work
// Title: What is a program?
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Understand that a program is a list of instructions for a machine to follow.
//   - The CPU runs a continuous fetch-decode-execute loop to carry those instructions out.
//
// WHY THIS MATTERS:
//   - Imagine a cook reading a recipe card. The recipe is the program, ingredients
//     are data, and the cook is the CPU. The cook follows one step at a time
//     without knowing your ultimate intent.
//
// RUN:
//   go run ./00-how-computers-work/1-what-is-a-program
//
// KEY TAKEAWAY:
//   - A program is passive data until the CPU fetches it and turns it into action.
// ============================================================================

package main

import "fmt"

// main (Function): entry point for the program. It prints the core lesson
// concept to stdout so the learner can see the fetch-decode-execute loop in action.
func main() {
	fmt.Println("A program is a list of instructions for the machine.")
	fmt.Println("The CPU keeps fetching, decoding, and executing those instructions.")
	fmt.Println("Even this printed text is just the visible effect of that loop.")

	// - High-level Go code eventually becomes numeric OpCodes (Machine Instructions).
	// - The CPU's Instruction Pointer (IP) fetches these OpCodes one by one.
	// - The CPU executes those instructions in a mechanical loop.
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.2 -> 00-how-computers-work/2-code-to-execution")
	fmt.Println("Run    : go run ./00-how-computers-work/2-code-to-execution")
	fmt.Println("Current: HC.1 (what-is-a-program)")
	fmt.Println("---------------------------------------------------")
}
