# Caching basics

## Learning objective

Build an in-memory cache in Go with TTL expiry and concurrent access safety, choose between `sync.Map` and `map` with `sync.RWMutex` for different access patterns, implement LRU eviction, and understand read-through vs write-through cache strategies.

## Why this matters

Every Go HTTP service with a database should have some form of in-memory cache. A cache hit is microseconds; a database query is milliseconds — three orders of magnitude difference. Without caching, a service that handles 10,000 req/s with a 10 ms DB query spends 100 seconds of CPU per second on database calls alone. Adding a cache with a 90% hit rate reduces that to 10 seconds. Caching is the single highest-impact performance optimization for most database-backed services.

## Mental model

A cache is a fast temporary storage that sits between the application and a slow data source. Think of it as a drawer next to your desk. The drawer holds frequently used tools (data). Reaching into the drawer takes one second; walking to the toolbox takes five minutes. The drawer has limited space — when it fills up, the least-used tool goes back to the toolbox (LRU eviction). The key insight: caching only helps if the access pattern has locality — the same keys are requested repeatedly within the eviction window.

## Core idea

A cache stores the result of an expensive operation so that future requests for the same key can be served faster. The fundamental operations are:

| Operation | What it does | Complexity |
|---|---|---|
| `Get(key)` | Return cached value or signal miss | O(1) amortized |
| `Set(key, value)` | Store a value, possibly with TTL | O(1) amortized |
| `Delete(key)` | Remove a value | O(1) amortized |
| Eviction | Remove old entries to free space | Varies by policy |

The three most common cache patterns in Go:

- **Cache-aside (lazy loading)**: Application checks cache first. On miss, reads from DB, writes to cache, returns. On write, updates DB and deletes from cache.
- **Read-through**: Cache automatically fetches from DB on miss. Application only talks to the cache.
- **Write-through**: Cache and DB are updated atomically on every write. Reads always hit the cache.

## Under the hood

A Go `map` used as a cache is a hash table with O(1) average lookup, but it is not goroutine-safe. Wrapping with `sync.RWMutex` provides safe concurrent access: `RLock` for reads, `Lock` for writes. `sync.Map` uses a combination of atomic loads for read-mostly workloads and a mutex-protected "dirty" map for writes, optimized for two patterns: append-only writes with many readers, and write-once/read-many.

LRU eviction requires a doubly-linked list plus a map: on access, the entry is moved to the front of the list; on eviction, the entry at the back is removed. This is O(1) per operation but requires pointer management. The `hashicorp/golang-lru` package implements this.

TTL expiration uses lazy eviction (check on `Get`) plus a periodic cleanup goroutine. Lazy eviction is cheaper per-op but can leave stale entries in memory longer. Periodic cleanup removes stale entries in batches but adds a mutex contention spike every N seconds. For very large caches, use a sharded design: divide the key space across N independent caches, each with its own lock, reducing lock contention by Nx.

## How Go uses it

- The standard library's `net/http` uses an internal cache for parsed HTTP header keys.
- `encoding/json` caches `reflect.Type` metadata for struct tags.
- The Go compiler uses a cache for `buildid` lookups and import graph resolution.
- Kubernetes API server uses an in-memory cache backed by etcd watches.
- Most Go web frameworks (Gin, Echo, Fiber) provide middleware for response caching.

## Go example

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	value  interface{}
	expiry time.Time
}

type TTLCache struct {
	mu       sync.RWMutex
	entries  map[string]CacheEntry
	ttl      time.Duration
	stopChan chan struct{}
}

func NewTTLCache(ttl time.Duration, cleanupInterval time.Duration) *TTLCache {
	c := &TTLCache{
		entries:  make(map[string]CacheEntry),
		ttl:      ttl,
		stopChan: make(chan struct{}),
	}
	go c.cleanupLoop(cleanupInterval)
	return c
}

func (c *TTLCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = CacheEntry{
		value:  value,
		expiry: time.Now().Add(c.ttl),
	}
}

func (c *TTLCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.expiry) {
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()
		return nil, false
	}
	return entry.value, true
}

func (c *TTLCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

func (c *TTLCache) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopChan:
			return
		}
	}
}

func (c *TTLCache) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, entry := range c.entries {
		if now.After(entry.expiry) {
			delete(c.entries, k)
		}
	}
}

func (c *TTLCache) Stop() {
	close(c.stopChan)
}

