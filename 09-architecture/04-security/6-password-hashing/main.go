// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: Password hashing
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn why passwords need one-way hashing with slow work factors instead of reversible encryption.
//
// WHY THIS MATTERS:
//   - Passwords should be verified, not decrypted.
//
// RUN:
//   go run ./09-architecture/04-security/6-password-hashing
//
// KEY TAKEAWAY:
//   - Learn why passwords need one-way hashing with slow work factors instead of reversible encryption.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== SEC.6 Password hashing ===")
	fmt.Println("Learn why passwords need one-way hashing with slow work factors instead of reversible encryption.")
	fmt.Println()
	fmt.Println("- Use dedicated password hashing algorithms, not generic hashes.")
	fmt.Println("- Store the hash and parameters, never the plaintext password.")
	fmt.Println("- Revisit work factors as hardware changes.")
	fmt.Println()
	fmt.Println("Hashing policy is security policy: weak parameters make a correct library call functionally useless.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.7")
	fmt.Println("Current: SEC.6 (password hashing)")
	fmt.Println("---------------------------------------------------")
}
