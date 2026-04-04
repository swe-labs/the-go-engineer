// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Section 15: Time & Scheduling — Timezones
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Understanding `time.Location` and IANA Time Zone database
//   - Converting a single unified `time.Time` into different locational offsets
//   - Best practices for system clocks (Always use UTC!)
//
// ENGINEERING DEPTH:
//   A massive architectural mistake is storing locally formatted timestamps in
//   a database (e.g. `2024-10-31 16:00:00 PST`). If the server moves, or DST
//   changes, the data is corrupted. A `time.Time` struct inherently represents
//   an absolute, universal instant in the cosmos. The "Timezone" is merely a
//   presentation layer formatting offset. ALWAYS write down and store absolute
//   UTC time in your databases. Only convert to a local `time.Location` at the
//   very last millisecond before rendering the UI for the user!
//
// RUN: go run ./15-time-and-scheduling/6-timezone
// ============================================================================

import (
	"fmt"
	"time"
)

func main() {

	now := time.Now()
	fmt.Printf("Current local time: %s\n", now)

	newYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Printf("Error loading New York timezone: %v\n", err)
		return
	}

	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Printf("Error loading Tokyo timezone: %v\n", err)
		return
	}

	london, err := time.LoadLocation("Europe/London")
	if err != nil {
		fmt.Printf("Error loading London timezone: %v\n", err)
		return
	}

	fmt.Println("\nSame moment in different timezones:")
	fmt.Printf("New York:  %s\n", now.In(newYork))
	fmt.Printf("Tokyo:     %s\n", now.In(tokyo))
	fmt.Printf("London:    %s\n", now.In(london))
	fmt.Printf("UTC:       %s\n", now.UTC())

	meetingTime := time.Date(2025, 12, 25, 10, 0, 0, 0, newYork)
	fmt.Printf("\nMeeting scheduled for: %s\n", meetingTime)
	fmt.Printf("That's %s in Tokyo\n", meetingTime.In(tokyo))
	fmt.Printf("That's %s in London\n", meetingTime.In(london))
	fmt.Printf("That's %s in UTC\n", meetingTime.UTC())

	timeStr := "2025-06-15T14:30:00-07:00"
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		fmt.Printf("Error parsing time: %v\n", err)
		return
	}
	fmt.Printf("\nParsed time: %s\n", parsedTime)
	fmt.Printf("In New York: %s\n", parsedTime.In(newYork))
	fmt.Printf("In Tokyo: %s\n", parsedTime.In(tokyo))

}
