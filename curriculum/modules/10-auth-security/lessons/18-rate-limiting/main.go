package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

type TokenBucket struct {
	mu         sync.Mutex
	capacity   int
	tokens     int
	refillRate float64
	lastRefill time.Time
}

func NewTokenBucket(capacity int, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.lastRefill = now

	tb.tokens += int(elapsed * tb.refillRate)
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

type RateLimiter struct {
	mu       sync.RWMutex
	buckets  map[string]*TokenBucket
	capacity int
	rate     float64
}

func NewRateLimiter(capacity int, rate float64) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*TokenBucket),
		capacity: capacity,
		rate:     rate,
	}
}

func (rl *RateLimiter) getBucket(key string) *TokenBucket {
	rl.mu.RLock()
	bucket, ok := rl.buckets[key]
	rl.mu.RUnlock()
	if ok {
		return bucket
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, ok = rl.buckets[key]
	if ok {
		return bucket
	}

	bucket = NewTokenBucket(rl.capacity, rl.rate)
	rl.buckets[key] = bucket
	return bucket
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		bucket := rl.getBucket(ip)

		if !bucket.Allow() {
			w.Header().Set("Retry-After", "1")
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.capacity))
			w.Header().Set("X-RateLimit-Remaining", "0")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.capacity))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", bucket.Tokens()))
		next.ServeHTTP(w, r)
	})
}

func (tb *TokenBucket) Tokens() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.tokens
}

func main() {
	fmt.Println("=== Token Bucket Rate Limiter ===")

	bucket := NewTokenBucket(3, 1.0)
	fmt.Printf("Capacity: %d, Refill rate: %.0f token/sec\n", bucket.capacity, bucket.refillRate)

	fmt.Println("\nAllowed requests:")
	for i := 0; i < 10; i++ {
		allowed := bucket.Allow()
		fmt.Printf("  Request %d: allowed=%v, tokens=%d\n", i+1, allowed, bucket.tokens)
		if !allowed {
			break
		}
	}

	fmt.Println("\n=== Rate Limiter Middleware ===")
	limiter := NewRateLimiter(2, 1.0)
	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	ts := httptest.NewServer(limiter.Middleware(mux))
	defer ts.Close()

	for i := 0; i < 5; i++ {
		resp, _ := http.Get(ts.URL + "/api")
		fmt.Printf("  Request %d: status=%d, remaining=%s\n",
			i+1, resp.StatusCode, resp.Header.Get("X-RateLimit-Remaining"))
		resp.Body.Close()
		time.Sleep(200 * time.Millisecond)
	}
}
