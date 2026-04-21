// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import "fmt"

// ============================================================================
// Stage 07: Concurrency — Concurrent Downloader (Exercise Starter)
// Level: Advanced
// ============================================================================
//
// EXERCISE: Build a Concurrent Multi-File Downloader
//
// REQUIREMENTS:
//  1. [ ] Define a `Download` struct with URL, Filename, Size fields
//  2. [ ] Implement `downloadFile(url string) ([]byte, error)` using http.Get
//  3. [ ] Use a `sync.WaitGroup` to wait for all downloads to complete
//  4. [ ] Use a channel-based semaphore (`chan struct{}`) to limit concurrency
//         to 3 simultaneous downloads (prevent OS file descriptor exhaustion)
//  5. [ ] Save each downloaded file to disk using os.WriteFile
//  6. [ ] Print progress for each download (started, completed, failed)
//
// HINTS:
//   - Semaphore pattern: sem := make(chan struct{}, 3) — acquire: sem <- struct{}{}
//   - Always defer wg.Done() and <-sem inside the goroutine
//   - Capture loop variables correctly in closures: use `url := url`
//   - Use defer resp.Body.Close() after every http.Get
//
// RUN: go run ./07-concurrency/01-concurrency/goroutines/7-downloader/_starter
// SOLUTION: See the main.go file in the parent directory
// ============================================================================

// TODO: Define your Download struct

// TODO: Implement downloadFile function

// TODO: Implement the concurrent download orchestrator

func main() {
	fmt.Println("=== Concurrent Downloader Exercise ===")
	fmt.Println()
	fmt.Println("TODO: Implement your concurrent downloader!")
	fmt.Println("See the REQUIREMENTS above for what to build.")
	fmt.Println()
	fmt.Println("When finished, compare your solution with ../main.go")
}
