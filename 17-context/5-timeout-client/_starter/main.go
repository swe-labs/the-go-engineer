// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Section 17: Context — Timeout-Aware API Client (Exercise Starter)
// Level: Advanced
// ============================================================================
//
// EXERCISE: Build a Timeout-Aware HTTP Client
//
// REQUIREMENTS:
//  1. [ ] Implement `fetchWithTimeout(url string, timeout time.Duration) (string, error)`
//  2. [ ] Use context.WithTimeout to create a cancellable context
//  3. [ ] Use http.NewRequestWithContext to attach the context to the request
//  4. [ ] If the timeout fires, return a clear error message
//  5. [ ] Test with a fast URL (should succeed) and an impossibly short timeout (should fail)
//
// RUN: go run ./17-context/5-timeout-client/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Implement fetchWithTimeout

func main() {
	fmt.Println("=== Timeout-Aware API Client Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your timeout-aware HTTP client!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
