package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCSRFTokenGeneration(t *testing.T) {
	p := NewCSRFProtector()

	token1 := p.GenerateToken("session-1")
	token2 := p.GenerateToken("session-1")

	if token1 == token2 {
		t.Error("expected different tokens for same session")
	}
	if len(token1) != 64 {
		t.Errorf("expected token length 64 hex chars, got %d", len(token1))
	}
}

func TestCSRFValidation(t *testing.T) {
	p := NewCSRFProtector()
	sessionID := "session-abc"

	token := p.GenerateToken(sessionID)

	tests := []struct {
		name      string
		token     string
		sessionID string
		want      bool
	}{
		{
			name:      "valid token for correct session",
			token:     token,
			sessionID: sessionID,
			want:      true,
		},
		{
			name:      "valid token for wrong session",
			token:     token,
			sessionID: "session-different",
			want:      false,
		},
		{
			name:      "empty token",
			token:     "",
			sessionID: sessionID,
			want:      false,
		},
		{
			name:      "tampered token",
			token:     token + "ff",
			sessionID: sessionID,
			want:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := p.ValidateToken(tc.token, tc.sessionID)
			if got != tc.want {
				t.Errorf("ValidateToken(%q, %q) = %v, want %v", tc.token, tc.sessionID, got, tc.want)
			}
		})
	}
}

func TestCSRFTokenOneTimeUse(t *testing.T) {
	p := NewCSRFProtector()
	sessionID := "session-one-time"

	token := p.GenerateToken(sessionID)

	if !p.ValidateToken(token, sessionID) {
		t.Error("expected first validation to succeed")
	}

	if p.ValidateToken(token, sessionID) {
		t.Error("expected second validation to fail (one-time use)")
	}
}

func TestCSRFMiddlewareOriginCheck(t *testing.T) {
	handler := middlewareCSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), NewCSRFProtector())

	tests := []struct {
		name       string
		method     string
		headers    map[string]string
		wantStatus int
	}{
		{
			name:       "GET request passes",
			method:     "GET",
			headers:    map[string]string{},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/", nil)
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tc.wantStatus)
			}
		})
	}
}
