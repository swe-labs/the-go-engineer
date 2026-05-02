// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 07: Concurrency
// Title: Deadlocks
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How circular dependencies cause deadlocks.
//   - How Go's runtime detects total deadlocks.
//   - How to prevent deadlocks using consistent lock ordering.
//
// WHY THIS MATTERS:
//   - A deadlock freezes your program, requiring a hard restart.
//   - They are difficult to debug because they often only happen under specific timing.
//
// RUN:
//   go run ./07-concurrency/01-concurrency/sync-primitives/6-deadlocks
//
// KEY TAKEAWAY:
//   - Always acquire multiple locks in the same order.
// ============================================================================

package main

import (
	"fmt"
	"sync"
	"time"
)

// Resource (Struct): groups the state used by the resource example boundary.
type Resource struct {
	Name string
	mu   sync.Mutex
}

func main() {
	fmt.Println("=== SY.6 Deadlocks ===")
	fmt.Println()

	resA := &Resource{Name: "Resource A"}
	resB := &Resource{Name: "Resource B"}

	fmt.Println("Scenario: Circular Wait Deadlock")
	fmt.Println("  [System] Worker 1 tries to lock A then B")
	fmt.Println("  [System] Worker 2 tries to lock B then A")
	fmt.Println("  [System] This will freeze the program!")
	fmt.Println()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("  [Worker 1] Locking A...")
		resA.mu.Lock()
		time.Sleep(100 * time.Millisecond) // Give worker 2 time to lock B

		fmt.Println("  [Worker 1] Trying to lock B...")
		resB.mu.Lock() // Will block forever!

		fmt.Println("  [Worker 1] Success!")
		resB.mu.Unlock()
		resA.mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		fmt.Println("  [Worker 2] Locking B...")
		resB.mu.Lock()
		time.Sleep(100 * time.Millisecond)

		fmt.Println("  [Worker 2] Trying to lock A...")
		resA.mu.Lock() // Will block forever!

		fmt.Println("  [Worker 2] Success!")
		resA.mu.Unlock()
		resB.mu.Unlock()
	}()

	// Note: We use a separate goroutine to monitor the deadlock
	// because wg.Wait() will also hang.
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("\n  !!! SYSTEM DETECTED HANG !!!")
		fmt.Println("  If this were a total deadlock (no other goroutines running),")
		fmt.Println("  Go would have panicked with 'fatal error: all goroutines are asleep'.")
		fmt.Println("  Because this monitor is running, the runtime doesn't panic yet.")
	}()

	fmt.Println("  [Main] Waiting for workers (expecting hang)...")
	// wg.Wait() // Uncommenting this will hang the whole program

	time.Sleep(3 * time.Second)
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CT.1 -> 07-concurrency/01-concurrency/context/1-background")
	fmt.Println("Current: SY.6 (deadlocks)")
	fmt.Println("Previous: SY.5 (goroutine-leaks)")
	fmt.Println("---------------------------------------------------")
}
