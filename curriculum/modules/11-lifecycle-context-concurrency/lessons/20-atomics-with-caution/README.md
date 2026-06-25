# Atomics with caution

## Learning objective

Use `sync/atomic` operations (`AddInt64`, `Load`, `Store`, `CompareAndSwap`) for lock-free counters and flags, and decide when atomics are appropriate versus mutexes.

## Why this matters

Atomic operations let you read, write, and modify a variable across goroutines without a mutex. They are faster than mutexes for simple operations like incrementing a counter or setting a flag. However, atomics are also easier to misuse: they protect only single-word operations, they do not compose, and they create subtle memory ordering bugs. Every production Go codebase uses atomics somewhere — metrics counters, shutdown flags, sequence generators — but experienced engineers use them sparingly and carefully.

## Mental model

An atomic operation is a hardware-guaranteed single instruction that reads and writes a memory location without interference. Think of it as a magic sticky note: you can increment the number on it, replace it with a new number, or swap it only if it still has the expected value — all in one indivisible step.

Unlike a mutex, there is no "lock" and "unlock." The atomic operation itself is the synchronization. Other goroutines see either the value before the atomic operation or after — never a partial write.

However, atomics only protect the single word they operate on. If you need to atomically update two related variables (e.g., balance and timestamp), a mutex is required.

## Core idea

```go
import "sync/atomic"

// Add
atomic.AddInt64(&counter, 1)

// Load / Store
val := atomic.LoadInt64(&counter)
atomic.StoreInt64(&counter, 42)

// Compare And Swap
swapped := atomic.CompareAndSwapInt64(&counter, old, new)
// If counter == old, set counter = new and return true.
// Otherwise, return false.

// Swap (always sets, returns old value)
old := atomic.SwapInt64(&counter, 100)
```

Newer typed API (Go 1.19+):

```go
var c atomic.Int64
c.Add(1)
c.Load()
c.Store(42)
c.CompareAndSwap(42, 100)
```

When to use atomics:

| Use case | Atomics | Mutex |
|---|---|---|
| Counter increment | Yes (atomic.Add) | Overkill |
| Boolean flag | Yes (atomic.Bool) | Overkill |
| Shutdown signal | Yes (atomic.Load/Store) | Overkill |
| Protect a map | No | Yes |
| Update multiple fields | No | Yes |
| Complex condition | No | Yes |

## Under the hood

Atomic operations compile to CPU instructions like `LOCK XADD` (x86) or `LDADD` (ARM). These instructions:

1. Lock the memory bus (or cache line) so no other core can read or write the same address concurrently.
2. Perform the operation (add, compare-and-swap, etc.).
3. Unlock.

The hardware guarantees that the operation is indivisible. No context switch or interrupt can occur during the operation.

`CompareAndSwap` is the foundation of most lock-free data structures. It loops until successful:

```go
for {
    old := atomic.LoadInt64(&ptr)
    new := compute(old)
    if atomic.CompareAndSwapInt64(&ptr, old, new) {
        break
    }
}
```

Memory ordering: atomic operations in Go provide sequentially consistent ordering. An `atomic.Store` before an `atomic.Load` on another goroutine is always visible. This is the strongest memory ordering guarantee, which prevents the CPU from reordering atomic operations.

## How Go uses it

- **`sync.WaitGroup`**: the counter is an atomic `uint64` packed in a single word.
- **`sync.Once`**: uses `atomic.LoadUint32` and `atomic.StoreUint32` for the done flag.
- **`sync.Mutex`**: uses atomic operations internally for the lock state.
- **`runtime` internals**: the Go scheduler, memory allocator, and garbage collector all use atomics extensively.
- **Metrics libraries**: `expvar`, `prometheus/client_golang` use atomics for counter and gauge updates to avoid mutex contention on hot paths.

## Go example

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type ServerMetrics struct {
	requests atomic.Int64
	errors   atomic.Int64
	active   atomic.Int64
}

func (m *ServerMetrics) RequestHandled() {
	m.requests.Add(1)
}

func (m *ServerMetrics) ErrorOccurred() {
	m.errors.Add(1)
}

func (m *ServerMetrics) ConnectionOpened() {
	m.active.Add(1)
}

func (m *ServerMetrics) ConnectionClosed() {
	m.active.Add(-1)
}

func (m *ServerMetrics) Snapshot() (requests, errors, active int64) {
	return m.requests.Load(), m.errors.Load(), m.active.Load()
}

