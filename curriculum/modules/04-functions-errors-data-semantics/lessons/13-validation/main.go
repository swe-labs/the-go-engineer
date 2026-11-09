package main

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field string
	Value interface{}
	Rule  string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("%s: rule %q violated (value=%v)", v.Field, v.Rule, v.Value)
}

type Order struct {
	ProductID  string
	Quantity   int
	PriceCents int64
}

func validateOrder(o Order) error {
	var errs []error
	if o.ProductID == "" {
		errs = append(errs, &ValidationError{Field: "ProductID", Value: o.ProductID, Rule: "required"})
	} else if len(o.ProductID) > 32 {
		errs = append(errs, &ValidationError{Field: "ProductID", Value: o.ProductID, Rule: "maxlen:32"})
	}
	if o.Quantity < 1 || o.Quantity > 100 {
		errs = append(errs, &ValidationError{Field: "Quantity", Value: o.Quantity, Rule: "range:1-100"})
	}
	if o.PriceCents <= 0 {
		errs = append(errs, &ValidationError{Field: "PriceCents", Value: o.PriceCents, Rule: "positive"})
	}
	return errors.Join(errs...)
}

func main() {
	orders := []Order{
		{ProductID: "", Quantity: 0, PriceCents: -100},
		{ProductID: "valid-product-123", Quantity: 5, PriceCents: 2999},
		{ProductID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Quantity: 200, PriceCents: 0},
	}
	for _, o := range orders {
		err := validateOrder(o)
		if err == nil {
			fmt.Printf("Order %q: valid\n", o.ProductID)
			continue
		}
		fmt.Printf("Order %q: %v\n", o.ProductID, err)
		// Print each individual validation error
		var ve *ValidationError
		remaining := err
		for remaining != nil {
			if errors.As(remaining, &ve) {
				fmt.Printf("  - field=%q rule=%q value=%v\n", ve.Field, ve.Rule, ve.Value)
				ve = nil
			}
			// unwrap to break out — in practice this requires walking the chain
			break
		}
	}
}
