// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 07: Concurrency
// Title: Goroutine Leaks
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How goroutines get stranded in memory (Leaks).
//   - How to detect leaks using runtime.NumGoroutine().
//   - How to prevent leaks using context cancellation.
//
// WHY THIS MATTERS:
//   - Leaked goroutines consume stack memory (2KB+ each).
//   - Thousands of leaks will eventually cause an OOM (Out of Memory) crash.
//
// RUN:
//   go run ./07-concurrency/01-concurrency/sync-primitives/5-goroutine-leaks
//
// KEY TAKEAWAY:
//   - Never start a goroutine without knowing how it will stop.
// ============================================================================

package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// leakGenerator creates a goroutine that blocks forever on a channel send.
// leakGenerator (Function): creates a goroutine that blocks forever on a channel send.
func leakGenerator() {
	ch := make(chan int)
	go func() {
		fmt.Println("  [Leak] Goroutine started, trying to send...")
		ch <- 42 // Blocks forever because no one is receiving!
		fmt.Println("  [Leak] This line will never be reached.")
	}()
}

// safeWorker creates a goroutine that respects context cancellation.
// safeWorker (Function): creates a goroutine that respects context cancellation.
func safeWorker(ctx context.Context) {
	go func() {
		fmt.Println("  [Safe] Goroutine started, waiting for work or cancel...")
		select {
		case <-ctx.Done():
			fmt.Println("  [Safe] Context cancelled, exiting gracefully.")
			return
		case <-time.After(1 * time.Hour):
			// simulate long work
		}
	}()
}

func main() {
	fmt.Println("=== SY.5 Goroutine Leaks ===")
	fmt.Println()

	initial := runtime.NumGoroutine()
	fmt.Printf("Initial Goroutines: %d\n", initial)

	// 1. Create a leak
	fmt.Println("\nScenario 1: Creating a leak...")
	leakGenerator()
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Current Goroutines: %d (Leaked: %d)\n", runtime.NumGoroutine(), runtime.NumGoroutine()-initial)

	// 2. Create a safe worker
	fmt.Println("\nScenario 2: Creating a safe worker with context...")
	ctx, cancel := context.WithCancel(context.Background())
	safeWorker(ctx)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Current Goroutines: %d\n", runtime.NumGoroutine())

	fmt.Println("Cancelling safe worker...")
	cancel()
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Final Goroutines: %d (Leaked remains, safe exited)\n", runtime.NumGoroutine())

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: SY.6 -> 07-concurrency/01-concurrency/sync-primitives/6-deadlocks")
	fmt.Println("Current: SY.5 (goroutine leaks)")
	fmt.Println("Previous: SY.4 (race conditions)")
	fmt.Println("---------------------------------------------------")
}
