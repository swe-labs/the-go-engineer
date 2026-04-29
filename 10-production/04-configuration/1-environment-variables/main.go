// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: Environment variables
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how environment variables shape runtime configuration without rebuilding the binary.
//
// WHY THIS MATTERS:
//   - Environment variables are process-level inputs provided by the runtime environment, not by the source code itself.
//
// RUN:
//   go run ./10-production/04-configuration/1-environment-variables
//
// KEY TAKEAWAY:
//   - Environment variables are late-bound process configuration.
//   - Missing or malformed values should fail fast.
//   - Keep names stable and documented.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== CFG.1 Environment variables ===")
	fmt.Println("Learn how environment variables shape runtime configuration without rebuilding the binary.")
	fmt.Println()
	fmt.Println("- Environment variables are late-bound process configuration.")
	fmt.Println("- Missing or malformed values should fail fast.")
	fmt.Println("- Keep names stable and documented.")
	fmt.Println()
	fmt.Println("Environment variables are simple and ubiquitous, but they need documentation and validation to stay reliable.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CFG.2")
	fmt.Println("Current: CFG.1 (environment variables)")
	fmt.Println("---------------------------------------------------")
}
