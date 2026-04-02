// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// Section 7: Strings & Text — String Operations
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Go strings are IMMUTABLE sequences of bytes (usually UTF-8)
//   - The strings package: your Swiss Army knife for text manipulation
//   - strings.Builder: efficient string concatenation
//   - Common operations: split, join, trim, replace, search
//   - Why string concatenation with + is slow (and what to use instead)
//
// ENGINEERING DEPTH:
//   In Go, a string is internally a 2-word struct:
//     Word 1: pointer to the underlying byte array
//     Word 2: length of the string (number of bytes)
//   Strings are IMMUTABLE — every modification creates a NEW string.
//   That's why s = s + "x" in a loop is O(n²) — use strings.Builder instead.
//
// RUN: go run ./07-strings-and-text/1-strings
// ============================================================================

func main() {
	fmt.Println("=== String Operations ===")
	fmt.Println()

	// --- CASE CONVERSION ---
	// strings.ToUpper / ToLower convert all characters.
	// These are Unicode-aware (work with non-English text too).
	domain := "GitHub.COM"
	fmt.Printf("  ToLower: %q → %q\n", domain, strings.ToLower(domain))
	fmt.Printf("  ToUpper: %q → %q\n", domain, strings.ToUpper(domain))
	fmt.Println()

	// --- TRIMMING WHITESPACE ---
	// strings.TrimSpace removes leading AND trailing whitespace.
	// Essential for user input: users often add accidental spaces.
	rawInput := "   rasel9t6@github.com   "
	cleaned := strings.TrimSpace(rawInput)
	fmt.Printf("  TrimSpace: %q → %q (len %d → %d)\n",
		rawInput, cleaned, len(rawInput), len(cleaned))
	fmt.Println()

	// --- SEARCHING ---
	// Contains, HasPrefix, HasSuffix check for substrings.
	// These are O(n) — the string library scans character by character.
	email := "rasel@devops.engineering"
	fmt.Println("  --- Searching ---")
	fmt.Printf("  Contains 'devops': %t\n", strings.Contains(email, "devops"))
	fmt.Printf("  HasPrefix 'rasel': %t\n", strings.HasPrefix(email, "rasel"))
	fmt.Printf("  HasSuffix '.engineering': %t\n", strings.HasSuffix(email, ".engineering"))
	fmt.Printf("  Index of '@': %d\n", strings.Index(email, "@")) // Position of first match
	fmt.Printf("  Count 'e': %d\n", strings.Count(email, "e"))    // How many times 'e' appears
	fmt.Println()

	// --- SPLIT & JOIN ---
	// Split breaks a string into a []string slice by a separator.
	// Join does the reverse: combines a []string with a separator.
	path := "/home/rasel/projects/go-bible"
	parts := strings.Split(path, "/")
	fmt.Printf("  Split %q by '/': %v\n", path, parts)
	// Note: Split("/home/...") produces ["", "home", "rasel", ...] — first element is empty

	// Fields splits by ANY whitespace (better than Split for words)
	logLine := "2025-06-15  INFO   server started on port 8080"
	fields := strings.Fields(logLine)
	fmt.Printf("  Fields: %v\n", fields)

	rejoined := strings.Join(fields, " | ")
	fmt.Printf("  Join: %s\n", rejoined)
	fmt.Println()

	// --- REPLACE ---
	// strings.Replace(s, old, new, n) replaces first n occurrences.
	// Use n = -1 to replace ALL occurrences.
	template := "Hello {name}, welcome to {name}'s dashboard"
	personalized := strings.ReplaceAll(template, "{name}", "Rasel")
	fmt.Printf("  ReplaceAll: %s\n", personalized)
	fmt.Println()

	// --- REPEAT ---
	// strings.Repeat repeats a string n times.
	separator := strings.Repeat("═", 40)
	fmt.Printf("  Repeat: %s\n", separator)
	fmt.Println()

	// --- strings.Builder: EFFICIENT CONCATENATION ---
	// Using + in a loop is O(n²) because strings are immutable.
	// Each += creates a new string and copies all previous content.
	//
	// strings.Builder allocates a growing buffer internally —
	// similar to how append() works for slices. It's O(n).
	//
	// RULE: For building strings in loops, ALWAYS use strings.Builder.
	fmt.Println("  --- strings.Builder ---")
	var b strings.Builder
	languages := []string{"Go", "Rust", "Python", "TypeScript"}

	for i, lang := range languages {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(lang) // Efficient: writes to internal buffer
	}

	result := b.String() // Convert buffer to string (one allocation)
	fmt.Printf("  Built: %s\n", result)
	fmt.Printf("  Builder allocated %d bytes for %d chars\n", b.Cap(), b.Len())

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - Strings are immutable byte sequences (modification = new string)")
	fmt.Println("  - strings.TrimSpace for user input cleanup")
	fmt.Println("  - strings.Fields splits by whitespace (better than Split for words)")
	fmt.Println("  - strings.Builder for efficient loop concatenation (not +)")
	fmt.Println("  - strings.Contains/HasPrefix/HasSuffix for searching")
	fmt.Println("  - strings.ReplaceAll for template-like substitutions")
}
