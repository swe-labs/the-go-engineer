// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package payment

import (
	"context"
	"fmt"
)

// Handler (Function): processes one payment job; returns error if processing fails
type Handler func(context.Context, Job) error

// Worker (Struct): consumes tenant-scoped payment jobs until the queue closes or context is cancelled
type Worker struct {
	handler Handler
	jobs    <-chan Job
}

// NewWorker (Constructor): creates a consumer that reads and processes payment jobs from a channel
func NewWorker(jobs <-chan Job, handler Handler) Worker {
	return Worker{
		jobs:    jobs,
		handler: handler,
	}
}

// Run (Method): blocks until context cancelled or jobs channel closed, executing handler for each job
func (w Worker) Run(ctx context.Context) error {
	if w.handler == nil {
		return fmt.Errorf("payment worker handler is not configured")
	}
	if w.jobs == nil {
		return fmt.Errorf("payment worker queue is not configured")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case job, ok := <-w.jobs:
			if !ok {
				return nil
			}
			if err := w.handler(ctx, job); err != nil {
				return fmt.Errorf("process payment job: %w", err)
			}
		}
	}
}
