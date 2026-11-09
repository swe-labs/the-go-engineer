package main

import (
	"errors"
	"fmt"
)

var ErrEmptyItem = errors.New("item cannot be empty")
var ErrNegativeQty = errors.New("quantity cannot be negative")

type ValidationError struct {
	Field string
	Rule  string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation: %s violates %q", e.Field, e.Rule)
}

func (e *ValidationError) Is(target error) bool {
	t, ok := target.(*ValidationError)
	if !ok {
		return false
	}
	return e.Field == t.Field && e.Rule == t.Rule
}

type OrderErrorCode int

const (
	ErrCodeEmptyOrder OrderErrorCode = iota + 1
	ErrCodeProcessing
)

type OrderError struct {
	Code OrderErrorCode
	Err  error
}

func (e *OrderError) Error() string {
	return fmt.Sprintf("order error [code=%d]: %v", e.Code, e.Err)
}

func (e *OrderError) Unwrap() error {
	return e.Err
}

func (e *OrderError) Is(target error) bool {
	t, ok := target.(*OrderError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

type Item struct {
	Name     string
	Quantity int
	Price    int
}

type OrderSummary struct {
	TotalItems int
	TotalPrice int
}

type Resource struct {
	locked bool
}

func acquire(name string) *Resource {
	fmt.Printf("acquired resource for %q\n", name)
	return &Resource{locked: true}
}

func (r *Resource) Close() error {
	if !r.locked {
		return nil
	}
	r.locked = false
	fmt.Println("released resource")
	return nil
}

func validateItem(it Item) error {
	var errs []error
	if it.Name == "" {
		errs = append(errs, &ValidationError{Field: "Name", Rule: "required"})
	}
	if it.Quantity < 0 {
		errs = append(errs, fmt.Errorf("item %q: %w", it.Name, ErrNegativeQty))
	}
	if it.Quantity == 0 {
		errs = append(errs, &ValidationError{Field: "Quantity", Rule: "min:1"})
	}
	if it.Price < 0 {
		errs = append(errs, &ValidationError{Field: "Price", Rule: "non_negative"})
	}
	return errors.Join(errs...)
}

func processOrder(items []Item) (OrderSummary, error) {
	if len(items) == 0 {
		return OrderSummary{}, &OrderError{Code: ErrCodeEmptyOrder, Err: errors.New("order has no items")}
	}

	res := acquire("order")
	defer res.Close()

	var summary OrderSummary
	for _, it := range items {
		if err := validateItem(it); err != nil {
			return OrderSummary{}, fmt.Errorf("process order item %q: %w", it.Name, err)
		}
		summary.TotalItems += it.Quantity
		summary.TotalPrice += it.Quantity * it.Price
	}
	return summary, nil
}

func main() {
	items := []Item{
		{Name: "widget", Quantity: 2, Price: 150},
		{Name: "gadget", Quantity: 1, Price: 350},
	}
	summary, err := processOrder(items)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Order: %d items, total price %d cents\n", summary.TotalItems, summary.TotalPrice)
}
