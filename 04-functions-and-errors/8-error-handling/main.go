// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// Section 4: Functions & Errors — Custom Error Types (Exercise)
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Creating custom error types by implementing the error interface
//   - How to add structured context to errors (operation, inputs, message)
//   - Using defer with multiple return values
//   - Building a complete "safe math" library with error handling
//
// THE ERROR INTERFACE:
//   type error interface {
//       Error() string
//   }
//   Any type with an Error() string method IS an error.
//   This is Go's duck typing in action (see Section 05: Interfaces).
//
// ENGINEERING DEPTH:
//   Because custom errors like `MathError` are essentially just Structs that
//   satisfy an Interface, the performance overhead of throwing a
//   `return &MathError{}` is astronomically lower than throwing an Exception
//   in Java. A Java Exception has to collect the entire dynamic Call Stack into
//   memory to generate a StackTrace. A Go custom error is just instantiating a
//   tiny 24-byte struct on the Heap. It is lightning fast.
//
// RUN: go run ./04-functions-and-errors/8-error-handling
// ============================================================================

// MathError is a CUSTOM ERROR TYPE.
// Unlike errors.New() which gives you a simple string, a custom error type
// carries STRUCTURED DATA — the operation, inputs, and a message.
// This lets the caller inspect the error programmatically instead of
// parsing error strings (which is fragile and an anti-pattern).
type MathError struct {
	Operation string // What operation failed (e.g., "Division")
	InputA    int    // First input
	InputB    int    // Second input
	Message   string // Human-readable description
}

// Constants for operation names — prevents typos and enables comparison.
const (
	division       = "Division"
	divisionErrMsg = "division by zero is not allowed"
)

// Error() implements the error interface.
// This is the ONLY method needed to make MathError a valid error.
// When you do fmt.Println(err), Go calls this method automatically.
func (e MathError) Error() string {
	var inputs []string
	if e.Operation == division {
		inputs = append(inputs, fmt.Sprintf("a=%d", e.InputA))
		inputs = append(inputs, fmt.Sprintf("b=%d", e.InputB))
	}

	return fmt.Sprintf("Math error in %s (%s): %s",
		e.Operation,
		strings.Join(inputs, ", "),
		e.Message)
}

// Sum demonstrates defer with return values.
// The deferred fmt.Println runs AFTER the return value is computed.
func Sum(numbers ...int) int {
	defer fmt.Println("  [defer] Sum computation finished")

	total := 0
	for _, n := range numbers {
		total += n
	}
	return total
}

// SafeDivision returns a custom MathError with full context when division fails.
// Notice: we return *MathError (pointer) instead of MathError (value).
// This is a convention — error values are almost always pointers because:
//  1. They may be nil (pointers can be nil, structs cannot)
//  2. Avoids copying the struct on every return
func SafeDivision(a, b int) (int, error) {
	if b == 0 {
		return 0, &MathError{
			Operation: division,
			InputA:    a,
			InputB:    b,
			Message:   divisionErrMsg,
		}
	}

	return a / b, nil
}

func main() {
	fmt.Println("=== Custom Error Types ===")
	fmt.Println()

	// --- DEFER WITH RETURN VALUES ---
	fmt.Println("Sum with defer:")
	result := Sum(1, 2, 3)
	fmt.Printf("  Result: %d\n\n", result)

	// --- CUSTOM ERROR HANDLING ---
	fmt.Println("SafeDivision:")
	value, err := SafeDivision(10, 0)
	if err != nil {
		// The error message contains structured info from our MathError
		fmt.Println("  Error:", err)

		// TYPE ASSERTION: Extract the concrete MathError from the error interface.
		// This gives us access to the structured fields (Operation, InputA, etc.)
		// You'll learn more about type assertions in Section 05.
		if mathErr, ok := err.(*MathError); ok {
			fmt.Printf("  Operation: %s\n", mathErr.Operation)
			fmt.Printf("  Inputs: a=%d, b=%d\n", mathErr.InputA, mathErr.InputB)
		}
	}

	fmt.Printf("  Value: %d\n\n", value) // 0 (the zero value returned on error)

	// --- SUCCESSFUL CALL ---
	fmt.Println("Successful division:")
	value, err = SafeDivision(10, 3)
	if err != nil {
		fmt.Println("  Error:", err)
	}
	fmt.Printf("  10 / 3 = %d\n", value)

	// KEY TAKEAWAY:
	// - Custom error types carry structured context (operation, inputs, message)
	// - Implement the error interface: just add Error() string to your type
	// - Return *MyError (pointer) so it can be nil on success
	// - Use type assertions to extract error details programmatically
	// - NEVER match error strings — always use structured types or sentinel errors
}
