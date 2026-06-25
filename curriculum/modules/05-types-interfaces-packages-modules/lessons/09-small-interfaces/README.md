# Small interfaces

## Learning objective

Design and consume small, focused interfaces (1-2 methods), apply the Interface Segregation Principle in Go, and recognize the standard library's canonical examples.

## Why this matters

Small interfaces are the building blocks of composable Go programs. `io.Reader` (one method) and `io.Writer` (one method) are the most successful interfaces in the standard library — they are implemented by countless types and composed into higher-level abstractions. Learning to think in small interfaces transforms how you design APIs and package boundaries.

## Mental model

A small interface is a contract that does one thing. "I can read bytes." "I can close a resource." "I can compare for order." Because these contracts are tiny, many types satisfy them, and they can be mixed and matched freely. This is the Go implementation of the **Interface Segregation Principle**: no client should be forced to depend on methods it does not use.

## Core idea

**io.Reader**:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

**io.Writer**:

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

These are the canonical small interfaces. Each has one method. Together they form the foundation of Go's I/O model.

**Interface Segregation Principle (ISP)**: Instead of one large `ReadWriteSeekCloser` interface, Go defines many small ones: `Reader`, `Writer`, `Seeker`, `Closer`. Types implement only what they need.

**Standard library examples of small interfaces**:

| Interface | Methods | Purpose |
|---|---|---|
| `io.Reader` | `Read` | Read bytes |
| `io.Writer` | `Write` | Write bytes |
| `io.Closer` | `Close` | Close resource |
| `io.Seeker` | `Seek` | Seek to position |
| `fmt.Stringer` | `String` | String representation |
| `error` | `Error` | Error message |
| `sort.Interface` | `Len, Less, Swap` | Sortable collection |

## Under the hood

Small interfaces are more efficient at runtime because the `itab` (interface table) is smaller, but the performance difference is negligible. The real benefit is at the type-system level: more types satisfy a small interface, enabling greater code reuse.

## How Go uses it

- **io.Copy** takes `io.Writer` and `io.Reader` — works with files, network connections, buffers, HTTP response bodies, compressed streams, etc.
- **http.Handler** is one method: `ServeHTTP(ResponseWriter, *Request)`.
- **net.Conn** implements `Reader`, `Writer`, `Closer` separately.
- **os.File** implements `Reader`, `Writer`, `Seeker`, `Closer`.
- **bytes.Buffer** implements `Reader`, `Writer`, `Seeker`.
- **crypto.Hash** is a `io.Writer` that computes a hash.

## Go example

```go
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

type UpperWriter struct {
	w io.Writer
}

func (u UpperWriter) Write(p []byte) (int, error) {
	return u.w.Write(bytes.ToUpper(p))
}

func main() {
	// io.Reader from a string
	r := strings.NewReader("hello world")
	io.Copy(os.Stdout, r)
	fmt.Println()

	// Compose writers: gzip → file
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	gw.Write([]byte("compressed data"))
	gw.Close()
	fmt.Printf("gzip wrote %d bytes\n", buf.Len())

	// Custom writer: upper-case transform
	UpperWriter{w: os.Stdout}.Write([]byte("hello\n"))
}
```

## Step-by-step execution

For `io.Copy(os.Stdout, r)` where `r = strings.NewReader("hello world")`:

1. `io.Copy` receives `dst io.Writer` = `os.Stdout` (type `*os.File`) and `src io.Reader` = `strings.Reader`.
2. Inside `Copy`, a 32KB buffer is allocated.
3. Loop: `r.Read(buf)` copies `"hello world"` into the buffer.
4. `os.Stdout.Write(buf[:n])` writes the buffer to stdout.
5. Loop ends when `r.Read` returns `io.EOF`.
6. `io.Copy` returns `(n=11, err=nil)`.

## Common mistakes

- Mistake: Defining a large interface when small interfaces would suffice.
  - Fix: Split into multiple 1-2 method interfaces. Compose them as needed.

- Mistake: Implementing `io.Reader` incorrectly (wrong return values, not handling EOF properly).
  - Fix: `Read` should return `(n int, err error)` where `n` is the number of bytes read, and `err` is `io.EOF` only when no bytes were read.

- Mistake: Accepting `*os.File` when `io.Reader` would be more general.
  - Fix: Accept the smallest interface you need. Your function becomes testable with `strings.Reader` or `bytes.Buffer`.

- Mistake: Adding methods to an interface that are not needed by all consumers.
  - Fix: If only some callers need the extra method, define a separate interface.

## Debugging walkthrough

This custom `io.Reader` has a bug:

```go
type MyReader struct {
    data []byte
}

func (r *MyReader) Read(p []byte) (int, error) {
    copy(p, r.data)
    return len(r.data), nil // BUG: infinite loop
}

func main() {
    r := &MyReader{data: []byte("abc")}
    io.Copy(os.Stdout, r) // runs forever
}
```

**Symptom**: `io.Copy` never returns, printing "abcabcabc...".

**Root cause**: `Read` always returns `len(r.data) = 3` and `nil` error. The caller never receives `io.EOF`, so it calls `Read` again.

**Fix**: Return `io.EOF` when all data has been consumed. Track position:

```go
func (r *MyReader) Read(p []byte) (int, error) {
    if len(r.data) == 0 {
        return 0, io.EOF
    }
    n := copy(p, r.data)
    r.data = r.data[n:]
    return n, nil
}
```

## Production notes

- **Accept interfaces, return structs**. Input parameters should be as abstract (small interface) as possible; return types should be as concrete as necessary.
- **Test with small interface implementations**. Use `strings.Reader` and `bytes.Buffer` to test code that accepts `io.Reader` / `io.Writer`.
- **The `fs.FS` interface** (Go 1.16) is another great small interface: `Open(name string) (File, error)`.
- **Middleware pattern**: Wrap `io.Writer` with transform layers (compression, encryption, validation).

## Performance implications

- **io.Copy** uses a 32KB buffer by default. For high-throughput streaming, consider `io.CopyBuffer` with a custom buffer.
- **Small interface call overhead** is the same as any interface call (~1-2 ns). The win is in code reuse, not raw speed.
- **Wrapping io.Writer** in many layers adds indirection. For hot paths, consider composing at the buffer level.

## Practice task

Implement a custom `io.Writer` called `CountingWriter` that wraps another `io.Writer` and counts the total bytes written. Implement `BytesWritten() int64` to retrieve the count. In `main()`, write to `CountingWriter` wrapping `os.Stdout`, print a message, and output the count. Verify with tests.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/09-small-interfaces
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/09-small-interfaces
```

## Review questions

1. What is the single-method interface in Go that represents "I can read bytes"?
2. Why is `io.Reader` a better parameter type than `*os.File`?
3. What is the Interface Segregation Principle, and how does Go's approach to interfaces support it?
4. Name three standard library interfaces with only one method.
5. How does `io.Copy` know when to stop reading?

## NEXT UP

Interface embedding — composing interfaces from smaller interfaces.
