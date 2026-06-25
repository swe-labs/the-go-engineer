# Sync.Once

## Learning objective

Use `sync.Once` for thread-safe lazy initialization, implement the singleton pattern, and understand when to use `sync.Pool` as an alternative.

## Why this matters

Lazy initialization is a common performance optimization: you defer expensive work until it is actually needed. But in concurrent Go code, lazy initialization must be thread-safe — two goroutines must not both run the initializer. `sync.Once` provides the canonical one-shot gate that runs a function exactly once, no matter how many goroutines call `Do` simultaneously. It is used for initializing configs, connection pools, caches, and registry entries.

## Mental model

`sync.Once` is a one-way turnstile. The first goroutine to call `Do(f)` passes through and runs `f`. After `f` completes, the turnstile locks permanently. Every subsequent goroutine that calls `Do(f)` sees the locked turnstile and passes through without running `f`. The turnstile can never be reset. Once a function has been executed (or is currently executing), no other goroutine will ever execute it again.

## Core idea

`sync.Once` has a single method:

```go
func (o *Once) Do(f func())
```

`Do` calls `f` exactly once. If multiple goroutines call `Do` concurrently, one runs `f` and the others block until `f` returns, then proceed without calling `f`. If `f` panics, `Do` still considers the function "done" — the panic propagates, but the `Once` is marked as executed and will never run `f` again. This is a deliberate design choice: if you need to retry after a panic, you cannot use `sync.Once`.

## Under the hood

The implementation uses a 32-bit (`uint32`) done flag and a mutex. The fast path is a single `atomic.LoadUint32(&o.done)`. If done is 1, `Do` returns immediately (no function call, no mutex acquisition). If done is 0, `Do` acquires the mutex, double-checks done (still 0), executes `f`, sets done to 1 with `atomic.StoreUint32`, and releases the mutex. This is double-checked locking, made safe by Go's memory model: the atomic store in `Do` happens-before any subsequent atomic load in another goroutine's call to `Do`.

## How Go uses it

The standard library uses `sync.Once` extensively:

- `sync.Pool` initialization: each pool's `New` function is called via `sync.Once`.
- `os/exec`: the `findExecutable` cache uses `sync.Once`.
- `net/http`: the `DefaultTransport` is initialized once.
- `crypto/tls`: the certificate cache uses `sync.Once`.
- `database/sql`: driver registration uses `sync.Once` internally.

## Go example

```go
package main

import (
	"fmt"
	"sync"
)

type Config struct {
	Addr string
	Port int
}

var (
	config     *Config
	configOnce sync.Once
)

func LoadConfig() *Config {
	configOnce.Do(func() {
		fmt.Println("Loading configuration...")
		config = &Config{
			Addr: "0.0.0.0",
			Port: 8080,
		}
	})
	return config
}

type Pool struct {
	items chan interface{}
	once  sync.Once
}

func NewPool(size int) *Pool {
	p := &Pool{}
	p.once.Do(func() {
		p.items = make(chan interface{}, size)
		for i := 0; i < size; i++ {
			p.items <- struct{}{}
		}
	})
	return p
}

func (p *Pool) Get() interface{} {
	return <-p.items
}

func (p *Pool) Put(x interface{}) {
	p.items <- x
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cfg := LoadConfig()
			fmt.Printf("Config: %+v\n", cfg)
		}()
	}
	wg.Wait()

	pool := NewPool(3)
	fmt.Printf("Pool acquired %v\n", pool.Get())
	pool.Put(struct{}{})
	fmt.Println("Done")
}
```

## Step-by-step execution

For `LoadConfig` with 10 goroutines:

1. Goroutine A calls `LoadConfig`. `configOnce.Do(f)` executes the fast path: `atomic.LoadUint32` reads 0.
2. A acquires the mutex inside `Do`. Double-check: done is still 0. A runs `f`, printing "Loading configuration..." and setting `config`.
3. A sets done to 1 atomically and releases the mutex. A returns `config`.
4. Goroutine B calls `LoadConfig`. Fast path: `atomic.LoadUint32` reads 1. `Do` returns immediately without calling `f`.
5. Goroutines C through J all see done=1, return immediately.
6. "Loading configuration..." is printed exactly once.

## Common mistakes

