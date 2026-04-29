// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: gRPC Streaming
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to define Server-side, Client-side, and Bidirectional streams.
//   - Real-world use cases for gRPC streaming (notifications, logs, metrics).
//   - How the 'stream' keyword changes the behavior of an RPC.
//
// WHY THIS MATTERS:
//   - Streaming allows you to handle massive amounts of data or real-time
//     updates without the overhead of opening new requests for every message.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/apis/6-grpc-streaming
//
// KEY TAKEAWAY:
//   - Streams turn a single Request/Response into a persistent conversation.
// ============================================================================

package main

import "fmt"

// Stage 06: APIs - gRPC Streaming
//
//   - Server Stream: One request, multiple responses (e.g., Live Logs)
//   - Client Stream: Multiple requests, one response (e.g., Bulk Upload)
//   - Bidirectional Stream: Real-time conversation (e.g., Chat)
//
// ENGINEERING DEPTH:
//   In standard HTTP/1.1, "streaming" is achieved via long-polling or
//   chunked encoding. In gRPC (HTTP/2), streams are first-class citizens.
//   The client and server can send messages whenever they want over the
//   same TCP connection. This reduces latency by eliminating the
//   "Handshake" overhead of new requests.

func main() {
	fmt.Println("=== gRPC Streaming (Conceptual) ===")
	fmt.Println()

	fmt.Println("  1. SERVER STREAMING")
	fmt.Println("     rpc GetLogs(LogRequest) returns (stream LogEntry);")
	fmt.Println("     Use Case: Real-time monitoring, Stock tickers.")
	fmt.Println()

	fmt.Println("  2. CLIENT STREAMING")
	fmt.Println("     rpc UploadMetrics(stream Metric) returns (UploadSummary);")
	fmt.Println("     Use Case: Batching data points, File uploads.")
	fmt.Println()

	fmt.Println("  3. BIDIRECTIONAL STREAMING")
	fmt.Println("     rpc Chat(stream ChatMsg) returns (stream ChatMsg);")
	fmt.Println("     Use Case: Multiplayer games, Collaborative editing.")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: API.7 grpc-interceptors")
	fmt.Println("Current: API.6 (grpc-streaming)")
	fmt.Println("Previous: API.5 (grpc-fundamentals)")
	fmt.Println("---------------------------------------------------")
}
