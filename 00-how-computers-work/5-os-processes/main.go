// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work — How the OS Manages Processes
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - A Go program runs inside a normal OS process
//   - The OS assigns process identity like a PID
//   - Signals, scheduling, and file descriptors belong to process reality
//
// WHY THIS MATTERS:
//   Production behavior is shaped by process boundaries, not only by business
//   logic written in the source code.
//
// RUN: go run ./00-how-computers-work/5-os-processes
// ============================================================================

package main

import (
	"fmt"
	"os"
)

func main() {
	pid := os.Getpid()
	fmt.Printf("This Go program is running as process ID %d.\n", pid)
	fmt.Println("The OS tracks this process separately from every other running program.")

	// KEY TAKEAWAY:
	// - A running program becomes an OS process with identity and resources.
	// - Signals and scheduling act on the process, not on abstract source code.
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.6 cpu-cache-and-performance")
	fmt.Println("Run    : go run ./00-how-computers-work/6-cpu-cache-and-performance")
	fmt.Println("Current: HC.5 (os-processes)")
	fmt.Println("---------------------------------------------------")
}
