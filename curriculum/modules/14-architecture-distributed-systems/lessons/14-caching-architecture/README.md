# Caching architecture

## Learning objective

Implement cache-aside, read-through, and write-behind caching strategies in Go, understand cache hierarchy in distributed systems, and measure hit rate, miss rate, and staleness tradeoffs.

## Why this matters

Caching is the single most effective performance optimization in distributed systems. A cache hit at L1 latency (~1ms) vs cache miss triggering a database query (~50ms) or an external API call (~500ms) is the difference between a snappy UI and a spinning loader. In production, every layer of the stack caches: CPU L1/L2/L3, application memory (Redis, in-memory map), CDN (CloudFront, Cloudflare), and database query cache. Understanding caching strategies is essential for any Go engineer building latency-sensitive services.

## Mental model

A cache is a fast, limited-size copy of a slower, authoritative data store. Think of a chef's counter (cache) vs the pantry (database). The chef keeps frequently used spices on the counter (cache-aside): if the spice is on the counter, use it. If not, go to the pantry, grab it, and put a refill on the counter. The counter has limited space -- less-used spices get evicted when new ones arrive.

In cache-aside (look-aside), the application code manages the cache. In read-through, the cache library automatically loads missing keys from the database. In write-behind, writes go to the cache first and are asynchronously flushed to the database. Write-through writes to both cache and database synchronously.

## Core idea

Caching strategies differ by who manages cache population and how writes propagate:

| Strategy | Read behavior | Write behavior | Consistency |
|---|---|---|---|
| Cache-aside | App checks cache, loads on miss | App invalidates cache on write | Eventual |
| Read-through | Cache loads from DB automatically | Cache invalidation on write | Eventual |
| Write-through | Same as read-through | Cache writes to DB synchronously | Strong |
| Write-behind | Same as read-through | Cache writes to DB asynchronously | Eventual |

Cache hierarchy uses multiple layers: L1 (in-process memory, fastest, smallest), L2 (Redis/memcached, fast, medium), L3 (database query cache, slowest, largest). The application checks L1 first, then L2, then falls back to the database.

## Under the hood

A cache-aside read flow:

1. Application calls `cache.Get(key)`.
2. Cache checks internal map. If found and not expired, return value (hit).
3. If not found (miss), application loads from database.
4. Application calls `cache.Set(key, value, ttl)`.
5. Return value to caller.

A cache-aside write flow:

1. Application writes to database.
2. Application calls `cache.Invalidate(key)`.
3. Next read for that key triggers a cache miss and reloads fresh data.

This invalidation-on-write approach ensures the cache does not serve stale data. The TTL acts as a safety net: even if invalidation is missed, the entry expires.

## How Go uses it

- **`hashicorp/golang-lru`**: a popular Go library for LRU (Least Recently Used) caches. Used by Consul, Vault, and Nomad.
- **`dgraph-io/ristretto`**: a high-performance Go cache with admission control (TinyLFU). Used by Dgraph and Badger.
- **`patrickmn/go-cache`**: an in-memory key:value store with TTL. Suitable for single-server applications.
- **`go-redis/redis`**: the de facto Go Redis client. Used for distributed caching with Redis cluster.
- **`groupcache`**: a distributed caching and cache-filling library from the Go authors. Used by Google.

## Go example

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type CacheEntry struct {
	Value     string
	ExpiresAt time.Time
}

type CacheAside struct {
	mu       sync.RWMutex
	entries  map[string]CacheEntry
	ttl      time.Duration
	loadFunc func(string) (string, error)
	hits     int
	misses   int
}

func NewCacheAside(ttl time.Duration, load func(string) (string, error)) *CacheAside {
	return &CacheAside{
		entries:  make(map[string]CacheEntry),
		ttl:      ttl,
		loadFunc: load,
	}
}

func (c *CacheAside) Get(key string) (string, error) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()

	if ok && time.Now().Before(entry.ExpiresAt) {
		c.mu.Lock()
		c.hits++
		c.mu.Unlock()
		return entry.Value, nil
	}

	c.mu.Lock()
	c.misses++
	c.mu.Unlock()

	value, err := c.loadFunc(key)
	if err != nil {
		return "", err
	}

	c.mu.Lock()
	c.entries[key] = CacheEntry{Value: value, ExpiresAt: time.Now().Add(c.ttl)}
	c.mu.Unlock()
	return value, nil
}

func (c *CacheAside) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

type WriteBehind struct {
	mu        sync.Mutex
	entries   map[string]string
	writeFunc func(string, string) error
}

func NewWriteBehind(write func(string, string) error) *WriteBehind {
	return &WriteBehind{entries: make(map[string]string), writeFunc: write}
}

func (w *WriteBehind) Set(key, value string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries[key] = value
}

func (w *WriteBehind) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	for k, v := range w.entries {
		if err := w.writeFunc(k, v); err != nil {
			return err
		}
		delete(w.entries, k)
	}
	return nil
}

func expensiveLoad(key string) (string, error) {
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("loaded-%s", key), nil
}

