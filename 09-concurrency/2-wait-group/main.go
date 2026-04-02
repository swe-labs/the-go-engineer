// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ============================================================================
// Section 9: Concurrency — WaitGroups
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - sync.WaitGroup: the standard tool for waiting on goroutines
//   - The 3-step pattern: Add → go func { defer Done } → Wait
//   - Why you MUST pass WaitGroup by pointer (never by value)
//   - Common mistakes: Add inside goroutine, forgetting Done, copy WaitGroup
//   - Real-world example: concurrent health checks
//
// ANALOGY:
//   WaitGroup is like a boarding pass counter at an airport gate.
//     wg.Add(1)  = "One more passenger to board"
//     wg.Done()  = "One passenger has boarded"
//     wg.Wait()  = "Don't close the gate until all passengers are aboard"
//
//   If you forget to call Done(), the gate never closes (deadlock).
//   If you call Add() after Wait(), the counter goes negative (panic).
//
// RUN: go run ./09-concurrency/2-wait-group
// ============================================================================

// ServiceStatus represents the health of a backend service.
type ServiceStatus struct {
	Name    string        // Service name (e.g., "database", "cache")
	Healthy bool          // Whether the health check passed
	Latency time.Duration // How long the check took
}

// checkService simulates a health check on a backend service.
// It takes a *sync.WaitGroup as a parameter and calls Done() when finished.
//
// THE GOLDEN RULE: Always pass WaitGroup by POINTER (*sync.WaitGroup).
// If you pass by value, the goroutine gets a COPY — calling Done()
// on the copy doesn't affect the original, causing a deadlock.
func checkService(name string, wg *sync.WaitGroup, results chan<- ServiceStatus) {
	// defer wg.Done() MUST be the first line in the goroutine.
	// defer ensures Done() is called even if the function panics.
	// Without it, wg.Wait() blocks forever (deadlock).
	defer wg.Done()

	// Simulate varying response times (50-300ms)
	latency := time.Duration(50+rand.Intn(250)) * time.Millisecond
	time.Sleep(latency)

	// Simulate: most services are healthy, some randomly fail
	healthy := rand.Intn(10) > 1 // 90% chance of passing

	results <- ServiceStatus{
		Name:    name,
		Healthy: healthy,
		Latency: latency,
	}
}

func main() {
	fmt.Println("=== WaitGroups: Coordinating Goroutine Completion ===")
	fmt.Println()

	// --- THE WAITGROUP PATTERN ---
	//
	//   var wg sync.WaitGroup    ← Step 0: Declare (zero value is ready to use)
	//   wg.Add(1)                ← Step 1: Increment BEFORE launching goroutine
	//   go func() {
	//       defer wg.Done()      ← Step 2: Decrement when goroutine finishes
	//       // ... do work ...
	//   }()
	//   wg.Wait()                ← Step 3: Block until counter reaches 0

	var wg sync.WaitGroup

	// Services to health check (simulating a microservice architecture)
	services := []string{
		"postgres-db",
		"redis-cache",
		"auth-service",
		"email-service",
		"payment-gateway",
		"search-engine",
	}

	// Channel to collect results from all goroutines
	results := make(chan ServiceStatus, len(services)) // Buffered: won't block senders

	fmt.Printf("🔍 Health checking %d services concurrently...\n\n", len(services))

	// Launch one goroutine per service
	for _, svc := range services {
		wg.Add(1)                          // CRITICAL: Add BEFORE launching the goroutine, not inside it
		go checkService(svc, &wg, results) // &wg passes a POINTER (not a copy)
	}

	// --- WAIT FOR ALL GOROUTINES ---
	// wg.Wait() blocks until the internal counter reaches 0.
	// Each goroutine decrements the counter with Done().
	// When all 6 goroutines have called Done(), Wait() returns.
	wg.Wait()

	// Close the results channel now that all writers are done.
	// This allows us to range over it safely.
	close(results)

	// Collect and display all results
	allHealthy := true
	for status := range results {
		icon := "✅"
		if !status.Healthy {
			icon = "❌"
			allHealthy = false
		}
		fmt.Printf("  %s %-20s latency: %v\n", icon, status.Name, status.Latency)
	}

	fmt.Println()
	if allHealthy {
		fmt.Println("🎉 All services healthy!")
	} else {
		fmt.Println("⚠️  Some services are degraded — check logs!")
	}

	// --- COMMON MISTAKES ---
	fmt.Println()
	fmt.Println("=== Common WaitGroup Mistakes ===")
	fmt.Println("  ❌ wg.Add(1) INSIDE the goroutine → race condition (main may exit first)")
	fmt.Println("  ❌ Passing wg by VALUE (not &wg)  → Done() on copy, main deadlocks")
	fmt.Println("  ❌ Forgetting defer wg.Done()     → counter never reaches 0, deadlock")
	fmt.Println("  ❌ Calling wg.Add() after Wait()  → panic (negative counter)")
	fmt.Println("  ✅ ALWAYS: Add() before go, &wg as pointer, defer Done() first line")
}
