// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 1. The Mocking Interface
// Instead of directly using `http.Client` (which fires real network connections),
// we define an interface that matches its method signature.
// At runtime, we can pass either a real `http.Client` OR a fake mock struct!
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// 2. Dependency Injection (DI) Component
// PostsClient depends on the `HTTPClient` interface, not a concrete struct.
// This is the core principle of SOLID architecture: "Depend on abstractions".
type PostsClient struct {
	httpClient HTTPClient
	baseURL    string
}

// NewPostsClient is our explicit DI constructor.
func NewPostsClient(httpClient HTTPClient, baseURL string) *PostsClient {
	return &PostsClient{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

type Post struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

type PostsResponse struct {
	Posts []Post `json:"posts"`
	Total int    `json:"total"`
}

func (c *PostsClient) FetchPosts(limit int) ([]Post, error) {

	url := fmt.Sprintf("%s/posts?limit=%d", c.baseURL, limit)

	// 3. Polymorphic Call
	// If running in production, this calls `http.Client.Get()`.
	// If running in a unit test, this calls `MockHTTPClient.Get()`.
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var postsResp PostsResponse
	if err := json.NewDecoder(resp.Body).Decode(&postsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return postsResp.Posts, nil
}
