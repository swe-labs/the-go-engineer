// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Empty interface, assertions, and nil interfaces
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how `any` works, how to extract concrete types safely, and why typed nil values can still make an interface non-nil.
//
// WHY THIS MATTERS:
//   - An interface value carries both a dynamic type and a dynamic value. Bugs appear when you forget that it needs both pieces.
//
// RUN:
//   go run ./04-types-design/11-dynamic-typing-with-any
//
// KEY TAKEAWAY:
//   - Learn how `any` works, how to extract concrete types safely, and why typed nil values can still make an interface non-nil.
// ============================================================================

package main

import "fmt"

//

type notifier interface {
	Notify()
}

type emailNotifier struct {
	address string
}

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
	fmt.Println("=== Dynamic Typing with any ===")
	inspect("orders")
	inspect(3)
	inspect(&emailNotifier{address: "ops@example.com"})

	var typedNil *emailNotifier
	var value any = typedNil
	fmt.Printf("typed nil inside any == nil? %t\n", value == nil)
	fmt.Println("The interface keeps type information even when the pointer payload is nil.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.12 -> 04-types-design/12-functional-options")
	fmt.Println("Current: TI.11 (dynamic-typing-with-any)")
	fmt.Println("Previous: TI.10 (payroll-processor)")
	fmt.Println("---------------------------------------------------")
}
