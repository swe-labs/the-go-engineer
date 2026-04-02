// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 1: Language Basics — Application Logger (Exercise)
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Combining constants, iota, custom types, and methods
//   - Implementing the fmt.Stringer interface (preview of Section 05)
//   - Bounds checking for safety
//   - Building a real-world utility from fundamentals
//
// ENGINEERING DEPTH:
//   Notice the `if l < LevelTrace || l > LevelError` bounds check in the String()
//   method. In C/C++, a missing bounds check on an array lookup like `levelNames[l]`
//   allows attackers to read arbitrary server memory. Go defends against this by
//   triggering a runtime Panic if you lookup an index out-of-bounds, terminating the
//   app instantly rather than silently leaking data. We handle it cleanly here before
//   the panic can even occur.
//
// RUN: go run ./01-language-basics/4-application-logger
// ============================================================================

// LogLevel is a custom type based on int.
// This makes it a TYPE-SAFE enum — you can't accidentally pass
// a random int where a LogLevel is expected.
type LogLevel int

// iota assigns auto-incrementing values:
// LevelTrace=0, LevelDebug=1, LevelInfo=2, LevelWarning=3, LevelError=4
const (
	LevelTrace   LogLevel = iota // 0 — Most verbose, for debugging internals
	LevelDebug                   // 1 — Debug info during development
	LevelInfo                    // 2 — Normal operational messages
	LevelWarning                 // 3 — Something unexpected, but not broken
	LevelError                   // 4 — Something failed, needs attention
)

// levelNames maps each LogLevel to its human-readable string.
// The index of each string matches the iota value of the constant.
// This is a common Go pattern: parallel array indexing.
var levelNames = []string{"Trace", "Debug", "Info", "Warning", "Error"}

// String implements the fmt.Stringer interface.
// When you pass a LogLevel to fmt.Println or fmt.Printf with %s,
// Go automatically calls this method to get the string representation.
//
// BOUNDS CHECK: We validate the input to prevent an out-of-range panic.
// If someone creates LogLevel(99), we return "Unknown" instead of crashing.
func (l LogLevel) String() string {
	if l < LevelTrace || l > LevelError {
		return "Unknown"
	}

	return levelNames[l]
}

// printLogLevel demonstrates using the custom type in a function.
// %d prints the integer value, .String() gives the human-readable name.
func printLogLevel(level LogLevel) {
	fmt.Printf("Log level: %d %s\n", level, level.String())
}

func main() {

	// Print all valid log levels
	printLogLevel(LevelTrace)   // Log level: 0 Trace
	printLogLevel(LevelDebug)   // Log level: 1 Debug
	printLogLevel(LevelInfo)    // Log level: 2 Info
	printLogLevel(LevelWarning) // Log level: 3 Warning
	printLogLevel(LevelError)   // Log level: 4 Error

	// Test the bounds check with an invalid level
	printLogLevel(10) // Log level: 10 Unknown — safely handled!

	// KEY TAKEAWAY:
	// This exercise combined:
	//   1. Custom types (LogLevel)
	//   2. iota for auto-incrementing constants
	//   3. fmt.Stringer interface (String() method)
	//   4. Defensive programming (bounds checking)
	// These patterns appear in every production Go codebase.
}
