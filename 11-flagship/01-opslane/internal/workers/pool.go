// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package workers

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/swe-labs/the-go-engineer/11-flagship/01-opslane/internal/events"
)

var (
	// ErrInvalidPoolConfig is returned when PoolConfig validation fails.
	ErrInvalidPoolConfig = errors.New("invalid worker pool config")
	// ErrPoolStopped is returned when attempting to submit to a stopped pool.
	ErrPoolStopped = errors.New("worker pool stopped")
	// ErrQueueFull is returned by TrySubmit when the job channel has no capacity.
	ErrQueueFull = errors.New("worker queue full")
)

// Handler (Function): processes a single event. It should return an error
// if processing fails, which will be passed to the ErrorHandler.
type Handler func(context.Context, events.Event) error

// ErrorHandler (Function): is called when a handler returns an error or panics.
type ErrorHandler func(events.Event, error)

// DurabilityMode (Type): defines the persistence strategy for the worker queue.
type DurabilityMode string

const (
	// DurabilityInMemory stores events in an in-memory channel (default).
	DurabilityInMemory DurabilityMode = "in_memory"
	// DurabilityDurable indicates events should be persisted to survive restarts.
	DurabilityDurable DurabilityMode = "durable"
)

// PoolConfig (Struct): configures a new worker pool instance.
type PoolConfig struct {
	Name           string
	Workers        int
	QueueSize      int
	Handler        Handler
	OnError        ErrorHandler
	DurabilityMode DurabilityMode
}

// Pool (Struct): manages a fixed-size set of concurrent workers that process
// events from a shared channel. It provides bounded queuing, graceful shutdown
// with drain support, and non-blocking submission options.
//
// Boundary: Workers coordinate through channels; stopCh signals shutdown.
// Failure mode: Panics in handlers are recovered and passed to OnError.
type Pool struct {
	name    string
	workers int
	jobs    chan events.Event
	handler Handler
	onError ErrorHandler
	mode    DurabilityMode

	mu      sync.RWMutex
	started bool
	stopped bool
	wg      sync.WaitGroup
	stopCh  chan struct{}
}

// NewPool (Function): creates a new worker pool with the given configuration.
// It validates that Workers > 0, QueueSize > 0, and Handler is not nil.
// Returns ErrInvalidPoolConfig if validation fails.
func NewPool(config PoolConfig) (*Pool, error) {
	if config.Workers <= 0 || config.QueueSize <= 0 || config.Handler == nil {
		return nil, ErrInvalidPoolConfig
	}
	if config.Name == "" {
		config.Name = "worker-pool"
	}
	if config.DurabilityMode == "" {
		config.DurabilityMode = DurabilityInMemory
	}
	if config.DurabilityMode != DurabilityInMemory && config.DurabilityMode != DurabilityDurable {
		return nil, ErrInvalidPoolConfig
	}

	return &Pool{
		name:    config.Name,
		workers: config.Workers,
		jobs:    make(chan events.Event, config.QueueSize),
		handler: config.Handler,
		onError: config.OnError,
		mode:    config.DurabilityMode,
		stopCh:  make(chan struct{}),
	}, nil
}

// Start (Method): launches worker goroutines and begins processing events
func (p *Pool) Start(ctx context.Context) error {
	if p == nil {
		return fmt.Errorf("worker pool is not configured")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return ErrPoolStopped
	}
	if p.started {
		p.mu.Unlock()
		return nil
	}
	p.started = true

	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
	}
	p.mu.Unlock()

	for i := 0; i < p.workers; i++ {
		go p.runWorker(ctx)
	}

	return nil
}

