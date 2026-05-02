// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 03: Functions and Errors
// Title: Multiple Return Values
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Defining multiple return types in parentheses `(int, bool)`.
//   - Returning multiple values separated by commas.
//   - How the caller must receive all returned values (or use the blank identifier).
//
// WHY THIS MATTERS:
//   - Go doesn't use "sentinel" values (like returning -1 for error) because
//     that can be ambiguous. Multiple return values allow a function to be
//     honest: "Here is the data, AND here is a flag saying if it's valid."
//     This is the foundation of Go's explicit error handling.
//
// RUN:
//   go run ./03-functions-errors/3-multiple-return-values
//
// KEY TAKEAWAY:
//   - Multiple returns replace ambiguous sentinel values.
// ============================================================================

package main

import (
	"fmt"
	"strings"
)

// Section 03: Functions & Errors - Multiple Return Values
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

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: FE.4 -> 03-functions-errors/4-errors-as-values")
	fmt.Println("Run    : go run ./03-functions-errors/4-errors-as-values")
	fmt.Println("Current: FE.3 (multiple-return-values)")
	fmt.Println("---------------------------------------------------")
}
