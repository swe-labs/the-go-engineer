// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"math"
)

// ============================================================================
// Section 5: Types & Interfaces — Interfaces
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What interfaces are: contracts that define behavior
//   - Implicit interface satisfaction (no "implements" keyword)
//   - Polymorphism: one function, many types
//   - The empty interface (any) and type assertions
//   - Interface internals: the 2-word struct (type + data pointers)
//   - Real-world design: "Accept interfaces, return structs"
//
// ANALOGY:
//   Think of a power outlet. The outlet defines a "contract":
//   "I accept anything with two prongs and a ground pin."
//   A lamp, a phone charger, and a refrigerator all satisfy this contract
//   differently, but the outlet doesn't care HOW they work internally —
//   only that they have the right "shape" (methods).
//
//   In Go: the outlet is an INTERFACE, and each appliance is a STRUCT
//   that implements those methods.
//
// ENGINEERING DEPTH:
//   An interface value is internally a 2-word struct:
//     Word 1: pointer to the TYPE DESCRIPTOR (what concrete type is stored)
//     Word 2: pointer to the DATA (the actual struct value)
//   This is why interface calls are slightly slower than direct calls —
//   the runtime must look up the method through the type descriptor.
//   This is called "dynamic dispatch" (same concept as virtual methods in C++).
//
// RUN: go run ./05-types-and-interfaces/3-interfaces
// ============================================================================

// Shape is an INTERFACE — it defines a CONTRACT.
// Any type that has BOTH Area() and Perimeter() methods
// automatically satisfies this interface.
//
// Notice: there is NO "implements" keyword.
// Go uses "implicit" or "structural" satisfaction:
//
//	If your type has the right methods → it satisfies the interface. Period.
//	This is called "duck typing": if it quacks like a duck, it IS a duck.
type Shape interface {
	Area() float64      // Calculate the total surface area
	Perimeter() float64 // Calculate the total perimeter (boundary length)
}

// --- CONCRETE TYPE 1: Rectangle ---

// Rectangle is a concrete type that will satisfy the Shape interface.
// It has no idea the Shape interface even exists — and that's the point.
type Rectangle struct {
	Width  float64 // Horizontal dimension
	Height float64 // Vertical dimension
}

// Area calculates the area of a rectangle: width × height.
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter calculates the perimeter: 2 × (width + height).
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// String implements fmt.Stringer for pretty printing.
func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle(%.1f × %.1f)", r.Width, r.Height)
}

// --- CONCRETE TYPE 2: Circle ---

// Circle is a completely different type that ALSO satisfies Shape.
// The interface doesn't care that Circle stores a Radius while
// Rectangle stores Width/Height — it only cares about the methods.
type Circle struct {
	Radius float64 // Distance from center to edge
}

// Area calculates the area of a circle: π × r².
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter calculates the circumference: 2 × π × r.
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// String implements fmt.Stringer.
func (c Circle) String() string {
	return fmt.Sprintf("Circle(r=%.1f)", c.Radius)
}

// --- CONCRETE TYPE 3: Triangle ---

// Triangle with three sides. Demonstrates that ANY type can join the party
// as long as it has Area() and Perimeter().
type Triangle struct {
	A, B, C float64 // The three side lengths
}

// Area uses Heron's formula: √(s(s-a)(s-b)(s-c)) where s = (a+b+c)/2.
func (t Triangle) Area() float64 {
	s := (t.A + t.B + t.C) / 2 // Semi-perimeter
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

// Perimeter of a triangle is the sum of all three sides.
func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

// String implements fmt.Stringer.
func (t Triangle) String() string {
	return fmt.Sprintf("Triangle(%.1f, %.1f, %.1f)", t.A, t.B, t.C)
}

// --- POLYMORPHIC FUNCTION ---

// printShapeInfo accepts ANY type that satisfies the Shape interface.
// This is POLYMORPHISM: one function works with Rectangle, Circle,
// Triangle, and any future type that implements Shape.
//
// This is the "Accept interfaces, return structs" principle:
//   - Functions ACCEPT interfaces (flexible, testable)
//   - Functions RETURN concrete types (precise, discoverable)
func printShapeInfo(s Shape) {
	fmt.Printf("  %-30s Area: %8.2f  Perimeter: %8.2f\n", s, s.Area(), s.Perimeter())
}

// totalArea demonstrates processing a collection of different shapes.
// The []Shape slice can hold Rectangles, Circles, AND Triangles
// all at the same time — that's the power of interfaces.
func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes { // Each s is a Shape — could be any concrete type
		total += s.Area()
	}
	return total
}

func main() {
	fmt.Println("=== Interfaces: Contracts for Behavior ===")
	fmt.Println()

	// --- CREATE CONCRETE TYPES ---
	// These structs have no idea about the Shape interface.
	// They just happen to have the right methods.
	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 7}
	tri := Triangle{A: 3, B: 4, C: 5} // A classic right triangle (3-4-5)

	// --- POLYMORPHISM IN ACTION ---
	// One function handles three completely different types.
	fmt.Println("Individual shapes:")
	printShapeInfo(rect)   // Works with Rectangle
	printShapeInfo(circle) // Works with Circle
	printShapeInfo(tri)    // Works with Triangle
	fmt.Println()

	// --- INTERFACE SLICES ---
	// A []Shape slice holds ANY type that satisfies Shape.
	// This is how Go achieves polymorphic collections without generics.
	allShapes := []Shape{rect, circle, tri}
	fmt.Printf("Total area of %d shapes: %.2f\n", len(allShapes), totalArea(allShapes))
	fmt.Println()

	// --- TYPE ASSERTIONS ---
	// Sometimes you need to extract the concrete type from an interface.
	// A type assertion checks if the interface value holds a specific type.
	//
	// Syntax: value, ok := interfaceValue.(ConcreteType)
	//   ok = true  → the assertion succeeded, value is the concrete type
	//   ok = false → the interface holds a different type
	fmt.Println("=== Type Assertions ===")
	var s Shape = circle // s is a Shape, but it holds a Circle

	// Safe assertion with comma-ok pattern (ALWAYS use this form)
	if c, ok := s.(Circle); ok {
		fmt.Printf("  It's a Circle! Radius = %.1f\n", c.Radius)
	}

	// Check for Rectangle (will fail because s holds a Circle)
	if _, ok := s.(Rectangle); !ok {
		fmt.Println("  Not a Rectangle")
	}
	fmt.Println()

	// --- TYPE SWITCH ---
	// When you need to handle multiple possible types, use a type switch.
	// This is cleaner than chaining type assertions.
	fmt.Println("=== Type Switch ===")
	for _, shape := range allShapes {
		switch v := shape.(type) { // v is the concrete type
		case Rectangle:
			fmt.Printf("  Rectangle: %.1f × %.1f\n", v.Width, v.Height)
		case Circle:
			fmt.Printf("  Circle: radius = %.1f\n", v.Radius)
		case Triangle:
			fmt.Printf("  Triangle: sides = %.1f, %.1f, %.1f\n", v.A, v.B, v.C)
		}
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Interfaces define WHAT a type can do (contract)")
	fmt.Println("  - Types satisfy interfaces IMPLICITLY (no 'implements' keyword)")
	fmt.Println("  - Polymorphism: one function, many types")
	fmt.Println("  - Accept interfaces, return structs (design principle)")
	fmt.Println("  - Use type assertions (value, ok) to extract concrete types")
	fmt.Println("  - Interfaces are Go's primary tool for abstraction and testing")
}
