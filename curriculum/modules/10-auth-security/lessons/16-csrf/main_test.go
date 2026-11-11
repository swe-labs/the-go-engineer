package main

import (
	"net/http"
	"net/http/httptest"
	"net/netip"
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
	}))

	tests := []struct {
		name       string
		method     string
		target     string
		headers    map[string]string
		wantStatus int
	}{
		{
			name:       "GET request passes",
			method:     http.MethodGet,
			target:     "http://example.com/",
			headers:    map[string]string{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "POST missing Origin and Referer fails",
			method:     http.MethodPost,
			target:     "http://example.com/",
			headers:    map[string]string{},
			wantStatus: http.StatusForbidden,
		},
		{
			name:   "POST matching http Origin passes",
			method: http.MethodPost,
			target: "http://example.com/",
			headers: map[string]string{
				"Origin": "http://example.com",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "POST matching https Origin passes",
			method: http.MethodPost,
			target: "https://example.com/",
			headers: map[string]string{
				"Origin": "https://example.com",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "POST cross-site Origin fails",
			method: http.MethodPost,
			target: "https://example.com/",
			headers: map[string]string{
				"Origin": "https://evil.example",
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:   "POST prefix-spoofed Origin fails",
			method: http.MethodPost,
			target: "http://example.com/",
			headers: map[string]string{
				"Origin": "http://example.com.evil.example",
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:   "POST scheme-mismatched Origin fails",
			method: http.MethodPost,
			target: "https://example.com/",
			headers: map[string]string{
				"Origin": "http://example.com",
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:   "POST null Origin fails",
			method: http.MethodPost,
			target: "http://example.com/",
			headers: map[string]string{
				"Origin": "null",
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:   "POST matching Referer passes when Origin absent",
			method: http.MethodPost,
			target: "https://example.com/",
			headers: map[string]string{
				"Referer": "https://example.com/transfer?step=1",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "POST cross-site Referer fails",
			method: http.MethodPost,
			target: "https://example.com/",
			headers: map[string]string{
				"Referer": "https://evil.example/attack",
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:   "POST untrusted Origin fails even with matching Referer",
			method: http.MethodPost,
			target: "http://example.com/",
			headers: map[string]string{
				"Origin":  "https://evil.example",
				"Referer": "http://example.com/form",
			},
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.target, nil)
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

func TestCSRFMiddlewareBehindTLSTermination(t *testing.T) {
	okNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("configured public origins", func(t *testing.T) {
		handler := middlewareCSRF(okNext, OriginConfig{
			TrustedOrigins: []string{"https://app.example.com"},
		})
		tests := []struct {
			name       string
			origin     string
			referer    string
			wantStatus int
		}{
			{name: "https origin allowed even with nil r.TLS and internal host", origin: "https://app.example.com", wantStatus: http.StatusOK},
			{name: "https referer allowed when Origin absent", referer: "https://app.example.com/transfer", wantStatus: http.StatusOK},
			{name: "http downgrade rejected", origin: "http://app.example.com", wantStatus: http.StatusForbidden},
			{name: "internal backend host origin rejected", origin: "http://backend:8080", wantStatus: http.StatusForbidden},
			{name: "evil origin rejected", origin: "https://evil.example", wantStatus: http.StatusForbidden},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "http://backend:8080/transfer", nil)
				if tc.origin != "" {
					req.Header.Set("Origin", tc.origin)
				}
				if tc.referer != "" {
					req.Header.Set("Referer", tc.referer)
				}
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, req)
				if w.Code != tc.wantStatus {
					t.Errorf("status = %d, want %d", w.Code, tc.wantStatus)
				}
			})
		}
	})

	t.Run("invalid configured origins fail closed", func(t *testing.T) {
		handler := middlewareCSRF(okNext, OriginConfig{
			TrustedOrigins: []string{"://not-a-valid-origin"},
		})
		req := httptest.NewRequest(http.MethodPost, "http://example.com/", nil)
		req.Header.Set("Origin", "http://example.com")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
		}
	})

	t.Run("trusted proxies X-Forwarded-Proto", func(t *testing.T) {
		handler := middlewareCSRF(okNext, OriginConfig{
			TrustedProxies: []netip.Prefix{netip.MustParsePrefix("127.0.0.1/32"), netip.MustParsePrefix("10.0.0.0/8")},
		})
		tests := []struct {
			name       string
			remoteAddr string
			xfProto    string
			origin     string
			wantStatus int
		}{
			{
				name:       "trusted proxy with https proto allows https Origin",
				remoteAddr: "10.0.0.5:43210",
				xfProto:    "https",
				origin:     "https://app.example.com",
				wantStatus: http.StatusOK,
			},
			{
				name:       "IPv4-mapped IPv6 trusted proxy allows https Origin",
				remoteAddr: "[::ffff:127.0.0.1]:43210",
				xfProto:    "https, http",
				origin:     "https://app.example.com",
				wantStatus: http.StatusOK,
			},
			{
				name:       "trusted proxy rejects http Origin when X-Forwarded-Proto is https",
				remoteAddr: "10.0.0.5:43210",
				xfProto:    "https",
				origin:     "http://app.example.com",
				wantStatus: http.StatusForbidden,
			},
			{
				name:       "untrusted remote addr cannot spoof X-Forwarded-Proto",
				remoteAddr: "203.0.113.99:5555",
				xfProto:    "https",
				origin:     "https://app.example.com",
				wantStatus: http.StatusForbidden,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "http://app.example.com/transfer", nil)
				req.RemoteAddr = tc.remoteAddr
				if tc.xfProto != "" {
					req.Header.Set("X-Forwarded-Proto", tc.xfProto)
				}
				req.Header.Set("Origin", tc.origin)
				w := httptest.NewRecorder()
				handler.ServeHTTP(w, req)
				if w.Code != tc.wantStatus {
					t.Errorf("status = %d, want %d", w.Code, tc.wantStatus)
				}
			})
		}
	})
}
