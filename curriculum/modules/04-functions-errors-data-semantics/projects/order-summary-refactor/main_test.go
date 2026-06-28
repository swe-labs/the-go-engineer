package main

import (
	"errors"
	"testing"
)

func TestProcessOrder_Success(t *testing.T) {
	items := []Item{
		{Name: "widget", Quantity: 2, Price: 150},
		{Name: "gadget", Quantity: 1, Price: 350},
	}
	processOrder(items)
}

func TestProcessOrder_EmptyItemsPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty items")
		}
	}()
	processOrder(nil)
}

func TestProcessOrder_EmptyNamePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty name")
		}
	}()
	processOrder([]Item{{Name: "", Quantity: 1, Price: 100}})
}

func TestProcessOrder_NegativeQuantityPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for negative quantity")
		}
	}()
	processOrder([]Item{{Name: "x", Quantity: -1, Price: 100}})
}

func TestValidateOrder_Success(t *testing.T) {
	err := validateOrder([]Item{{Name: "a", Quantity: 1, Price: 100}})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateOrder_EmptyName(t *testing.T) {
	err := validateOrder([]Item{{Name: "", Quantity: 1, Price: 100}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrEmptyItem) {
		t.Error("expected errors.Is to match ErrEmptyItem")
	}
}

func TestValidateOrder_NegativeQuantity(t *testing.T) {
	err := validateOrder([]Item{{Name: "a", Quantity: -1, Price: 100}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNegativeQty) {
		t.Error("expected errors.Is to match ErrNegativeQty")
	}
}

func TestValidateOrder_NoErrorsForValid(t *testing.T) {
	tests := []struct {
		name  string
		items []Item
	}{
		{"single valid item", []Item{{Name: "a", Quantity: 1, Price: 100}}},
		{"multiple valid items", []Item{{Name: "a", Quantity: 2, Price: 50}, {Name: "b", Quantity: 1, Price: 200}}},
		{"zero price is allowed", []Item{{Name: "free", Quantity: 1, Price: 0}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOrder(tt.items)
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}
		})
	}
}
