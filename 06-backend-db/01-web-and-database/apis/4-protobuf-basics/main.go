// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 06: Backend, APIs & Databases
// Title: Protobuf Basics
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - The difference between IDL-based schemas and dynamic JSON.
//   - Why binary serialization is faster and smaller than text.
//   - How to read a '.proto' file.
//
// WHY THIS MATTERS:
//   - Protocol Buffers (Protobuf) is the industry standard for high-performance
//     microservices. It ensures that both the client and server agree on
//     the exact data shape before a single byte is sent.
//
// RUN:
//   go run ./06-backend-db/01-web-and-database/apis/4-protobuf-basics
//
// KEY TAKEAWAY:
//   - Protobuf is Schema-First. JSON is Data-First.
// ============================================================================

package main

import "fmt"

// Stage 06: APIs - Protobuf Basics
//
//   - The .proto file: Our source of truth
//   - Field Tags: 1, 2, 3 (The binary identity)
//   - Scalar Types: string, int32, bool
//
// ENGINEERING DEPTH:
//   In JSON, the field names ("first_name") are sent with every single
//   object. In Protobuf, the field names are replaced by small numbers
//   (tags). The receiver uses the shared schema to know that tag #1
//   means "first_name". This reduces payload size by up to 80% and
//   eliminates CPU-heavy string parsing.

func main() {
	fmt.Println("=== Protobuf Basics (Conceptual) ===")
	fmt.Println()

	fmt.Println("  1. THE SCHEMA (.proto)")
	fmt.Println("     message User {")
	fmt.Println("       int32 id = 1;")
	fmt.Println("       string email = 2;")
	fmt.Println("       bool is_active = 3;")
	fmt.Println("     }")
	fmt.Println()

	fmt.Println("  2. THE COMPARISON")
	fmt.Println("     JSON:     {\"id\": 1, \"email\": \"a@b.com\", \"is_active\": true}")
	fmt.Println("     PROTO:    [08 01 12 07 61 40 62 2E 63 6F 6D 18 01]")
	fmt.Println("               (Small, binary, and efficient!)")
	fmt.Println()

	fmt.Println("  3. KEY CONCEPTS")
	fmt.Println("     - Strictly Typed: No more 'maybe this is a string, maybe an int'.")
	fmt.Println("     - Forward/Backward Compatibility: Rules for adding/removing fields.")
	fmt.Println("     - Code Generation: Generate Go, Java, Python, and TS from one file.")

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: API.5 grpc-fundamentals")
	fmt.Println("Current: API.4 (protobuf-basics)")
	fmt.Println("Previous: API.3 (pagination-and-filtering)")
	fmt.Println("---------------------------------------------------")
}
