package main

import (
	"strings"
	"testing"
)

func TestBuildHTTPRequest_GET(t *testing.T) {
	headers := map[string]string{"Host": "example.com"}
	raw := buildHTTPRequest("GET", "/", headers, "")
	if !strings.HasPrefix(raw, "GET / HTTP/1.1\r\n") {
		t.Errorf("expected GET request line, got %q", raw)
	}
	if !strings.Contains(raw, "Host: example.com") {
		t.Errorf("expected Host header, got %q", raw)
	}
	if strings.Contains(raw, "Content-Length") {
		t.Errorf("GET with no body should not have Content-Length")
	}
	if !strings.Contains(raw, "\r\n\r\n") {
		t.Errorf("expected empty line separating headers from body")
	}
}

func TestBuildHTTPRequest_POST(t *testing.T) {
	headers := map[string]string{"Content-Type": "application/json"}
	body := `{"name":"Alice"}`
	raw := buildHTTPRequest("POST", "/api/users", headers, body)
	if !strings.HasPrefix(raw, "POST /api/users HTTP/1.1\r\n") {
		t.Errorf("expected POST request line, got %q", raw)
	}
	if !strings.Contains(raw, "Content-Length: 16") {
		t.Errorf("expected Content-Length 16, got %q", raw)
	}
	if !strings.HasSuffix(raw, body) {
		t.Errorf("expected body suffix %q", body)
	}
}

func TestClassifyStatus(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{100, "Informational"},
		{200, "Success"},
		{201, "Success"},
		{301, "Redirection"},
		{400, "Client Error"},
		{404, "Client Error"},
		{500, "Server Error"},
		{503, "Server Error"},
		{99, "Unknown"},
		{600, "Unknown"},
	}
	for _, tc := range tests {
		got := classifyStatus(tc.code)
		if got != tc.want {
			t.Errorf("classifyStatus(%d) = %q; want %q", tc.code, got, tc.want)
		}
	}
}
