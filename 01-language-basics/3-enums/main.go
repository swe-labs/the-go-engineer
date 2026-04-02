// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 1: Language Basics — Enums with iota
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Go doesn't have a native "enum" keyword — it uses iota instead
//   - How iota auto-increments within a const block
//   - Creating type-safe enums with custom types
//   - Why type-safe enums prevent bugs
//
// ENGINEERING DEPTH:
//   Go controversially omitted the `enum` keyword found in Java or C++.
//   Instead it utilizes custom types (`type LogLevel int`) paired with `iota`.
//   This forces the developer to define behavior natively on the type itself
//   using method receivers (like the `String()` method), blending standard
//   constants directly into Go's powerful Interface ecosystem smoothly.
//
// RUN: go run ./01-language-basics/3-enums
// ============================================================================

// --- BASIC IOTA: Auto-incrementing Constants ---
//
// "iota" is a special constant generator. Inside a const block:
//   - The first constant gets iota = 0
//   - Each subsequent constant gets iota + 1
//   - iota resets to 0 at every new const block
//
// Here, we start at iota + 1, so Sunday = 1, Monday = 2, etc.
// Starting at 1 (not 0) is a common pattern because 0 is the zero value
// for int — you might want 0 to mean "unset" or "unknown".
const (
	Sunday    = iota + 1 // 1
	Monday               // 2 (Go repeats the expression: iota + 1)
	Tuesday              // 3
	Wednesday            // 4
	Thursday             // 5
	Friday               // 6
	Saturday             // 7
)

// --- TYPE-SAFE ENUMS: The Production Pattern ---
//
// In production Go code, you almost NEVER use raw int constants.
// Instead, you create a NAMED TYPE and bind constants to it.
//
// Why? Type safety. With raw ints, nothing prevents you from passing
// a random number like 999 where a LogLevel is expected.
// With a named type, the compiler helps catch misuse.
type LogLevel int

const (
	LogError LogLevel = iota // 0 — Most severe
	LogWarn                  // 1
	LogInfo                  // 2
	LogDebug                 // 3
	LogFatal                 // 4 — Program cannot continue
)

// String() makes LogLevel implement the fmt.Stringer interface.
// This means fmt.Println(LogError) will print "ERROR" instead of "0".
// You'll learn more about interfaces in Section 05.
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

	// --- Basic iota demo ---
	fmt.Println("=== Days of the Week (iota + 1) ===")
	fmt.Println("Sunday:   ", Sunday)    // 1
	fmt.Println("Monday:   ", Monday)    // 2
	fmt.Println("Tuesday:  ", Tuesday)   // 3
	fmt.Println("Wednesday:", Wednesday) // 4
	fmt.Println("Thursday: ", Thursday)  // 5
	fmt.Println("Friday:   ", Friday)    // 6
	fmt.Println("Saturday: ", Saturday)  // 7

	fmt.Println()

	// --- Type-safe enum demo ---
	fmt.Println("=== Log Levels (type-safe enum) ===")

	// Because LogLevel has a String() method, fmt.Println automatically
	// calls it, printing "ERROR" instead of the raw integer 0.
	fmt.Println("LogError:", LogError) // Prints: ERROR (not 0)
	fmt.Println("LogWarn: ", LogWarn)  // Prints: WARN  (not 1)
	fmt.Println("LogInfo: ", LogInfo)  // Prints: INFO  (not 2)

	// You can still access the underlying integer value with explicit conversion
	fmt.Printf("LogError as int: %d\n", int(LogError))

	// KEY TAKEAWAY:
	// - iota generates sequential integer constants automatically.
	// - Always use a named type (like LogLevel) for production enums.
	// - Add a String() method so your enums print readable names.
	// - Start at iota+1 when 0 should mean "unset/unknown".
}
