package main

import (
	"errors"
	"sync"
	"testing"
)

func TestBug1_ConcurrentCounter(t *testing.T) {
	counter := NewSafeCounter()

	const workers = 50
	const iterations = 100
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				counter.Inc()
			}
		}()
	}

	wg.Wait()

	expected := workers * iterations
	if counter.Value() != expected {
		t.Fatalf("expected counter %d, got %d", expected, counter.Value())
	}
}

func TestBug2_NilMapCache(t *testing.T) {
	// Zero-value cache without constructor initialization
	var cache KeyValueCache

	// Must not panic on write to nil map
	cache.Set("service", "billing")

	val, ok := cache.Get("service")
	if !ok || val != "billing" {
		t.Fatalf("expected 'billing', got %q (ok=%v)", val, ok)
	}
}

func TestBug3_SlicePaginationBounds(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}

	testCases := []struct {
		name      string
		offset    int
		limit     int
		wantItems []string
	}{
		{
			name:      "normal page",
			offset:    1,
			limit:     2,
			wantItems: []string{"b", "c"},
		},
		{
			name:      "limit exceeding slice length",
			offset:    3,
			limit:     10,
			wantItems: []string{"d", "e"},
		},
		{
			name:      "offset beyond slice length",
			offset:    10,
			limit:     5,
			wantItems: []string{},
		},
		{
			name:      "negative offset and limit",
			offset:    -2,
			limit:     -5,
			wantItems: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := PaginateSlice(items, tc.offset, tc.limit)
			if len(got) != len(tc.wantItems) {
				t.Fatalf("expected length %d, got %d (%v)", len(tc.wantItems), len(got), got)
			}
			for i := range got {
				if got[i] != tc.wantItems[i] {
					t.Fatalf("item at index %d = %q; want %q", i, got[i], tc.wantItems[i])
				}
			}
		})
	}
}

func TestBug4_ShadowedError(t *testing.T) {
	// Valid user
	user, err := ValidateAndFormatUser("u1", "Grace")
	if err != nil {
		t.Fatalf("unexpected error for valid user: %v", err)
	}
	if user.Name != "User: Grace" {
		t.Fatalf("expected formatted name 'User: Grace', got %q", user.Name)
	}

	// Empty ID returns ErrInvalidInput
	_, errBadID := ValidateAndFormatUser("", "Grace")
	if !errors.Is(errBadID, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for empty ID, got %v", errBadID)
	}

	// Empty Name must return ErrInvalidInput and NOT swallow error via shadowing!
	_, errBadName := ValidateAndFormatUser("u2", "")
	if errBadName == nil {
		t.Fatal("expected error for empty name, but got nil (error was swallowed by shadowing)")
	}
	if !errors.Is(errBadName, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for empty name, got %v", errBadName)
	}
}
