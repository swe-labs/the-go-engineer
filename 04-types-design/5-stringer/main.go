// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Stringer
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Implementing the `fmt.Stringer` interface for custom types.
//   - Controlling the visual representation of data in logs and output.
//   - Defining new types based on existing primitives.
//   - The mechanics of how the `fmt` package detects interface satisfaction.
//
// WHY THIS MATTERS:
//   - The `fmt.Stringer` interface is the most ubiquitous interface in the
//     Go ecosystem. Correct implementation ensures that your types are
//     self-documenting in logs, error messages, and debugger views,
//     drastically reducing troubleshooting time.
//
// RUN:
//   go run ./04-types-design/5-stringer
//
// KEY TAKEAWAY:
//   - The `String()` method provides a standardized textual view of a type.
// ============================================================================

// See LICENSE for usage terms.

package main

import "fmt"

// Section 04: Types & Design - Stringer

// HTTPStatus represents a network response status code and message.
// HTTPStatus (Struct): represents a network response status code and message.
type HTTPStatus struct {
	Code    int
	Message string
}

// String implements the fmt.Stringer interface for HTTPStatus.
// HTTPStatus.String (Method): implements the fmt.Stringer interface for HTTPStatus.
func (s HTTPStatus) String() string {
	return fmt.Sprintf("HTTP %d: %s", s.Code, s.Message)
}

// Weekday represents a day of the week as an enumerated integer.
// Weekday (Type): represents a day of the week as an enumerated integer.
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

// String implements the fmt.Stringer interface for Weekday.
// Weekday.String (Method): implements the fmt.Stringer interface for Weekday.
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

// Weekday.IsWeekend (Method): applies the is weekend operation to receiver state at a visible boundary.
func (d Weekday) IsWeekend() bool {
	return d == Saturday || d == Sunday
}

func main() {
	fmt.Println("=== Stringer: Customizing Textual Representation ===")
	fmt.Println()

	// 1. Automatic interface detection.
	// fmt.Println and related functions use reflection to check for the String() method.
	fmt.Println("--- HTTP Status Documentation ---")
	ok := HTTPStatus{Code: 200, Message: "OK"}
	notFound := HTTPStatus{Code: 404, Message: "Not Found"}
	serverErr := HTTPStatus{Code: 500, Message: "Internal Server Error"}

	fmt.Println(" ", ok)
	fmt.Println(" ", notFound)
	fmt.Println(" ", serverErr)
	fmt.Println()

	// 2. Custom types on primitives.
	// Weekday is an int, but because it has its own String() method, it
	// behaves like a rich object when printed.
	fmt.Println("--- Enumerated Types (int-based) ---")
	today := Wednesday
	fmt.Printf("  Today is: %s\n", today)
	fmt.Printf("  Is weekend: %t\n", today.IsWeekend())

	weekend := Saturday
	fmt.Printf("  %s is weekend: %t\n", weekend, weekend.IsWeekend())
	fmt.Println()

	// 3. Iterating and inspecting custom types.
	// The Stringer implementation handles the translation from internal state (int)
	// to external representation (string) during the loop.
	fmt.Println("--- Weekly Schedule ---")
	for d := Sunday; d <= Saturday; d++ {
		fmt.Printf("  Day %d: %s\n", d, d)
	}

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.6 -> 04-types-design/6-type-switch")
	fmt.Println("Run    : go run ./04-types-design/6-type-switch")
	fmt.Println("Current: TI.5 (stringer)")
	fmt.Println("---------------------------------------------------")
}
