// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Custom Error Types
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Defining custom types that satisfy the `error` interface.
//   - Adding structured metadata to errors for programmatic inspection.
//   - Using `errors.As` and `errors.Is` for reliable type-safe error checking.
//
// WHY THIS MATTERS:
//   - String-only errors are insufficient for complex systems. Custom
//     error types allow you to carry machine-readable context (like
//     status codes or field names), enabling sophisticated recovery
//     logic and better observability in production.
//
// RUN:
//   go run ./04-types-design/8-custom-errors
//
// KEY TAKEAWAY:
//   - Custom errors provide structured context for failure analysis.
// ============================================================================

package main

import (
	"errors"
	"fmt"
)

// Section 04: Types & Design - Custom Error Types

// ValidationError is returned when input data fails domain-specific rules.
// ValidationError (Struct): is returned when input data fails domain-specific rules.
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

// Error implements the built-in error interface for ValidationError.
// ValidationError.Error (Method): implements the built-in error interface for ValidationError.
func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s=%v - %s", e.Field, e.Value, e.Message)
}

// NotFoundError is returned when a requested resource does not exist.
// NotFoundError (Struct): is returned when a requested resource does not exist.
type NotFoundError struct {
	Resource string
	ID       string
}

// Error implements the built-in error interface for NotFoundError.
// NotFoundError.Error (Method): implements the built-in error interface for NotFoundError.
func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

// validate (Function): runs the validate step and keeps its inputs, outputs, or errors visible.
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

// findUser (Function): runs the find user step and keeps its inputs, outputs, or errors visible.
func findUser(id string) error {
	if id == "" {
		return ValidationError{Field: "id", Value: id, Message: "cannot be empty"}
	}
	if id != "123" {
		return NotFoundError{Resource: "user", ID: id}
	}
	return nil
}

// handleError (Function): runs the handle error step and keeps its inputs, outputs, or errors visible.
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
	fmt.Println("=== Custom Error Types: Structured Context ===")
	fmt.Println()

	// 1. Validation Failures.
	// Custom error types carry domain-specific context (like Field names)
	// for precise error reporting.
	fmt.Println("--- Domain Validation ---")
	handleError(validate("", 25))
	handleError(validate("John", -5))
	handleError(validate("John", 200))
	fmt.Println()

	// 2. Resource Discovery.
	// Distinct error types allow callers to differentiate between client
	// errors (validation) and system/data errors (not found).
	fmt.Println("--- Resource Lookup ---")
	handleError(findUser(""))
	handleError(findUser("999"))
	fmt.Println()

	// 3. Error Inspection (errors.As).
	// Using errors.As allows for type-safe unwrapping of structured error
	// data, supporting modern error wrapping patterns.
	fmt.Println("--- Type-Safe Inspection ---")
	err := validate("Ava", -1)
	var vErr ValidationError
	if errors.As(err, &vErr) {
		fmt.Printf("  Caught validation error on field: %s\n", vErr.Field)
	}

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: TI.9 -> 04-types-design/9-generics")
	fmt.Println("Run    : go run ./04-types-design/9-generics")
	fmt.Println("Current: TI.8 (custom-errors)")
	fmt.Println("---------------------------------------------------")
}
