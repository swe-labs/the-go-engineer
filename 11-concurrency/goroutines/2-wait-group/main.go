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
// Section 11: Concurrency � WaitGroups
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - sync.WaitGroup: the standard tool for waiting on goroutines
//   - The 3-step pattern: Add ? go func { defer Done } ? Wait
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
// RUN: go run ./11-concurrency/goroutines/2-wait-group
// ============================================================================

type ServiceStatus struct {
	Name    string
	Healthy bool
	Latency time.Duration
}

func checkService(name string, wg *sync.WaitGroup, results chan<- ServiceStatus) {
	defer wg.Done()

	latency := time.Duration(50+rand.Intn(250)) * time.Millisecond
	time.Sleep(latency)
	healthy := rand.Intn(10) > 1

	results <- ServiceStatus{
		Name:    name,
		Healthy: healthy,
		Latency: latency,
	}
}

func main() {
	fmt.Println("=== WaitGroups: Coordinating Goroutine Completion ===")
	fmt.Println()

	var wg sync.WaitGroup
	services := []string{
		"postgres-db",
		"redis-cache",
		"auth-service",
		"email-service",
		"payment-gateway",
		"search-engine",
	}

	results := make(chan ServiceStatus, len(services))

	fmt.Printf("?? Health checking %d services concurrently...\n\n", len(services))

	for _, svc := range services {
		wg.Add(1)
		go checkService(svc, &wg, results)
	}

	wg.Wait()
	close(results)

	allHealthy := true
	for status := range results {
		icon := "?"
		if !status.Healthy {
			icon = "?"
			allHealthy = false
		}
		fmt.Printf("  %s %-20s latency: %v\n", icon, status.Name, status.Latency)
	}

	fmt.Println()
	if allHealthy {
		fmt.Println("?? All services healthy!")
	} else {
		fmt.Println("??  Some services are degraded � check logs!")
	}

	fmt.Println()
	fmt.Println("=== Common WaitGroup Mistakes ===")
	fmt.Println("  ? wg.Add(1) INSIDE the goroutine ? race condition")
	fmt.Println("  ? Passing wg by VALUE (not &wg)  ? Done() on copy, main deadlocks")
	fmt.Println("  ❌ Forgetting defer wg.Done()     → counter never reaches 0, deadlock")
	fmt.Println("  ❌ Calling wg.Add() after Wait()  → panic (negative counter)")
	fmt.Println("  ✅ ALWAYS: Add() before go, &wg as pointer, defer Done() first line")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: GC.3 channels (unbuffered)")
	fmt.Println("   Current: GC.2 (WaitGroups)")
	fmt.Println("---------------------------------------------------")
}
