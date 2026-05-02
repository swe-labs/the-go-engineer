// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Application Logger
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Combining named types, constants, and iota into a system.
//   - Mapping numeric states to human-readable strings via methods.
//   - Basic bounds checking for lookup safety.
//
// WHY THIS MATTERS:
//   - Production logging is more than just printing text; it requires a stable
//     foundation of levels (Trace, Debug, Info, etc.) that can be parsed by
//     machines and read by humans.
//
// RUN:
//   go run ./02-language-basics/4-application-logger
//
// KEY TAKEAWAY:
//   - Meaningful code composition starts with choosing the right types and
//     providing clear ways to inspect the system's state.
// ============================================================================

package main

import "fmt"

// LogLevel (Type): names the log level concept so the lesson can pass it as a first-class value.
type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
)

// A package-level slice of strings. Slices are covered later, but here
// it acts as a simple lookup table to map our integer enum to a string.
// levelNames (Slice): keeps ordered lesson state visible for iteration and comparison.
var levelNames = []string{"Trace", "Debug", "Info", "Warning", "Error"}

// String() is a special method signature. When you print a LogLevel, Go will
// automatically call this method to get the string representation.
// LogLevel.String (Method): () is a special method signature. When you print a LogLevel, Go will.
func (l LogLevel) String() string {
	// Simple bounds checking to prevent index-out-of-range errors on our slice.
	if l < LevelTrace || l > LevelError {
		return "Unknown"
	}
	return levelNames[l]
}

// printLogLevel (Function): runs the print log level step and keeps its inputs, outputs, or errors visible.
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

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CF.1 -> 02-language-basics/03-control-flow/1-if-else")
	fmt.Println("Run    : go run ./02-language-basics/03-control-flow/1-if-else")
	fmt.Println("Current: LB.4 (application-logger)")
	fmt.Println("---------------------------------------------------")
}
