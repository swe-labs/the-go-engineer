// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 01: Getting Started
// Title: Installation Verification
// Level: Foundation
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Confirm that Go is installed and that this machine can run a real Go program from this repo.
//
// WHY THIS MATTERS:
//   - `go run` is a short pipeline: 1. Read the source files. 2. Compile them into a runnable program. 3. Execute that program. 4. Show the output in the...
//
// RUN:
//   go run ./01-getting-started/1-installation
//
// KEY TAKEAWAY:
//   - Confirm that Go is installed and that this machine can run a real Go program from this repo.
// ============================================================================

package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Go installation looks healthy.")
	fmt.Printf("Go version:   %s\n", runtime.Version())
	fmt.Printf("OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Logical CPUs: %d\n", runtime.NumCPU())
	fmt.Println()
	fmt.Println("---------------------------------------------------")
	fmt.Println("NEXT UP: GT.2 -> 01-getting-started/2-hello-world")
	fmt.Println("Current: GT.1 (installation)")
	fmt.Println("---------------------------------------------------")
}
