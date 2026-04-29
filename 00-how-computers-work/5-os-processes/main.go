// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work
// Title: How the OS manages processes
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Understand that your Go program runs inside an operating-system process and crosses a syscall boundary whenever it needs OS help.
//
// WHY THIS MATTERS:
//   - A process is the OS sandbox for a running program. Syscalls are the doors that let that sandbox ask for files, network access, clocks, and other ex...
//
// RUN:
//   go run ./00-how-computers-work/5-os-processes
//
// KEY TAKEAWAY:
//   - Understand that your Go program runs inside an operating-system process and crosses a syscall boundary whenever it needs OS help.
// ============================================================================

package main

import (
	"fmt"
	"os"
	"runtime"
)

//

func main() {
	fmt.Println("=== OS Processes and Syscalls ===")
	fmt.Printf("PID: %d\n", os.Getpid())
	fmt.Printf("Parent PID: %d\n", os.Getppid())
	host, err := os.Hostname()
	if err != nil {
		fmt.Printf("hostname lookup failed: %v\n", err)
	} else {
		fmt.Printf("Hostname: %s\n", host)
	}
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Println()
	fmt.Println("A process owns memory, file descriptors, and execution state.")
	fmt.Println("When it needs files, clocks, or network access, it crosses the syscall boundary to ask the OS for help.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: GT.1 installation")
	fmt.Println("Current: HC.5 (os processes)")
	fmt.Println("---------------------------------------------------")
}
