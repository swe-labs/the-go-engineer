package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestCircuitBreakerClosedInitial(t *testing.T) {
	cb := NewCircuitBreaker(5, time.Second)
	if cb.State() != StateClosed {
		t.Errorf("expected initial state closed, got %v", cb.State())
	}
}

func TestCircuitBreakerOpensAfterThreshold(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Minute)

	for i := 0; i < 3; i++ {
		err := cb.Call(context.Background(), func(ctx context.Context) error {
			return errors.New("fail")
		})
		if err == nil {
			t.Fatalf("expected error on call %d", i+1)
		}
	}

	if cb.State() != StateOpen {
		t.Errorf("expected state open after threshold failures, got %v", cb.State())
	}
}

func TestCircuitBreakerRejectsWhenOpen(t *testing.T) {
	cb := NewCircuitBreaker(2, time.Minute)

	for i := 0; i < 2; i++ {
		cb.Call(context.Background(), func(ctx context.Context) error {
			return errors.New("fail")
		})
	}

	err := cb.Call(context.Background(), func(ctx context.Context) error {
		return nil
	})
	if err == nil || err.Error() != "circuit breaker open" {
		t.Errorf("expected 'circuit breaker open' error, got %v", err)
	}
}

func TestCircuitBreakerHalfOpenSuccess(t *testing.T) {
	cb := NewCircuitBreaker(3, 50*time.Millisecond)

	for i := 0; i < 3; i++ {
		cb.Call(context.Background(), func(ctx context.Context) error {
			return errors.New("fail")
		})
	}

	time.Sleep(60 * time.Millisecond)

	for i := 0; i < 4; i++ {
		err := cb.Call(context.Background(), func(ctx context.Context) error {
			return nil
		})
		if err != nil {
			t.Fatalf("expected success in half-open state, got %v", err)
		}
	}

	if cb.State() != StateClosed {
		t.Errorf("expected state closed after half-open successes, got %v", cb.State())
	}
}

func TestCircuitBreakerHalfOpenFailure(t *testing.T) {
	cb := NewCircuitBreaker(2, 50*time.Millisecond)

	for i := 0; i < 2; i++ {
		cb.Call(context.Background(), func(ctx context.Context) error {
			return errors.New("fail")
		})
	}

	time.Sleep(60 * time.Millisecond)

	err := cb.Call(context.Background(), func(ctx context.Context) error {
		return errors.New("still failing")
	})
	if err == nil {
		t.Fatal("expected error on half-open probe failure")
	}

	if cb.State() != StateOpen {
		t.Errorf("expected state open after half-open failure, got %v", cb.State())
	}
}

func TestBulkheadLimitsConcurrency(t *testing.T) {
	bh := NewBulkhead(2)

	var wg sync.WaitGroup
	started := make(chan struct{}, 3)
	block := make(chan struct{})

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()
			bh.Execute(ctx, func(ctx context.Context) error {
				started <- struct{}{}
				<-block
				return nil
			})
		}()
	}

	<-started
	<-started

	select {
	case <-started:
		t.Error("third call should not start (bulkhead limit is 2)")
	case <-time.After(50 * time.Millisecond):
	}

	close(block)
	wg.Wait()
}

func TestBulkheadContextCancellation(t *testing.T) {
	bh := NewBulkhead(1)

	bh.Execute(context.Background(), func(ctx context.Context) error {
		_ = ctx
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := bh.Execute(ctx, func(ctx context.Context) error {
		return nil
	})
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestResilientClientDegradation(t *testing.T) {
	callCount := 0
	client := NewResilientClient(
		func(ctx context.Context) (string, error) {
			callCount++
			return "", errors.New("primary failed")
		},
		func(ctx context.Context) (string, error) {
			return "cached", nil
		},
	)

	result, err := client.Fetch(context.Background())
	if err != nil {
		t.Fatalf("expected fallback to succeed, got %v", err)
	}
	if result != "cached (degraded)" {
		t.Errorf("expected 'cached (degraded)', got %q", result)
	}
}

func TestResilientClientPrimarySuccess(t *testing.T) {
	client := NewResilientClient(
		func(ctx context.Context) (string, error) {
			return "fresh", nil
		},
		func(ctx context.Context) (string, error) {
			return "cached", nil
		},
	)

	result, err := client.Fetch(context.Background())
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if result != "fresh" {
		t.Errorf("expected 'fresh', got %q", result)
	}
}

func TestBreakerIndependent(t *testing.T) {
	cb1 := NewCircuitBreaker(2, time.Minute)
	cb2 := NewCircuitBreaker(2, time.Minute)

	for i := 0; i < 2; i++ {
		cb1.Call(context.Background(), func(ctx context.Context) error {
			return errors.New("fail")
		})
	}

	if cb1.State() != StateOpen {
		t.Errorf("expected cb1 open, got %v", cb1.State())
	}
	if cb2.State() != StateClosed {
		t.Errorf("expected cb2 closed, got %v", cb2.State())
	}
}
