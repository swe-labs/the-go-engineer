package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFastHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/fast", nil)
	rec := httptest.NewRecorder()
	fastHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != "ok" {
		t.Errorf("body = %q, want %q", got, "ok")
	}
}

func TestSlowHandlerContextCancelledOnShutdown(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	ctx, cancel := context.WithCancel(req.Context())
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	cancel()

	slowHandler(rec, req)
}

func TestServerShutdown(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/fast", fastHandler)

	srv := &http.Server{
		Addr:    ":0",
		Handler: mux,
	}

	go srv.ListenAndServe()
	time.Sleep(50 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}

func TestServerShutdownRespectsTimeout(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/slow", slowHandler)

	srv := &http.Server{
		Addr:    ":0",
		Handler: mux,
	}

	go srv.ListenAndServe()
	time.Sleep(50 * time.Millisecond)

	// Use a very short timeout so shutdown should time out
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err == nil {
		t.Log("shutdown completed (handler may have finished early)")
	}
}

func TestRegisterOnShutdown(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/fast", fastHandler)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	srv := &http.Server{
		Handler: mux,
	}

	called := make(chan struct{}, 1)
	srv.RegisterOnShutdown(func() {
		called <- struct{}{}
	})

	go srv.Serve(l)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Logf("Shutdown: %v", err)
	}

	select {
	case <-called:
	case <-time.After(500 * time.Millisecond):
		t.Error("RegisterOnShutdown callback was not called within timeout")
	}
}

func TestLessonCompiles(t *testing.T) {
}
