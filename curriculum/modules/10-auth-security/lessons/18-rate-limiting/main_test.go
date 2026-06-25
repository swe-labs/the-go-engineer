package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestTokenBucketInitialTokens(t *testing.T) {
	bucket := NewTokenBucket(10, 5.0)
	if bucket.Tokens() != 10 {
		t.Errorf("expected 10 initial tokens, got %d", bucket.Tokens())
	}
}

func TestTokenBucketAllowsUpToCapacity(t *testing.T) {
	bucket := NewTokenBucket(3, 1.0)

	for i := 0; i < 3; i++ {
		if !bucket.Allow() {
			t.Errorf("expected request %d to be allowed", i+1)
		}
	}

	if bucket.Allow() {
		t.Error("expected 4th request to be denied")
	}
}

func TestTokenBucketRefill(t *testing.T) {
	bucket := NewTokenBucket(3, 10.0)

	bucket.Allow()
	bucket.Allow()
	bucket.Allow()

	_ = bucket.Allow()

	if bucket.Allow() {
		t.Error("expected request after capacity to be denied without refill")
	}
}

func TestRateLimiterMiddleware(t *testing.T) {
	limiter := NewRateLimiter(2, 10.0)
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	var mu sync.Mutex
	statuses := make([]int, 0, 3)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			mu.Lock()
			statuses = append(statuses, rec.Code)
			mu.Unlock()
		}()
	}
	wg.Wait()

	allowed := 0
	denied := 0
	for _, s := range statuses {
		if s == http.StatusOK {
			allowed++
		}
		if s == http.StatusTooManyRequests {
			denied++
		}
	}

	if allowed != 2 {
		t.Errorf("expected 2 allowed, got %d", allowed)
	}
	if denied != 1 {
		t.Errorf("expected 1 denied, got %d", denied)
	}
}

func TestRateLimiterHeaders(t *testing.T) {
	limiter := NewRateLimiter(5, 10.0)
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("X-RateLimit-Limit") == "" {
		t.Error("expected X-RateLimit-Limit header")
	}
	if rec.Header().Get("X-RateLimit-Remaining") == "" {
		t.Error("expected X-RateLimit-Remaining header")
	}
}
