// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Receiver Sets
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn the difference between value and pointer receivers and understand how method sets determine interface satisfaction.
//
// WHY THIS MATTERS:
//   - Think of a type's method set like a menu. A Counter value has only the Get() menu item. A *Counter pointer has the full menu: Get(), Inc(), Reset(). Interface satisfaction depends on which menu is visible.
//
// RUN:
//   go run ./04-types-design/7-receiver-sets
//
// KEY TAKEAWAY:
//   - Learn the difference between value and pointer receivers and understand how method sets determine interface satisfaction.
// ============================================================================

package main

import "fmt"

//
//   - Value vs pointer receivers
//   - Method sets: what methods a type actually provides
//   - How receiver type affects interface satisfaction
//

type Counter struct {
	Value int
}

func (c Counter) Get() int {
	return c.Value
}

func (c *Counter) Inc() {
	c.Value++
}

func (c *Counter) Reset() {
	c.Value = 0
}

type Reader interface {
	Get() int
}

func printValue(r Reader) {
	fmt.Printf("  Value: %d\n", r.Get())
}

func main() {
	fmt.Println("=== Receiver Sets Drill ===")
	fmt.Println()

	fmt.Println("--- Counter Value (value receiver only) ---")
	c := Counter{Value: 42}
	fmt.Printf("  Counter value: %d\n", c.Get())

	fmt.Println()
	fmt.Println("--- Counter Pointer (pointer receivers) ---")
	c2 := &Counter{Value: 100}
	c2.Inc()
	c2.Inc()
	fmt.Printf("  After two Inc(): %d\n", c2.Get())
	c2.Reset()
	fmt.Printf("  After Reset(): %d\n", c2.Get())

	fmt.Println()
	fmt.Println("--- Method Set: Counter vs *Counter ---")
	counter := Counter{Value: 0}

	fmt.Println("  Counter (value type):")
	fmt.Printf("    Has Get(): YES\n")
	fmt.Printf("    Has Inc(): NO (pointer receiver)\n")

	fmt.Println()
	fmt.Println("  *Counter (pointer type):")
	fmt.Printf("    Has Get(): YES (inherited from value)\n")
	fmt.Printf("    Has Inc(): YES\n")
	fmt.Printf("    Has Reset(): YES\n")

	fmt.Println()
	fmt.Println("--- Interface Satisfaction ---")
	var r Reader = Counter{Value: 10}
	printValue(r)

	var r2 Reader = &Counter{Value: 20}
	printValue(r2)

	fmt.Println()
	fmt.Println("--- Important: Can't pass value where pointer needed ---")
	counter = Counter{Value: 50}
	counter.Inc()

	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Value receiver methods work on copies")
	fmt.Println("  - Pointer receiver methods need the original")
	fmt.Println("  - A type's method set depends on whether you use value or pointer")
	fmt.Println("  - Interface satisfaction requires the method to exist on the type you pass")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.8 -> 04-types-design/8-custom-errors")
	fmt.Println("Current: TI.7 (receiver-sets)")
	fmt.Println("Previous: TI.6 (type-switch)")
	fmt.Println("---------------------------------------------------")
}
