// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: 12-Factor App principles
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn the configuration and deployment discipline behind the 12-Factor mindset.
//
// WHY THIS MATTERS:
//   - 12-Factor is a set of operational habits that keep apps portable and environment-aware.
//
// RUN:
//   go run ./10-production/04-configuration/4-twelve-factor-principles
//
// KEY TAKEAWAY:
//   - Store config in the environment, not in code.
//   - Keep build, release, and run concerns separate.
//   - Treat logs and backing services as platform concerns.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== CFG.4 12-Factor App principles ===")
	fmt.Println("Learn the configuration and deployment discipline behind the 12-Factor mindset.")
	fmt.Println()
	fmt.Println("- Store config in the environment, not in code.")
	fmt.Println("- Keep build, release, and run concerns separate.")
	fmt.Println("- Treat logs and backing services as platform concerns.")
	fmt.Println()
	fmt.Println("The most useful part of 12-Factor for Go teams is the config and deployment discipline, not ritual compliance with every slogan.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CFG.5 -> 10-production/04-configuration/5-config-validation-on-boot")
	fmt.Println("Current: CFG.4 (12-factor app principles)")
	fmt.Println("---------------------------------------------------")
}
