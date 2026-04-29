// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: Dockerised Service
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Package one service shape with config, container build, and rollout thinking into a single exercise surface.
//
// WHY THIS MATTERS:
//   - A deployment exercise proves that the code, container, and operating assumptions agree with each other.
//
// RUN:
//   go run ./10-production/03-docker-and-deployment/6-dockerised-service
//
// KEY TAKEAWAY:
//   - A deployment exercise proves code, container, and ops assumptions work together.
//   - Package config, build, and rollout into one coherent surface.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== DEPLOY.3 Dockerised Service ===")
	fmt.Println("Package one service shape with config, container build, and rollout thinking into a single exercise surface.")
	fmt.Println()
	fmt.Println("- Package a service and its runtime assumptions together.")
	fmt.Println("- Make startup, config, and health behavior explicit.")
	fmt.Println("- Treat the exercise as a delivery proof, not just a build artifact.")
	fmt.Println()
	fmt.Println("Production readiness is a property of the whole delivery path, not just the application source.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: CG.1")
	fmt.Println("   Current: DEPLOY.3 (dockerised service)")
	fmt.Println("---------------------------------------------------")
}
