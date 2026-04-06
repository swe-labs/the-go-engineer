// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./08-modules-and-packages/4-build-tags
// Try also: go test -v -tags=integration ./08-modules-and-packages/4-build-tags
package main

import (
	"fmt"
	"runtime"
)

// ============================================================================
// Section 08: Modules & Packages — Build Tags
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
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
// ============================================================================

func main() {
	fmt.Println("=== Build Tags Demonstration ===")
	fmt.Printf("Compiled for: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// GetSystemDetails is implemented differently depending on the OS!
	// Look at os_windows.go and os_unix.go. The compiler automatically
	// selects the correct file based on the implicit build tags of the target OS.
	fmt.Println(GetSystemDetails())
}
