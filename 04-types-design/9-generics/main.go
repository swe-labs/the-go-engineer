// Copyright (c) 2026 Rasel Hossen

// ============================================================================
// Section 04: Types and Design
// Title: Generics
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Defining generic functions using type parameters and constraints.
//   - Using built-in constraints: `any` and `comparable`.
//   - Implementing custom type constraints using interface unions.
//   - The mechanics of monomorphization and compile-time type safety.
//
// WHY THIS MATTERS:
//   - Generics solve the "duplicate code" problem for data structures and
//     algorithms. They allow you to write reusable logic that operates
//     on multiple types without sacrificing the performance or safety of
//     static typing.
//
// RUN:
//   go run ./04-types-design/9-generics
//
// KEY TAKEAWAY:
//   - Generics enable type-safe logic reuse across different data types.
// ============================================================================

// See LICENSE for usage terms.

package main

import (
	"fmt"
)

// Section 04: Types & Design - Generics

// Numeric defines a type constraint for all standard integer and floating-point types.
type Numeric interface {
	int | int8 | int16 | int32 | int64 |
		float32 | float64
}

// Sum calculates the arithmetic total of a slice of numeric values.
func Sum[T Numeric](numbers []T) T {
	var total T
	for _, n := range numbers {
		total += n
	}
	return total
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map transforms a slice of type T into a slice of type U using the transform function.
func Map[T any, U any](items []T, transform func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = transform(item)
	}
	return result
}

// Contains returns true if the target element exists within the slice.
// It requires T to be comparable (supports == and !=).
func Contains[T comparable](items []T, target T) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("=== Generics: Type-Safe Logic Reuse ===")
	fmt.Println()

	// 1. Numeric Constraints.
	// The Sum function works on any type satisfying the Numeric interface union.
	// The compiler generates specific machine code for each concrete type used.
	fmt.Println("--- Arithmetic Sums (int vs float64) ---")
	fmt.Printf("  Sum[int]:     %d\n", Sum([]int{10, 20, 30}))
	fmt.Printf("  Sum[float64]: %.2f\n", Sum([]float64{1.5, 2.7, 3.3}))
	fmt.Println()

	// 2. Functional Programming with 'any'.
	// Filter and Map operate on slices of any type. The type T is narrowed
	// based on the input slice at the call site.
	fmt.Println("--- Functional Utilities (Filter/Map) ---")
	numbers := []int{1, 2, 3, 4, 5, 6}
	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Printf("  Evens: %v\n", evens)

	words := []string{"go", "rust", "python"}
	lengths := Map(words, func(w string) int { return len(w) })
	fmt.Printf("  Word Lengths: %v\n", lengths)
	fmt.Println()

	// 3. Comparable types.
	// The Contains function requires T to support equality checks (==).
	// This prevents using it with types that cannot be compared (like slices or maps).
	fmt.Println("--- Set Membership (comparable) ---")
	fmt.Printf("  Has 'go':   %t\n", Contains(words, "go"))
	fmt.Printf("  Has 5:      %t\n", Contains(numbers, 5))
	fmt.Printf("  Has 'ruby': %t\n", Contains(words, "ruby"))

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.14 -> 04-types-design/14-complex-generic-constraints")
	fmt.Println("Run    : go run ./04-types-design/14-complex-generic-constraints")
	fmt.Println("Current: TI.9 (generics)")
	fmt.Println("---------------------------------------------------")
}
