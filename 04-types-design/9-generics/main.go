// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// Section 6: Types & Interfaces — Generics
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
// RUN: go run ./04-types-design/9-generics
// ============================================================================

type Numeric interface {
	int | int8 | int16 | int32 | int64 |
		float32 | float64
}

func Sum[T Numeric](numbers []T) T {
	var total T
	for _, n := range numbers {
		total += n
	}
	return total
}

func Filter[T any](items []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, U any](items []T, transform func(T) U) []U {
	result := make([]U, len(items))
	for i, item := range items {
		result[i] = transform(item)
	}
	return result
}

func Contains[T comparable](items []T, target T) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("=== Generics: One Function, Many Types ===")
	fmt.Println()

	fmt.Println("--- Sum (Numeric constraint) ---")
	intSum := Sum([]int{10, 20, 30, 40})
	floatSum := Sum([]float64{1.5, 2.5, 3.0})
	fmt.Printf("  Sum[int]:     %d\n", intSum)
	fmt.Printf("  Sum[float64]: %.1f\n", floatSum)
	fmt.Println()

	fmt.Println("--- Filter (any constraint) ---")
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
	fmt.Printf("  Evens: %v\n", evens)

	words := []string{"go", "rust", "python", "java", "c", "typescript"}
	longWords := Filter(words, func(w string) bool { return len(w) > 3 })
	fmt.Printf("  Long words: %v\n", longWords)
	fmt.Println()

	fmt.Println("--- Map (two type parameters) ---")
	labels := Map(numbers, func(n int) string {
		return fmt.Sprintf("#%d", n)
	})
	fmt.Printf("  Labels: %v\n", labels)

	lengths := Map(words, func(w string) int { return len(w) })
	fmt.Printf("  Word lengths: %v\n", lengths)
	fmt.Println()

	fmt.Println("--- Contains (comparable constraint) ---")
	upper := Map(words, func(w string) string { return strings.ToUpper(w) })
	fmt.Printf("  Has 'GO': %t\n", Contains(upper, "GO"))
	fmt.Printf("  Has 'RUBY': %t\n", Contains(upper, "RUBY"))
	fmt.Printf("  Has 5: %t\n", Contains(numbers, 5))
	fmt.Printf("  Has 99: %t\n", Contains(numbers, 99))

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Generics eliminate duplicate code for different types")
	fmt.Println("  - [T constraint] declares a type parameter with limits")
	fmt.Println("  - 'any' = accepts all types, 'comparable' = supports ==")
	fmt.Println("  - Custom constraints (int | float64) restrict to specific types")
	fmt.Println("  - Use generics for data structures, algorithms, utilities")
	fmt.Println("  - Use interfaces for behavior abstraction (polymorphism)")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.10 payroll processor")
	fmt.Println("   Current: TI.9 (generics)")
	fmt.Println("---------------------------------------------------")
}
