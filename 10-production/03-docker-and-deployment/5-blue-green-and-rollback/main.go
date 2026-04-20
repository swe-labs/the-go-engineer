package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 10: Production Operations - Blue/green and rollback
//
// Run: go run ./10-production/03-docker-and-deployment/5-blue-green-and-rollback

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
