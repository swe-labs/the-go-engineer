// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: OWASP Top 10 for Go
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Turn the OWASP Top 10 into a practical checklist for Go services and review conversations.
//
// WHY THIS MATTERS:
//   - The OWASP list is a prioritization aid, not a complete security strategy.
//
// RUN:
//   go run ./09-architecture/04-security/10-owasp-top-10-for-go
//
// KEY TAKEAWAY:
//   - Turn the OWASP Top 10 into a practical checklist for Go services and review conversations.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== SEC.10 OWASP Top 10 for Go ===")
	fmt.Println("Turn the OWASP Top 10 into a practical checklist for Go services and review conversations.")
	fmt.Println()
	fmt.Println("- Use checklists to drive review questions.")
	fmt.Println("- Map common risks to concrete Go boundary decisions.")
	fmt.Println("- Revisit the list whenever a new public surface is added.")
	fmt.Println()
	fmt.Println("Security checklists matter most when they are embedded into normal engineering work instead of living in a forgotten audit document.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.11")
	fmt.Println("Current: SEC.10 (owasp top 10 for go)")
	fmt.Println("---------------------------------------------------")
}
