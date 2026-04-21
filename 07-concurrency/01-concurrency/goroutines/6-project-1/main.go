// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Stage 07: Concurrency - Pipeline Project
// Level: Advanced
// ============================================================================
//
// RUN: go run ./07-concurrency/01-concurrency/goroutines/6-project-1
// ============================================================================

import (
	"context"
	"fmt"
	"time"
)

// ping acts as a goroutine "actor" in an infinite loop.
// It shares the `ch` channel memory reference with the `pong` actor.
func ping(ctx context.Context, ch chan string) {
	for {
		// 1. The select multiplexer
		// `select` blocks the goroutine until one of its cases becomes unblocked.
		// If both are unblocked, Go picks one randomly.
		select {
		case <-ctx.Done():
			// 2. Context cancellation
			// If the parent calls cancel() or times out, the Done channel is closed.
			// Reading from a closed channel instantly returns, allowing this
			// goroutine to exit gracefully and avoid leaks.
			return
		case ch <- fmt.Sprintf("ping: %v", time.Now()):
			// 3. Unbuffered write blocking
			// Writing to `ch` blocks this goroutine until `main` reads from it.
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

	// 4. Background coordinator goroutine
	go func() {
		// time.After creates a timer that delivers on its channel after 5 seconds.
		timeout := time.After(5 * time.Second)
		for {
			select {
			case <-timeout:
				fmt.Println("operation completed")
				// 5. Channel teardown
				// Context cancellation handles shutdown gracefully.
				// We do not close pingerCh here to avoid send-on-closed-channel panics.
				done <- struct{}{}
				return
			case msg := <-pingerCh:
				fmt.Println(msg)
			}
		}
	}()

	<-done
	fmt.Println("done")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: GC.7 concurrent downloader")
	fmt.Println("   Current: GC.6 (pipeline project)")
	fmt.Println("---------------------------------------------------")
}
