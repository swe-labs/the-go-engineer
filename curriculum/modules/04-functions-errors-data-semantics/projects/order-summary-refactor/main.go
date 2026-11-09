package main

import (
	"errors"
	"fmt"
)

// TODO: Replace flat errors.New with custom error types (ValidationError, OrderError).
// TODO: Add Unwrap() and Is() methods to custom types.
var ErrEmptyItem = errors.New("item cannot be empty")
var ErrNegativeQty = errors.New("quantity cannot be negative")

type Item struct {
	Name     string
	Quantity int
	Price    int
}

type Order struct {
	Items []Item
}

// TODO: Return (OrderSummary, error) instead of panicking.
// TODO: Use defer for cleanup simulation.
// TODO: Use fmt.Errorf with %w for all error returns.
// TODO: Use errors.Join to collect multiple validation errors.
// TODO: Replace panics with (result, error) return pattern.
func processOrder(items []Item) {
	if len(items) == 0 {
		panic("order must have at least one item")
	}
	for _, it := range items {
		if it.Name == "" {
			panic(errors.New("item name is required"))
		}
		if it.Quantity < 0 {
			panic(errors.New("item quantity cannot be negative"))
		}
		if it.Quantity == 0 {
			panic(errors.New("item quantity must be at least 1"))
		}
		if it.Price < 0 {
			panic(errors.New("item price cannot be negative"))
		}
	}
	totalItems := 0
	totalPrice := 0
	for _, it := range items {
		totalItems += it.Quantity
		totalPrice += it.Quantity * it.Price
	}
	fmt.Printf("Order: %d items, total price %d cents\n", totalItems, totalPrice)
}

// TODO: Add proper (OrderSummary, error) return signature.
// TODO: Use errors.Join to aggregate validation failures.
func validateOrder(items []Item) error {
	for _, it := range items {
		if it.Name == "" {
			return fmt.Errorf("validation failed: %w", ErrEmptyItem)
		}
		if it.Quantity < 0 {
			return fmt.Errorf("validation failed: %w", ErrNegativeQty)
		}
	}
	return nil
}

func main() {
	items := []Item{
		{Name: "widget", Quantity: 2, Price: 150},
		{Name: "gadget", Quantity: 1, Price: 350},
	}
	// TODO: Handle error instead of panicking inside processOrder.
	processOrder(items)
	err := validateOrder(items)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
