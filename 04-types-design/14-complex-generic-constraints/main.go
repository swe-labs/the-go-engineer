// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Complex Generic Constraints
// Level: Stretch
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Defining complex type constraints using method-based interfaces.
//   - Composing multiple behaviors into a single generic requirement.
//   - Utilizing the `comparable` constraint for map-based utilities.
//   - The mechanics of static satisfaction vs dynamic dispatch.
//
// WHY THIS MATTERS:
//   - Simple type unions are often insufficient for library authors.
//     Complex constraints allow you to enforce behavioral requirements
//     on generic types, ensuring that any T passed to your logic
//     possesses the necessary methods to execute correctly and safely.
//
// RUN:
//   go run ./04-types-design/14-complex-generic-constraints
//
// KEY TAKEAWAY:
//   - Generic constraints define the behavioral contract for type parameters.
// ============================================================================

// See LICENSE for usage terms.

package main

import (
	"fmt"
)

// Section 04: Types & Design - Complex Generic Constraints

// Numeric defines a requirement for types that support addition and multiplication.
// Numeric (Interface): defines a requirement for types that support addition and multiplication.
type Numeric interface {
	Add(other int) int
	Multiply(other int) int
}

// ScaleAll uses a generic type T constrained by the Numeric interface.
// ScaleAll (Function): uses a generic type T constrained by the Numeric interface.
func ScaleAll[T Numeric](item T, factor int) int {
	return item.Multiply(factor)
}

// CustomInt satisfies the Numeric constraint by providing Add and Multiply.
// CustomInt (Type): satisfies the Numeric constraint by providing Add and Multiply.
type CustomInt int

// Add returns the sum of the CustomInt and another integer.
// CustomInt.Add (Method): returns the sum of the CustomInt and another integer.
func (c CustomInt) Add(other int) int {
	return int(c) + other
}

// Multiply returns the product of the CustomInt and another integer.
// CustomInt.Multiply (Method): returns the product of the CustomInt and another integer.
func (c CustomInt) Multiply(other int) int {
	return int(c) * other
}

// Stringable defines a requirement for types that provide a textual view.
// Stringable (Interface): defines a requirement for types that provide a textual view.
type Stringable interface {
	ToString() string
}

// PrintAll processes a slice of any type T that satisfies the Stringable constraint.
// PrintAll (Function): processes a slice of any type T that satisfies the Stringable constraint.
func PrintAll[T Stringable](items []T) {
	for _, item := range items {
		fmt.Println("  ", item.ToString())
	}
}

// User represents a system entity that implements Stringable.
// User (Struct): represents a system entity that implements Stringable.
type User struct {
	Name string
	Age  int
}

// ToString provides a formatted string representation of the User.
// User.ToString (Method): provides a formatted string representation of the User.
func (u User) ToString() string {
	return fmt.Sprintf("%s (%d)", u.Name, u.Age)
}

// GetOrSet retrieves a value from a map or initializes it with a default if missing.
// GetOrSet (Function): retrieves a value from a map or initializes it with a default if missing.
func GetOrSet[K comparable, V any](m map[K]V, key K, defaultVal V) V {
	if val, ok := m[key]; ok {
		return val
	}
	m[key] = defaultVal
	return defaultVal
}

// JSONMarshaler (Interface): captures the behavior boundary the json marshaler example depends on.
type JSONMarshaler interface {
	MarshalJSON() ([]byte, error)
}

// Data (Struct): groups the state used by the data example boundary.
type Data struct {
	Value int
}

// Data.MarshalJSON (Method): applies the marshal json operation to receiver state at a visible boundary.
func (d Data) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"value":%d}`, d.Value)), nil
}

// MarshalAll (Function): runs the marshal all step and keeps its inputs, outputs, or errors visible.
func MarshalAll[T JSONMarshaler](items []T) [][]byte {
	results := make([][]byte, len(items))
	for i, item := range items {
		results[i], _ = item.MarshalJSON()
	}
	return results
}

func main() {
	fmt.Println("=== Complex Generic Constraints: Behavioral Requirements ===")
	fmt.Println()

	// 1. Method-based constraints.
	// The Numeric interface requires Add and Multiply methods. CustomInt satisfies this.
	fmt.Println("--- Arithmetic Behavior (CustomInt) ---")
	ci := CustomInt(5)
	result := ScaleAll(ci, 4)
	fmt.Printf("  Scaled CustomInt: %d\n", result)
	fmt.Println()

	// 2. Behavioral abstraction.
	// PrintAll requires any type that provides a ToString() method.
	fmt.Println("--- Behavioral Abstraction (ToString) ---")
	users := []User{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}
	PrintAll(users)
	fmt.Println()

	// 3. Built-in constraints (comparable).
	// The 'comparable' keyword allows generic logic to use equality operators (==, !=),
	// making it essential for map-key operations.
	fmt.Println("--- Utility Logic (comparable map keys) ---")
	cache := map[string]int{"active": 1}
	val := GetOrSet(cache, "init", 42)
	fmt.Printf("  Cache Result: %d -> Store: %v\n", val, cache)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.15 -> 04-types-design/15-generic-data-structures")
	fmt.Println("Run    : go run ./04-types-design/15-generic-data-structures")
	fmt.Println("Current: TI.14 (complex-generic-constraints)")
	fmt.Println("---------------------------------------------------")
}
