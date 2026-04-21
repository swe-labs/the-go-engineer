// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import "fmt"

// ============================================================================
// Section 6: Types & Interfaces — The Stringer Interface
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The fmt.Stringer interface: the most commonly implemented interface
//   - How fmt.Println automatically calls String() on your types
//   - Adding custom types to int, string, etc. with type definitions
//   - Why Stringer matters for logging, debugging, and user-facing output
//
// RUN: go run ./04-types-design/5-stringer
// ============================================================================

type HTTPStatus struct {
	Code    int
	Message string
}

func (s HTTPStatus) String() string {
	return fmt.Sprintf("HTTP %d: %s", s.Code, s.Message)
}

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func (d Weekday) String() string {
	names := [...]string{
		"Sunday", "Monday", "Tuesday", "Wednesday",
		"Thursday", "Friday", "Saturday",
	}
	if d < Sunday || d > Saturday {
		return "Unknown"
	}
	return names[d]
}

func (d Weekday) IsWeekend() bool {
	return d == Saturday || d == Sunday
}

func main() {
	fmt.Println("=== The Stringer Interface ===")
	fmt.Println()

	fmt.Println("--- HTTP Status codes ---")
	ok := HTTPStatus{Code: 200, Message: "OK"}
	notFound := HTTPStatus{Code: 404, Message: "Not Found"}
	serverErr := HTTPStatus{Code: 500, Message: "Internal Server Error"}

	fmt.Println(" ", ok)
	fmt.Println(" ", notFound)
	fmt.Println(" ", serverErr)
	fmt.Println()

	fmt.Println("--- Weekdays (custom type on int) ---")
	today := Wednesday
	fmt.Printf("  Today is: %s\n", today)
	fmt.Printf("  Is weekend: %t\n", today.IsWeekend())

	weekend := Saturday
	fmt.Printf("  %s is weekend: %t\n", weekend, weekend.IsWeekend())
	fmt.Println()

	fmt.Println("--- All weekdays ---")
	for d := Sunday; d <= Saturday; d++ {
		marker := "  "
		if d.IsWeekend() {
			marker = "  "
		}
		fmt.Printf("  %s %s (%d)\n", marker, d, int(d))
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Implement String() string to control how your type prints")
	fmt.Println("  - fmt.Println, format verbs, and log functions all call String()")
	fmt.Println("  - Custom types (type X int) create new types with their own methods")
	fmt.Println("  - Stringer is Go's most commonly implemented interface")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.6 type switch")
	fmt.Println("   Current: TI.5 (Stringer)")
	fmt.Println("---------------------------------------------------")
}
