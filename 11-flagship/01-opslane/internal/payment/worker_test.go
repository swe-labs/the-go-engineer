package payment

import (
	"context"
	"errors"
	"testing"
)

func TestWorkerProcessesJobsUntilQueueCloses(t *testing.T) {
	t.Parallel()

	jobs := make(chan Job, 2)
	jobs <- Job{TenantID: 7, OrderID: 101, ProviderReference: "pay_1", AmountCents: 2500}
	jobs <- Job{TenantID: 7, OrderID: 102, ProviderReference: "pay_2", AmountCents: 4000}
	close(jobs)

	processed := 0
	worker := NewWorker(jobs, func(context.Context, Job) error {
		processed++
		return nil
	})

	if err := worker.Run(context.Background()); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if processed != 2 {
		t.Fatalf("processed = %d, want 2", processed)
	}
}

func TestWorkerReturnsHandlerError(t *testing.T) {
	t.Parallel()

	handlerErr := errors.New("gateway down")
	jobs := make(chan Job, 1)
	jobs <- Job{TenantID: 7, OrderID: 101, ProviderReference: "pay_1", AmountCents: 2500}
	close(jobs)

	worker := NewWorker(jobs, func(context.Context, Job) error {
		return handlerErr
	})

	err := worker.Run(context.Background())
	if !errors.Is(err, handlerErr) {
		t.Fatalf("Run error = %v, want handlerErr", err)
	}
}
