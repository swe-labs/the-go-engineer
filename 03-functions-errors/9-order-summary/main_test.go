package main

import (
	"strings"
	"testing"
)

func TestValidateOrderNameRejectsBlank(t *testing.T) {
	err := validateOrderName("   ")
	if err == nil {
		t.Fatal("expected blank order name to fail validation")
	}
}

func TestValidatePricesRejectsNegative(t *testing.T) {
	err := validatePrices([]int{10, -2, 5})
	if err == nil {
		t.Fatal("expected negative price to fail validation")
	}
}

func TestProcessOrderSuccess(t *testing.T) {
	summary, err := processOrder("starter cart", []int{12, 18, 25}, 10)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if !strings.Contains(summary, "adjusted subtotal: 55") {
		t.Fatalf("expected adjusted subtotal 55 in summary, got %q", summary)
	}

	if !strings.Contains(summary, "total: 65") {
		t.Fatalf("expected total 65 in summary, got %q", summary)
	}
}

func TestProcessOrderAppliesPricingRule(t *testing.T) {
	discount := makeMinimumSubtotalDiscount(50, 5)

	summary, err := processOrder("starter cart", []int{12, 18, 25}, 10, discount)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if !strings.Contains(summary, "adjusted subtotal: 50") {
		t.Fatalf("expected adjusted subtotal 50 in summary, got %q", summary)
	}

	if !strings.Contains(summary, "total: 60") {
		t.Fatalf("expected discounted total 60 in summary, got %q", summary)
	}
}

func TestProcessOrderRejectsNegativeShipping(t *testing.T) {
	_, err := processOrder("starter cart", []int{12, 18, 25}, -1)
	if err == nil {
		t.Fatal("expected negative shipping to fail")
	}
}
