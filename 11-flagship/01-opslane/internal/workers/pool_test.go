package workers

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/events"
)

func TestPoolRejectsWhenQueueIsFull(t *testing.T) {
	t.Parallel()

	pool, err := NewPool(PoolConfig{
		Name:      "test",
		Workers:   1,
		QueueSize: 1,
		Handler:   func(context.Context, events.Event) error { return nil },
	})
	if err != nil {
		t.Fatalf("NewPool returned error: %v", err)
	}

	if err := pool.TrySubmit(events.Event{Type: events.TypeOrderCreated, TenantID: 7}); err != nil {
		t.Fatalf("first TrySubmit returned error: %v", err)
	}

	err = pool.TrySubmit(events.Event{Type: events.TypePaymentRequested, TenantID: 7})
	if !errors.Is(err, ErrQueueFull) {
		t.Fatalf("second TrySubmit error = %v, want ErrQueueFull", err)
	}

	pool.Stop()
}

func TestPoolDrainsQueuedWorkOnStop(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	processed := 0

	pool, err := NewPool(PoolConfig{
		Name:      "drain",
		Workers:   2,
		QueueSize: 4,
		Handler: func(context.Context, events.Event) error {
			mu.Lock()
			processed++
			mu.Unlock()
			return nil
		},
	})
	if err != nil {
		t.Fatalf("NewPool returned error: %v", err)
	}
	if err := pool.Start(context.Background()); err != nil {
		t.Fatalf("Start returned error: %v", err)
	}

	for i := 0; i < 4; i++ {
		if err := pool.Submit(context.Background(), events.Event{Type: events.TypeOrderCreated, TenantID: 7}); err != nil {
			t.Fatalf("Submit returned error: %v", err)
		}
	}

	pool.Stop()

	mu.Lock()
	defer mu.Unlock()
	if processed != 4 {
		t.Fatalf("processed = %d, want 4", processed)
	}
}

func TestPoolReportsHandlerErrors(t *testing.T) {
	t.Parallel()

	handlerErr := errors.New("worker failed")
	errs := make(chan error, 1)

	pool, err := NewPool(PoolConfig{
		Name:      "errors",
		Workers:   1,
		QueueSize: 1,
		Handler: func(context.Context, events.Event) error {
			return handlerErr
		},
		OnError: func(_ events.Event, err error) {
			errs <- err
		},
	})
	if err != nil {
		t.Fatalf("NewPool returned error: %v", err)
	}
	if err := pool.Start(context.Background()); err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	if err := pool.Submit(context.Background(), events.Event{Type: events.TypeOrderCreated, TenantID: 7}); err != nil {
		t.Fatalf("Submit returned error: %v", err)
	}

	select {
	case err := <-errs:
		if !errors.Is(err, handlerErr) {
			t.Fatalf("reported error = %v, want handlerErr", err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for worker error")
	}

	pool.Stop()
}

func TestPoolRejectsSubmitAfterStop(t *testing.T) {
	t.Parallel()

	pool, err := NewPool(PoolConfig{
		Name:      "stopped",
		Workers:   1,
		QueueSize: 1,
		Handler:   func(context.Context, events.Event) error { return nil },
	})
	if err != nil {
		t.Fatalf("NewPool returned error: %v", err)
	}

	pool.Stop()

	err = pool.TrySubmit(events.Event{Type: events.TypeOrderCreated, TenantID: 7})
	if !errors.Is(err, ErrPoolStopped) {
		t.Fatalf("TrySubmit error = %v, want ErrPoolStopped", err)
	}
}

// TestSubmitToStoppedPool verifies that Submit returns ErrPoolStopped when trying to submit to a stopped pool
func TestSubmitToStoppedPool(t *testing.T) {
	t.Parallel()

	pool, err := NewPool(PoolConfig{
		Name:      "submit-stopped",
		Workers:   1,
		QueueSize: 1,
		Handler:   func(context.Context, events.Event) error { return nil },
	})
	if err != nil {
		t.Fatalf("NewPool returned error: %v", err)
	}

	// Stop the pool
	pool.Stop()

	// Try to submit to the stopped pool
	err = pool.Submit(context.Background(), events.Event{Type: events.TypeOrderCreated, TenantID: 7})
	if !errors.Is(err, ErrPoolStopped) {
		t.Fatalf("Submit error = %v, want ErrPoolStopped", err)
	}
}

// TestSubmitToStoppedPoolWithContext verifies that Submit respects context cancellation even when pool is stopped
func TestSubmitToStoppedPoolWithContext(t *testing.T) {
	t.Parallel()

	pool, err := NewPool(PoolConfig{
		Name:      "submit-stopped-context",
		Workers:   1,
		QueueSize: 1,
		Handler:   func(context.Context, events.Event) error { return nil },
	})
	if err != nil {
		t.Fatalf("NewPool returned error: %v", err)
	}

	// Stop the pool
	pool.Stop()

	// Create a context that gets cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Try to submit to the stopped pool with a cancellable context
	err = pool.Submit(ctx, events.Event{Type: events.TypeOrderCreated, TenantID: 7})
	// Should return ErrPoolStopped because the pool is stopped, not context.Canceled
	if !errors.Is(err, ErrPoolStopped) {
		t.Fatalf("Submit error = %v, want ErrPoolStopped", err)
	}
}

// TestPublishToClosedBus verifies that Publish returns ErrBusClosed when trying to publish to a closed bus
func TestPublishToClosedBus(t *testing.T) {
	t.Parallel()

	bus := events.NewBus(1)
	bus.Close()

	err := bus.Publish(context.Background(), events.Event{Type: events.TypeOrderCreated, TenantID: 7})
	if !errors.Is(err, events.ErrBusClosed) {
		t.Fatalf("Publish error = %v, want ErrBusClosed", err)
	}
}

// TestPublishToClosedBusWithContext verifies that Publish respects context cancellation even when bus is closed
func TestPublishToClosedBusWithContext(t *testing.T) {
	t.Parallel()

	bus := events.NewBus(1)
	bus.Close()

	// Create a context that gets cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Try to publish to the closed bus with a cancellable context
	err := bus.Publish(ctx, events.Event{Type: events.TypeOrderCreated, TenantID: 7})
	// Should return ErrBusClosed because the bus is closed, not context.Canceled
	if !errors.Is(err, events.ErrBusClosed) {
		t.Fatalf("Publish error = %v, want ErrBusClosed", err)
	}
}
