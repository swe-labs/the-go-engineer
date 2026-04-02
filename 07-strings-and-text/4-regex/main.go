// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"regexp"
)

// ============================================================================
// Section 7: Strings & Text — Regular Expressions
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - regexp.Compile vs regexp.MustCompile (compile-time vs runtime safety)
//   - MatchString: does the text match the pattern?
//   - FindString / FindAllString: extracting matches
//   - FindStringSubmatch: capturing groups (data extraction)
//   - ReplaceAllString: search-and-replace with patterns
//   - Named capture groups: (?P<name>pattern)
//
// ENGINEERING DEPTH:
//   Go's regex engine uses RE2 (linear time, guaranteed no backtracking).
//   This means some features from Perl/Python regex DON'T exist in Go:
//     ❌ Lookaheads/lookbehinds (?=...) (?<=...)
//     ❌ Backreferences \1
//   But RE2 guarantees O(n) performance — no catastrophic backtracking.
//   Always prefer MustCompile for compile-time-known patterns (panics on bad regex).
//
// RUN: go run ./07-strings-and-text/4-regex
// ============================================================================

func main() {
	fmt.Println("=== Regular Expressions ===")
	fmt.Println()

	// =====================================================================
	// 1. MustCompile vs Compile
	// =====================================================================
	// MustCompile PANICS if the regex is invalid — use for hardcoded patterns.
	// Compile returns an ERROR — use for user-provided patterns.
	emailPattern := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	// The ` ` (backtick) string avoids escaping issues. Prefer backticks for regex.

	// =====================================================================
	// 2. MatchString: Does text contain a match?
	// =====================================================================
	fmt.Println("1️⃣  MatchString (true/false check):")
	testEmails := []string{
		"rasel@devops.io",
		"not-an-email",
		"admin@github.com",
		"bad@",
	}
	for _, e := range testEmails {
		match := emailPattern.MatchString(e)
		icon := "❌"
		if match {
			icon = "✅"
		}
		fmt.Printf("   %s %q\n", icon, e)
	}
	fmt.Println()

	// =====================================================================
	// 3. FindString / FindAllString: Extract matches
	// =====================================================================
	fmt.Println("2️⃣  FindAllString (extract all matches):")
	text := "Contact us at support@gomastery.dev or sales@gomastery.dev. Invalid: @broken"

	// FindAllString returns all non-overlapping matches.
	// The second argument limits results (-1 = all matches).
	allEmails := emailPattern.FindAllString(text, -1)
	for i, email := range allEmails {
		fmt.Printf("   Match %d: %s\n", i+1, email)
	}
	fmt.Println()

	// =====================================================================
	// 4. FindStringSubmatch: Capturing groups (data extraction)
	// =====================================================================
	// Parentheses in regex create CAPTURE GROUPS.
	// FindStringSubmatch returns: [fullMatch, group1, group2, ...]
	fmt.Println("3️⃣  Capture groups (extract structured data):")

	// Pattern to extract parts of a log line:
	//   2025-06-15 INFO Server started on port 8080
	logPattern := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})\s+(\w+)\s+(.+)`)
	logLines := []string{
		"2025-06-15 INFO Server started on port 8080",
		"2025-06-15 ERROR Database connection timeout",
		"2025-06-16 WARN Memory usage above 90%",
	}

	for _, line := range logLines {
		match := logPattern.FindStringSubmatch(line)
		if match != nil {
			// match[0] = full match, match[1] = date, match[2] = level, match[3] = message
			fmt.Printf("   Date: %s  Level: %-5s  Msg: %s\n", match[1], match[2], match[3])
		}
	}
	fmt.Println()

	// =====================================================================
	// 5. Named capture groups
	// =====================================================================
	// (?P<name>pattern) creates a named group accessible by name.
	fmt.Println("4️⃣  Named capture groups:")
	urlPattern := regexp.MustCompile(`(?P<protocol>https?)://(?P<host>[^/:]+)(?::(?P<port>\d+))?`)

	url := "https://api.gomastery.dev:8443"
	match := urlPattern.FindStringSubmatch(url)
	if match != nil {
		for i, name := range urlPattern.SubexpNames() {
			if name != "" && i < len(match) {
				fmt.Printf("   %s = %q\n", name, match[i])
			}
		}
	}
	fmt.Println()

	// =====================================================================
	// 6. ReplaceAllString: Search and replace
	// =====================================================================
	fmt.Println("5️⃣  ReplaceAllString:")
	// Redact all email addresses in text
	redacted := emailPattern.ReplaceAllString(text, "[REDACTED]")
	fmt.Printf("   Original: %s\n", text)
	fmt.Printf("   Redacted: %s\n", redacted)
	fmt.Println()

	// Replace with a function for dynamic replacement
	sensitiveData := "SSN: 123-45-6789, Phone: 555-867-5309"
	ssnPattern := regexp.MustCompile(`\d{3}-\d{2}-\d{4}`)
	masked := ssnPattern.ReplaceAllStringFunc(sensitiveData, func(match string) string {
		return "XXX-XX-" + match[len(match)-4:]
	})
	fmt.Printf("   Masked: %s\n", masked)

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - MustCompile for hardcoded patterns, Compile for user input")
	fmt.Println("  - MatchString checks, FindAllString extracts, ReplaceAllString replaces")
	fmt.Println("  - Capture groups ( ) extract structured data from text")
	fmt.Println("  - Named groups (?P<name>) improve readability")

}
