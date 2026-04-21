// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import (
	"fmt"
	"math"
)

// ============================================================================
// Section 6: Types & Interfaces — Interfaces
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
// RUN: go run ./04-types-design/3-interfaces
// ============================================================================

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle(%.1f × %.1f)", r.Width, r.Height)
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
	fmt.Println("=== Interfaces: Contracts for Behavior ===")
	fmt.Println()

	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 7}
	tri := Triangle{A: 3, B: 4, C: 5}

	fmt.Println("Individual shapes:")
	printShapeInfo(rect)
	printShapeInfo(circle)
	printShapeInfo(tri)
	fmt.Println()

	allShapes := []Shape{rect, circle, tri}
	fmt.Printf("Total area of %d shapes: %.2f\n", len(allShapes), totalArea(allShapes))
	fmt.Println()

	fmt.Println("=== Type Assertions ===")
	var s Shape = circle

	if c, ok := s.(Circle); ok {
		fmt.Printf("  It's a Circle! Radius = %.1f\n", c.Radius)
	}

	if _, ok := s.(Rectangle); !ok {
		fmt.Println("  Not a Rectangle")
	}
	fmt.Println()

	fmt.Println("=== Type Switch ===")
	for _, shape := range allShapes {
		switch v := shape.(type) {
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
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.4 interface embedding")
	fmt.Println("   Current: TI.3 (interfaces)")
	fmt.Println("---------------------------------------------------")
}
