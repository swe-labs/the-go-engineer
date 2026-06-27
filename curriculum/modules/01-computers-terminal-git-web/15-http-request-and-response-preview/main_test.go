package main

import (
	"strings"
	"testing"
)

func TestHTTPResponseStatus(t *testing.T) {
	resp := httpResponse(200, "OK")
	if !strings.Contains(resp, "200 OK") {
		t.Errorf("expected 200 OK in response, got %q", resp)
	}
}

func TestHTTPRequestRoot(t *testing.T) {
	resp := httpRequest("GET", "/", "")
	if !strings.Contains(resp, "Hello from Go server!") {
		t.Errorf("expected greeting in response, got %q", resp)
	}
}

func TestHTTPRequestAPI(t *testing.T) {
	resp := httpRequest("GET", "/api/hello", "")
	if !strings.Contains(resp, "Hello, API!") {
		t.Errorf("expected API greeting in response, got %q", resp)
	}
}

func TestHTTPRequestNotFound(t *testing.T) {
	resp := httpRequest("GET", "/nonexistent", "")
	if !strings.Contains(resp, "404") {
		t.Errorf("expected 404 in response, got %q", resp)
	}
}

func TestHTTPResponseFormat(t *testing.T) {
	resp := httpResponse(200, "test body")
	if !strings.Contains(resp, "HTTP/1.1") {
		t.Errorf("expected HTTP/1.1 in response, got %q", resp)
	}
	if !strings.Contains(resp, "Content-Type") {
		t.Errorf("expected Content-Type header, got %q", resp)
	}
}
