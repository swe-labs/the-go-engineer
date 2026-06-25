package main

import (
	"sync"
	"testing"
)

func TestIdempotencyStore(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      string
		wantSet    bool
		wantGetOK  bool
		wantGetVal string
	}{
		{"first set succeeds", "key1", "result1", true, true, "result1"},
		{"duplicate set fails", "key1", "result2", false, true, "result1"},
		{"new key succeeds", "key2", "result3", true, true, "result3"},
	}

	store := NewIdempotencyStore()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "first set succeeds" || tt.name == "new key succeeds" {
				if got := store.Set(tt.key, tt.value); got != tt.wantSet {
					t.Errorf("Set() = %v, want %v", got, tt.wantSet)
				}
			}
			gotVal, gotOK := store.Get(tt.key)
			if gotOK != tt.wantGetOK || gotVal != tt.wantGetVal {
				t.Errorf("Get() = (%q, %v), want (%q, %v)", gotVal, gotOK, tt.wantGetVal, tt.wantGetOK)
			}
		})
	}
}

func TestIdempotencyStoreConcurrent(t *testing.T) {
	store := NewIdempotencyStore()
	var wg sync.WaitGroup
	n := 100
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			store.Set("concurrent_key", "value")
		}()
	}
	wg.Wait()
	val, ok := store.Get("concurrent_key")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if val != "value" {
		t.Errorf("got %q, want 'value'", val)
	}
}

func TestProcessPayment(t *testing.T) {
	store := NewIdempotencyStore()
	req := PaymentRequest{ID: "pay_1", Amount: 100, Currency: "USD"}

	result1, dedup1 := processPayment(store, req)
	result2, dedup2 := processPayment(store, req)

	if dedup1 {
		t.Error("first call should not be dedup")
	}
	if result1 != result2 {
		t.Error("results should match")
	}
	if !dedup2 {
		t.Error("second call should be dedup")
	}
}
