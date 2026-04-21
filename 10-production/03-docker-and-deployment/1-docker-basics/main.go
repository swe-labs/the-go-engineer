package main

import "fmt"

// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// Section 10: Production Operations - Docker basics
//
// Run: go run ./10-production/03-docker-and-deployment/1-docker-basics

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
	fmt.Println("NEXT UP: DOCKER.2")
	fmt.Println("Current: DOCKER.1 (docker basics)")
	fmt.Println("---------------------------------------------------")
}
