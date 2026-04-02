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
// Section 16: HTTP Clients & Mocking — Table-Driven Mocks
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Combining Table-Driven tests with Function Mocks
//   - Testing success, network errors, and JSON errors in one loop
//
// ENGINEERING DEPTH:
//   This represents the holy grail of standard Library Go testing.
//   We define a table of test cases. Each case provides its own `mockResponse`
//   and `expectedError`.
//   This pattern yields 100% test coverage with minimal duplicate code.
//
// RUN: go test -v ./16-http-clients-and-mocking/5-manual-mock-table-driven
// ============================================================================

func main() {
	fmt.Println("=== Table-Driven Mocking ===")

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

	fmt.Println("\n(See client_test.go to see the ultimate Table-Driven mock pattern!)")
}
