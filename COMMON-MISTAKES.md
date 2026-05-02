# Common Go Mistakes - Reference Guide

This guide lists common Go mistakes, why they happen, and where the correct pattern is taught.

Section references use the locked v2.1 section map and lesson IDs from `curriculum.v2.json`.

Stable architecture means the section map and curriculum contract are locked. Individual lesson depth may still grow during post-v2.1 implementation.

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

**Taught in:** s03 FE.9, s07 GC.1

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

**Taught in:** s07 GC.2

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

**Taught in:** s06 DB.3

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

**Taught in:** s06 DB.3, s02 CF.5

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

**Taught in:** s04 ST.4, s08 PR.5

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

**Taught in:** s04 ST.1

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

**Taught in:** s07 GC.2

## 8. Sending on a closed channel

**The bug:**

```go
close(ch)
ch <- value
```

**The fix:** only the sending owner closes the channel, or redesign ownership so one goroutine is responsible for closing.

**Taught in:** s07 GC.5, s07 SY.2

## 9. Not using `%w` for error wrapping

**The bug:**

```go
return fmt.Errorf("query failed: %v", err)
```

**The fix:**

```go
return fmt.Errorf("query failed: %w", err)
```

**Taught in:** s03 FE.4, s04 TI.8

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

**Taught in:** s06 HS.7

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

**Taught in:** s06 HS.1, s07 CT.5

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

**Taught in:** s05 FS.1, s02 CF.5

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

**Taught in:** s07 CT.3

## 14. Concurrent map read/write without synchronization

**The bug:**

```go
cache["key"] = "value"
_ = cache["key"]
```

from different goroutines.

**The fix:** protect the map with `sync.Mutex`, `sync.RWMutex`, channels, or `sync.Map` when appropriate.

**Taught in:** s07 SY.1, s07 SY.4

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

**Taught in:** s03 FE.10, s10 GS.1

## Required Checks

```bash
go build ./...
go vet ./...
unformatted=$(gofmt -l .); test -z "$unformatted" || (echo "$unformatted" && exit 1)
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
go run ./scripts/validate_curriculum.go
```

On PowerShell, quote the coverage flag if needed:

```powershell
go test "-coverprofile=coverage.out" ./...
```

Recommended when installed:

```bash
staticcheck ./...
```
