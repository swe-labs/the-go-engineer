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
// Section 4: Functions & Errors — Error Wrapping
// Level: Intermediate → Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - fmt.Errorf with %w — wrapping errors with additional context
//   - errors.Unwrap — peeling back layers of wrapped errors
//   - errors.Is — checking if ANY error in a chain matches a target
//   - errors.As — extracting a specific error type from a chain
//   - errors.Join — combining multiple errors (Go 1.20+)
//   - Why string matching errors is an ANTI-PATTERN
//
// WHY ERROR WRAPPING MATTERS:
//   Without wrapping, you lose context:
//     return err → "file not found"  (WHERE? WHICH file? WHAT was happening?)
//   With wrapping, you get a trace:
//     return fmt.Errorf("loading config: %w", err)
//     → "loading config: opening /etc/app.conf: file not found"
//
// ENGINEERING DEPTH:
//   When you use `fmt.Errorf("%w")`, Go does not just concatenate a string.
//   It returns a special private struct type inside the `fmt` package called
//   `wrapError`. This struct holds BOTH the formatted string AND a pointer to
//   the original underlying error. When you call `errors.Is()`, the function
//   dynamically unwraps the struct layer by layer, inspecting the pointers
//   until it finds a memory match. This makes Go's error handling as powerful
//   as Java Exceptions, without the hidden control-flow nightmares.
//
// RUN: go run ./04-functions-and-errors/5b-error-wrapping
// ============================================================================

// --- SENTINEL ERRORS ---
// Sentinel errors are package-level variables that represent specific conditions.
// Other packages compare against these using errors.Is().
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

// ValidationError is a custom error type with structured data.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on '%s': %s", e.Field, e.Message)
}

// --- SIMULATED LAYERS OF AN APPLICATION ---

// findUser simulates a database lookup that may fail.
func findUser(id int) (string, error) {
	if id <= 0 {
		// Return a sentinel error — the caller can check with errors.Is()
		return "", ErrNotFound
	}
	if id == 999 {
		return "", &ValidationError{Field: "id", Message: "reserved ID"}
	}
	return "Alice", nil
}

// getUser wraps the error from findUser with additional context.
// The %w verb creates a WRAPPED error — the original error is preserved
// inside, and errors.Is() / errors.As() can still find it.
func getUser(id int) (string, error) {
	user, err := findUser(id)
	if err != nil {
		// WRAPPING with %w: adds context while preserving the original error
		return "", fmt.Errorf("getUser(id=%d): %w", id, err)
	}
	return user, nil
}

// handleRequest wraps again — building a chain of context.
func handleRequest(userID int) error {
	_, err := getUser(userID)
	if err != nil {
		return fmt.Errorf("handleRequest at %s: %w",
			time.Now().Format("15:04:05"), err)
	}
	return nil
}

func main() {
	fmt.Println("=== Error Wrapping with %w ===")
	fmt.Println()

	// --- BUILD AN ERROR CHAIN ---
	err := handleRequest(0)
	fmt.Printf("Full error chain:\n  %v\n\n", err)
	// Output: "handleRequest at 12:00:00: getUser(id=0): not found"
	// Each layer adds context — you can trace exactly where it went wrong.

	// --- errors.Is: Check if ANY error in the chain matches ---
	// errors.Is traverses the entire wrapped chain looking for a match.
	// This is why wrapping with %w is powerful — you don't lose the original error.
	fmt.Println("=== errors.Is (chain traversal) ===")
	if errors.Is(err, ErrNotFound) {
		fmt.Println("  ✅ errors.Is(err, ErrNotFound) = true")
		fmt.Println("     Even though err is wrapped multiple times!")
	}
	if !errors.Is(err, ErrUnauthorized) {
		fmt.Println("  ✅ errors.Is(err, ErrUnauthorized) = false (correct)")
	}
	fmt.Println()

	// --- errors.As: Extract a specific error TYPE from the chain ---
	// errors.As traverses the chain and extracts the FIRST error
	// that can be assigned to the target type.
	fmt.Println("=== errors.As (type extraction) ===")
	err2 := handleRequest(999) // This triggers a ValidationError
	fmt.Printf("Error: %v\n", err2)

	var valErr *ValidationError
	if errors.As(err2, &valErr) {
		fmt.Printf("  ✅ Found ValidationError in chain:\n")
		fmt.Printf("     Field:   %s\n", valErr.Field)
		fmt.Printf("     Message: %s\n", valErr.Message)
	}
	fmt.Println()

	// --- errors.Join: Combine multiple errors (Go 1.20+) ---
	// When multiple operations fail and you want to report ALL errors,
	// not just the first one.
	fmt.Println("=== errors.Join (multiple errors) ===")
	err3 := errors.Join(
		errors.New("database connection failed"),
		errors.New("cache unavailable"),
		errors.New("config file missing"),
	)
	fmt.Printf("  Joined error: %v\n", err3)
	fmt.Println()

	// --- ANTI-PATTERN: String matching ---
	fmt.Println("=== Anti-Pattern (DO NOT DO THIS) ===")
	fmt.Println("  ❌ if err.Error() == \"not found\" { ... }    ← FRAGILE")
	fmt.Println("  ❌ if strings.Contains(err.Error(), ...) { } ← BREAKS WITH WRAPPING")
	fmt.Println("  ✅ if errors.Is(err, ErrNotFound) { ... }    ← CORRECT")
	fmt.Println("  ✅ if errors.As(err, &myErr) { ... }          ← CORRECT")

	// KEY TAKEAWAY:
	// - Wrap errors with fmt.Errorf("context: %w", err) — ALWAYS add context
	// - Use errors.Is() to check for sentinel errors through wrapped chains
	// - Use errors.As() to extract typed errors through wrapped chains
	// - Use errors.Join() when you have multiple independent errors
	// - NEVER match error strings — they break when errors are wrapped
}
