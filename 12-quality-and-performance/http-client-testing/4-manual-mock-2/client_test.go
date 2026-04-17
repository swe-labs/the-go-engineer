// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

// 1. Dynamic Function Mocking
// By storing a function signature `func(...)` inside a struct field,
// we can define Custom logic per-test without creating a dozen different Mock structs.
type MockHTTPClient struct {
	GetFunc func(url string) (resp *http.Response, err error)

	// We can also track state: how many times was this mock called?
	GetCalls []string
}

// Get executes whatever anonymous function is stored in `GetFunc`.
func (c *MockHTTPClient) Get(url string) (*http.Response, error) {
	c.GetCalls = append(c.GetCalls, url) // Record the spy call
	if c.GetFunc != nil {
		return c.GetFunc(url) // Execute the dynamically injected test logic
	}

	return nil, errors.New("GetFunc not implemented yet")
}

func mockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}

func TestPostsClient_FetchPosts_Success(t *testing.T) {

	mockClient := &MockHTTPClient{
		GetFunc: func(url string) (resp *http.Response, err error) {
			return mockResponse(http.StatusOK, dummyPosts), nil
		},
	}

	client := NewPostsClient(mockClient, "http://example.com")
	posts, err := client.FetchPosts(2)

	if err != nil {
		t.Errorf("Error while fetching posts: %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}

	if posts[0].Title != "His mother had always taught him" {
		t.Errorf("post at index 0 should have this title: His mother had always taught him")
	}

	expectedURL := "http://example.com/posts?limit=2"
	if len(mockClient.GetCalls) != 1 {
		t.Errorf("Get called %v times, want %v", len(mockClient.GetCalls), 1)
	}

	if mockClient.GetCalls[0] != expectedURL {
		t.Errorf("Get called %v times, want %v", len(mockClient.GetCalls), 1)
	}

}

func TestPostsClient_FetchPosts_NetworkError(t *testing.T) {

	mockClient := &MockHTTPClient{
		GetFunc: func(url string) (resp *http.Response, err error) {
			return nil, errors.New("network error")
		},
	}

	client := NewPostsClient(mockClient, "http://example.com")
	posts, err := client.FetchPosts(2)

	if err == nil {
		t.Errorf("expected error not to be nil")
	}

	if posts != nil {
		t.Errorf("posts should be nil")
	}

}

func TestPostsClient_FetchPosts_NotOK(t *testing.T) {

	mockClient := &MockHTTPClient{
		GetFunc: func(url string) (resp *http.Response, err error) {
			return mockResponse(http.StatusInternalServerError, ""), nil
		},
	}

	client := NewPostsClient(mockClient, "http://example.com")
	posts, err := client.FetchPosts(2)

	if err == nil {
		t.Fatalf("expected error not to be nil")
	}

	if err.Error() != "unexpected status code: 500" {
		t.Errorf("got error %v, want %v", err.Error(), "unexpected status code: 200")
	}

	if posts != nil {
		t.Errorf("posts should be nil")
	}

}

func TestPostsClient_FetchPosts_InvalidJSON(t *testing.T) {

	mockClient := &MockHTTPClient{
		GetFunc: func(url string) (resp *http.Response, err error) {
			return mockResponse(http.StatusOK, "{"), nil
		},
	}

	client := NewPostsClient(mockClient, "http://example.com")
	posts, err := client.FetchPosts(2)

	if err == nil {
		t.Fatalf("expected error not to be nil")
	}

	if err.Error() != "failed to decode response: unexpected EOF" {
		t.Errorf("got error %v, want %v", err.Error(), "failed to decode response: unexpected EOF")
	}

	if posts != nil {
		t.Errorf("posts should be nil")
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
