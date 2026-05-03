// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work
// Title: Memory basics - stack vs heap
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The difference between stack and heap memory.
//   - How Go's compiler performs escape analysis to choose between them.
//
// WHY THIS MATTERS:
//   - The stack is like a neat stack of plates: fast, ordered, and self-cleaning.
//     The heap is a flexible but messy shared room that requires the Garbage
//     Collector to clean up once objects are no longer used.
//
// RUN:
//   go run ./00-how-computers-work/3-memory-basics
//
// KEY TAKEAWAY:
//   - Memory behavior shapes performance. Heap allocation is flexible but comes
//     with the cost of GC management.
// ============================================================================

package main

import "fmt"

// noEscape (Function): returns a plain value whose lifetime stays inside the caller boundary.
func noEscape() int {
	x := 42
	return x
}

// escapes (Function): returns an address, forcing the local value to outlive its stack frame.
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

	// - Every goroutine has its own growable stack (starting at 2KB).
	// - Escape analysis is a compiler-time decision, not a runtime one.
	// - Reducing heap escapes is a primary strategy for Go performance tuning.
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.4 -> 00-how-computers-work/4-terminal-confidence")
	fmt.Println("Run    : go run ./00-how-computers-work/4-terminal-confidence")
	fmt.Println("Current: HC.3 (memory-basics)")
	fmt.Println("---------------------------------------------------")
}
