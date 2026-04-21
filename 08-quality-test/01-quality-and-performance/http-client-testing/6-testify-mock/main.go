// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

// RUN: go run ./08-quality-test/01-quality-and-performance/http-client-testing/6-testify-mock
package main

import (
	"fmt"
	"net/http"
	"os"
)

// ============================================================================
// Section 16: HTTP Clients & Mocking — Testify Mock Module
// Level: Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Using 3rd party mocking libraries (`github.com/stretchr/testify/mock`)
//   - `mock.On("Get", url).Return(resp, err)` syntax
//
// ENGINEERING DEPTH:
//   While Go purists prefer manual mocking (Section 16-5), large enterprise
//   projects with huge interfaces often use code generators (like `mockery`)
//   combined with `testify/mock`.
//   It allows extreme expressive assertions: `mock.AssertNumberOfCalls(t, "Get", 1)`
//
// RUN: go test -v ./08-quality-test/01-quality-and-performance/http-client-testing/6-testify-mock
// ============================================================================

func main() {
	fmt.Println("=== Testify Mock Engine ===")

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

	fmt.Println("\n(See client_test.go to see testify/mock in action!)")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: PR.1 CPU profile")
	fmt.Println("   Current: HM.4 (testify/mock)")
	fmt.Println("---------------------------------------------------")
}
