# CP.3 sync.Pool: Recycled Memory

## Mission

Master high-performance memory management using `sync.Pool`. Learn how to reduce Garbage Collector (GC) pressure by recycling temporary objects, understand the `Get -> Reset -> Put` lifecycle, and learn why pooling is the secret weapon of Go's standard library.

## Prerequisites

- `CP.2` errgroup-context

## Mental Model

Think of `sync.Pool` as **A Library's Reference Section**.

1. **The Request (`Get`)**: You need a dictionary (a byte buffer). You go to the library and ask if there's one available on the desk.
2. **The New Object (`New`)**: If no dictionaries are on the desk, the librarian goes to the back and buys a brand-new one for you.
3. **The Usage**: You use the dictionary to do your work.
4. **The Cleanup (`Reset`)**: Before you return the dictionary, you erase all your notes from the pages.
5. **The Return (`Put`)**: You put the dictionary back on the desk for the next person to use.
6. **The Garbage Man (`GC`)**: Every once in a while, the librarian clears off the desk and throws everything away. You must be prepared to get a new dictionary next time.

## Visual Model

```mermaid
graph LR
    A[App] -- "Get()" --> P[sync.Pool]
    P -- "Found" --> O[Existing Object]
    P -- "Empty" --> N[New() Function]
    O --> U[Use]
    N --> U
    U -- "Reset" --> R[Clean State]
    R -- "Put()" --> P
```

## Machine View

- **GC Interaction**: `sync.Pool` is tightly integrated with the Go Garbage Collector. During every GC cycle, all objects in all pools are automatically cleared. This prevents the pool from growing indefinitely and leaking memory.
- **CPU Locality**: Internally, `sync.Pool` maintains a separate cache for each CPU core (P). This means that in a multi-threaded system, goroutines on the same core can get/put objects without ever needing a Mutex lock, making it extremely fast.
- **Allocation Savings**: If a web server allocates a 4KB buffer for every request, at 10,000 req/s, that's 40MB/s of garbage. `sync.Pool` can reduce this to nearly zero.

## Run Instructions

```bash
go run ./07-concurrency/02-concurrency-patterns/3-sync-pool
```

## Code Walkthrough

### `New` Function
This is the factory. It is only called if the pool is empty. It should return a pointer to a freshly allocated object with its initial capacity set.

### `Get()`
Retrieves an item from the pool. It returns `any`, so you must use a type assertion: `buf := pool.Get().(*bytes.Buffer)`.

### `Reset()`
**CRITICAL**: You must clear the state of an object before returning it. If you put a buffer containing secret data back into the pool without resetting it, the next user might see that data!

### `Put()`
Returns the object to the pool. Note the "Guard" pattern in the example: if a buffer has grown too large (e.g., to 10MB), we don't put it back, as it would waste memory. We let the GC reclaim it instead.

## Try It

1. Write a benchmark comparing `buildHTTPResponseWithPool` and `buildHTTPResponseWithoutPool` using `go test -bench`.
2. Implement a pool for a custom `struct` that contains a map. Don't forget to clear the map during `Reset`!
3. What happens if you forget to call `Put()`? (Hint: The object just gets garbage collected like normal).

## Verification Surface

Observe the successful reuse of objects and the lack of allocations in high-traffic simulations:

```text
Response via pool:
HTTP/1.1 200 OK
Content-Length: 11
{"ok":true}

Processed: processed request req_001 for user usr_42
Processed: processed request req_002 for user usr_99
```

## In Production
**Don't over-use pools.**
Pooling adds complexity and can lead to subtle bugs (like data leaking between requests). Only use `sync.Pool` for:
- Large objects that are expensive to allocate.
- Small objects that are allocated millions of times per second (e.g., log entries, metadata structs).
- Situations where profiling shows that GC is a bottleneck.

## Thinking Questions
1. Why does `sync.Pool` clear itself during every GC cycle?
2. Why is it better to store pointers in a pool instead of raw values?
3. How can pooling cause "Cache Pollution" if you don't reset objects correctly?

## Next Step

We've mastered the building blocks. Now let's combine them into a high-performance architecture. Continue to [CP.4 Bounded Pipeline](../4-bounded-pipeline-exercise/README.md).
