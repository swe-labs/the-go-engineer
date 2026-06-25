# Race conditions

## Learning objective

Define a data race, identify race conditions in concurrent Go code, and fix them using a mutex or a channel.

## Why this matters

Data races are the most common class of concurrency bug in Go. A single unsynchronized write from one goroutine while another reads or writes the same variable produces undefined behavior: the program may crash, produce wrong output, or silently corrupt data. Races are notoriously hard to reproduce because they depend on the Go scheduler's interleaving, which changes across runs, machines, and CPU counts. Every Go engineer must recognize the pattern on sight and know the two canonical fixes: mutex and channel.

## Mental model

A data race is like two people editing the same whiteboard cell at the same time. One writes "42", the other writes "99", and a third person reading at that instant might see "42", "99", or even "4" followed by "2" because the write was not atomic. The Go memory model guarantees nothing about reads that overlap with unsynchronized writes. The only way to make the outcome deterministic is to ensure that at most one goroutine accesses the variable at a time (mutex) or that access happens through a single owner that communicates via a channel.

## Core idea

A data race occurs when two or more goroutines access the same variable concurrently, and at least one access is a write. All such unsynchronized accesses are data races. The Go memory model specifies that a read `r` of a variable `v` is allowed to observe a write `w` to `v` only if `r` does not happen before `w` and `w` does not happen before `r`. Without synchronization, there is no happens-before relation between goroutines.

There are only two correct ways to prevent data races in Go:

1. **Mutex**: `sync.Mutex` or `sync.RWMutex` ensures exclusive access to the shared variable.
2. **Channel**: communication via a channel creates a happens-before edge: a send happens before the corresponding receive completes.

## Under the hood

The Go compiler and CPU can reorder memory accesses within a single goroutine as long as the goroutine's own semantics are preserved. On a multicore machine, each core has its own cache. Without synchronization, a write on core 1 may never be visible to a read on core 2, or may be visible in a different order. The `sync.Mutex` uses a memory barrier (full fence) on unlock and acquire semantics on lock. The channel send/receive pair also establishes a memory barrier via the runtime's internal lock around the channel buffer.

The race is not just about torn writes. Even reading a 64-bit value on a 32-bit platform can tear, and a loop like `for i < len(slice)` can see a stale length after another goroutine appends.

## How Go uses it

Go's standard library and runtime are carefully synchronized. The `net/http` server handles each request in a new goroutine, and the `sync.Pool`, `database/sql`, and `log` packages all rely on mutexes internally. Go's philosophy is to share memory by communicating, not to communicate by sharing memory, but in practice mutexes are used extensively when state is too complex or performance-sensitive for channels.

## Go example

```go
package main

import (
	"fmt"
	"sync"
)

type RacyCounter struct {
	value int
}

func (c *RacyCounter) Add(n int) {
	c.value += n // DATA RACE: concurrent writes
}

func (c *RacyCounter) Value() int {
	return c.value // DATA RACE: concurrent read with write
}

type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Add(n int) {
	c.mu.Lock()
	c.value += n
	c.mu.Unlock()
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	var wg sync.WaitGroup
	racy := &RacyCounter{}
	safe := &SafeCounter{}
	n := 1000

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			racy.Add(1)
			safe.Add(1)
		}()
	}
	wg.Wait()

	fmt.Printf("Expected: %d\n", n)
	fmt.Printf("Racy counter: %d (may differ)\n", racy.Value())
	fmt.Printf("Safe counter: %d (correct)\n", safe.Value())
}
```

## Step-by-step execution

When 1000 goroutines execute `racy.Add(1)`:

1. Two goroutines A and B both read `c.value` at nearly the same time (both see 42).
2. A computes `42 + 1 = 43` and writes 43 back.
3. B computes `42 + 1 = 43` and writes 43 back.
4. Two increments produced only one net increase. The final value is unpredictable.

With `safe.Add(1)`:

