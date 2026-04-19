// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics — Lesson CF.7: Pricing Checkout
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Integrating switch and if/else logic
//   - Using for loops to process collections
//   - Managing state within a loop
//
// WHY THIS MATTERS:
//   This exercise combines everything learned in the Control Flow subsystem.
//
// RUN: go run ./02-language-basics/03-control-flow/7-pricing-checkout
// ============================================================================

package main

import "fmt"

func main() {
	cart := []string{"TSHIRT", "MUG", "HAT", "BOOK", "KEYBOARD"}

	var subtotal float64

	fmt.Println("Processing checkout:")

	for _, item := range cart {
		var price float64

		switch item {
		case "TSHIRT":
			price = 20.00
		case "MUG":
			price = 12.50
		case "HAT":
			price = 18.00
		case "BOOK":
			price = 25.99
		}

		if price == 0 {
			fmt.Printf("skip %s: unknown item\n", item)
			continue
		}

		if item == "BOOK" {
			originalPrice := price
			price = price * 0.90
			fmt.Printf("%s promo: %.2f -> %.2f\n", item, originalPrice, price)
		} else {
			fmt.Printf("%s: %.2f\n", item, price)
		}

		subtotal += price
	}

	fmt.Printf("subtotal: %.2f\n", subtotal)

	// KEY TAKEAWAY:
	// - Use switch for discrete value matching.
	// - Use if/else for conditional logic and ranges.
	// - Use continue to skip invalid or unhandled items in a loop.

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: DS.1 arrays")
	fmt.Println("Current: CF.7 (pricing-checkout)")
	fmt.Println("---------------------------------------------------")
}
