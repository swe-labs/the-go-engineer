package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func main() {
	fmt.Println("=== HTTP Request / Response Lifecycle Preview ===")

	// 1. Setup local HTTP router with two endpoints
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

	// 2. Bind local listener on ephemeral port (port 0)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Printf("Failed to bind listener: %v\n", err)
		return
	}
	defer listener.Close()

	server := &http.Server{Handler: mux}
	go func() {
		_ = server.Serve(listener)
	}()
	defer server.Shutdown(context.Background())

	serverURL := fmt.Sprintf("http://%s/api/hello", listener.Addr().String())
	fmt.Printf("Server listening actively on: %s\n\n", serverURL)

	// 3. Issue client HTTP GET request
	client := &http.Client{Timeout: 3 * time.Second}
	fmt.Printf("[Client] Sending GET request to: %s\n", serverURL)

	resp, err := client.Get(serverURL)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 4. Deconstruct HTTP Response
	fmt.Println("\n=== Deconstructing HTTP Response Wire Data ===")
	fmt.Printf("HTTP Protocol Version:  %s\n", resp.Proto)
	fmt.Printf("Status Line:            %s (Code: %d)\n", resp.Status, resp.StatusCode)
	fmt.Println("Headers:")
	for key, values := range resp.Header {
		fmt.Printf("  %s: %s\n", key, values[0])
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	fmt.Printf("\nResponse Body Payload:\n%s\n", string(bodyBytes))
	fmt.Println("HTTP Lifecycle complete: Client Request -> Server Route -> Status Code + Headers + Body.")
}
