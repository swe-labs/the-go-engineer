// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Custom Error Types
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how to define custom error types that carry structured information for better error handling.
//
// WHY THIS MATTERS:
//   - Think of a boarding pass. A simple "flight delayed" message is not enough. A good error includes: flight number, original time, new time, reason, a...
//
// RUN:
//   go run ./04-types-design/8-custom-errors
//
// KEY TAKEAWAY:
//   - Learn how to define custom error types that carry structured information for better error handling.
// ============================================================================

package main

import (
	"errors"
	"fmt"
)

//
//   - Defining custom error types
//   - Adding structured information to errors
//   - Type assertions for specific error handling
//

type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s=%v - %s", e.Field, e.Value, e.Message)
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

func validate(name string, age int) error {
	if name == "" {
		return ValidationError{Field: "name", Value: name, Message: "cannot be empty"}
	}
	if age < 0 {
		return ValidationError{Field: "age", Value: age, Message: "cannot be negative"}
	}
	if age > 150 {
		return ValidationError{Field: "age", Value: age, Message: "seems unrealistic"}
	}
	return nil
}

func findUser(id string) error {
	if id == "" {
		return ValidationError{Field: "id", Value: id, Message: "cannot be empty"}
	}
	if id != "123" {
		return NotFoundError{Resource: "user", ID: id}
	}
	return nil
}

func handleError(err error) {
	if err == nil {
		fmt.Println("  No error")
		return
	}

	var validationErr ValidationError
	if errors.As(err, &validationErr) {
		fmt.Printf("  Validation failed: %s\n", validationErr.Message)
		fmt.Printf("    Field: %s, Value: %v\n", validationErr.Field, validationErr.Value)
		return
	}

	var notFoundErr NotFoundError
	if errors.As(err, &notFoundErr) {
		fmt.Printf("  Not found: %s\n", notFoundErr.Error())
		return
	}

	fmt.Printf("  Unknown error: %s\n", err.Error())
}

func main() {
	fmt.Println("=== Custom Error Types ===")
	fmt.Println()

	fmt.Println("--- Validation Errors ---")
	err := validate("", 25)
	handleError(err)

	err = validate("John", -5)
	handleError(err)

	err = validate("John", 200)
	handleError(err)

	err = validate("John", 30)
	handleError(err)

	fmt.Println()
	fmt.Println("--- Not Found Errors ---")
	err = findUser("")
	handleError(err)

	err = findUser("999")
	handleError(err)

	err = findUser("123")
	handleError(err)

	fmt.Println()
	fmt.Println("--- Wrapped Errors (Go 1.13+) ---")
	err = validate("", 0)
	fmt.Printf("  Original: %v\n", err)

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Custom errors are structs that implement Error() string")
	fmt.Println("  - Add fields to carry structured error information")
	fmt.Println("  - Use errors.As() to check specific error types")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: TI.9 -> 04-types-design/9-generics")
	fmt.Println("Current: TI.8 (custom-errors)")
	fmt.Println("Previous: TI.7 (receiver-sets)")
	fmt.Println("---------------------------------------------------")
}
