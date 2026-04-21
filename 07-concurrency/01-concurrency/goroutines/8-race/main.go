// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// Stage 07: Concurrency � Race Conditions & sync.Mutex
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What a race condition is: concurrent access to shared data
//   - Why races are dangerous (data corruption, crashes)
//   - sync.Mutex: mutual exclusion lock
//   - sync/atomic: lock-free atomic operations (fastest option)
//   - The race detector: go run -race (finds races at runtime)
//   - "Share memory by communicating" � channels as the preferred alternative
//
// ANALOGY:
//   Imagine two cashiers updating the same cash register at the same time.
//   Cashier A reads $100, adds $50 ? writes $150.
//   Cashier B reads $100 (before A writes!), adds $30 ? writes $130.
//   A's $50 deposit VANISHED. This is a race condition.
//
//   A MUTEX is like a lock on the register. Only one cashier can open
//   it at a time. The other must wait until the first one is done.
//
// ENGINEERING DEPTH:
//   - Go's race detector uses ThreadSanitizer (from LLVM/Google)
//   - It instruments memory accesses at compile time
//   - Run with: go run -race ./... or go test -race ./...
//   - ALWAYS test with -race in CI/CD pipelines
//
// RUN: go run ./07-concurrency/01-concurrency/goroutines/8-race
// ============================================================================

// unsafeCounter demonstrates what happens WITHOUT synchronization.
// Multiple goroutines increment the same variable concurrently.
// The final count will be WRONG because of lost writes.
func unsafeCounter() int {
	counter := 0 // SHARED STATE � accessed by multiple goroutines
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// THE RACE:
			// 1. Goroutine A reads counter (100)
			// 2. Goroutine B reads counter (100) � same value!
			// 3. A writes counter = 101
			// 4. B writes counter = 101 ? A's write is LOST!
			// This is called a "lost write" or "read-modify-write" race.
			counter++ // NOT THREAD-SAFE � this is THREE operations: read, add, write
		}()
	}
	wg.Wait()
	return counter
}

// mutexCounter uses sync.Mutex to protect the shared variable.
// Only one goroutine can hold the mutex at a time.
func mutexCounter() int {
	counter := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	return counter
}

// atomicCounter uses sync/atomic for lock-free increment.
func atomicCounter() int64 {
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()
	return counter
}

func main() {
	fmt.Println("=== Race Conditions & Synchronization ===")
	fmt.Println()

	fmt.Println("1??  Unsafe (no protection):")
	for i := 0; i < 3; i++ {
		result := unsafeCounter()
		status := "?"
		if result != 1000 {
			status = "? RACE!"
		}
		fmt.Printf("   Run %d: count = %d (expected 1000) %s\n", i+1, result, status)
	}
	fmt.Println()

	fmt.Println("2??  Mutex (sync.Mutex):")
	start := time.Now()
	result := mutexCounter()
	elapsed := time.Since(start)
	fmt.Printf("   count = %d ? (took %v)\n", result, elapsed)
	fmt.Println()

	fmt.Println("3??  Atomic (sync/atomic):")
	start = time.Now()
	atomicResult := atomicCounter()
	elapsed = time.Since(start)
	fmt.Printf("   count = %d ? (took %v)\n", atomicResult, elapsed)
	fmt.Println()

	fmt.Println("=== When to Use Each ===")
	fmt.Println("  +-----------------------------------------------------------+")
	fmt.Println("  � Approach         � Use When                               �")
	fmt.Println("  +------------------+----------------------------------------�")
	fmt.Println("  � sync.Mutex       � Protecting complex shared state        �")
	fmt.Println("  � sync.RWMutex     � Many readers, few writers              �")
	fmt.Println("  � sync/atomic      � Simple counters, flags, single values  �")
	fmt.Println("  � Channels         � Communicating between goroutines       �")
	fmt.Println("  +-----------------------------------------------------------+")
	fmt.Println()

	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Race condition = concurrent unsynchronized access ? data corruption")
	fmt.Println("  - sync.Mutex: Lock/Unlock to protect shared state (most common)")
	fmt.Println("  - sync/atomic: lock-free CPU operations (fastest for simple values)")
	fmt.Println("  - ALWAYS test with: go run -race or go test -race")
	fmt.Println("  - Prefer channels over shared memory when possible")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("?? NEXT UP: GC.10 sync primitives")
	fmt.Println("   Current: GC.8 (race conditions)")
	fmt.Println("---------------------------------------------------")
}
