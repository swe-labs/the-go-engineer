// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: Prometheus integration
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn the scrape-based model behind Prometheus and how application metrics become time series.
//
// WHY THIS MATTERS:
//   - Prometheus pulls metrics from your service on a regular interval instead of waiting for the service to push them.
//
// RUN:
//   go run ./10-production/05-observability/2-prometheus-integration
//
// KEY TAKEAWAY:
//   - Prometheus pulls metrics on a regular interval (scrape model).
//   - Expose a scrape-friendly /metrics endpoint.
//   - Label discipline keeps data bounded.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== OPS.2 Prometheus integration ===")
	fmt.Println("Learn the scrape-based model behind Prometheus and how application metrics become time series.")
	fmt.Println()
	fmt.Println("- Expose a scrape-friendly metrics endpoint.")
	fmt.Println("- Choose labels that stay bounded over time.")
	fmt.Println("- Metric names should describe both unit and domain.")
	fmt.Println()
	fmt.Println("Prometheus is simple to adopt, but label discipline and bucket design matter far more than whether the scrape endpoint exists.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: OPS.3")
	fmt.Println("Current: OPS.2 (prometheus integration)")
	fmt.Println("---------------------------------------------------")
}
