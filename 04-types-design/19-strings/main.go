// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Strings
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Understanding string immutability and the underlying byte-slice representation.
//   - Utilizing the 'strings' package for efficient text transformation.
//   - Leveraging 'strings.Builder' to avoid O(n^2) concatenation performance traps.
//   - Mastering search, split, join, and replacement operations.
//
// WHY THIS MATTERS:
//   - Strings are one of the most frequently used types in Go, yet their
//     immutable nature often leads to performance bottlenecks if handled
//     incorrectly. Every modification to a string creates a new allocation
//     in memory. Mastering the 'strings' package and 'strings.Builder'
//     ensures that your applications remain memory-efficient and
//     performant when processing large volumes of text data.
//
// RUN:
//   go run ./04-types-design/19-strings
//
// KEY TAKEAWAY:
//   - Go strings are immutable; use 'strings.Builder' for efficient mutations.
// ============================================================================

// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"strings"
)

//   - Internal memory representation of strings as two-word headers.
//   - Utilizing 'strings.Builder' for efficient allocation-aware mutation.
//   - Search, split, and transformation patterns in the standard library.
//
// TECHNICAL RATIONALE:
//   In Go, a string is a read-only header containing a pointer to an
//   underlying byte array and a length. Because strings are immutable,
//   every modification (like concatenation with '+') results in a
//   new memory allocation and a full copy of the existing data. To
//   avoid O(n^2) performance degradation during bulk text
//   construction, developers must use 'strings.Builder', which
//   manages a growing internal buffer to minimize allocations.
//

func main() {
	fmt.Println("=== Strings: Immutable Byte Sequences ===")
	fmt.Println()

	// 1. Normalization & Cleanup.
	// Standardizing input by converting case and trimming whitespace.
	fmt.Println("--- Normalization ---")
	rawInput := "   USER-input@Domain.COM   "
	normalized := strings.ToLower(strings.TrimSpace(rawInput))
	fmt.Printf("  Original:   %q\n", rawInput)
	fmt.Printf("  Normalized: %q\n", normalized)
	fmt.Println()

	// 2. Pattern Matching & Searching.
	// Efficient O(n) scans for prefixes, suffixes, and substrings.
	fmt.Println("--- Pattern Matching ---")
	uri := "https://api.github.com/v1/repos"
	fmt.Printf("  Secure: %t (HasPrefix)\n", strings.HasPrefix(uri, "https"))
	fmt.Printf("  API:    %t (Contains)\n", strings.Contains(uri, "/api/"))
	fmt.Printf("  Index of 'github': %d\n", strings.Index(uri, "github"))
	fmt.Println()

	// 3. Transformation & Tokenization.
	// Splitting strings into slices and joining them back.
	fmt.Println("--- Tokenization ---")
	logEntry := "2025-06-15  INFO   service.start  port=8080"
	tokens := strings.Fields(logEntry) // Splits by any whitespace
	fmt.Printf("  Tokens: %q\n", tokens)
	fmt.Printf("  Join:   %s\n", strings.Join(tokens, " | "))
	fmt.Println()

	// 4. Efficient Concatenation.
	// b (strings.Builder) manages an internal byte buffer for allocation-efficient string construction.
	var b strings.Builder
	items := []string{"core", "auth", "cache", "db"}

	for i, item := range items {
		if i > 0 {
			b.WriteString(" -> ")
		}
		b.WriteString(item)
	}
	fmt.Printf("  Pipeline: %s\n", b.String())

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ST.2 -> 04-types-design/20-formatting")
	fmt.Println("Run    : go run ./04-types-design/20-formatting")
	fmt.Println("Current: ST.1 (strings)")
	fmt.Println("---------------------------------------------------")
}
