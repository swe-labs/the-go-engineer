// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"sync"
)

// ============================================================================
// Section 9: Concurrency — Closing Channels
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why and when to close channels
//   - The comma-ok pattern: value, ok := <-ch (detecting closed channels)
//   - range over channels — automatically stops when channel closes
//   - Rules: only the SENDER should close, never the receiver
//   - What happens when you read from / write to a closed channel
//   - Closing as a broadcast signal (notifying multiple goroutines)
//
// ANALOGY:
//   A channel is like a conveyor belt in a factory.
//   close(ch) is like hitting the "STOP" button on the belt.
//   Workers downstream (receivers) see the belt stop and know
//   no more items are coming. They finish processing what's left
//   and go home.
//
// CRITICAL RULES:
//   1. ONLY the sender should close the channel
//   2. Sending to a closed channel causes a PANIC
//   3. Receiving from a closed channel returns zero values immediately
//   4. Closing is optional — only needed when receivers must know "we're done"
//   5. Closing a nil channel causes a PANIC
//
// RUN: go run ./11-concurrency/concurrency/5-channels-closing
// ============================================================================

func main() {
	fmt.Println("=== Closing Channels ===")
	fmt.Println()

	// =====================================================================
	// 1. range over channel — the cleanest pattern
	// =====================================================================
	// range automatically stops when the channel is closed.
	// Without close(), range blocks forever (deadlock).
	fmt.Println("1️⃣  range over channel:")

	taskQueue := make(chan string, 5)

	// Producer: send tasks, then close to signal "done"
	go func() {
		tasks := []string{
			"Compile source code",
			"Run unit tests",
			"Build Docker image",
			"Push to registry",
			"Deploy to staging",
		}
		for _, task := range tasks {
			taskQueue <- task // Send each task
		}
		close(taskQueue) // Signal: no more tasks will be sent
	}()

	// Consumer: range reads until close
	step := 1
	for task := range taskQueue { // Stops automatically when taskQueue is closed
		fmt.Printf("   Step %d: %s ✅\n", step, task)
		step++
	}
	fmt.Println("   Pipeline complete!")
	fmt.Println()

	// =====================================================================
	// 2. comma-ok pattern — detecting closure manually
	// =====================================================================
	// value, ok := <-ch
	//   ok == true  → value is real data
	//   ok == false → channel is closed, value is zero value
	fmt.Println("2️⃣  comma-ok pattern:")

	signals := make(chan int)
	go func() {
		signals <- 42
		signals <- 99
		close(signals) // Close after sending 2 values
	}()

	for {
		value, ok := <-signals // ok tells us if the channel is still open
		if !ok {
			fmt.Println("   Channel closed — no more data")
			break
		}
		fmt.Printf("   Received: %d (ok=%t)\n", value, ok)
	}
	fmt.Println()

	// =====================================================================
	// 3. close() as a BROADCAST signal
	// =====================================================================
	// Closing a channel wakes up ALL goroutines blocked on <-ch.
	// This is a powerful pattern for signaling "shutdown" to many workers.
	fmt.Println("3️⃣  close() as broadcast signal:")

	shutdown := make(chan struct{}) // Empty struct uses zero memory
	var wg sync.WaitGroup

	// Start 3 worker goroutines
	for id := 1; id <= 3; id++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			fmt.Printf("   Worker %d: waiting for shutdown signal...\n", workerID)
			<-shutdown // Blocks until channel is closed
			fmt.Printf("   Worker %d: received shutdown, cleaning up ✅\n", workerID)
		}(id)
	}

	// Close the channel — ALL three workers wake up simultaneously
	fmt.Println("   Main: sending shutdown signal (closing channel)...")
	close(shutdown) // All 3 workers unblock at the same time!

	wg.Wait()
	fmt.Println("   All workers shut down gracefully")

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - close(ch) signals receivers that no more data is coming")
	fmt.Println("  - range ch: reads until closed (cleanest consumer pattern)")
	fmt.Println("  - v, ok := <-ch: ok=false means closed (manual detection)")
	fmt.Println("  - close() wakes ALL blocked receivers (broadcast signal)")
	fmt.Println("  - ONLY senders close channels — NEVER receivers")
	fmt.Println("  - Sending to a closed channel = PANIC (unrecoverable)")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: GC.6 pipeline project")
	fmt.Println("   Current: GC.5 (closing channels)")
	fmt.Println("---------------------------------------------------")
}
