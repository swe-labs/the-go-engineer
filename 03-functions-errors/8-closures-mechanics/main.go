// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: Closures - mechanics
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How closures "capture" variables from their surrounding scope.
//   - The difference between capturing by reference vs. capturing by value.
//   - How the Go compiler moves captured variables to the heap.
//
// WHY THIS MATTERS:
//   - Closures allow you to "carry" state with a function. This is how
//     middleware keeps configuration data, how observers keep track of
//     subject state, and how factories generate specialized helpers.
//
// RUN:
//   go run ./03-functions-errors/8-closures-mechanics
//
// KEY TAKEAWAY:
//   - A closure is a function paired with its surrounding environment.
// ============================================================================

package main

import "fmt"

// Section 03: Functions & Errors - Closures Mechanics
//
// Mental model:
// A closure remembers variables from the scope where it was created.
//

// counter returns a function that holds onto its own internal state.
// counter (Function): returns a function that holds onto its own internal state.
func counter() func() int {
	// 'count' is defined outside the returned function
	count := 0

	// The returned function "closes over" the 'count' variable
	return func() int {
		count++
		return count
	}
}

func main() {
	// 1. Closures capture variables from outer scope
	message := "Hello"

	sayHello := func() {
		// sayHello remembers and can use the 'message' variable
		fmt.Printf("Message from closure: %s\n", message)
	}
	sayHello()

	// If the outer variable changes, the closure sees the change
	message = "Goodbye"
	sayHello()

	// 2. Captured state stays alive as long as the closure can still use it
	nextID := counter() // nextID gets its own 'count' variable initialized to 0

	fmt.Printf("first call: %d\n", nextID())  // 1
	fmt.Printf("second call: %d\n", nextID()) // 2

	// If we create a new counter, it gets a fresh copy of state
	anotherCounter := counter()
	fmt.Printf("another counter first call: %d\n", anotherCounter()) // 1

	// 3. Loop variables must be rebound when each closure needs its own copy
	// Version boundary: Go 1.22+ changed loop variable scoping to capture per-iteration,
	// making the classic "loop variable capture bug" far less common,
	// but it is still important to understand the boundary.
	funcs := []func(){}

	for i := range 3 {
		// explicit rebinding ensures the closure captures exactly what we expect
		val := i
		funcs = append(funcs, func() {
			fmt.Printf("closure captured: %d\n", val)
		})
	}

	for _, f := range funcs {
		f()
	}

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: FE.7 -> 03-functions-errors/9-order-summary")
	fmt.Println("Run    : go run ./03-functions-errors/9-order-summary")
	fmt.Println("Current: FE.9 (closures-mechanics)")
	fmt.Println("---------------------------------------------------")
}
