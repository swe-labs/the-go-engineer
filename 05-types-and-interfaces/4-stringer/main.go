// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 5: Types & Interfaces — The Stringer Interface
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The fmt.Stringer interface: the most commonly implemented interface
//   - How fmt.Println automatically calls String() on your types
//   - Adding custom types to int, string, etc. with type definitions
//   - Why Stringer matters for logging, debugging, and user-facing output
//
// ANALOGY:
//   When you hand someone a business card, the card shows a carefully
//   formatted summary — not raw data. The String() method is your type's
//   "business card". Without it, Go prints the raw struct fields.
//   With it, you control exactly how your type is presented.
//
// THE STRINGER INTERFACE (from package fmt):
//   type Stringer interface {
//       String() string
//   }
//   Any type with a String() string method will be used by fmt.Println,
//   fmt.Printf with %s or %v, log.Println, and error messages.
//
// RUN: go run ./05-types-and-interfaces/4-stringer
// ============================================================================

// --- EXAMPLE 1: Stringer on a struct ---

// HTTPStatus represents an HTTP response status code with its meaning.
type HTTPStatus struct {
	Code    int    // Numeric code (e.g., 200, 404, 500)
	Message string // Human-readable message (e.g., "OK", "Not Found")
}

// String implements fmt.Stringer for HTTPStatus.
// When you pass an HTTPStatus to fmt.Println, Go calls this method.
//
// WITHOUT String(): {200 OK}               ← raw struct dump
// WITH String():    HTTP 200: OK            ← clean, readable
func (s HTTPStatus) String() string {
	return fmt.Sprintf("HTTP %d: %s", s.Code, s.Message)
}

// --- EXAMPLE 2: Stringer on a custom type ---

// Weekday is a NAMED TYPE based on int.
// This creates a completely new type — Weekday and int are NOT interchangeable.
// This is Go's answer to enums (combined with iota from Section 01).
type Weekday int

// Define the weekday constants using iota.
const (
	Sunday    Weekday = iota // 0
	Monday                   // 1
	Tuesday                  // 2
	Wednesday                // 3
	Thursday                 // 4
	Friday                   // 5
	Saturday                 // 6
)

// String implements fmt.Stringer for Weekday.
// Now fmt.Println(Monday) prints "Monday" instead of "1".
func (d Weekday) String() string {
	// A slice indexed by the Weekday value — fast O(1) lookup.
	names := [...]string{
		"Sunday", "Monday", "Tuesday", "Wednesday",
		"Thursday", "Friday", "Saturday",
	}
	// Safety check: prevent out-of-bounds panic
	if d < Sunday || d > Saturday {
		return "Unknown"
	}
	return names[d]
}

// IsWeekend returns true if the day is Saturday or Sunday.
// This shows that custom types can have multiple methods, not just String().
func (d Weekday) IsWeekend() bool {
	return d == Saturday || d == Sunday
}

func main() {
	fmt.Println("=== The Stringer Interface ===")
	fmt.Println()

	// --- Struct with Stringer ---
	fmt.Println("--- HTTP Status codes ---")
	ok := HTTPStatus{Code: 200, Message: "OK"}
	notFound := HTTPStatus{Code: 404, Message: "Not Found"}
	serverErr := HTTPStatus{Code: 500, Message: "Internal Server Error"}

	// fmt.Println automatically calls String() on each value
	fmt.Println(" ", ok)        // HTTP 200: OK
	fmt.Println(" ", notFound)  // HTTP 404: Not Found
	fmt.Println(" ", serverErr) // HTTP 500: Internal Server Error
	fmt.Println()

	// --- Custom type with Stringer ---
	fmt.Println("--- Weekdays (custom type on int) ---")
	today := Wednesday
	fmt.Printf("  Today is: %s\n", today)               // Wednesday (not "3")
	fmt.Printf("  Is weekend: %t\n", today.IsWeekend()) // false

	weekend := Saturday
	fmt.Printf("  %s is weekend: %t\n", weekend, weekend.IsWeekend()) // true
	fmt.Println()

	// --- Iterate all weekdays ---
	fmt.Println("--- All weekdays ---")
	for d := Sunday; d <= Saturday; d++ {
		marker := "  " // Indentation
		if d.IsWeekend() {
			marker = "🏖️" // Mark weekends
		}
		fmt.Printf("  %s %s (%d)\n", marker, d, int(d))
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Implement String() string to control how your type prints")
	fmt.Println("  - fmt.Println, format verbs, and log functions all call String()")
	fmt.Println("  - Custom types (type X int) create new types with their own methods")
	fmt.Println("  - Stringer is Go's most commonly implemented interface")
}
