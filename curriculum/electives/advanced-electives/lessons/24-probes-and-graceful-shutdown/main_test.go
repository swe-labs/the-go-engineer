package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHealthzAlwaysOK(t *testing.T) {
	s := NewServer(nil)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	s.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("healthz status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestReadyzReflectsReadiness(t *testing.T) {
	s := NewServer(nil)
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)

	rec := httptest.NewRecorder()
	s.routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("ready status = %d, want %d", rec.Code, http.StatusOK)
	}

	s.NotReady()
	rec = httptest.NewRecorder()
	s.routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("draining status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestShutdownMarksNotReady(t *testing.T) {
	s := NewServer(nil)
	if !s.Ready() {
		t.Fatal("new server is not ready")
	}

	if err := s.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	if s.Ready() {
		t.Fatal("server is still ready after shutdown")
	}
}

func TestShutdownWaitsForInflightRequest(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once

	work := func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() { close(started) })
		<-release
		w.WriteHeader(http.StatusOK)
	}

	s := NewServer(work)
	if err := s.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("start: %v", err)
	}

	requestDone := make(chan error, 1)
	go func() {
		resp, err := http.Get("http://" + s.Addr() + "/work")
		if err == nil {
			resp.Body.Close()
		}
		requestDone <- err
	}()

	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	shutdownDone := make(chan error, 1)
	go func() { shutdownDone <- s.Shutdown(ctx) }()

	select {
	case <-shutdownDone:
		t.Fatal("shutdown returned before the in-flight request completed")
	case <-time.After(100 * time.Millisecond):
	}

	if s.Ready() {
		t.Fatal("server is still ready during shutdown")
	}

	close(release)

	if err := <-requestDone; err != nil {
		t.Fatalf("request: %v", err)
	}
	if err := <-shutdownDone; err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}
