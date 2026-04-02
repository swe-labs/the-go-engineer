// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
)

// ============================================================================
// Section 4: Functions & Errors — Basic Functions
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to define a function with parameters and return types
//   - Parameter type grouping (e.g. `a, b int`)
//   - The fundamental structure of Go programs
//
// ENGINEERING DEPTH:
//   Go manages a small initial contiguous Stack (usually 2KB) per goroutine.
//   When a function is called, variables are pushed onto this Stack.
//   All arguments in Go are passed by value (copied). If you pass a 100MB
//   struct to a function, Go literally copies 100MB of RAM to the stack!
//   This is why we use pointers (passing an 8-byte memory address) for large objects.
//
// RUN: go run ./04-functions-and-errors/1-function
// ============================================================================

// greet takes a single string parameter and returns nothing.
// Notice that the type (`string`) comes AFTER the variable name (`name`).
func greet(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// add demonstrates type grouping.
// Instead of writing `a int, b int`, Go lets you group identical types: `a, b int`.
func add(a, b int) {
	fmt.Printf("%d + %d = %d\n", a, b, a+b)
}

// calculateArea takes parameters and explicitly declares a return type (`float64`).
// Every code path in this function MUST return a float64.
func calculateArea(width float64, height float64) float64 {

	// Defensive programming: guard against invalid input.
	if width < 0 || height < 0 {
		fmt.Println("Error: width and height must be positive")
		// We return 0.0 because the function signature demands a float64 return.
		// In later lessons, we will learn how to return actual Errors instead of 0.
		return 0.0
	}

	// The return keyword hands the data back to the caller.
	return width * height
}

func main() {
	fmt.Println("=== Basic Functions ===")

	// 1. Invoking a simple function
	greet("Bob Wonderland")

	// 2. Invoking with grouped parameters
	add(1, 2)

	// 3. Capturing a return value
	// We assign the returned float64 to the variable `area`.
	area := calculateArea(4.0, 4.0)
	fmt.Printf("Calculated Area: %.2f\n", area)

	// KEY TAKEAWAY:
	// - Type declarations always come AFTER the variable name.
	// - If consecutive parameters share a type, you can omit the type until the last one.
	// - Functions are pass-by-value (copies are made).
}
