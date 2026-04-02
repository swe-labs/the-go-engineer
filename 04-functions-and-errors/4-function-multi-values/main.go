// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"errors"
	"fmt"
	"strings"
)

// ============================================================================
// Section 4: Functions & Errors — Multiple Return Values & Named Returns
// Level: Beginner → Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Functions returning multiple values (Go's signature pattern)
//   - The (result, error) convention — how Go handles errors
//   - Named return values (pre-declared in the function signature)
//   - The "naked return" and when to use it (sparingly)
//
// ENGINEERING DEPTH:
//   Java and Python use Exceptions (e.g. `throw new Error()`) to handle failures.
//   Exceptions implicitly hijack the control-flow of the application, jumping up
//   call-stacks in a way that is invisible to the reader. Go designers rejected this.
//   By returning errors as normal second-values (`val, err := func()`), Go forces
//   developers to treat errors as standard data. This makes the control-flow
//   100% visible, reducing complex production bugs drastically.
//
// RUN: go run ./04-functions-and-errors/4-function-multi-values
// ============================================================================

// divide returns TWO values: the result and an error.
//
// THE (RESULT, ERROR) PATTERN:
// This is Go's MOST IMPORTANT convention. Instead of throwing exceptions
// (like Java/Python), Go functions return an error as the last value.
// The caller MUST check the error before using the result.
//
//	result, err := someFunction()
//	if err != nil {
//	    // handle the error
//	}
//	// use result safely
//
// This pattern makes error handling EXPLICIT. You can trace every possible
// failure path by reading the code — no hidden exception throws.
func divide(a, b int) (int, error) {
	if b == 0 {
		// errors.New creates a simple error with a message.
		// Errors in Go are just values — they implement the error interface:
		//   type error interface { Error() string }
		return 0, errors.New("divide by zero")
	}

	return a / b, nil // nil means "no error" — everything went fine
}

// splitName demonstrates NAMED RETURN VALUES.
// The return values (firstName, lastName) are declared in the signature.
// They act as pre-declared variables initialized to their zero values.
//
// NAKED RETURN: A bare "return" (without arguments) returns the current
// values of the named return variables.
//
// CAUTION: Naked returns are controversial. They're convenient for short
// functions but hurt readability in longer functions. Google's style guide
// recommends using them only in functions shorter than ~5 lines.
func splitName(fullName string) (firstName, lastName string) {
	parts := strings.Split(fullName, " ")
	firstName = parts[0] // Assigns to the named return variable
	lastName = parts[1]  // Assigns to the named return variable

	return // Naked return — returns firstName and lastName implicitly
}

func main() {

	// --- MULTIPLE RETURN VALUES: The Error Pattern ---
	fmt.Println("=== Multiple Return Values ===")

	// CORRECT: Always check the error before using the result.
	value, err := divide(10, 0)
	if err != nil {
		fmt.Println("  Error:", err) // "divide by zero"
	} else {
		fmt.Println("  Result:", value)
	}

	// Successful call
	value, err = divide(10, 3)
	if err != nil {
		fmt.Println("  Error:", err)
	} else {
		fmt.Println("  10 / 3 =", value) // 3 (integer division truncates)
	}

	fmt.Println()

	// --- NAMED RETURN VALUES ---
	fmt.Println("=== Named Return Values ===")
	firstName, lastName := splitName("Joseph Abah")
	fmt.Printf("  First: %s, Last: %s\n", firstName, lastName)

	// KEY TAKEAWAY:
	// - Go functions can return multiple values — (result, error) is the standard
	// - ALWAYS check error before using the result
	// - Named returns are useful for documentation but use naked returns sparingly
	// - nil error means success; non-nil error means something went wrong
	// - This pattern replaces try/catch exceptions found in other languages
}
