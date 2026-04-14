// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"errors"
	"fmt"
	"strings"
)

// 05 Functions and Errors - Order Summary (Exercise)
//
// Mental model:
// Smaller helpers validate and calculate, and one orchestration function keeps
// the whole flow readable.
//
// Run: go run ./01-foundations/05-functions-and-errors/7-order-summary

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

func buildSummary(name string, subtotal int, shipping int) string {
	total := subtotal + shipping
	return fmt.Sprintf("%s -> subtotal: %d, shipping: %d, total: %d", name, subtotal, shipping, total)
}

func processOrder(name string, prices []int, shipping int) (string, error) {
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
	return buildSummary(name, subtotal, shipping), nil
}

func main() {
	fmt.Println("=== Order Summary ===")

	summary, err := processOrder("starter cart", []int{12, 18, 25}, 10)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(summary)
	}

	summary, err = processOrder(" ", []int{12, 18, 25}, 10)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(summary)
	}

	summary, err = processOrder("starter cart", []int{12, -3, 25}, 10)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(summary)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("SECTION COMPLETE: FE.7 order-summary")
	fmt.Println("NEXT UP: 02-engineering-core")
	fmt.Println("---------------------------------------------------")
}
