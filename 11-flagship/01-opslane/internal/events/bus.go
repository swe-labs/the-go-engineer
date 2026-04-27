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
	ErrInvalidEvent = errors.New("invalid event")
	ErrQueueFull    = errors.New("event queue full")
	ErrBusClosed    = errors.New("event bus closed")
)

type Bus struct {
	mu     sync.RWMutex
	events chan Event
	closed bool
	now    func() time.Time
}

func NewBus(capacity int) *Bus {
	if capacity <= 0 {
		capacity = 1
	}

	return &Bus{
		events: make(chan Event, capacity),
		now:    time.Now,
	}
}

func (b *Bus) Subscribe() <-chan Event {
	if b == nil {
		return nil
	}

	return b.events
}

func (b *Bus) TryPublish(event Event) error {
	if b == nil {
		return fmt.Errorf("event bus is not configured")
	}

	event, err := b.prepare(event)
	if err != nil {
		return err
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return ErrBusClosed
	}

	select {
	case b.events <- event:
		return nil
	default:
		return ErrQueueFull
	}
}

func (b *Bus) Publish(ctx context.Context, event Event) error {
	if b == nil {
		return fmt.Errorf("event bus is not configured")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	event, err := b.prepare(event)
	if err != nil {
		return err
	}

	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return ErrBusClosed
	}

	select {
	case b.events <- event:
		b.mu.RUnlock()
		return nil
	default:
	}
	b.mu.RUnlock()

	select {
	case b.events <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *Bus) Close() {
	if b == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	b.closed = true
	close(b.events)
}

func (b *Bus) prepare(event Event) (Event, error) {
	if event.Type == "" || event.TenantID <= 0 {
		return Event{}, ErrInvalidEvent
	}

	return event.WithOccurredAt(b.now().UTC()), nil
}
