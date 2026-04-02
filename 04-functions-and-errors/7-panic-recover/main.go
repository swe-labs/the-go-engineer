// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 4: Functions & Errors — Panic & Recover
// Level: Intermediate → Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - panic() immediately stops the current function and unwinds the stack
//   - recover() catches a panic, preventing a crash
//   - recover() ONLY works inside a deferred function
//   - When to use panic vs returning errors (almost never use panic)
//
// IMPORTANT RULE:
//   In production Go code, panic is almost NEVER appropriate.
//   Use error returns (value, error) for expected failure cases.
//   Panic is reserved for truly unrecoverable situations:
//     - Programming bugs (nil dereference, index out of range)
//     - Initialization failures (can't load config, can't connect to DB at startup)
//     - Invariant violations (conditions that "should be impossible")
//
// ENGINEERING DEPTH:
//   When `panic()` is called, the Go runtime halts standard execution and begins
//   "Stack Unwinding". It walks backward up the Call Stack, executing any active
//   deferred functions it finds. If it unwinds all the way back to the root of the
//   Goroutine without encountering a `recover()`, the runtime sends an immediate
//   `SIGABRT` to the OS, violently crashing the entire application and taking down
//   every other active Goroutine with it. NEVER use panic for control flow.
//
// RUN: go run ./04-functions-and-errors/7-panic-recover
// ============================================================================

// mightPanic demonstrates the panic() function.
// When panic is called:
//  1. The current function STOPS immediately (remaining lines are skipped)
//  2. All deferred functions in the current goroutine execute (LIFO)
//  3. The program crashes with a stack trace — UNLESS recover() catches it
func mightPanic(shouldPanic bool) {
	if shouldPanic {
		panic("something went terribly wrong") // Crash! (unless recovered)
	}

	fmt.Println("  This function executed without a panic")
}

// recoverable demonstrates how to catch a panic using defer + recover.
//
// THE RECOVERY PATTERN:
//
//	func safe() {
//	    defer func() {
//	        if r := recover(); r != nil {
//	            // Handle the panic (log it, return an error, etc.)
//	        }
//	    }()
//	    riskyOperation()
//	}
//
// recover() returns:
//   - The value passed to panic() — if a panic was in progress
//   - nil — if no panic occurred
//
// recover() ONLY works when called from a DEFERRED function.
// Calling recover() in normal code does nothing.
func recoverable() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("  ✅ Recovered from panic:", r)
		}
	}()

	fmt.Println("  About to call a function that panics...")
	mightPanic(true) // This will panic
	// This line NEVER executes — panic stops the function immediately:
	fmt.Println("  This line is UNREACHABLE after a panic")
}

func main() {
	fmt.Println("=== Panic & Recover Demo ===")
	fmt.Println()

	// Call the recoverable function — the panic is caught internally
	recoverable()

	// Execution continues here! Because recover() caught the panic,
	// it didn't crash the program.
	fmt.Println()
	fmt.Println("  Program continues after recovered panic ✅")

	// Show the non-panic path
	fmt.Println()
	fmt.Println("=== Normal Execution (No Panic) ===")
	mightPanic(false) // No panic — prints normally

	// KEY TAKEAWAY:
	// - panic = nuclear option. Use errors.New() for expected failures.
	// - recover() only works inside a deferred function.
	// - The defer+recover pattern is used in HTTP servers to prevent
	//   one bad request from crashing the entire server.
	//   (See Section 13: Web Masterclass — Middleware)
	// - If you find yourself writing panic() in application code,
	//   ask: "Can I return an error instead?" The answer is almost always yes.
}
