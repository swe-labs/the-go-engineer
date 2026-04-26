package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
	"time"
)

func TestCORSHandlesPreflight(t *testing.T) {
	t.Parallel()

	handler := CORS(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler should not run for preflight")
	}))
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/orders", nil)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNoContent)
	}

	if res.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing CORS allow-origin header")
	}
}

func TestRateLimitRejectsRequestsOverLimit(t *testing.T) {
	t.Parallel()

	handler := RateLimit(1, time.Minute, nil)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	first := httptest.NewRecorder()
	handler.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil))
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, want %d", first.Code, http.StatusOK)
	}

	second := httptest.NewRecorder()
	handler.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil))
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want %d", second.Code, http.StatusTooManyRequests)
	}

	if second.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("content type = %q, want application/json", second.Header().Get("Content-Type"))
	}

	wantBody := `{"error":{"code":"rate_limited","message":"rate limit exceeded"}}`
	if second.Body.String() != wantBody {
		t.Fatalf("body = %q, want %q", second.Body.String(), wantBody)
	}
}

func TestRateLimitUsesForwardedClientIP(t *testing.T) {
	t.Parallel()

	handler := RateLimit(1, time.Minute, []netip.Prefix{netip.MustParsePrefix("10.0.0.10/32")})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	firstReq := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	firstReq.RemoteAddr = "10.0.0.10:1234"
	firstReq.Header.Set("X-Forwarded-For", "203.0.113.10")
	first := httptest.NewRecorder()
	handler.ServeHTTP(first, firstReq)

	secondReq := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	secondReq.RemoteAddr = "10.0.0.10:5678"
	secondReq.Header.Set("X-Forwarded-For", "203.0.113.11")
	second := httptest.NewRecorder()
	handler.ServeHTTP(second, secondReq)

	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, want %d", first.Code, http.StatusOK)
	}

	if second.Code != http.StatusOK {
		t.Fatalf("second status = %d, want %d", second.Code, http.StatusOK)
	}
}

func TestRateLimitIgnoresForwardedHeadersFromUntrustedPeer(t *testing.T) {
	t.Parallel()

	handler := RateLimit(1, time.Minute, []netip.Prefix{netip.MustParsePrefix("127.0.0.1/32")})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	firstReq := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	firstReq.RemoteAddr = "198.51.100.20:1234"
	firstReq.Header.Set("X-Forwarded-For", "203.0.113.10")
	first := httptest.NewRecorder()
	handler.ServeHTTP(first, firstReq)

	secondReq := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	secondReq.RemoteAddr = "198.51.100.20:5678"
	secondReq.Header.Set("X-Forwarded-For", "203.0.113.11")
	second := httptest.NewRecorder()
	handler.ServeHTTP(second, secondReq)

	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, want %d", first.Code, http.StatusOK)
	}

	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want %d", second.Code, http.StatusTooManyRequests)
	}
}

func TestPruneExpiredClientWindows(t *testing.T) {
	t.Parallel()

	now := time.Now()
	clients := map[string]clientWindow{
		"expired": {count: 1, resetAt: now.Add(-time.Second)},
		"active":  {count: 1, resetAt: now.Add(time.Second)},
	}

	pruneExpiredClientWindows(clients, now)

	if _, ok := clients["expired"]; ok {
		t.Fatal("expired client window should be evicted")
	}

	if _, ok := clients["active"]; !ok {
		t.Fatal("active client window should remain")
	}
}
