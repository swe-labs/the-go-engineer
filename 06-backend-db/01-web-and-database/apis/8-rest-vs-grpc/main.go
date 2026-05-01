// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: REST vs gRPC - The Trade-off
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - A side-by-side comparison of REST and gRPC.
//   - The technical trade-offs (Payload, Latency, Tooling).
//   - How to choose the right transport for your specific use case.
//
// WHY THIS MATTERS:
//   - There is no "Best" protocol. An engineer's job is to choose the right
//     tool for the constraints of the project.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/apis/8-rest-vs-grpc
//
// KEY TAKEAWAY:
//   - REST is for Public Ease. gRPC is for Internal Speed.
// ============================================================================

package main

import "fmt"

// Stage 06: APIs - REST vs gRPC
//
//   - REST: Resource-oriented, JSON, HTTP/1.1, Browser-native
//   - gRPC: Action-oriented, Protobuf, HTTP/2, Machine-optimized
//
// ENGINEERING DEPTH:
//   The choice often comes down to the "Cost of Integration".
//   If your API is consumed by random developers on the web, REST
//   is superior because everyone has a browser and a JSON parser.
//   If your API is consumed by your own microservices in a data
//   center, gRPC is superior because it maximizes hardware
//   efficiency and provides compile-time safety.

func main() {
	fmt.Println("=== REST vs gRPC: The Decision Matrix ===")
	fmt.Println()

	fmt.Println("  FEATURE          | REST (JSON)         | gRPC (Protobuf)")
	fmt.Println("  -----------------|---------------------|-------------------")
	fmt.Println("  Payload          | Text (Large)        | Binary (Small)")
	fmt.Println("  Transport        | HTTP/1.1 or 2       | HTTP/2 Only")
	fmt.Println("  Contract         | Loose (Doc-based)   | Strict (.proto)")
	fmt.Println("  Browser Support  | Excellent (Native)  | Limited (Requires Proxy)")
	fmt.Println("  Streaming        | Unidirectional      | Bidirectional")
	fmt.Println("  Performance      | Medium              | Very High")
	fmt.Println()

	fmt.Println("  WHICH ONE SHOULD YOU CHOOSE?")
	fmt.Println()
	fmt.Println("  ✅ Choose REST if:")
	fmt.Println("     - Building a public API for third-party developers.")
	fmt.Println("     - Primary consumers are web browsers.")
	fmt.Println("     - You need easy debugging with simple tools like 'curl' or 'Postman'.")
	fmt.Println()
	fmt.Println("  ✅ Choose gRPC if:")
	fmt.Println("     - Building internal microservices (Service-to-Service).")
	fmt.Println("     - You need the absolute highest performance and lowest latency.")
	fmt.Println("     - You have a large team and need strict type safety across languages.")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: API.9 -> 06-backend-db/01-web-and-database/apis/9-grpc-service-exercise")
	fmt.Println("Current: API.8 (rest-vs-grpc)")
	fmt.Println("Previous: API.7 (grpc-interceptors)")
	fmt.Println("---------------------------------------------------")
}
