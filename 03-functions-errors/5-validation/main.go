// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: Validation
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how a function rejects bad input before the program does the real work.
//
// WHY THIS MATTERS:
//   - Validation is the first gate before useful work begins. If the input is clearly wrong, the function should say so immediately instead of pretending...
//
// RUN:
//   go run ./03-functions-errors/5-validation
//
// KEY TAKEAWAY:
//   - Learn how a function rejects bad input before the program does the real work.
// ============================================================================

package main

import (
	"errors"
	"fmt"
	"strings"
)

// 05 Functions and Errors - Validation
//
// Mental model:
// Validation rejects bad input early so the rest of the program can stay honest.
//

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

func main() {
	err := validateCartName("starter cart")
	fmt.Println("name check:", err)

	err = validateCartName("   ")
	fmt.Println("blank name check:", err)

	err = validatePrices([]int{12, 18, 25})
	fmt.Println("price check:", err)

	err = validatePrices([]int{12, -3, 25})
	fmt.Println("negative price check:", err)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: FE.6 -> 03-functions-errors/6-orchestration")
	fmt.Println("Current: FE.5 (validation)")
	fmt.Println("Previous: FE.4 (errors-as-values)")
	fmt.Println("---------------------------------------------------")
}
