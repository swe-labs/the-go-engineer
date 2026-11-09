package main

import (
	"net/http"
	"testing"
	"time"
)

func TestStartServer(t *testing.T) {
	srv, addr, err := startServer("127.0.0.1:0")
	if err != nil {
		t.Fatalf("startServer failed: %v", err)
	}
	defer srv.Close()

	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://" + addr.String() + "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestStartServerInvalidAddr(t *testing.T) {
	_, _, err := startServer("999.999.999.999:0")
	if err == nil {
		t.Error("expected error for invalid address")
	}
}
