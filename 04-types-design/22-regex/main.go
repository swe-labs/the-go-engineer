// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Regex
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Compiling regular expressions using Compile and MustCompile.
//   - Applying pattern matching, extraction, and replacement operations.
//   - Utilizing capture groups and named sub-expressions for data parsing.
//   - Understanding the performance characteristics of Go's RE2 engine.
//
// WHY THIS MATTERS:
//   - Regular expressions provide a powerful, standardized mechanism for
//     validating input, parsing complex logs, and transforming text data.
//     Go's 'regexp' package uses the RE2 engine, which guarantees linear-time
//     complexity and prevents catastrophic backtracking vulnerabilities.
//     Mastering regex is essential for building robust data validation and
//     text-processing pipelines in high-scale systems.
//
// RUN:
//   go run ./04-types-design/22-regex
//
// KEY TAKEAWAY:
//   - Go's RE2 engine provides linear-time pattern matching and safe data extraction.
// ============================================================================

// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"regexp"
)

//   - Compiling state machines with Compile and MustCompile.
//   - Linear-time matching guarantees with the RE2 engine.
//   - Utilizing capture groups for structured text decomposition.
//   - Pattern-based redaction and transformation.
//
// TECHNICAL RATIONALE:
//   Go's 'regexp' package implements the RE2 syntax, which uses a
//   deterministic finite automaton (DFA) approach. Unlike
//   backtracking-based engines, RE2 guarantees O(n) execution time relative
//   to the input length. While this omits features like lookaheads,
//   it prevents "Regular Expression Denial of Service" (ReDoS) attacks.
//   Compilation is a significant one-time cost, as Go builds a
//   complex internal state machine to represent the pattern.
//

func main() {
	fmt.Println("=== Regular Expressions: Pattern Matching ===")
	fmt.Println()

	// 1. Compilation & Performance.
	// emailRegex (regexp.Regexp) encapsulates a compiled state machine for linear-time email validation.
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	fmt.Printf("  Pattern initialized: %s\n", emailRegex.String())
	fmt.Println()

	// 2. Validation & Predicates.
	// MatchString returns a boolean indicating if the pattern exists in the input.
	fmt.Println("--- Input Validation ---")
	inputs := []string{"admin@opslane.io", "malformed-email", "dev@github.com"}
	for _, input := range inputs {
		isValid := emailRegex.MatchString(input)
		fmt.Printf("  %-18s | Valid: %t\n", input, isValid)
	}
	fmt.Println()

	// 3. Structured Data Extraction.
	// logPattern (regexp.Regexp) decomposes log entries into timestamp, level, and message components.
	logPattern := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})\s+\[(\w+)\]\s+(.+)$`)
	logEntry := "2025-06-15 [INFO] cache.warmup completed"
	match := logPattern.FindStringSubmatch(logEntry)
	if len(match) == 4 {
		// Index 0 is the full match; indices 1-3 are capture groups.
		fmt.Printf("  Timestamp: %s\n", match[1])
		fmt.Printf("  Level:     %s\n", match[2])
		fmt.Printf("  Message:   %s\n", match[3])
	}
	fmt.Println()

	// 4. Pattern-Based Transformation.
	// ReplaceAllString enables bulk redaction or transformation.
	fmt.Println("--- Transformation & Redaction ---")
	rawText := "Internal notification sent to sysadmin@opslane.io and secondary@opslane.io"
	redacted := emailRegex.ReplaceAllString(rawText, "[REDACTED]")
	fmt.Printf("  Original: %s\n", rawText)
	fmt.Printf("  Redacted: %s\n", redacted)

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ST.5 -> 04-types-design/23-text-template")
	fmt.Println("Run    : go run ./04-types-design/23-text-template")
	fmt.Println("Current: ST.4 (regex)")
	fmt.Println("---------------------------------------------------")
}