func main() {
	cache := NewCacheAside(5*time.Minute, expensiveLoad)

	val, _ := cache.Get("user:42")
	fmt.Println("first:", val)
	val, _ = cache.Get("user:42")
	fmt.Println("cached:", val)

	cache.Invalidate("user:42")
	fmt.Println("invalidated")

	wb := NewWriteBehind(func(k, v string) error {
		fmt.Printf("writing %s=%s\n", k, v)
		return nil
	})
	wb.Set("count", "100")
	wb.Set("total", "500")
	wb.Flush()
}
```

## Step-by-step execution

For `cache.Get("user:42")` on first call:

1. `c.mu.RLock()` acquires read lock.
2. Check `c.entries["user:42"]` -- not found (empty cache).
3. Release read lock.
4. Increment `c.misses`.
5. Call `c.loadFunc("user:42")` → calls `expensiveLoad` → sleeps 10ms → returns `"loaded-user:42"`.
6. Acquire write lock. Store `entries["user:42"] = {Value: "loaded-user:42", ExpiresAt: now + 5min}`.
7. Release write lock. Return `"loaded-user:42"`.

On second call:
1. Entry exists and is not expired.
2. Return `"loaded-user:42"` immediately. No `loadFunc` call. This is the cache hit.

## Common mistakes

- **Caching without TTL**: entries live forever. If the source data changes, the cache serves stale data forever. Always set a TTL, even if it is hours or days.
- **Cache invalidation on every write**: invalidating the cache entry on every write is correct but expensive if the entry is rarely read. Consider write-through caching for high-read, high-write data.
- **Cache stampede (thundering herd)**: when a popular key expires, N concurrent requests all trigger a cache miss and all hit the database simultaneously. Use request coalescing (singleflight) to allow only one request to load the value.
- **Distributed cache key collisions**: using short or ambiguous keys in a shared Redis instance causes cross-feature collisions. Prefix keys with the feature name: `user:profile:42` instead of `42`.
- **Ignoring cache hit ratio**: a cache with a 10% hit rate adds latency overhead without benefit. Measure hit rate and tune TTL, size, and eviction policy to achieve 90%+ for read-heavy workloads.

## Debugging walkthrough

A recommendation service is slow after deployment:

```go
func GetRecommendations(userID string) ([]Recommendation, error) {
	val, err := cache.Get(fmt.Sprintf("recs:%s", userID))
	if err == nil {
		return val, nil
	}
	return loadFromDB(userID)
}
```

**Symptom**: Response times are 500ms instead of the expected 50ms. Database CPU is at 100%.

**Investigation**: Check cache hit rate:

```go
stats := cache.Stats()
fmt.Printf("hits=%d misses=%d ratio=%.2f%%\n", stats.Hits, stats.Misses,
	float64(stats.Hits)/float64(stats.Hits+stats.Misses)*100)
```

**Root cause**: The cache TTL was accidentally set to 10 seconds during deployment config. After 10 seconds, every request misses and hits the database. The 1000 req/s overwhelms the database.

**Fix**: Set TTL to 10 minutes. The cache hit ratio returns to 95%. Response times drop to 50ms.

## Production notes

- **Cache warming**: on deployment, the cache is empty. Preload popular keys to avoid a cold-start storm. Start warming in a background goroutine before accepting traffic.
- **Monitoring**: track hit rate, miss rate, eviction count, cache size, and latency per Get/Set. Alert when hit rate drops below 80%.
- **Memory limits**: set a maximum cache size with an eviction policy (LRU, LFU, TTL). An unbounded in-memory cache causes OOM kills.
- **Serialization**: when caching structs in Redis, choose efficient serialization: protobuf or msgpack over JSON. JSON is 3-5x slower to parse.
- **Distributed caching consistency**: Redis cluster uses slot-based sharding. Writes are not immediately visible across all nodes. Be prepared for stale reads after a write.

## Performance implications

- **In-memory cache (L1)**: 100-500ns per Get. Sub-microsecond. Thread-safe with `sync.RWMutex`.
- **Redis cache (L2)**: 0.5-2ms per round trip (network). 100x slower than L1 but shared across instances.
- **Cache miss (database)**: 10-100ms. 1000x slower than L1. The cache turns a 100ms operation into a 1us operation.
- **Cache size vs GC**: a large in-memory map with millions of entries causes GC pressure because the GC scans every entry pointer. Use `ristretto` or `freecache` which minimize GC scanning by storing values off-heap.

## Practice task

Build a two-level cache hierarchy: L1 is an in-memory map with 1-second TTL, L2 is a simulated Redis (another map) with 60-second TTL. On a miss at L1, check L2. On a miss at L2, call a slow function `loadFromDB`. Return the value and populate both caches. Measure the latency of each level.

## Tests / verification

```bash
go run ./curriculum/modules/14-architecture-distributed-systems/lessons/14-caching-architecture
go test ./curriculum/modules/14-architecture-distributed-systems/lessons/14-caching-architecture
```

The tests verify cache hits/misses, TTL expiry, invalidation, write-behind flush, and stats tracking.

## Review questions

1. What is the difference between cache-aside and read-through caching?
2. Why is cache invalidation on write necessary for cache-aside but not for write-through?
3. What is a cache stampede and how can it be prevented?
4. Why might an in-memory cache with a 1-hour TTL still serve stale data?
5. Under what conditions would a write-behind cache lose data?

## NEXT UP

When to split services.
