package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 09: Architecture & Security - TLS and HTTPS in Go
//
// Run: go run ./09-architecture/04-security/8-tls-and-https-in-go

func main() {
	fmt.Println("=== SEC.8 TLS and HTTPS in Go ===")
	fmt.Println("Learn the transport-level rules that turn plain HTTP into encrypted, identity-checked HTTPS.")
	fmt.Println()
	fmt.Println("- HTTPS is HTTP over a verified TLS channel.")
	fmt.Println("- Certificate validation is part of the trust model, not an optional extra.")
	fmt.Println("- Reasonable defaults are safer than hand-rolling exotic TLS settings.")
	fmt.Println()
	fmt.Println("Transport security is not optional on hostile networks, and misconfiguration can silently remove the safety you thought you had.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: SEC.9")
	fmt.Println("Current: SEC.8 (tls and https in go)")
	fmt.Println("---------------------------------------------------")
}
