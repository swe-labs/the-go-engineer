// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: Blue/green and rollback
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn rollout strategies that reduce downtime and give operators a clear path back when a release is bad.
//
// WHY THIS MATTERS:
//   - A deployment strategy is really a risk-management strategy for switching traffic between application versions.
//
// RUN:
//   go run ./10-production/03-docker-and-deployment/5-blue-green-and-rollback
//
// KEY TAKEAWAY:
//   - Blue/green: run two versions, switch traffic atomically.
//   - Rollback: switch back to previous version when release is bad.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== DEPLOY.2 Blue/green and rollback ===")
	fmt.Println("Learn rollout strategies that reduce downtime and give operators a clear path back when a release is bad.")
	fmt.Println()
	fmt.Println("- Traffic switching is safer when the old version is still available.")
	fmt.Println("- Health checks and drain windows shape safe cutovers.")
	fmt.Println("- Rollback is a first-class deployment step, not a wish.")
	fmt.Println()
	fmt.Println("Zero-downtime claims are only meaningful when rollback, health checks, and drain behavior are planned together.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: DEPLOY.3")
	fmt.Println("Current: DEPLOY.2 (blue/green and rollback)")
	fmt.Println("---------------------------------------------------")
}
