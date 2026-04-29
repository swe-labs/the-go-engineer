// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: Service layer pattern
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how service layers coordinate domain operations, error handling, and side effects above repositories.
//
// WHY THIS MATTERS:
//   - A service layer owns use-case orchestration across multiple dependencies.
//
// RUN:
//   go run ./09-architecture/03-architecture-patterns/5-service-layer-pattern
//
// KEY TAKEAWAY:
//   - Learn how service layers coordinate domain operations, error handling, and side effects above repositories.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== ARCH.5 Service layer pattern ===")
	fmt.Println("Learn how service layers coordinate domain operations, error handling, and side effects above repositories.")
	fmt.Println()
	fmt.Println("- Keep handlers thin and repositories storage-focused.")
	fmt.Println("- Services coordinate domain behavior across dependencies.")
	fmt.Println("- Error classification becomes clearer when service methods own the use-case boundary.")
	fmt.Println()
	fmt.Println("Service layers are where cross-entity rules, retries, and error classification often belong.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ARCH.6")
	fmt.Println("Current: ARCH.5 (service layer pattern)")
	fmt.Println("---------------------------------------------------")
}
