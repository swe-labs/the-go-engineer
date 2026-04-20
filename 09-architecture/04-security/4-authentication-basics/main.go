package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 09: Architecture & Security - Authentication basics
//
// Run: go run ./09-architecture/04-security/4-authentication-basics

func main() {
	fmt.Println("=== SEC.4 Authentication basics ===")
	fmt.Println("Learn the core differences between proving identity with sessions, tokens, and surrounding account rules.")
	fmt.Println()
	fmt.Println("- Authentication establishes identity.")
	fmt.Println("- Authorization should use that identity explicitly.")
	fmt.Println("- Session-based and token-based systems trade different operational costs.")
	fmt.Println()
	fmt.Println("Auth systems fail when identity, session storage, and permission checks blur together without clear ownership.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.5")
	fmt.Println("Current: SEC.4 (authentication basics)")
	fmt.Println("---------------------------------------------------")
}
