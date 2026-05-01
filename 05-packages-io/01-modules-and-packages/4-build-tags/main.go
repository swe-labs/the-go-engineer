// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section 05: Packages and I/O
// Title: Build Tags
// Level: Stretch
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - How to use build constraints to conditionally compile code and separate
//     platform-specific or environment-specific logic.
//
// WHY THIS MATTERS:
//   - Build tags allow you to manage platform-specific code (Windows vs Linux)
//     or isolate slow integration tests from the fast unit-test suite.
//
// RUN:
//   go run ./05-packages-io/01-modules-and-packages/4-build-tags
//
// KEY TAKEAWAY:
//   - Use '//go:build' to control which files are included in your binary.
// ============================================================================

// Commercial use is prohibited without permission.

// Try also: go test -v -tags=integration ./05-packages-io/01-modules-and-packages/4-build-tags
package main

import (
	"fmt"
	"runtime"
)

// Stage 05: Modules & Packages - Build Tags
//
//   - How `//go:build` constraints conditionally compile code.
//   - How to separate platform-specific logic (Windows vs Linux).
//   - How to isolate slow integration tests behind custom tags.
//
// ENGINEERING DEPTH:
//   Go has native support for conditional compilation via comments!
//   A comment like `//go:build linux` at the top of a file tells the Go compiler
//   "only include this file if the target OS is Linux."
//
//   Custom tags are heavily used for testing. `//go:build integration` ensures
//   that a test file isn't executed during a normal `go test ./...` unless the
//   developer explicitly opted in using `go test -tags=integration`. This keeps
//   local testing fast by default while allowing CI to run the heavy tests.

func main() {
	fmt.Println("=== Build Tags Demonstration ===")
	fmt.Printf("Compiled for: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// GetSystemDetails is implemented differently depending on the OS!
	// Look at os_windows.go and os_unix.go. The compiler automatically
	// selects the correct file based on the implicit build tags of the target OS.
	fmt.Println(GetSystemDetails())

	fmt.Println("\n---------------------------------------------------")
	fmt.Println("NEXT UP: CL.1 -> 05-packages-io/02-io-and-cli/cli-tools/1-args")
	fmt.Println("Current: MP.4 (build-tags)")
	fmt.Println("Previous: MP.3 (versioning-workshop)")
	fmt.Println("---------------------------------------------------")
}
