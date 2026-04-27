// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package workers

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/events"
)

var (
	ErrInvalidPoolConfig = errors.New("invalid worker pool config")
	ErrPoolStopped       = errors.New("worker pool stopped")
	ErrQueueFull         = errors.New("worker queue full")
)

type Handler func(context.Context, events.Event) error
type ErrorHandler func(events.Event, error)

type PoolConfig struct {
	Name      string
	Workers   int
	QueueSize int
	Handler   Handler
	OnError   ErrorHandler
}

type Pool struct {
	name    string
	workers int
	jobs    chan events.Event
	handler Handler
	onError ErrorHandler

	mu      sync.RWMutex
	started bool
	stopped bool
	wg      sync.WaitGroup
}

func NewPool(config PoolConfig) (*Pool, error) {
	if config.Workers <= 0 || config.QueueSize <= 0 || config.Handler == nil {
		return nil, ErrInvalidPoolConfig
	}
	if config.Name == "" {
		config.Name = "worker-pool"
	}

	return &Pool{
		name:    config.Name,
		workers: config.Workers,
		jobs:    make(chan events.Event, config.QueueSize),
		handler: config.Handler,
		onError: config.OnError,
	}, nil
}

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
	p.mu.Unlock()

	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.runWorker(ctx)
	}

	return nil
}

func (p *Pool) Submit(ctx context.Context, event events.Event) error {
	if p == nil {
		return fmt.Errorf("worker pool is not configured")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	p.mu.RLock()
	if p.stopped {
		p.mu.RUnlock()
		return ErrPoolStopped
	}

	select {
	case p.jobs <- event:
		p.mu.RUnlock()
		return nil
	default:
	}
	p.mu.RUnlock()

	select {
	case p.jobs <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

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
	close(p.jobs)
	p.mu.Unlock()

	p.wg.Wait()
}

func (p *Pool) QueueLength() int {
	if p == nil {
		return 0
	}

	return len(p.jobs)
}

func (p *Pool) Name() string {
	if p == nil {
		return ""
	}

	return p.name
}

func (p *Pool) runWorker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
		case event, ok := <-p.jobs:
			if !ok {
				return
			}
			if err := p.handler(ctx, event); err != nil && p.onError != nil {
				p.onError(event, err)
			}
			continue
		}

		for {
			select {
			case event, ok := <-p.jobs:
				if !ok {
					return
				}
				if err := p.handler(ctx, event); err != nil && p.onError != nil {
					p.onError(event, err)
				}
			default:
				return
			}
		}
	}
}
