// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: Parameters and Returns
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Declaring function parameters (input) and return types (output).
//   - How `int` values are passed by copy.
//   - How slice headers are copied into function boundaries.
//
// WHY THIS MATTERS:
//   - Functions are more than just named blocks; they are data processors.
//     Understanding how data moves in and out of a function is the key to
//     building reusable logic that doesn't rely on global state.
//
// RUN:
//   go run ./03-functions-errors/2-parameters-and-returns
//
// KEY TAKEAWAY:
//   - Parameters = Input; Return values = Output.
// ============================================================================

package main

import "fmt"

// Section 03: Functions & Errors - Parameters and Returns
//
// Mental model:
// Parameters bring values into a function, and return values send results back.
//

// announceCart prints a notification for the current processing target.
func announceCart(name string) {
	fmt.Printf("Processing Cart: %s\n", name)
}

// sumPrices calculates the arithmetic sum of a slice of integers.
func sumPrices(prices []int) int {
	total := 0
	for _, price := range prices {
		total += price
	}
	return total
}

// labelPrice formats a descriptive string combining a label and a total value.
func labelPrice(name string, total int) string {
	return fmt.Sprintf("Summary -> %s: $%d", name, total)
}

func main() {
	// 1. Parameter pass-by-copy.
	// When we pass 'starter cart' or the 'prices' slice, Go copies the
	// value (or the slice header) into the function's local scope.
	prices := []int{12, 18, 25}
	announceCart("starter cart")

	// 2. Data processing and returns.
	// Functions transform inputs into outputs, keeping the caller's
	// scope clean of intermediate calculation state.
	total := sumPrices(prices)
	summary := labelPrice("starter cart", total)

	fmt.Println(summary)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: FE.3 -> 03-functions-errors/3-multiple-return-values")
	fmt.Println("Run    : go run ./03-functions-errors/3-multiple-return-values")
	fmt.Println("Current: FE.2 (parameters-and-returns)")
	fmt.Println("---------------------------------------------------")
}
