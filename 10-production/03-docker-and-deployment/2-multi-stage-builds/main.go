// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: Multi-stage builds
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn why Go binaries fit naturally into multi-stage builds that separate compile tools from runtime images.
//
// WHY THIS MATTERS:
//   - One stage builds the binary; a later stage runs only the binary and minimal runtime assets.
//
// RUN:
//   go run ./10-production/03-docker-and-deployment/2-multi-stage-builds
//
// KEY TAKEAWAY:
//   - One stage builds the binary; a later stage runs only the binary.
//   - Multi-stage builds keep images small and build tools out of runtime.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== DOCKER.2 Multi-stage builds ===")
	fmt.Println("Learn why Go binaries fit naturally into multi-stage builds that separate compile tools from runtime images.")
	fmt.Println()
	fmt.Println("- Use a builder image for compilation.")
	fmt.Println("- Copy only the resulting binary into the runtime image.")
	fmt.Println("- Smaller runtime images usually mean fewer attack surface and faster pulls.")
	fmt.Println()
	fmt.Println("Multi-stage builds reduce both artifact size and the number of unnecessary tools shipped into production.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: DOCKER.3")
	fmt.Println("Current: DOCKER.2 (multi-stage builds)")
	fmt.Println("---------------------------------------------------")
}
