package main

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestCancellationStopsGoroutine(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var count atomic.Int32

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				count.Add(1)
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()

	time.Sleep(5 * time.Millisecond)
	cancel()
	time.Sleep(5 * time.Millisecond)

	before := count.Load()
	time.Sleep(10 * time.Millisecond)
	after := count.Load()

	if before != after {
		t.Errorf("goroutine continued after cancel: before=%d after=%d", before, after)
	}
}

func TestErrAfterCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if ctx.Err() != context.Canceled {
		t.Errorf("Err() = %v, want Canceled", ctx.Err())
	}
}

func TestCancellationIsBroadcast(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	results := make(chan error, 3)
	for i := 0; i < 3; i++ {
		go func() {
			<-ctx.Done()
			results <- ctx.Err()
		}()
	}

	cancel()

	for i := 0; i < 3; i++ {
		if err := <-results; err != context.Canceled {
			t.Errorf("goroutine %d: got %v, want Canceled", i, err)
		}
	}
}

func TestGracefulShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		defer close(done)
		<-ctx.Done()
		time.Sleep(5 * time.Millisecond)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(50 * time.Millisecond):
		t.Fatal("goroutine did not shut down gracefully")
	}
}
