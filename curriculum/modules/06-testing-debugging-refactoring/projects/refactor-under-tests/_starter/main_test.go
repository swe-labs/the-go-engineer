package main

import (
	"errors"
	"strings"
	"testing"
)

func TestProcessInvoice_ValidCalculations(t *testing.T) {
	testCases := []struct {
		name         string
		payload      string
		wantTotal    int
		wantSubtotal int
		wantDiscount int
		wantTax      int
	}{
		{
			name: "Gold customer in CA (15% discount, 9% tax)",
			payload: `ORD-001,CA,GOLD
Keyboard,2,5000
Mouse,1,3000`,
			// Subtotal: (2*5000) + 3000 = 13000 cents
			// Discount (15%): 13000 * 0.15 = 1950 cents
			// Taxable: 13000 - 1950 = 11050 cents
			// Tax (9%): int(11050 * 0.09) = 994 cents
			// Total: 11050 + 994 = 12044 cents
			wantSubtotal: 13000,
			wantDiscount: 1950,
			wantTax:      994,
			wantTotal:    12044,
		},
		{
			name: "Silver customer in NY (10% discount, 8% tax)",
			payload: `ORD-002,NY,SILVER
Monitor,1,20000`,
			// Subtotal: 20000 cents
			// Discount (10%): 2000 cents
			// Taxable: 18000 cents
			// Tax (8%): 18000 * 0.08 = 1440 cents
			// Total: 18000 + 1440 = 19440 cents
			wantSubtotal: 20000,
			wantDiscount: 2000,
			wantTax:      1440,
			wantTotal:    19440,
		},
		{
			name: "Standard customer in TX (0% discount, 6.25% tax)",
			payload: `ORD-003,TX,STANDARD
Cable,4,1000`,
			// Subtotal: 4000 cents
			// Discount: 0 cents
			// Tax (6.25%): int(4000 * 0.0625) = 250 cents
			// Total: 4250 cents
			wantSubtotal: 4000,
			wantDiscount: 0,
			wantTax:      250,
			wantTotal:    4250,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			summary, total, err := ProcessInvoice(tc.payload)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if total != tc.wantTotal {
				t.Fatalf("expected total %d cents, got %d", tc.wantTotal, total)
			}

			if !strings.Contains(summary, "Order") || !strings.Contains(summary, "Total") {
				t.Fatalf("malformed summary string: %q", summary)
			}
		})
	}
}

func TestProcessInvoice_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name    string
		payload string
		wantErr error
	}{
		{
			name:    "empty input",
			payload: "",
			wantErr: ErrInvalidInvoice,
		},
		{
			name:    "missing header columns",
			payload: "ORD-001,CA",
			wantErr: ErrInvalidInvoice,
		},
		{
			name: "negative quantity",
			payload: `ORD-001,CA,GOLD
Item,-2,1000`,
			wantErr: ErrInvalidQuantity,
		},
		{
			name: "negative price",
			payload: `ORD-001,CA,GOLD
Item,2,-1000`,
			wantErr: ErrInvalidPrice,
		},
		{
			name: "header only with no items",
			payload: `ORD-001,CA,GOLD`,
			wantErr: ErrInvalidInvoice,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := ProcessInvoice(tc.payload)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error wrapping %v, got %v", tc.wantErr, err)
			}
		})
	}
}
