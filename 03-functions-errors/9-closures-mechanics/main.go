// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: Closures - mechanics
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how closures capture variables, why that extends lifetimes, and where the loop-variable trap comes from.
//
// WHY THIS MATTERS:
//   - A closure remembers variables from the scope where it was created, not just the values you expected in that moment.
//
// RUN:
//   go run ./03-functions-errors/9-closures-mechanics
//
// KEY TAKEAWAY:
//   - Learn how closures capture variables, why that extends lifetimes, and where the loop-variable trap comes from.
// ============================================================================

package main

import "fmt"

//
// Mental model:
// A closure remembers variables from the scope where it was created.
//

// counter returns a function that holds onto its own internal state.
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
	// Note: Go 1.22+ changed loop variable scoping to capture per-iteration,
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

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: FE.10 -> 03-functions-errors/10-panic-and-recover")
	fmt.Println("Current: FE.9 (closures - mechanics)")
	fmt.Println("Previous: FE.8 (first-class-functions)")
	fmt.Println("---------------------------------------------------")
}
