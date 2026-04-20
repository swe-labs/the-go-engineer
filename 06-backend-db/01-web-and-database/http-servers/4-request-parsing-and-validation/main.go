package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - Request parsing and validation
//
// Run: go run ./06-backend-db/01-web-and-database/http-servers/4-request-parsing-and-validation

func main() {
	fmt.Println("=== HS.4 Request parsing and validation ===")
	fmt.Println("Learn how to decode request data deliberately and reject malformed input early.")
	fmt.Println()
	fmt.Println("- Parse one request surface at a time.")
	fmt.Println("- Limit body size before decoding it.")
	fmt.Println("- Return validation feedback before business logic runs.")
	fmt.Println()
	fmt.Println("Most API bugs start at the boundary. Size limits, decode errors, and validation checks should fail fast.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: HS.5")
	fmt.Println("Current: HS.4 (request parsing and validation)")
	fmt.Println("---------------------------------------------------")
}
