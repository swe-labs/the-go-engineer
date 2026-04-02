// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"errors"
	"fmt"
	"time"
)

// ============================================================================
// Section 4: Functions & Errors — Custom Errors & Sentinel Errors
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Declaring package-level Sentinel Errors (e.g. `ErrNotFound`)
//   - Defining custom Error structs to hold rich metadata (timestamps, codes)
//   - Implementing the `error` interface on custom structs
//   - Using `errors.Is` to check for specific error identities
//
// ENGINEERING DEPTH:
//   Go's `error` type is not a primitive—it is a built-in Interface:
//   `type error interface { Error() string }`. Any struct that implements
//   an `Error() string` method automatically becomes a valid error that can
//   be returned. This "Duck Typing" means errors in Go aren't just strings;
//   they are data structures. You can attach HTTP status codes, SQL query logs,
//   and stack traces directly to your error structs and pass them up the stack!
//
// RUN: go run ./04-functions-and-errors/5-custom-error
// ============================================================================

// --- SENTINEL ERRORS ---
// A sentinel error is an expected, package-level error variable.
// By convention, they always start with "Err". They act like enums for errors.
var ErrDivisionByZero = errors.New("division by zero")
var ErrNumTooLarge = errors.New("number too large")

// --- CUSTOM ERROR STRUCT ---
// OpError holds rich operational metadata about a failure.
// Instead of just passing "failed to connect", we pass the exact operation,
// an internal error code, and the exact timestamp of the failure.
type OpError struct {
	Op      string
	Code    int
	Message string
	Time    time.Time
}

// Error implements the `error` interface for OpError.
// Notice we use a value receiver `(op OpError)`, but in production,
// pointer receivers `(op *OpError)` are usually preferred to allow `nil` returns.
func (op OpError) Error() string {
	return fmt.Sprintf("[%s] Error %d at %s: %s", op.Op, op.Code, op.Time.Format(time.RFC3339), op.Message)
}

// NewOpError is a constructor function for our custom error.
// It returns a pointer *OpError so we can return nil on success paths.
func NewOpError(op string, code int, message string, t time.Time) *OpError {
	return &OpError{
		Op:      op,
		Code:    code,
		Message: message,
		Time:    t,
	}
}

// DoSomething simulates a complex operation failing and returning our custom error.
// The signature returns the interface `error`, but we return our concrete struct!
func DoSomething() error {
	return NewOpError("database_connect", 503, "connection timeout", time.Now())
}

// divide simulates an operation that guards against specific mathematical rules,
// returning our Sentinel Errors.
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}

	if a > 1000 {
		return 0, ErrNumTooLarge
	}

	return a / b, nil
}

func main() {

	fmt.Println("=== Checking Sentinel Errors ===")
	// Attempt a forbidden division
	value, err := divide(1001, 1)

	if err != nil {
		// --- errors.Is ---
		// We use errors.Is() to compare the returned error against our Sentinels.
		// Never use `err.Error() == "number too large"`. String matching is fragile!
		if errors.Is(err, ErrDivisionByZero) {
			fmt.Println("❌ Caught a Division By Zero Exception")
		} else if errors.Is(err, ErrNumTooLarge) {
			fmt.Println("❌ Caught a Number Too Large Exception")
		}
	} else {
		fmt.Println("✅ Success: ", value)
	}

	fmt.Println("\n=== Rich Custom Errors ===")
	// Simulate an operational failure
	opErr := DoSomething()
	if opErr != nil {
		// Because OpError implements the error interface, fmt.Println automatically
		// calls opErr.Error()!
		fmt.Printf("❌ %v\n", opErr)
	}

	// KEY TAKEAWAY:
	// - Sentinel Errors (`ErrSomething`) are used for simple, expected state checks.
	// - Custom Error Structs (`OpError`) are used when you need to attach metadata (like HTTP codes).
	// - ALWAYS use `errors.Is(err, ErrVar)` instead of comparing strings.
}
