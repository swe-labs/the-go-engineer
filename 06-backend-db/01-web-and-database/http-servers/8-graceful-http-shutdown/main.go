package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - Graceful HTTP shutdown
//
// Run: go run ./06-backend-db/01-web-and-database/http-servers/8-graceful-http-shutdown

func main() {
	fmt.Println("=== HS.8 Graceful HTTP shutdown ===")
	fmt.Println("Learn how to stop accepting new work while allowing in-flight requests to finish cleanly.")
	fmt.Println()
	fmt.Println("- Stop accepting new requests before terminating the process.")
	fmt.Println("- Use contexts and deadlines to bound the drain window.")
	fmt.Println("- Design handlers so cancellation and shutdown signals can be respected.")
	fmt.Println()
	fmt.Println("Rolling deploys, autoscaling, and restarts all depend on shutdown behavior that does not drop useful work unnecessarily.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: HS.9")
	fmt.Println("Current: HS.8 (graceful http shutdown)")
	fmt.Println("---------------------------------------------------")
}
