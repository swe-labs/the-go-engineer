// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: Errors as Values
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why Go returns errors as values instead of using exceptions.
//   - Using `errors.New` and `fmt.Errorf` to create error values.
//   - The `if err != nil` idiom for checking failures.
//   - Why the result is unreliable when an error is present.
//
// WHY THIS MATTERS:
//   - In many languages, a failure "blows up" the stack with an exception.
//     In Go, a failure is just data. This forces engineers to handle
//     problems explicitly, leading to more robust systems that don't
//     crash unexpectedly in production.
//
// RUN:
//   go run ./03-functions-errors/4-errors-as-values
//
// KEY TAKEAWAY:
//   - Errors are data, not magic. Handle them explicitly.
// ============================================================================

package main

import (
	"errors"
	"fmt"
)

// Section 03: Functions & Errors - Errors as Values
//
// Mental model:
// An error is a returned value that tells the caller the work did not succeed.
//

func divide(total int, parts int) (int, error) {
	if parts == 0 {
		return 0, errors.New("cannot divide by zero")
	}

	return total / parts, nil
}

func lookupPrice(catalog map[string]int, item string) (int, error) {
	price, exists := catalog[item]
	if !exists {
		return 0, fmt.Errorf("price not found for %q", item)
	}

	return price, nil
}

func main() {
	result, err := divide(12, 0)
	if err != nil {
		fmt.Println("divide error:", err)
	} else {
		fmt.Println("divide result:", result)
	}

	catalog := map[string]int{
		"bread": 40,
		"tea":   25,
	}

	teaPrice, err := lookupPrice(catalog, "tea")
	if err != nil {
		fmt.Println("lookup error:", err)
	} else {
		fmt.Println("tea price:", teaPrice)
	}

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: FE.5 -> 03-functions-errors/5-validation")
	fmt.Println("Run    : go run ./03-functions-errors/5-validation")
	fmt.Println("Current: FE.4 (errors-as-values)")
	fmt.Println("---------------------------------------------------")
}
