package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 09: Architecture & Security - Rate limiting patterns
//
// Run: go run ./09-architecture/04-security/7-rate-limiting-patterns

func main() {
	fmt.Println("=== SEC.7 Rate limiting patterns ===")
	fmt.Println("Learn how token bucket and similar controls protect shared resources from abuse and accidental overload.")
	fmt.Println()
	fmt.Println("- Token buckets handle bursts better than flat per-second counters.")
	fmt.Println("- Choose rate keys carefully: IP, user, API key, or tenant.")
	fmt.Println("- Limiters need observability so operators can tell abuse from bad defaults.")
	fmt.Println()
	fmt.Println("Rate limits protect both your service and the dependencies behind it, especially on auth, search, and write-heavy endpoints.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.8")
	fmt.Println("Current: SEC.7 (rate limiting patterns)")
	fmt.Println("---------------------------------------------------")
}
