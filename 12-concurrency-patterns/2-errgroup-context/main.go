// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

// ============================================================================
// Section 12: Concurrency Patterns � errgroup with Context Cancellation
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
//   the "I launched 10 goroutines but one failed � how do I stop the other 9?"
//   problem.
//
// RUN: go run ./12-concurrency-patterns/2-errgroup-context
// ============================================================================

type WorkItem struct {
	URL      string
	Priority int
}

type Result struct {
	URL      string
	StatusOK bool
	Latency  time.Duration
}

func producer(ctx context.Context, jobs chan<- WorkItem) error {
	defer close(jobs)

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
			slog.Warn("producer cancelled", "reason", ctx.Err(), "sent", i)
			return ctx.Err()
		case jobs <- WorkItem{URL: url, Priority: i}:
			slog.Info("work item queued", "url", url)
			time.Sleep(20 * time.Millisecond)
		}
	}

	return nil
}

func consumer(ctx context.Context, id int, jobs <-chan WorkItem, results chan<- Result) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case item, ok := <-jobs:
			if !ok {
				return nil
			}

			start := time.Now()
			slog.Info("worker processing", "worker_id", id, "url", item.URL)

			if id == 2 && rand.Intn(3) == 0 {
				return fmt.Errorf("worker %d: connection reset on %s", id, item.URL)
			}

			time.Sleep(time.Duration(30+rand.Intn(50)) * time.Millisecond)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case results <- Result{URL: item.URL, StatusOK: true, Latency: time.Since(start)}:
			}
		}
	}
}

func resultCollector(ctx context.Context, results <-chan Result) error {
	var count int
	for {
		select {
		case <-ctx.Done():
			return nil
		case r, ok := <-results:
			if !ok {
				slog.Info("collection complete", "count", count)
				return nil
			}
			count++
			status := "?"
			if !r.StatusOK {
				status = "?"
			}
			slog.Info("result", "status", status, "url", r.URL, "latency", r.Latency.Round(time.Millisecond))
		}
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	fmt.Println("=== Fan-out pipeline with errgroup.WithContext ===")
	start := time.Now()

	g, ctx := errgroup.WithContext(context.Background())

	jobs := make(chan WorkItem, 10)
	results := make(chan Result, 10)

	g.Go(func() error {
		return producer(ctx, jobs)
	})

	for i := 1; i <= 3; i++ {
		i := i
		g.Go(func() error {
			return consumer(ctx, i, jobs, results)
		})
	}

	go func() {
		g.Wait()
		close(results)
	}()

	var totalResults int
	for r := range results {
		totalResults++
		_ = r
	}

	if err := g.Wait(); err != nil && err != context.Canceled {
		fmt.Printf("? Pipeline failed: %v\n", err)
	} else {
		fmt.Printf("? Pipeline complete: %d results in %v\n", totalResults, time.Since(start).Round(time.Millisecond))
	}

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("?? NEXT UP: CP.3 sync.Pool")
	fmt.Println("   Current: CP.2 (errgroup + context)")
	fmt.Println("---------------------------------------------------")
}
