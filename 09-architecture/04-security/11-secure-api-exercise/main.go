// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: Secure API
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Build a small API surface that applies validation, auth boundaries, secret-safe behavior, and rate limiting ideas together.
//
// WHY THIS MATTERS:
//   - Security is rarely one big feature. It is the sum of many small boundary decisions made consistently.
//
// RUN:
//   go run ./09-architecture/04-security/11-secure-api-exercise
//
// KEY TAKEAWAY:
//   - Build a small API surface that applies validation, auth boundaries, secret-safe behavior, and rate limiting ideas together.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== SEC.11 Secure API ===")
	fmt.Println("Build a small API surface that applies validation, auth boundaries, secret-safe behavior, and rate limiting ideas together.")
	fmt.Println()
	fmt.Println("- Apply validation, auth, and response discipline together.")
	fmt.Println("- Protect the boundary before business logic runs.")
	fmt.Println("- Treat the exercise as proof that security design is an engineering concern, not a checklist afterthought.")
	fmt.Println()
	fmt.Println("Secure systems are built from layered controls that keep working even when one layer is stressed or misused.")
	fmt.Println("NEXT UP: SL.1 -> 10-production/01-structured-logging/1-slog-basics")
}
