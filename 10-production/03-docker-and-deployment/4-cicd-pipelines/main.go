// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: CI/CD pipelines
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how automated build, test, package, and deploy stages turn repository changes into controlled releases.
//
// WHY THIS MATTERS:
//   - A pipeline is a repeatable release process written down as automation instead of tribal knowledge.
//
// RUN:
//   go run ./10-production/03-docker-and-deployment/4-cicd-pipelines
//
// KEY TAKEAWAY:
//   - Pipeline stages: build -> test -> package -> deploy.
//   - Each stage must be repeatable and automatable.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== DEPLOY.1 CI/CD pipelines ===")
	fmt.Println("Learn how automated build, test, package, and deploy stages turn repository changes into controlled releases.")
	fmt.Println()
	fmt.Println("- Pipelines make release steps repeatable.")
	fmt.Println("- Quality gates should fail early and clearly.")
	fmt.Println("- Artifact creation and deployment should be explicit stages.")
	fmt.Println()
	fmt.Println("CI/CD is valuable because it removes hidden release steps and makes quality gates visible before production.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: DEPLOY.2 -> 10-production/03-docker-and-deployment/5-blue-green-and-rollback")
	fmt.Println("Current: DEPLOY.1 (ci/cd pipelines)")
	fmt.Println("---------------------------------------------------")
}
