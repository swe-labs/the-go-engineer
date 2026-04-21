package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 09: Architecture & Security - SQL injection prevention
//
// Run: go run ./09-architecture/04-security/2-sql-injection-prevention

func main() {
	fmt.Println("=== SEC.2 SQL injection prevention ===")
	fmt.Println("Learn why parameterized queries are the baseline defense against SQL injection in Go.")
	fmt.Println()
	fmt.Println("- Never concatenate untrusted values into SQL strings.")
	fmt.Println("- Use parameter placeholders and argument binding.")
	fmt.Println("- Validate input shape even when the query is parameterized.")
	fmt.Println()
	fmt.Println("String-building SQL with user input is a security bug even when it looks harmless in tests.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.3")
	fmt.Println("Current: SEC.2 (sql injection prevention)")
	fmt.Println("---------------------------------------------------")
}
