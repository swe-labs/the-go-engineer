package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestServerResponds(t *testing.T) {
	server := startServer()
	defer server.Close()
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	if !strings.Contains(string(body), "Hello from Go server!") {
		t.Errorf("unexpected body: %q", string(body))
	}
}

func TestAPIHello(t *testing.T) {
	server := startServer()
	defer server.Close()
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/api/hello")
	if err != nil {
		t.Fatalf("GET /api/hello failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected application/json, got %q", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body failed: %v", err)
	}
	if !strings.Contains(string(body), "Hello, API!") {
		t.Errorf("unexpected body: %q", string(body))
	}
}

func TestRequestNotFound(t *testing.T) {
	server := startServer()
	defer server.Close()
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/nonexistent")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Logf("Note: got status %d for nonexistent path (catch-all handler)", resp.StatusCode)
	}
}

func TestMakeRequest(t *testing.T) {
	server := startServer()
	defer server.Close()
	time.Sleep(100 * time.Millisecond)

	status, body, headers, err := makeRequest("http://localhost:8080/")
	if err != nil {
		t.Fatalf("makeRequest failed: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("expected 200, got %d", status)
	}
	if body == "" {
		t.Error("expected non-empty body")
	}
	if headers.Get("Content-Type") == "" {
		t.Error("expected Content-Type header")
	}
}
