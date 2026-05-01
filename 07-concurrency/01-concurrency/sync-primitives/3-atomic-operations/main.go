// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 07: Concurrency
// Title: Atomic Operations
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use sync/atomic for lock-free state updates.
//   - The difference between Mutex and Hardware-level atomics.
//   - How to use atomic.Value for complex read-only snapshots.
//
// WHY THIS MATTERS:
//   - Atomics are the fastest way to sync small counters or flags.
//   - They avoid the overhead of the OS scheduler (parking threads).
//
// RUN:
//   go run ./07-concurrency/01-concurrency/sync-primitives/3-atomic-operations
//
// KEY TAKEAWAY:
//   - Use atomics for simple counters; use Mutex for multiple-field invariants.
// ============================================================================

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Metrics struct {
	requests int64
	errors   int64
	active   int64
}

func main() {
	fmt.Println("=== SY.3 Atomic Operations ===")
	fmt.Println()

	// 1. Primitive Atomics (int64)
	var metrics Metrics
	var wg sync.WaitGroup

	fmt.Println("Scenario 1: High-Frequency Counter (sync/atomic)")
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&metrics.requests, 1)
			atomic.AddInt64(&metrics.active, 1)
			// ... simulate some work ...
			atomic.AddInt64(&metrics.active, -1)
		}()
	}
	wg.Wait()
	fmt.Printf("  Total Requests: %d\n", atomic.LoadInt64(&metrics.requests))
	fmt.Printf("  Active Workers: %d (Expected: 0)\n", atomic.LoadInt64(&metrics.active))

	// 2. Atomic Value (Snapshotting)
	// atomic.Value is great for replacing a whole config or map at once.
	fmt.Println("\nScenario 2: Config Snapshotting (atomic.Value)")
	var config atomic.Value
	config.Store("Initial Config")

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond)
		config.Store("Updated Production Config")
		fmt.Println("  [System] Config rotated!")
	}()

	for i := 0; i < 3; i++ {
		fmt.Printf("  Read %d: %v\n", i+1, config.Load())
		time.Sleep(40 * time.Millisecond)
	}
	wg.Wait()

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: SY.4 -> 07-concurrency/01-concurrency/sync-primitives/4-race-conditions")
	fmt.Println("Current: SY.3 (atomic operations)")
	fmt.Println("Previous: SY.2 (sync.once and sync.map)")
	fmt.Println("---------------------------------------------------")
}
