package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 07: Concurrency - Goroutine leaks
//
// Run: go run ./07-concurrency/01-concurrency/sync-primitives/5-goroutine-leaks

func main() {
	fmt.Println("=== SY.5 Goroutine leaks ===")
	fmt.Println("Learn how goroutines get stranded and why lifetime ownership matters as much as spawning work.")
	fmt.Println()
	fmt.Println("- Every launched goroutine needs a clear owner.")
	fmt.Println("- Contexts and channel closure define lifetime boundaries.")
	fmt.Println("- Leaks often hide behind blocked sends, blocked receives, or forgotten timers.")
	fmt.Println()
	fmt.Println("Leak prevention is mostly design discipline: every goroutine needs a stop condition that another part of the program actually controls.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SY.6")
	fmt.Println("Current: SY.5 (goroutine leaks)")
	fmt.Println("---------------------------------------------------")
}
