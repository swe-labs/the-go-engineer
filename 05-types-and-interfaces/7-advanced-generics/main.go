// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./05-types-and-interfaces/7-advanced-generics
package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// Section 05: Types & Interfaces — Advanced Generics
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Using Generics to build reusable functional patterns (Map / Filter).
//   - Defining operations on generic collections.
//   - Understanding Type Inference and its limits in Go.
//
// ENGINEERING DEPTH:
//   Prior to Go 1.18, utility functions like `Map` or `Filter` had to be re-written
//   for `int`, `string`, `struct`, etc., or rely on heavy `reflect` packages which
//   bypassed compile-time type safety.
//   Now, we can write type-safe, generic slice operators that compile down to
//   highly optimized monomorphized code.
// ============================================================================

// 1. Generic Filter Function
// Takes a slice of any type T, and a predicate function that returns a boolean.
// It returns a new slice containing only elements where the predicate returned true.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// 2. Generic Map Function
// Takes a slice of type T, and a transform function that converts T to type R.
// It returns a new slice of type R.
// Notice that Map requires TWO type parameters.
func Map[T any, R any](slice []T, transform func(T) R) []R {
	// Pre-allocate the result slice since we know the exact length.
	// This avoids expensive underlying array reallocations!
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = transform(v)
	}
	return result
}

type User struct {
	ID     int
	Name   string
	Active bool
}

func main() {
	// --- Example A: Filtering Integers ---
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
	evenNumbers := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("Evens from %v: %v\n", numbers, evenNumbers)

	// --- Example B: Mapping Strings to Uppercase ---
	words := []string{"go", "generics", "are", "cool"}
	shouted := Map(words, func(s string) string {
		return strings.ToUpper(s) + "!"
	})
	fmt.Printf("\nOriginal words: %v\n", words)
	fmt.Printf("Mapped words: %v\n", shouted)

	// --- Example C: Working with Custom Structs ---
	users := []User{
		{ID: 1, Name: "Alice", Active: true},
		{ID: 2, Name: "Bob", Active: false},
		{ID: 3, Name: "Charlie", Active: true},
	}

	// 1. Filter only active users
	activeUsers := Filter(users, func(u User) bool {
		return u.Active
	})

	// 2. Map to extract only their names
	activeUserNames := Map(activeUsers, func(u User) string {
		return u.Name
	})

	fmt.Printf("\nExtracted Active User Names: %v\n", activeUserNames)
}
