// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: JWT - implementation and risks
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn what a JWT contains, how signing works, and why tokens still create real operational risk when used carelessly.
//
// WHY THIS MATTERS:
//   - A JWT is a signed claim set, not a trust system by itself.
//
// RUN:
//   go run ./09-architecture/04-security/5-jwt-implementation-and-risks
//
// KEY TAKEAWAY:
//   - Learn what a JWT contains, how signing works, and why tokens still create real operational risk when used carelessly.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== SEC.5 JWT - implementation and risks ===")
	fmt.Println("Learn what a JWT contains, how signing works, and why tokens still create real operational risk when used carelessly.")
	fmt.Println()
	fmt.Println("- Signing proves integrity, not that every claim is safe to trust blindly.")
	fmt.Println("- Validate issuer, audience, expiry, and algorithm policy.")
	fmt.Println("- Treat tokens as credentials and logs as hostile to secret material.")
	fmt.Println()
	fmt.Println("JWT mistakes are rarely library mistakes. They are usually trust-boundary mistakes like weak key policy, bad expiry handling, or missing audience checks.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.6 -> 09-architecture/04-security/6-password-hashing")
	fmt.Println("Current: SEC.5 (jwt - implementation and risks)")
	fmt.Println("---------------------------------------------------")
}
