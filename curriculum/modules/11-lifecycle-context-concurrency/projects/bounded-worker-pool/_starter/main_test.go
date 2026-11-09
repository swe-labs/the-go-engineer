package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool_ConcurrencyBound(t *testing.T) {
	const workerCount = 4
	const totalJobs = 50

	pool := NewWorkerPool(PoolConfig{Workers: workerCount, QueueSize: totalJobs})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	var running atomic.Int32
	var maxObserved atomic.Int32
	var completed atomic.Int32

	for i := 0; i < totalJobs; i++ {
		err := pool.Submit(func(ctx context.Context) error {
			current := running.Add(1)
			defer running.Add(-1)

			for {
				max := maxObserved.Load()
				if current <= max || maxObserved.CompareAndSwap(max, current) {
					break
				}
			}

			time.Sleep(5 * time.Millisecond)
			completed.Add(1)
			return nil
		})
		if err != nil {
			t.Fatalf("failed to submit job: %v", err)
		}
	}

	pool.Stop()

	if completed.Load() != totalJobs {
		t.Fatalf("expected %d jobs completed, got %d", totalJobs, completed.Load())
	}

	if maxObserved.Load() > workerCount {
		t.Fatalf("max concurrent workers (%d) exceeded configured limit (%d)", maxObserved.Load(), workerCount)
	}

	if pool.PeakWorkers() > workerCount {
		t.Fatalf("pool.PeakWorkers() (%d) exceeded limit (%d)", pool.PeakWorkers(), workerCount)
	}
}

func TestWorkerPool_PanicContainment(t *testing.T) {
	pool := NewWorkerPool(PoolConfig{Workers: 2, QueueSize: 10})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	var completed atomic.Int32

	// Job 1: Panics!
	_ = pool.Submit(func(ctx context.Context) error {
		panic("fatal panic in job")
	})

	// Job 2 & 3: Normal jobs executed after panic
	for i := 0; i < 2; i++ {
		_ = pool.Submit(func(ctx context.Context) error {
			completed.Add(1)
			return nil
		})
	}

	pool.Stop()

	stats := pool.Stats()
	if stats.Panicked != 1 {
		t.Fatalf("expected 1 panicked job, got %d", stats.Panicked)
	}
	if stats.Completed != 2 {
		t.Fatalf("expected 2 completed jobs despite prior panic, got %d", stats.Completed)
	}
}

func TestWorkerPool_BackpressureAndQueueFull(t *testing.T) {
	// 1 worker, 1 queue slot
	pool := NewWorkerPool(PoolConfig{Workers: 1, QueueSize: 1})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	started := make(chan struct{})
	release := make(chan struct{})

	// Occupy worker
	if err := pool.Submit(func(ctx context.Context) error {
		close(started)
		<-release
		return nil
	}); err != nil {
		t.Fatalf("failed to submit worker job: %v", err)
	}

	<-started

	// Occupy queue buffer
	err := pool.Submit(func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected second job in queue, got error: %v", err)
	}

	// Third job must return ErrQueueFull
	errThird := pool.Submit(func(ctx context.Context) error {
		return nil
	})
	if !errors.Is(errThird, ErrQueueFull) {
		t.Fatalf("expected ErrQueueFull, got %v", errThird)
	}

	// Submit with timeout context must timeout
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer timeoutCancel()
	errTimeout := pool.SubmitWithContext(timeoutCtx, func(ctx context.Context) error {
		return nil
	})
	if !errors.Is(errTimeout, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", errTimeout)
	}

	// Release worker
	close(release)
	pool.Stop()
}

func TestWorkerPool_GracefulShutdown(t *testing.T) {
	pool := NewWorkerPool(PoolConfig{Workers: 3, QueueSize: 20})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	var processed atomic.Int32
	for i := 0; i < 15; i++ {
		_ = pool.Submit(func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			processed.Add(1)
			return nil
		})
	}

	// Stop must wait for all 15 submitted jobs to finish
	pool.Stop()

	if processed.Load() != 15 {
		t.Fatalf("expected all 15 queued jobs to be processed before shutdown, got %d", processed.Load())
	}

	// Further submits must return ErrPoolClosed
	err := pool.Submit(func(ctx context.Context) error { return nil })
	if !errors.Is(err, ErrPoolClosed) {
		t.Fatalf("expected ErrPoolClosed after stop, got %v", err)
	}
}

func TestWorkerPool_RaceConditionFreedom(t *testing.T) {
	pool := NewWorkerPool(PoolConfig{Workers: 8, QueueSize: 100})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)

	var wg sync.WaitGroup
	for producer := 0; producer < 5; producer++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				_ = pool.Submit(func(ctx context.Context) error {
					return nil
				})
			}
		}()
	}

	wg.Wait()
	pool.Stop()
}
