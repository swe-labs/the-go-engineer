// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work — What Is a Program?
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - A program is a list of instructions
//   - The CPU runs instructions through a fetch-decode-execute loop
//   - Printed output is just the visible effect of machine instructions running
//
// WHY THIS MATTERS:
//   Engineers write better code when they can reason about what the machine is
//   doing instead of treating execution like magic.
//
// RUN: go run ./00-how-computers-work/1-what-is-a-program
// ============================================================================

package main

import "fmt"

func main() {
	fmt.Println("A program is a list of instructions for the machine.")
	fmt.Println("The CPU keeps fetching, decoding, and executing those instructions.")
	fmt.Println("Even this printed text is just the visible effect of that loop.")

	// KEY TAKEAWAY:
	// - High-level Go code eventually becomes machine instructions.
	// - The CPU executes those instructions one step at a time.
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.2 code-to-execution")
	fmt.Println("Run    : go run ./00-how-computers-work/2-code-to-execution")
	fmt.Println("Current: HC.1 (what-is-a-program)")
	fmt.Println("---------------------------------------------------")
}
