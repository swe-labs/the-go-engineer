package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - Middleware - the pattern
//
// Run: go run ./06-backend-db/01-web-and-database/http-servers/3-middleware-pattern

func main() {
	fmt.Println("=== HS.3 Middleware - the pattern ===")
	fmt.Println("Learn how middleware wraps handlers to add logging, auth, recovery, and cross-cutting rules.")
	fmt.Println()
	fmt.Println("- Middleware keeps cross-cutting behavior out of endpoint handlers.")
	fmt.Println("- Ordering matters because outer layers run before and after inner ones.")
	fmt.Println("- Recovery middleware is a boundary that turns panics into controlled responses.")
	fmt.Println()
	fmt.Println("Keep middleware small and composable so one stack does not become a hidden framework.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: HS.4")
	fmt.Println("Current: HS.3 (middleware - the pattern)")
	fmt.Println("---------------------------------------------------")
}