// Submit enqueues an event for processing. It blocks until the event is
// accepted, the pool is stopped, or the context is cancelled.
//
// Note: p.jobs is never closed while the pool is running, so there is no
// risk of a "send on closed channel" panic here. Workers exit by watching
// stopCh, not by channel closure.
// Submit (Method): enqueues an event for processing, blocking until accepted, pool stopped, or context cancelled
func (p *Pool) Submit(ctx context.Context, event events.Event) error {
	if p == nil {
		return fmt.Errorf("worker pool is not configured")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	// Fast-path: reject if pool is already stopped. Without this check, a
	// stopped pool with remaining buffer capacity would non-deterministically
	// accept the send in the select below (Go selects randomly among ready
	// cases). The pre-check makes rejection deterministic.
	select {
	case <-p.stopCh:
		return ErrPoolStopped
	default:
	}

	select {
	case p.jobs <- event:
		return nil
	case <-p.stopCh:
		return ErrPoolStopped
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TrySubmit attempts a non-blocking enqueue. Returns ErrQueueFull immediately
// if the job channel has no capacity, or ErrPoolStopped if the pool has been
// shut down.
//
// The stopped check and the channel send are both guarded by the read-lock so
// that Stop (which holds the write-lock when it sets p.stopped) cannot sneak
// in between the two operations. Because p.jobs is never closed while the pool
// is running, the non-blocking send is safe.
// TrySubmit (Method): non-blocking enqueue; returns ErrQueueFull or ErrPoolStopped immediately
func (p *Pool) TrySubmit(event events.Event) error {
	if p == nil {
		return fmt.Errorf("worker pool is not configured")
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.stopped {
		return ErrPoolStopped
	}

	select {
	case p.jobs <- event:
		return nil
	default:
		return ErrQueueFull
	}
}

// Stop signals all workers to exit and waits for them to finish.
//
// Crucially, p.jobs is NOT closed here. Closing a channel that other
// goroutines may still be sending to causes a panic. Instead workers
// watch stopCh and drain any remaining buffered events before returning.
// Stop (Method): signals all workers to exit and waits for them to finish draining buffered events
func (p *Pool) Stop() {
	if p == nil {
		return
	}

	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return
	}
	p.stopped = true
	close(p.stopCh)
	p.mu.Unlock()

	p.wg.Wait()
}

// QueueLength (Method): returns the number of events currently buffered in the job channel
func (p *Pool) QueueLength() int {
	if p == nil {
		return 0
	}

	return len(p.jobs)
}

// Name (Method): returns the pool's configured name
func (p *Pool) Name() string {
	if p == nil {
		return ""
	}

	return p.name
}

// DurabilityMode (Method): returns the pool's configured durability mode
func (p *Pool) DurabilityMode() DurabilityMode {
	if p == nil {
		return DurabilityInMemory
	}
	return p.mode
}

// runWorker is the per-goroutine event loop. It exits when either:
//   - stopCh is closed (graceful shutdown): remaining buffered events are drained
//     so that no accepted work is silently dropped.
//   - ctx is cancelled: the worker stops immediately without draining, because
//     the caller's context signals an abort, not a graceful stop.
func (p *Pool) runWorker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			// Caller context cancelled - stop immediately, do not drain.
			return

		case <-p.stopCh:
			// Graceful shutdown: drain any already-buffered events before exiting
			// so that work that was accepted prior to Stop() was not lost.
			// A fresh background context is used because the original ctx may
			// already be cancelled, which would make every handler call fail.
			drainCtx := context.Background()
			for {
				select {
				case event, ok := <-p.jobs:
					if !ok {
						return
					}
					p.doHandle(drainCtx, event)
				default:
					return
				}
			}

		case event, ok := <-p.jobs:
			if !ok {
				return
			}
			p.doHandle(ctx, event)
		}
	}
}

// doHandle (Method): executes the handler with panic recovery
func (p *Pool) doHandle(ctx context.Context, event events.Event) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("panic: %v", r)
			}
			if p.onError != nil {
				p.onError(event, err)
			}
		}
	}()

	if err := p.handler(ctx, event); err != nil && p.onError != nil {
		p.onError(event, err)
	}
}
