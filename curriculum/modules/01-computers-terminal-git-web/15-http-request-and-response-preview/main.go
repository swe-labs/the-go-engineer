package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello from Go server!\nMethod: %s\nPath: %s\n", r.Method, r.URL.Path)
}

func startServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message": "Hello, API!", "method": "%s"}`, r.Method)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	return server
}

func makeRequest(url string) (int, string, http.Header, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, "", nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", nil, fmt.Errorf("read failed: %w", err)
	}

	return resp.StatusCode, string(body), resp.Header, nil
}

func main() {
	server := startServer()
	defer server.Close()
	time.Sleep(100 * time.Millisecond)

	fmt.Println("=== HTTP Request/Response Preview ===")
	fmt.Println()
	fmt.Println("Server started on http://localhost:8080")
	fmt.Println()

	urls := []string{
		"http://localhost:8080/",
		"http://localhost:8080/api/hello",
	}

	for _, url := range urls {
		fmt.Printf("--- GET %s ---\n", url)
		status, body, headers, err := makeRequest(url)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("Status: %d %s\n", status, http.StatusText(status))
		fmt.Printf("Content-Type: %s\n", headers.Get("Content-Type"))
		fmt.Printf("Body: %s\n", body)
		fmt.Println()
	}
}
