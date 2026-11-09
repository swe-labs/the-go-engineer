package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSAllowedOrigin(t *testing.T) {
	handler := corsMiddleware([]string{"https://example.com"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name       string
		origin     string
		wantStatus int
		wantACAO   string
	}{
		{
			name:       "allowed origin returns ACAO header",
			origin:     "https://example.com",
			wantStatus: http.StatusOK,
			wantACAO:   "https://example.com",
		},
		{
			name:       "disallowed origin returns no ACAO",
			origin:     "https://evil.com",
			wantStatus: http.StatusOK,
			wantACAO:   "",
		},
		{
			name:       "no origin header returns no ACAO",
			origin:     "",
			wantStatus: http.StatusOK,
			wantACAO:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tc.wantStatus)
			}
			gotACAO := w.Header().Get("Access-Control-Allow-Origin")
			if gotACAO != tc.wantACAO {
				t.Errorf("ACAO = %q, want %q", gotACAO, tc.wantACAO)
			}
		})
	}
}

func TestCORSPreflight(t *testing.T) {
	handler := corsMiddleware([]string{"https://example.com"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("preflight status = %d, want %d", w.Code, http.StatusNoContent)
	}
	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("expected Access-Control-Allow-Methods header on preflight")
	}
	if w.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Error("expected Access-Control-Allow-Headers header on preflight")
	}
}

func TestCORSWithCredentials(t *testing.T) {
	handler := corsMiddlewareWithCredentials([]string{"https://example.com"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name     string
		origin   string
		wantACAO string
		wantACAC string
	}{
		{
			name:     "allowed origin with credentials",
			origin:   "https://example.com",
			wantACAO: "https://example.com",
			wantACAC: "true",
		},
		{
			name:     "disallowed origin no credentials",
			origin:   "https://evil.com",
			wantACAO: "",
			wantACAC: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			gotACAO := w.Header().Get("Access-Control-Allow-Origin")
			if gotACAO != tc.wantACAO {
				t.Errorf("ACAO = %q, want %q", gotACAO, tc.wantACAO)
			}
			gotACAC := w.Header().Get("Access-Control-Allow-Credentials")
			if gotACAC != tc.wantACAC {
				t.Errorf("ACAC = %q, want %q", gotACAC, tc.wantACAC)
			}
		})
	}
}

func TestValidateOriginHeader(t *testing.T) {
	tests := []struct {
		name           string
		origin         string
		allowedOrigins []string
		want           string
	}{
		{
			name:           "empty origin returns empty",
			origin:         "",
			allowedOrigins: []string{"https://example.com"},
			want:           "",
		},
		{
			name:           "exact match returns origin",
			origin:         "https://example.com",
			allowedOrigins: []string{"https://example.com"},
			want:           "https://example.com",
		},
		{
			name:           "wildcard returns wildcard",
			origin:         "https://example.com",
			allowedOrigins: []string{"*"},
			want:           "*",
		},
		{
			name:           "no match returns empty",
			origin:         "https://evil.com",
			allowedOrigins: []string{"https://example.com"},
			want:           "",
		},
		{
			name:           "case insensitive match",
			origin:         "HTTPS://EXAMPLE.COM",
			allowedOrigins: []string{"https://example.com"},
			want:           "HTTPS://EXAMPLE.COM",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := validateOriginHeader(tc.origin, tc.allowedOrigins)
			if got != tc.want {
				t.Errorf("validateOriginHeader(%q, %v) = %q, want %q", tc.origin, tc.allowedOrigins, got, tc.want)
			}
		})
	}
}
