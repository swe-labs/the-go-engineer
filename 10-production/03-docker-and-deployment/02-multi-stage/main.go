package main

/*
Multi-Stage Dockerfile

This lesson demonstrates the multi-stage build pattern.
Instead of one large image with compiler tools, we use two stages:
1. Build stage with golang:alpine (uses compiler)
2. Runtime stage with alpine (only the binary)

Result: Final image is  ~15MB instead of ~350MB!
"""

// See the Dockerfile in this directory for implementation
//
// RUN: docker build -t myapp:latest .
// RUN: docker run myapp:latest
*/

func main() {
	// This lesson demonstrates multi-stage Docker builds.
	// See the Dockerfile in this directory for the implementation.
	// Run: docker build -t myapp:latest .
	// Then: docker run myapp:latest
	println("Multi-stage Dockerfile lesson - see Dockerfile in this directory")
}
