# Cache invalidation

## Learning objective

Implement cache invalidation strategies in Go, including TTL-based expiry, write-invalidate, write-update, and stale-while-revalidate patterns, and prevent cache stampede using singleflight and probabilistic early expiration.

## Why this matters

Cache invalidation is one of the two hard things in computer science (along with naming and off-by-one errors). A cache with stale data is worse than no cache at all — it serves incorrect data silently. Invalidation strategy determines the staleness window, consistency guarantees, and operational complexity of a caching layer. Choosing the wrong strategy leads to data inconsistency bugs that are subtle, hard to reproduce, and expensive to fix.

## Mental model

Cache invalidation is the problem of keeping a copy (cache) consistent with the source of truth (DB). Imagine a whiteboard that lists office phone numbers. Someone updates the central directory (DB) but forgets to erase the old number on the whiteboard (cache). Callers reach the wrong person. If the erasure is forgotten, the whiteboard shows old data. If the erasure happens before the directory update, the whiteboard is empty even though the directory update failed. If two people erase and write simultaneously, the final whiteboard state may match neither intended update. The three mechanisms are: TTL (the whiteboard self-erases after a timeout), active invalidation (someone erases when the directory changes), and write-through (someone updates the whiteboard as part of the directory update). Each has tradeoffs.

## Core idea

The four main cache invalidation strategies:

| Strategy | Mechanism | Staleness window | Complexity |
|---|---|---|---|
| TTL-only | Cache entries expire after a fixed duration | Up to TTL duration | Low |
| Write-invalidate | On DB write, delete the cache entry | Milliseconds | Medium |
| Write-update | On DB write, update the cache entry | Zero (theoretical) | High |
| Stale-while-revalidate | Serve stale data while refreshing in background | Configurable | Medium |

Each strategy is appropriate for different consistency requirements:

- **TTL-only**: Social media feeds, top-10 lists, configuration data that changes rarely.
- **Write-invalidate**: Most CRUD APIs where consistency matters but eventual consistency is acceptable.
- **Write-update**: Systems requiring read-your-writes consistency (user profile updates).
- **Stale-while-revalidate**: High-read/low-write data where latency matters more than perfect freshness.

## Under the hood

In Go, cache invalidation typically uses one of three patterns:

1. **Cache-aside with delete**: The application reads from cache first. On cache miss, reads from DB and populates the cache. On writes, it updates the DB and deletes the cache entry (or updates it). This is the most common pattern because it handles concurrent writes safely — deleting is idempotent, and the next read fetches fresh data.

2. **Write-through**: The application updates both cache and DB atomically. Go has no built-in distributed transaction support, so this requires coordinating a two-phase commit between cache (Redis) and DB (PostgreSQL) — typically implemented with the outbox pattern: write to a DB table, a separate process reads the outbox and updates the cache.

3. **Cache stampede prevention**: When a popular key expires and multiple requests simultaneously try to refresh it, they all hit the DB. `golang.org/x/sync/singleflight` coalesces concurrent calls: only one goroutine performs the DB query; the rest wait for its result. Probabilistic early expiration (PEE) adds jitter to TTLs: each entry expires at `baseTTL + rand(0, jitterWindow)` so not all entries expire simultaneously.

## How Go uses it

- `golang.org/x/sync/singleflight` is the standard library-adjacent package for coalescing concurrent cache refresh requests.
- The `hashicorp/golang-lru` package provides LRU eviction with optional `evicted` callbacks for write-invalidate patterns.
- Redis (`go-redis`) supports SET with `NX` (set if not exists) for distributed write-invalidate and EXPIRE for TTL.
- etcd's watch mechanism is the foundation for Kubernetes' cache invalidation: clients watch for resource changes and invalidate local caches.
- The Go community's `groupcache` implements a distributed cache with consistent hashing and automatic invalidation via peer communication.

## Go example

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type staleEntry struct {
	value      interface{}
	expiry     time.Time
	staleUntil time.Time
}

type StaleWhileRevalidateCache struct {
	mu       sync.RWMutex
	entries  map[string]staleEntry
	ttl      time.Duration
	staleTTL time.Duration
	renew    func(key string) (interface{}, error)
}

func NewStaleWhileRevalidateCache(ttl, staleTTL time.Duration, renew func(key string) (interface{}, error)) *StaleWhileRevalidateCache {
	return &StaleWhileRevalidateCache{
		entries:  make(map[string]staleEntry),
		ttl:      ttl,
		staleTTL: staleTTL,
		renew:    renew,
	}
}

func (c *StaleWhileRevalidateCache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()

	if ok {
		if time.Now().Before(entry.expiry) {
			return entry.value, nil
		}
		if time.Now().Before(entry.staleUntil) {
			go c.refresh(key)
			return entry.value, nil
		}
	}

	return c.refresh(key)
}

func (c *StaleWhileRevalidateCache) refresh(key string) (interface{}, error) {
	value, err := c.renew(key)
	if err != nil {
		c.mu.RLock()
		entry, ok := c.entries[key]
		c.mu.RUnlock()
		if ok && time.Now().Before(entry.staleUntil) {
			return entry.value, nil
		}
		return nil, err
	}

	c.mu.Lock()
	c.entries[key] = staleEntry{
		value:      value,
		expiry:     time.Now().Add(c.ttl),
		staleUntil: time.Now().Add(c.ttl + c.staleTTL),
	}
	c.mu.Unlock()
	return value, nil
}

