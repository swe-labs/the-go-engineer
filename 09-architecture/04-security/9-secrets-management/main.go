// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: Secrets management
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how to keep credentials, keys, and tokens out of source control, logs, and casual developer workflows.
//
// WHY THIS MATTERS:
//   - A secret is any value that grants access or trust if an attacker obtains it.
//
// RUN:
//   go run ./09-architecture/04-security/9-secrets-management
//
// KEY TAKEAWAY:
//   - Learn how to keep credentials, keys, and tokens out of source control, logs, and casual developer workflows.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== SEC.9 Secrets management ===")
	fmt.Println("Learn how to keep credentials, keys, and tokens out of source control, logs, and casual developer workflows.")
	fmt.Println()
	fmt.Println("- Load secrets from dedicated config channels, not checked-in files.")
	fmt.Println("- Never log secret material intentionally or accidentally.")
	fmt.Println("- Rotate secrets and scope them narrowly.")
	fmt.Println()
	fmt.Println("Most secret leaks are boring process failures, not cryptographic failures: committed files, copied logs, and reused credentials.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.10")
	fmt.Println("Current: SEC.9 (secrets management)")
	fmt.Println("---------------------------------------------------")
}
