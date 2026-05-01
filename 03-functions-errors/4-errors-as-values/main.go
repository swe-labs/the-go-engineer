// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: Errors as Values
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn the core Go rule for ordinary failure: return an error value instead of hiding the failure.
//
// WHY THIS MATTERS:
//   - In Go, an error is a value returned to the caller. That means: - the function does not hide failure - the caller sees the failure directly - the ca...
//
// RUN:
//   go run ./03-functions-errors/4-errors-as-values
//
// KEY TAKEAWAY:
//   - Learn the core Go rule for ordinary failure: return an error value instead of hiding the failure.
// ============================================================================

package main

import (
	"errors"
	"fmt"
)

// 05 Functions and Errors - Errors as Values
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

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: FE.5 -> 03-functions-errors/5-validation")
	fmt.Println("Current: FE.4 (errors as values)")
	fmt.Println("Previous: FE.3 (multiple-return-values)")
	fmt.Println("---------------------------------------------------")
}
