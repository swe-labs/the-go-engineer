package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 06: Backend, APIs & Databases - Health and readiness probes
//
// Run: go run ./06-backend-db/01-web-and-database/http-servers/9-health-and-readiness-probes

func main() {
	fmt.Println("=== HS.9 Health and readiness probes ===")
	fmt.Println("Learn why liveness and readiness are different signals and why each needs a clear contract.")
	fmt.Println()
	fmt.Println("- Keep health endpoints simple and deterministic.")
	fmt.Println("- Readiness should reflect dependencies required for safe traffic.")
	fmt.Println("- Do not overload probe endpoints with expensive checks.")
	fmt.Println()
	fmt.Println("A misleading readiness check can black-hole traffic or mask a degraded dependency long enough to cause a larger incident.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: HS.10")
	fmt.Println("Current: HS.9 (health and readiness probes)")
	fmt.Println("---------------------------------------------------")
}