1. A calls `Lock()` and acquires the mutex.
2. B calls `Lock()` and blocks (spins or sleeps) waiting for A.
3. A reads `c.value` (42), increments, writes 43, then calls `Unlock()`.
4. B's `Lock()` now succeeds. B reads 43, increments to 44, writes, unlocks.
5. Every increment is visible. The final value is exactly 1000.

## Common mistakes

- Mistake: Using a mutex on the `Add` method but not on `Value`, assuming reads are safe.
  - Why it happens: The developer thinks a read without a write is harmless, but the read sees stale data.
  - Fix: Lock in all methods that access the field, including reads.

- Mistake: Using a mutex value instead of a pointer (copying a `sync.Mutex`).
  - Why it happens: `sync.Mutex` must not be copied after first use. Passing a struct by value copies its mutex, defeating synchronization.
  - Fix: Always use `*sync.Mutex` or embed `sync.Mutex` in a struct passed by pointer.

- Mistake: Writing `go func() { counter.Add(1) }()` without a `WaitGroup`, then checking the counter before the goroutine runs.
  - Why it happens: The main goroutine exits before the child goroutine executes.
  - Fix: Use `sync.WaitGroup` to wait for completion, or use a channel to signal done.

- Mistake: Assuming `atomic.LoadInt64` and `atomic.StoreInt64` are always the right fix.
  - Why it happens: Atomic operations prevent torn reads but do not prevent race windows between a read and a dependent write.
  - Fix: Use atomics only for counters and flags. For compound operations (read-modify-write), use a mutex or `atomic.Add*`.

## Debugging walkthrough

Consider this buggy program:

```go
type Cache struct {
	items map[string]string
}

func (c *Cache) Get(key string) string {
	return c.items[key]
}

func (c *Cache) Set(key, val string) {
	c.items[key] = val
}
```

Symptom: intermittent panics with "assignment to entry in nil map" or "concurrent map writes". The map is written from one goroutine while another reads. Go's map is not safe for concurrent use.

Investigation: run with `go run -race .`:

```
WARNING: DATA RACE
Write at 0x... by goroutine 7:
  main.(*Cache).Set(...)
Previous read at 0x... by goroutine 8:
  main.(*Cache).Get(...)
```

Root cause: `Get` and `Set` access the map with no synchronization.

Fix: add a mutex:

```go
type Cache struct {
	mu    sync.Mutex
	items map[string]string
}

func (c *Cache) Get(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.items[key]
}

func (c *Cache) Set(key, val string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = val
}
```

## Production notes

In real services, race conditions often appear in lazy caches, connection pools, and shared metrics accumulators. Common patterns that hide races: map writes in HTTP handlers, slice appends in background workers, and shared struct field updates in tests that only use a single goroutine. Always lock the full critical section. Never assume a single-machine test proves correctness. Run CI with `-race` on every test suite.

## Performance implications

A mutex `Lock`/`Unlock` pair on a contended path costs roughly 25-100 ns uncontended and can reach microseconds under contention as goroutines spin and park. Channels add similar overhead for the happens-before guarantee. However, the cost of a data race in production (corrupt data, crashes, silent billing errors) is orders of magnitude higher. Correctness first, then profile and optimize.

## Practice task

Rewrite the `RacyCounter` to be safe using a channel instead of a mutex. Create a `ChannelCounter` with a private channel. The `Add` method sends on the channel, and a dedicated goroutine receives on the channel and updates the value. The `Value` method must use a separate request/reply channel pattern or a closure. Verify with the same concurrent access pattern.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/21-race-conditions
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/21-race-conditions
go test -race ./curriculum/modules/11-lifecycle-context-concurrency/lessons/21-race-conditions
```

## Review questions

1. What three conditions must hold for a data race to occur?
2. Can a data race cause a program crash, or only produce wrong values?
3. Why is reading a shared variable without a mutex still a race even if no goroutine writes at that exact instant?
4. What happens if you copy a struct that contains a `sync.Mutex`?
5. Why does using a channel prevent a data race but a buffered channel with capacity 0 can still cause a deadlock?

## NEXT UP

Race detector
