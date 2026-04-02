// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 1: Language Basics — Application Logger (Exercise Starter)
// Level: Beginner
// ============================================================================
//
// EXERCISE: Build an Application Logger with Severity Levels
//
// REQUIREMENTS:
//  1. [ ] Define a custom type `LogLevel` based on `int`
//  2. [ ] Create constants using `iota`: LevelTrace, LevelDebug, LevelInfo,
//         LevelWarning, LevelError
//  3. [ ] Create a `levelNames` slice mapping each level to its string name
//  4. [ ] Implement the `fmt.Stringer` interface (a `String()` method)
//         with bounds checking to prevent out-of-range panics
//  5. [ ] Create a `Log(level LogLevel, message string)` function that prints
//         formatted log messages like: [INFO] Server started
//  6. [ ] Test it in main() by logging messages at different levels
//
// HINTS:
//   - iota starts at 0 and auto-increments
//   - The Stringer interface: type Stringer interface { String() string }
//   - Use the levelNames slice index to convert a LogLevel to a string
//
// RUN: go run ./01-language-basics/4-application-logger
// SOLUTION: See the main.go file in this same directory
// ============================================================================

// TODO: Define your LogLevel type here

// TODO: Define your iota constants here

// TODO: Define your levelNames slice here

// TODO: Implement the String() method on LogLevel

// TODO: Implement the Log function

func main() {
	fmt.Println("=== Application Logger Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your logger and test it here!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with main.go in this directory.")
}
