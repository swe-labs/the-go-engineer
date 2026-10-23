package main

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool_BoundedConcurrencyAndCompletion(t *testing.T) {
	// Core practice task requirement:
	// Submit 100 jobs to a worker pool with 5 workers,
	// verify that at most 5 jobs run concurrently,
	// and confirm that all 100 jobs complete.
	const totalJobs = 100
	const maxWorkers = 5

	pool := NewBoundedWorkerPool(maxWorkers, totalJobs)

	var runningCount atomic.Int32
	var peakConcurrency atomic.Int32
	var completedCount atomic.Int32

	for i := 1; i <= totalJobs; i++ {
		jobID := fmt.Sprintf("bench-job-%d", i)
		err := pool.Submit(&SimpleJob{
			JobID: jobID,
			Work: func(ctx context.Context) error {
				curr := runningCount.Add(1)

				// Update observed peak concurrency
				for {
					peak := peakConcurrency.Load()
					if curr <= peak {
						break
					}
					if peakConcurrency.CompareAndSwap(peak, curr) {
						break
					}
				}

				// Simulate realistic I/O work
				time.Sleep(5 * time.Millisecond)

				runningCount.Add(-1)
				completedCount.Add(1)
				return nil
			},
		})
		if err != nil {
			t.Fatalf("failed to submit job %d: %v", i, err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Shutdown(ctx); err != nil {
		t.Fatalf("worker pool shutdown error: %v", err)
	}

	// 1. Verify all 100 jobs completed
	if completedCount.Load() != totalJobs {
		t.Fatalf("expected all %d jobs to complete, but only %d completed", totalJobs, completedCount.Load())
	}

	// 2. Verify concurrency was strictly bounded by maxWorkers
	observedPeak := peakConcurrency.Load()
	if observedPeak > maxWorkers {
		t.Fatalf("CONCURRENCY INVARIANT VIOLATION: expected peak concurrency <= %d, but observed %d!", maxWorkers, observedPeak)
	}
	if observedPeak == 0 {
		t.Errorf("observed peak concurrency was 0")
	}
}

func TestWorkerPool_JobPanicRecovery(t *testing.T) {
	pool := NewBoundedWorkerPool(2, 10)

	var normalJobCompleted atomic.Bool

	// Job 1 panics
	_ = pool.Submit(&SimpleJob{
		JobID: "panic-job",
		Work: func(ctx context.Context) error {
			panic("fatal unhandled error in worker job")
		},
	})

	// Job 2 succeeds normally
	_ = pool.Submit(&SimpleJob{
		JobID: "normal-job",
		Work: func(ctx context.Context) error {
			normalJobCompleted.Store(true)
			return nil
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = pool.Shutdown(ctx)

	if !normalJobCompleted.Load() {
		t.Errorf("expected normal job to complete despite prior panic in worker pool")
	}

	active, _, completed := pool.Stats()
	if active != 0 {
		t.Errorf("expected 0 active workers after shutdown, got %d", active)
	}
	if completed != 2 {
		t.Errorf("expected 2 completed/recovered jobs in stats, got %d", completed)
	}
}

func TestWorkerPool_ShutdownRejectsNewJobs(t *testing.T) {
	pool := NewBoundedWorkerPool(2, 5)

	_ = pool.Shutdown(context.Background())

	err := pool.Submit(&SimpleJob{JobID: "late-job", Work: func(ctx context.Context) error { return nil }})
	if !errors.Is(err, ErrPoolClosed) {
		t.Errorf("expected ErrPoolClosed when submitting to shut-down pool, got: %v", err)
	}
}
