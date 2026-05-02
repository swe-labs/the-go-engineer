// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: Validation
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Implementing "Guard Clauses" to reject bad data early.
//   - Using `nil` as the success signal for validation functions.
//   - Returning descriptive error messages with index context.
//
// WHY THIS MATTERS:
//   - Software engineering is about building "Defensive" code. By validating
//     inputs at the boundary, you prevent bad data from "poisoning" your
//     database or causing crashes in the core logic of your application.
//
// RUN:
//   go run ./03-functions-errors/5-validation
//
// KEY TAKEAWAY:
//   - Validate at the boundary; fail fast and explicitly.
// ============================================================================

package main

import (
	"errors"
	"fmt"
	"strings"
)

// Section 03: Functions & Errors - Validation
//
// Mental model:
// Validation rejects bad input early so the rest of the program can stay honest.
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

func main() {
	err := validateCartName("starter cart")
	fmt.Println("name check:", err)

	err = validateCartName("   ")
	fmt.Println("blank name check:", err)

	err = validatePrices([]int{12, 18, 25})
	fmt.Println("price check:", err)

	err = validatePrices([]int{12, -3, 25})
	fmt.Println("negative price check:", err)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: FE.6 -> 03-functions-errors/6-orchestration")
	fmt.Println("Run    : go run ./03-functions-errors/6-orchestration")
	fmt.Println("Current: FE.5 (validation)")
	fmt.Println("---------------------------------------------------")
}