func main() {
	var m ServerMetrics
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.ConnectionOpened()
			m.RequestHandled()
			if i%10 == 0 {
				m.ErrorOccurred()
			}
			m.ConnectionClosed()
		}()
	}
	wg.Wait()

	req, err, act := m.Snapshot()
	fmt.Printf("requests=%d errors=%d active=%d\n", req, err, act)
}
```

Output:

```
requests=100 errors=10 active=0
```

## Step-by-step execution

For `counter.Add(1)` in 3 concurrent goroutines:

1. Goroutine A executes `LOCK XADD [counter], 1`. CPU locks the cache line containing `counter`. Reads `counter=0`, adds 1, writes `counter=1`. Unlocks.
2. Goroutine B executes `LOCK XADD [counter], 1`. Reads `counter=1`, writes `counter=2`.
3. Goroutine C does the same. Result: `counter=3`.

No goroutine ever sees a stale or intermediate value. If this were a regular `counter++` (which is load, increment, store — 3 instructions), the interleaving could produce 2 instead of 3.

## Common mistakes

- **Using atomics for non-atomic needs**: `counter++` where `counter` is `atomic.Int64` — even though the variable is atomic, the `++` is not the atomic add. You must use `counter.Add(1)`. The `atomic.Int64` type only makes the `Load`/`Store`/`Add` methods atomic, not regular arithmetic.
- **Forgetting memory ordering**: if goroutine A writes multiple values (one atomic, one not), another goroutine B might see the atomic write but not the non-atomic write. Use `atomic.Store` for all synchronized variables, not just one.
- **Using atomics when a mutex is needed**: atomic operations cannot atomically update two related fields. If `x` and `y` must be consistent, use a mutex.
- **CAS in a tight loop (contended)**: if `CompareAndSwap` fails repeatedly (high contention), the loop burns CPU. This is called "livelock." Mutexes are better under high contention.
- **Using atomics with non-word-sized types**: atomics only work with integer types, pointers, and `bool`. You cannot atomically update a `struct`, `string`, or `slice`.

## Debugging walkthrough

```go
var counter int64

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++  // not atomic!
        }()
    }
    wg.Wait()
    fmt.Println(counter)
}
```

**Symptom**: prints values like 997, 998, 999 — never 1000. `go test -race` reports a data race.

**Investigation**: `counter++` is not atomic. It compiles to: load `counter` into register, increment register, store register to `counter`. Between load and store of one goroutine, another goroutine loads the same value.

**Root cause**: concurrent non-atomic increment.

**Fix**: use `atomic.AddInt64(&counter, 1)`.

## Production notes

- **Zero values are safe**: `atomic.Int64` with zero value is ready to use. No constructor needed.
- **`atomic.Value` is for arbitrary types**: `atomic.Value` can store and load any type atomically, but it must be the same type for all stores, or it panics. It is used for configuration reload, cache swaps, and lock-free reads of immutable data.
- **Profiling**: `pprof` shows atomic operations as `runtime/internal/atomic` functions. If atomics dominate CPU profiles, consider batching.
- **Avoid premature optimization**: write the simple mutex version first. Profile. If mutex contention appears on the profile, replace with atomics.
- **`go vet` is limited**: `go vet` does not detect the `counter++` race. Always run `go test -race`.

## Performance implications

| Operation | Approximate cost |
|---|---|
| `atomic.AddInt64` (uncontested) | ~5-10 ns |
| `atomic.Load` | ~2-5 ns |
| `atomic.Store` | ~5-10 ns |
| `atomic.CompareAndSwap` (success) | ~10-15 ns |
| `sync.Mutex.Lock + Unlock` (uncontested) | ~20-30 ns |

Under low contention, atomics are 2-5x faster than mutexes. Under high contention, the difference narrows because CPU cache-line bouncing dominates. On NUMA systems, contended atomics can be very expensive (interconnect traffic).

Atomics do not block goroutines — they spin in hardware. This means they do not cause goroutine scheduling overhead, but they also do not yield. A contended CAS loop can starve other goroutines if it never succeeds.

## Practice task

Implement a concurrent sequence generator using `atomic.Int64`:

```go
type Sequencer struct {
    counter atomic.Int64
}

func (s *Sequencer) Next() int64
func (s *Sequencer) Peek() int64   // returns next value without consuming
func (s *Sequencer) Reset()        // sets back to 0 (use with caution)
```

Then implement `ResetIf` that resets only if the current value equals the given parameter (use `CompareAndSwap`). Write a test that launches 100 goroutines each calling `Next()` 100 times and verifies no duplicates and the final value is 10000.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/20-atomics-with-caution
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/20-atomics-with-caution
```

## Review questions

1. Why does `counter++` with `atomic.Int64` not work as expected?
2. What is the difference between `atomic.AddInt64` and `atomic.CompareAndSwapInt64`?
3. When would a mutex be a better choice than an atomic operation?
4. Why is `atomic.Value` useful for configuration reload patterns?
5. What does `go test -race` detect that atomics do not fix on their own?

## NEXT UP

Race conditions — how data races occur, how to detect them with the race detector, and patterns to prevent them.
