// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0
// Commercial use is prohibited without permission.

package user

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Stage 08: Quality and Performance - HTTP Handlers (`httptest`)
// Level: Intermediate
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Testing HTTP handlers without starting a real web server
//   - Using httptest.NewRecorder() to capture responses
//   - Using httptest.NewRequest() to mock incoming requests
//
// ENGINEERING DEPTH:
//   Starting a real web server (e.g. `http.ListenAndServe`) inside a unit
//   test is slow, brittle, and can cause port conflicts.
//   Instead, the `httptest` package provides a fake ResponseWriter
//   (`ResponseRecorder`) that records exactly what the handler writes to it.
// ============================================================================

// HelloWorldHandler simply returns a static greeting.
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!"))
}

// EchoHandler reads the request POST body and sends it right back.
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func TestHelloWorldHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	recorder := httptest.NewRecorder()

	HelloWorldHandler(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	bodyBytes, err := io.ReadAll(recorder.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", string(bodyBytes))
}

func TestEchoHandler(t *testing.T) {
	payload := `{"id": 123, "action": "echo"}`
	bodyReader := bytes.NewBufferString(payload)

	req := httptest.NewRequest(http.MethodPost, "/echo", bodyReader)
	recorder := httptest.NewRecorder()

	EchoHandler(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, payload, string(responseBody))
}

func TestEchoHandler_WrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	recorder := httptest.NewRecorder()

	EchoHandler(recorder, req)

	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)

	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, "Method not allowed\n", string(responseBody))
}
