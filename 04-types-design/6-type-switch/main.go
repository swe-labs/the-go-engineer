// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import "fmt"

// ============================================================================
// Section 6: Types & Interfaces — Type Switch
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Type switch syntax and usage
//   - Handling multiple concrete types from an interface
//   - The comma-ok pattern in type switches
//
// RUN: go run ./04-types-design/6-type-switch
// ============================================================================

type Shape interface{}

type Circle struct{ Radius float64 }
type Rectangle struct{ Width, Height float64 }
type Triangle struct{ A, B, C float64 }

func describeShape(s Shape) string {
	switch v := s.(type) {
	case Circle:
		return fmt.Sprintf("Circle with radius %.2f", v.Radius)
	case Rectangle:
		return fmt.Sprintf("Rectangle %.2f x %.2f", v.Width, v.Height)
	case Triangle:
		return fmt.Sprintf("Triangle with sides %.2f, %.2f, %.2f", v.A, v.B, v.C)
	default:
		return "Unknown shape"
	}
}

func getArea(s Shape) interface{} {
	switch v := s.(type) {
	case Circle:
		return 3.14 * v.Radius * v.Radius
	case Rectangle:
		return v.Width * v.Height
	case Triangle:
		ss := (v.A + v.B + v.C) / 2
		return ss * (ss - v.A) * (ss - v.B) * (ss - v.C)
	default:
		return "Unknown"
	}
}

func main() {
	fmt.Println("=== Type Switch ===")
	fmt.Println()

	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 10, Height: 5},
		Triangle{A: 3, B: 4, C: 5},
	}

	fmt.Println("--- Describing Shapes ---")
	for _, s := range shapes {
		fmt.Printf("  %s\n", describeShape(s))
	}

	fmt.Println()
	fmt.Println("--- Getting Areas ---")
	for _, s := range shapes {
		area := getArea(s)
		fmt.Printf("  Area: %v\n", area)
	}

	fmt.Println()
	fmt.Println("--- Type Switch with Interface{} ---")
	var anything interface{} = "hello"
	switch v := anything.(type) {
	case string:
		fmt.Printf("  String: %s (len=%d)\n", v, len(v))
	case int:
		fmt.Printf("  Int: %d\n", v)
	case bool:
		fmt.Printf("  Bool: %t\n", v)
	default:
		fmt.Printf("  Unknown type: %T\n", v)
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Type switch checks multiple types in one statement")
	fmt.Println("  - value.(type) extracts the concrete type")
	fmt.Println("  - Default case handles unknown types")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.7 custom errors")
	fmt.Println("   Current: TI.6 (type switch)")
	fmt.Println("---------------------------------------------------")
}
