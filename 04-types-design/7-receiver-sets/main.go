// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Receiver Sets
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Defining method sets for value types vs. pointer types.
//   - How receiver choice determines interface implementation.
//   - The mechanics of method inheritance for pointer types.
//
// WHY THIS MATTERS:
//   - Method sets are the strict rules the compiler uses to verify
//     interface satisfaction. Understanding these rules prevents common
//     "type does not implement interface" errors during system
//     integration.
//
// RUN:
//   go run ./04-types-design/7-receiver-sets
//
// KEY TAKEAWAY:
//   - Pointer types possess a superset of the methods available to value types.
// ============================================================================

package main

import "fmt"

// Section 04: Types & Design - Receiver Sets
//   - Value vs pointer receivers
//   - Method sets: what methods a type actually provides
//   - How receiver type affects interface satisfaction
//

// Counter represents a simple stateful integer.
// Counter (Struct): represents a simple stateful integer.
type Counter struct {
	Value int
}

// Get returns the current counter value. It uses a value receiver.
// Counter.Get (Method): returns the current counter value. It uses a value receiver.
func (c Counter) Get() int {
	return c.Value
}

// Inc increments the counter. It requires a pointer receiver to modify the state.
// Counter.Inc (Method): increments the counter. It requires a pointer receiver to modify the state.
func (c *Counter) Inc() {
	c.Value++
}

// Reset clears the counter. It requires a pointer receiver.
// Counter.Reset (Method): clears the counter. It requires a pointer receiver.
func (c *Counter) Reset() {
	c.Value = 0
}

// Reader defines the behavioral contract for reading a value.
// Reader (Interface): defines the behavioral contract for reading a value.
type Reader interface {
	Get() int
}

// printValue (Function): runs the print value step and keeps its inputs, outputs, or errors visible.
func printValue(r Reader) {
	fmt.Printf("  Value: %d\n", r.Get())
}

func main() {
	fmt.Println("=== Receiver Sets: Interface Satisfaction Rules ===")
	fmt.Println()

	// 1. Counter using a value receiver.
	// Both Counter and *Counter satisfy Reader because Get() has a value receiver.
	var c Counter
	var r Reader = c
	fmt.Printf("Reader (value): %d\n", r.Get())

	r = &c
	fmt.Printf("Reader (pointer): %d\n", r.Get())
	fmt.Println()

	// 2. Counter using a pointer receiver.
	// A value type (T) does NOT satisfy an interface if the method uses a pointer receiver (*T).
	// UNCOMMENTING the line below will cause a compile-time error:
	// "Counter does not implement Writer (Inc method has pointer receiver)"
	// var w Writer = c

	var ptr *Counter = &c
	ptr.Inc()
	fmt.Printf("After Inc() via pointer: %d\n", c.Get())
	fmt.Println()

	// 3. Automatic address-taking (Syntactic Sugar).
	// Go allows calling a pointer method on a value variable if it is addressable.
	// The compiler implicitly transforms 'c.Reset()' into '(&c).Reset()'.
	c.Reset()
	fmt.Printf("After Reset() via value: %d\n", c.Get())

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.3 -> 04-types-design/3-interfaces")
	fmt.Println("Run    : go run ./04-types-design/3-interfaces")
	fmt.Println("Current: TI.7 (receiver-sets)")
	fmt.Println("---------------------------------------------------")
}
