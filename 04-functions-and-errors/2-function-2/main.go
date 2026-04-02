// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 4: Functions & Errors — Recursion & Closures
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Recursive functions (functions that call themselves)
//   - Closures: anonymous functions that capture surrounding state
//   - How closures "close over" variables (they share, not copy)
//   - First-class functions: functions as variables
//
// ENGINEERING DEPTH:
//   When a closure captures a local variable (like `i` in `intSeq`), the Go compiler
//   performs "Escape Analysis". Because the returned function outlives the
//   scope of the parent function, `i` can no longer safely live on the Stack
//   (which gets destroyed when `intSeq` returns). The compiler automatically
//   allocates `i` on the Heap so it survives. This is powerful, but Heap
//   allocations trigger Garbage Collection, making closures microscopically
//   slower than standard functions.
//
// RUN: go run ./04-functions-and-errors/2-function-2
// ============================================================================

// factorial computes n! using recursion.
// Every recursive function MUST have a BASE CASE to stop the recursion.
// Without it, the function calls itself forever → stack overflow crash.
//
// How it works:
//
//	factorial(5) → 5 * factorial(4)
//	                  → 4 * factorial(3)
//	                       → 3 * factorial(2)
//	                            → 2 * factorial(1)
//	                                 → 1 (base case!)
//	Result: 5 * 4 * 3 * 2 * 1 = 120
func factorial(n int) int {
	if n <= 1 {
		return 1 // BASE CASE: stop recursion here
	}
	return n * factorial(n-1) // RECURSIVE CASE: call ourselves with n-1
}

// intSeq demonstrates a CLOSURE — a function that captures external state.
//
// intSeq returns a FUNCTION (not a value). Each call to the returned function
// increments the captured variable "i" and returns its new value.
//
// The variable "i" lives in intSeq's scope, but the returned function
// "closes over" it — maintaining access even after intSeq returns.
// Each call to intSeq() creates a NEW, independent counter.
func intSeq() func() int {
	i := 0 // This variable is captured by the closure below
	return func() int {
		i++ // Modifies the SAME "i" on every call (shared reference)
		return i
	}
}

func main() {

	// --- RECURSION ---
	fmt.Println("=== Recursion ===")
	fmt.Printf("  5! = %d\n", factorial(5)) // 120
	fmt.Println()

	// --- CLOSURES ---
	fmt.Println("=== Closures ===")

	// intSeq() returns a function. Each invocation of the returned function
	// increments and returns the SAME counter.
	nextInt := intSeq()
	fmt.Println("  Counter 1:", nextInt()) // 1
	fmt.Println("  Counter 1:", nextInt()) // 2
	fmt.Println("  Counter 1:", nextInt()) // 3
	fmt.Println("  Counter 1:", nextInt()) // 4

	// Creating a NEW closure gives us an INDEPENDENT counter.
	// The two counters don't interfere with each other.
	anotherCounter := intSeq()
	fmt.Println("  Counter 2:", anotherCounter()) // 1 (independent!)
	fmt.Println()

	// --- ANONYMOUS FUNCTIONS ---
	// Functions in Go are "first class" — they can be:
	//   1. Assigned to variables
	//   2. Passed as arguments to other functions
	//   3. Returned from other functions
	//
	// An anonymous function is a function without a name.
	// Useful for one-off operations, callbacks, and goroutines.
	logger := func(msg string) {
		fmt.Printf("  [LOG] %s\n", msg)
	}

	logger("Hello World")
	logger("Application started")

	// KEY TAKEAWAY:
	// - Recursive functions: always define a base case to prevent infinite loops
	// - Closures capture variables by REFERENCE, not by copy
	// - Each closure instance maintains its own state
	// - Functions are first-class: store them in variables, pass them around
	// - Anonymous functions are essential for goroutines (Section 09)
}
