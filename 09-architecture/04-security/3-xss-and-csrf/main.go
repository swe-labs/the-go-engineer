package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 09: Architecture & Security - XSS and CSRF
//
// Run: go run ./09-architecture/04-security/3-xss-and-csrf

func main() {
	fmt.Println("=== SEC.3 XSS and CSRF ===")
	fmt.Println("Learn the difference between output-encoding bugs and cross-site request trust bugs in web systems.")
	fmt.Println()
	fmt.Println("- Escape untrusted output before it becomes HTML.")
	fmt.Println("- CSRF defenses prove the request came from the expected origin or flow.")
	fmt.Println("- Sessions and browser clients need different care than pure machine APIs.")
	fmt.Println()
	fmt.Println("Security bugs at the browser boundary depend on precise rules about output escaping, cookie scope, and request intent.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.4")
	fmt.Println("Current: SEC.3 (xss and csrf)")
	fmt.Println("---------------------------------------------------")
}
