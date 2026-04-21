// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"sync"
)

// ============================================================================
// Stage 07: Concurrency - Closing Channels
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why and when to close channels
//   - The comma-ok pattern: value, ok := <-ch (detecting closed channels)
//   - range over channels: automatically stops when the channel closes
//   - Rules: only the sender should close, never the receiver
//   - What happens when you read from or write to a closed channel
//   - Closing as a broadcast signal for multiple goroutines
//
// ANALOGY:
//   A channel is like a conveyor belt in a factory.
//   close(ch) is like hitting the STOP button on the belt.
//   Workers downstream see the belt stop and know no more items are coming.
//
// CRITICAL RULES:
//   1. Only the sender should close the channel
//   2. Sending to a closed channel causes a panic
//   3. Receiving from a closed channel returns zero values immediately
//   4. Closing is optional unless receivers must know "we are done"
//   5. Closing a nil channel causes a panic
//
// RUN: go run ./07-concurrency/01-concurrency/goroutines/5-channels-closing
// ============================================================================

func main() {
	fmt.Println("=== Closing Channels ===")
	fmt.Println()

	fmt.Println("1) range over channel:")
	taskQueue := make(chan string, 5)

	go func() {
		tasks := []string{
			"Compile source code",
			"Run unit tests",
			"Build Docker image",
			"Push to registry",
			"Deploy to staging",
		}
		for _, task := range tasks {
			taskQueue <- task
		}
		close(taskQueue)
	}()

	step := 1
	for task := range taskQueue {
		fmt.Printf("   Step %d: %s -> done\n", step, task)
		step++
	}
	fmt.Println("   Pipeline complete!")
	fmt.Println()

	fmt.Println("2) comma-ok pattern:")
	signals := make(chan int)
	go func() {
		signals <- 42
		signals <- 99
		close(signals)
	}()

	for {
		value, ok := <-signals
		if !ok {
			fmt.Println("   Channel closed - no more data")
			break
		}
		fmt.Printf("   Received: %d (ok=%t)\n", value, ok)
	}
	fmt.Println()

	fmt.Println("3) close() as broadcast signal:")
	shutdown := make(chan struct{})
	var wg sync.WaitGroup

	for id := 1; id <= 3; id++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			fmt.Printf("   Worker %d: waiting for shutdown signal...\n", workerID)
			<-shutdown
			fmt.Printf("   Worker %d: received shutdown, cleaning up\n", workerID)
		}(id)
	}

	fmt.Println("   Main: sending shutdown signal (closing channel)...")
	close(shutdown)

	wg.Wait()
	fmt.Println("   All workers shut down gracefully")

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - close(ch) signals receivers that no more data is coming")
	fmt.Println("  - range ch: reads until closed (cleanest consumer pattern)")
	fmt.Println("  - v, ok := <-ch: ok=false means closed (manual detection)")
	fmt.Println("  - close() wakes all blocked receivers (broadcast signal)")
	fmt.Println("  - Only senders close channels - never receivers")
	fmt.Println("  - Sending to a closed channel = panic (unrecoverable)")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: GC.6 pipeline project")
	fmt.Println("   Current: GC.5 (closing channels)")
	fmt.Println("---------------------------------------------------")
}
