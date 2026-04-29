// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: Multiple Return Values
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how one function can return more than one value and why that matters before errors enter the picture.
//
// WHY THIS MATTERS:
//   - Sometimes one result is not enough. A function may need to return: - a value and a success signal - two related values - a result and some extra co...
//
// RUN:
//   go run ./03-functions-errors/3-multiple-return-values
//
// KEY TAKEAWAY:
//   - Learn how one function can return more than one value and why that matters before errors enter the picture.
// ============================================================================

package main

import (
	"fmt"
	"strings"
)

// 05 Functions and Errors - Multiple Return Values
//
// Mental model:
// A function can return more than one value when one result is not enough.
//

func findItem(items []string, target string) (int, bool) {
	for i, item := range items {
		if item == target {
			return i, true
		}
	}

	return -1, false
}

func splitName(fullName string) (string, string) {
	parts := strings.SplitN(fullName, " ", 2)
	if len(parts) < 2 {
		return fullName, ""
	}

	return parts[0], parts[1]
}

func main() {
	items := []string{"bread", "tea", "rice"}

	index, found := findItem(items, "tea")
	fmt.Printf("tea found? %t at index %d\n", found, index)

	firstName, lastName := splitName("Ava Stone")
	fmt.Printf("first=%s last=%s\n", firstName, lastName)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: FE.4 errors-as-values")
	fmt.Println("Current: FE.3 (multiple return values)")
	fmt.Println("Previous: FE.2 (parameters-and-returns)")
	fmt.Println("---------------------------------------------------")
}
