// Copyright (c) 2026 Rasel Hossen

// ============================================================================
// Section 04: Types and Design
// Title: Complex Generic Constraints
// Level: Stretch
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn advanced constraint patterns including parameterized constraints, using interfaces as constraints, and creating reusable generic utilities.
//
// WHY THIS MATTERS:
//   - Think of a vending machine that accepts only certain payment methods. The constraint is not just "some type"-it's "anything with a Pay() method." Complex constraints allow you to define these nuanced rules.
//
// RUN:
//   go run ./04-types-design/14-complex-generic-constraints
//
// KEY TAKEAWAY:
//   - Learn advanced constraint patterns including parameterized constraints, using interfaces as constraints, and creating reusable generic utilities.
// ============================================================================

// See LICENSE for usage terms.

package main

import (
	"fmt"
)

//
//   - Interface constraints
//   - Multiple interface requirements
//   - Comparable constraint
//   - Parameterized constraints
//

type Numeric interface {
	Add(other int) int
	Multiply(other int) int
}

func ScaleAll[T Numeric](item T, factor int) int {
	return item.Multiply(factor)
}

type CustomInt int

func (c CustomInt) Add(other int) int {
	return int(c) + other
}

func (c CustomInt) Multiply(other int) int {
	return int(c) * other
}

type Stringable interface {
	ToString() string
}

func PrintAll[T Stringable](items []T) {
	for _, item := range items {
		fmt.Println("  ", item.ToString())
	}
}

type User struct {
	Name string
	Age  int
}

func (u User) ToString() string {
	return fmt.Sprintf("%s (%d)", u.Name, u.Age)
}

type CacheKey comparable

func GetOrSet[K comparable, V any](m map[K]V, key K, defaultVal V) V {
	if val, ok := m[key]; ok {
		return val
	}
	m[key] = defaultVal
	return defaultVal
}

type JSONMarshaler interface {
	MarshalJSON() ([]byte, error)
}

type Data struct {
	Value int
}

func (d Data) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"value":%d}`, d.Value)), nil
}

func MarshalAll[T JSONMarshaler](items []T) [][]byte {
	results := make([][]byte, len(items))
	for i, item := range items {
		results[i], _ = item.MarshalJSON()
	}
	return results
}

func main() {
	fmt.Println("=== Complex Generic Constraints ===")
	fmt.Println()

	fmt.Println("--- Interface Constraint with Custom Type ---")
	ci := CustomInt(5)
	fmt.Printf("  CustomInt 5 * 3 = %d\n", ci.Multiply(3))
	fmt.Printf("  CustomInt 5 + 2 = %d\n", ci.Add(2))
	result := ScaleAll(ci, 4)
	fmt.Printf("  ScaleAll(5, 4) = %d\n", result)

	fmt.Println()
	fmt.Println("--- Multiple Interface Requirements ---")
	users := []User{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}
	fmt.Println("  All users:")
	PrintAll(users)

	fmt.Println()
	fmt.Println("--- Comparable Constraint ---")
	cache := map[string]int{"cached": 100}
	val := GetOrSet(cache, "missing", 42)
	fmt.Printf("  Got value: %d (cache now: %v)\n", val, cache)

	fmt.Println()
	fmt.Println("--- JSON Marshaler Constraint ---")
	data := []Data{{Value: 1}, {Value: 2}}
	jsons := MarshalAll(data)
	for i, j := range jsons {
		fmt.Printf("  JSON %d: %s\n", i, string(j))
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Interfaces can be generic constraints")
	fmt.Println("  - Embed multiple interfaces for complex requirements")
	fmt.Println("  - comparable constraint allows == and !=")
	fmt.Println("  - Type must implement all methods in constraint")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.15 -> 04-types-design/15-generic-data-structures")
	fmt.Println("Current: TI.14 (complex-generic-constraints)")
	fmt.Println("Previous: TI.13 (method-values)")
	fmt.Println("---------------------------------------------------")
}
