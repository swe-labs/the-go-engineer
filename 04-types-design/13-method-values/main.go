// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Method Values
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Extracting methods from receivers to create first-class function values.
//   - The mechanics of receiver binding in method values.
//   - Using method values as callbacks and event handlers.
//   - Storing method-bound behaviors in maps and slices.
//
// WHY THIS MATTERS:
//   - In Go, a method can be "detached" from its object and used as a
//     standalone function value while still retaining a reference to
//     its original receiver. This is the primary mechanism for
//     implementing flexible event systems and clean callbacks
//     without manually wrapping methods in closures.
//
// RUN:
//   go run ./04-types-design/13-method-values
//
// KEY TAKEAWAY:
//   - Method values are functions with a permanently bound receiver.
// ============================================================================

// See LICENSE for usage terms.

package main

import "fmt"

// Section 04: Types & Design - Method Values

// Counter represents a simple integer state with mutation methods.
// Counter (Struct): represents a simple integer state with mutation methods.
type Counter struct {
	Value int
}

// Increment increases the counter value by 1.
// Counter.Increment (Method): increases the counter value by 1.
func (c *Counter) Increment() {
	c.Value++
}

// Decrement decreases the counter value by 1.
// Counter.Decrement (Method): decreases the counter value by 1.
func (c *Counter) Decrement() {
	c.Value--
}

// GetValue returns the current state of the counter.
// Counter.GetValue (Method): returns the current state of the counter.
func (c *Counter) GetValue() int {
	return c.Value
}

// Handler represents an event target with specific callbacks.
// Handler (Struct): represents an event target with specific callbacks.
type Handler struct {
	Name string
}

// OnClick simulates a button click event.
// Handler.OnClick (Method): simulates a button click event.
func (h *Handler) OnClick() {
	fmt.Printf("  [Event] %s triggered click logic\n", h.Name)
}

// OnHover simulates a mouse hover event.
// Handler.OnHover (Method): simulates a mouse hover event.
func (h *Handler) OnHover() {
	fmt.Printf("  [Event] %s triggered hover logic\n", h.Name)
}

// runHandler accepts a generic callback function and executes it.
// runHandler (Function): accepts a generic callback function and executes it.
func runHandler(name string, handler func()) {
	fmt.Printf("Executing Handler: %s\n", name)
	handler()
}

func main() {
	fmt.Println("=== Method Values: First-Class Behavior ===")
	fmt.Println()

	// 1. Basic Method Extraction.
	// Assigning 'counter.Increment' to a variable captures the '*Counter'
	// receiver 'counter' and stores it inside the function value.
	fmt.Println("--- Receiver Binding ---")
	counter := &Counter{Value: 10}
	inc := counter.Increment
	inc()
	inc()
	fmt.Printf("  Counter value after detached calls: %d\n", counter.Value)
	fmt.Println()

	// 2. Methods as Callbacks.
	// Any method matching the signature 'func()' can be passed where
	// a plain function is expected.
	fmt.Println("--- Event Delegation ---")
	handler := &Handler{Name: "SubmitButton"}
	runHandler("click_event", handler.OnClick)
	fmt.Println()

	// 3. Instance Independence.
	// Method values are tied to specific instances. Calling fn1 increments
	// c1, while calling fn2 increments c2.
	fmt.Println("--- Instance Independence ---")
	c1 := &Counter{Value: 100}
	c2 := &Counter{Value: 200}
	fn1 := c1.Increment
	fn2 := c2.Increment
	fn1()
	fn2()
	fmt.Printf("  c1: %d, c2: %d\n", c1.Value, c2.Value)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.10 -> 04-types-design/10-payroll-processor")
	fmt.Println("Run    : go run ./04-types-design/10-payroll-processor")
	fmt.Println("Current: TI.13 (method-values)")
	fmt.Println("---------------------------------------------------")
}
