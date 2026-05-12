// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHelloWorldHandler (Test): demonstrates httptest.NewRecorder for capturing handler output.
func TestHelloWorldHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()

	HelloWorldHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	body, _ := io.ReadAll(rec.Body)
	assert.Equal(t, "Hello World!", string(body))
}

// TestEchoHandler (Test): demonstrates testing a handler that reads the request body.
func TestEchoHandler(t *testing.T) {
	payload := `{"id": 123, "action": "echo"}`
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()

	EchoHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	body, _ := io.ReadAll(rec.Body)
	assert.Equal(t, payload, string(body))
}

// TestEchoHandler_WrongMethod (Test): demonstrates testing HTTP method rejection.
func TestEchoHandler_WrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	rec := httptest.NewRecorder()

	EchoHandler(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

	body, _ := io.ReadAll(rec.Body)
	assert.Equal(t, "Method not allowed\n", string(body))
}
