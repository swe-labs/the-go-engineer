# Mutex and RWMutex

## Learning objective

Protect shared state with `sync.Mutex` and `sync.RWMutex`, understand critical sections, and choose between exclusive and shared locks based on read/write patterns.

## Why this matters

Goroutines share memory — maps, slices, counters, caches, connection pools. Without synchronization, concurrent writes to shared memory produce data races: corrupt data, non-deterministic crashes, and hard-to-reproduce bugs. `sync.Mutex` provides exclusive access to a critical section. `sync.RWMutex` allows multiple concurrent readers while preserving exclusive access for writers. Choosing the right mutex type optimizes throughput for read-heavy workloads. These primitives are the foundation of all shared-state concurrency in Go.

## Mental model

A mutex is a lock on a door to a room containing shared data.

- **`sync.Mutex`**: one person enters the room, locks the door, does work, unlocks the door, and leaves. No one else can enter while the room is occupied. This is **exclusive access**.

- **`sync.RWMutex`**: the door has two locks. The "reader lock" allows unlimited people to enter the room as long as they only look (read). The "writer lock" kicks everyone out, locks both doors, and lets one person modify the data. This is **multiple readers / single writer**.

The protocol is cooperative: all goroutines must follow the locking convention. A single goroutine that reads or writes without locking breaks the safety guarantee.

## Core idea

```go
var mu sync.Mutex
mu.Lock()
// critical section: only one goroutine at a time
x = x + 1
mu.Unlock()

var rwmu sync.RWMutex
rwmu.RLock()
// read section: multiple goroutines can read simultaneously
v := cache[key]
rwmu.RUnlock()

rwmu.Lock()
// write section: exclusive access
cache[key] = value
rwmu.Unlock()
```

Rules:

- Every `Lock()` must have a matching `Unlock()`. Use `defer` for simple cases.
- `RLock()`/`RUnlock()` pairs can overlap with each other but not with `Lock()`/`Unlock()`.
- A `sync.Mutex` or `sync.RWMutex` must not be copied after first use.
- A locked mutex is not reentrant. If a goroutine locks a mutex and then tries to lock it again (same goroutine), it deadlocks.

## Under the hood

`sync.Mutex` in Go is implemented with two modes: normal and starvation.

- **Normal mode**: goroutines are queued in FIFO order, but a newly woken goroutine competes with goroutines that have been running (in a spin loop). If the running goroutine acquires the lock first, the woken goroutine goes back to the tail of the queue. This gives higher throughput under contention.

- **Starvation mode** (since Go 1.9): if a goroutine waits longer than 1ms, the mutex enters starvation mode. In starvation mode, the goroutine at the head of the queue is guaranteed the lock as soon as it is released. The running goroutine does not spin — it goes to the tail. This prevents tail-latency outliers.

`sync.RWMutex` uses an atomic counter for the reader count. `RLock` atomically increments. `Lock` sets a "writer pending" flag and blocks until the reader count is zero and no other writer holds the lock. `RUnlock` decrements; if it was the last reader and a writer is pending, it signals the writer.

Both mutexes use the `runtime_Semacquire`/`runtime_Semrelease` primitives for parking and waking goroutines.

## How Go uses it

- **Maps**: `sync.Map` uses mutexes internally for some operations, but hand-rolled `map[Mutex` is common for custom needs.
- **Connection pools**: `database/sql` uses mutexes to protect the pool of `*sql.Conn`.
- **Caches**: in-memory caches use `sync.RWMutex` for concurrent reads with exclusive writes.
- **Counters**: metrics counters, rate limiters, credit trackers protect state with mutexes.
- **Configuration reload**: hot-reloaded configs use `sync.RWMutex` — many readers, rare writer (reload).

## Go example

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Metrics struct {
	mu     sync.RWMutex
	visits map[string]int64
}

func (m *Metrics) RecordVisit(page string) {
	m.mu.Lock()
	m.visits[page]++
	m.mu.Unlock()
}

func (m *Metrics) Report() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for page, count := range m.visits {
		fmt.Printf("%s: %d\n", page, count)
	}
}

