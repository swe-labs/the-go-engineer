// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// 05 Functions and Errors - Order Summary (Exercise Starter)
//
// Mental model:
// Build small helpers first, then let one function coordinate the whole flow.
// Pricing rules are function values, and the discount helper is a closure that
// captures threshold and amount.
//
// Run: go run ./03-functions-errors/7-order-summary/_starter

type pricingRule func(int) int

// TODO: Implement validateOrderName(name string) error
// TODO: Implement validatePrices(prices []int) error
// TODO: Implement validateShipping(shipping int) error
// TODO: Implement sumPrices(prices []int) int
// TODO: Implement applyPricingRules(subtotal int, rules ...pricingRule) int
// TODO: Implement makeMinimumSubtotalDiscount(threshold int, amount int) pricingRule
// TODO: Implement buildSummary(name string, subtotal int, adjustedSubtotal int, shipping int) string
// TODO: Implement processOrder(name string, prices []int, shipping int, rules ...pricingRule) (string, error)

func main() {
	fmt.Println("=== Order Summary Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Build the order summary using validation helpers, pricing-rule callbacks, and a closure-based discount rule.")
	fmt.Println("Start with the validators, then add sumPrices, applyPricingRules, makeMinimumSubtotalDiscount, and processOrder.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
