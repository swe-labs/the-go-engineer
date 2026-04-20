package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 10: Production Operations - Feature flags
//
// Run: go run ./10-production/05-observability/4-feature-flags

func main() {
	fmt.Println("=== OPS.4 Feature flags ===")
	fmt.Println("Learn how flag-controlled rollout changes the deployment story by separating shipping from activating behavior.")
	fmt.Println()
	fmt.Println("- Flags separate rollout from deployment.")
	fmt.Println("- Targeting rules should be simple and inspectable.")
	fmt.Println("- Remove flags once the rollout decision is finished.")
	fmt.Println()
	fmt.Println("Flags should have owners, expiry expectations, and observability. Permanent forgotten flags become dead architecture.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: OPS.5")
	fmt.Println("Current: OPS.4 (feature flags)")
	fmt.Println("---------------------------------------------------")
}
