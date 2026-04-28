// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: First-class functions
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn that functions are ordinary values in Go, which makes callbacks and higher-order helpers possible.
//
// WHY THIS MATTERS:
//   - A function value is just another tool you can store, pass, and call later.
//
// RUN:
//   go run ./03-functions-errors/8-first-class-functions
//
// KEY TAKEAWAY:
//   - Learn that functions are ordinary values in Go, which makes callbacks and higher-order helpers possible.
// ============================================================================

package main

import "fmt"

//
// Mental model:
// A function is just a value. You can assign it to a variable or pass it to another function.
//

// 1. Assigning functions to variables
func add(a, b int) int {
	return a + b
}

func multiply(a, b int) int {
	return a * b
}

// 2. Passing behavior into other functions (callbacks)
// The parameter 'operation' is a function that takes two ints and returns one int.
func calculate(a int, b int, operation func(int, int) int) int {
	return operation(a, b)
}

func main() {
	// A function can be assigned to a variable without calling it (no parentheses)
	var mathFunc func(int, int) int
	
	mathFunc = add
	result1 := mathFunc(5, 3)
	fmt.Printf("Using add: 5 + 3 = %d\n", result1)

	mathFunc = multiply
	result2 := mathFunc(5, 3)
	fmt.Printf("Using multiply: 5 * 3 = %d\n", result2)

	// Passing functions as arguments (callbacks)
	fmt.Printf("calculate with add: %d\n", calculate(10, 4, add))
	fmt.Printf("calculate with multiply: %d\n", calculate(10, 4, multiply))

	// You can also define an anonymous function inline
	subtract := func(a, b int) int {
		return a - b
	}
	fmt.Printf("calculate with anonymous subtract: %d\n", calculate(10, 4, subtract))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: FE.9 closures-mechanics")
	fmt.Println("Current: FE.8 (first-class functions)")
	fmt.Println("Previous: FE.6 (orchestration)")
	fmt.Println("---------------------------------------------------")
}
