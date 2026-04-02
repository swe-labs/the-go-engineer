// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Section 9: Concurrency — Pipeline Project
// Level: Advanced
// ============================================================================
//
// RUN: go run ./09-concurrency/6-project-1
// ============================================================================

import (
	"context"
	"fmt"
	"time"
)

// ping acts as a Goroutine "Actor" in an infinite loop.
// It shares the `ch` channel memory reference with the `pong` actor.
func ping(ctx context.Context, ch chan string) {
	for {
		// 1. The Select Multiplexer
		// `select` blocks the Goroutine until one of its cases becomes unblocked.
		// If both are unblocked, Go picks one randomly.
		select {
		case <-ctx.Done():
			// 2. Context Cancellation
			// If the parent calls `cancel()` or times out, the `Done()` channel
			// is closed. Reading from a closed channel instantly returns, allowing
			// this Goroutine to gracefully exit (preventing memory leaks).
			return
		case ch <- fmt.Sprintf("ping: %v", time.Now()):
			// 3. Unbuffered Write Blocking
			// Writing to `ch` BLOCKS this Goroutine until `main` reads from it!
			// After the write succeeds, we sleep to throttle the output.
			time.Sleep(1 * time.Second)
		}
	}
}

func pong(ctx context.Context, ch chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case ch <- fmt.Sprintf("pong: %v", time.Now()):
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pingerCh := make(chan string)
	done := make(chan struct{})

	go ping(ctx, pingerCh)
	go pong(ctx, pingerCh)

	// 4. Background Coordinator Goroutine
	// This anonymous function spins up to coordinate the system's shutdown loop.
	go func() {
		// time.After creates a hardware timer that fires a message into the channel
		// after exactly 5 seconds.
		timeout := time.After(5 * time.Second)
		for {
			select {
			case <-timeout:
				fmt.Println("operation completed")
				// 5. Channel Teardown
				// Context cancellation (cancel() above) handled shutdown gracefully.
				// We DO NOT close pingerCh here to avoid data races and "send on closed channel" panics
				// from ping/pong goroutines that are concurrently trying to write.
				done <- struct{}{} // Signal the main thread we are finished
				return
			case msg := <-pingerCh:
				// As long as `ping` and `pong` write to `pingerCh`, we print it here.
				fmt.Println(msg)
			}
		}
	}()

	<-done
	fmt.Println("done")

}
