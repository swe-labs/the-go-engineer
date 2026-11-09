package main

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrPoolClosed      = errors.New("worker pool is closed")
	ErrQueueFull       = errors.New("job queue is full")
	ErrNilJob          = errors.New("job cannot be nil")
	ErrShutdownTimeout = errors.New("shutdown timed out")
)

// Job represents an executable unit of work.
type Job func(ctx context.Context) error

// PoolConfig specifies worker pool capacity limits.
type PoolConfig struct {
	Workers   int
	QueueSize int
}

// PoolStats holds diagnostic metrics.
type PoolStats struct {
	ActiveWorkers int32
	Completed     int64
	Failed        int64
	Panicked      int64
}

// WorkerPool manages a bounded set of worker goroutines.
type WorkerPool struct {
	workers  int
	jobs     chan Job
	wg       sync.WaitGroup
	mu       sync.RWMutex
	closed   bool
	stopOnce sync.Once
	doneChan chan struct{}

	activeWorkers int32
	peakWorkers   int32
	completed     atomic.Int64
	failed        atomic.Int64
	panicked      atomic.Int64
}

// NewWorkerPool initializes a bounded pool with workers and queue capacity.
func NewWorkerPool(cfg PoolConfig) *WorkerPool {
	if cfg.Workers <= 0 {
		cfg.Workers = 1
	}
	if cfg.QueueSize < 0 {
		cfg.QueueSize = 0
	}

	return &WorkerPool{
		workers:  cfg.Workers,
		jobs:     make(chan Job, cfg.QueueSize),
		doneChan: make(chan struct{}),
	}
}

// Start launches worker goroutines.
func (p *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.workerLoop(ctx, i)
	}
}

// workerLoop is the execution loop for each worker goroutine.
func (p *WorkerPool) workerLoop(ctx context.Context, workerID int) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			// Context canceled: drain any remaining jobs non-blockingly or return
			return
		case job, ok := <-p.jobs:
			if !ok {
				// Job queue closed: exit worker cleanly
				return
			}
			p.executeJob(ctx, job)
		}
	}
}

// executeJob wraps job invocation with panic containment and metric tracking.
func (p *WorkerPool) executeJob(ctx context.Context, job Job) {
	current := atomic.AddInt32(&p.activeWorkers, 1)
	defer atomic.AddInt32(&p.activeWorkers, -1)

	// Track peak concurrent workers
	for {
		peak := atomic.LoadInt32(&p.peakWorkers)
		if current <= peak || atomic.CompareAndSwapInt32(&p.peakWorkers, peak, current) {
			break
		}
	}

	defer func() {
		if rec := recover(); rec != nil {
			p.panicked.Add(1)
			slog.Error("worker job recovered from panic", "panic", rec)
		}
	}()

	if err := job(ctx); err != nil {
		p.failed.Add(1)
	} else {
		p.completed.Add(1)
	}
}

// Submit queues a job. Returns ErrQueueFull immediately if the queue buffer is full.
func (p *WorkerPool) Submit(job Job) error {
	if job == nil {
		return ErrNilJob
	}

	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return ErrPoolClosed
	}
	p.mu.RUnlock()

	select {
	case p.jobs <- job:
		return nil
	default:
		return ErrQueueFull
	}
}

// SubmitWithContext queues a job waiting until context expires or slot becomes available.
func (p *WorkerPool) SubmitWithContext(ctx context.Context, job Job) error {
	if job == nil {
		return ErrNilJob
	}

	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return ErrPoolClosed
	}
	p.mu.RUnlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.jobs <- job:
		return nil
	}
}

// Stop stops accepting new jobs and gracefully waits for queued jobs to complete.
func (p *WorkerPool) Stop() {
	p.stopOnce.Do(func() {
		p.mu.Lock()
		p.closed = true
		close(p.jobs)
		p.mu.Unlock()

		p.wg.Wait()
		close(p.doneChan)
	})
}

// StopWithTimeout stops the pool and waits up to timeout duration for workers to finish.
func (p *WorkerPool) StopWithTimeout(timeout time.Duration) error {
	p.mu.Lock()
	if !p.closed {
		p.closed = true
		close(p.jobs)
	}
	p.mu.Unlock()

	waitCh := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		return nil
	case <-time.After(timeout):
		return ErrShutdownTimeout
	}
}

// Stats returns current operational metrics.
func (p *WorkerPool) Stats() PoolStats {
	return PoolStats{
		ActiveWorkers: atomic.LoadInt32(&p.activeWorkers),
		Completed:     p.completed.Load(),
		Failed:        p.failed.Load(),
		Panicked:      p.panicked.Load(),
	}
}

// PeakWorkers returns highest observed simultaneous workers.
func (p *WorkerPool) PeakWorkers() int32 {
	return atomic.LoadInt32(&p.peakWorkers)
}

func main() {
	cfg := PoolConfig{Workers: 4, QueueSize: 20}
	pool := NewWorkerPool(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool.Start(ctx)
	slog.Info("worker pool started", "workers", cfg.Workers, "queue_size", cfg.QueueSize)

	// Submit sample workload
	for i := 1; i <= 15; i++ {
		jobID := i
		err := pool.Submit(func(ctx context.Context) error {
			time.Sleep(20 * time.Millisecond)
			slog.Info("job finished", "job_id", jobID)
			return nil
		})
		if err != nil {
			slog.Warn("submission failed", "job_id", jobID, "error", err)
		}
	}

	// Graceful shutdown
	pool.Stop()
	stats := pool.Stats()
	slog.Info("pool shutdown complete",
		"completed", stats.Completed,
		"failed", stats.Failed,
		"panicked", stats.Panicked,
		"peak_workers", pool.PeakWorkers(),
	)
}
