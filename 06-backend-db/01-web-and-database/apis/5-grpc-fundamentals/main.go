package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - gRPC fundamentals
//
// Run: go run ./06-backend-db/01-web-and-database/apis/5-grpc-fundamentals

func main() {
	fmt.Println("=== API.5 gRPC fundamentals ===")
	fmt.Println("Learn the core request/response model behind gRPC services and generated contracts.")
	fmt.Println()
	fmt.Println("- Services are defined in proto files and implemented in Go.")
	fmt.Println("- Unary RPCs look like strongly typed function calls across the network.")
	fmt.Println("- Status codes and metadata form the transport contract around the message payload.")
	fmt.Println()
	fmt.Println("gRPC shines most when typed contracts, internal services, and generated clients are a better fit than public HTTP ergonomics.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: API.6")
	fmt.Println("Current: API.5 (grpc fundamentals)")
	fmt.Println("---------------------------------------------------")
}
