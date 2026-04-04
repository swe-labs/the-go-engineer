// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

// ============================================================================
// Section 24: errgroup & sync.Pool — errgroup with Context Cancellation
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - errgroup.WithContext: automatic cancellation on first error
//   - How to implement a fan-out/fan-in pipeline with errgroup
//   - Producer + multiple consumers pattern (parallel processing pipeline)
//
// ENGINEERING DEPTH:
//   errgroup.WithContext creates a context that is cancelled automatically
//   the moment any goroutine returns a non-nil error. This is the answer to
//   the "I launched 10 goroutines but one failed — how do I stop the other 9?"
//   problem.
//
//   The pattern:
//     g, ctx := errgroup.WithContext(parent)
//     g.Go(func() error { return producer(ctx) })
//     g.Go(func() error { return consumer(ctx) })
//     return g.Wait() // cancel() called automatically on first error
//
//   This replaces the common anti-pattern of manually done channels:
//     done := make(chan struct{})
//     var once sync.Once
//     cancel := func() { once.Do(func() { close(done) }) }
//     // ... 20 lines of error propagation ...
//
// RUN: go run ./24-errgroup-and-pools/2-errgroup-context
// ============================================================================

// WorkItem represents a URL to be crawled.
type WorkItem struct {
	URL      string
	Priority int
}

// Result represents the outcome of processing a WorkItem.
type Result struct {
	URL      string
	StatusOK bool
	Latency  time.Duration
}

// producer sends work items into the jobs channel.
// It respects context cancellation — if the context is cancelled (e.g., a
// consumer failed), the producer stops immediately instead of loading more work.
func producer(ctx context.Context, jobs chan<- WorkItem) error {
	defer close(jobs) // Signal: no more work coming

	urls := []string{
		"https://api.example.com/users",
		"https://api.example.com/orders",
		"https://api.example.com/products",
		"https://api.example.com/inventory",
		"https://api.example.com/payments",
		"https://api.example.com/analytics",
	}

	for i, url := range urls {
		select {
		case <-ctx.Done():
			// Context cancelled (a consumer errored out). Stop producing.
			slog.Warn("producer cancelled", "reason", ctx.Err(), "sent", i)
			return ctx.Err()
		case jobs <- WorkItem{URL: url, Priority: i}:
			slog.Info("work item queued", "url", url)
			time.Sleep(20 * time.Millisecond) // Simulate item production pace
		}
	}

	return nil
}

// consumer reads from jobs, processes each item, and sends results.
// It uses ctx to detect early cancellation from other goroutines.
func consumer(ctx context.Context, id int, jobs <-chan WorkItem, results chan<- Result) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case item, ok := <-jobs:
			if !ok {
				return nil // Channel closed — no more work
			}

			start := time.Now()
			slog.Info("worker processing", "worker_id", id, "url", item.URL)

			// Simulate random failure in worker 2 for demonstration
			if id == 2 && rand.Intn(3) == 0 {
				return fmt.Errorf("worker %d: connection reset on %s", id, item.URL)
			}

			// Simulate HTTP request latency
			time.Sleep(time.Duration(30+rand.Intn(50)) * time.Millisecond)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case results <- Result{URL: item.URL, StatusOK: true, Latency: time.Since(start)}:
			}
		}
	}
}

// resultCollector drains the results channel until it is closed.
func resultCollector(ctx context.Context, results <-chan Result) error {
	var count int
	for {
		select {
		case <-ctx.Done():
			return nil // Allow clean exit
		case r, ok := <-results:
			if !ok {
				slog.Info("collection complete", "count", count)
				return nil
			}
			count++
			status := "✅"
			if !r.StatusOK {
				status = "❌"
			}
			slog.Info("result", "status", status, "url", r.URL, "latency", r.Latency.Round(time.Millisecond))
		}
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(nil, &slog.HandlerOptions{Level: slog.LevelInfo})))

	fmt.Println("=== Fan-out pipeline with errgroup.WithContext ===")
	start := time.Now()

	// errgroup.WithContext returns:
	//   g — the group. g.Wait() returns the first error.
	//   ctx — cancelled automatically when any goroutine returns an error.
	g, ctx := errgroup.WithContext(context.Background())

	// Buffered channels decouple producer pace from consumer pace.
	jobs := make(chan WorkItem, 10)
	results := make(chan Result, 10)

	// Launch producer
	g.Go(func() error {
		return producer(ctx, jobs)
	})

	// Launch 3 consumers, each reading from the same jobs channel.
	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		i := i
		g.Go(func() error {
			return consumer(ctx, i, jobs, results)
		})
	}

	// Close results channel AFTER all producers and consumers finish.
	// This requires a separate goroutine — g.Wait() blocks, so we can't call it
	// and then close results in the same goroutine.
	go func() {
		g.Wait()
		close(results)
	}()

	// Collect results in the main goroutine (or a dedicated collector goroutine).
	var totalResults int
	for r := range results {
		totalResults++
		_ = r
	}

	// g.Wait() has already been called above. Call it again to get the error.
	// Calling Wait() multiple times is safe — it always returns the same error.
	if err := g.Wait(); err != nil && err != context.Canceled {
		fmt.Printf("❌ Pipeline failed: %v\n", err)
	} else {
		fmt.Printf("✅ Pipeline complete: %d results in %v\n", totalResults, time.Since(start).Round(time.Millisecond))
	}

	// KEY TAKEAWAY:
	// - errgroup.WithContext: automatic context cancellation on first error
	// - Producer + consumers + result collector = idiomatic fan-out pipeline
	// - Separate goroutine closes the results channel after g.Wait()
	// - ctx.Done() in select cases makes goroutines respond to cancellation
	// - This replaces hundreds of lines of manual done-channel machinery
}
