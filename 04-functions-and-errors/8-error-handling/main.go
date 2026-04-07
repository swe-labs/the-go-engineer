// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"errors"
	"fmt"
	"math"
)

// ============================================================================
// Section 4: Functions & Errors - Custom Error Types (Exercise)
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Creating custom error types that carry structured context
//   - Returning explicit failures from small math operations
//   - Wrapping low-level failures with extra context using %w
//   - Inspecting errors safely with errors.As after wrapping
//   - Using defer for small cleanup/completion reporting without hiding flow
//
// ENGINEERING DEPTH:
//   In production Go, ordinary failures should stay in the return path.
//   A custom error value is cheap to construct, easy to inspect, and much
//   easier to test than exception-style control flow. The goal is not "fancy"
//   errors. The goal is stable, inspectable failure behavior.
//
// RUN: go run ./04-functions-and-errors/8-error-handling
// ============================================================================

const (
	divisionOp = "division"
	moduloOp   = "modulo"
	sqrtOp     = "sqrt"
)

// MathError is a custom error type that carries structured failure details.
// Callers can inspect these fields with errors.As instead of parsing strings.
type MathError struct {
	Operation string  // Which operation failed.
	InputA    float64 // The primary input.
	InputB    float64 // The secondary input when one exists.
	Message   string  // A human-readable explanation of the failure.
}

// Error satisfies the error interface, which makes MathError a valid error
// value anywhere Go expects an error.
func (e MathError) Error() string {
	switch e.Operation {
	case sqrtOp:
		return fmt.Sprintf("math error in %s (n=%.2f): %s", e.Operation, e.InputA, e.Message)
	default:
		return fmt.Sprintf("math error in %s (a=%.2f, b=%.2f): %s", e.Operation, e.InputA, e.InputB, e.Message)
	}
}

// safeDivide returns a structured error instead of panicking when the caller
// tries to divide by zero.
func safeDivide(a, b int) (float64, error) {
	if b == 0 {
		return 0, &MathError{
			Operation: divisionOp,
			InputA:    float64(a),
			InputB:    float64(b),
			Message:   "division by zero is not allowed",
		}
	}

	return float64(a) / float64(b), nil
}

// safeModulo follows the same explicit-error pattern for integer remainder.
func safeModulo(a, b int) (int, error) {
	if b == 0 {
		return 0, &MathError{
			Operation: moduloOp,
			InputA:    float64(a),
			InputB:    float64(b),
			Message:   "modulo by zero is not allowed",
		}
	}

	return a % b, nil
}

// safeSqrt returns an error for negative input instead of pretending the
// operation succeeded.
func safeSqrt(n float64) (float64, error) {
	if n < 0 {
		return 0, &MathError{
			Operation: sqrtOp,
			InputA:    n,
			Message:   "square root of a negative number is not allowed",
		}
	}

	return math.Sqrt(n), nil
}

// buildSqrtReport wraps a lower-level math failure with higher-level context.
// The wrapped error still preserves the underlying MathError for errors.As.
func buildSqrtReport(n float64) (float64, error) {
	result, err := safeSqrt(n)
	if err != nil {
		return 0, fmt.Errorf("build square-root report: %w", err)
	}

	return result, nil
}

// printMathErrorDetails demonstrates structured error inspection. The error
// string is still useful for logs, but errors.As gives callers stable access
// to the underlying data even after wrapping.
func printMathErrorDetails(err error) {
	var mathErr *MathError
	if !errors.As(err, &mathErr) {
		return
	}

	fmt.Printf("  Operation: %s\n", mathErr.Operation)
	switch mathErr.Operation {
	case sqrtOp:
		fmt.Printf("  Inspectable input: n=%.2f\n", mathErr.InputA)
	default:
		fmt.Printf("  Inspectable inputs: a=%.0f, b=%.0f\n", mathErr.InputA, mathErr.InputB)
	}
}

// runFloatOperation keeps the repeated reporting logic small and clear.
// The deferred message is intentionally tiny: it marks the end of one attempt
// without changing the actual success or error flow.
func runFloatOperation(label string, operation func() (float64, error)) {
	defer fmt.Printf("  [defer] %s attempt finished\n", label)

	value, err := operation()
	if err != nil {
		fmt.Println("  Error:", err)
		printMathErrorDetails(err)
		fmt.Println()
		return
	}

	fmt.Printf("  Result: %.2f\n\n", value)
}

// runIntOperation is the integer counterpart for modulo.
func runIntOperation(label string, operation func() (int, error)) {
	defer fmt.Printf("  [defer] %s attempt finished\n", label)

	value, err := operation()
	if err != nil {
		fmt.Println("  Error:", err)
		printMathErrorDetails(err)
		fmt.Println()
		return
	}

	fmt.Printf("  Result: %d\n\n", value)
}

func main() {
	fmt.Println("=== Safe Math Library Exercise ===")
	fmt.Println()

	fmt.Println("Division:")
	runFloatOperation("division success", func() (float64, error) {
		return safeDivide(10, 4)
	})
	runFloatOperation("division failure", func() (float64, error) {
		return safeDivide(10, 0)
	})

	fmt.Println("Modulo:")
	runIntOperation("modulo success", func() (int, error) {
		return safeModulo(10, 3)
	})
	runIntOperation("modulo failure", func() (int, error) {
		return safeModulo(10, 0)
	})

	fmt.Println("Square Root:")
	runFloatOperation("sqrt success", func() (float64, error) {
		return buildSqrtReport(16)
	})
	runFloatOperation("sqrt failure", func() (float64, error) {
		return buildSqrtReport(-4)
	})

	// KEY TAKEAWAY:
	// - Ordinary failures should return explicit error values, not panic.
	// - Custom errors make failure meaning inspectable and stable.
	// - Wrapping adds context without destroying the original failure identity.
	// - errors.As lets the caller recover structured details safely.
	// - defer is useful for small cleanup or completion behavior, not as a
	//   replacement for normal error handling.
}
