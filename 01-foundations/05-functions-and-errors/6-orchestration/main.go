// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"errors"
	"fmt"
	"strings"
)

// 05 Functions and Errors - Orchestration
//
// Mental model:
// One function can coordinate several helpers and stop early when a step fails.
//
// Run: go run ./01-foundations/05-functions-and-errors/6-orchestration

func validateCartName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("cart name is required")
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

func sumPrices(prices []int) int {
	total := 0

	for _, price := range prices {
		total += price
	}

	return total
}

func buildSummary(name string, total int) string {
	return fmt.Sprintf("%s total: %d", name, total)
}

func processCart(name string, prices []int) (string, error) {
	if err := validateCartName(name); err != nil {
		return "", err
	}

	if err := validatePrices(prices); err != nil {
		return "", err
	}

	total := sumPrices(prices)
	summary := buildSummary(name, total)

	return summary, nil
}

func main() {
	summary, err := processCart("starter cart", []int{12, 18, 25})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(summary)
	}

	summary, err = processCart("", []int{12, 18, 25})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(summary)
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: FE.7 order-summary")
	fmt.Println("Current: FE.6 (orchestration)")
	fmt.Println("---------------------------------------------------")
}
