// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 1: Language Basics — Constants
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to declare constants with `const`
//   - The difference between constants and variables
//   - Grouped constant declarations
//   - Package-level vs function-level constants
//   - Why constants matter for configuration and safety
//
// ENGINEERING DEPTH:
//   Unlike variables which occupy memory addresses at runtime, constants are
//   structurally different. During the compilation process, the Go compiler literally
//   takes the value of the constant and "inlines" it directly into the Machine Code
//   instructions. They do not exist at runtime, which makes them inherently faster
//   and mathematically impossible to mutate by a rogue pointer!
//
// RUN: go run ./01-language-basics/2-constants
// ============================================================================

// --- PACKAGE-LEVEL CONSTANTS ---
// Constants declared outside of functions are accessible to the entire package.
// By convention, configuration values (hosts, ports, default settings) are
// declared as package-level constants.
//
// Grouped const declarations use parentheses for clean organization.
// These values are set at COMPILE TIME and can never change at runtime.
const (
	Host = "127.0.0.1" // Server hostname
	Port = ":8080"     // Server port (note: Go convention includes the colon)
	User = "root"      // Default user
)

// Package-level variables are declared with "var" (not :=).
// They are initialized to their zero value if no value is provided.
var (
	isRunning bool // zero value: false
)

func main() {

	// --- FUNCTION-LEVEL VARIABLES ---
	// Inside functions, you can use := for short declarations.
	// Note: This is a VARIABLE, not a constant. It CAN be reassigned.
	AppName := "Go"
	fmt.Println(AppName)

	// --- FUNCTION-LEVEL CONSTANTS ---
	// Constants inside functions follow the same rules as package-level ones.
	// They cannot be changed after declaration.
	//
	// TYPED CONSTANTS: When you specify a type (like float64), the constant
	// is locked to that type and cannot be used where a different type is expected.
	const pi float64 = 3.1415926
	fmt.Println(pi)

	const rate float32 = 5.2
	fmt.Println(rate)

	// --- WHY CONSTANTS MATTER ---
	//
	// 1. SAFETY: Constants can't be accidentally modified at runtime.
	//    Imagine if someone changed "Host" mid-execution — chaos.
	//
	// 2. PERFORMANCE: The compiler inlines constants directly into machine code.
	//    No memory allocation, no runtime lookup.
	//
	// 3. DOCUMENTATION: Constants tell readers "this value is fixed by design."
	//
	// RULE: If a value should never change, make it a constant.
	//       If it might change, make it a variable.

	// Print the package-level constants to show they're accessible
	fmt.Printf("Server: %s%s (User: %s)\n", Host, Port, User)
	fmt.Printf("Running: %t\n", isRunning)

	// KEY TAKEAWAY:
	// - Constants are immutable values set at compile time.
	// - Use `const` for configuration, mathematical constants, and fixed values.
	// - Constants can be typed (float64) or untyped (Go infers from usage).
	// - Group related constants with const ( ... ) blocks.
}
