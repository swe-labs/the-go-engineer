package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - Pagination and filtering
//
// Run: go run ./06-backend-db/01-web-and-database/apis/3-pagination-and-filtering

func main() {
	fmt.Println("=== API.3 Pagination and filtering ===")
	fmt.Println("Learn why large result sets need stable pagination and explicit filtering rules.")
	fmt.Println()
	fmt.Println("- Always define an ordering when paginating.")
	fmt.Println("- Offset pagination is simple but weak on large, changing datasets.")
	fmt.Println("- Cursor pagination is harder up front but often safer at scale.")
	fmt.Println()
	fmt.Println("Pagination bugs often become latency bugs, memory bugs, or duplicated work for clients.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: API.4")
	fmt.Println("Current: API.3 (pagination and filtering)")
	fmt.Println("---------------------------------------------------")
}
