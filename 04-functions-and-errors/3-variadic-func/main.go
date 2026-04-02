// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 4: Functions & Errors — Variadic Functions
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Variadic functions accept a variable number of arguments
//   - The ...Type syntax and how it creates a slice internally
//   - Real-world use: fmt.Println itself is variadic!
//   - Optional parameter patterns using variadics
//
// ENGINEERING DEPTH:
//   When you call a variadic function like `sum(1, 2, 3)`, the Go compiler
//   actually creates a hidden `[]int` slice behind the scenes, populates it with
//   `[1, 2, 3]`, and passes the 24-byte Slice Header into the function. This
//   is why you can iterate over `numbers` inside the function exactly like a slice
//   because, internally, it *is* a slice.
//
// RUN: go run ./04-functions-and-errors/3-variadic-func
// ============================================================================

// sum accepts zero or more int arguments.
// The "...int" syntax means "any number of ints".
//
// Inside the function, "numbers" is a []int slice — you can range over it,
// check its length, access by index, etc.
//
// REAL EXAMPLES in the standard library:
//
//	fmt.Println(a ...any)              — prints any number of values
//	append(slice, elems ...Type)       — adds any number of elements
//	errors.Join(errs ...error)         — combines multiple errors
func sum(numbers ...int) int {
	total := 0
	for _, number := range numbers {
		total += number
	}
	return total
}

// config demonstrates optional parameters using variadics.
// Since Go doesn't have default parameter values (like Python's def f(x=5)),
// the variadic pattern is a clean alternative:
//
//	config()        — use default
//	config(42)      — use provided value
//	config(1, 2, 3) — use all provided values
func config(numbers ...int) {
	if len(numbers) > 0 {
		first := numbers[0]
		fmt.Println("  Custom config:", first)
	} else {
		fmt.Println("  Default config applied")
	}
}

func main() {

	// --- CALLING VARIADIC FUNCTIONS ---
	fmt.Println("=== Variadic Functions ===")

	// You can pass any number of arguments (including zero)
	fmt.Println("  sum():        ", sum())           // 0
	fmt.Println("  sum(5):       ", sum(5))          // 5
	fmt.Println("  sum(1,2,3,4): ", sum(1, 2, 3, 4)) // 10

	fmt.Println()

	// --- SPREADING A SLICE ---
	// If you already have a slice, use ... to "spread" it into the arguments.
	// This is like JavaScript's spread operator (...array).
	nums := []int{10, 20, 30}
	fmt.Println("  sum(slice...): ", sum(nums...)) // 60

	fmt.Println()

	// --- OPTIONAL PARAMETERS PATTERN ---
	fmt.Println("=== Optional Parameters ===")
	config()   // Uses default
	config(42) // Uses provided value

	// KEY TAKEAWAY:
	// - ...Type makes a function accept any number of arguments
	// - Inside the function, the variadic parameter is a regular slice
	// - Use slice... to spread an existing slice into a variadic call
	// - This pattern replaces default/optional parameters (which Go doesn't have)
}
