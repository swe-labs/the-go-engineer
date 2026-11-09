package main

import (
	"reflect"
	"testing"
	"time"
)

func TestDeduplicatePreserveOrder(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "consecutive and non-consecutive duplicates",
			input:    []string{"apple", "banana", "apple", "cherry", "banana", "date"},
			expected: []string{"apple", "banana", "cherry", "date"},
		},
		{
			name:     "all identical elements",
			input:    []string{"x", "x", "x", "x"},
			expected: []string{"x"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := DeduplicatePreserveOrder(tc.input)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("DeduplicatePreserveOrder(%v) = %v; want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestPartitionActive(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6}
	count := PartitionActive(numbers, func(n int) bool { return n%2 == 0 })

	if count != 3 {
		t.Fatalf("expected count 3, got %d", count)
	}

	for i := 0; i < count; i++ {
		if numbers[i]%2 != 0 {
			t.Errorf("expected even number at index %d, got %d", i, numbers[i])
		}
	}
}

func TestFindPriceTier(t *testing.T) {
	thresholds := []int{1, 10, 50, 100}
	prices := []float64{20.0, 15.0, 10.0, 5.0}

	tests := []struct {
		qty       int
		wantPrice float64
		wantFound bool
	}{
		{qty: 0, wantPrice: 0, wantFound: false},
		{qty: 1, wantPrice: 20.0, wantFound: true},
		{qty: 9, wantPrice: 20.0, wantFound: true},
		{qty: 10, wantPrice: 15.0, wantFound: true},
		{qty: 49, wantPrice: 15.0, wantFound: true},
		{qty: 50, wantPrice: 10.0, wantFound: true},
		{qty: 150, wantPrice: 5.0, wantFound: true},
	}

	for _, tc := range tests {
		price, found := FindPriceTier(thresholds, prices, tc.qty)
		if found != tc.wantFound || price != tc.wantPrice {
			t.Errorf("FindPriceTier(qty=%d) = (%v, %v); want (%v, %v)",
				tc.qty, price, found, tc.wantPrice, tc.wantFound)
		}
	}
}

func TestSortTransactions(t *testing.T) {
	t1 := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	txs := []Transaction{
		{ID: "low-late", Priority: 1, Timestamp: t3},
		{ID: "high-mid", Priority: 3, Timestamp: t2},
		{ID: "low-early", Priority: 1, Timestamp: t1},
		{ID: "high-early", Priority: 3, Timestamp: t1},
	}

	SortTransactions(txs)

	expectedOrder := []string{"high-early", "high-mid", "low-early", "low-late"}
	for i, wantID := range expectedOrder {
		if txs[i].ID != wantID {
			t.Errorf("at index %d: expected %s, got %s", i, wantID, txs[i].ID)
		}
	}
}
