# Opslane Milestone 11 Assessment — Caching Layer: Questions

Total: 10 questions (6 multiple choice × 1 pt, 4 short answer × 2 pt = 14 points). Passing score: 10/14 (71%).

---

## Multiple Choice (1 point each)

**1.** What is a "Cache Stampede" (also known as a Thundering Herd)?
- a) When cache memory runs out and triggers an Out-Of-Memory panic.
- b) When a popular cached key expires, and hundreds of concurrent incoming requests simultaneously detect a cache miss and query the database at the exact same moment, overwhelming the database server.
- c) When Redis deletes all keys during a reboot.
- d) When multiple goroutines write to a map without a mutex.

**2.** How does `golang.org/x/sync/singleflight` protect against cache stampedes?
- a) It encrypts database queries so they run faster.
- b) It suppresses duplicate in-flight function executions with the same key, ensuring only one execution runs while all other concurrent callers wait and share the return value.
- c) It doubles the memory cache capacity automatically.
- d) It sets HTTP cache-control headers on client browsers.

**3.** In `CachedOrderService.UpdateOrder`, why is the database updated *before* invalidating the cache (`s.cache.Delete`)?
- a) Database operations are slower than in-memory operations.
- b) If the cache were deleted first and a concurrent read occurred before the database committed, the read would fetch old data from the database and re-populate the cache with stale data.
- c) SQL requires cache keys to exist.
- d) SQLite cannot delete keys.

**4.** Why does `TTLCache.Get` take a read lock (`mu.RLock()`) rather than a full write lock (`mu.Lock()`)?
- a) Read locks allow multiple concurrent goroutines to read cached entries simultaneously without blocking each other, maximizing read throughput.
- b) Write locks do not support expiration checks.
- c) Go maps require `RLock` for all operations.
- d) Read locks consume less power on CPU cores.

**5.** In `GetOrder`, why is there a double-check cache lookup inside the `s.sfGroup.Do` callback function?
- a) To verify that the JSON syntax is valid.
- b) If a preceding singleflight execution completed and populated the cache just before this flight acquired its turn, the double-check hits the warm cache and avoids an unnecessary database query.
- c) To refresh the TTL timestamp.
- d) To log the query execution time.

**6.** What happens if an entry in `TTLCache` has an `ExpiresAt` timestamp that is before `time.Now()`?
- a) `Get` deletes the entry from the map immediately.
- b) `Get` returns `nil, false`, treating the expired entry as an immediate cache miss.
- c) `Get` returns the expired data with a warning.
- d) Go triggers a panic.

---

## Short Answer (2 points each)

**7.** Contrast Cache Invalidation (Cache-Aside / Delete on Write) with Write-Through caching. Why is Cache Invalidation generally preferred in distributed microservices?

**8.** What is Cache Penetration, and how does caching null values or negative results (`nil, nil`) mitigate it?

**9.** Explain the memory leak risk in a passive TTL cache that only checks expiration during `Get()`. How do production LRU (Least Recently Used) eviction policies or active sweepers resolve this?

**10.** In `UnderlyingStore.FetchOrder`, why does the mock store return a clone of the struct (`clone := *o; return &clone, nil`) rather than returning a pointer directly to its internal map value?
