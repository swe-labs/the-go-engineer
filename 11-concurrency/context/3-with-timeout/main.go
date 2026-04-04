// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"context"
	"fmt"
	"time"
)

// ============================================================================
// Section 17: Context — WithTimeout & WithDeadline
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - context.WithTimeout — auto-cancel after a duration
//   - context.WithDeadline — auto-cancel at a specific time
//   - How to use timeouts for database queries, API calls, etc.
//   - Detecting timeout vs manual cancellation
//   - The production pattern: always set timeouts on I/O operations
//
// ENGINEERING DEPTH:
//   `WithTimeout` is fundamentally built on `WithDeadline`. When you call
//   `WithTimeout(ctx, 5 * time.Second)`, Go computes `time.Now().Add(5s)` and
//   passes it internally to `WithDeadline`. The Go runtime registers a timer in
//   its internal timer heap (a min-heap data structure managed per-P). When the
//   deadline arrives, the runtime's timer goroutine fires the context's internal
//   `cancel()` function, closing the Done() channel and unblocking all waiters.
//
// RUN: go run ./17-context/3-with-timeout
// ============================================================================

func main() {
	fmt.Println("=== Context: WithTimeout ===")
	fmt.Println()

	// --- WithTimeout ---
	// Creates a context that automatically cancels after the specified duration.
	// This is the MOST COMMONLY USED context derivation in production code.
	//
	// Real-world usage:
	//   ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	//   defer cancel()
	//   rows, err := db.QueryContext(ctx, "SELECT ...")
	//   // If the query takes > 5 seconds, it's automatically cancelled

	// Create a context that expires in 200ms
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel() // Always call cancel, even with timeouts

	// Simulate a slow operation
	fmt.Println("  Starting slow operation (timeout: 200ms)...")
	result, err := slowOperation(ctx, 500*time.Millisecond) // Takes 500ms — will timeout!
	if err != nil {
		fmt.Printf("  ❌ Operation failed: %v\n", err)
		// Check WHICH type of error occurred
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("  Reason: Timeout — deadline exceeded")
		} else if ctx.Err() == context.Canceled {
			fmt.Println("  Reason: Manually cancelled")
		}
	} else {
		fmt.Printf("  ✅ Result: %s\n", result)
	}

	fmt.Println()

	// --- Successful operation (finishes before timeout) ---
	ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel2()

	fmt.Println("  Starting fast operation (timeout: 500ms)...")
	result, err = slowOperation(ctx2, 100*time.Millisecond) // Takes 100ms — will succeed!
	if err != nil {
		fmt.Printf("  ❌ Failed: %v\n", err)
	} else {
		fmt.Printf("  ✅ Result: %s\n", result)
	}

	fmt.Println()

	// --- WithDeadline ---
	// WithDeadline is like WithTimeout but takes an absolute time instead of a duration.
	// WithTimeout(ctx, 5*time.Second) is equivalent to:
	// WithDeadline(ctx, time.Now().Add(5*time.Second))
	//
	// Use WithDeadline when you have a fixed deadline (e.g., "must finish by 5:00 PM").
	// Use WithTimeout when you have a duration (e.g., "must finish within 5 seconds").
	deadline := time.Now().Add(150 * time.Millisecond)
	ctx3, cancel3 := context.WithDeadline(context.Background(), deadline)
	defer cancel3()

	fmt.Println("  Starting operation with absolute deadline...")
	dl, ok := ctx3.Deadline()
	fmt.Printf("  Deadline set: %v (has deadline: %v)\n",
		dl.Format("15:04:05.000"), ok)

	result, err = slowOperation(ctx3, 300*time.Millisecond)
	if err != nil {
		fmt.Printf("  ❌ Failed: %v\n", err)
	} else {
		fmt.Printf("  ✅ Result: %s\n", result)
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAYS:")
	fmt.Println("  1. WithTimeout auto-cancels after a duration — use for I/O timeouts")
	fmt.Println("  2. WithDeadline auto-cancels at an absolute time")
	fmt.Println("  3. ALWAYS defer cancel() — even with auto-cancellation")
	fmt.Println("  4. Check ctx.Err() to distinguish DeadlineExceeded vs Canceled")
	fmt.Println("  5. PRODUCTION RULE: Never call a DB or API without a timeout context")
}

// slowOperation simulates an operation that takes the specified duration.
// It respects context cancellation — if the context expires before the
// operation completes, it returns immediately with an error.
//
// This is the STANDARD PATTERN for any cancellation-aware function:
//
//	select {
//	case <-ctx.Done(): return ctx.Err()   // Cancelled!
//	case result := <-work: return result  // Completed!
//	}
func slowOperation(ctx context.Context, duration time.Duration) (string, error) {
	select {
	case <-ctx.Done():
		// Context expired — return the cancellation error
		return "", ctx.Err()
	case <-time.After(duration):
		// Operation completed before timeout
		return "operation completed successfully", nil
	}
}
