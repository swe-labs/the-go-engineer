// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: Distributed tracing basics
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn how trace and span IDs follow a request through multiple boundaries so latency can be explained, not guessed.
//
// WHY THIS MATTERS:
//   - Tracing is request storytelling with timestamps and parent-child relationships.
//
// RUN:
//   go run ./10-production/05-observability/3-distributed-tracing-basics
//
// KEY TAKEAWAY:
//   - Trace shows one request across many services.
//   - Span explains where latency was spent inside that request.
//   - Context propagation carries trace ID through boundaries.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== OPS.3 Distributed tracing basics ===")
	fmt.Println("Learn how trace and span IDs follow a request through multiple boundaries so latency can be explained, not guessed.")
	fmt.Println()
	fmt.Println("- Traces show one request across many services.")
	fmt.Println("- Spans explain where latency was spent inside that request.")
	fmt.Println("- Context propagation is the transport for trace identity.")
	fmt.Println()
	fmt.Println("Tracing is expensive enough that you should know which paths matter and how sampling changes what you see.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: OPS.4")
	fmt.Println("Current: OPS.3 (distributed tracing basics)")
	fmt.Println("---------------------------------------------------")
}
