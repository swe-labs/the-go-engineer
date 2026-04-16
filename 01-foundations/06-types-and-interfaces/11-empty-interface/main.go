// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// ============================================================================
// Section 6: Types & Interfaces — Empty Interface
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The empty interface (any/interface{}) and its uses
//   - Type assertions with empty interface
//   - When to use empty interface vs generics
//
// RUN: go run ./01-foundations/06-types-and-interfaces/11-empty-interface
// ============================================================================

func printAny(v any) {
	fmt.Printf("  Type: %T, Value: %v\n", v, v)
}

func processAny(v any) {
	switch val := v.(type) {
	case string:
		fmt.Printf("  String: %s (upper: %s)\n", val, fmt.Sprintf("%s", val))
	case int:
		fmt.Printf("  Int: %d (doubled: %d)\n", val, val*2)
	case float64:
		fmt.Printf("  Float: %.2f (rounded: %d)\n", val, int(val))
	case []int:
		sum := 0
		for _, n := range val {
			sum += n
		}
		fmt.Printf("  []int: %v (sum: %d)\n", val, sum)
	default:
		fmt.Printf("  Unknown type: %T\n", v)
	}
}

func acceptAnything(items ...any) {
	fmt.Println("  Accepting any number of arguments:")
	for i, item := range items {
		fmt.Printf("    [%d] Type: %T\n", i, item)
	}
}

func main() {
	fmt.Println("=== Empty Interface (any) ===")
	fmt.Println()

	fmt.Println("--- Storing Any Type ---")
	var x any
	x = 42
	printAny(x)
	x = "hello"
	printAny(x)
	x = []int{1, 2, 3}
	printAny(x)
	x = map[string]int{"a": 1}
	printAny(x)

	fmt.Println()
	fmt.Println("--- Type Switch with any ---")
	values := []any{42, "hello", 3.14, []int{1, 2}, true}
	for _, v := range values {
		processAny(v)
	}

	fmt.Println()
	fmt.Println("--- Variadic any ---")
	acceptAnything(1, "two", 3.0, []int{4, 5})

	fmt.Println()
	fmt.Println("--- Type Assertion from any ---")
	var z any = "converted"
	if s, ok := z.(string); ok {
		fmt.Printf("  Asserted string: %s\n", s)
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - any (interface{}) can hold any Go value")
	fmt.Println("  - Use type assertions or switches to extract values")
	fmt.Println("  - Prefer generics when possible for type safety")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: TI.12 type assertions")
	fmt.Println("   Current: TI.11 (empty interface)")
	fmt.Println("---------------------------------------------------")
}
