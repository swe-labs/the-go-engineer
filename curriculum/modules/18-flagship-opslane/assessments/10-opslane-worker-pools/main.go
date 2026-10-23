package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var ErrPoolClosed = errors.New("worker pool is closed to new jobs")

// Job defines a concurrent unit of executable work.
type Job interface {
	ID() string
	Execute(ctx context.Context) error
}

// SimpleJob provides a basic Job implementation.
type SimpleJob struct {
	JobID string
	Work  func(ctx context.Context) error
}

func (j *SimpleJob) ID() string                        { return j.JobID }
func (j *SimpleJob) Execute(ctx context.Context) error { return j.Work(ctx) }

// BoundedWorkerPool manages a fixed set of concurrent worker goroutines.
type BoundedWorkerPool struct {
	workerCount        int
	queueCapacity      int
	jobs               chan Job
	wg                 sync.WaitGroup
	ctx                context.Context
	cancel             context.CancelFunc
	activeWorkers      atomic.Int32
	maxConcurrencySeen atomic.Int32
	completedJobs      atomic.Int64
	closed             atomic.Bool
}

func NewBoundedWorkerPool(workers, queueCap int) *BoundedWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &BoundedWorkerPool{
		workerCount:   workers,
		queueCapacity: queueCap,
		jobs:          make(chan Job, queueCap),
		ctx:           ctx,
		cancel:        cancel,
	}

	for i := 1; i <= workers; i++ {
		pool.wg.Add(1)
		go pool.workerLoop(i)
	}

	return pool
}

func (p *BoundedWorkerPool) workerLoop(workerID int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			// Drain remaining buffered jobs before quitting
			for {
				select {
				case job, ok := <-p.jobs:
					if !ok {
						return
					}
					p.executeJob(job)
				default:
					return
				}
			}

		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			p.executeJob(job)
		}
	}
}

func (p *BoundedWorkerPool) executeJob(j Job) {
	active := p.activeWorkers.Add(1)

	// Track peak concurrent execution
	for {
		currentMax := p.maxConcurrencySeen.Load()
		if active <= currentMax {
			break
		}
		if p.maxConcurrencySeen.CompareAndSwap(currentMax, active) {
			break
		}
	}

	defer func() {
		p.activeWorkers.Add(-1)
		p.completedJobs.Add(1)
		if r := recover(); r != nil {
			log.Printf("[WorkerPool] Recovered worker panic on job %s: %v", j.ID(), r)
		}
	}()

	_ = j.Execute(p.ctx)
}

// Submit enqueues a job into the bounded pool.
func (p *BoundedWorkerPool) Submit(j Job) error {
	if p.closed.Load() {
		return ErrPoolClosed
	}

	select {
	case <-p.ctx.Done():
		return ErrPoolClosed
	case p.jobs <- j:
		return nil
	}
}

// Shutdown initiates graceful termination of the worker pool.
func (p *BoundedWorkerPool) Shutdown(ctx context.Context) error {
	if p.closed.Swap(true) {
		return nil // Already shutting down
	}

	close(p.jobs)

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		p.cancel() // Force terminate remaining
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (p *BoundedWorkerPool) Stats() (active, maxSeen int32, completed int64) {
	return p.activeWorkers.Load(), p.maxConcurrencySeen.Load(), p.completedJobs.Load()
}

func main() {
	pool := NewBoundedWorkerPool(4, 50)

	for i := 1; i <= 20; i++ {
		id := fmt.Sprintf("job-%d", i)
		_ = pool.Submit(&SimpleJob{
			JobID: id,
			Work: func(ctx context.Context) error {
				time.Sleep(10 * time.Millisecond)
				return nil
			},
		})
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = pool.Shutdown(shutdownCtx)

	_, maxSeen, completed := pool.Stats()
	fmt.Printf("Worker pool completed %d jobs (Max concurrency: %d)\n", completed, maxSeen)
}
