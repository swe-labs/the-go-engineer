// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import "fmt"

// ============================================================================
// Section 6: Types & Interfaces — Method Values
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Extracting methods as values
//   - Binding receiver to method
//   - Using method values as callbacks
//
// RUN: go run ./01-foundations/06-types-and-interfaces/15-method-values
// ============================================================================

type Counter struct {
	Value int
}

func (c *Counter) Increment() {
	c.Value++
}

func (c *Counter) Decrement() {
	c.Value--
}

func (c *Counter) GetValue() int {
	return c.Value
}

type Handler struct {
	Name string
}

func (h *Handler) OnClick() {
	fmt.Printf("  %s clicked!\n", h.Name)
}

func (h *Handler) OnHover() {
	fmt.Printf("  %s hovered!\n", h.Name)
}

func runHandler(name string, handler func()) {
	fmt.Printf("Running handler for %s\n", name)
	handler()
}

func main() {
	fmt.Println("=== Method Values ===")
	fmt.Println()

	fmt.Println("--- Basic Method Value ---")
	counter := &Counter{Value: 10}
	inc := counter.Increment
	inc()
	inc()
	fmt.Printf("  After two increments: %d\n", counter.Value)

	fmt.Println()
	fmt.Println("--- Method Value as Callback ---")
	handler := &Handler{Name: "SubmitButton"}
	clickHandler := handler.OnClick
	runHandler("button", clickHandler)

	fmt.Println()
	fmt.Println("--- Storing Method Values in Map ---")
	handlers := map[string]func(){
		"click": handler.OnClick,
		"hover": handler.OnHover,
	}
	for _, fn := range handlers {
		fn()
	}

	fmt.Println()
	fmt.Println("--- Method Values with Different Receivers ---")
	c1 := &Counter{Value: 100}
	c2 := &Counter{Value: 200}
	fn1 := c1.Increment
	fn2 := c2.Increment
	fn1()
	fn2()
	fmt.Printf("  c1: %d, c2: %d\n", c1.Value, c2.Value)

	fmt.Println()
	fmt.Println("--- Method Value Preserves Receiver ---")
	original := &Counter{Value: 5}
	methodFn := original.Increment
	methodFn()
	methodFn()
	fmt.Printf("  Original after method value calls: %d\n", original.Value)

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Method values capture their receiver")
	fmt.Println("  - Use method values as callbacks and event handlers")
	fmt.Println("  - Each instance's method value is independent")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.16 complex generic constraints")
	fmt.Println("   Current: TI.15 (method values)")
	fmt.Println("---------------------------------------------------")
}
