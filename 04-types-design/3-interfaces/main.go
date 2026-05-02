// Copyright (c) 2026 Rasel Hossen

// ============================================================================
// Section 04: Types and Design
// Title: Interfaces
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Defining behavioral contracts using the `interface` keyword.
//   - Implicit interface satisfaction (structural typing).
//   - Achieving polymorphism without class inheritance.
//   - The internal memory representation of interface values (itab).
//
// WHY THIS MATTERS:
//   - Interfaces are Go's primary tool for abstraction and decoupling.
//     They enable the creation of interchangeable components and
//     facilitate rigorous testing through dependency injection.
//
// RUN:
//   go run ./04-types-design/3-interfaces
//
// KEY TAKEAWAY:
//   - Interfaces decouple requirements from concrete implementations.
// ============================================================================

// See LICENSE for usage terms.

package main

import (
	"fmt"
	"math"
)

// Section 04: Types & Design - Interfaces
//   - What interfaces are: contracts that define behavior
//   - Implicit interface satisfaction (no "implements" keyword)
//   - Polymorphism: one function, many types
//   - The empty interface (any) and type assertions
//   - Interface internals: the 2-word struct (type + data pointers)
//   - Real-world design: "Accept interfaces, return structs"
//

// Shape defines the behavioral contract for geometry calculations.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle represents a four-sided polygon with right angles.
type Rectangle struct {
	Width  float64
	Height float64
}

// Area calculates the surface area of the rectangle.
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the distance around the rectangle.
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle(%.1f x %.1f)", r.Width, r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle(r=%.1f)", c.Radius)
}

type Triangle struct {
	A, B, C float64
}

func (t Triangle) Area() float64 {
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

func (t Triangle) String() string {
	return fmt.Sprintf("Triangle(%.1f, %.1f, %.1f)", t.A, t.B, t.C)
}

func printShapeInfo(s Shape) {
	fmt.Printf("  %-30s Area: %8.2f  Perimeter: %8.2f\n", s, s.Area(), s.Perimeter())
}

func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

func main() {
	fmt.Println("=== Interfaces: Decoupling Logic from Implementation ===")
	fmt.Println()

	// 1. Concrete implementations.
	// Different types provide their own specific logic for the same behavior.
	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 7}
	tri := Triangle{A: 3, B: 4, C: 5}

	fmt.Println("--- Individual Implementation Details ---")
	printShapeInfo(rect)
	printShapeInfo(circle)
	printShapeInfo(tri)
	fmt.Println()

	// 2. Polymorphic behavior.
	// Any type that implements the required methods can be treated as a 'Shape'.
	// This allows writing functions like totalArea that operate on the abstraction.
	allShapes := []Shape{rect, circle, tri}
	fmt.Printf("Total area of %d shapes: %.2f\n", len(allShapes), totalArea(allShapes))
	fmt.Println()

	// 3. Type Discovery.
	// Interfaces carry dynamic type information. Use assertions to recover concrete types.
	fmt.Println("--- Runtime Type Discovery ---")
	var s Shape = circle

	// Safe type assertion with comma-ok idiom.
	if c, ok := s.(Circle); ok {
		fmt.Printf("  Identity confirmed: Circle with radius %.1f\n", c.Radius)
	}

	// Type switch for branching based on multiple implementations.
	for _, shape := range allShapes {
		switch v := shape.(type) {
		case Rectangle:
			fmt.Printf("  Processing Rectangle: %.1fx%.1f\n", v.Width, v.Height)
		case Circle:
			fmt.Printf("  Processing Circle: radius %.1f\n", v.Radius)
		case Triangle:
			fmt.Printf("  Processing Triangle: sides %.1f, %.1f, %.1f\n", v.A, v.B, v.C)
		}
	}

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.4 -> 04-types-design/4-interface-embedding")
	fmt.Println("Run    : go run ./04-types-design/4-interface-embedding")
	fmt.Println("Current: TI.3 (interfaces)")
	fmt.Println("---------------------------------------------------")
}
