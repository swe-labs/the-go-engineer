// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Application Logger
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Build a small logger that combines variables, constants, `iota`, and methods into one readable program.
//
// WHY THIS MATTERS:
//   - This exercise turns separate language pieces into one compact system: - a named type models the log level - `iota` creates ordered constants - a me...
//
// RUN:
//   go run ./02-language-basics/4-application-logger
//
// KEY TAKEAWAY:
//   - Combining custom types, iota constants, and package-level variables
//     allows us to build a robust, type-safe logging foundation that maps
//     numeric state to readable strings.
// ============================================================================

package main

import "fmt"

type LogLevel int

// Backward reference:
// We use 'iota' just like we did in the Enums lesson: ../3-enums/README.md
const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
)

// A package-level slice of strings. Slices are covered later, but here
// it acts as a simple lookup table to map our integer enum to a string.
var levelNames = []string{"Trace", "Debug", "Info", "Warning", "Error"}

// String() is a special method signature. When you print a LogLevel, Go will
// automatically call this method to get the string representation.
func (l LogLevel) String() string {
	// Simple bounds checking to prevent index-out-of-range errors on our slice.
	if l < LevelTrace || l > LevelError {
		return "Unknown"
	}
	return levelNames[l]
}

func printLogLevel(level LogLevel) {
	fmt.Printf("Log level: %d %s\n", level, level.String())
}

func main() {
	printLogLevel(LevelTrace)
	printLogLevel(LevelDebug)
	printLogLevel(LevelInfo)
	printLogLevel(LevelWarning)
	printLogLevel(LevelError)

	printLogLevel(10)

	// Forward reference:
	// Now that you understand the basic data types, we will move on to Control Flow
	// to see how 'if', 'for', and 'switch' direct the program's path.
	// See: ../03-control-flow/1-if-else/README.md
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CF.1 -> 02-language-basics/03-control-flow/1-if-else")
	fmt.Println("Current: LB.4 (application-logger)")
	fmt.Println("Previous: LB.3 (enums)")
	fmt.Println("---------------------------------------------------")
}
