// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Section 07: Concurrency - Sync Primitives
// Level: Core
// ============================================================================
//
// RUN: go run ./07-concurrency/01-concurrency/sync-primitives/2-once-and-sync-map
// ============================================================================

import (
	"fmt"
	"sync"
)

// ============================================================================
// Section 07: sync.Once and sync.Map
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - sync.Once for one-time initialization (singletons)
//   - sync.Map for concurrent map access without explicit locking
//   - When to use sync.Map vs regular map + sync.RWMutex
//
// WHY THIS MATTERS:
//   Production services often need safe one-time initialization and shared read-heavy state.
// ============================================================================

// --- sync.Once ---
// sync.Once ensures a function is executed exactly once, even when
// called from multiple goroutines simultaneously.
// Use cases: database connection init, config loading, logger setup.

// DBConnection (Struct): groups the state used by the db connection example boundary.
type DBConnection struct {
	Host string
	Port int
}

var (
	dbInstance *DBConnection
	dbOnce     sync.Once
)

// GetDB returns a singleton database connection.
// No matter how many goroutines call this, the connection is created exactly once.
// GetDB (Function): returns a singleton database connection.
func GetDB() *DBConnection {
	dbOnce.Do(func() {
		fmt.Println("  Initializing database connection (only once!)")
		dbInstance = &DBConnection{
			Host: "localhost",
			Port: 5432,
		}
	})
	return dbInstance
}

// --- sync.Map ---
// sync.Map is a concurrent map that doesn't require explicit locking.
//
// When to use sync.Map:
//   1. When keys are stable (read-heavy, few writes)
//   2. When different goroutines read/write disjoint key sets
//
// When NOT to use sync.Map:
//   1. When you need to iterate often (Range is O(n) and holds no lock)
//   2. When operations on multiple keys must be atomic
//   3. When key set changes frequently - use regular map + RWMutex instead

func main() {
	fmt.Println("=== sync.Once Demo ===")

	// Call GetDB from multiple goroutines - init happens exactly once
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			db := GetDB()
			fmt.Printf("  Goroutine %d: got DB at %s:%d\n", id, db.Host, db.Port)
		}(i)
	}
	wg.Wait()

	fmt.Println("\n=== sync.Map Demo ===")

	var m sync.Map

	// Store values from multiple goroutines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id)
			m.Store(key, id*10) // Thread-safe store
			fmt.Printf("  Stored: %s = %d\n", key, id*10)
		}(i)
	}
	wg.Wait()

	// Load values
	fmt.Println("\nLoading values:")
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key-%d", i)
		if val, ok := m.Load(key); ok {
			fmt.Printf("  %s = %v\n", key, val)
		}
	}

	// LoadOrStore: atomic "get or set" operation
	actual, loaded := m.LoadOrStore("key-0", 999)
	fmt.Printf("\nLoadOrStore key-0: value=%v, loaded=%v (existing value kept)\n", actual, loaded)

	// Range: iterate over all entries
	fmt.Println("\nAll entries:")
	m.Range(func(key, value any) bool {
		fmt.Printf("  %v -> %v\n", key, value)
		return true // return false to stop iteration
	})

	// CompareAndDelete (Go 1.20+): atomic conditional delete
	fmt.Println("\n=== Regular Map + RWMutex (Preferred for Most Cases) ===")

	type SafeMap struct {
		mu sync.RWMutex
		m  map[string]int
	}

	sm := &SafeMap{m: make(map[string]int)}

	// Write
	sm.mu.Lock()
	sm.m["counter"] = 42
	sm.mu.Unlock()

	// Read (multiple concurrent readers allowed)
	sm.mu.RLock()
	fmt.Printf("  counter = %d\n", sm.m["counter"])
	sm.mu.RUnlock()

	fmt.Println("\n  ⚡ Use sync.Map for read-heavy, stable key sets")
	fmt.Println("  ⚡ Use map + RWMutex for everything else")
	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - sync.Once protects one-time initialization across goroutines.")
	fmt.Println("  - sync.Map is specialized for read-heavy, stable-key concurrent access.")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: SY.3 -> 07-concurrency/01-concurrency/sync-primitives/3-atomic-operations")
	fmt.Println("   Current: SY.2 (sync.once and sync.map)")
	fmt.Println("---------------------------------------------------")
}
