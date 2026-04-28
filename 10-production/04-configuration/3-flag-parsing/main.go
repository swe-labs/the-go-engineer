// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: Flag parsing
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn when startup flags are the right config surface and how they interact with file and env-based config.
//
// WHY THIS MATTERS:
//   - Flags are runtime arguments for one process start, which makes them great for overrides and local tooling.
//
// RUN:
//   go run ./10-production/04-configuration/3-flag-parsing
//
// KEY TAKEAWAY:
//   - Flags are explicit one-start overrides.
//   - Do not hide critical startup behavior behind undocumented flags.
//   - Define precedence clearly across flags, files, and environment.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== CFG.3 Flag parsing ===")
	fmt.Println("Learn when startup flags are the right config surface and how they interact with file and env-based config.")
	fmt.Println()
	fmt.Println("- Flags are explicit one-start overrides.")
	fmt.Println("- Do not hide critical startup behavior behind undocumented flags.")
	fmt.Println("- Define precedence clearly across flags, files, and environment.")
	fmt.Println()
	fmt.Println("Flags are best when the operator should see and choose the value at launch time instead of through ambient environment state.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CFG.4")
	fmt.Println("Current: CFG.3 (flag parsing)")
	fmt.Println("---------------------------------------------------")
}
