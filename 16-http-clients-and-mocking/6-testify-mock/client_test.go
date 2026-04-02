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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 1. Testify Mock Embedding
// By anonymously embedding `mock.Mock`, MockHTTPClient inherits all of Testify's
// argument-capturing and reflection assertion methods.
type MockHTTPClient struct {
	mock.Mock
}

func (c *MockHTTPClient) Get(url string) (*http.Response, error) {

	// 2. Argument Recording
	// `.Called()` records that this method was invoked, tracking the injected `url` string.
	args := c.Called(url)

	// Testify returns an `Arguments` array. We must type-assert index 0 to *http.Response.
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*http.Response), args.Error(1)
}

func mockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}

func TestPostsClient_FetchPosts_Success(t *testing.T) {

	mockClient := new(MockHTTPClient)

	mockRes := mockResponse(http.StatusOK, dummyPosts)

	// 3. Testify Assertion Chain
	// `.On()` tells the mock: "If 'Get' is called exactly with this URL string..."
	// `.Return()` tells the mock: "...then immediately return this response and nil error."
	mockClient.On("Get", "http://example.com/posts?limit=1").Return(mockRes, nil)

	client := NewPostsClient(mockClient, "http://example.com")
	posts, err := client.FetchPosts(1)

	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.Len(t, posts, 1)
	assert.Equal(t, "His mother had always taught him", posts[0].Title)

}

func TestPostsClient_FetchPosts_NetworkError(t *testing.T) {

	mockClient := new(MockHTTPClient)

	mockClient.On("Get", "http://example.com/posts?limit=1").Return(nil, errors.New("network error"))

	client := NewPostsClient(mockClient, "http://example.com")
	posts, err := client.FetchPosts(1)

	assert.Error(t, err)
	assert.Nil(t, posts)

}

func TestPostsClient_FetchPosts_NotOK(t *testing.T) {

	mockClient := new(MockHTTPClient)

	mockRes := mockResponse(http.StatusInternalServerError, "")

	mockClient.On("Get", "http://example.com/posts?limit=1").Return(mockRes, nil)

	client := NewPostsClient(mockClient, "http://example.com")
	posts, err := client.FetchPosts(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
	assert.Nil(t, posts)

}

func TestPostsClient_FetchPosts_InvalidJSON(t *testing.T) {

	mockClient := new(MockHTTPClient)

	mockRes := mockResponse(http.StatusOK, "{")

	mockClient.On("Get", "http://example.com/posts?limit=1").Return(mockRes, nil)

	client := NewPostsClient(mockClient, "http://example.com")
	posts, err := client.FetchPosts(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode response: unexpected EOF")
	assert.Nil(t, posts)

}

const dummyPosts = `{
"posts": [
{"id": 1, "title": "His mother had always taught him", "body": "Content 1", "userId": 1}
],
"total": 1
}`
