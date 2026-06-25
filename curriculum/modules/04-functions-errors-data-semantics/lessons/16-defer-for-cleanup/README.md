# defer for cleanup

## Learning objective

Use `defer` to release resources safely: close files, unlock mutexes, commit or roll back transactions, and clean up in error paths. Recognize and avoid defer-in-loop pitfalls.

## Why this matters

A file handle not closed, a mutex not unlocked, or a database transaction not rolled back corrupts state or leaks resources. Before `defer`, cleanup code was duplicated at every `return` statement, and adding a new return path meant adding new cleanup -- a formula for leaks. `defer` eliminates this by binding cleanup to function exit. It is the single most important tool for writing reliable resource-handling code in Go.

## Mental model

`defer` is a contract you make with the function: "When I acquire this resource, I promise to release it before I leave." By placing the release immediately after the acquire, you never forget. The cleanup is guaranteed to run regardless of how the function exits -- normal return, error return, or panic. This is RAII (Resource Acquisition Is Initialization) adapted to Go's explicit error-handling style.

## Core idea

The canonical cleanup pattern in Go:

```go
f, err := os.Open(path)
if err != nil {
    return err
}
defer f.Close()
// ... use f ...
```

This pattern works because:
1. `f.Close()` is deferred before any further code that might fail.
2. It runs on normal return, on error return, and on panic.
3. The `err` check happens before the defer, so `f` is never nil when `Close` is called.

This generalizes to mutexes, transactions, HTTP responses, temporary directories -- any pair of acquire/release operations.

## Under the hood

The compiler tracks `defer` statements and inserts cleanup code in the function epilogue. When the function returns (normally or via panic), the deferred functions execute in LIFO order. This means cleanup happens in reverse acquisition order -- exactly what you want when later resources depend on earlier ones (e.g., a transaction depends on a connection).

The runtime guarantee: `defer` always runs, even through panics. This is implemented in `runtime.gopanic`, which iterates the goroutine's defer list and executes each deferred function during stack unwinding.

## How Go uses it

```go
// Mutex cleanup
mu.Lock()
defer mu.Unlock()

// Transaction commit/rollback
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback() // no-op if Commit succeeds
if err := doWork(ctx, tx); err != nil {
    return err
}
return tx.Commit()

// HTTP response body
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()

// Temporary directory cleanup
dir, err := os.MkdirTemp("", "example")
if err != nil {
    return err
}
defer os.RemoveAll(dir)
```

## Go example

```go
package main

import (
	"errors"
	"fmt"
	"sync"
)

type Resource struct {
	Name string
}

func (r *Resource) Open() error {
	fmt.Printf("opening %s\n", r.Name)
	return nil
}

func (r *Resource) Close() error {
	fmt.Printf("closing %s\n", r.Name)
	return nil
}

type ResourceManager struct {
	mu  sync.Mutex
	res map[string]*Resource
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{res: make(map[string]*Resource)}
}

func (rm *ResourceManager) Acquire(name string) (*Resource, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if r, ok := rm.res[name]; ok {
		return r, nil
	}
	r := &Resource{Name: name}
	if err := r.Open(); err != nil {
		return nil, fmt.Errorf("open resource: %w", err)
	}
	rm.res[name] = r
	return r, nil
}

func processResource(rm *ResourceManager, name string) (err error) {
	r, err := rm.Acquire(name)
	if err != nil {
		return err
	}
	defer fmt.Println("deferred cleanup for", r.Name)

	// Simulate processing
	if name == "bad" {
		return errors.New("processing failed")
	}
	return nil
}

func main() {
	rm := NewResourceManager()
	for _, name := range []string{"good", "bad"} {
		err := processResource(rm, name)
		if err != nil {
			fmt.Printf("processResource(%q): %v\n", name, err)
		} else {
			fmt.Printf("processResource(%q): success\n", name)
		}
	}
}
```

## Step-by-step execution

For `processResource(rm, "bad")`:

1. `rm.Acquire("bad")` is called.
2. Inside Acquire: `rm.mu.Lock()`, then `defer rm.mu.Unlock()`.
3. The resource is created and stored.
4. `rm.Acquire` returns the resource.
5. Back in `processResource`: `defer fmt.Println("deferred cleanup for", r.Name)` is registered.
6. The `if name == "bad"` branch triggers, returning `errors.New("processing failed")`.
7. Before the function actually returns to the caller, the deferred statement runs: prints "deferred cleanup for bad".
8. The function returns the error to `main`.

