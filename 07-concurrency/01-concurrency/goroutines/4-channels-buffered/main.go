// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"time"
)

// ============================================================================
// Stage 07: Concurrency - Buffered Channels
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Buffered vs unbuffered channels: the critical distinction
//   - make(chan T, capacity): creating a channel with buffer space
//   - When buffered channels block (only when full or empty)
//   - Use cases: batch processing, rate limiting, producer-consumer
//   - When to use buffered vs unbuffered channels
//
// ANALOGY:
//   Unbuffered = a phone call. Both parties must be on the line at the same time.
//   The sender waits until the receiver picks up.
//
//   Buffered = a mailbox with N slots. The sender can drop messages and keep going
//   until the mailbox is full. Then the sender waits.
//
// ENGINEERING DEPTH:
//   Buffered channels use a circular ring buffer internally.
//   - Send blocks only when the buffer is full
//   - Receive blocks only when the buffer is empty
//   - Buffer size is set at creation and cannot be resized
//   - cap(ch) returns the buffer capacity
//   - len(ch) returns the number of items currently in the buffer
//
// RUN: go run ./07-concurrency/01-concurrency/goroutines/4-channels-buffered
// ============================================================================

type logEvent struct {
	Level   string
	Message string
}

func main() {
	fmt.Println("=== Buffered Channels ===")
	fmt.Println()

	events := make(chan logEvent, 3)
	events <- logEvent{"INFO", "Server started on :8080"}
	events <- logEvent{"INFO", "Connected to database"}
	events <- logEvent{"WARN", "Cache miss rate above 50%"}

	fmt.Printf("  Buffer: %d/%d items\n\n", len(events), cap(events))

	fmt.Println("  1) Basic buffered channel (capacity=3):")
	for i := 0; i < 3; i++ {
		e := <-events
		fmt.Printf("     [%s] %s\n", e.Level, e.Message)
	}
	fmt.Println()

	fmt.Println("  2) Producer-consumer pattern:")
	jobs := make(chan int, 5)

	go func() {
		for i := 1; i <= 8; i++ {
			fmt.Printf("     -> Producing job #%d\n", i)
			jobs <- i
		}
		close(jobs)
	}()

	time.Sleep(50 * time.Millisecond)
	for job := range jobs {
		fmt.Printf("     -> Processing job #%d\n", job)
		time.Sleep(30 * time.Millisecond)
	}
	fmt.Println()

	fmt.Println("  3) Buffered vs unbuffered:")
	fmt.Println("     - Unbuffered: make(chan T) -> sender waits for receiver")
	fmt.Println("     - Buffered:   make(chan T, N) -> sender waits only when buffer is full")
	fmt.Println("     - Unbuffered emphasizes synchronization")
	fmt.Println("     - Buffered emphasizes bounded asynchronous flow")

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Buffered: make(chan T, N) -> N items can be sent without blocking")
	fmt.Println("  - Unbuffered: make(chan T) -> sender waits for receiver (synchronous)")
	fmt.Println("  - Use buffered channels to decouple fast producers from slow consumers")
	fmt.Println("  - Buffer size should be tuned based on throughput needs")
	fmt.Println("  - When in doubt, start unbuffered -> add buffer only for performance")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: GC.5 closing channels")
	fmt.Println("   Current: GC.4 (buffered channels)")
	fmt.Println("---------------------------------------------------")
}
