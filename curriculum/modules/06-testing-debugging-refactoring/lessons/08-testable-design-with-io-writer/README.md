# Testable design with io.Writer

## Learning objective

Design functions that accept `io.Writer` for testable output generation, use `bytes.Buffer` as a test spy, and verify log output by capturing writes.

## Why this matters

Functions that write to stdout (`fmt.Print`, `os.Stdout`) are difficult to test. You cannot capture their output without global state manipulation (`os.Stdout = ...`), which is fragile and non-parallel-safe. By accepting an `io.Writer` parameter, the same function becomes testable: in production you pass `os.Stdout`, in tests you pass a `bytes.Buffer`. This is the single most impactful design change for testability in Go.

## Mental model

`io.Writer` is a pipe. Data goes in one end. What happens at the other end depends on what you attach:

- `os.Stdout` → data appears in the terminal.
- `bytes.Buffer` → data is captured in memory for assertions.
- `io.Discard` → data is thrown away (for benchmarks or when output is irrelevant).
- `net.Conn` → data goes over the network.
- `*os.File` → data is written to disk.

By accepting `io.Writer`, your function does not care what is on the receiving end. It only needs to satisfy `Write(p []byte) (n int, err error)`.

```
Function → Write(data) → io.Writer → os.Stdout (prod)
                                  → bytes.Buffer (test)
```

## Core idea

Instead of:

```go
func Greet(name string) {
    fmt.Printf("Hello, %s!\n", name) // hard to test
}
```

Write:

```go
func Greet(w io.Writer, name string) {
    fmt.Fprintf(w, "Hello, %s!\n", name) // testable
}
```

The `fmt.Fprint*` family of functions is the key: they write to any `io.Writer`. In production:

```go
Greet(os.Stdout, "Alice")
```

In tests:

```go
var buf bytes.Buffer
Greet(&buf, "Alice")
if buf.String() != "Hello, Alice!\n" {
    t.Errorf(...)
}
```

This pattern applies to any output-producing code: logging, report generation, serialization, template execution.

## Under the hood

`io.Writer` is an interface with one method:

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

`bytes.Buffer` implements `io.Writer` by appending to an internal byte slice. `os.File` implements it by writing to the underlying file descriptor. Your function calls `Write` (possibly through `fmt.Fprintf`) and the implementation handles storage or transmission.

`bytes.Buffer` also implements `fmt.Stringer` via `String()`, which returns the accumulated data. This makes assertions straightforward.

## How Go uses it

The principle "accept interfaces, return structs" is central to Go design. `io.Writer` appears in:

- `fmt.Fprintf(w, ...)` — formatted output to any writer.
- `log.New(w, ...)` — logger that writes to a writer.
- `net/http.ResponseWriter` — the HTTP response writer is an `io.Writer`.
- `io.MultiWriter(w1, w2)` — writes to multiple destinations simultaneously.

Production Go codebases routinely design functions to accept `io.Writer` or `io.Reader` for testability.

## Go example

```go
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

// Greet writes a greeting to w.
func Greet(w io.Writer, name string) {
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

// NewLogger creates a logger that writes to w with a prefix.
func NewLogger(w io.Writer, prefix string) *log.Logger {
	return log.New(w, prefix, log.LstdFlags)
}

func main() {
	Greet(os.Stdout, "Alice")

	var buf bytes.Buffer
	logger := NewLogger(&buf, "APP: ")
	logger.Println("started")
	fmt.Print("Captured log: ", buf.String())
}
```

```go
// main_test.go
package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestGreet(t *testing.T) {
	var buf bytes.Buffer
	Greet(&buf, "Alice")

	want := "Hello, Alice!\n"
	got := buf.String()
	if got != want {
		t.Errorf("Greet = %q; want %q", got, want)
	}
}

func TestGreetEmptyName(t *testing.T) {
	var buf bytes.Buffer
	Greet(&buf, "")

	want := "Hello, !\n"
	got := buf.String()
	if got != want {
		t.Errorf("Greet(\"\") = %q; want %q", got, want)
	}
}

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, "TEST: ")
	logger.Println("hello world")

	got := buf.String()
	if !strings.Contains(got, "TEST: hello world") {
		t.Errorf("log = %q; want containing %q", got, "TEST: hello world")
	}
}

func TestMultiWrite(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	fmt.Fprint(&buf1, "data1")
	fmt.Fprint(&buf2, "data2")

	if buf1.String() != "data1" {
		t.Errorf("buf1 = %q; want %q", buf1.String(), "data1")
	}
	if buf2.String() != "data2" {
		t.Errorf("buf2 = %q; want %q", buf2.String(), "data2")
	}
}
```

## Step-by-step execution

Running `TestGreet`:

