package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - REST vs gRPC - the trade-off
//
// Run: go run ./06-backend-db/01-web-and-database/apis/8-rest-vs-grpc

func main() {
	fmt.Println("=== API.8 REST vs gRPC - the trade-off ===")
	fmt.Println("Compare when REST is a better fit and when gRPC earns its extra contract machinery.")
	fmt.Println()
	fmt.Println("- REST is simpler for browsers and public integration.")
	fmt.Println("- gRPC is strong for typed internal service-to-service communication.")
	fmt.Println("- Pick the transport that matches clients and operational constraints.")
	fmt.Println()
	fmt.Println("Teams waste time when they choose gRPC for public human-facing APIs or JSON over HTTP for every high-volume internal service without thinking.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: API.9")
	fmt.Println("Current: API.8 (rest vs grpc - the trade-off)")
	fmt.Println("---------------------------------------------------")
}
