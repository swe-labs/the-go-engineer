package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareSetsDynamicHeadersOnAllow(t *testing.T) {
	limiter := New(Config{BurstSize: 2, WindowSeconds: 60})
	mw := Middleware(limiter, func(*http.Request) string { return "client-a" }, nil)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("X-RateLimit-Limit"); got != "2" {
		t.Fatalf("limit header=%q want %q", got, "2")
	}
	if got := rec.Header().Get("X-RateLimit-Remaining"); got != "1" {
		t.Fatalf("remaining header=%q want %q", got, "1")
	}
}

func TestMiddlewareReturnsJSONErrorOnBlock(t *testing.T) {
	limiter := New(Config{BurstSize: 1, WindowSeconds: 60})
	mw := Middleware(limiter, func(*http.Request) string { return "client-b" }, nil)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil))

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status=%d want=%d", rec.Code, http.StatusTooManyRequests)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content-type=%q want application/json", got)
	}
	want := "{\"error\":{\"code\":\"rate_limited\",\"message\":\"rate limit exceeded\"}}\n"
	if rec.Body.String() != want {
		t.Fatalf("body=%q want=%q", rec.Body.String(), want)
	}
	if rec.Header().Get("X-RateLimit-Limit") != "1" {
		t.Fatalf("missing dynamic limit header")
	}
}
