// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

// ============================================================================
// Stage 07: Time & Scheduling — Time Formatting
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The unique Go standard Time Layout string: "2006-01-02 15:04:05"
//   - Built-in constants like `time.RFC3339`
//   - Parsing strings back into `time.Time` structs
//
// ENGINEERING DEPTH:
//   Most languages use token-based formatting like `YYYY-MM-DD HH:MM:SS`.
//   Go developers realized this requires constantly looking up cheat sheets.
//   Instead, Go uses a single, memorable "Reference Time":
//   Mon Jan 2 15:04:05 MST 2006. (1-2-3-4-5-6-7).
//   You simply write down what that specific reference time would look like
//   in your desired format, and the engine tokenizes it underneath.
//
// RUN: go run ./07-concurrency/01-concurrency/time-and-scheduling/2-formatting
// ============================================================================

import (
	"fmt"
	"time"
)

func main() {

	now := time.Now()

	fmt.Printf("Current time (default Go format): %s\n", now)

	fmt.Printf("Formatted as YYYY-MM-DD: %s\n", now.Format("2006-01-02"))

	fmt.Printf("Formatted as MM/DD/YYYY hh:mm:ss PM: %s\n", now.Format("01/02/2006 03:04:05 PM"))

	fmt.Printf("Formatted as Day, Month Date, Year: %s\n", now.Format("Mon, Jan 2, 2006"))

	fmt.Printf("Formatted as RFC3339: %s\n", now.Format(time.RFC3339))
	fmt.Printf("Formatted as ANSIC: %s\n", now.Format(time.ANSIC))

	dateStr1 := "2025-07-15"
	layout1 := "2006-01-02"

	parsedTime1, err := time.Parse(layout1, dateStr1)
	if err != nil {
		fmt.Printf("Error parsing '%s' with layout '%s': %v\n", dateStr1, layout1, err)
	} else {
		fmt.Printf("Parsed '%s' -> Year: %d, Month: %s, Day: %d\n",
			dateStr1, parsedTime1.Year(), parsedTime1.Month(), parsedTime1.Day())
	}

	rfc3339Str := "2025-12-25T10:00:00+01:00"
	parsedTimeRFC, err := time.Parse(time.RFC3339, rfc3339Str)
	if err != nil {
		fmt.Printf("Error parsing RFC3339 string '%s': %v\n", rfc3339Str, err)
	} else {
		fmt.Printf("Parsed RFC3339: %s (Location: %s, Offset: %d seconds)\n",
			parsedTimeRFC, parsedTimeRFC.Location(), getOffsetInSeconds(parsedTimeRFC))
		fmt.Printf("Parsed RFC3339 (in UTC): %s\n", parsedTimeRFC.UTC())
	}
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: TM.3 timers & tickers")
	fmt.Println("   Current: TM.2 (formatting)")
	fmt.Println("---------------------------------------------------")
}

func getOffsetInSeconds(t time.Time) int {
	_, offset := t.Zone()
	return offset
}
