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
	summary, err := processOrder(items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.TotalItems != 3 {
		t.Errorf("expected 3 total items, got %d", summary.TotalItems)
	}
	if summary.TotalPrice != 650 {
		t.Errorf("expected 650 total price, got %d", summary.TotalPrice)
	}
}

func TestProcessOrder_EmptyOrder(t *testing.T) {
	_, err := processOrder(nil)
	if err == nil {
		t.Fatalf("expected error for empty order, got nil")
	}

	var orderErr *OrderError
	if !errors.As(err, &orderErr) {
		t.Fatalf("expected error to be of type *OrderError, got %T", err)
	}
	if orderErr.Code != ErrCodeEmptyOrder {
		t.Errorf("expected ErrCodeEmptyOrder, got %d", orderErr.Code)
	}
}

func TestValidateItem_ValidationErrorExtraction(t *testing.T) {
	item := Item{Name: "", Quantity: 1, Price: 100}
	err := validateItem(item)
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected errors.As to extract *ValidationError, got %T", err)
	}
	if valErr.Field != "Name" || valErr.Rule != "required" {
		t.Errorf("unexpected validation error details: %+v", valErr)
	}
}

func TestValidateItem_SentinelMatching(t *testing.T) {
	item := Item{Name: "widget", Quantity: -5, Price: 100}
	err := validateItem(item)
	if err == nil {
		t.Fatalf("expected error for negative quantity, got nil")
	}

	if !errors.Is(err, ErrNegativeQty) {
		t.Errorf("expected errors.Is to match ErrNegativeQty")
	}
}

func TestValidateItem_MultipleErrorsJoin(t *testing.T) {
	// Item with empty name, negative quantity, and negative price
	item := Item{Name: "", Quantity: -2, Price: -50}
	err := validateItem(item)
	if err == nil {
		t.Fatalf("expected multi-error, got nil")
	}

	// Should match ErrNegativeQty
	if !errors.Is(err, ErrNegativeQty) {
		t.Errorf("expected errors.Is to match ErrNegativeQty inside joined errors")
	}

	// Should also extract ValidationError for Name
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Errorf("expected errors.As to extract ValidationError from joined errors")
	}
}