Important: `defer` does NOT close the resource -- it just prints a message. In real code, you would defer `r.Close()` or similar. The key observation is that the deferred cleanup runs even on the error path.

## Common mistakes

- **Deferring close on a nil resource**: If `os.Open` returns an error, `f` is nil. `defer f.Close()` would panic because `Close` is called on a nil pointer. Always check the error before deferring.

```go
f, err := os.Open(path)
defer f.Close() // BUG: panics if err != nil
if err != nil {
    return err
}
```

- **Defer in a loop**: Each iteration defers a cleanup that runs only when the function returns, not at the end of the iteration. For large loops, this accumulates O(n) open resources.

```go
for _, path := range paths {
    f, err := os.Open(path)
    if err != nil { return err }
    defer f.Close() // all closes run at function return, not per iteration
    // process f
}
```

- **Shadowed variable in defer closure**: The loop variable capture bug:

```go
for _, name := range names {
    f, err := os.Open(name)
    if err != nil { return err }
    defer func() {
        f.Close() // captures f by reference — f changes each iteration!
    }()
}
```

- **Assuming defer runs before `os.Exit`**: `os.Exit` skips deferred functions. Use `log.Fatal` with care.

## Debugging walkthrough

```go
package main

import "fmt"

type DB struct {
	connected bool
}

func (db *DB) Connect() error {
	db.connected = true
	return nil
}

func (db *DB) Close() error {
	db.connected = false
	return nil
}

func process(db *DB) error {
	if err := db.Connect(); err != nil {
		return err
	}
	defer db.Close()
	return fmt.Errorf("processing error")
}

func main() {
	db := &DB{}
	err := process(db)
	fmt.Println("error:", err)
	fmt.Println("connected:", db.connected) // should be false
}
```

**Symptom**: `db.connected` is `false` after `process` returns, even though the function returned an error. The `defer db.Close()` ran.

**Investigation**: The deferred `Close()` runs after the `return fmt.Errorf(...)` sets the return value but before the function returns to the caller.

**Root cause**: None -- this is correct behavior. The defer ensures cleanup happens even on error.

**Fix**: None needed. This is the intended design.

## Production notes

- **Always pair acquire and defer in the same function**: When you acquire a resource, defer its release in the same function before doing anything that might cause a return. If you pass the resource to another function, that function owns the cleanup.

- **Transaction pattern**: `defer tx.Rollback()` before `tx.Commit()`. If `Commit` succeeds, `Rollback` becomes a no-op. If `Commit` fails, `Rollback` cleans up. This is the standard database transaction pattern in Go.

- **Mutex pattern**: `mu.Lock(); defer mu.Unlock()`. Never use a mutex without defer. The only exception is when the lock must be held across a function call boundary.

- **Error handling in deferred functions**: A deferred `Close()` that returns an error is typically logged, not returned (the function's return value is already set). If you need to handle close errors, use a named return:

```go
func readFile(path string) (data []byte, err error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = cerr
        }
    }()
    return io.ReadAll(f)
}
```

## Performance implications

- Each `defer` in a function without panic/recover is inlined at each return site (Go 1.14+), making defer essentially free for the common case.
- Deferring in a loop still defers at function scope, not loop scope. For iterating over many items, close resources explicitly at the end of each iteration instead of using defer.
- A deferred function call has the same cost as a direct function call (after the optimization). No additional allocation.

## Practice task

Write a function `writeToFile(filename, content string) error` that:

1. Creates the file with `os.Create`.
2. Defers `file.Close()`.
3. Writes the content with `file.WriteString`.
4. Returns any error from `WriteString`.
5. Returns `nil` on success.

Then write a function `safeProcess(items []string) error` that processes each item by calling `writeToFile`. For each item, open and close the file within the loop iteration (NOT using defer in the loop).

## Tests / verification

```bash
go run ./curriculum/modules/04-functions-errors-data-semantics/lessons/16-defer-for-cleanup
go test ./curriculum/modules/04-functions-errors-data-semantics/lessons/16-defer-for-cleanup
```

## Review questions

1. Why should `defer` be placed immediately after resource acquisition?
2. What happens if you defer `Close()` on a nil resource?
3. How does the transaction pattern `defer tx.Rollback()` followed by `tx.Commit()` work?
4. Why is defer in a loop problematic?
5. Can a deferred function modify the return value of the enclosing function?

## NEXT UP

panic and recover
