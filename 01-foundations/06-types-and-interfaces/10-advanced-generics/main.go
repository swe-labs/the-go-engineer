// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// Section 6: Types & Interfaces — Advanced Generics
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Using Generics to build reusable functional patterns (Map / Filter).
//   - Defining operations on generic collections.
//   - Understanding Type Inference and its limits in Go.
//
// RUN: go run ./01-foundations/06-types-and-interfaces/10-advanced-generics
// ============================================================================

func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func Map[T any, R any](slice []T, transform func(T) R) []R {
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
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
	evenNumbers := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("Evens from %v: %v\n", numbers, evenNumbers)

	words := []string{"go", "generics", "are", "cool"}
	shouted := Map(words, func(s string) string {
		return strings.ToUpper(s) + "!"
	})
	fmt.Printf("\nOriginal words: %v\n", words)
	fmt.Printf("Mapped words: %v\n", shouted)

	users := []User{
		{ID: 1, Name: "Alice", Active: true},
		{ID: 2, Name: "Bob", Active: false},
		{ID: 3, Name: "Charlie", Active: true},
	}

	activeUsers := Filter(users, func(u User) bool {
		return u.Active
	})

	activeUserNames := Map(activeUsers, func(u User) string {
		return u.Name
	})

	fmt.Printf("\nExtracted Active User Names: %v\n", activeUserNames)

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: TI.11 empty interface")
	fmt.Println("   Current: TI.10 (advanced generics)")
	fmt.Println("---------------------------------------------------")
}
