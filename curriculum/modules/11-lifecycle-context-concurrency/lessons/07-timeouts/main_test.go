package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTimeoutExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := slowOperation(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestTimeoutNotExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- slowOperation(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("operation did not complete within 3s")
	}
}

func TestTimeoutCancelBeforeExpiry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	cancel()

	err := slowOperation(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected Canceled, got %v", err)
	}
}

func TestDeadlineExceededError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), -1*time.Nanosecond)
	defer cancel()

	<-ctx.Done()
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Errorf("Err() = %v, want DeadlineExceeded", ctx.Err())
	}
}
