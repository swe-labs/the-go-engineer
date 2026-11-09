package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// TokenResult represents the outcome of a rate-limit check.
type TokenResult struct {
	Allowed   bool
	Remaining int
	ResetIn   time.Duration
}

// ClientBucket maintains request timestamps for a client.
// UNOPTIMIZED: Appends to unbounded slice, causing memory leak and garbage collection pressure.
type ClientBucket struct {
	timestamps []time.Time
}

// GlobalRateLimiter
// UNOPTIMIZED: Uses a single global mutex across all keys, creating massive lock contention under concurrency.
// UNOPTIMIZED: CheckLimit allocates *TokenResult on heap instead of returning by value.
// TODO: Task 1 - Refactor to sharded mutex array (e.g. 64 shards partitioned by FNV key hash)
// TODO: Task 2 - Replace unbounded slice with fixed-capacity circular ring buffer
// TODO: Task 3 - Return TokenResult by value to eliminate heap allocation escapes
type ShardedRateLimiter struct {
	mu        sync.Mutex
	buckets   map[string]*ClientBucket
	burstCap  int
	windowDur time.Duration
}

func NewShardedRateLimiter(burstCapacity int, window time.Duration) *ShardedRateLimiter {
	return &ShardedRateLimiter{
		buckets:   make(map[string]*ClientBucket),
		burstCap:  burstCapacity,
		windowDur: window,
	}
}

// CheckLimit checks if the client key is permitted.
func (l *ShardedRateLimiter) CheckLimit(key string) *TokenResult {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	threshold := now.Add(-l.windowDur)

	bucket, exists := l.buckets[key]
	if !exists {
		bucket = &ClientBucket{timestamps: make([]time.Time, 0)}
		l.buckets[key] = bucket
	}

	// Purge expired
	valid := make([]time.Time, 0, len(bucket.timestamps))
	for _, t := range bucket.timestamps {
		if t.After(threshold) {
			valid = append(valid, t)
		}
	}
	bucket.timestamps = valid

	if len(bucket.timestamps) >= l.burstCap {
		return &TokenResult{
			Allowed:   false,
			Remaining: 0,
			ResetIn:   l.windowDur,
		}
	}

	bucket.timestamps = append(bucket.timestamps, now)
	remaining := l.burstCap - len(bucket.timestamps)

	return &TokenResult{
		Allowed:   true,
		Remaining: remaining,
		ResetIn:   l.windowDur,
	}
}

func main() {
	limiter := NewShardedRateLimiter(5, time.Second)

	fmt.Println("Legacy Unoptimized Rate Limiter Demonstration:")
	for i := 1; i <= 7; i++ {
		res := limiter.CheckLimit("user-123")
		fmt.Printf("Request %d: Allowed=%v, Remaining=%d\n", i, res.Allowed, res.Remaining)
	}

	go func() {
		_ = http.ListenAndServe(":6061", nil)
	}()
	_ = os.Stdout.Sync()
}
