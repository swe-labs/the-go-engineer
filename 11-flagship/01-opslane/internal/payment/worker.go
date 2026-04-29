// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package payment

import (
	"context"
	"fmt"
)

// Handler processes one payment job.
type Handler func(context.Context, Job) error

// Worker consumes tenant-scoped payment jobs until the queue closes or context stops.
type Worker struct {
	handler Handler
	jobs    <-chan Job
}

// NewWorker instantiates a consumer that reads payment jobs from a channel
// and processes them sequentially using the provided handler.
func NewWorker(jobs <-chan Job, handler Handler) Worker {
	return Worker{
		jobs:    jobs,
		handler: handler,
	}
}

// Run blocks until the context is cancelled or the jobs channel is closed.
// It executes the handler for each job received.
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
