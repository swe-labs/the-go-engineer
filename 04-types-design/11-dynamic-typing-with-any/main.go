// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Empty interface, assertions, and nil interfaces
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Using the `any` keyword to handle untyped data.
//   - Recovering concrete types safely using assertions.
//   - The "Typed Nil" pitfall: why non-nil interfaces can contain nil values.
//
// WHY THIS MATTERS:
//   - The empty interface (`any`) is Go's tool for extreme flexibility.
//     However, misuse leads to runtime panics and subtle nil-checking
//     bugs. Understanding its internal representation is critical for
//     building robust generic-like systems and handling dynamic data.
//
// RUN:
//   go run ./04-types-design/11-dynamic-typing-with-any
//
// KEY TAKEAWAY:
//   - An interface value is only nil if BOTH its type and data are nil.
// ============================================================================

package main

import "fmt"

// Section 04: Types & Design - Dynamic Typing with Any

// notifier defines a simple notification behavior.
type notifier interface {
	Notify()
}

// emailNotifier implements the notifier interface for email delivery.
type emailNotifier struct {
	address string
}

// Notify prints a notification message to stdout.
func (e *emailNotifier) Notify() {
	fmt.Printf("notifying %s\n", e.address)
}

func inspect(v any) {
	switch value := v.(type) {
	case string:
		fmt.Printf("string: %q\n", value)
	case int:
		fmt.Printf("int: %d\n", value)
	case notifier:
		fmt.Println("notifier value detected")
		value.Notify()
	default:
		fmt.Printf("value of type %T\n", value)
	}
}

func main() {
	fmt.Println("=== Dynamic Typing: Working with the 'any' Interface ===")
	fmt.Println()

	// 1. Heterogeneous data.
	// 'any' allows us to pass values of any type into a single processing point.
	fmt.Println("--- Inspecting Varied Types ---")
	inspect("production-logs")
	inspect(500)
	inspect(&emailNotifier{address: "admin@systems.net"})
	fmt.Println()

	// 2. The 'Typed Nil' pitfall.
	// An interface variable stores a pair: (Type, Value).
	// If the Type is non-nil, the interface ITSELF is non-nil, even if the
	// underlying Value pointer is nil.
	fmt.Println("--- The Interface Nil Pitfall ---")
	var typedNil *emailNotifier = nil
	var interfaceValue any = typedNil

	fmt.Printf("  Concrete Pointer is nil: %t\n", typedNil == nil)
	fmt.Printf("  Interface variable is nil: %t\n", interfaceValue == nil)
	fmt.Println("  RESULT: Interface is NOT nil because it carries the '*emailNotifier' type info.")

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.12 -> 04-types-design/12-functional-options")
	fmt.Println("Run    : go run ./04-types-design/12-functional-options")
	fmt.Println("Current: TI.11 (dynamic-typing-with-any)")
	fmt.Println("---------------------------------------------------")
}
