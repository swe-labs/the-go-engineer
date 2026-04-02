// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "testing"

// ============================================================================
// Tests for: Pricing Calculator
// ============================================================================
//
// These tests verify the calculateItemPrice function handles
// regular items, sale items, and unknown items correctly.
//
// RUN: go test -v ./02-control-flow/4-pricing-calculator
// ============================================================================

func TestCalculateItemPrice(t *testing.T) {
	tests := []struct {
		name      string
		itemCode  string
		wantPrice float64
		wantFound bool
	}{
		{
			name:      "regular item found",
			itemCode:  "TSHIRT",
			wantPrice: 20.00,
			wantFound: true,
		},
		{
			name:      "sale item applies 10 percent discount",
			itemCode:  "MUG_SALE",
			wantPrice: 12.50 * 0.90, // 10% off base price
			wantFound: true,
		},
		{
			name:      "unknown item returns not found",
			itemCode:  "KEYBOARD",
			wantPrice: 0.0,
			wantFound: false,
		},
		{
			name:      "unknown sale item returns not found",
			itemCode:  "LAPTOP_SALE",
			wantPrice: 0.0,
			wantFound: false,
		},
		{
			name:      "another regular item",
			itemCode:  "BOOK",
			wantPrice: 25.99,
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrice, gotFound := calculateItemPrice(tt.itemCode)

			if gotFound != tt.wantFound {
				t.Errorf("calculateItemPrice(%q) found = %v, want %v",
					tt.itemCode, gotFound, tt.wantFound)
			}

			if gotPrice != tt.wantPrice {
				t.Errorf("calculateItemPrice(%q) price = %.2f, want %.2f",
					tt.itemCode, gotPrice, tt.wantPrice)
			}
		})
	}
}
