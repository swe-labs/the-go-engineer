// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 4: Functions & Errors — Custom Error Types (Exercise Starter)
// Level: Intermediate
// ============================================================================
//
// EXERCISE: Build a Safe Math Library with Custom Errors
//
// REQUIREMENTS:
//  1. [ ] Define a `MathError` struct with Operation, InputA, InputB, Message fields
//  2. [ ] Implement the `error` interface: `Error() string` method on MathError
//  3. [ ] Implement `safeDivide(a, b int) (float64, error)` — returns error on div by zero
//  4. [ ] Implement `safeModulo(a, b int) (int, error)` — returns error on mod by zero
//  5. [ ] Implement `safeSqrt(n float64) (float64, error)` — returns error on negative input
//  6. [ ] Test all operations in main(), handling both success and error cases
//
// HINTS:
//   - The error interface: type error interface { Error() string }
//   - Return &MathError{...} to create a pointer to your error struct
//   - Use fmt.Sprintf to format the error message
//
// RUN: go run ./04-functions-and-errors/8-error-handling/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Define your MathError struct here

// TODO: Implement the Error() string method

// TODO: Implement safeDivide, safeModulo, safeSqrt

func main() {
	fmt.Println("=== Safe Math Library Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your custom error types and math functions!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
