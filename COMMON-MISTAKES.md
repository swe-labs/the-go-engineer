# Common Go Mistakes — Reference Guide

This guide catalogs common Go engineering mistakes, underlying mechanisms, and defensive patterns taught throughout the curriculum.

Defensive design principles and failure modes are formally tracked in [`curriculum/failures.json`](./curriculum/failures.json) and demonstrated throughout canonical modules.

## 1. Capturing loop variables incorrectly in closures

**The bug:**

```go
for _, url := range urls {
	go func() {
		fmt.Println(url)
	}()
}
```

**The fix:**

```go
for _, url := range urls {
	url := url
	go func() {
		fmt.Println(url)
	}()
}
```

Go 1.22+ fixes this for `for`/`range` loop variables, but the pattern remains useful for understanding older code.

**Taught in:** Module 04: Functions (`core-04-08`), Module 11: Goroutines (`core-11-11`)

## 2. Using `time.Sleep` to synchronize goroutines

**The bug:**

```go
go doWork()
time.Sleep(time.Second)
```

**The fix:**

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
	defer wg.Done()
	doWork()
}()
wg.Wait()
```

**Taught in:** Module 11: WaitGroups & Goroutine Coordination (`core-11-12`)

## 3. Not checking `rows.Err()`

**The bug:**

```go
for rows.Next() {
	// scan rows
}
```

**The fix:**

```go
for rows.Next() {
	// scan rows
}
if err := rows.Err(); err != nil {
	return fmt.Errorf("iterate rows: %w", err)
}
```

**Taught in:** Module 09: SQL Queries & Error Inspection (`core-09-10`, `core-09-13`)

## 4. Forgetting `rows.Close()`

**The bug:**

```go
rows, err := db.QueryContext(ctx, query)
if err != nil {
	return err
}
// missing rows.Close()
```

**The fix:**

```go
rows, err := db.QueryContext(ctx, query)
if err != nil {
	return err
}
defer rows.Close()
```

**Taught in:** Module 04: Defer for Cleanup (`core-04-16`), Module 09: Connection Pool Lifecycle (`core-09-11`, `core-09-13`)

## 5. Compiling regex inside a loop

**The bug:**

```go
for _, line := range lines {
	re := regexp.MustCompile(`ERROR|WARN`)
	_ = re.MatchString(line)
}
```

**The fix:**

```go
var alertPattern = regexp.MustCompile(`ERROR|WARN`)
```

**Taught in:** Module 07: CLI, Files & Text Processing (`core-07-06`), Module 13: Benchmark-Driven Optimization (`core-13-03`)

## 6. String concatenation with `+=` in a loop

**The bug:**

```go
var result string
for _, word := range words {
	result += word
}
```

**The fix:**

```go
var builder strings.Builder
for _, word := range words {
	builder.WriteString(word)
}
result := builder.String()
```

**Taught in:** Module 03: String and Memory Layouts (`core-03-19`), Module 13: Zero-Allocation Engineering (`core-13-04`)

## 7. Passing a `sync.WaitGroup` by value

**The bug:**

```go
func work(wg sync.WaitGroup) {
	defer wg.Done()
}
```

**The fix:**

```go
func work(wg *sync.WaitGroup) {
	defer wg.Done()
}
```

**Taught in:** Module 05: Pointer vs Value Receivers (`core-05-03`), Module 11: WaitGroups (`core-11-12`)

## 8. Sending on a closed channel

**The bug:**

```go
close(ch)
ch <- value
```

**The fix:** only the sending owner closes the channel, or redesign ownership so one goroutine is responsible for closing.

**Taught in:** Module 11: Closing Channels & Channel Ownership (`core-11-15`, `core-11-16`)

## 9. Not using `%w` for error wrapping

**The bug:**

```go
return fmt.Errorf("query failed: %v", err)
```

**The fix:**

```go
return fmt.Errorf("query failed: %w", err)
```

**Taught in:** Module 04: Errors as Values & Custom Error Types (`core-04-09`, `core-04-18`)

## 10. HTTP server with no timeouts

**The bug:**

```go
http.ListenAndServe(":8080", mux)
```

**The fix:**

```go
server := &http.Server{
	Addr:              ":8080",
	Handler:           mux,
	ReadHeaderTimeout: 3 * time.Second,
	ReadTimeout:       5 * time.Second,
	WriteTimeout:      30 * time.Second,
	IdleTimeout:       120 * time.Second,
}
```

**Taught in:** Module 08: Server Timeouts & Production HTTP (`core-08-14`)

## 11. HTTP client with no timeout

**The bug:**

```go
resp, err := http.Get(url)
```

**The fix:**

```go
client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Get(url)
```

**Taught in:** Module 08: HTTP Clients (`core-08-02`, `core-08-15`), Module 11: Context Deadlines & Cancellation (`core-11-06`)

## 12. Ignoring close errors on write paths

**The bug:**

```go
defer file.Close()
```

**The fix for write paths:**

```go
defer func() {
	if closeErr := file.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
}()
```

**Taught in:** Module 04: Defer Mechanics (`core-04-16`), Module 07: Files and Filesystem I/O (`core-07-07`)

## 13. Not canceling a context

**The bug:**

```go
ctx, _ := context.WithTimeout(parent, 5*time.Second)
```

**The fix:**

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()
```

**Taught in:** Module 11: Context Cancellation and Timeouts (`core-11-06`, `core-11-07`)

## 14. Concurrent map read/write without synchronization

**The bug:**

```go
cache["key"] = "value"
_ = cache["key"]
```

from different goroutines.

**The fix:** protect the map with `sync.Mutex`, `sync.RWMutex`, channels, or `sync.Map` when appropriate.

**Taught in:** Module 03: Maps (`core-03-17`), Module 11: Mutex, RWMutex & Race Detection (`core-11-19`, `core-11-21`)

## 15. `log.Fatal` inside a goroutine

**The bug:**

```go
go func() {
	if err := doWork(); err != nil {
		log.Fatal(err)
	}
}()
```

`log.Fatal` calls `os.Exit(1)` and skips deferred cleanup.

**The fix:** report the error to a coordinator and let `main` decide how to shut down.

**Taught in:** Module 04: Panic and Recover Mechanics (`core-04-17`), Module 11: Graceful Shutdown & Goroutine Leaks (`core-11-04`, `core-11-23`)

## Required Checks

```bash
go build ./...
go vet ./...
go test ./...
go test -race ./...
go run ./scripts/validate_curriculum.go
```

Recommended when installed:

```bash
staticcheck ./...
```
