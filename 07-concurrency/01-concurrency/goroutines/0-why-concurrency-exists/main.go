// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 07: Concurrency
// Title: Why concurrency exists
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why overlapping waits is the core reason for concurrency.
//   - The difference between CPU-bound and I/O-bound tasks.
//   - How Go's scheduler exploits idle time to improve throughput.
//
// WHY THIS MATTERS:
//   - Most backend work is spent waiting (for DBs, APIs, or disks).
//   - Concurrency allows your program to make progress while waiting.
//
// RUN:
//   go run ./07-concurrency/01-concurrency/goroutines/0-why-concurrency-exists
//
// KEY TAKEAWAY:
//   - Concurrency is not about running faster; it's about overlapping waits.
// ============================================================================

package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== GC.0 Why Concurrency Exists ===")
	fmt.Println()

	// In a real backend, these might be database queries or API calls.
	tasks := []string{"DB Query", "API Fetch", "Disk Read"}
	taskDuration := 100 * time.Millisecond

	fmt.Println("Scenario 1: Sequential Execution (Traditional)")
	start := time.Now()
	for _, task := range tasks {
		performTask(task, taskDuration)
	}
	fmt.Printf("Total Time Sequential: %v\n", time.Since(start))
	fmt.Println("Wait time is cumulative: 100ms + 100ms + 100ms = 300ms")

	fmt.Println("\nScenario 2: Concurrent Execution (Overlap)")
	start = time.Now()

	// NOTE: This is a simplified "fake" concurrency example to show the concept.
	// In the next lesson (GC.1), we will learn how to use real goroutines.
	// Here, we simulate that all three tasks start at the same time.
	fmt.Println("  [System] Launching all tasks simultaneously...")

	// We simulate the overlap by only sleeping once for the duration of the longest task.
	time.Sleep(taskDuration)

	fmt.Printf("Total Time Concurrent: ~%v\n", time.Since(start))
	fmt.Println("Wait time is overlapping: max(100ms, 100ms, 100ms) = 100ms")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: GC.1 -> 07-concurrency/01-concurrency/goroutines/1-goroutine")
	fmt.Println("Current: GC.0 (why concurrency exists)")
	fmt.Println("Previous: DB.8 (query-timeouts-via-context)")
	fmt.Println("---------------------------------------------------")
}

// performTask (Function): runs the perform task step and keeps its inputs, outputs, or errors visible.
func performTask(name string, duration time.Duration) {
	fmt.Printf("  [Wait] Starting %s...\n", name)
	time.Sleep(duration) // Simulating a network/disk wait
	fmt.Printf("  [Done] Finished %s\n", name)
}
