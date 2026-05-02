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
// RUN: go run ./03-functions-errors/9-order-summary/_starter

// pricingRule (Type): names the pricing rule concept so the lesson can pass it as a first-class value.
type pricingRule func(int) int

// Learner task: implement validateOrderName(name string) error.
// Learner task: implement validatePrices(prices []int) error.
// Learner task: implement validateShipping(shipping int) error.
// Learner task: implement sumPrices(prices []int) int.
// Learner task: implement applyPricingRules(subtotal int, rules ...pricingRule) int.
// Learner task: implement makeMinimumSubtotalDiscount(threshold int, amount int) pricingRule.
// Learner task: implement buildSummary(name string, subtotal int, adjustedSubtotal int, shipping int) string.
// Learner task: implement processOrder(name string, prices []int, shipping int, rules ...pricingRule) (string, error).

func main() {
	fmt.Println("=== Order Summary Exercise ===")
	fmt.Println()
	fmt.Println("Build the order summary using validation helpers,")
	fmt.Println("pricing-rule callbacks, and a closure-based discount rule.")
	fmt.Println("Start with validators, then add sumPrices, applyPricingRules,")
	fmt.Println("makeMinimumSubtotalDiscount, and processOrder.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
