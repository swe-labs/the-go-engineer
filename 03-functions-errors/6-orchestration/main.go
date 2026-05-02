// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: Orchestration
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to combine multiple specialized functions into a single workflow.
//   - Delegating small responsibilities to helpers.
//   - Managing "Short-Circuit" logic to stop execution on the first error.
//
// WHY THIS MATTERS:
//   - In production, a single "Request" often triggers a chain of actions
//     (Validate -> Auth -> DB Lookup -> Calc -> Log). Orchestration is the
//     art of linking these steps without turning your code into a messy
//     "God Function" that does everything.
//
// RUN:
//   go run ./03-functions-errors/6-orchestration
//
// KEY TAKEAWAY:
//   - High-level orchestrators delegate details to low-level helpers.
// ============================================================================

package main

import (
	"errors"
	"fmt"
	"strings"
)

// Section 03: Functions & Errors - Orchestration
//
// Mental model:
// One function can coordinate several helpers and stop early when a step fails.
//

// validateCartName (Function): runs the validate cart name step and keeps its inputs, outputs, or errors visible.
func validateCartName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("cart name is required")
	}

	return nil
}

// validatePrices (Function): runs the validate prices step and keeps its inputs, outputs, or errors visible.
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

// sumPrices (Function): runs the sum prices step and keeps its inputs, outputs, or errors visible.
func sumPrices(prices []int) int {
	total := 0

	for _, price := range prices {
		total += price
	}

	return total
}

// buildSummary (Function): runs the build summary step and keeps its inputs, outputs, or errors visible.
func buildSummary(name string, total int) string {
	return fmt.Sprintf("%s total: %d", name, total)
}

// processCart (Function): runs the process cart step and keeps its inputs, outputs, or errors visible.
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

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: FE.8 -> 03-functions-errors/7-first-class-functions")
	fmt.Println("Run    : go run ./03-functions-errors/7-first-class-functions")
	fmt.Println("Current: FE.6 (orchestration)")
	fmt.Println("---------------------------------------------------")
}
