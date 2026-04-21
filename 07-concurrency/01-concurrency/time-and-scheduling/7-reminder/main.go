// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./07-concurrency/01-concurrency/time-and-scheduling/7-reminder
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// ============================================================================
// Stage 07: Time & Scheduling — Console Reminder (Exercise)
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Using time.NewTicker for repeating interval events
//   - Using time.AfterFunc for one-shot delayed execution
//   - Parsing durations from command-line arguments
//   - Clean resource management with defer ticker.Stop()
//
// ENGINEERING DEPTH:
//   `time.NewTicker` creates a channel-based timer that fires at regular intervals.
//   Internally, the Go runtime inserts the ticker into its per-P timer heap. Unlike
//   `time.Sleep` (which blocks the goroutine), a Ticker is non-blocking — you receive
//   tick events on a channel, allowing you to `select` between the tick and other
//   signals (like cancellation). Always call `ticker.Stop()` when done, or the
//   runtime will keep the timer alive, leaking a small amount of memory.
//
// USAGE: go run ./07-concurrency/01-concurrency/time-and-scheduling/7-reminder <seconds> <message>
// EXAMPLE: go run ./07-concurrency/01-concurrency/time-and-scheduling/7-reminder 5 "Take a break!"
// ============================================================================

func main() {
	fmt.Println("=== Console Reminder ===")
	fmt.Println()

	// Parse command-line arguments
	// os.Args[0] = program name, os.Args[1] = seconds, os.Args[2] = message
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run . <seconds> <message>")
		fmt.Println("Example: go run . 5 \"Take a break!\"")
		fmt.Println()
		fmt.Println("Running demo mode (3 second reminder)...")
		runReminder(3, "Time's up! This is your reminder.")
		return
	}

	seconds, err := strconv.Atoi(os.Args[1])
	if err != nil || seconds <= 0 {
		fmt.Println("Error: first argument must be a positive integer (seconds)")
		os.Exit(1)
	}

	message := os.Args[2]
	runReminder(seconds, message)
}

func runReminder(seconds int, message string) {
	duration := time.Duration(seconds) * time.Second
	fmt.Printf("⏰ Reminder set for %v from now.\n\n", duration)

	// Create a done channel to signal when the reminder fires.
	done := make(chan struct{})

	// time.AfterFunc schedules a function to run after the duration.
	// It runs the function in its own goroutine.
	// The returned Timer can be used to cancel the reminder if needed.
	time.AfterFunc(duration, func() {
		fmt.Printf("\n🔔 REMINDER: %s\n", message)
		close(done)
	})

	// Create a 1-second ticker for the countdown display.
	// The Ticker fires every interval, sending the current time on its channel.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // Always stop tickers to prevent resource leaks!

	elapsed := 0
	for {
		select {
		case <-ticker.C:
			// Ticker fired — update the countdown.
			elapsed++
			remaining := seconds - elapsed
			if remaining > 0 {
				fmt.Printf("  ⏳ %d seconds remaining...\n", remaining)
			}
		case <-done:
			// Reminder fired — exit the loop.
			fmt.Println()
			fmt.Println("KEY TAKEAWAY:")
			fmt.Println("  - time.AfterFunc schedules a one-shot delayed function")
			fmt.Println("  - time.NewTicker fires at regular intervals on a channel")
			fmt.Println("  - Always defer ticker.Stop() to prevent memory leaks")
			fmt.Println("  - Use select to multiplex between ticker and done signals")
			return
		}
	}
}
