// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 04: Types and Design
// Title: Unicode and Runes
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Differentiating between bytes, runes, and UTF-8 encoding.
//   - Understanding why len() is misleading for multi-byte text.
//   - Utilizing the 'unicode' and 'unicode/utf8' packages for character processing.
//   - Iterating safely over multi-byte strings using the 'for range' loop.
//
// WHY THIS MATTERS:
//   - In a globalized world, software must handle international text
//     correctly. Go treats strings as immutable byte slices, which
//     means that operations like slicing and indexing can break
//     multi-byte characters if handled naively. Mastering runes and
//     UTF-8 mechanics ensures that your applications are Unicode-compliant
//     and robust against character corruption in international contexts.
//
// RUN:
//   go run ./04-types-design/21-unicode
//
// KEY TAKEAWAY:
//   - Go strings are UTF-8 byte sequences; use runes for character-level logic.
// ============================================================================

// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

//   - Relationship between bytes, runes, and variable-width UTF-8.
//   - Understanding rune as an alias for int32 (Unicode code point).
//   - Safe iteration patterns using 'for range' and byte offset awareness.
//   - Utilizing 'unicode' and 'unicode/utf8' for validation.
//
// TECHNICAL RATIONALE:
//   Go strings are immutable byte slices (`[]byte`) expected to hold
//   UTF-8 encoded text. Because UTF-8 is a variable-width encoding,
//   characters can occupy between 1 and 4 bytes. Using `len()` on a
//   string returns the byte count, not the character count. To
//   process text correctly, developers must work with **runes**
//   (representing Unicode code points) and use the `for range`
//   loop, which handles the complex decoding of multi-byte sequences
//   at the runtime level.
//

func main() {
	fmt.Println("=== Unicode: Bytes vs. Runes ===")
	fmt.Println()

	// 1. Length Discrepancy.
	// len() returns the number of bytes. Multi-byte characters (like 'é')
	// increase the byte count beyond the visible character count.
	fmt.Println("--- Understanding String Length ---")
	text := "café"
	fmt.Printf("  Literal: %q\n", text)
	fmt.Printf("  Bytes:   %d (len())\n", len(text))
	fmt.Printf("  Runes:   %d (utf8.RuneCount)\n", utf8.RuneCountInString(text))
	fmt.Println()

	// 2. Multi-Byte Breakdown.
	// We can inspect the underlying bytes that form a single Unicode character.
	fmt.Println("--- Byte Representation ---")
	emoji := "🚀"
	fmt.Printf("  Emoji: %s\n", emoji)
	fmt.Printf("  Width: %d bytes\n", len(emoji))
	fmt.Printf("  Hex:   %x\n", emoji)
	fmt.Println()

	// 3. Safe Iteration.
	// for range (Loop) decodes UTF-8 byte sequences into individual runes during iteration.
	greeting := "Go 世界!"
	for i, r := range greeting {
		// i is the byte offset where the rune starts.
		fmt.Printf("  Offset %2d: Character '%c' (U+%04X)\n", i, r, r)
	}
	fmt.Println()

	// 4. Character Classification.
	// samples (Slice) provides a collection of diverse runes for classification demonstration.
	samples := []rune{'A', '9', '!', '中'}
	for _, r := range samples {
		isLetter := unicode.IsLetter(r)
		isDigit := unicode.IsDigit(r)
		fmt.Printf("  '%c' -> Letter: %t, Digit: %t\n", r, isLetter, isDigit)
	}

	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ST.4 -> 04-types-design/22-regex")
	fmt.Println("Run    : go run ./04-types-design/22-regex")
	fmt.Println("Current: ST.3 (unicode)")
	fmt.Println("---------------------------------------------------")
}
