package main

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestDedupStoreRejectsDuplicateKey(t *testing.T) {
	store := NewDedupStore(nil, nil)
	if !store.TryProcess("key-1") {
		t.Fatal("expected first process to succeed")
	}
	if store.TryProcess("key-1") {
		t.Fatal("expected duplicate process to be rejected")
	}
}

func TestDedupStoreAllowsDifferentKeys(t *testing.T) {
	store := NewDedupStore(nil, nil)
	for _, k := range []string{"a", "b", "c"} {
		if !store.TryProcess(k) {
			t.Fatalf("expected first process of %q to succeed", k)
		}
	}
}

func TestDedupStoreCallbacks(t *testing.T) {
	var missCount, dupCount int32
	store := NewDedupStore(
		func(k string) { atomic.AddInt32(&missCount, 1) },
		func(k string) { atomic.AddInt32(&dupCount, 1) },
	)
	store.TryProcess("x")
	store.TryProcess("x")
	store.TryProcess("y")
	store.TryProcess("x")
	if missCount != 2 {
		t.Errorf("expected 2 misses, got %d", missCount)
	}
	if dupCount != 2 {
		t.Errorf("expected 2 dups, got %d", dupCount)
	}
}

func TestDedupStoreSeen(t *testing.T) {
	store := NewDedupStore(nil, nil)
	store.TryProcess("seen-key")
	if !store.Seen("seen-key") {
		t.Fatal("expected Seen to return true after processing")
	}
	if store.Seen("unseen") {
		t.Fatal("expected Seen to return false for unprocessed key")
	}
}

func TestDedupStoreCount(t *testing.T) {
	store := NewDedupStore(nil, nil)
	store.TryProcess("a")
	store.TryProcess("b")
	store.TryProcess("a")
	if store.Count() != 2 {
		t.Errorf("expected count 2, got %d", store.Count())
	}
}

func TestRetryWithBackoffSucceeds(t *testing.T) {
	var attempts int32
	op := func() error {
		atomic.AddInt32(&attempts, 1)
		if atomic.LoadInt32(&attempts) < 3 {
			return errors.New("not yet")
		}
		return nil
	}
	cfg := RetryConfig{MaxAttempts: 5, BaseDelay: 1 * time.Millisecond, MaxDelay: 5 * time.Millisecond}
	result := RetryWithBackoff(cfg, op)
	if !result.Success {
		t.Fatal("expected success")
	}
	if result.Attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", result.Attempts)
	}
}

func TestRetryWithBackoffFails(t *testing.T) {
	op := func() error {
		return errors.New("always fails")
	}
	cfg := RetryConfig{MaxAttempts: 3, BaseDelay: 1 * time.Millisecond, MaxDelay: 5 * time.Millisecond}
	result := RetryWithBackoff(cfg, op)
	if result.Success {
		t.Fatal("expected failure")
	}
	if result.Attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", result.Attempts)
	}
	if result.Err == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestRetryWithBackoffTable(t *testing.T) {
	tests := []struct {
		name    string
		maxAtt  int
		needOK  int
		wantAtt int
		wantOK  bool
	}{
		{name: "succeeds_first_try", maxAtt: 3, needOK: 1, wantAtt: 1, wantOK: true},
		{name: "succeeds_third_try", maxAtt: 5, needOK: 3, wantAtt: 3, wantOK: true},
		{name: "exhausts_all_attempts", maxAtt: 2, needOK: 5, wantAtt: 2, wantOK: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var count int32
			op := func() error {
				atomic.AddInt32(&count, 1)
				if int(atomic.LoadInt32(&count)) < tc.needOK {
					return errors.New("not yet")
				}
				return nil
			}
			cfg := RetryConfig{MaxAttempts: tc.maxAtt, BaseDelay: 1 * time.Millisecond, MaxDelay: 5 * time.Millisecond}
			result := RetryWithBackoff(cfg, op)
			if result.Success != tc.wantOK {
				t.Errorf("expected success=%v, got %v", tc.wantOK, result.Success)
			}
			if result.Attempts != tc.wantAtt {
				t.Errorf("expected attempts=%d, got %d", tc.wantAtt, result.Attempts)
			}
		})
	}
}
