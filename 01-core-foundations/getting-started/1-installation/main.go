// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"runtime"
)

// ============================================================================
// Section 0: Getting Started — Installation Verification
// Level: Beginner
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to verify that Go is installed correctly
//   - How to inspect your Go environment
//   - What GOROOT, GOPATH, and GOOS mean
//
// RUN: go run ./00-getting-started/1-installation
// ============================================================================

func main() {
	// If you can run this program, Go is installed correctly!
	fmt.Println("✅ Go is installed and working!")
	fmt.Println()

	// runtime.Version() returns the Go version this binary was built with.
	// This should match what "go version" prints in your terminal.
	fmt.Printf("Go Version:    %s\n", runtime.Version())

	// runtime.GOOS returns the operating system target.
	// Common values: "linux", "darwin" (macOS), "windows"
	fmt.Printf("OS:            %s\n", runtime.GOOS)

	// runtime.GOARCH returns the CPU architecture target.
	// Common values: "amd64" (Intel/AMD 64-bit), "arm64" (Apple M-series, ARM servers)
	fmt.Printf("Architecture:  %s\n", runtime.GOARCH)

	// runtime.NumCPU() returns the number of logical CPUs available.
	// Go's concurrency scheduler (Section 09) uses these to run goroutines in parallel.
	fmt.Printf("CPUs:          %d\n", runtime.NumCPU())

	// runtime.GOROOT() returns where Go is installed on your system.
	// You rarely need to change this — it's set during installation.
	//lint:ignore SA1019 runtime.GOROOT is deprecated in Go 1.24 but helpful for beginners to see.
	fmt.Printf("GOROOT:        %s\n", runtime.GOROOT())

	fmt.Println()
	fmt.Println("🎉 Your Go development environment is ready to hack!")
	fmt.Println("   Next step: go run ./00-getting-started/2-hello-world")
}
