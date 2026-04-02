// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// Section 2: Control Flow — Pricing Calculator (Exercise)
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Combining maps, loops, if/else, and functions
//   - Real-world data processing: looking up prices, applying discounts
//   - The comma-ok pattern for safe map access
//   - String manipulation with the strings package
//
// ENGINEERING DEPTH:
//   Notice how we strip suffixes using `strings.TrimSuffix()`. In Go, strings are
//   immutable read-only byte slices. When you "modify" a string, Go actually
//   allocates a brand new block of heap memory and copies the bytes over.
//   For high-performance systems handling millions of logs, string allocations
//   become a massive CPU bottleneck, which is why byte slices (`[]byte`) are
//   preferred in hot paths.
//
// RUN: go run ./02-control-flow/4-pricing-calculator
// ============================================================================

// productPrices is a package-level map acting as a simple "database" of products.
// In a real app, this would come from a database or API.
// Map type: map[string]float64 — product code (string) maps to price (float64).
var productPrices = map[string]float64{
	"TSHIRT": 20.00,
	"MUG":    12.50,
	"HAT":    18.00,
	"BOOK":   25.99,
}

// calculateItemPrice looks up a product's price and handles sale items.
//
// RETURN VALUES: (float64, bool)
//   - float64: the calculated price (0.0 if not found)
//   - bool: whether the product was found
//
// This (value, ok) pattern mirrors how Go maps work internally,
// giving the caller control over how to handle missing products.
func calculateItemPrice(itemCode string) (float64, bool) {
	// First, try a direct lookup in the product map.
	basePrice, found := productPrices[itemCode]
	if !found {
		// Not found directly — check if it's a sale item (ends with "_SALE").
		// strings.HasSuffix checks if a string ends with a given suffix.
		if strings.HasSuffix(itemCode, "_SALE") {
			// Strip the "_SALE" suffix to find the original product code.
			originalItemCode := strings.TrimSuffix(itemCode, "_SALE")
			basePrice, found = productPrices[originalItemCode]
			if found {
				// Apply a 10% discount for sale items.
				salePrice := basePrice * 0.90
				fmt.Printf("  📦 %s (Sale! $%.2f → $%.2f)\n",
					originalItemCode, basePrice, salePrice)
				return salePrice, true
			}
		}

		// Product not found at all — neither regular nor sale.
		fmt.Printf("  ❌ %s (not found)\n", itemCode)
		return 0.0, false
	}

	// Regular (non-sale) product found.
	fmt.Printf("  📦 %s — $%.2f\n", itemCode, basePrice)
	return basePrice, true
}

func main() {
	fmt.Println("=== Sales Order Processor ===")
	fmt.Println()

	// Simulate a customer order with a mix of regular and sale items.
	// Notice "MUG_SALE" — our function will detect the suffix and apply 10% off.
	orderItems := []string{
		"TSHIRT", "MUG_SALE", "HAT", "BOOK", "KEYBOARD",
	}

	var subtotal float64
	fmt.Println("Processing order:")

	// Range over the order items and calculate each price.
	for _, item := range orderItems {
		price, found := calculateItemPrice(item)
		if found {
			subtotal += price
		}
	}

	fmt.Println()
	fmt.Printf("Subtotal: $%.2f\n", subtotal)

	// KEY TAKEAWAY:
	// This exercise combines everything from Section 02:
	//   - for/range loop to iterate the order
	//   - if/else branching for sale detection
	//   - map lookups with comma-ok pattern
	//   - Functions with multiple return values
	//   - String operations (HasSuffix, TrimSuffix)
}
