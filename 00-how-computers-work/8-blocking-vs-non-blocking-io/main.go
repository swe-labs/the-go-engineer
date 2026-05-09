// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work — Blocking vs Non-Blocking I/O
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Blocking I/O makes one execution path wait
//   - Concurrency can overlap waiting with other useful work
//   - I/O-bound programs are often slow because they wait, not because they compute
//
// WHY THIS MATTERS:
//   Real services spend a lot of time waiting on databases, networks, and files,
//   so understanding waiting is as important as understanding computation.
//
// RUN: go run ./00-how-computers-work/8-blocking-vs-non-blocking-io
// ============================================================================

package main

import (
	"fmt"
	"time"
)

func slowIO() {
	time.Sleep(200 * time.Millisecond)
}

func main() {
	start := time.Now()
	slowIO()
	time.Sleep(50 * time.Millisecond)
	blockingDuration := time.Since(start)

	done := make(chan struct{})
	start = time.Now()
	go func() {
		slowIO()
		close(done)
	}()
	time.Sleep(50 * time.Millisecond)
	<-done
	concurrentDuration := time.Since(start)

	fmt.Printf("Blocking path   : %s\n", blockingDuration)
	fmt.Printf("Concurrent path : %s\n", concurrentDuration)
	fmt.Println("Both paths wait on the same slow I/O. Only one overlaps the waiting.")

	// KEY TAKEAWAY:
	// - I/O often dominates runtime because waiting is expensive.
	// - Concurrency helps useful work continue while another task waits.
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: GT.1 installation")
	fmt.Println("Run    : go run ./01-getting-started/1-installation")
	fmt.Println("Current: HC.8 (blocking-vs-non-blocking-io)")
	fmt.Println("---------------------------------------------------")
}
