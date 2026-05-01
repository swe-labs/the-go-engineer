// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work
// Title: Memory basics - stack vs heap
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Understand how memory is allocated and managed during execution, especially the difference between stack and heap memory in Go.
//
// WHY THIS MATTERS:
//   - The stack is like a neat stack of plates: fast, ordered, and easy to clean up. The heap is like a shared storage room: flexible and useful, but it ...
//
// RUN:
//   go run ./00-how-computers-work/3-memory-basics
//
// KEY TAKEAWAY:
//   - Understand how memory is allocated and managed during execution, especially the difference between stack and heap memory in Go.
// ============================================================================

//   Memory behavior shapes performance, correctness, and the kind of bugs a
//   Go program can have under load.

package main

import "fmt"

func noEscape() int {
	x := 42
	return x
}

func escapes() *int {
	x := 99
	return &x
}

func main() {
	value := noEscape()
	pointer := escapes()

	fmt.Printf("Value-returning function result: %d\n", value)
	fmt.Printf("Pointer-returning function result: %d (address %p)\n", *pointer, pointer)
	fmt.Println("Returning a pointer means the pointed-to value must outlive the function call.")

	// - Stack and heap serve different lifetime needs.
	// - Escape analysis helps Go decide which one a value needs.
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.4 -> 00-how-computers-work/4-terminal-confidence")
	fmt.Println("Run    : go run ./00-how-computers-work/4-terminal-confidence")
	fmt.Println("Current: HC.3 (memory-basics)")
	fmt.Println("---------------------------------------------------")
}
