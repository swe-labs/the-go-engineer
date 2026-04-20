package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - Response writing patterns
//
// Run: go run ./06-backend-db/01-web-and-database/http-servers/5-response-writing-patterns

func main() {
	fmt.Println("=== HS.5 Response writing patterns ===")
	fmt.Println("Learn how status codes, headers, bodies, and streaming responses form one response contract.")
	fmt.Println()
	fmt.Println("- Write headers intentionally instead of relying on defaults.")
	fmt.Println("- Keep success and error payload shapes predictable.")
	fmt.Println("- Streaming changes when and how bytes reach the client.")
	fmt.Println()
	fmt.Println("Response helpers reduce duplication and keep error cases consistent across the API.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: HS.6")
	fmt.Println("Current: HS.5 (response writing patterns)")
	fmt.Println("---------------------------------------------------")
}
