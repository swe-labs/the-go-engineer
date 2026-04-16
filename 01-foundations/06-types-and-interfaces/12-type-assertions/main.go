// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import "fmt"

// ============================================================================
// Section 6: Types & Interfaces — Type Assertions
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Basic type assertions
//   - Comma-ok pattern for safe assertions
//   - Type switches vs type assertions
//
// RUN: go run ./01-foundations/06-types-and-interfaces/12-type-assertions
// ============================================================================

type Shape interface {
	Area() float64
}

type Circle struct{ Radius float64 }

func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

func printArea(s Shape) {
	fmt.Printf("  Shape area: %.2f\n", s.Area())
}

func main() {
	fmt.Println("=== Type Assertions ===")
	fmt.Println()

	fmt.Println("--- Basic Type Assertion ---")
	var i interface{} = "hello"
	s := i.(string)
	fmt.Printf("  Asserted string: %s\n", s)

	fmt.Println()
	fmt.Println("--- Comma-ok Pattern (Safe) ---")
	var vals []interface{} = []interface{}{42, "hello", 3.14, true}
	for _, v := range vals {
		if str, ok := v.(string); ok {
			fmt.Printf("  String: %s\n", str)
		} else if num, ok := v.(int); ok {
			fmt.Printf("  Int: %d\n", num)
		} else {
			fmt.Printf("  Other: %T\n", v)
		}
	}

	fmt.Println()
	fmt.Println("--- Type Assertion with Interface ---")
	var shape Shape = Circle{Radius: 5}
	c := shape.(Circle)
	fmt.Printf("  Circle radius from interface: %.2f\n", c.Radius)

	fmt.Println()
	fmt.Println("--- Type Assertion Panic Example ---")
	var unknown interface{} = 123
	// str := unknown.(string) // This would panic!
	// Safe way:
	if str, ok := unknown.(string); ok {
		fmt.Printf("  String: %s\n", str)
	} else {
		fmt.Printf("  Not a string, got: %T\n", unknown)
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Type assertion extracts concrete type from interface")
	fmt.Println("  - Use comma-ok pattern to avoid panics")
	fmt.Println("  - Type switch handles multiple types cleanly")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.13 nil interfaces")
	fmt.Println("   Current: TI.12 (type assertions)")
	fmt.Println("---------------------------------------------------")
}