func main() {
	m := Metrics{visits: make(map[string]int64)}

	var wg sync.WaitGroup
	pages := []string{"/", "/about", "/contact"}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			m.RecordVisit(pages[n%len(pages)])
		}(i)
	}
	wg.Wait()

	fmt.Println("Final report:")
	m.Report()
}
```

Output:

```
Final report:
/: 34
/about: 33
/contact: 33
```

The RWMutex allows concurrent reads in `Report()` while `RecordVisit` gets exclusive access for writes.

## Step-by-step execution

For `RecordVisit("/")` under contention:

1. Goroutine A calls `m.mu.Lock()`. The mutex's state flips from unlocked to locked (atomic CAS). Goroutine A enters the critical section.
2. Goroutine B calls `m.mu.Lock()`. The CAS fails. B's goroutine is added to the FIFO queue and parks via `runtime_Semacquire`.
3. A increments `m.visits["/"]`, calls `m.mu.Unlock()`. The unlock sets state to unlocked and wakes the next goroutine (B) via `runtime_Semrelease`.
4. B wakes, acquires the lock, enters the critical section, increments, unlocks.
5. Goroutine C calls `m.mu.RLock()` (on `Report()`). The reader count goes from 0 to 1. Since no writer holds the lock, C enters the read section.
6. Multiple `RLock`s can happen simultaneously. The reader count is atomically incremented for each.
7. When a `Lock()` (write) arrives, it sets a writer-pending flag. New `RLock` calls wait until the writer acquires and releases the lock.

## Common mistakes

- **Forgetting to unlock**: if a function acquires a lock and returns early without unlocking (e.g., due to a `return` before `Unlock`), the mutex stays locked forever.
  - Fix: `defer mu.Unlock()` immediately after `Lock()`.
- **Copying a mutex**: passing a struct with a mutex by value copies the mutex, which is a data race.
  - Fix: pass `*T` or embed the mutex with a pointer receiver.
- **Deadlock with channels**: holding a mutex while sending on a channel that the receiver needs the same mutex for.
  - Fix: lock, read data, unlock, then send.
- **Using Mutex when RWMutex is better**: if reads dominate (e.g., 99% reads, 1% writes), `sync.RWMutex` significantly improves throughput.
  - Fix: use `sync.RWMutex` and protect reads with `RLock`.
- **Locking in defer with conditional logic**: `if cond { mu.Lock(); defer mu.Unlock() }` — the defer defers until the function returns, not the block.
  - Fix: use explicit blocks: `func() { mu.Lock(); defer mu.Unlock(); ... }()`.

## Debugging walkthrough

```go
type Counter struct {
    mu sync.Mutex
    n  int
}

func (c Counter) Inc() {  // value receiver — copies the mutex!
    c.mu.Lock()
    c.n++
    c.mu.Unlock()
}
```

**Symptom**: `go test -race` reports data races. Counter value is wrong under concurrency.

**Investigation**: `Counter` is passed by value. Each call to `Inc()` operates on a copy. The mutex on the copy locks and unlocks, but the original `c` and the copy are different mutexes.

**Root cause**: value receiver method on a struct containing a mutex.

**Fix**: use pointer receiver: `func (c *Counter) Inc()`.

## Production notes

- **Guard the right scope**: lock should surround the critical section, not the entire function. Keep the lock duration minimal. Acquire late, release early.
- **Prefer `sync.Map` for simple cases** but understand its limitations: it is optimized for append-only or read-heavy with infrequent writes. For general use, `map + sync.RWMutex` is often faster and more flexible.
- **Avoid locking in hot paths**: if a frequently called function locks a contended mutex, consider sharding (multiple mutexes, one per shard). For example, `map[string]int64` can be sharded by hash of key.
- **`go vet` detects some mutex copy issues**. Always run `go vet` on concurrent code.
- **Deadlock detection**: the runtime does not detect deadlocks between mutexes. Use tools like `go-deadlock` in testing.

## Performance implications

| Operation | Approximate cost (uncontended) |
|---|---|
| `sync.Mutex.Lock` + `Unlock` | ~20-30 ns |
| `sync.RWMutex.RLock` + `RUnlock` | ~10-15 ns |
| `sync.RWMutex.Lock` + `Unlock` | ~30-50 ns |
| Contended Lock (park + wake) | ~1-10 µs |

Under contention, `sync.Mutex` is cheaper than `sync.RWMutex.Lock` because RWMutex must manage reader state. However, `RWMutex.RLock` is cheaper than `Mutex.Lock` for readers.

Rule of thumb: if more than 80% of operations are reads and writes are rare, use RWMutex. Otherwise, use Mutex — it is simpler and has lower write overhead.

## Practice task

Implement a thread-safe key-value store with TTL support:

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]item
}

type item struct {
    value  string
    expiry time.Time
}
```

Methods: `Get(key string) (string, bool)`, `Set(key, value string, ttl time.Duration)`, `DeleteExpired()`.

`DeleteExpired` iterates all items and removes expired ones (hold write lock only when removing). Use RLock for iteration and upgrade to Lock only when deleting.

Write a test with 10 concurrent goroutines setting and getting 1000 keys each.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/19-mutex-and-rwmutex
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/19-mutex-and-rwmutex
```

## Review questions

1. What is the difference between `sync.Mutex` and `sync.RWMutex`?
2. Can a goroutine lock the same mutex twice? What happens?
3. Why should you use `defer mu.Unlock()` immediately after `mu.Lock()`?
4. When would you choose `sync.RWMutex` over `sync.Mutex`?
5. What does `go test -race` detect that a regular test does not?

## NEXT UP

Atomics with caution — lock-free atomic operations for counters and flags, and when to prefer mutexes over atomics.
