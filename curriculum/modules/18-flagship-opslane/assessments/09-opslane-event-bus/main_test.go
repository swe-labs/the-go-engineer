package main

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestEventBus_SubscribeAndPublishAsync(t *testing.T) {
	// Core practice task requirement:
	// Add a new subscriber to the event bus, publish an event,
	// and write a test verifying that the subscriber receives the event asynchronously.
	bus := NewEventBus()
	receivedCh := make(chan Event, 1)

	bus.Subscribe("order.created", func(ctx context.Context, e Event) error {
		receivedCh <- e
		return nil
	})

	testEvent := BaseEvent{
		EventType: "order.created",
		Tenant:    "tenant-acme",
		Data:      "ord-test-99",
		Timestamp: time.Now(),
	}

	bus.Publish(context.Background(), testEvent)

	select {
	case evt := <-receivedCh:
		if evt.Type() != "order.created" || evt.TenantID() != "tenant-acme" || evt.Payload() != "ord-test-99" {
			t.Errorf("unexpected event payload: %+v", evt)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timed out waiting for async event delivery")
	}

	_ = bus.Close(context.Background())
}

func TestEventBus_MultipleSubscribers(t *testing.T) {
	bus := NewEventBus()
	var sub1Received, sub2Received atomic.Bool

	bus.Subscribe("payment.processed", func(ctx context.Context, e Event) error {
		sub1Received.Store(true)
		return nil
	})

	bus.Subscribe("payment.processed", func(ctx context.Context, e Event) error {
		sub2Received.Store(true)
		return nil
	})

	bus.Publish(context.Background(), BaseEvent{
		EventType: "payment.processed",
		Tenant:    "tenant-1",
		Timestamp: time.Now(),
	})

	time.Sleep(50 * time.Millisecond)

	if !sub1Received.Load() || !sub2Received.Load() {
		t.Errorf("expected both subscribers to receive the event, got sub1=%t, sub2=%t", sub1Received.Load(), sub2Received.Load())
	}

	_ = bus.Close(context.Background())
}

func TestEventBus_PanicRecovery(t *testing.T) {
	bus := NewEventBus()
	var siblingReceived atomic.Bool

	// Sub 1 panics
	bus.Subscribe("user.registered", func(ctx context.Context, e Event) error {
		panic("simulated fatal subscriber failure")
	})

	// Sub 2 should still run safely
	bus.Subscribe("user.registered", func(ctx context.Context, e Event) error {
		siblingReceived.Store(true)
		return nil
	})

	bus.Publish(context.Background(), BaseEvent{
		EventType: "user.registered",
		Tenant:    "tenant-1",
		Timestamp: time.Now(),
	})

	time.Sleep(50 * time.Millisecond)

	if !siblingReceived.Load() {
		t.Errorf("expected sibling subscriber to process event despite prior panic")
	}

	_ = bus.Close(context.Background())
}

func TestEventBus_Unsubscribe(t *testing.T) {
	bus := NewEventBus()
	var callCount atomic.Int32

	subID := bus.Subscribe("invoice.sent", func(ctx context.Context, e Event) error {
		callCount.Add(1)
		return nil
	})

	bus.Publish(context.Background(), BaseEvent{EventType: "invoice.sent"})
	time.Sleep(20 * time.Millisecond)

	if callCount.Load() != 1 {
		t.Fatalf("expected 1 call before unsubscribe, got %d", callCount.Load())
	}

	bus.Unsubscribe("invoice.sent", subID)

	bus.Publish(context.Background(), BaseEvent{EventType: "invoice.sent"})
	time.Sleep(20 * time.Millisecond)

	if callCount.Load() != 1 {
		t.Errorf("expected no additional calls after unsubscribe, got %d", callCount.Load())
	}

	_ = bus.Close(context.Background())
}
