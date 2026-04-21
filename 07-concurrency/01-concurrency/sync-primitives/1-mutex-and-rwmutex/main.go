package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 07: Concurrency - sync.Mutex and RWMutex
//
// Run: go run ./07-concurrency/01-concurrency/sync-primitives/1-mutex-and-rwmutex

func main() {
	fmt.Println("=== SY.1 sync.Mutex and RWMutex ===")
	fmt.Println("Learn how mutual exclusion protects shared state and when read/write locks are worth the extra rules.")
	fmt.Println()
	fmt.Println("- Mutex protects read-modify-write state changes.")
	fmt.Println("- RWMutex only helps when many readers and few writers exist.")
	fmt.Println("- Always know which state a lock owns before using it.")
	fmt.Println()
	fmt.Println("Use locks to protect ownership, not to hide unclear data flow. The safest critical sections stay small and obvious.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SY.2")
	fmt.Println("Current: SY.1 (sync.mutex and rwmutex)")
	fmt.Println("---------------------------------------------------")
}
