// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package events

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrInvalidEvent (Error): returned when an event is missing required routing fields
	ErrInvalidEvent = errors.New("invalid event")
	// ErrQueueFull (Error): returned by TryPublish when the bus buffer is saturated
	ErrQueueFull = errors.New("event queue full")
	// ErrBusClosed (Error): returned when attempting to publish to a permanently shut-down bus
	ErrBusClosed = errors.New("event bus closed")
)

// Bus (Struct): single-channel event bus with buffered channel, thread-safe and idempotent close
//
// Design note: the closed signal is carried solely by the `closed` channel.
// A non-blocking select on a closed channel is O(1) and allocation-free, so
// a separate boolean flag + mutex fast-path is unnecessary overhead.
type Bus struct {
	events chan Event
	closed chan struct{}
	once   sync.Once
	now    func() time.Time
}

// NewBus (Constructor): creates a new thread-safe event bus with the given channel buffer capacity
func NewBus(capacity int) *Bus {
	if capacity <= 0 {
		capacity = 1
	}

	return &Bus{
		events: make(chan Event, capacity),
		closed: make(chan struct{}),
		now:    time.Now,
	}
}

// Subscribe (Method): returns the read-only event channel for ordered delivery; channel is never closed by the bus
func (b *Bus) Subscribe() <-chan Event {
	if b == nil {
		return nil
	}

	return b.events
}

// TryPublish (Method): non-blocking publish; returns ErrQueueFull or ErrBusClosed immediately
func (b *Bus) TryPublish(event Event) error {
	if b == nil {
		return fmt.Errorf("event bus is not configured")
	}

	// Fast closed check - a receive on a closed channel is non-blocking and
	// allocation-free, so no mutex is needed.
	select {
	case <-b.closed:
		return ErrBusClosed
	default:
	}

	prepared, err := b.prepare(event)
	if err != nil {
		return err
	}

	select {
	case b.events <- prepared:
		return nil
	case <-b.closed:
		return ErrBusClosed
	default:
		return ErrQueueFull
	}
}

// Publish (Method): enqueues an event, blocking until accepted, bus closed, or context cancelled
func (b *Bus) Publish(ctx context.Context, event Event) error {
	if b == nil {
		return fmt.Errorf("event bus is not configured")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	// Fast closed check before the more expensive prepare call.
	select {
	case <-b.closed:
		return ErrBusClosed
	default:
	}

	prepared, err := b.prepare(event)
	if err != nil {
		return err
	}

	select {
	case b.events <- prepared:
		return nil
	case <-b.closed:
		return ErrBusClosed
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close (Method): permanently shuts the bus; idempotent and safe for concurrent use
//
// Note: b.events is intentionally NOT closed here. The bus does not own the
// consumer side of the channel, so closing it would risk a panic in any
// subscriber that attempts a receive after a hypothetical re-open, and would
// also race with in-flight sends from Publish. Subscribers should stop reading
// by watching their own done/context signal alongside the Subscribe channel.
func (b *Bus) Close() {
	if b == nil {
		return
	}
	// sync.Once guarantees close(b.closed) is called exactly once, making
	// Close safe to call concurrently from multiple goroutines.
	b.once.Do(func() {
		close(b.closed)
	})
}

// prepare (Method): validates and timestamps an event before publishing
func (b *Bus) prepare(event Event) (Event, error) {
	if event.Type == "" || event.TenantID <= 0 {
		return Event{}, ErrInvalidEvent
	}

	if event.Metadata != nil {
		metadataCopy := make(map[string]string, len(event.Metadata))
		for k, v := range event.Metadata {
			metadataCopy[k] = v
		}
		event.Metadata = metadataCopy
	}

	return event.WithOccurredAt(b.now().UTC()), nil
}
