package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - gRPC interceptors
//
// Run: go run ./06-backend-db/01-web-and-database/apis/7-grpc-interceptors

func main() {
	fmt.Println("=== API.7 gRPC interceptors ===")
	fmt.Println("Learn how interceptors provide middleware-style hooks for auth, logging, metrics, and recovery in gRPC.")
	fmt.Println()
	fmt.Println("- Interceptors wrap handlers before and after business logic runs.")
	fmt.Println("- Recovery, auth, and logging are common boundary concerns.")
	fmt.Println("- Ordering matters because outer interceptors see the whole call lifecycle.")
	fmt.Println()
	fmt.Println("Transport policy should live in interceptors or shared boundary helpers, not be copy-pasted into every RPC method.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: API.8")
	fmt.Println("Current: API.7 (grpc interceptors)")
	fmt.Println("---------------------------------------------------")
}
