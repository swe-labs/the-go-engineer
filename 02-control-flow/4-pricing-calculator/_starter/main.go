// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 2: Control Flow — Pricing Calculator (Exercise Starter)
// Level: Beginner
// ============================================================================
//
// EXERCISE: Build a Sales Order Processor
//
// REQUIREMENTS:
//  1. [ ] Create a package-level map `productPrices` with product codes and prices
//  2. [ ] Implement `calculateItemPrice(itemCode string) (float64, bool)` that:
//         - Looks up the item in the map
//         - If it ends with "_SALE", strip the suffix and apply a 10% discount
//         - Returns (price, found)
//  3. [ ] In main(), iterate over a slice of order items and calculate the subtotal
//
// HINTS:
//   - Use the comma-ok pattern: price, found := productPrices[code]
//   - strings.HasSuffix checks if a string ends with a suffix
//   - strings.TrimSuffix removes a suffix from a string
//
// RUN: go run ./02-control-flow/4-pricing-calculator/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Define your productPrices map here

// TODO: Implement calculateItemPrice function here

func main() {
	fmt.Println("=== Sales Order Processor Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your pricing calculator!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
