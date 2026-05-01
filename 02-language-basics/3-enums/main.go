// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Enums with iota
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how Go models enum-like values with named types and `iota`.
//
// WHY THIS MATTERS:
//   - Go does not have an `enum` keyword. Instead, it combines: - a named type - a `const` block - `iota` for ordered values That gives you fixed related...
//
// RUN:
//   go run ./02-language-basics/3-enums
//
// KEY TAKEAWAY:
//   - Go simulates enums using custom types and the 'iota' keyword inside
//     grouped constant blocks. This provides type-safe, auto-incrementing
//     values without needing a dedicated 'enum' keyword.
// ============================================================================

package main

import "fmt"

// 'iota' is a special constant generator that starts at 0 and increments by 1
// for each subsequent line in a grouped const block.
// Here, we add 1 to the first iota so Sunday starts at 1 instead of 0.
const (
	Sunday = iota + 1
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// Backward reference:
// We create a custom type based on 'int'. This builds on the basic 'int'
// type we learned in: ../1-variables/README.md
type LogLevel int

// Using our custom type 'LogLevel' along with 'iota' enforces type safety.
// Functions can now require a 'LogLevel' instead of a generic 'int'.
const (
	LogError LogLevel = iota
	LogWarn
	LogInfo
	LogDebug
	LogFatal
)

func (l LogLevel) String() string {
	switch l {
	case LogError:
		return "ERROR"
	case LogWarn:
		return "WARN"
	case LogInfo:
		return "INFO"
	case LogDebug:
		return "DEBUG"
	case LogFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func main() {
	fmt.Println("=== Days of the Week (iota + 1) ===")
	fmt.Println("Sunday:   ", Sunday)
	fmt.Println("Monday:   ", Monday)
	fmt.Println("Tuesday:  ", Tuesday)
	fmt.Println("Wednesday:", Wednesday)
	fmt.Println("Thursday: ", Thursday)
	fmt.Println("Friday:   ", Friday)
	fmt.Println("Saturday: ", Saturday)

	fmt.Println()

	fmt.Println("=== Log Levels (type-safe enum) ===")
	fmt.Println("LogError:", LogError)
	fmt.Println("LogWarn: ", LogWarn)
	fmt.Println("LogInfo: ", LogInfo)
	fmt.Printf("LogError as int: %d\n", int(LogError))

	// Forward reference:
	// We will use these exact logging enum techniques to build a real
	// application logger in the next exercise: ../4-application-logger/README.md
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: LB.4 -> 02-language-basics/4-application-logger")
	fmt.Println("Current: LB.3 (enums)")
	fmt.Println("Previous: LB.2 (constants)")
	fmt.Println("---------------------------------------------------")
}
