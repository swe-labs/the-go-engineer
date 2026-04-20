package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - gRPC Service
//
// Run: go run ./06-backend-db/01-web-and-database/apis/9-grpc-service-exercise

func main() {
	fmt.Println("=== API.9 gRPC Service ===")
	fmt.Println("Build a small service contract that combines protobuf design, unary calls, and boundary behavior into one exercise.")
	fmt.Println()
	fmt.Println("- Define a service contract that feels coherent to a caller.")
	fmt.Println("- Keep transport concerns explicit at the boundary.")
	fmt.Println("- Treat the schema as a shared artifact, not a generated side effect.")
	fmt.Println()
	fmt.Println("Exercise work is where schema design stops being abstract and becomes something clients must live with.")
}
