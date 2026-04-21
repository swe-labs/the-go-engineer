# Common Go Mistakes — Reference Guide

The 15 mistakes every Go engineer makes at least once, with the fix and the lesson where the correct pattern is taught.

> Section references use the format `Section ID: slug` from `curriculum.v2.json`.
> All cross-references have been updated to match the v2.1 architecture.
> Some referenced lessons are planned but not yet created. See `ROADMAP.md` for current status.

---

## 1. Capturing loop variables incorrectly in closures

**The bug:**

```go
for _, url := range urls {
    go func() {
        fmt.Println(url) // Prints the LAST url for every goroutine
    }()
}
```

**Why it happens:** All goroutines share the same `url` variable. By the time they run, the loop has finished and `url` holds the last value.

**The fix:**

```go
// Option 1: Shadow the variable — creates a new binding per iteration
for _, url := range urls {
    url := url
    go func() {
        fmt.Println(url)
    }()
}

// Option 2: Pass as a parameter (cleaner)
for _, url := range urls {
    go func(u string) {
        fmt.Println(u)
    }(url)
}
```

**Note:** Go 1.22+ fixes this for all `for`/`range` loop forms — each iteration now creates a new variable automatically. The explicit shadowing pattern is still useful to understand for pre-1.22 codebases and for clarity when reviewing older code.

**Taught in:** s03 (functions-errors) — FE.9: Closures; s07 (concurrency) — GC.1: Goroutines

---

## 2. Using `time.Sleep` to synchronise goroutines

**The bug:**

```go
go doWork()
time.Sleep(1 * time.Second) // Hoping the goroutine finishes in time
```

**Why it happens:** `time.Sleep` looks like "waiting." It isn't — it delays by an arbitrary duration with no guarantee.

**The fix:**

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    doWork()
}()
wg.Wait() // Blocks until Done() is called
```

**Taught in:** s07 (concurrency) — GC.2: WaitGroups

---

## 3. Not checking `rows.Err()` after iterating a database result set

**The bug:**

```go
rows, _ := db.Query("SELECT id FROM users")
defer rows.Close()
for rows.Next() { ... }
// Missing: rows.Err() check
```

**Why it happens:** The loop ends either when there are no more rows OR when an error occurs mid-stream (network drop, query cancelled). Without `rows.Err()` you silently return incomplete data.

**The fix:**

```go
for rows.Next() { ... }
if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("iterating rows: %w", err)
}
```

**Taught in:** s06 (backend-apis-databases) — DB.3: Select Queries

---

## 4. Forgetting `rows.Close()` — exhausting the connection pool

**The bug:**

```go
rows, err := db.Query("SELECT ...")
if err != nil { return err }
// No defer rows.Close() — connection held until GC
for rows.Next() { ... }
```

**Why it happens:** `rows` holds a live database connection. Without `Close()`, the connection is never returned to the pool. Under load, the pool empties and all new queries hang.

**The fix:**

```go
rows, err := db.Query("SELECT ...")
if err != nil { return err }
defer rows.Close() // Immediately after the nil-error check
```

**Taught in:** s06 (backend-apis-databases) — DB.3: Select Queries; s02 (language-basics) — CF.5: Defer

---

## 5. Compiling regex inside a loop

**The bug:**

```go
for _, line := range lines {
    re := regexp.MustCompile(`ERROR|WARN`) // Compiled on every iteration
    if re.MatchString(line) { ... }
}
```

**Why it happens:** `MustCompile` builds a DFA on every call. Inside a loop of 1 million lines this is 1 million DFA constructions.

**The fix:**

```go
var alertPattern = regexp.MustCompile(`ERROR|WARN`) // Package-level: compiled once

for _, line := range lines {
    if alertPattern.MatchString(line) { ... }
}
```

**Taught in:** s04 (types-design) — ST.4: Regex; s08 (quality-test) — PR.5: Benchmark-Driven Dev

---

## 6. String concatenation with `+=` in a loop

**The bug:**

```go
var result string
for _, word := range words {
    result += word + " " // O(N²): allocates a new string every iteration
}
```

**Why it happens:** Strings are immutable. Each `+=` creates a new string, copies all previous bytes, and releases the old string to GC.

**The fix:**

```go
var sb strings.Builder
for _, word := range words {
    sb.WriteString(word)
    sb.WriteByte(' ')
}
result := sb.String() // One allocation
```

**Taught in:** s04 (types-design) — ST.1: Strings

---

## 7. Passing a `sync.WaitGroup` by value

**The bug:**

```go
func doWork(wg sync.WaitGroup) { // COPY of the WaitGroup
    defer wg.Done()              // Decrements the copy, not the original
}
wg.Add(1)
go doWork(wg) // Original wg.Wait() blocks forever
```

**Why it happens:** `sync.WaitGroup` contains an internal counter. Passing by value copies the counter — changes inside the function don't affect the original.

**The fix:**

```go
func doWork(wg *sync.WaitGroup) { // POINTER to the WaitGroup
    defer wg.Done()
}
go doWork(&wg)
```

**Taught in:** s07 (concurrency) — GC.2: WaitGroups

---

## 8. Sending on a closed channel

**The bug:**

```go
close(ch)
ch <- value // panic: send on closed channel
```

**Why it happens:** The rule is: only the sender closes the channel. When multiple goroutines can send to the same channel, any one closing it can cause another to panic on the next send.

**The fix:**

```go
// Use sync.Once to ensure the channel is closed exactly once:
var once sync.Once
closeOnce := func() { once.Do(func() { close(ch) }) }

