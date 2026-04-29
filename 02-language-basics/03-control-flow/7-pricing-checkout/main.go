// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Pricing Checkout
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Build a small checkout flow that combines branching, loops, `switch`, and `continue` into one runnable program.
//
// WHY THIS MATTERS:
//   - This milestone is a miniature rule engine: - loop over each cart item - classify the item with `switch` - apply extra rules with `if` - skip bad da...
//
// RUN:
//   go run ./02-language-basics/03-control-flow/7-pricing-checkout
//
// KEY TAKEAWAY:
//   - Control flow constructs (loops, switches, conditions, and breaks) compose
//     together naturally to build resilient business logic, like this miniature
//     shopping cart rule engine.
// ============================================================================

package main

import "fmt"

func main() {
	cart := []string{"TSHIRT", "MUG", "HAT", "BOOK", "KEYBOARD"}

	var subtotal float64

	fmt.Println("Processing checkout:")

	// Loop through each item in the slice
	for _, item := range cart {
		var price float64

		// Classify the item using 'switch'
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

		// If the price wasn't updated, the item is unknown.
		// We use 'continue' to skip the rest of this loop iteration safely.
		if price == 0 {
			fmt.Printf("skip %s: unknown item\n", item)
			continue
		}

		// Apply a conditional discount specifically for books.
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

	// - Use switch for discrete value matching.
	// - Use if/else for conditional logic and ranges.
	// - Use continue to skip invalid or unhandled items in a loop.

	// Forward reference:
	// We've been using slices (like 'cart = []string{...}') without explaining
	// them fully. Next, we enter the Data Structures section to learn exactly
	// how Arrays and Slices work in memory.
	// See: ../../04-data-structures/1-array/README.md
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: DS.1 arrays")
	fmt.Println("Current: CF.7 (pricing-checkout)")
	fmt.Println("Previous: CF.6 (defer-use-cases)")
	fmt.Println("---------------------------------------------------")
}
