// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./12-concurrency-patterns/6-worker-pool
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ============================================================================
// Section 12: Concurrency Patterns — Robust Worker Pool
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to build a bounded worker pool handling a queue of jobs.
//   - How to gracefully shut down the pool (drain the queue or cancel midway).
//   - How to recover from worker panics so one bad job doesn't kill the app.
//
// ENGINEERING DEPTH:
//   Unbounded concurrency (spawning a Goroutine per job) is dangerous if the
//   queue is large (e.g., parsing 100,000 URLs). You can easily exhaust file
//   descriptors, RAM, or database connections.
//   A Worker Pool limits the maximum concurrency explicitly.
// ============================================================================

// Job represents a unit of work.
type Job struct {
	ID  int
	URL string
}

// Result represents the outcome of a Job.
type Result struct {
	JobID int
	Data  string
	Err   error
}

// 1. The Worker Function
// It reads from `jobs` channel and sends to `results` channel.
// Notice it takes a context to allow cancellation.
func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		// If the context is canceled, the worker stops processing new jobs.
		case <-ctx.Done():
			log.Printf("Worker %d shutting down (context canceled)", id)
			return

		// If a job is available, process it.
		// If the `jobs` channel is closed, `ok` will be false, and the worker exits.
		case j, ok := <-jobs:
			if !ok {
				log.Printf("Worker %d shutting down (queue empty & closed)", id)
				return
			}
			processJobSafe(id, j, results)
		}
	}
}

// 2. Safe Job Processing
// We isolate the actual work into a function with a `defer recover()`
// to ensure a panicked job doesn't kill the entire worker (or application).
func processJobSafe(workerID int, j Job, results chan<- Result) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in worker %d processing job %d: %v", workerID, j.ID, r)
			results <- Result{JobID: j.ID, Err: err}
		}
	}()

	log.Printf("Worker %d starting job %d (%s)", workerID, j.ID, j.URL)

	// Simulate work (e.g., HTTP GET or heavy parsing)
	time.Sleep(500 * time.Millisecond)

	// Simulate an unexpected panic on a specific nasty job
	if j.ID == 3 {
		panic("simulated nasty bug in parser")
	}

	results <- Result{
		JobID: j.ID,
		Data:  fmt.Sprintf("Scraped data from %s", j.URL),
		Err:   nil,
	}
}

func main() {
	// Setup Context mapping to OS Interrupts (Ctrl+C).
	// This allows us to gracefully shut down the worker pool.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	const numWorkers = 3
	const numJobs = 10

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup

	// Step 1: Start the fixed pool of workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(ctx, w, jobs, results, &wg)
	}

	// Step 2: Queue up all the jobs.
	// In a real app, you might be reading URLs from a DB or RabbitMQ here.
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, URL: fmt.Sprintf("https://example.com/page/%d", j)}
	}

	// Step 3: Tell the workers no more jobs are coming (by closing the channel).
	// They will finish whatever is left in the buffer and then exit.
	close(jobs)

	// Step 4: Wait for all workers to finish IN A SEPARATE GOROUTINE.
	// This way, we can safely close the `results` channel without deadlocking the main thread
	// which is busy reading from `results` below.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Step 5: Consume the results
	success := 0
	failures := 0
	for r := range results {
		if r.Err != nil {
			log.Printf("❌ Failed job %d: %v", r.JobID, r.Err)
			failures++
		} else {
			log.Printf("✅ Success job %d: %s", r.JobID, r.Data)
			success++
		}
	}

	fmt.Printf("\nDone! Success: %d, Failures: %d\n", success, failures)
}
