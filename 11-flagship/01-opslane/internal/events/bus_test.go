package events

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBusPublishesEventWithTimestamp(t *testing.T) {
	t.Parallel()

	bus := NewBus(1)
	bus.now = func() time.Time {
		return time.Date(2026, 4, 27, 10, 0, 0, 0, time.UTC)
	}

	if err := bus.TryPublish(Event{Type: TypeOrderCreated, TenantID: 7, OrderID: 101}); err != nil {
		t.Fatalf("TryPublish returned error: %v", err)
	}

	event := <-bus.Subscribe()
	if event.OccurredAt.IsZero() {
		t.Fatal("OccurredAt should be set")
	}
	if event.OrderID != 101 {
		t.Fatalf("OrderID = %d, want 101", event.OrderID)
	}
}

func TestBusRejectsWhenQueueIsFull(t *testing.T) {
	t.Parallel()

	bus := NewBus(1)
	first := Event{Type: TypeOrderCreated, TenantID: 7, OrderID: 101}
	second := Event{Type: TypePaymentRequested, TenantID: 7, OrderID: 101}

	if err := bus.TryPublish(first); err != nil {
		t.Fatalf("first TryPublish returned error: %v", err)
	}

	err := bus.TryPublish(second)
	if !errors.Is(err, ErrQueueFull) {
		t.Fatalf("second TryPublish error = %v, want ErrQueueFull", err)
	}
}

func TestBusPublishHonorsContextCancellation(t *testing.T) {
	t.Parallel()

	bus := NewBus(1)
	if err := bus.TryPublish(Event{Type: TypeOrderCreated, TenantID: 7}); err != nil {
		t.Fatalf("TryPublish returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := bus.Publish(ctx, Event{Type: TypePaymentRequested, TenantID: 7})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Publish error = %v, want context.Canceled", err)
	}
}

func TestBusRejectsPublishAfterClose(t *testing.T) {
	t.Parallel()

	bus := NewBus(1)
	bus.Close()

	err := bus.TryPublish(Event{Type: TypeOrderCreated, TenantID: 7})
	if !errors.Is(err, ErrBusClosed) {
		t.Fatalf("TryPublish error = %v, want ErrBusClosed", err)
	}
}
