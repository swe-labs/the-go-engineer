package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - net/http basics
//
// Run: go run ./06-backend-db/01-web-and-database/http-servers/1-net-http-basics

func main() {
	fmt.Println("=== HS.1 net/http basics ===")
	fmt.Println("Learn the smallest useful shape of an HTTP server in Go.")
	fmt.Println()
	fmt.Println("- Handlers are ordinary functions or interface implementations.")
	fmt.Println("- ServeMux routes requests by pattern.")
	fmt.Println("- ListenAndServe owns the accept loop for incoming connections.")
	fmt.Println()
	fmt.Println("The stdlib server is enough for most internal services and a strong baseline for external APIs too.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: HS.2")
	fmt.Println("Current: HS.1 (net/http basics)")
	fmt.Println("---------------------------------------------------")
}
