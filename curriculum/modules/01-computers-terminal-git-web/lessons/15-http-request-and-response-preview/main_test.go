package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestHTTPServerEndpoints(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Welcome to The Go Engineer HTTP Server!")
	})
	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "Hello from Go!", "status": "ok"}`)
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to bind listener: %v", err)
	}
	defer listener.Close()

	server := &http.Server{Handler: mux}
	go func() {
		_ = server.Serve(listener)
	}()
	defer server.Shutdown(context.Background())

	baseURL := fmt.Sprintf("http://%s", listener.Addr().String())
	client := &http.Client{Timeout: 3 * time.Second}

	// 1. Test Root Endpoint (text/plain)
	respRoot, err := client.Get(baseURL + "/")
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	defer respRoot.Body.Close()

	if respRoot.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", respRoot.StatusCode)
	}
	if !strings.Contains(respRoot.Header.Get("Content-Type"), "text/plain") {
		t.Errorf("expected text/plain Content-Type, got %s", respRoot.Header.Get("Content-Type"))
	}
	bodyRoot, _ := io.ReadAll(respRoot.Body)
	if !strings.Contains(string(bodyRoot), "Welcome to The Go Engineer") {
		t.Errorf("unexpected body content: %s", string(bodyRoot))
	}

	// 2. Test API Endpoint (application/json)
	respAPI, err := client.Get(baseURL + "/api/hello")
	if err != nil {
		t.Fatalf("GET /api/hello failed: %v", err)
	}
	defer respAPI.Body.Close()

	if respAPI.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", respAPI.StatusCode)
	}
	if respAPI.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json Content-Type, got %s", respAPI.Header.Get("Content-Type"))
	}

	var payload struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	if err := json.NewDecoder(respAPI.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}
	if payload.Status != "ok" || payload.Message != "Hello from Go!" {
		t.Errorf("unexpected JSON payload: %+v", payload)
	}
}

func TestMainExecution(t *testing.T) {
	main()
}

