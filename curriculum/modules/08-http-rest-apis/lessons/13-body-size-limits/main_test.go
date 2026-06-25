package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCommentHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantPrefix string
	}{
		{
			name:       "valid small body",
			body:       `{"text":"hello"}`,
			wantStatus: http.StatusOK,
			wantPrefix: "comment received: hello",
		},
		{
			name:       "body over limit",
			body:       `{"text":"` + strings.Repeat("a", 200) + `"}`,
			wantStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:       "malformed json",
			body:       `not json`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/comment", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			commentHandler(rec, req)

			resp := rec.Result()
			resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
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

func TestMaxBytesMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid within limit",
			body:       `{"text":"ok"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "exceeds limit",
			body:       `{"text":"` + strings.Repeat("x", 200) + `"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	handler := maxBytesMiddleware(50)(http.HandlerFunc(secureCommentHandler))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/secure-comment", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			resp := rec.Result()
			resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestCommentHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/comment", nil)
	rec := httptest.NewRecorder()

	commentHandler(rec, req)

	resp := rec.Result()
	resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestLessonCompiles(t *testing.T) {
}