func (c *StaleWhileRevalidateCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

func (c *StaleWhileRevalidateCache) Invalidate(key string) {
	c.Delete(key)
}

func main() {
	renewCount := 0
	cache := NewStaleWhileRevalidateCache(
		50*time.Millisecond,
		100*time.Millisecond,
		func(key string) (interface{}, error) {
			renewCount++
			return fmt.Sprintf("value-%d", renewCount), nil
		},
	)

	val, _ := cache.Get("key")
	fmt.Println("First get:", val)

	time.Sleep(60 * time.Millisecond)

	val, _ = cache.Get("key")
	fmt.Println("Stale get (served stale, refreshing):", val)

	time.Sleep(60 * time.Millisecond)

	val, _ = cache.Get("key")
	fmt.Println("After revalidate:", val)

	fmt.Println("Total renew calls:", renewCount)
	cache.Delete("key")
	fmt.Println("After delete, size:", cache.Len())
}
```

## Step-by-step execution

1. `NewStaleWhileRevalidateCache(50ms, 100ms, renew)` creates a cache with 50 ms TTL and 100 ms stale window.
2. First `Get("key")`: no entry exists. `refresh` is called synchronously. `renew` returns `"value-1"`. Entry stored with expiry = now+50ms, staleUntil = now+150ms.
3. `time.Sleep(60ms)` — the entry is past its TTL but within the stale window.
4. Second `Get("key")`: entry exists but `time.Now()` is after `expiry`. Since `time.Now()` is before `staleUntil`, the stale value `"value-1"` is returned, and `refresh` is started as a background goroutine.
5. The background goroutine calls `renew`, which returns `"value-2"`. The entry is updated with the new value and expiry.
6. Third `Get("key")`: the fresh entry is returned.
7. `Delete("key")` removes the entry from the map, simulating a write-invalidate.

## Common mistakes

- **Setting a long TTL as the only invalidation mechanism** — A 30-minute TTL means updates are invisible for up to 30 minutes. TTL-only caches are acceptable only for data that changes rarely or where staleness is acceptable.
- **Deleting cache entries before writing the new value to DB** — If the cache entry is deleted and the DB write fails, the next request sees a cache miss and triggers a new DB query that also sees the old data (or no data). The cache is now empty for this key, increasing DB load. Fix: delete after the DB write succeeds.
- **Invalidating an entry that was never cached** — Calling `Delete(key)` for every DB write is wasteful if the key was never cached. Use a bloom filter or set of cached keys to track what has been cached.
- **Not using singleflight for hot keys** — When a popular key expires, thousands of goroutines simultaneously hit the DB. This is cache stampede. Use `singleflight.Group` to coalesce concurrent refreshes to a single DB query.

## Debugging walkthrough

A service shows inconsistent data. Some requests return stale values minutes after a write:

Check cache configuration:

```go
cache := NewTTLCache(5 * time.Minute, 30 * time.Second)
```

The TTL is 5 minutes. Writes use write-invalidate (delete), but the delete happens before the DB write. If the DB write fails, the cache has no entry and the next request fetches the old data from DB. Fix: delete after the DB write:

```go
func UpdateProfile(db *sql.DB, cache *TTLCache, id int, name string) error {
	_, err := db.Exec("UPDATE users SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return err
	}
	cache.Delete(fmt.Sprintf("user:%d", id))
	return nil
}
```

After the fix, consistency is restored: the cache is invalidated only after the DB write succeeds. The next `Get` triggers a cache miss and fetches the new data.

## Production notes

- Use `singleflight.Group` to prevent cache stampede on popular keys. This is a single line of code: `g.Do(key, renewFn)`.
- For write-invalidate, ensure the cache deletion is idempotent. Deleting a key that does not exist is safe.
- For distributed caches (Redis), use a message queue (Redis pub/sub, RabbitMQ) to broadcast invalidation events to all service instances.
- Monitor invalidation rate, miss rate, and stale-hit rate as metrics. A sudden increase in stale hits may indicate a TTL that is too short.
- Consider using a version vector (monotonic counter) for each cache entry. Only update the cache if the version from DB is higher than the cached version.

## Performance implications

Cache invalidation itself is cheap — a map deletion is O(1). The expensive part is the DB write that triggers invalidation. Write-invalidate is preferable to write-update for most systems because:

- Deleting is idempotent and race-condition-free. Two concurrent invalidation calls produce the same result.
- Write-update requires a read of the new data (the DB write returns the row, but constructing the cached value may require joins or computation).
- Write-update has race conditions: a concurrent reader may get the old value while the update is in progress.

For systems with very high read rates (>100,000 req/s), stale-while-revalidate is the best choice because it never blocks a read on a DB query. The stale window absorbs traffic spikes without cache stampede.

## Practice task

Implement a `WriteInvalidateCache` that wraps the `TTLCache` from lesson 07:

1. `Set(key string, value interface{})` writes to both a "database" (internal map) and the cache.
2. `Delete(key string)` deletes from the "database" and invalidates the cache.
3. A background goroutine simulates delayed cache invalidation (100 ms delay).
4. Write a test that verifies the cache returns the new value immediately after `Set` (write-update) and returns a miss after `Delete` (write-invalidate).

## Tests / verification

```bash
go run ./curriculum/modules/13-performance-memory-engineering/lessons/08-cache-invalidation
go test ./curriculum/modules/13-performance-memory-engineering/lessons/08-cache-invalidation
```

## Review questions

1. What is the cache stampede problem, and how does `singleflight.Group` solve it?
2. Why is write-invalidate (cache delete) preferred over write-update for most production systems?
3. What is the stale-while-revalidate pattern, and when should you use it?
4. Why should you delete a cache entry after writing to the database, not before?
5. How does probabilistic early expiration prevent thundering herd on cache refresh?

## NEXT UP

Congratulations on completing Module 13! You now understand profiling, benchmarking, memory optimization, caching, and performance engineering in Go. Next up: Module 14 — Event-Driven Systems.