- Mistake: Using `sync.Once` inside a loop or for functions that must run again.
  - Why it happens: `Do` runs `f` only on the first call. Subsequent calls are no-ops silently.
  - Fix: Use a `sync.Mutex` with a guard variable if the function may need to run again.

- Mistake: Wrapping `sync.Once.Do(f)` in an external mutex.
  - Why it happens: The developer thinks `Do` needs external synchronization.
  - Fix: `sync.Once` is itself goroutine-safe. An external mutex is redundant and defeats performance.

- Mistake: Calling `Do` with a function that panics, expecting a retry.
  - Why it happens: After the panic, `Once` is marked done. The function will never run again.
  - Fix: If the function can fail, handle errors inside `f` (e.g., by retrying internally) or use a different pattern.

- Mistake: Copying a `sync.Once` value.
  - Why it happens: `sync.Once` must not be copied after first use. Its internal mutex and atomic are stateful.
  - Fix: Always use `*sync.Once` or embed it in a struct passed by pointer.

## Debugging walkthrough

Buggy program:

```go
var once sync.Once
var conn *sql.DB
var connErr error

func getConn() (*sql.DB, error) {
	once.Do(func() {
		conn, connErr = sql.Open("driver", "dsn")
	})
	return conn, connErr
}
```

Symptom: if `sql.Open` fails, subsequent calls to `getConn` return the same error without retrying. The once gate does not reset on error.

Investigation: `sql.Open` may return an error (invalid DSN). Once `Do` returns, `once` is marked executed. All future calls return the nil `conn` and the error without re-initializing.

Fix: Do not use `sync.Once` for initialization that can fail. Instead, use `sync.Mutex` with a guard:

```go
var mu sync.Mutex
var conn *sql.DB

func getConn() (*sql.DB, error) {
	mu.Lock()
	defer mu.Unlock()
	if conn != nil {
		return conn, nil
	}
	var err error
	conn, err = sql.Open("driver", "dsn")
	return conn, err
}
```

## Production notes

Use `sync.Once` for one-time initialization that cannot fail: configuration loading from a file that exists, TLS certificate caching, Prometheus metric registration, and global variable initialization. For initialization that can fail (database connections, network calls), use a `sync.Mutex` with an `init` guard, or use `sync.Once` with the initialization done at startup before accepting requests.

`sync.Pool` is a related pattern: it maintains a pool of reusable objects. The `New` field is a function that `sync.Pool` calls via `sync.Once` per pool instance when `Get` finds no available objects. `sync.Pool`'s contents are automatically cleared at each GC cycle, making it suitable for temporary objects (like buffers) but not for persistent singletons.

## Performance implications

`sync.Once.Do` on the fast path (after initialization) is a single atomic load — approximately 1-3 ns. On the slow path (first call), it acquires a mutex and executes `f`. The cost is entirely dominated by `f`. After initialization, `sync.Once` is essentially free. This makes it the cheapest thread-safe initialization mechanism in Go.

## Practice task

Implement a `LazyValue` type that wraps `sync.Once`:

```go
type LazyValue struct {
	once  sync.Once
	value interface{}
}

func NewLazyValue(fn func() interface{}) *LazyValue
func (lv *LazyValue) Get() interface{}
```

Write a table-driven test that verifies: multiple goroutines calling `Get` all see the same value, the function is called exactly once, and the function is not called on the second `Get`. Then write a benchmark comparing `LazyValue` vs. a mutex-based lazy init.

## Tests / verification

```bash
go run ./curriculum/modules/11-lifecycle-context-concurrency/lessons/29-sync-once
go test ./curriculum/modules/11-lifecycle-context-concurrency/lessons/29-sync-once
```

## Review questions

1. What is the time complexity of `sync.Once.Do` after initialization?
2. What happens if the function passed to `Do` panics? Can the function be retried?
3. Why is it a mistake to copy a `sync.Once` value?
4. How does `sync.Pool` differ from `sync.Once` in terms of lifecycle?
5. What pattern should you use instead of `sync.Once` when initialization can fail?

## NEXT UP

Congratulations on completing Module 11! You now understand time, context, goroutines, channels, synchronization, and patterns for robust concurrent programs. Next up: Module 12 — Backend Architecture.
