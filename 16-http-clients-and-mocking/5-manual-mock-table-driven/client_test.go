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
)

type MockHTTPClient struct {
	GetFunc func(url string) (resp *http.Response, err error)

	GetCalls []string
}

func (c *MockHTTPClient) Get(url string) (*http.Response, error) {
	c.GetCalls = append(c.GetCalls, url)
	if c.GetFunc != nil {
		return c.GetFunc(url)
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

func TestPostsClient_FetchPosts(t *testing.T) {

	// 1. Table-Driven Tests
	// By defining an anonymous `struct` and immediately instantiating a slice `[]struct{}`
	// we avoid writing 15 different nearly-identical test functions.
	tests := []struct {
		name            string                                   // Used for naming the t.Run subtest
		getFunc         func(url string) (*http.Response, error) // The dynamic mock payload
		limit           int
		wantErr         bool
		errContains     string
		wantPostsNil    bool
		wantPostsLen    int
		expectedURL     string
		expectedGetCall int
	}{
		{
			name: "success",
			getFunc: func(url string) (*http.Response, error) {
				return mockResponse(http.StatusOK, dummyPosts), nil
			},
			limit:           2,
			wantErr:         false,
			wantPostsNil:    false,
			wantPostsLen:    2,
			expectedURL:     "http://example.com/posts?limit=2",
			expectedGetCall: 1,
		},
		{
			name: "network error",
			getFunc: func(url string) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			limit:           2,
			wantErr:         true,
			errContains:     "network error",
			wantPostsNil:    true,
			expectedURL:     "http://example.com/posts?limit=2",
			expectedGetCall: 1,
		},

		{
			name: "non-200 status",
			getFunc: func(url string) (*http.Response, error) {
				return mockResponse(http.StatusInternalServerError, ""), nil
			},
			limit:           2,
			wantErr:         true,
			errContains:     "unexpected status code: 500",
			wantPostsNil:    true,
			expectedURL:     "http://example.com/posts?limit=2",
			expectedGetCall: 1,
		},

		{
			name: "invalid JSON",
			getFunc: func(url string) (*http.Response, error) {
				return mockResponse(http.StatusOK, "{"), nil
			},
			limit:           2,
			wantErr:         true,
			errContains:     "failed to decode response: unexpected EOF",
			wantPostsNil:    true,
			expectedURL:     "http://example.com/posts?limit=2",
			expectedGetCall: 1,
		},
	}

	for _, tt := range tests {
		// 2. Sub-testing
		// `t.Run()` executes each iteration as a fully isolated subtest.
		// If one case panics, the other cases in the table still execute!
		t.Run(tt.name, func(t *testing.T) { // sub-test
			mockClient := &MockHTTPClient{
				GetFunc: tt.getFunc,
			}
			client := NewPostsClient(mockClient, "http://example.com")

			posts, err := client.FetchPosts(tt.limit)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			if tt.wantPostsNil {
				assert.Nil(t, posts)
			} else {
				if assert.NotNil(t, posts) {
					assert.Len(t, posts, tt.wantPostsLen)
				}
			}

			if tt.expectedGetCall > 0 {
				assert.Len(t, mockClient.GetCalls, tt.expectedGetCall)
				if tt.expectedURL != "" && len(mockClient.GetCalls) > 0 {
					assert.Equal(t, tt.expectedURL, mockClient.GetCalls[0])
				}
			}

		})
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
