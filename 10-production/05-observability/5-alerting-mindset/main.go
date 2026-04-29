// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: Alerting mindset
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how useful alerts differ from noisy alerts and why runbooks and service objectives matter.
//
// WHY THIS MATTERS:
//   - Alerts are requests for human attention, which means they are expensive by default.
//
// RUN:
//   go run ./10-production/05-observability/5-alerting-mindset
//
// KEY TAKEAWAY:
//   - Alert on symptoms that matter to users or service health.
//   - Tie alerts to runbooks or clear first actions.
//   - Use service objectives to decide what deserves a page.
// ============================================================================

package main

import "fmt"

//

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

// ---------------------------------------------------
// NEXT UP: DOCKER.1
// ---------------------------------------------------
