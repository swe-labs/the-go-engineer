// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work — CPU Cache and Performance
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Sequential memory access is usually friendlier to CPU caches
//   - Random access creates more cache misses and waiting
//   - Performance often depends on data access patterns, not only operation count
//
// WHY THIS MATTERS:
//   Real-world performance work often starts with memory behavior before it
//   starts with clever arithmetic tricks.
//
// RUN: go run ./00-how-computers-work/6-cpu-cache-and-performance
// ============================================================================

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	const size = 1 << 18

	data := make([]int, size)
	indices := make([]int, size)
	for i := range data {
		data[i] = i
		indices[i] = i
	}

	rng := rand.New(rand.NewSource(42))
	rng.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	start := time.Now()
	sumSequential := 0
	for _, value := range data {
		sumSequential += value
	}
	sequentialDuration := time.Since(start)

	start = time.Now()
	sumRandom := 0
	for _, index := range indices {
		sumRandom += data[index]
	}
	randomDuration := time.Since(start)

	fmt.Printf("Sequential access: sum=%d took %s\n", sumSequential, sequentialDuration)
	fmt.Printf("Shuffled access  : sum=%d took %s\n", sumRandom, randomDuration)
	fmt.Println("The math is the same. The access pattern is not.")

	// KEY TAKEAWAY:
	// - Cache-friendly access patterns reduce waiting on memory.
	// - Performance depends on where data lives and how code touches it.
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.7 syscalls")
	fmt.Println("Run    : go run ./00-how-computers-work/7-syscalls")
	fmt.Println("Current: HC.6 (cpu-cache-and-performance)")
	fmt.Println("---------------------------------------------------")
}
