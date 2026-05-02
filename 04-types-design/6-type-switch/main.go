// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Type Switch
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Using type switches to branch logic based on concrete types.
//   - Extracting concrete values from interfaces using the `value.(type)` syntax.
//   - Handling multiple types and default cases in a single block.
//
// WHY THIS MATTERS:
//   - Interfaces abstract behavior, but sometimes a specific operation
//     requires access to the underlying concrete type. Type switches
//     provide a type-safe, readable mechanism for inspecting interface
//     contents at runtime.
//
// RUN:
//   go run ./04-types-design/6-type-switch
//
// KEY TAKEAWAY:
//   - Type switches facilitate runtime type discovery and branching.
// ============================================================================

// See LICENSE for usage terms.

package main

import "fmt"

// Section 04: Types & Design - Type Switch

// Shape is an empty interface representing any geometric shape.
// Shape (Interface): is an empty interface representing any geometric shape.
type Shape interface{}

// Circle represents a circular geometry.
// Circle (Struct): represents a circular geometry.
type Circle struct{ Radius float64 }

// Rectangle represents a four-sided polygon.
// Rectangle (Struct): represents a four-sided polygon.
type Rectangle struct{ Width, Height float64 }

// Triangle represents a three-sided polygon.
// Triangle (Struct): represents a three-sided polygon.
type Triangle struct{ A, B, C float64 }

// describeShape (Function): runs the describe shape step and keeps its inputs, outputs, or errors visible.
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

// getArea (Function): runs the get area step and keeps its inputs, outputs, or errors visible.
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
	fmt.Println("=== Type Switch: Runtime Type Discovery ===")
	fmt.Println()

	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 10, Height: 5},
		Triangle{A: 3, B: 4, C: 5},
	}

	// 1. Differentiated behavior based on concrete types.
	// We use the type switch to generate descriptive text for each specific implementation.
	fmt.Println("--- Describing Shapes ---")
	for _, s := range shapes {
		fmt.Printf("  %s\n", describeShape(s))
	}

	// 2. State-dependent logic.
	// Type switches allow for type-safe access to fields (like Radius or Width)
	// that are not part of the common interface.
	fmt.Println()
	fmt.Println("--- Area Calculation ---")
	for _, s := range shapes {
		area := getArea(s)
		fmt.Printf("  Area: %v\n", area)
	}

	// 3. Handling primitive types in empty interfaces.
	// Type switches are commonly used to process untyped 'any' data from
	// external sources (JSON, DB, etc.).
	fmt.Println()
	fmt.Println("--- Processing 'any' Data ---")
	var anything any = "dynamic content"
	switch v := anything.(type) {
	case string:
		fmt.Printf("  Found String: %q (len=%d)\n", v, len(v))
	case int:
		fmt.Printf("  Found Int: %d\n", v)
	case bool:
		fmt.Printf("  Found Bool: %t\n", v)
	default:
		fmt.Printf("  Unknown Type: %T\n", v)
	}

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.11 -> 04-types-design/11-dynamic-typing-with-any")
	fmt.Println("Run    : go run ./04-types-design/11-dynamic-typing-with-any")
	fmt.Println("Current: TI.6 (type-switch)")
	fmt.Println("---------------------------------------------------")
}
