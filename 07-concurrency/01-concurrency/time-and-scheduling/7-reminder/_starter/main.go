// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Stage 07: Time & Scheduling — Console Reminder (Exercise Starter)
// Level: Intermediate
// ============================================================================
//
// EXERCISE: Build a Console Reminder with Countdown
//
// REQUIREMENTS:
//  1. [ ] Parse seconds and message from os.Args
//  2. [ ] Use time.AfterFunc to schedule the reminder message
//  3. [ ] Use time.NewTicker(1*time.Second) for a countdown display
//  4. [ ] Use select to multiplex between ticker.C and a done channel
//  5. [ ] Always defer ticker.Stop() to prevent leaks
//
// RUN: go run ./07-concurrency/01-concurrency/time-and-scheduling/7-reminder/_starter 5 "Hello!"
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Implement runReminder(seconds int, message string)

func main() {
	fmt.Println("=== Console Reminder Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your reminder app!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