1. `var buf bytes.Buffer` creates an empty buffer.
2. `Greet(&buf, "Alice")` is called.
3. Inside `Greet`, `fmt.Fprintf(w, "Hello, %s!\n", "Alice")` calls `w.Write([]byte("Hello, Alice!\n"))`.
4. `bytes.Buffer.Write` appends `"Hello, Alice!\n"` to its internal slice.
5. `buf.String()` returns `"Hello, Alice!\n"`.
6. Assertion: `got == want`. Pass.

Compare with a non-testable implementation that called `fmt.Print` directly — testing would require capturing `os.Stdout` with a replacement pipe, which is fragile and cannot run in parallel.

## Common mistakes

- **Not accepting `io.Writer` at the deepest level.** If `Greet` calls another function that also writes to stdout, both functions should accept `io.Writer` and pass it down.

- **Capturing `os.Stdout` by reassignment.** Some developers do `old := os.Stdout; os.Stdout = buf; defer func() { os.Stdout = old }()`. This is fragile, not parallel-safe, and unreliable on Windows. Accept `io.Writer` instead.

- **Forgetting to check write errors.** `io.Writer.Write` can return an error. In most tests, you can ignore it, but in production code, check and handle the error.

- **Using `io.Writer` when `io.WriteString` or `io.WriterTo` is more appropriate.** For string-heavy output, prefer `io.WriteString(w, s)` which avoids allocation.

- **Not resetting the buffer between subtests.** If you reuse a buffer across subtests, call `buf.Reset()` to clear it.

## Debugging walkthrough

A test captures log output but the assertion fails:

```go
func TestLogOutput(t *testing.T) {
    var buf bytes.Buffer
    logger := log.New(&buf, "MYAPP: ", 0)
    logger.Print("test message")
    if buf.String() != "MYAPP: test message\n" {
        t.Errorf("got %q", buf.String())
    }
}
```

**Symptom**: The test fails. `buf.String()` shows `MYAPP: test message` without a newline, or with a different format.

**Investigation**: `log.Print` adds a newline only when the message does not end with one. `log.Println` always adds a newline. The `log.New` flags (third argument) control prefix format — `0` means no date/time. The assertion expects a newline but `log.Print` behavior differs.

**Root cause**: `log.Print` appends `\n` only if the message lacks it. The test expectation is mismatched with `log` package behavior.

**Fix**: Use `log.Println` for consistent newline behavior, or adjust the test expectation:
```go
// log.Println adds a newline
logger.Println("test message")
// log.Print adds only if missing
logger.Print("test message\n")
```

## Production notes

- **Interface parameters are documentation.** Accepting `io.Writer` signals that the function produces output and does not care where it goes.
- **`os.Stdout` is a default, not a requirement.** Provide convenience constructors that set `io.Writer` to `os.Stdout`, but always allow overriding.
- **Loggers are typically configured once.** Pass a `*log.Logger` (which wraps an `io.Writer`) rather than `io.Writer` directly when log formatting is needed.
- **`io.MultiWriter`** is useful for writing to both a file and stdout simultaneously.
- **Failing to write** (disk full, broken pipe) should be handled. In CLI tools, check write errors. In tests, `bytes.Buffer.Write` never fails.

## Performance implications

- `bytes.Buffer` grows dynamically. For large output (MB+), consider pre-allocating with `bytes.NewBuffer(make([]byte, 0, estimatedSize))`.
- `fmt.Fprintf` allocates. For performance-critical paths, use `w.Write([]byte(...))` directly.
- `bytes.Buffer.String()` does not allocate — it returns a view of the underlying slice.
- `io.Discard` is a no-op writer for benchmarks where output is irrelevant. Use it to measure the computation cost without IO.

## Practice task

Write a function `WriteTable(w io.Writer, rows [][2]string)` that writes a two-column table:
```
Name    | Age
--------|-----
Alice   | 30
Bob     | 25
```

Use `fmt.Fprintf`. Write tests that:
- Verify the table header and separator.
- Verify an empty rows list produces only headers.
- Verify alignment for varying name lengths.

Use `bytes.Buffer` to capture output in all tests.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/08-testable-design-with-io-writer
go test -v ./curriculum/modules/06-testing-debugging-refactoring/lessons/08-testable-design-with-io-writer
```

## Review questions

1. Why is `fmt.Printf` harder to test than `fmt.Fprintf`?
2. What interface does `bytes.Buffer` implement that makes it useful for capturing output?
3. How would you refactor a function that writes to stdout to be testable without changing its public API?
4. What is `io.Discard` and when would you use it?
5. What are the risks of capturing output by reassigning `os.Stdout`?

## NEXT UP

Fakes before mocks — in-memory fake implementations as an alternative to mocking frameworks.
