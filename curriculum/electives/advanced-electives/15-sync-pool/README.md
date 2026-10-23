# sync.Pool

## When you need this

Go applications that allocate many short-lived objects (buffers, structs for serialization) put pressure on the garbage collector, increasing GC pause time and reducing throughput. sync.Pool solves this by reusing objects across goroutines, reducing allocation frequency and GC pressure.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-14

## Mental model

sync.Pool is a temporary object cache per-P (per-processor). Each goroutine accesses its local P's pool with minimal synchronization. Get retrieves an object from the pool (or allocates a new one if empty). Put returns an object to the pool for reuse. The pool is drained during GC — objects are freed, and subsequent Get() calls allocate fresh objects. The pool is NOT a cache with predictable retention; it is an allocation amortization mechanism.

## How it works

sync.Pool maintains per-P private and shared lists. Get() first checks the P's private slot (fastest, no lock), then the P's shared list (fast, per-P lock), then steals from other P's shared lists (slower, cross-P lock), and finally calls New() to allocate. Put() stores in the P's private slot if empty, otherwise appends to the shared list. During each GC cycle, the pool's local caches are drained: all objects are freed unless referenced elsewhere.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: sync.Pool")
}
```

## Tradeoffs

- Fast path Get (per-P private hit): ~3ns — a single atomic load; faster than heap allocation
- Slow path Get (New function called): ~100-300ns (includes allocation) — still faster than a fresh allocation because Pool spreads allocation across GC cycles
- Put is always fast: ~5-10ns — stores to per-P private slot or shared list with minimal synchronization
- Pool reduces GC pressure: objects in the Pool are alive and scanned by GC — a large pool increases GC scan time; balance pool size against GC frequency

## Common mistakes

- Assuming sync.Pool is a persistent cache — objects in the pool are removed during GC; use a regular cache (e.g., LRU) for persistent caching.
- Putting objects without resetting them — the next Get() caller receives an object with residual state; always reset (buf.Reset(), obj.Clear()) before use.
- Using sync.Pool for objects that hold external resources — the pool does not call Close or cleanup; objects holding file handles or network connections leak if not explicitly closed.
- Not calling Put in a defer — if the function panics or returns early, the object is not returned to the pool, increasing GC pressure.
- Using sync.Pool for large objects — large objects in the pool increase GC scan time; pool only objects that are expensive to allocate AND cheap to reset.

## Practice task

The learner must write a benchmark that allocates 1000 bytes.Buffers without sync.Pool and with sync.Pool, report the allocs/op and ns/op for each approach, and explain why sync.Pool reduces GC pressure despite not reducing total allocation count.

## Production notes

sync.Pool is used in every Go production service for: bytes.Buffer pooling (JSON/protobuf serialization), message struct pooling (event publishing), SQL row scanning buffers, template rendering buffers, and any hot-path allocation-heavy operation. Opslane uses sync.Pool for order event marshaling buffers and HTTP response writer buffers.

## Related core lessons

- Module 16.14 — sync.Map
- Module 17.13 — REST vs gRPC tradeoffs
