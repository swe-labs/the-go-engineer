// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ============================================================================
// Section 17: Context — Timeout-Aware API Client (Exercise)
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Attaching context.WithTimeout to HTTP requests
//   - How context cancellation propagates to the HTTP transport layer
//   - The production pattern: NEVER make I/O calls without a timeout
//   - Distinguishing between timeout errors and network errors
//
// ENGINEERING DEPTH:
//   When you create a request with `http.NewRequestWithContext(ctx, ...)`, the
//   HTTP client's internal transport layer monitors `ctx.Done()`. If the context
//   expires before the TCP connection completes, the transport immediately closes
//   the underlying socket and returns `context.DeadlineExceeded`. This prevents
//   your application from hanging indefinitely on a slow or unresponsive server —
//   a critical defense against cascading failures in microservice architectures.
//
// RUN: go run ./17-context/5-timeout-client
// ============================================================================

// fetchWithTimeout makes an HTTP GET request with a timeout context.
// If the request doesn't complete within the timeout, it's automatically cancelled.
func fetchWithTimeout(url string, timeout time.Duration) (string, error) {
	// Create a context with a timeout.
	// After `timeout` duration, the context automatically calls cancel().
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // Always call cancel to release resources!

	// Create the HTTP request WITH the context attached.
	// This is the critical step — linking the context to the HTTP client.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	// Execute the request. The HTTP client monitors ctx.Done() internally.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Check if the error was caused by our timeout
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("request timed out after %v: %w", timeout, err)
		}
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	return string(body), nil
}

func main() {
	fmt.Println("=== Timeout-Aware API Client ===")
	fmt.Println()

	// Example 1: Fast request (should succeed)
	fmt.Println("1️⃣  Fetching httpbin.org with 5s timeout...")
	body, err := fetchWithTimeout("https://httpbin.org/get", 5*time.Second)
	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		// Print just the first 200 chars to keep output clean
		if len(body) > 200 {
			body = body[:200] + "..."
		}
		fmt.Printf("   ✅ Response (%d bytes): %s\n", len(body), body)
	}

	fmt.Println()

	// Example 2: Deliberately short timeout (should fail)
	fmt.Println("2️⃣  Fetching with impossibly short timeout (1ms)...")
	_, err = fetchWithTimeout("https://httpbin.org/delay/3", 1*time.Millisecond)
	if err != nil {
		fmt.Printf("   ❌ Expected timeout: %v\n", err)
	} else {
		fmt.Println("   ✅ Unexpected success!")
	}

	fmt.Println()
	fmt.Println("KEY TAKEAWAY:")
	fmt.Println("  - ALWAYS use context.WithTimeout for HTTP requests")
	fmt.Println("  - Use http.NewRequestWithContext to attach the context")
	fmt.Println("  - Check ctx.Err() == context.DeadlineExceeded for timeout detection")
	fmt.Println("  - Production default: 5-30 seconds for API calls, 1-5 seconds for health checks")
}
