# Opslane Milestone 11 Assessment — Answer Key

## Multiple Choice (1 point each)

**1. b) When a popular cached key expires, and hundreds of concurrent incoming requests simultaneously detect a cache miss and query the database at the exact same moment, overwhelming the database server.**  
*Explanation:* A cache stampede occurs when an expired or missing key receives high concurrent read traffic. Every caller misses the cache and attempts to populate it simultaneously, knocking down backend databases with redundant queries.

**2. b) It suppresses duplicate in-flight function executions with the same key, ensuring only one execution runs while all other concurrent callers wait and share the return value.**  
*Explanation:* `singleflight.Group.Do(key, fn)` checks if a call for `key` is already running. If so, subsequent callers block until the active call finishes and receive a shared duplicate of the return value and error.

**3. b) If the cache were deleted first and a concurrent read occurred before the database committed, the read would fetch old data from the database and re-populate the cache with stale data.**  
*Explanation:* If Cache is deleted first: (1) Reader misses cache; (2) Reader reads old value from DB before Writer commits; (3) Writer commits new value to DB; (4) Reader writes old value back into Cache. The cache now holds stale data permanently until TTL expiration. Updating DB first mitigates this race window.

**4. a) Read locks allow multiple concurrent goroutines to read cached entries simultaneously without blocking each other, maximizing read throughput.**  
*Explanation:* In read-heavy workloads (95%+ reads), a standard `sync.Mutex` serializes all reads into single-file execution. `sync.RWMutex.RLock()` permits hundreds of concurrent readers across all CPU cores.

**5. b) If a preceding singleflight execution completed and populated the cache just before this flight acquired its turn, the double-check hits the warm cache and avoids an unnecessary database query.**  
*Explanation:* In high-concurrency environments, double-checking inside the barrier guarantees that if the key was populated by a previous worker right as this execution was scheduled, the database fetch is bypassed.

**6. b) `Get` returns `nil, false`, treating the expired entry as an immediate cache miss.**  
*Explanation:* Passive expiration checks `time.Now().After(entry.ExpiresAt)` on access. If expired, `Get` signals a cache miss so the application can trigger a fresh lookup.

---

## Short Answer (2 points each)

**7. Cache Invalidation (Cache-Aside) vs. Write-Through Caching:**
- **Cache-Aside (Invalidate on Write):** The application updates the database and deletes the cache key. Next time the resource is read, it is re-fetched and cached lazily.  
  *Pros:* Simple concurrency semantics; avoids race conditions where two concurrent writes update the cache in reverse order; keys that are never read again don't waste cache memory.
- **Write-Through:** The application updates both database and cache simultaneously in the same call.  
  *Cons:* Added write latency; concurrent writes can overwrite the cache with out-of-order stale data; writes clutter cache with data that may not be read.
*Grading: 1 pt for explaining delete/lazy fetch vs write-through; 1 pt for analyzing race condition avoidance in cache-aside.*

**8. Cache Penetration and mitigation:**
- **What is Cache Penetration:** When clients repeatedly query for non-existent keys (e.g. `order:non_existent_id`). Because the key is never in the cache, every request bypasses the cache and hits the database, potentially overwhelming it.
- **Mitigation:** Cache negative results (null values) with a short TTL (e.g. 30 seconds), or use a Bloom filter in front of the cache to immediately reject keys that definitely do not exist in the database.
*Grading: 1 pt for defining non-existent key barrage; 1 pt for explaining caching nulls or Bloom filter.*

**9. Memory leak risk in passive TTL caches:**
- **Risk:** A passive TTL cache only evaluates expiration when `Get(key)` is invoked. If an application caches 1,000,000 keys with a 5-minute TTL, but never queries those keys again, they will remain in memory indefinitely, eventually causing an Out-Of-Memory (OOM) crash.
- **Resolution:**
  1. **Active Sweepers:** A background goroutine periodically samples random keys and purges expired entries (the Redis approach).
  2. **LRU Eviction (Least Recently Used):** A fixed memory/entry ceiling where adding new items automatically evicts the oldest unaccessed items regardless of TTL.
*Grading: 1 pt for explaining unaccessed expired keys accumulating in RAM; 1 pt for explaining active sweepers or LRU eviction.*

**10. Pointer aliasing and struct cloning in mock store:**
If `FetchOrder` returns a raw pointer `*OrderRecord` directly referencing the internal map value, the caller now holds a direct memory pointer to the store's state. If the caller mutates `order.Amount = 999`, the internal store's data is corrupted without acquiring the store's mutex lock, and concurrent readers will experience data races. Returning a clone (`clone := *o; return &clone, nil`) ensures pointer isolation.
*Grading: 1 pt for explaining memory aliasing / mutation without lock; 1 pt for explaining data race prevention through cloning.*
