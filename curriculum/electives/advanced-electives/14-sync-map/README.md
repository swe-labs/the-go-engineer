# sync.Map

## When you need this

Go's built-in map type is not safe for concurrent use — concurrent reads and writes cause a fatal runtime panic. Using a sync.Mutex to protect map access adds contention that slows down read-heavy workloads. sync.Map solves this with an optimized concurrent map that uses an atomic read path for the common case and a write-locked path for mutations.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-13

## Mental model

sync.Map has two internal maps: a read map (accessed atomically, no lock) and a dirty map (accessed with a mutex). Reads first check the read map — if the key exists, the read returns instantly with no lock. If the key is not in the read map, the read falls back to the dirty map with a lock. Writes always go to the dirty map with a lock. Periodically, the dirty map is promoted to become the new read map (amended=true).

## How it works

The read map is a pointer to a readOnly struct that contains a map and an amended flag. Load() atomically loads the read map pointer and checks for the key. If found (and not deleted), return the value without locking. If not found and amended is true, acquire the mutex, promote the dirty map to read, and retry. Store() acquires the mutex, writes to the dirty map, and sets amended=true. Delete() marks the key as deleted in the read map (without removing it) and removes it from the dirty map under the mutex.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: sync.Map")
}
```

## Tradeoffs

- sync.Map read with read-map hit: a single atomic pointer load — ~1-2ns, comparable to regular map read with no lock contention
- sync.Map read with read-map miss (dirty fallback): mutex Lock + map lookup + Unlock — ~30-50ns, comparable to RWMutex.RLock
- sync.Map write: mutex Lock + dirty map write + potential promotion (O(n) copy) — promotion cost is proportional to total entry count, not just written entries
- sync.Map delete does not reclaim memory: the key is marked deleted in the read map but the entry remains — keys accumulate over time

## Common mistakes

- Using sync.Map when a regular map with sync.Mutex is simpler — sync.Map has a complex API (Load, Store, LoadOrStore, Delete, Range) that is harder to use than a regular map with a mutex.
- Using sync.Map for write-heavy workloads — sync.Map's performance degrades when writes are frequent because every write requires a mutex lock and periodic map promotion.
- Using sync.Map for small maps — a regular map with a mutex performs better for maps with fewer than 10 entries because the atomic read path overhead is not justified.
- Not using Range for iteration — sync.Map does not support for range; use Range(f func(key, value any) bool) instead, and return false to stop iteration.
- Using sync.Map as a typed map — Load returns any, requiring type assertions. Consider a typed wrapper around sync.Map or a regular map with a mutex for type safety.

## Practice task

The learner must implement a typed cache using sync.Map that stores and retrieves values with type-safe Load and Store wrapper methods, write a concurrent test that proves 10 goroutines reading and 2 goroutines writing do not cause race conditions, and compare the performance against a map with sync.RWMutex using a benchmark.

## Production notes

sync.Map is used in Go for: connection pools (HTTP client, database), configuration caches (read-mostly, updated infrequently), feature flag registries (set at startup, read by every request), metrics registries (registered once, read periodically by the metrics exporter), and DNS or IP caches.

## Related core lessons

- Module 16.13 — REST vs gRPC tradeoffs
- Module 17.12 — gRPC interceptors
