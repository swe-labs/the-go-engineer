// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// Section 5: Types & Interfaces — Generics
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What generics are: functions/types that work with MULTIPLE types
//   - Type parameters: the [T constraint] syntax
//   - Type constraints: limiting what types T can be
//   - Built-in constraints: comparable, any
//   - Custom constraints with interface unions (int | float64 | ...)
//   - When to use generics vs interfaces
//
// ANALOGY:
//   A vending machine is generic. It doesn't care if it dispenses
//   sodas, snacks, or toys — the mechanism is the same.
//   The "type parameter" is what's in each slot.
//   The "constraint" says "must fit in the slot" (size/shape limit).
//
// HISTORY:
//   Generics were added in Go 1.18 (March 2022).
//   Before generics, you had to write duplicate functions for each type
//   or use interface{}/any and lose type safety.
//
// RUN: go run ./05-types-and-interfaces/5-generics
// ============================================================================

// --- TYPE CONSTRAINTS ---

// Numeric is a CUSTOM TYPE CONSTRAINT.
// It uses the | (union) operator to list all types that are allowed.
// Only these types can be used as T when calling functions constrained by Numeric.
type Numeric interface {
	int | int8 | int16 | int32 | int64 |
		float32 | float64
}

// --- GENERIC FUNCTIONS ---

// Sum adds up a slice of any Numeric type.
//
// SYNTAX BREAKDOWN:
//
//	func Sum[T Numeric](numbers []T) T
//	^^^^^^^^ ^^^^^^^^^  ^^^^^^^^^^  ^
//	func     type param  parameter  return
//	        [Name Constraint]
//
// T is a PLACEHOLDER for the actual type. When you call Sum[int](...),
// the compiler replaces T with int everywhere.
func Sum[T Numeric](numbers []T) T {
	var total T // Zero value of T (0 for int, 0.0 for float64)
	for _, n := range numbers {
		total += n // Works because all Numeric types support +
	}
	return total
}

// Filter returns only elements that satisfy the predicate function.
// This is a generic utility that works with ANY type (constraint: any).
// "any" is an alias for "interface{}" — it accepts all types.
func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0) // Empty slice of type T
	for _, item := range items {
		if predicate(item) { // Call the predicate function for each item
			result = append(result, item)
		}
	}
	return result
}

// Map transforms a slice of type T into a slice of type U.
// This demonstrates TWO type parameters in one function.
// T is the input type, U is the output type — they can be different.
func Map[T any, U any](items []T, transform func(T) U) []U {
	result := make([]U, len(items)) // Pre-allocate output slice
	for i, item := range items {
		result[i] = transform(item)
	}
	return result
}

// Contains checks if a slice contains a specific value.
// The "comparable" constraint means T must support == and != operators.
// Slices, maps, and functions are NOT comparable — this constraint prevents
// passing them, which would cause a runtime panic.
func Contains[T comparable](items []T, target T) bool {
	for _, item := range items {
		if item == target { // == only works with comparable types
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("=== Generics: One Function, Many Types ===")
	fmt.Println()

	// --- GENERIC SUM ---
	// One Sum function works with int, float64, int32, etc.
	// Before generics, you'd need SumInt(), SumFloat64(), SumInt32()...
	fmt.Println("--- Sum (Numeric constraint) ---")
	intSum := Sum([]int{10, 20, 30, 40})           // T inferred as int
	floatSum := Sum([]float64{1.5, 2.5, 3.0})      // T inferred as float64
	fmt.Printf("  Sum[int]:     %d\n", intSum)     // 100
	fmt.Printf("  Sum[float64]: %.1f\n", floatSum) // 7.0
	fmt.Println()

	// --- FILTER ---
	fmt.Println("--- Filter (any constraint) ---")
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Filter even numbers
	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Printf("  Evens: %v\n", evens) // [2 4 6 8 10]

	// Filter long words
	words := []string{"go", "rust", "python", "java", "c", "typescript"}
	longWords := Filter(words, func(w string) bool { return len(w) > 3 })
	fmt.Printf("  Long words: %v\n", longWords) // [rust python java typescript]
	fmt.Println()

	// --- MAP (Transform) ---
	fmt.Println("--- Map (two type parameters) ---")
	// Transform []int → []string
	labels := Map(numbers, func(n int) string {
		return fmt.Sprintf("#%d", n)
	})
	fmt.Printf("  Labels: %v\n", labels) // [#1 #2 #3 ... #10]

	// Transform []string → []int (word lengths)
	lengths := Map(words, func(w string) int { return len(w) })
	fmt.Printf("  Word lengths: %v\n", lengths) // [2 4 6 4 1 10]
	fmt.Println()

	// --- CONTAINS (comparable) ---
	fmt.Println("--- Contains (comparable constraint) ---")
	upper := Map(words, func(w string) string { return strings.ToUpper(w) })
	fmt.Printf("  Has 'GO': %t\n", Contains(upper, "GO"))     // true
	fmt.Printf("  Has 'RUBY': %t\n", Contains(upper, "RUBY")) // false
	fmt.Printf("  Has 5: %t\n", Contains(numbers, 5))         // true
	fmt.Printf("  Has 99: %t\n", Contains(numbers, 99))       // false

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Generics eliminate duplicate code for different types")
	fmt.Println("  - [T constraint] declares a type parameter with limits")
	fmt.Println("  - 'any' = accepts all types, 'comparable' = supports ==")
	fmt.Println("  - Custom constraints (int | float64) restrict to specific types")
	fmt.Println("  - Use generics for data structures, algorithms, utilities")
	fmt.Println("  - Use interfaces for behavior abstraction (polymorphism)")
}
