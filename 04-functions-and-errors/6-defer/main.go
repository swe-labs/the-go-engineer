// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"os"
)

// ============================================================================
// Section 4: Functions & Errors — Defer
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - defer schedules a function call to run AFTER the surrounding function returns
//   - Deferred calls execute in LIFO (Last In, First Out) order — like a stack
//   - defer is essential for cleanup: file.Close(), mutex.Unlock(), rows.Close()
//   - Arguments to deferred functions are evaluated IMMEDIATELY (not at execution)
//
// ENGINEERING DEPTH:
//   Historically, `defer` was slow because it allocated a linked-list node on
//   the Heap for every `defer` call. As of Go 1.14, the Go compiler uses "Open
//   Coded Defers". The compiler literally inspects your function, finds all the
//   `defer` statements, and injects the cleanup code directly into every single
//   exit path (`return`, `panic`) of the compiled Machine Code. This makes
//   `defer` virtually zero-cost at runtime.
//
// RUN: go run ./04-functions-and-errors/6-defer
// ============================================================================

// simpleDefer shows that a deferred function runs AFTER the function body.
// Even though "defer" is called early, its execution is delayed to the end.
func simpleDefer() {
	fmt.Println("  Start")
	defer fmt.Println("  Deferred (runs LAST)") // Scheduled for later
	fmt.Println("  Middle 1")
	fmt.Println("  Middle 2")
	// Output order: Start → Middle 1 → Middle 2 → Deferred
}

// lifoSimpleDefer shows LIFO order: multiple defers execute in reverse order.
// Think of it as a stack: the LAST defer pushed is the FIRST to execute.
//
// This is critical when cleanup must happen in reverse order:
//
//	defer mutex.Lock()     ← runs second
//	defer mutex.Unlock()   ← runs first (reverse of acquisition)
func lifoSimpleDefer() {
	fmt.Println("\n  Start")
	defer fmt.Println("  First defer  (runs SECOND — LIFO)")
	defer fmt.Println("  Second defer (runs FIRST — LIFO)")
	fmt.Println("  Middle")
	// Output: Start → Middle → Second defer → First defer
}

func main() {
	fmt.Println("=== Simple Defer ===")
	simpleDefer()

	fmt.Println("\n=== LIFO Order ===")
	lifoSimpleDefer()

	// --- THE MOST COMMON USE CASE: Resource Cleanup ---
	// Open a resource, then IMMEDIATELY defer its cleanup.
	// This guarantees cleanup even if the function panics or returns early.
	//
	// Pattern:
	//   resource, err := open()
	//   if err != nil { return err }
	//   defer resource.Close()    ← ALWAYS right after successful open
	//   // ... use resource ...
	//
	// This pattern is used for:
	//   - Files:       defer file.Close()
	//   - DB rows:     defer rows.Close()
	//   - Mutexes:     defer mutex.Unlock()
	//   - HTTP bodies: defer resp.Body.Close()
	//   - Temp files:  defer os.Remove(tmpFile)

	fmt.Println("\n=== File Cleanup with Defer ===")
	file, err := os.CreateTemp("", "defer-demo-*.txt")
	if err != nil {
		fmt.Println("  Error:", err)
		return
	}
	defer os.Remove(file.Name()) // Clean up: delete the temp file when main() exits
	defer file.Close()           // Clean up: close the file when main() exits

	// Write to the file — we know Close() will happen no matter what
	fmt.Fprintf(file, "Hello from defer demo!")
	fmt.Printf("  Wrote to temp file: %s\n", file.Name())
	fmt.Println("  (file will be closed and deleted when main() exits)")

	// --- ARGUMENT EVALUATION TIMING ---
	// Arguments to deferred function calls are evaluated IMMEDIATELY,
	// not when the deferred function actually runs.
	fmt.Println("\n=== Argument Evaluation ===")
	x := 10
	defer fmt.Printf("  Deferred with x=%d (captured at defer time, not at execution)\n", x)
	x = 99 // This change does NOT affect the deferred call (x was already captured as 10)
	fmt.Printf("  x is now %d\n", x)

	// KEY TAKEAWAY:
	// - defer delays execution until the surrounding function returns
	// - Multiple defers execute in LIFO (reverse) order
	// - ALWAYS defer cleanup immediately after opening a resource
	// - Arguments are evaluated when defer is called, not when it executes
	// - This is the #1 pattern for preventing resource leaks in Go
}
