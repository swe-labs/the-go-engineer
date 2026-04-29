// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: Metrics basics
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn what metrics answer that logs do not and why cardinality discipline matters early.
//
// WHY THIS MATTERS:
//   - Metrics are numeric summaries over time, which makes them good for trend and saturation questions.
//
// RUN:
//   go run ./10-production/05-observability/1-metrics-basics
//
// KEY TAKEAWAY:
//   - Counters answer "how much total work happened".
//   - Gauges answer "what value is true right now".
//   - Histograms answer "how work is distributed".
//   - Cardinality discipline keeps metrics useful.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== OPS.1 Metrics basics ===")
	fmt.Println("Learn what metrics answer that logs do not and why cardinality discipline matters early.")
	fmt.Println()
	fmt.Println("- Counters answer how much total work happened.")
	fmt.Println("- Gauges answer what value is true right now.")
	fmt.Println("- Histograms answer how work is distributed, not just how much happened.")
	fmt.Println()
	fmt.Println("Observability starts with choosing stable dimensions. Cardinality explosions make a correct metrics program operationally useless.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: OPS.2")
	fmt.Println("Current: OPS.1 (metrics basics)")
	fmt.Println("---------------------------------------------------")
}
