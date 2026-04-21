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
// Section 16: HTTP Clients & Mocking — Basic Manual Mock
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Implementing a fake struct that satisfies an interface
//   - Hardcoding mock responses for tests
//
// ENGINEERING DEPTH:
//   This is the simplest form of mocking. We create a `MockHTTPClient` struct
//   that implements `Get(url)`. In our tests, we set `expectedBody` and
//   `expectedCode` on the struct, and pass it to our API client.
//
//   DOWNSIDE: If you need to test 10 different URLs with 10 different responses,
//   this basic struct becomes very unwieldy.
//
// RUN: go run ./08-quality-test/01-quality-and-performance/http-client-testing/3-manual-mock
// ============================================================================

func main() {
	fmt.Println("=== Basic Manual Mocking ===")

	baseURL := "https://dummyjson.com"
	postsClient := NewPostsClient(http.DefaultClient, baseURL)

	posts, err := postsClient.FetchPosts(2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for _, post := range posts {
		fmt.Printf("[%d] %s\n", post.ID, post.Title)
	}

	fmt.Println("\n(See client_test.go to understand how this is mocked!)")
	fmt.Println("\n---------------------------------------------------")
	fmt.Println("🚀 NEXT UP: HM.2 function-injection mock")
	fmt.Println("   Current: HM.1 (manual mock)")
	fmt.Println("---------------------------------------------------")
}
