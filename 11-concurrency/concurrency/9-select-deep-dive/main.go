// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Section 9: Concurrency — Select Deep Dive
// Level: Advanced
// ============================================================================
//
// RUN: go run ./11-concurrency/concurrency/9-select-deep-dive
// ============================================================================

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ============================================================================
// Section 9: Select Statement Deep Dive
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - select for multiplexing channel operations
//   - Timeouts with time.After
//   - Non-blocking channel operations with default
//   - Context-based cancellation
//   - Fan-in pattern (merging multiple channels)
// ============================================================================

func main() {
	fmt.Println("=== Select Statement Deep Dive ===")
	fmt.Println()

	// 1. Basic select — wait on multiple channels
	basicSelect()

	// 2. Timeout pattern — prevent waiting forever
	timeoutPattern()

	// 3. Non-blocking operations with default
	nonBlockingSelect()

	// 4. Context cancellation — production-standard cancellation
	contextCancellation()

	// 5. Fan-in — merging multiple channels into one
	fanInPattern()
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: GC.10 sync primitives")
	fmt.Println("   Current: GC.9 (select deep dive)")
	fmt.Println("---------------------------------------------------")
}

// basicSelect demonstrates waiting on multiple channels.
// select blocks until ONE of the cases can proceed.
func basicSelect() {
	fmt.Println("--- 1. Basic Select ---")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "from channel 1"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "from channel 2"
	}()

	// Receive from whichever channel is ready first
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Printf("  Received: %s\n", msg)
		case msg := <-ch2:
			fmt.Printf("  Received: %s\n", msg)
		}
	}
	fmt.Println()
}

// timeoutPattern prevents goroutines from waiting forever.
func timeoutPattern() {
	fmt.Println("--- 2. Timeout Pattern ---")

	slowOperation := make(chan string)

	go func() {
		time.Sleep(2 * time.Second) // Simulates a slow operation
		slowOperation <- "done"
	}()

	select {
	case result := <-slowOperation:
		fmt.Printf("  Got result: %s\n", result)
	case <-time.After(500 * time.Millisecond):
		// time.After returns a channel that receives after the duration.
		// This is the standard timeout pattern in Go.
		fmt.Println("  ⏰ Operation timed out after 500ms")
	}
	fmt.Println()
}

// nonBlockingSelect uses the default case to avoid blocking.
func nonBlockingSelect() {
	fmt.Println("--- 3. Non-blocking Select ---")

	ch := make(chan int, 1) // Buffered channel

	// Try to receive without blocking
	select {
	case val := <-ch:
		fmt.Printf("  Received: %d\n", val)
	default:
		fmt.Println("  Channel empty — no blocking!")
	}

	// Try to send without blocking
	ch <- 42
	select {
	case ch <- 100: // Buffer is full, this won't execute
		fmt.Println("  Sent 100")
	default:
		fmt.Println("  Channel full — no blocking!")
	}
	fmt.Println()
}

// contextCancellation demonstrates production-standard cancellation.
func contextCancellation() {
	fmt.Println("--- 4. Context Cancellation ---")

	// context.WithTimeout creates a context that auto-cancels after duration.
	// This is the PRODUCTION pattern (preferred over time.After for goroutines).
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel() // Always call cancel to release resources

	resultCh := make(chan string)

	go func() {
		// Simulate work
		time.Sleep(500 * time.Millisecond)
		resultCh <- "work complete"
	}()

	select {
	case result := <-resultCh:
		fmt.Printf("  Result: %s\n", result)
	case <-ctx.Done():
		// ctx.Done() returns a channel that closes when the context is cancelled.
		fmt.Printf("  ❌ Cancelled: %v\n", ctx.Err())
	}
	fmt.Println()
}

// fanInPattern merges multiple input channels into a single output channel.
func fanInPattern() {
	fmt.Println("--- 5. Fan-In Pattern ---")

	// Create 3 producer channels
	producers := make([]<-chan string, 3)
	for i := 0; i < 3; i++ {
		producers[i] = produce(i)
	}

	// Merge all into one channel
	merged := fanIn(producers...)

	// Read all results from the merged channel
	for i := 0; i < 6; i++ { // 3 producers × 2 messages each
		fmt.Printf("  %s\n", <-merged)
	}
	fmt.Println()
}

func produce(id int) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for i := 0; i < 2; i++ {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			ch <- fmt.Sprintf("Producer %d: message %d", id, i)
		}
	}()
	return ch
}

// fanIn merges multiple channels into one using a WaitGroup for cleanup.
func fanIn(channels ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan string) {
			defer wg.Done()
			for msg := range c {
				out <- msg
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
