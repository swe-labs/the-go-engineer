// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 2: Control Flow — For Loops
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Go has ONLY ONE loop keyword: "for" (no while, no do-while)
//   - 4 forms of the for loop: C-style, while-style, infinite, range
//   - break and continue flow control
//   - The range keyword for iterating over collections
//
// ENGINEERING DEPTH:
//   Go controversially removed the `while` and `do-while` keywords that exist in
//   almost every other C-family language. Why? Because a `while` loop is simply
//   a `for` loop with the initialization and post-statements omitted. By forcing
//   developers to use a single reserved keyword (`for`), the compiler's parsing
//   engine is drastically simplified, and codebases remain visually consistent
//   across millions of lines of code.
//
// RUN: go run ./02-control-flow/1-for-loop
// ============================================================================

func main() {

	// --- FORM 1: C-Style for loop ---
	// Syntax: for init; condition; post { body }
	//
	// - init: runs once before the loop starts (i := 1)
	// - condition: checked before each iteration (i <= 10)
	// - post: runs after each iteration (i++)
	//
	// The variable "i" is scoped to the loop — it doesn't exist outside.
	fmt.Println("=== C-Style Loop ===")
	for i := 1; i <= 5; i++ {
		fmt.Printf("  i = %d\n", i)
	}
	// fmt.Println(i) ← COMPILE ERROR: i is not defined here

	// --- FORM 2: While-style loop ---
	// Go doesn't have a "while" keyword. Instead, use for with only a condition.
	// Syntax: for condition { body }
	//
	// This is equivalent to "while (k > 0)" in C or Python.
	fmt.Println("\n=== While-Style Loop ===")
	k := 3
	for k > 0 {
		fmt.Printf("  k = %d\n", k)
		k-- // Decrement k by 1 each iteration
	}

	// --- FORM 3: Infinite loop ---
	// Syntax: for { body }
	//
	// This runs FOREVER until you explicitly "break" out.
	// Common uses: servers waiting for connections, event loops, retry logic.
	fmt.Println("\n=== Infinite Loop (with break) ===")
	counter := 0
	for {
		fmt.Printf("  counter = %d\n", counter)
		counter++
		if counter >= 5 {
			break // Exit the loop when counter reaches 5
		}
	}

	// --- CONTINUE: Skip the current iteration ---
	// "continue" jumps to the next iteration, skipping the remaining body.
	// Here, we skip even numbers and only print odd numbers.
	fmt.Println("\n=== Continue (odd numbers only) ===")
	for i := 1; i <= 10; i++ {
		if i%2 == 0 { // i%2 is the remainder of i÷2. Even numbers have remainder 0.
			continue // Skip to the next iteration
		}
		fmt.Printf("  %d\n", i)
	}

	// --- FORM 4: Range loop ---
	// The "range" keyword iterates over collections (arrays, slices, maps, strings).
	// It returns TWO values: the index (position) and the value at that position.
	//
	// If you don't need the index, use _ (blank identifier) to discard it.
	// If you don't need the value, just omit it: for i := range items
	fmt.Println("\n=== Range Loop ===")
	items := [3]string{"Go", "Python", "Java"}

	// Using both index and value
	for index, value := range items {
		fmt.Printf("  items[%d] = %s\n", index, value)
	}

	// Using only the value (discard index with _)
	fmt.Println("\n  Values only:")
	for _, value := range items {
		fmt.Printf("  %s\n", value)
	}

	// KEY TAKEAWAY:
	// Go has ONE loop keyword (for) with FOUR forms:
	//   1. for i := 0; i < n; i++ { }     — C-style
	//   2. for condition { }                — while-style
	//   3. for { }                          — infinite loop
	//   4. for index, value := range x { }  — iterate collections
}
