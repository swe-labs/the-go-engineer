// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 10: Production Operations
// Title: Docker basics
// Level: Core
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Build a Docker image from a Go application using a Dockerfile.
//   - Understand images, containers, layers, and the Docker build cache.
//
// WHY THIS MATTERS:
//   - A Docker image is a packaged filesystem and startup contract; a container
//     is one running instance of that package. Images are build outputs;
//     containers are running instances.
//
// RUN:
//   go run ./10-production/03-docker-and-deployment/01-docker-basics
//
// Build with Docker:
//   docker build -t docker-basics ./10-production/03-docker-and-deployment/01-docker-basics/
//   docker run -p 8080:8080 docker-basics
//
// KEY TAKEAWAY:
//   - Images are build outputs; containers are running instances.
//   - Dockerfiles describe how the image is assembled.
//   - Layer order influences rebuild speed and image cleanliness.
// ============================================================================

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Docker! Path: %s\n", r.URL.Path)
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "healthy")
	})

	port := ":8080"
	fmt.Printf("=== DOCKER.1 Docker basics ===\n")
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println()
	fmt.Println("Build with Docker:")
	fmt.Println("  docker build -t docker-basics ./10-production/03-docker-and-deployment/01-docker-basics/")
	fmt.Println("  docker run -p 8080:8080 docker-basics")
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: DOCKER.2 -> 10-production/03-docker-and-deployment/02-multi-stage-builds")
	fmt.Println("Current: DOCKER.1 (docker basics)")
	fmt.Println("---------------------------------------------------")

	log.Fatal(http.ListenAndServe(port, mux))
}
