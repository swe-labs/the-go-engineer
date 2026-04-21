// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

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
// Run: go run ./03-functions-errors/3-multiple-return-values

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
	fmt.Println("---------------------------------------------------")
}