// Or: redesign so only one goroutine sends — close is then trivial.
```

**Taught in:** s07 (concurrency) — GC.5: Closing Channels; s07 (concurrency) — SY.2: sync.Once

---

## 9. Not using `fmt.Errorf("%w", err)` for error wrapping

**The bug:**

```go
if err != nil {
    return fmt.Errorf("database query failed: %v", err) // %v, not %w
}
// Caller cannot use errors.Is(err, sql.ErrNoRows)
```

**Why it happens:** `%v` formats the error as a string. `%w` wraps the original error so callers can inspect it with `errors.Is()` and `errors.As()`.

**The fix:**

```go
return fmt.Errorf("database query failed: %w", err) // %w preserves the chain
```

**Taught in:** s03 (functions-errors) — FE.4: Errors as Values; s04 (types-design) — TI.8: Custom Error Types

---

## 10. HTTP server with no timeouts

**The bug:**

```go
http.ListenAndServe(":8080", mux) // No timeouts — Slowloris vulnerability
```

**Why it happens:** The default `http.Server` has no read or write timeout. A client that sends headers slowly holds a connection open indefinitely.

**The fix:**

```go
server := &http.Server{
    Addr:              ":8080",
    Handler:           mux,
    ReadHeaderTimeout: 3 * time.Second,  // Prevent Slowloris
    ReadTimeout:       5 * time.Second,
    WriteTimeout:      30 * time.Second,
    IdleTimeout:       120 * time.Second,
}
server.ListenAndServe()
```

**Taught in:** s06 (backend-apis-databases) — HS.7: Server Timeouts

---

## 11. `http.DefaultClient` with no timeout

**The bug:**

```go
resp, err := http.Get(url) // Uses DefaultClient — no timeout
// If the server hangs, this goroutine hangs forever
```

**Why it happens:** `http.Get` uses `http.DefaultClient` which has `Timeout: 0` (infinite).

**The fix:**

```go
client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Get(url)
```

**Taught in:** s06 (backend-apis-databases) — HS.1: net/http Basics; s07 (concurrency) — CT.5: Timeout-Aware API Client

---

## 12. Ignoring the error from `defer file.Close()`

**The bug:**

```go
defer file.Close() // Error silently discarded
```

**Why it happens:** `Close()` flushes buffers. If the disk is full, `Close()` returns an error indicating data was not written — but `defer` discards return values.

**The fix for write paths:**

```go
defer func() {
    if cerr := file.Close(); cerr != nil && err == nil {
        err = cerr // Propagate close error only if no other error occurred
    }
}()
```

**For read paths:** `defer file.Close()` is fine — reads don't buffer data that can be lost on close.

**Taught in:** s05 (packages-io) — FS.1: Files; s02 (language-basics) — CF.5: Defer

---

## 13. Not cancelling a context after `WithTimeout`

**The bug:**

```go
ctx, _ := context.WithTimeout(parent, 5*time.Second) // cancel discarded
// The timer goroutine leaks until the timeout fires
```

**Why it happens:** Even when the timeout fires, the context's internal timer goroutine holds resources until `cancel()` is called. Discarding it leaks the goroutine.

**The fix:**

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel() // Always. Even for timeouts.
```

**Taught in:** s07 (concurrency) — CT.3: WithTimeout & WithDeadline

---

## 14. Concurrent map read/write without synchronisation

**The bug:**

```go
var cache = map[string]string{}
go func() { cache["key"] = "value" }() // Write
go func() { _ = cache["key"] }()       // Read — DATA RACE
```

**Why it happens:** Go maps are not safe for concurrent use. A concurrent read and write causes a fatal runtime error: `concurrent map read and map write`.

**The fix:**

```go
var mu sync.RWMutex
var cache = map[string]string{}

// Write:
mu.Lock()
cache["key"] = "value"
mu.Unlock()

// Read:
mu.RLock()
_ = cache["key"]
mu.RUnlock()
```

**Always run with `-race` to detect this:**

```bash
go test -race ./...
go run -race main.go
```

**Taught in:** s07 (concurrency) — SY.1: sync.Mutex & RWMutex; s07 (concurrency) — SY.4: Race Conditions

---

## 15. `log.Fatal` inside a goroutine or deferred function

**The bug:**

```go
go func() {
    if err := doWork(); err != nil {
        log.Fatal(err) // Calls os.Exit(1) — skips ALL deferred functions
    }
}()
```

**Why it happens:** `log.Fatal` calls `os.Exit(1)` immediately. This bypasses all `defer` calls — database connections, file handles, and in-flight requests are not cleaned up.

**The fix:**

```go
go func() {
    if err := doWork(); err != nil {
        log.Printf("worker error: %v", err) // Log, don't Fatal
        errCh <- err                        // Send to a coordinator
        return
    }
}()
```

Use `log.Fatal` only in `main()` during initialisation, before any deferred cleanup exists.

**Taught in:** s03 (functions-errors) — FE.10: panic & recover; s10 (production-operations) — GS.1: signal.NotifyContext

---

## Run These Checks Before Every Commit

```bash
go vet ./...          # Catch suspicious code patterns
go test -race ./...   # Catch data races
go build ./...        # Verify everything compiles
staticcheck ./...     # Install: go install honnef.co/go/tools/cmd/staticcheck@latest
```
