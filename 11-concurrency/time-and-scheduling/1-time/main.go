// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Section 15: Time & Scheduling — Time Basics
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The `time.Time` struct and how it stores moments in time
//   - Durations (`time.Duration`) representing intervals
//   - Parsing and formatting standard dates
//   - Time mathematics (Add, Sub, Before, After)
//
// ENGINEERING DEPTH:
//   Under the hood, a `time.Time` struct does not just store a single timestamp.
//   It actually contains BOTH a "Wall Clock" (subject to Timezone shifts, DST,
//   and NTP adjustments) AND a "Monotonic Clock" (which counts hardware ticks
//   since boot and NEVER goes backwards). When you do `t2.Sub(t1)`, Go
//   intelligently uses the Monotonic clock to guarantee the Duration is perfectly
//   accurate, even if the OS suddenly synced its Wall Clock backward by an hour!
//
// RUN: go run ./15-time-and-scheduling/1-time
// ============================================================================

import (
	"fmt"
	"time"
)

func main() {

	now := time.Now() // The current time
	fmt.Printf("Current time: %s\n", now)
	fmt.Printf("Year: %d, Month: %s, Day: %d\n", now.Year(), now.Month(), now.Day())
	fmt.Printf("Hour: %d, Minute: %d, Second: %d\n", now.Hour(), now.Minute(), now.Second())
	fmt.Printf("Weekday: %s\n", now.Weekday().String())

	goLaunchDate := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	fmt.Printf("\nGo's approximate launch date (UTC): %s\n", goLaunchDate)

	oneSecond := time.Second
	fiveMinutes := 5 * time.Minute
	oneHourThirtyMinutes := (1 * time.Hour) + (30 * time.Minute)
	fmt.Printf("\nOne second: %v\n", oneSecond)
	fmt.Printf("Five minutes in nanoseconds: %d ns\n", fiveMinutes)
	fmt.Printf("Five minutes as string: %s\n", fiveMinutes.String())
	fmt.Printf("1h30m as string: %s\n", oneHourThirtyMinutes.String())

	futureTime := now.Add(2 * time.Hour)
	fmt.Printf("Current time + 2 hours: %s\n", futureTime)

	pastTime := now.Add(-30 * time.Minute)
	fmt.Printf("Current time + 30 minutes: %s\n", pastTime)

	durationSinceGoLaunch := now.Sub(goLaunchDate)

	fmt.Printf("Time since Go launch (approx): %s\n", durationSinceGoLaunch)
	fmt.Printf("Approx hours since Go launch: %.0f hours\n", durationSinceGoLaunch.Hours())

	if futureTime.After(now) {
		fmt.Println("futureTime is indeed after now.")
	}

	if pastTime.Before(now) {
		fmt.Println("pastTime is indeed before now.")
	}

	fmt.Println("\nSleeping for 100 milliseconds...")
	shortDuration := 100 * time.Millisecond
	time.Sleep(shortDuration)

	fmt.Println("Awake after sleep!")

}
