// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: gRPC Fundamentals
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The difference between REST and gRPC.
//   - How to define a service with Remote Procedure Calls (RPC).
//   - The benefits of HTTP/2 and multiplexing.
//
// WHY THIS MATTERS:
//   - gRPC allows you to call functions on a remote server as if they were
//     local functions, with full type safety and high performance.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/apis/5-grpc-fundamentals
//
// KEY TAKEAWAY:
//   - gRPC is for calling Functions. REST is for managing Resources.
// ============================================================================

package main

import "fmt"

// Stage 06: APIs - gRPC Fundamentals
//
//   - Unary RPC: One request, one response
//   - Service Definition: rpc GetUser(UserRequest) returns (UserResponse)
//   - HTTP/2: Binary framing and multiplexing
//
// ENGINEERING DEPTH:
//   gRPC uses HTTP/2 as its transport. Unlike HTTP/1.1, which opens a
//   new connection (or waits) for every request, HTTP/2 allows many
//   requests to fly over a single TCP connection at the same time.
//   Combined with Protobuf, this makes gRPC significantly lower latency
//   than traditional REST APIs.

func main() {
	fmt.Println("=== gRPC Fundamentals (Conceptual) ===")
	fmt.Println()

	fmt.Println("  1. THE SERVICE DEFINITION (.proto)")
	fmt.Println("     service UserService {")
	fmt.Println("       rpc GetUser(UserRequest) returns (UserResponse);")
	fmt.Println("       rpc CreateUser(CreateUserRequest) returns (User);")
	fmt.Println("     }")
	fmt.Println()

	fmt.Println("  2. THE FLOW")
	fmt.Println("     Client: calls generated 'GetUser()' function.")
	fmt.Println("     Stub:   serializes request to Protobuf bytes.")
	fmt.Println("     Trans:  sends bytes over HTTP/2 stream.")
	fmt.Println("     Server: deserializes and calls your implementation.")
	fmt.Println()

	fmt.Println("  3. WHY gRPC?")
	fmt.Println("     - Strongly Typed: No more 'is this 404 or 500?' ambiguity.")
	fmt.Println("     - Auto-Generated: No need to write 'json.Marshal' ever again.")
	fmt.Println("     - Streaming: Supports real-time data out of the box.")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: API.6 grpc-streaming")
	fmt.Println("Current: API.5 (grpc-fundamentals)")
	fmt.Println("Previous: API.4 (protobuf-basics)")
	fmt.Println("---------------------------------------------------")
}
