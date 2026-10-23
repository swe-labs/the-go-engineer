# Opslane caching layer

## Mission

Implement a high-performance in-memory caching subsystem for Opslane with Time-To-Live (TTL) expiration, Cache-Aside / Read-Through lookups, mutation invalidation, and Singleflight cache stampede defense.

## Prerequisites

- `opslane-10`

## Mental Model

Direct database queries for frequently read, infrequently updated entities (such as tenant metadata or order summaries) add unnecessary latency and exhaust database connection pools. Caching stores serialized payloads in memory for fast retrieval. However, naive caching introduces two major operational hazards:
1. **Cache Stampede (Thundering Herd)**: When a popular key expires under heavy load, hundreds of concurrent requests simultaneously detect a cache miss and execute identical database queries, overwhelming the database server.
2. **Stale Reads**: When an entity is updated in the database, retaining the cached version serves outdated data to users.

Opslane mitigates stampedes using Go's `singleflight.Group`: only one goroutine executes the database query while other concurrent callers wait and share the result. Stale reads are prevented by strict cache invalidation on write (`cache.Delete`).

## Visual Model

```text
Concurrent Client Requests (100 simultaneous calls for order:ord-1)
  │   │   │   │
  └───┴───┼───┘
          ▼
┌─────────────────────────────────────────────────────────────┐
│ TTLCache.Get("order:ord-1")                                 │
└──────────────┬───────────────────────────────┬──────────────┘
               │ Hit                           │ Miss (Expired / Absent)
               ▼                               ▼
       Return Cached Data       ┌─────────────────────────────┐
                                │ singleflight.Group.Do       │
                                │   (1 DB query executed)     │
                                └──────────────┬──────────────┘
                                               │
                                               ▼
                                ┌─────────────────────────────┐
                                │ Store in TTLCache           │
                                │ Broadcast Result to All 100 │
                                └─────────────────────────────┘
```

## Machine View

Under the hood:
1. `TTLCache` protects its internal map with `sync.RWMutex`. Reading keys takes `c.mu.RLock()`, checks `time.Now().After(entry.ExpiresAt)`, and returns `hit == false` if expired.
2. When a cache miss occurs, `s.sfGroup.Do(cacheKey, fn)` registers the in-flight flight. If 50 goroutines call `GetOrder` simultaneously for `"ord-1"`, exactly *one* goroutine enters `fn()` to query `s.store.FetchOrder(ctx, id)`.
3. Inside `fn()`, a double-check pattern runs: if another concurrent flight just populated the cache, it returns early.
4. When `fn()` returns, `singleflight` broadcasts the return value and error to all waiting goroutines simultaneously, and `s.dbReads` increments by exactly 1.
5. In `UpdateOrder`, updating the database source of truth is executed first; upon success, `s.cache.Delete(cacheKey)` removes the key so subsequent reads fetch fresh data.

## Run Instructions

```bash
# Run caching layer test suite
go test -v ./curriculum/modules/18-flagship-opslane/assessments/11-opslane-caching-layer

# Run standalone demo
go run ./curriculum/modules/18-flagship-opslane/assessments/11-opslane-caching-layer
```

## Code Walkthrough

1. **`TTLCache` & Expiration (lines 18–66)**: Thread-safe map holding `CacheEntry` with `ExpiresAt`; `Get` returns false if expired; `Set` stores with TTL; `Delete` purges keys.
2. **`OrderRecord` & `UnderlyingStore` (lines 68–80)**: Entity representation and interface decoupling database access.
3. **`CachedOrderService` (lines 82–98)**: Encapsulates cache, underlying store, TTL duration, and `singleflight.Group`.
4. **`GetOrder()` & Stampede Protection (lines 101–146)**:
   - Evaluates fast-path cache hit.
   - On miss: invokes `s.sfGroup.Do(cacheKey, ...)`.
   - Double-checks cache; fetches from store; serializes to JSON; sets cache; returns order.
5. **`UpdateOrder()` & Invalidation (lines 149–163)**: Updates database first, then purges cache key via `s.cache.Delete(cacheKey)`.

## Try It

1. Run the test suite: verify that 100 concurrent requests for an uncached key trigger only 1 database read.
2. Test cache expiration: wait for TTL to expire and verify subsequent request refetches from database.
3. Test invalidation: update an order and verify that subsequent `GetOrder` returns updated data immediately.

Common mistakes:
- Updating the cache with stale data before committing the database transaction, causing race condition overwrites.
- Not using `singleflight`, allowing cache stampedes to crash production databases on key expiration.
- Forgetting to clone structs or return immutable bytes, leading to data races if callers mutate cached pointers.

## In Production

In distributed microservices, in-memory caches are paired with Redis clusters in a multi-tier caching architecture: L1 (in-process `sync.Map` / LRU with sub-microsecond latency) and L2 (centralized Redis with sub-millisecond latency). Singleflight is equally vital when querying L2 Redis or downstream databases.

## Thinking Questions

1. What is a "Cache Stampede" (Thundering Herd) and how does `singleflight.Group` prevent it?
2. Why is Cache Invalidation (`Delete`) preferred over Cache Update (`Set`) when modifying entities in the database?
3. What are the memory risks of an in-memory TTL cache that does not run periodic background cleanup of expired keys?

## Next Step

Proceed to [Milestone 12 — Observability](../12-opslane-observability/README.md) to integrate Prometheus dimensional metrics, OpenTelemetry distributed tracing, and health check dashboards.