func main() {
	cache := NewTTLCache(100*time.Millisecond, 50*time.Millisecond)
	defer cache.Stop()

	cache.Set("name", "Alice")
	if val, ok := cache.Get("name"); ok {
		fmt.Println("Got:", val)
	} else {
		fmt.Println("Miss")
	}

	time.Sleep(200 * time.Millisecond)

	if _, ok := cache.Get("name"); !ok {
		fmt.Println("Expired after TTL")
	}

	cache.Set("count", 42)
	fmt.Println("Cache size:", cache.Len())
}
```

## Step-by-step execution

1. `NewTTLCache(100ms, 50ms)` creates a cache with 100 ms TTL and a cleanup goroutine that runs every 50 ms.
2. `cache.Set("name", "Alice")` acquires the write lock and stores a `CacheEntry` with expiry = `now + 100ms`.
3. `cache.Get("name")` acquires the read lock, finds the entry, checks expiry (not yet expired), and returns `"Alice"`.
4. `time.Sleep(200ms)` — the cleanup goroutine runs after 50 ms and 100 ms, deleting expired entries.
5. `cache.Get("name")` after sleep: the read lock finds the entry, but the expiry check fails (TTL expired). The write lock is acquired to delete the entry. Returns `nil, false`.
6. `cache.Set("count", 42)` stores a new entry.
7. `cache.Stop()` closes the stop channel, causing the cleanup goroutine to exit its loop.

## Common mistakes

- **Using `sync.Map` for read-heavy workloads with simple keys** — `sync.Map` is optimized for two specific patterns. For a standard read-heavy workload (90% reads, single writer), a regular `map` with `sync.RWMutex` is 2-5x faster because `sync.Map`'s atomic `Load` overhead outweighs mutex contention at low concurrency.
- **Adding a cache without a TTL or eviction policy** — An unbounded cache grows until the process OOMs. Every cache must have: max size (item count or byte limit), eviction policy (LRU, LFU, TTL), and a plan for stale data.
- **Caching at the wrong layer** — Cache the expensive-to-compute value, not the raw data. If the expensive work is `json.Unmarshal`, cache the parsed Go struct, not the raw JSON bytes.
- **Not handling cache stampede** — When a popular key expires and multiple requests simultaneously try to refresh it, they all hit the DB. Use singleflight (`golang.org/x/sync/singleflight`) or a mutex per key.

## Debugging walkthrough

A service's response time increased after adding a cache. The cache hit rate is 80%, but response times are higher than without caching:

Check cache serialization cost:

```go
func BenchmarkCacheGetSerialization(b *testing.B) {
	cache := NewTTLCache(time.Minute, time.Minute)
	cache.Set("key", generateLargeStruct())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, ok := cache.Get("key")
		if !ok {
			b.Fatal("miss")
		}
	}
}
```

If the benchmark shows high `ns/op`, the problem is in the cache implementation (e.g., deep copy on `Get`, JSON serialization). The fix: store pointers instead of values, or pre-serialize to bytes.

If the hit rate is low, use `cache.Set` to track how many times each key is set vs missed:

```go
cache.mu.Lock()
cache.missCount++
cache.mu.Unlock()
```

A miss rate >20% suggests the TTL is too short or the working set exceeds the cache size.

## Production notes

- Set `GOMAXPROCS` to match the number of CPU cores. Cache sharding (N shards × 1 lock each) provides the best throughput.
- Monitor cache hit rate, size, and eviction count as metrics (`expvar`, Prometheus).
- Use `hashicorp/golang-lru` for LRU instead of reimplementing it. The library is battle-tested.
- In Kubernetes, cache size should be a fraction of the container memory limit (e.g., 25%).
- For distributed caching, use Redis with `go-redis` and configure connection pooling.

## Performance implications

The performance of an in-memory cache depends on:

- **Lock contention**: A single mutex protects the entire cache. At high concurrency (>1000 goroutines), lock contention dominates. Use sharding: N shards, each with its own mutex, partitioned by `hash(key) % N`.
- **Allocation**: Each `Set` allocates a `CacheEntry` struct and stores the key (string). For high-throughput caches, pre-allocate or use `sync.Pool`.
- **Eviction cost**: Periodic cleanup scans the entire cache. For 1 million entries, a full scan takes 10-50 ms. Use probabilistic expiration (each entry has a random TTL) to spread cleanup work.

A well-tuned in-memory cache has `Get` latency under 100 ns and `Set` latency under 200 ns (excluding allocation). For comparison, a Redis `GET` via loopback takes 500 µs — 5000x slower.

## Practice task

Extend the `TTLCache` implementation to support:

1. `GetOrSet(key string, ttl time.Duration, fn func() (interface{}, error)) (interface{}, error)` — atomically get or compute-and-set.
2. A maximum size limit with LRU eviction (remove the entry with the earliest expiry when the cache exceeds capacity).

Write tests for the LRU eviction behavior.

## Tests / verification

```bash
go run ./curriculum/modules/13-performance-memory-engineering/lessons/07-caching-basics
go test ./curriculum/modules/13-performance-memory-engineering/lessons/07-caching-basics
```

## Review questions

1. What is the difference between cache-aside and read-through caching?
2. When would you choose `sync.Map` over `map` with `sync.RWMutex`?
3. Why is it dangerous to have a cache with no eviction policy or TTL?
4. What is cache stampede, and how does singleflight prevent it?
5. How does sharding reduce lock contention in a concurrent cache?

## NEXT UP

Cache invalidation — strategies for keeping cached data consistent with the source of truth, including TTL-based, write-invalidate, write-update, cache stampede prevention, and stale-while-revalidate patterns.
