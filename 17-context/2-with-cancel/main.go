// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"context"
	"fmt"
	"time"
)

// ============================================================================
// Section 17: Context — WithCancel
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - context.WithCancel creates a cancellable context
//   - The cancel function: how and when to call it
//   - Listening for cancellation with ctx.Done()
//   - Using select to handle cancellation in goroutines
//   - Why you MUST always call cancel() (preventing goroutine leaks)
//
// ENGINEERING DEPTH:
//   When you invoke `context.WithCancel(parent)`, Go creates a `cancelCtx` struct
//   under the hood. This struct contains a Mutex and a `done` channel. It also
//   recursively climbs up the Context tree until it finds the first cancellable
//   parent, and appends itself to that parent's internal array of children!
//   When you call `cancel()`, Go closes the `done` channel, broadcasts the close
//   signal to all children in the array, and then removes itself from the parent
//   to allow the Garbage Collector to sweep the dead goroutines.
//
// RUN: go run ./17-context/2-with-cancel
// ============================================================================

func main() {
	fmt.Println("=== Context: WithCancel ===")
	fmt.Println()

	// --- CREATING A CANCELLABLE CONTEXT ---
	// WithCancel returns two values:
	//   1. ctx — a new context that can be cancelled
	//   2. cancel — a function that triggers the cancellation
	//
	// CRITICAL RULE: You MUST call cancel() when you're done with the context.
	// If you don't, the context's resources (goroutines) are never freed.
	// Use defer cancel() immediately after creation.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Always defer cancel — even if you call it explicitly later

	// Start a background worker that listens for cancellation
	results := make(chan string)
	go worker(ctx, results)

	// Collect some results from the worker
	for i := 0; i < 3; i++ {
		fmt.Printf("  Received: %s\n", <-results)
	}

	// Now cancel the context — this signals the worker to stop
	fmt.Println("\n  Calling cancel()...")
	cancel()

	// Give the worker a moment to notice the cancellation
	time.Sleep(100 * time.Millisecond)

	// After cancellation, ctx.Err() returns a non-nil error
	fmt.Printf("  Context error after cancel: %v\n", ctx.Err())
	// ctx.Err() returns context.Canceled

	fmt.Println()

	// --- NESTED CANCELLATION ---
	// When a parent context is cancelled, ALL children are cancelled too.
	fmt.Println("=== Nested Cancellation ===")
	parentCtx, parentCancel := context.WithCancel(context.Background())

	// Create child contexts from the parent
	childCtx1, childCancel1 := context.WithCancel(parentCtx)
	childCtx2, childCancel2 := context.WithCancel(parentCtx)
	defer childCancel1()
	defer childCancel2()

	fmt.Printf("  Before cancel — parent err: %v, child1 err: %v, child2 err: %v\n",
		parentCtx.Err(), childCtx1.Err(), childCtx2.Err())

	// Cancel the PARENT — both children are cancelled automatically
	parentCancel()

	fmt.Printf("  After parent cancel — parent err: %v, child1 err: %v, child2 err: %v\n",
		parentCtx.Err(), childCtx1.Err(), childCtx2.Err())

	fmt.Println()
	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  1. WithCancel returns (ctx, cancel) — ALWAYS defer cancel()")
	fmt.Println("  2. Listen for cancellation with <-ctx.Done() in a select")
	fmt.Println("  3. After cancellation, ctx.Err() returns context.Canceled")
	fmt.Println("  4. Cancelling a parent cancels ALL children automatically")
	fmt.Println("  5. Forgetting cancel() causes GOROUTINE LEAKS (memory leak)")
	fmt.Println()
	fmt.Println("   Next: go run ./17-context/3-with-timeout")
}

// worker simulates a long-running task that checks for cancellation.
// This is the standard pattern for cancellation-aware goroutines.
//
// The select statement waits for EITHER:
//   - ctx.Done() — the context was cancelled (stop working)
//   - default/time — normal work continues
func worker(ctx context.Context, results chan<- string) {
	i := 0
	for {
		select {
		case <-ctx.Done():
			// The context was cancelled — clean up and return
			fmt.Printf("  Worker stopped: %v\n", ctx.Err())
			return
		default:
			// Context is still active — do work
			i++
			results <- fmt.Sprintf("result-%d", i)
			time.Sleep(50 * time.Millisecond) // Simulate work
		}
	}
}
