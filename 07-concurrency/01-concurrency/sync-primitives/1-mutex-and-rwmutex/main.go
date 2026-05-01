// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 07: Concurrency
// Title: sync.Mutex and RWMutex
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to protect shared state using Mutual Exclusion (Mutex).
//   - The difference between Mutex and RWMutex.
//   - When to prioritize read performance with RWMutex.
//
// WHY THIS MATTERS:
//   - Shared state is the source of 90% of concurrency bugs.
//   - Mutexes are the "Seatbelt" of concurrent Go code.
//
// RUN:
//   go run ./07-concurrency/01-concurrency/sync-primitives/1-mutex-and-rwmutex
//
// KEY TAKEAWAY:
//   - Lock early, unlock late (defer), and keep critical sections small.
// ============================================================================

package main

import (
	"fmt"
	"sync"
	"time"
)

// BankAccount uses a Mutex because every operation (Deposit/Withdraw)
// is a Write operation that must be exclusive.
type BankAccount struct {
	mu      sync.Mutex
	balance int
}

func (a *BankAccount) Deposit(amount int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
}

// ConfigCache uses an RWMutex because we expect many readers
// but very few writers (updates to config).
type ConfigCache struct {
	mu     sync.RWMutex
	values map[string]string
}

func (c *ConfigCache) Get(key string) string {
	c.mu.RLock() // Multiple readers can hold this at once!
	defer c.mu.RUnlock()
	return c.values[key]
}

func (c *ConfigCache) Set(key, value string) {
	c.mu.Lock() // Only ONE writer can hold this. Blocks all readers.
	defer c.mu.Unlock()
	c.values[key] = value
}

func main() {
	fmt.Println("=== SY.1 Mutex & RWMutex ===")
	fmt.Println()

	// 1. Mutex Demo (Bank Account)
	account := &BankAccount{balance: 100}
	var wg sync.WaitGroup

	fmt.Println("Scenario 1: Concurrent Deposits (sync.Mutex)")
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			account.Deposit(1)
		}()
	}
	wg.Wait()
	fmt.Printf("  Final Balance: %d (Expected: 1100)\n", account.balance)

	// 2. RWMutex Demo (Config Cache)
	cache := &ConfigCache{values: map[string]string{"theme": "dark"}}

	fmt.Println("\nScenario 2: Read-Heavy Workload (sync.RWMutex)")
	start := time.Now()

	// Launch many readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = cache.Get("theme")
				time.Sleep(1 * time.Millisecond)
			}
		}(i)
	}

	// Launch one writer
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond)
		cache.Set("theme", "light")
		fmt.Println("  [Writer] Theme updated to light")
	}()

	wg.Wait()
	fmt.Printf("  Final Theme: %s (took %v)\n", cache.Get("theme"), time.Since(start))

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: SY.2 -> 07-concurrency/01-concurrency/sync-primitives/2-once-and-sync-map")
	fmt.Println("Current: SY.1 (sync.mutex and rwmutex)")
	fmt.Println("Previous: GC.7 (concurrent-downloader)")
	fmt.Println("---------------------------------------------------")
}
