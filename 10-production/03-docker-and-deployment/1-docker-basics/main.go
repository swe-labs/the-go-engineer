// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: Docker basics
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Learn the basic building blocks of images, containers, layers, and Dockerfiles.
//
// WHY THIS MATTERS:
//   - A Docker image is a packaged filesystem and startup contract; a container is one running instance of that package.
//
// RUN:
//   go run ./10-production/03-docker-and-deployment/1-docker-basics
//
// KEY TAKEAWAY:
//   - Images are build outputs; containers are running instances.
//   - Dockerfiles describe how the image is assembled.
//   - Layer order influences rebuild speed and image cleanliness.
// ============================================================================

package main

import "fmt"

//

func main() {
	fmt.Println("=== DOCKER.1 Docker basics ===")
	fmt.Println("Learn the basic building blocks of images, containers, layers, and Dockerfiles.")
	fmt.Println()
	fmt.Println("- Images are build outputs; containers are running instances.")
	fmt.Println("- Dockerfiles describe how the image is assembled.")
	fmt.Println("- Layer order influences rebuild speed and image cleanliness.")
	fmt.Println()
	fmt.Println("Container basics matter because deployment packaging changes startup, debugging, and runtime assumptions.")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: DOCKER.2 -> 10-production/03-docker-and-deployment/2-multi-stage-builds")
	fmt.Println("Current: DOCKER.1 (docker basics)")
	fmt.Println("---------------------------------------------------")
}
