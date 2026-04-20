package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 10: Production Operations - Alerting mindset
//
// Run: go run ./10-production/05-observability/5-alerting-mindset

func main() {
	fmt.Println("=== OPS.5 Alerting mindset ===")
	fmt.Println("Learn how useful alerts differ from noisy alerts and why runbooks and service objectives matter.")
	fmt.Println()
	fmt.Println("- Alert on symptoms that matter to users or service health.")
	fmt.Println("- Tie alerts to runbooks or clear first actions.")
	fmt.Println("- Use service objectives to decide what deserves a page.")
	fmt.Println()
	fmt.Println("A noisy alert stream teaches engineers to ignore signals, which is the opposite of observability.")
}
