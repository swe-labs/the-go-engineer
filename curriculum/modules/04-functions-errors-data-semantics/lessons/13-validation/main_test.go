package main

import (
	"errors"
	"testing"
)

func TestValidateOrder_Valid(t *testing.T) {
	err := validateOrder(Order{ProductID: "abc", Quantity: 5, PriceCents: 2999})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidateOrder_AllInvalid(t *testing.T) {
	err := validateOrder(Order{ProductID: "", Quantity: 0, PriceCents: -1})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	count := 0
	var ve *ValidationError
	// Check that joined errors contain at least one ValidationError
	if errors.As(err, &ve) {
		count++
	}
	if count == 0 {
		t.Error("expected at least one ValidationError")
	}
}

func TestValidateOrder_ProductIDTooLong(t *testing.T) {
	longID := ""
	for i := 0; i < 40; i++ {
		longID += "a"
	}
	err := validateOrder(Order{ProductID: longID, Quantity: 1, PriceCents: 100})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ve *ValidationError
	if errors.As(err, &ve) {
		if ve.Field != "ProductID" {
			t.Errorf("expected ProductID field, got %s", ve.Field)
		}
	}
}

func TestValidateOrder_QuantityOutOfRange(t *testing.T) {
	err := validateOrder(Order{ProductID: "abc", Quantity: 200, PriceCents: 100})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestValidateOrder_PriceCentsZero(t *testing.T) {
	err := validateOrder(Order{ProductID: "abc", Quantity: 1, PriceCents: 0})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ve *ValidationError
	if errors.As(err, &ve) {
		if ve.Rule != "positive" {
			t.Errorf("expected rule 'positive', got %q", ve.Rule)
		}
	}
}
