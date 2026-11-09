package main

import (
	"context"
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
	if got := strings.TrimSpace(rec.Body.String()); got != "fast response" {
		t.Errorf("body = %q, want %q", got, "fast response")
	}
}

func TestSlowHandlerContextCancelled(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	ctx, cancel := context.WithCancel(req.Context())
	req = req.WithContext(ctx)
	cancel()

	rec := httptest.NewRecorder()
	slowHandler(rec, req)

	// handler returns without writing when context is cancelled
}

func TestEchoHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantPrefix string
	}{
		{"echo text", "hello", http.StatusOK, "echo: hello"},
		{"echo json", `{"key":"val"}`, http.StatusOK, "echo: {\"key\":\"val\"}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			echoHandler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantPrefix != "" {
				got := strings.TrimSpace(rec.Body.String())
				if !strings.HasPrefix(got, tt.wantPrefix) {
					t.Errorf("body = %q, want prefix %q", got, tt.wantPrefix)
				}
			}
		})
	}
}

func TestNewServerConfig(t *testing.T) {
	srv := newServer(":0", 5*time.Second, 10*time.Second, 60*time.Second, 2*time.Second)

	if srv.ReadTimeout != 5*time.Second {
		t.Errorf("ReadTimeout = %v, want %v", srv.ReadTimeout, 5*time.Second)
	}
	if srv.WriteTimeout != 10*time.Second {
		t.Errorf("WriteTimeout = %v, want %v", srv.WriteTimeout, 10*time.Second)
	}
	if srv.IdleTimeout != 60*time.Second {
		t.Errorf("IdleTimeout = %v, want %v", srv.IdleTimeout, 60*time.Second)
	}
	if srv.ReadHeaderTimeout != 2*time.Second {
		t.Errorf("ReadHeaderTimeout = %v, want %v", srv.ReadHeaderTimeout, 2*time.Second)
	}
}

func TestFetchWithTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	_, err := fetchWithTimeout(server.URL, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFetchWithTimeoutExceeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	_, err := fetchWithTimeout(server.URL, 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestLessonCompiles(t *testing.T) {
}
