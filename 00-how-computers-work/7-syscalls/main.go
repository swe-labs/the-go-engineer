// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 00: How Computers Work — Syscalls
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - User-space programs ask the OS to perform privileged work
//   - File operations eventually cross into the kernel through syscalls
//   - OS-backed operations are a boundary between code and reality
//
// WHY THIS MATTERS:
//   Files, sockets, and many I/O surfaces are not “just function calls” — they
//   depend on the operating system to perform protected work.
//
// RUN: go run ./00-how-computers-work/7-syscalls
// ============================================================================

package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.CreateTemp("", "hc7-*.txt")
	if err != nil {
		fmt.Println("create temp file failed:", err)
		return
	}
	defer os.Remove(file.Name())
	defer file.Close()

	payload := []byte("syscalls connect user code to the operating system")
	if _, err := file.Write(payload); err != nil {
		fmt.Println("write failed:", err)
		return
	}

	readBack, err := os.ReadFile(file.Name())
	if err != nil {
		fmt.Println("read failed:", err)
		return
	}

	fmt.Printf("Created file : %s\n", file.Name())
	fmt.Printf("Wrote bytes  : %d\n", len(payload))
	fmt.Printf("Read content : %q\n", string(readBack))
	fmt.Println("The Go helpers are friendly. The OS work underneath is still real.")

	// KEY TAKEAWAY:
	// - Programs reach protected OS resources through syscalls.
	// - File and network work crosses the user/kernel boundary.
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: HC.8 blocking-vs-non-blocking-io")
	fmt.Println("Run    : go run ./00-how-computers-work/8-blocking-vs-non-blocking-io")
	fmt.Println("Current: HC.7 (syscalls)")
	fmt.Println("---------------------------------------------------")
}
