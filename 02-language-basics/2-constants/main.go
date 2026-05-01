// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 02: Language Basics
// Title: Constants
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how Go represents values that should never change at runtime.
//
// WHY THIS MATTERS:
//   - A variable can change while the program runs. A constant cannot. Constants communicate: - this value is fixed by design - the compiler can treat it...
//
// RUN:
//   go run ./02-language-basics/2-constants
//
// KEY TAKEAWAY:
//   - Constants are evaluated entirely at compile-time. They enforce immutability
//     and communicate fixed design decisions safely without runtime overhead.
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
	// Backward reference:
	// Notice how variables like 'AppName' can be modified later, while
	// the 'Host' and 'Port' constants above cannot.
	// Learn more in the previous lesson: ../1-variables/README.md
	AppName := "Go"
	fmt.Println(AppName)

	// Constants can be explicitly typed if needed to prevent implicit conversions.
	const pi float64 = 3.1415926
	fmt.Println(pi)

	const rate float32 = 5.2
	fmt.Println(rate)

	fmt.Printf("Server: %s%s (User: %s)\n", Host, Port, User)
	fmt.Printf("Running: %t\n", isRunning)

	// Forward reference:
	// Go doesn't have an 'enum' keyword. Next, we will use grouped constants
	// along with the 'iota' keyword to create robust enumerations.
	// See: ../3-enums/README.md
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: LB.3 -> 02-language-basics/3-enums")
	fmt.Println("Current: LB.2 (constants)")
	fmt.Println("Previous: LB.1 (variables)")
	fmt.Println("---------------------------------------------------")
}
