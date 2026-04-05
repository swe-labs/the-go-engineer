// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import "fmt"

// ============================================================================
// Section 24: errgroup & sync.Pool — URL Health Checker (Exercise Starter)
// Level: Advanced
// ============================================================================
//
// EXERCISE: Build a Concurrent URL Health Checker using errgroup
//
// REQUIREMENTS:
//  1. [ ] Define a `CheckResult` struct with URL, StatusCode, Latency, Error
//  2. [ ] Implement `checkURL(ctx context.Context, url string) CheckResult`
//         using http.NewRequestWithContext and a 5-second timeout
//  3. [ ] Use errgroup.WithContext to run all checks concurrently
//         with a maximum of 5 simultaneous requests (use SetLimit)
//  4. [ ] Cancel all remaining checks if any returns a non-2xx status
//  5. [ ] Collect and print results sorted by latency (fastest first)
//  6. [ ] Add a sync.Pool of http.Clients to avoid allocating a new client
//         per request (use a shared pool with Transport.MaxIdleConnsPerHost=10)
//
// HINTS:
//   - g, ctx := errgroup.WithContext(context.Background())
//   - g.SetLimit(5) caps concurrency
//   - results[i] is safe to write from goroutine i (unique index)
//   - sort.Slice(results, func(i,j int) bool { return results[i].Latency < results[j].Latency })
//   - Pool New: func() any { return &http.Client{Timeout: 5*time.Second} }
//
// URLS TO TEST:
//   var urls = []string{
//     "https://go.dev",
//     "https://pkg.go.dev",
//     "https://github.com",
//     "https://api.github.com",
//     "https://httpbin.org/status/200",
//     "https://httpbin.org/delay/1",
//   }
//
// RUN: go run ./24-errgroup-and-pools/5-url-checker-exercise/_starter
// SOLUTION: See the main.go in the parent directory (5-url-checker-exercise)
// ============================================================================

// TODO: Define CheckResult struct

// TODO: Declare a sync.Pool of *http.Client

// TODO: Implement checkURL(ctx context.Context, url string) CheckResult

// TODO: Implement main with errgroup.WithContext + SetLimit

func main() {
	fmt.Println("=== URL Health Checker Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your concurrent URL health checker!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
