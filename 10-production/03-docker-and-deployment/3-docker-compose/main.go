package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 10: Production Operations - Docker Compose
//
// Run: go run ./10-production/03-docker-and-deployment/3-docker-compose

func main() {
	fmt.Println("=== DOCKER.3 Docker Compose ===")
	fmt.Println("Learn how Compose coordinates multiple services, shared networks, and local environment defaults for development.")
	fmt.Println()
	fmt.Println("- Compose defines multi-container local environments.")
	fmt.Println("- Networks and named services remove manual wiring work.")
	fmt.Println("- Keep local orchestration close to how the system actually starts in practice.")
	fmt.Println()
	fmt.Println("Compose is helpful when a service is only meaningful next to its database, cache, or queue and local startup should be boring.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: DEPLOY.1")
	fmt.Println("Current: DOCKER.3 (docker compose)")
	fmt.Println("---------------------------------------------------")
}
