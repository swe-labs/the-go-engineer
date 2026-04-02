// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ============================================================================
// Section 7: Strings & Text — Unicode & Runes
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - bytes vs runes: why len("café") is 5, not 4
//   - What a rune is: an alias for int32, representing a Unicode code point
//   - How Go stores strings as UTF-8 bytes (variable-width encoding)
//   - Iterating over strings correctly (for range uses runes, not bytes)
//   - The unicode package: character classification (letter, digit, etc.)
//   - utf8.RuneCountInString: counting actual characters, not bytes
//
// ENGINEERING DEPTH:
//   Go strings are just byte slices ([]byte) under the hood.
//   UTF-8 encoding uses 1-4 bytes per character:
//     'A' = 1 byte (ASCII compatible)
//     'é' = 2 bytes (Latin with accent)
//     '中' = 3 bytes (CJK character)
//     '🚀' = 4 bytes (emoji)
//
//   This is why len("café") returns 5 (4 letters + 1 extra byte for 'é').
//   Use utf8.RuneCountInString() for the actual character count.
//
// RUN: go run ./07-strings-and-text/3-unicode
// ============================================================================

func main() {
	fmt.Println("=== Unicode & Runes ===")
	fmt.Println()

	// --- BYTES vs RUNES ---
	// len() counts BYTES, not characters. This surprises many beginners.
	text := "café"
	fmt.Println("--- Bytes vs Runes ---")
	fmt.Printf("  Text:           %q\n", text)
	fmt.Printf("  len() (bytes):  %d\n", len(text))                    // 5 bytes
	fmt.Printf("  RuneCount:      %d\n", utf8.RuneCountInString(text)) // 4 characters
	fmt.Println()

	// Why? Because 'é' is encoded as 2 bytes in UTF-8:
	fmt.Println("  Byte breakdown:")
	for i, b := range []byte(text) {
		fmt.Printf("    byte[%d] = %d (0x%02x)", i, b, b)
		if b < 128 {
			fmt.Printf(" → '%c' (ASCII)\n", b)
		} else {
			fmt.Printf(" → part of multi-byte character\n")
		}
	}
	fmt.Println()

	// --- RUNE: int32 alias for Unicode code points ---
	// A rune represents a single Unicode character, regardless of byte width.
	// Declaration: var r rune = '中' or just r := '中'
	fmt.Println("--- Rune Examples ---")
	examples := []rune{'G', 'o', 'é', '中', '🚀'}
	for _, r := range examples {
		fmt.Printf("  '%c' → Unicode: U+%04X, Bytes: %d\n",
			r, r, utf8.RuneLen(r))
	}
	fmt.Println()

	// --- FOR RANGE iterates by RUNE (correct) ---
	// A regular for loop (i=0; i<len; i++) iterates by BYTE — breaks on multi-byte chars.
	// for range iterates by RUNE — the correct way for text processing.
	greeting := "Hello, 世界! 🌍"
	fmt.Println("--- for range (iterates by rune) ---")
	fmt.Printf("  Text: %q\n", greeting)
	for i, r := range greeting {
		// i = byte offset (not character index!)
		// r = the rune (character) at that offset
		if !unicode.IsSpace(r) && !unicode.IsPunct(r) {
			fmt.Printf("    byte_offset=%2d  rune='%c'  code=U+%04X\n", i, r, r)
		}
	}
	fmt.Println()

	// --- UNICODE PACKAGE: Character Classification ---
	// The unicode package provides functions to classify characters.
	// Essential for input validation, parsing, and text processing.
	fmt.Println("--- unicode package ---")
	testChars := []rune{'A', '7', ' ', '!', 'é', '中', '🚀'}
	for _, r := range testChars {
		props := []string{}
		if unicode.IsLetter(r) {
			props = append(props, "letter")
		}
		if unicode.IsDigit(r) {
			props = append(props, "digit")
		}
		if unicode.IsSpace(r) {
			props = append(props, "space")
		}
		if unicode.IsPunct(r) {
			props = append(props, "punct")
		}
		if unicode.IsUpper(r) {
			props = append(props, "upper")
		}
		if unicode.IsLower(r) {
			props = append(props, "lower")
		}
		if len(props) == 0 {
			props = append(props, "symbol/other")
		}
		fmt.Printf("  '%c' (U+%04X): %s\n", r, r, strings.Join(props, ", "))
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Go strings are UTF-8 byte slices (1-4 bytes per character)")
	fmt.Println("  - len() counts BYTES, utf8.RuneCountInString() counts CHARACTERS")
	fmt.Println("  - for range iterates by rune (correct), for i iterates by byte (wrong)")
	fmt.Println("  - rune = int32 alias for a Unicode code point")
	fmt.Println("  - unicode.IsLetter/IsDigit/IsSpace for character classification")
}
