// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Constants
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Why constants are immutable (cannot change at runtime).
//   - How to group constants for cleaner configuration.
//   - The difference between compile-time constants and runtime variables.
//
// WHY THIS MATTERS:
//   - Encoding stable facts (like ports, mathematical constants, or limits) into
//     constants prevents accidental runtime drift and allows the compiler to
//     optimize the binary by inlining values.
//
// RUN:
//   go run ./02-language-basics/2-constants
//
// KEY TAKEAWAY:
//   - Use constants for stability; use variables for state. Immutability is a
//     powerful tool for reducing system complexity.
// ============================================================================

package main

import "fmt"

// Constants can be grouped. They are resolved by the compiler, not at runtime.
// This is perfect for defining fixed infrastructure ports, connection limits,
// or state codes that never change during the lifecycle of the application.
const (
	Host = "127.0.0.1"
	Port = ":8080"
	User = "root"
)

var (
	isRunning bool
)

func main() {
	AppName := "Go"
	fmt.Println(AppName)

	// Constants can be explicitly typed if needed to prevent implicit conversions.
	const pi float64 = 3.1415926
	fmt.Println(pi)

	const rate float32 = 5.2
	fmt.Println(rate)

	fmt.Printf("Server: %s%s (User: %s)\n", Host, Port, User)
	fmt.Printf("Running: %t\n", isRunning)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: LB.3 -> 02-language-basics/3-enums")
	fmt.Println("Run    : go run ./02-language-basics/3-enums")
	fmt.Println("Current: LB.2 (constants)")
	fmt.Println("---------------------------------------------------")
}
