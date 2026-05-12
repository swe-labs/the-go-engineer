// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work
// Title: How the OS manages processes
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - A Go program is a process managed by the operating system.
//   - Syscalls are the boundary between your code and the OS kernel.
//
// WHY THIS MATTERS:
//   - A process is a protected sandbox. Every time your program wants to touch
//     the "outside world" (files, network, screen), it must cross the syscall
//     boundary. This boundary has costs and security implications.
//
// RUN:
//   go run ./00-how-computers-work/5-os-processes
//
// KEY TAKEAWAY:
//   - Understanding the process model is the key to mastering signals, performance,
//     and distributed systems.
// ============================================================================

package main

import (
	"fmt"
	"os"
	"runtime"
)

// main (Function): entry point for the program. It queries the OS for process
// metadata (PID, parent PID, hostname) and prints them to demonstrate that a Go
// program runs as an OS-managed process with syscall boundaries.
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
	fmt.Println("A process owns private virtual memory, file descriptors, and execution state.")
	fmt.Println("To touch hardware, the CPU must switch from User Mode to Kernel Mode via a syscall.")
	fmt.Println("This boundary protects the system from buggy or malicious application code.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: GT.1 -> 01-getting-started/1-installation")
	fmt.Println("Run    : go run ./01-getting-started/1-installation")
	fmt.Println("Current: HC.5 (os-processes)")
	fmt.Println("---------------------------------------------------")
}
