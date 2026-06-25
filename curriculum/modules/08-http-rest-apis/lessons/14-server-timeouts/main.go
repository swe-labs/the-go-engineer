package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func slowHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(5 * time.Second):
		fmt.Fprintln(w, "done")
	case <-r.Context().Done():
		return
	}
}

func fastHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "fast response")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read error", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "echo: %s", body)
}

func newServer(addr string, readTimeout, writeTimeout, idleTimeout, readHeaderTimeout time.Duration) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/fast", fastHandler)
	mux.HandleFunc("/slow", slowHandler)
	mux.HandleFunc("/echo", echoHandler)

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}
}

func main() {
	srv := newServer(":8080", 5*time.Second, 10*time.Second, 60*time.Second, 2*time.Second)

	log.Println("Server starting on :8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func fetchWithTimeout(url string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
