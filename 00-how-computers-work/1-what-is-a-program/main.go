// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work
// Title: What is a program?
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Understand that a program is a list of instructions for a machine to follow, and that the CPU runs a continuous fetch-decode-execute loop to carry ...
//
// WHY THIS MATTERS:
//   - Imagine a cook reading a recipe card. - the recipe card is the program - the ingredients are the data - the cook is the CPU The cook does not under...
//
// RUN:
//   go run ./00-how-computers-work/1-what-is-a-program
//
// KEY TAKEAWAY:
//   - Understand that a program is a list of instructions for a machine to follow, and that the CPU runs a continuous fetch-decode-execute loop to carry ...
// ============================================================================

//   Engineers write better code when they can reason about what the machine is
//   doing instead of treating execution like magic.

package main

import "fmt"

func main() {
	fmt.Println("A program is a list of instructions for the machine.")
	fmt.Println("The CPU keeps fetching, decoding, and executing those instructions.")
	fmt.Println("Even this printed text is just the visible effect of that loop.")

	// - High-level Go code eventually becomes machine instructions.
	// - The CPU executes those instructions one step at a time.
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.2 -> 00-how-computers-work/2-code-to-execution")
	fmt.Println("Run    : go run ./00-how-computers-work/2-code-to-execution")
	fmt.Println("Current: HC.1 (what-is-a-program)")
	fmt.Println("---------------------------------------------------")
}
