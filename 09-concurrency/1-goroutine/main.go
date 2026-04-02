// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// Section 9: Concurrency — Goroutines
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - What a goroutine is: a lightweight thread managed by Go's runtime
//   - The "go" keyword: launching concurrent execution
//   - Why time.Sleep is a BAD way to wait (and what to use instead)
//   - Goroutine lifecycle: creation, execution, termination
//   - How goroutines differ from OS threads
//
// ANALOGY:
//   Imagine a restaurant kitchen. The head chef (main goroutine) has a
//   list of orders to prepare. Instead of making each dish ONE AT A TIME,
//   the chef says "go make the salad" to one cook, "go grill the steak"
//   to another, and "go prepare dessert" to a third.
//
//   Each cook works INDEPENDENTLY and CONCURRENTLY. The head chef doesn't
//   wait for the salad to finish before starting the steak.
//   That's exactly what goroutines do — concurrent, independent work.
//
// ENGINEERING DEPTH:
//   - A goroutine starts with only ~2KB of stack (vs 1-8MB for an OS thread)
//   - Go can run MILLIONS of goroutines on a single machine
//   - The Go runtime multiplexes goroutines onto a small number of OS threads
//     using an M:N scheduler (M goroutines on N OS threads)
//   - Goroutines are NOT parallelism — they're concurrency (tasks can interleave)
//   - True parallelism requires multiple CPU cores (GOMAXPROCS > 1)
//
// RUN: go run ./09-concurrency/1-goroutine
// ============================================================================

// processOrder simulates a kitchen worker preparing a dish.
// Each dish takes a different amount of time.
func processOrder(dish string, prepTime time.Duration) {
	fmt.Printf("  🍳 Started preparing: %s\n", dish)
	time.Sleep(prepTime) // Simulate work (cooking time)
	fmt.Printf("  ✅ Finished: %s (took %v)\n", dish, prepTime)
}

func main() {
	fmt.Println("=== Goroutines: Lightweight Concurrent Execution ===")
	fmt.Println()

	// Record the start time to measure total duration.
	start := time.Now()

	// --- SEQUENTIAL EXECUTION (Without goroutines) ---
	// If we called these WITHOUT "go", they'd run one after another:
	//   processOrder("Salad", 1*time.Second)        ← 1 second
	//   processOrder("Steak", 2*time.Second)        ← 2 seconds
	//   processOrder("Pasta", 1500*time.Millisecond) ← 1.5 seconds
	//   Total: ~4.5 seconds (serial, one after another)

	// --- CONCURRENT EXECUTION (With goroutines) ---
	// The "go" keyword starts a function in a NEW goroutine.
	// The current function continues IMMEDIATELY — it doesn't wait.
	// All three dishes cook AT THE SAME TIME.
	fmt.Println("Kitchen is open! Starting all orders concurrently...")
	fmt.Println()

	// --- THE RIGHT WAY: sync.WaitGroup ---
	// WaitGroup tracks running goroutines. It's Go's standard tool for
	// waiting for a group of goroutines to finish.
	//
	// Pattern:
	//   var wg sync.WaitGroup
	//   wg.Add(1)         ← "One more goroutine to wait for"
	//   go func() {
	//       defer wg.Done() ← "This goroutine is done"
	//       // ... do work ...
	//   }()
	//   wg.Wait()         ← "Block until all goroutines call Done()"
	var wg sync.WaitGroup

	// Launch 4 goroutines — each prepares a dish concurrently
	dishes := []struct {
		name string
		time time.Duration
	}{
		{"Caesar Salad", 500 * time.Millisecond},
		{"Grilled Steak", 800 * time.Millisecond},
		{"Mushroom Pasta", 600 * time.Millisecond},
		{"Chocolate Cake", 700 * time.Millisecond},
	}

	for _, dish := range dishes {
		wg.Add(1) // Tell WaitGroup: "one more goroutine to track"

		// IMPORTANT: We pass dish.name and dish.time as parameters.
		// If we used the loop variable directly, ALL goroutines would
		// see the LAST value of the loop (a common closure bug).
		go func(name string, prepTime time.Duration) {
			defer wg.Done() // When this goroutine finishes, decrement the counter
			processOrder(name, prepTime)
		}(dish.name, dish.time) // Pass current values into the goroutine
	}

	// wg.Wait() blocks until ALL goroutines have called wg.Done().
	// This is MUCH better than time.Sleep — it waits exactly as long as needed.
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println()
	fmt.Printf("🎉 All orders complete! Total time: %v\n", elapsed)
	fmt.Println("   (Sequential would have been ~2.6s — goroutines did it concurrently!)")

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - 'go functionName()' starts a function in a new goroutine")
	fmt.Println("  - Goroutines are lightweight (~2KB each) — launch millions of them")
	fmt.Println("  - NEVER use time.Sleep to wait — use sync.WaitGroup or channels")
	fmt.Println("  - Pass loop variables as parameters to avoid closure bugs")
	fmt.Println("  - The main goroutine must wait, or it will exit (killing all goroutines)")
	fmt.Println("  - Next: go run ./09-concurrency/2-wait-group (deeper WaitGroup patterns)")
}
