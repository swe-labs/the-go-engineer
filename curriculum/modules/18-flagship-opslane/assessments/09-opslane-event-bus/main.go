package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Event interface represents an asynchronous domain notification.
type Event interface {
	Type() string
	TenantID() string
	Payload() any
	OccurredAt() time.Time
}

// EventHandler callback for incoming events.
type EventHandler func(ctx context.Context, e Event) error

// BaseEvent provides common event properties.
type BaseEvent struct {
	EventType string
	Tenant    string
	Data      any
	Timestamp time.Time
}

func (b BaseEvent) Type() string           { return b.EventType }
func (b BaseEvent) TenantID() string       { return b.Tenant }
func (b BaseEvent) Payload() any           { return b.Data }
func (b BaseEvent) OccurredAt() time.Time  { return b.Timestamp }

// Subscription represents an active event listener.
type Subscription struct {
	id      int
	topic   string
	handler EventHandler
}

// EventBus provides decoupled in-memory pub-sub messaging.
type EventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]Subscription
	nextSubID   int
	wg          sync.WaitGroup
	closed      bool
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]Subscription),
	}
}

// Subscribe registers a handler for a specific event topic.
func (b *EventBus) Subscribe(topic string, handler EventHandler) int {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.nextSubID++
	sub := Subscription{
		id:      b.nextSubID,
		topic:   topic,
		handler: handler,
	}

	b.subscribers[topic] = append(b.subscribers[topic], sub)
	return sub.id
}

// Unsubscribe removes a registered listener.
func (b *EventBus) Unsubscribe(topic string, subID int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subs, exists := b.subscribers[topic]
	if !exists {
		return
	}

	filtered := make([]Subscription, 0, len(subs))
	for _, s := range subs {
		if s.id != subID {
			filtered = append(filtered, s)
		}
	}
	b.subscribers[topic] = filtered
}

// Publish dispatches the event asynchronously to all subscribed listeners.
func (b *EventBus) Publish(ctx context.Context, e Event) {
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return
	}

	subs := append([]Subscription(nil), b.subscribers[e.Type()]...)
	// Also invoke wildcard subscribers if any
	wildcardSubs := append([]Subscription(nil), b.subscribers["*"]...)
	b.mu.RUnlock()

	allSubs := append(subs, wildcardSubs...)

	for _, sub := range allSubs {
		b.wg.Add(1)
		go func(s Subscription) {
			defer b.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[EventBus] Recovered panic in handler for topic %s: %v", s.topic, r)
				}
			}()

			_ = s.handler(ctx, e)
		}(sub)
	}
}

// Close waits for all in-flight event handlers to complete.
func (b *EventBus) Close(ctx context.Context) error {
	b.mu.Lock()
	b.closed = true
	b.mu.Unlock()

	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func main() {
	bus := NewEventBus()
	bus.Subscribe("order.created", func(ctx context.Context, e Event) error {
		fmt.Printf("Received async event %s for tenant %s\n", e.Type(), e.TenantID())
		return nil
	})

	bus.Publish(context.Background(), BaseEvent{
		EventType: "order.created",
		Tenant:    "tenant-acme",
		Data:      map[string]any{"order_id": "ord-100"},
		Timestamp: time.Now(),
	})

	_ = bus.Close(context.Background())
}
