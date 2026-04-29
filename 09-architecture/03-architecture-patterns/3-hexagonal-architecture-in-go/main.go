// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: Hexagonal architecture in Go
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how ports and adapters keep domain rules independent of transport and storage details.
//
// WHY THIS MATTERS:
//   - Hexagonal architecture isolates the core from delivery and persistence concerns through explicit ports.
//
// RUN:
//   go run ./09-architecture/03-architecture-patterns/3-hexagonal-architecture-in-go
//
// KEY TAKEAWAY:
//   - Learn how ports and adapters keep domain rules independent of transport and storage details.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== ARCH.3 Hexagonal architecture in Go ===")
	fmt.Println("Learn how ports and adapters keep domain rules independent of transport and storage details.")
	fmt.Println()
	fmt.Println("- Ports describe what the core needs from the outside world.")
	fmt.Println("- Adapters translate external details into the port contract.")
	fmt.Println("- The domain should not know whether it is called by HTTP, gRPC, or a job runner.")
	fmt.Println()
	fmt.Println("The value of hexagonal design is not the shape of the diagram. It is the ability to change boundaries without rewriting domain logic.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ARCH.4")
	fmt.Println("Current: ARCH.3 (hexagonal architecture in go)")
	fmt.Println("---------------------------------------------------")
}
