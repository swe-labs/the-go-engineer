// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: Input validation patterns
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how allow-lists, normalization, and fail-fast checks turn raw input into trustworthy domain values.
//
// WHY THIS MATTERS:
//   - Validation is boundary work that decides whether input is acceptable before business logic depends on it.
//
// RUN:
//   go run ./09-architecture/04-security/01-input-validation-patterns
//
// KEY TAKEAWAY:
//   - Learn how allow-lists, normalization, and fail-fast checks turn raw input into trustworthy domain values.
// ============================================================================

package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// allowListedUsernames (Map): allow-list of acceptable usernames for access control.
var allowListedUsernames = map[string]bool{
	"alice":   true,
	"bob":     true,
	"charlie": true,
}

// validEmailRegex matches a basic email format.
var validEmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// allowedPathCharsRegex rejects paths that contain characters outside the allow-list.
var allowedPathCharsRegex = regexp.MustCompile(`^[a-zA-Z0-9_./-]+$`)

// validateUsername (Function): normalizes and checks raw input against the username allow-list.
func validateUsername(raw string) (string, bool) {
	normalized := strings.TrimSpace(strings.ToLower(raw))
	if !allowListedUsernames[normalized] {
		return normalized, false
	}
	return normalized, true
}

// validateEmail (Function): normalizes and validates email format against a regex pattern.
func validateEmail(raw string) (string, bool) {
	normalized := strings.TrimSpace(strings.ToLower(raw))
	if !validEmailRegex.MatchString(normalized) {
		return normalized, false
	}
	return normalized, true
}

// validateFilePath (Function): normalizes and checks file path characters against an allow-list regex.
func validateFilePath(raw string) (string, bool) {
	normalized := strings.TrimSpace(raw)
	if !allowedPathCharsRegex.MatchString(normalized) {
		return normalized, false
	}
	return normalized, true
}

func main() {
	fmt.Println("=== SEC.1 Input validation patterns ===")
	fmt.Println()

	cases := []struct {
		label string
		valid func(string) (string, bool)
		input string
	}{
		{"username allow-list", validateUsername, "Alice"},
		{"username allow-list", validateUsername, "mallory"},
		{"username allow-list", validateUsername, " bob "},
		{"email format", validateEmail, "Alice@Example.COM"},
		{"email format", validateEmail, "not-an-email"},
		{"email format", validateEmail, " user@example.com "},
		{"file path chars", validateFilePath, "../safe/path-v2.go"},
		{"file path chars", validateFilePath, "unsafe<path>.go"},
	}

	for _, c := range cases {
		normalized, ok := c.valid(c.input)
		if ok {
			fmt.Printf("  PASS  %-22s %q -> %q\n", c.label, c.input, normalized)
		} else {
			fmt.Printf("  FAIL  %-22s %q -> normalized=%q\n", c.label, c.input, normalized)
		}
	}

	log.Println("Validation patterns demonstrated successfully.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.2 -> 09-architecture/04-security/02-sql-injection-prevention")
	fmt.Println("Current: SEC.1 (input validation patterns)")
	fmt.Println("---------------------------------------------------")
}
