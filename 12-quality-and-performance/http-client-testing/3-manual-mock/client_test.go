// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

// 1. The Fake Struct
// This struct manually implements the `HTTPClient` interface.
type MockHTTPClient struct {
	expectedBody string
	expectedCode int
}

// 2. Spoofing the Response
// By attaching the `Get()` method, MockHTTPClient satisfies `HTTPClient`.
func (c *MockHTTPClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		// 3. Spoofing the Network Stream
		// The `http.Response.Body` requires an `io.ReadCloser` interface.
		// `bytes.NewBufferString` creates a readable memory stream from a raw string.
		// `io.NopCloser` wraps it, adding a fake `Close()` method that does nothing
		// so it perfectly matches the `io.ReadCloser` signature.
		Body:       io.NopCloser(bytes.NewBufferString(c.expectedBody)),
		StatusCode: c.expectedCode,
		Header:     make(http.Header),
	}, nil
}

func TestPostsClient_FetchPosts_Success(t *testing.T) {

	// 4. Injecting the Mock
	// We create an instance of our fake, supplying a static JSON string.
	mockClient := &MockHTTPClient{
		expectedBody: dummyPosts,
		expectedCode: http.StatusOK,
	}
	// We pass the fake struct directly into the constructor!
	client := NewPostsClient(mockClient, "http://example.com")
	posts, err := client.FetchPosts(2)
	if err != nil {
		t.Fatalf("failed to fetch posts: %v", err)
	}

	if posts == nil {
		t.Fatalf("posts should not be nil")
	}

	if len(posts) != 2 {
		t.Fatalf("posts should have 2 posts, but has %d", len(posts))
	}

}

func TestPostsClient_FetchPosts_NotOK(t *testing.T) {

	mockClient := &MockHTTPClient{
		expectedBody: dummyPosts,
		expectedCode: http.StatusBadRequest,
	}
	client := NewPostsClient(mockClient, "http://example.com")
	posts, err := client.FetchPosts(2)
	if err == nil {
		t.Fatalf("should have errored")
	}

	if posts != nil {
		t.Fatalf("posts should be nil")
	}

}

func TestPostsClient_FetchPosts_JSON_Error(t *testing.T) {

	mockClient := &MockHTTPClient{
		expectedBody: `{`,
		expectedCode: http.StatusOK,
	}
	client := NewPostsClient(mockClient, "http://example.com")
	posts, err := client.FetchPosts(2)
	if err == nil {
		t.Fatalf("should have errored")
	}

	if posts != nil {
		t.Fatalf("posts should be nil")
	}

}

const dummyPosts = `{
  "posts": [
    {
      "id": 1,
      "title": "His mother had always taught him",
      "body": "His mother had always taught him not to ever think of himself as better than others. He'd tried to live by this motto. He never looked down on those who were less fortunate or who had less money than him. But the stupidity of the group of people he was talking to made him change his mind.",
      "tags": [
        "history",
        "american",
        "crime"
      ],
      "reactions": {
        "likes": 192,
        "dislikes": 25
      },
      "views": 305,
      "userId": 121
    },
    {
      "id": 2,
      "title": "He was an expert but not in a discipline",
      "body": "He was an expert but not in a discipline that anyone could fully appreciate. He knew how to hold the cone just right so that the soft server ice-cream fell into it at the precise angle to form a perfect cone each and every time. It had taken years to perfect and he could now do it without even putting any thought behind it.",
      "tags": [
        "french",
        "fiction",
        "english"
      ],
      "reactions": {
        "likes": 859,
        "dislikes": 32
      },
      "views": 4884,
      "userId": 91
    }
  ],
  "total": 251,
  "skip": 0,
  "limit": 2
}`
