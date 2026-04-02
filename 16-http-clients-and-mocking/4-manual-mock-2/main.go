// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"fmt"
	"net/http"
	"os"
)

// ============================================================================
// Section 16: HTTP Clients & Mocking — Function Injection Mocking
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Improving manual mocks using first-class functions
//   - Tracking function invocations (`GetCalls`)
//
// ENGINEERING DEPTH:
//   A massive step up from basic structs is storing a FUNCTION on your mock struct.
//   `GetFunc func(url string) (*http.Response, error)`
//
//   Now, inside the test, you define exactly what `GetFunc` does! This allows
//   you to write custom logic per-test (e.g., returning errors based on URL).
//
// RUN: go test -v ./16-http-clients-and-mocking/4-manual-mock-2
// ============================================================================

func main() {
	fmt.Println("=== Function-Injection Mocking ===")

	baseURL := "https://dummyjson.com"
	postsClient := NewPostsClient(http.DefaultClient, baseURL)

	posts, err := postsClient.FetchPosts(3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for _, post := range posts {
		fmt.Printf("[%d] %s\n", post.ID, post.Title)
	}

	fmt.Println("\n(See client_test.go to see the GetFunc pattern in action!)")
}
