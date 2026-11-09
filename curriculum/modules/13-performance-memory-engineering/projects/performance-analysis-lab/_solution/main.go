package main

import (
	"fmt"
	"hash/fnv"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"
)

const (
	NumShards      = 64
	WindowDuration = time.Second
)

// TokenResult represents the outcome of a rate-limit check.
// Returned by value to guarantee stack allocation and zero heap escapes.
type TokenResult struct {
	Allowed   bool
	Remaining int
	ResetIn   time.Duration
}

// CircularRingWindow maintains fixed-capacity sliding window timestamp records without allocations.
type CircularRingWindow struct {
	timestamps []int64 // Unix nano timestamps
	head       int
	count      int
	capacity   int
}

func NewCircularRingWindow(capacity int) *CircularRingWindow {
	return &CircularRingWindow{
		timestamps: make([]int64, capacity),
		capacity:   capacity,
	}
}

func (w *CircularRingWindow) Allow(now int64, windowNano int64) (bool, int) {
	threshold := now - windowNano

	// Purge expired entries from the circular window
	validCount := 0
	for i := 0; i < w.count; i++ {
		idx := (w.head - w.count + i + w.capacity) % w.capacity
		if w.timestamps[idx] > threshold {
			validCount = w.count - i
			break
		}
	}
	w.count = validCount

	if w.count >= w.capacity {
		return false, 0
	}

	w.timestamps[w.head] = now
	w.head = (w.head + 1) % w.capacity
	w.count++

	remaining := w.capacity - w.count
	return true, remaining
}

type clientShard struct {
	mu      sync.Mutex
	buckets map[string]*CircularRingWindow
}

// ShardedRateLimiter partitions the keyspace across 64 independent mutex shards,
// virtually eliminating lock contention in high-throughput concurrent workloads.
type ShardedRateLimiter struct {
	shards    [NumShards]clientShard
	burstCap  int
	windowDur time.Duration
}

func NewShardedRateLimiter(burstCapacity int, window time.Duration) *ShardedRateLimiter {
	limiter := &ShardedRateLimiter{
		burstCap:  burstCapacity,
		windowDur: window,
	}
	for i := 0; i < NumShards; i++ {
		limiter.shards[i].buckets = make(map[string]*CircularRingWindow, 256)
	}
	return limiter
}

func (l *ShardedRateLimiter) getShard(key string) *clientShard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	idx := h.Sum32() % NumShards
	return &l.shards[idx]
}

// CheckLimit checks if the given client key is permitted under the rate limit.
// Designed with zero-allocation semantics: no heap escapes on return.
func (l *ShardedRateLimiter) CheckLimit(key string) TokenResult {
	shard := l.getShard(key)
	now := time.Now().UnixNano()
	windowNano := l.windowDur.Nanoseconds()

	shard.mu.Lock()
	defer shard.mu.Unlock()

	bucket, exists := shard.buckets[key]
	if !exists {
		bucket = NewCircularRingWindow(l.burstCap)
		shard.buckets[key] = bucket
	}

	allowed, remaining := bucket.Allow(now, windowNano)

	return TokenResult{
		Allowed:   allowed,
		Remaining: remaining,
		ResetIn:   l.windowDur,
	}
}

func main() {
	limiter := NewShardedRateLimiter(5, time.Second)

	fmt.Println("Sharded Rate Limiter Demonstration:")
	for i := 1; i <= 7; i++ {
		res := limiter.CheckLimit("user-123")
		fmt.Printf("Request %d: Allowed=%v, Remaining=%d\n", i, res.Allowed, res.Remaining)
	}

	// Mount pprof for profiling
	go func() {
		_ = http.ListenAndServe(":6061", nil)
	}()
	_ = os.Stdout.Sync()
}
