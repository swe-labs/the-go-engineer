package main

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestPaymentPipeline_IdempotentReplay(t *testing.T) {
	// Core practice task requirement:
	// Send two identical payment requests with the same idempotency key and
	// verify that the second request returns the cached result without charging twice.
	gateway := &FlakyGatewaySimulator{TransientErrors: 0}
	pipeline := NewPaymentPipeline(gateway, 3, 5*time.Millisecond)

	ctx := context.Background()
	idemKey := "order-checkout-uuid-888"

	// 1. Initial Charge
	receipt1, err := pipeline.Charge(ctx, idemKey, "ord-1", 120.00)
	if err != nil {
		t.Fatalf("first charge failed: %v", err)
	}

	if atomic.LoadInt32(&gateway.CallCount) != 1 {
		t.Fatalf("expected gateway to be called exactly 1 time, got %d", gateway.CallCount)
	}

	// 2. Duplicate Replay Charge with identical key
	receipt2, err := pipeline.Charge(ctx, idemKey, "ord-1", 120.00)
	if err != nil {
		t.Fatalf("second idempotent charge failed: %v", err)
	}

	// Gateway must NOT have been called a second time
	if atomic.LoadInt32(&gateway.CallCount) != 1 {
		t.Fatalf("CRITICAL IDEMPOTENCY FAILURE: Gateway called %d times! Second request re-charged!", gateway.CallCount)
	}

	// Transaction ID must be identical
	if receipt1.TransactionID != receipt2.TransactionID {
		t.Errorf("expected transaction IDs to match: %s vs %s", receipt1.TransactionID, receipt2.TransactionID)
	}
}

func TestPaymentPipeline_TransientRetrySuccess(t *testing.T) {
	// Gateway fails 2 times with transient network error, succeeds on 3rd attempt
	gateway := &FlakyGatewaySimulator{TransientErrors: 2}
	pipeline := NewPaymentPipeline(gateway, 3, 5*time.Millisecond)

	receipt, err := pipeline.Charge(context.Background(), "idem-retry-1", "ord-retry", 50.0)
	if err != nil {
		t.Fatalf("expected charge to succeed after retries, got: %v", err)
	}

	if atomic.LoadInt32(&gateway.CallCount) != 3 {
		t.Errorf("expected 3 attempts (2 failures + 1 success), got %d", gateway.CallCount)
	}
	if receipt.Status != "COMPLETED" {
		t.Errorf("expected receipt status COMPLETED, got %s", receipt.Status)
	}
}

func TestPaymentPipeline_NonTransientFailure(t *testing.T) {
	// Non-retryable error (card declined) should abort immediately
	gateway := &FlakyGatewaySimulator{PermanentFail: true}
	pipeline := NewPaymentPipeline(gateway, 3, 5*time.Millisecond)

	_, err := pipeline.Charge(context.Background(), "idem-fail-1", "ord-bad-card", 75.0)
	if !errors.Is(err, ErrCardDeclined) {
		t.Fatalf("expected ErrCardDeclined, got: %v", err)
	}

	if atomic.LoadInt32(&gateway.CallCount) != 1 {
		t.Errorf("expected non-retryable error to stop on 1st attempt, got %d", gateway.CallCount)
	}
}
