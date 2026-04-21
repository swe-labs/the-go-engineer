// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"errors"
	"fmt"
	"strings"
)

// 05 Functions and Errors - Order Summary (Exercise)
//
// Mental model:
// Smaller helpers validate and calculate, one orchestration function owns the
// sequence, and pricing rules are passed in as function values.
//
// Run: go run ./03-functions-errors/7-order-summary

type pricingRule func(int) int

func validateOrderName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("order name is required")
	}

	return nil
}

func validatePrices(prices []int) error {
	if len(prices) == 0 {
		return errors.New("at least one price is required")
	}

	for i, price := range prices {
		if price < 0 {
			return fmt.Errorf("price at index %d cannot be negative", i)
		}
	}

	return nil
}

func validateShipping(shipping int) error {
	if shipping < 0 {
		return errors.New("shipping cannot be negative")
	}

	return nil
}

func sumPrices(prices []int) int {
	total := 0

	for _, price := range prices {
		total += price
	}

	return total
}

func applyPricingRules(subtotal int, rules ...pricingRule) int {
	adjusted := subtotal

	for _, rule := range rules {
		adjusted = rule(adjusted)
		if adjusted < 0 {
			adjusted = 0
		}
	}

	return adjusted
}

func makeMinimumSubtotalDiscount(threshold int, amount int) pricingRule {
	return func(subtotal int) int {
		if subtotal < threshold {
			return subtotal
		}

		adjusted := subtotal - amount
		if adjusted < 0 {
			return 0
		}

		return adjusted
	}
}

func buildSummary(name string, subtotal int, adjustedSubtotal int, shipping int) string {
	total := adjustedSubtotal + shipping
	return fmt.Sprintf(
		"%s -> subtotal: %d, adjusted subtotal: %d, shipping: %d, total: %d",
		name,
		subtotal,
		adjustedSubtotal,
		shipping,
		total,
	)
}

func processOrder(name string, prices []int, shipping int, rules ...pricingRule) (string, error) {
	if err := validateOrderName(name); err != nil {
		return "", err
	}

	if err := validatePrices(prices); err != nil {
		return "", err
	}

	if err := validateShipping(shipping); err != nil {
		return "", err
	}

	subtotal := sumPrices(prices)
	adjustedSubtotal := applyPricingRules(subtotal, rules...)

	return buildSummary(name, subtotal, adjustedSubtotal, shipping), nil
}

func main() {
	fmt.Println("=== Order Summary ===")

	vipDiscount := makeMinimumSubtotalDiscount(50, 5)

	summary, err := processOrder("starter cart", []int{12, 18, 25}, 10, vipDiscount)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(summary)
	}

	summary, err = processOrder("small cart", []int{12, 8}, 5, vipDiscount)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(summary)
	}

	summary, err = processOrder(" ", []int{12, 18, 25}, 10, vipDiscount)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(summary)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: FE.10")
	fmt.Println("Current: FE.7 (order summary)")
	fmt.Println("---------------------------------------------------")
}
