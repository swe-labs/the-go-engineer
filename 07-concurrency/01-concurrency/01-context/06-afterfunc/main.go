// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 07: Concurrency
// Title: context.AfterFunc
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Schedule a callback function to run when a context is cancelled.
//   - Use context.AfterFunc to trigger cleanup, logging, or fallback actions.
//   - Understand the difference between AfterFunc and goroutines on ctx.Done().
//
// WHY THIS MATTERS:
//   context.AfterFunc (Go 1.21+) lets you attach side effects to context
//   cancellation without starting a dedicated goroutine. It's more efficient
//   than a select-on-ctx.Done() pattern when you only need a single callback.
//
// RUN:
//   go run ./07-concurrency/01-concurrency/01-context/06-afterfunc
//
// KEY TAKEAWAY:
//   - AfterFunc schedules one callback when a context is cancelled.
//   - The callback runs in its own goroutine — do not block on it.
//   - If the context is already done, the callback fires immediately.
//   - Call the returned stop function to cancel the callback if it hasn't fired.
// ============================================================================

package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// Example 1: AfterFunc with WithTimeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan struct{})

	stop := context.AfterFunc(ctx, func() {
		fmt.Println("[AfterFunc] Timeout reached — cleaning up resources...")
		close(done)
	})
	defer stop()

	fmt.Println("Waiting for timeout...")
	<-done
	fmt.Println("Main: AfterFunc completed.")
	fmt.Println()

	// Example 2: AfterFunc fires immediately if context is already done.
	cancelledCtx, cancelNow := context.WithCancel(context.Background())
	cancelNow() // ctx is already cancelled

	fired := false
	stop2 := context.AfterFunc(cancelledCtx, func() {
		fired = true
	})
	defer stop2()

	// The AfterFunc should fire almost immediately since the context is done.
	time.Sleep(10 * time.Millisecond)
	fmt.Printf("AfterFunc on cancelled context fired: %v\n", fired)
	fmt.Println()

	// Example 3: Cancelling the AfterFunc via stop.
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	didNotFire := true
	stop3 := context.AfterFunc(ctx2, func() {
		didNotFire = false
	})
	stop3() // Cancel the callback before the context is cancelled.
	cancel2()
	time.Sleep(5 * time.Millisecond)
	fmt.Printf("Stopped AfterFunc did not fire: %v\n", didNotFire)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TM.1 -> 07-concurrency/01-concurrency/04-time-and-scheduling/01-time")
	fmt.Println("Run    : go run ./07-concurrency/01-concurrency/04-time-and-scheduling/01-time")
	fmt.Println("Current: CT.6 (afterfunc)")
	fmt.Println("---------------------------------------------------")
}
