package main

import "testing"

func TestNewItem(t *testing.T) {
	item := NewItem("Test", 1234)
	if item.Name != "Test" {
		t.Errorf("NewItem name = %q; want %q", item.Name, "Test")
	}
	if item.PriceCents != 1234 {
		t.Errorf("NewItem PriceCents = %d; want %d", item.PriceCents, 1234)
	}
}

func TestNewItemZero(t *testing.T) {
	item := NewItem("", 0)
	if item.Name != "" {
		t.Errorf("NewItem name = %q; want %q", item.Name, "")
	}
	if item.PriceCents != 0 {
		t.Errorf("NewItem PriceCents = %d; want %d", item.PriceCents, 0)
	}
}

func TestTaxRate(t *testing.T) {
	tests := []struct {
		country string
		want    float64
	}{
		{"US", 0.07},
		{"GB", 0.20},
		{"DE", 0.19},
		{"JP", 0.10},
		{"CA", 0.0},
		{"", 0.0},
	}
	for _, tc := range tests {
		got := TaxRate(tc.country)
		if got != tc.want {
			t.Errorf("TaxRate(%q) = %f; want %f", tc.country, got, tc.want)
		}
	}
}

func TestDiscountPercent(t *testing.T) {
	tests := []struct {
		quantity int
		want     float64
	}{
		{200, 0.15},
		{100, 0.15},
		{75, 0.10},
		{50, 0.10},
		{25, 0.05},
		{10, 0.05},
		{5, 0.0},
		{1, 0.0},
		{0, 0.0},
	}
	for _, tc := range tests {
		got := DiscountPercent(tc.quantity)
		if got != tc.want {
			t.Errorf("DiscountPercent(%d) = %f; want %f", tc.quantity, got, tc.want)
		}
	}
}

func TestFormatPrice(t *testing.T) {
	tests := []struct {
		cents int
		want  string
	}{
		{0, "$0.00"},
		{1, "$0.01"},
		{99, "$0.99"},
		{100, "$1.00"},
		{1234, "$12.34"},
		{999999, "$9999.99"},
		{-1, "-$0.01"},
		{-1234, "-$12.34"},
	}
	for _, tc := range tests {
		got := FormatPrice(tc.cents)
		if got != tc.want {
			t.Errorf("FormatPrice(%d) = %q; want %q", tc.cents, got, tc.want)
		}
	}
}

func TestCalculateTotal(t *testing.T) {
	items := []Item{
		NewItem("A", 1000),
		NewItem("B", 2000),
	}
	total, err := CalculateTotal(items, "US")
	if err != nil {
		t.Fatalf("CalculateTotal unexpected error: %v", err)
	}
	// (1000 + 2000) * 1.07 = 3210
	if total != 3210 {
		t.Errorf("CalculateTotal = %d; want 3210", total)
	}
}

func TestCalculateTotalEmpty(t *testing.T) {
	_, err := CalculateTotal(nil, "US")
	if err == nil {
		t.Error("CalculateTotal(nil) expected error")
	}
	_, err = CalculateTotal([]Item{}, "US")
	if err == nil {
		t.Error("CalculateTotal(empty) expected error")
	}
}

func TestCalculateTotalNoTax(t *testing.T) {
	items := []Item{NewItem("X", 5000)}
	total, err := CalculateTotal(items, "XX")
	if err != nil {
		t.Fatalf("CalculateTotal unexpected error: %v", err)
	}
	if total != 5000 {
		t.Errorf("CalculateTotal = %d; want 5000", total)
	}
}

func TestCalculateTotalBulkDiscount(t *testing.T) {
	items := make([]Item, 100)
	for i := 0; i < 100; i++ {
		items[i] = NewItem("Item", 100)
	}
	total, err := CalculateTotal(items, "XX")
	if err != nil {
		t.Fatalf("CalculateTotal unexpected error: %v", err)
	}
	// 100 * 100 = 10000, discount 15% = 1500 => 8500, no tax
	if total != 8500 {
		t.Errorf("CalculateTotal bulk = %d; want 8500", total)
	}
}
