// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 09: Architecture & Security
// Title: Event-driven architecture
// Level: Production
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn the benefits and failure modes of publishing domain events instead of wiring every side effect inline.
//
// WHY THIS MATTERS:
//   - Events decouple producers from consumers by letting one change be observed by many subscribers later.
//
// RUN:
//   go run ./09-architecture/03-architecture-patterns/6-event-driven-architecture
//
// KEY TAKEAWAY:
//   - Learn the benefits and failure modes of publishing domain events instead of wiring every side effect inline.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== ARCH.6 Event-driven architecture ===")
	fmt.Println("Learn the benefits and failure modes of publishing domain events instead of wiring every side effect inline.")
	fmt.Println()
	fmt.Println("- Events describe something that already happened.")
	fmt.Println("- Publishers should not need to know every consumer.")
	fmt.Println("- Delivery guarantees and idempotency matter more than the event bus brand.")
	fmt.Println()
	fmt.Println("Event-driven designs help when many downstream reactions exist, but they demand discipline around idempotency and debugging.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: ARCH.7")
	fmt.Println("Current: ARCH.6 (event-driven architecture)")
	fmt.Println("---------------------------------------------------")
}
