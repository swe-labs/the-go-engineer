# Deadlocks

## Learning objective

Define deadlock using the four Coffman conditions, detect deadlocks in Go programs, and prevent them with consistent mutex ordering and channel discipline.

## Why this matters

A deadlock freezes a program permanently. All affected goroutines are blocked waiting for each other, no progress is possible, and the only recourse is killing the process. In production, a deadlock causes complete service unavailability until a monitor kills and restarts the process. Unlike data races, deadlocks are deterministic given the same execution path: once the circular wait is established, the program is stuck.

## Mental model

A deadlock is a traffic gridlock at an intersection. Four cars each need to cross, but each is blocking the car to its right. No car can move because each is waiting for the car in front. In Go terms, goroutine A holds mutex 1 and waits for mutex 2, while goroutine B holds mutex 2 and waits for mutex 1. Neither can proceed. The only way out is to break the cycle: impose a global order on resource acquisition.

## Core idea

The four Coffman conditions, all of which must hold for a deadlock to occur:

1. **Mutual exclusion**: at least one resource is held in a non-shareable mode (a mutex is locked).
2. **Hold and wait**: a goroutine holds a resource while waiting for another.
3. **No preemption**: a resource cannot be forcibly taken from a goroutine (you cannot unlock a mutex from another goroutine).
4. **Circular wait**: there exists a cycle of goroutines where each holds a resource the next needs.

Prevention strategy: break condition 4 with a consistent lock ordering. If every goroutine acquires mutexes in the same global order, a cycle cannot form.

## Under the hood

The Go runtime detects deadlocks only in limited cases: if all goroutines are blocked and the program would make no further progress, the runtime panics with "fatal error: all goroutines are asleep - deadlock!". However, this detection only works when *all* goroutines are blocked. If at least one goroutine can make progress, the runtime does not detect the deadlock. The runtime also detects single-goroutine deadlocks: if a goroutine locks a non-recursive mutex twice, it panics with "fatal error: sync: unlock of unlocked mutex".

## How Go uses it

The `net/http` package avoids deadlocks by never holding locks across I/O. The `database/sql` package uses a connection pool with a mutex and condition variable, acquiring locks in a strict order. Go's `sync` package explicitly documents that `Mutex` is not recursive and that copying a mutex after use causes undefined behavior.

## Go example

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func deadlockExample() {
	var mu1, mu2 sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		mu1.Lock()
		time.Sleep(10 * time.Millisecond) // force interleaving
		mu2.Lock()
		mu2.Unlock()
		mu1.Unlock()
	}()

	go func() {
		defer wg.Done()
		mu2.Lock()
		time.Sleep(10 * time.Millisecond)
		mu1.Lock() // deadlock: mu2 holds, mu1 held by other goroutine
		mu1.Unlock()
		mu2.Unlock()
	}()
	wg.Wait()
}

func safeOrderedLock() {
	var mu1, mu2 sync.Mutex
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		mu1.Lock()
		mu2.Lock()
		mu2.Unlock()
		mu1.Unlock()
	}()

	go func() {
		defer wg.Done()
		mu1.Lock() // same order: mu1 before mu2
		mu2.Lock()
		mu2.Unlock()
		mu1.Unlock()
	}()
	wg.Wait()
}

func main() {
	fmt.Println("Deadlocks: Coffman conditions")
	fmt.Println("Run tests to verify safe locking order")
}
```

## Step-by-step execution

For `deadlockExample`:

1. Goroutine A locks `mu1`.
2. Goroutine B locks `mu2`.
3. Both sleep for 10 ms (ensuring the other has locked its first mutex).
4. A tries to lock `mu2`. B holds `mu2`, so A blocks.
5. B tries to lock `mu1`. A holds `mu1`, so B blocks.
6. Circular wait: A waits for `mu2`, B waits for `mu1`. Deadlock.

For `safeOrderedLock`:

1. Goroutine A locks `mu1` then `mu2`.
2. Goroutine B tries to lock `mu1`. It is held by A, so B blocks.
3. A finishes its work, unlocks `mu2`, unlocks `mu1`.
4. B's lock on `mu1` succeeds. B locks `mu2`, works, unlocks both.
5. No cycle because both follow the same order.

## Common mistakes

- Mistake: Assuming mutex ordering only matters with three or more mutexes.
  - Why it happens: A cycle needs only two mutexes and two goroutines.
  - Fix: Always acquire locks in a globally consistent order.

- Mistake: Using `sync.Mutex` as if it were recursive (re-entrant).
  - Why it happens: A function that locks a mutex and calls another function that locks the same mutex deadlocks itself. Unlike Java's `synchronized` or Python's `threading.RLock`, Go's `sync.Mutex` is non-recursive by design.
  - Fix: Restructure code so the same mutex is not acquired twice in the same goroutine, or use a `sync.RWMutex` with separate read and write methods that do not call each other.

- Mistake: Sending on an unbuffered channel from a goroutine while the main goroutine waits for the send to complete via `sync.WaitGroup`.
  - Why it happens: The sender blocks until a receiver is ready, but the receiver is waiting for the sender to finish.
  - Fix: Use a buffered channel or reorder the send and receive.

- Mistake: Forgetting that `select` with all channels nil blocks forever.
  - Why it happens: A `select` with no ready case and no `default` blocks until a channel becomes ready. If all channels are nil or closed, the goroutine blocks permanently.
  - Fix: Always include a `default` case or ensure at least one channel is ready.

## Debugging walkthrough

Buggy program:

```go
func main() {
	ch := make(chan int)
	ch <- 1
	fmt.Println(<-ch)
}
```

Symptom: program hangs at `ch <- 1`. The runtime eventually prints "fatal error: all goroutines are asleep - deadlock!".

Investigation: add a goroutine to send:

```go
go func() { ch <- 1 }()
fmt.Println(<-ch) // now works
```

Root cause: an unbuffered channel blocks the sender until a receiver is ready. With only the main goroutine, the send blocks permanently.

Fix: use a goroutine for the sender, or use a buffered channel.

## Production notes

Deadlocks in production are rare in well-reviewed code but happen most often during refactoring when a function that acquires a lock starts calling another function that also acquires a lock. A common trigger is adding instrumentation or logging inside a locked section: if the logger itself acquires a lock, you can create a deadlock between the application lock and the logger lock. Enforce lock ordering by convention and code review. For complex locking, document the lock hierarchy as a comment above the mutex declarations. Use `go test -race` to detect potential deadlocks (the race detector can catch some lock-order inversions). For channels, use `select` with `default` to avoid blocking when no goroutine is ready to communicate.

## Performance implications

A deadlocked program uses zero CPU (all goroutines are blocked) but retains all memory. The process sits idle until a watchdog kills it. The cost is measured in downtime, not CPU cycles. Preventing deadlocks with lock ordering adds no runtime overhead. Adding a `default` case to `select` adds a single branch that is predicted correctly nearly 100% of the time.

## Practice task

Write three goroutines that attempt to lock three mutexes (`muA`, `muB`, `muC`) in an order that creates a cycle. Verify the deadlock by running with a timeout. Then implement a fix using consistent lock ordering. Write a table-driven test that passes for the fixed version and fails (times out) for the deadlocking version.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/24-deadlocks
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/24-deadlocks
```

## Review questions

1. What are the four Coffman conditions for a deadlock?
2. Which condition does consistent lock ordering break?
3. Can the Go runtime always detect a deadlock? Why or why not?
4. What happens when a goroutine tries to lock a `sync.Mutex` it already holds?
5. How does a channel deadlock differ from a mutex deadlock?

## NEXT UP

Errgroup
