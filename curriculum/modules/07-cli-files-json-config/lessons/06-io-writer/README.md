# io.Writer

## Learning objective

Implement and compose `io.Writer` — Go's universal interface for writing data — and use combinators like `io.MultiWriter`, `bufio.Writer`, and `io.Copy`.

## Why this matters

`io.Writer` complements `io.Reader` as the output half of Go's I/O model. Files, network connections, HTTP responses, compression encoders, and hash functions all implement `io.Writer`. Composing writers lets you write to multiple destinations, buffer output, and compute checksums simultaneously with zero extra code.

## Mental model

An `io.Writer` is a sink. You pour data in with Write, and the writer decides where it goes: a file, a network buffer, a string builder, or a hash function. Like a funnel — the same liquid goes in, but the destination varies.

```
p []byte --> Writer --> [file | buffer | hash | network]
           Write(p)
```

## Core idea

The `io.Writer` interface is one method:

```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

- `Write` writes `len(p)` bytes from `p` to the underlying data stream.
- Returns the number of bytes written (`n`). If `n < len(p)`, it returns a non-nil error.
- Must not modify `p` or retain `p` after returning.
- Callers should always check the error and handle partial writes.

## Under the hood

For file and network writers, `Write` issues a write system call. The kernel copies data from user-space to the kernel buffer (or the device). The call blocks until at least one byte can be written. For in-memory writers (`bytes.Buffer`), `Write` is a `copy` into the buffer's backing array with possible reallocation.

`bufio.Writer` wraps an `io.Writer` and buffers writes in an internal 4096-byte buffer. `Flush()` must be called to ensure all buffered data is written to the underlying writer. This dramatically reduces system call overhead for small writes.

## How Go uses it

- `bytes.Buffer` implements `io.Writer` — useful for building strings.
- `os.File` implements `io.Writer` — write to files and stdout/stderr.
- `crypto/sha256.digest` implements `io.Writer` — write data to compute hash.
- `io.MultiWriter(w1, w2, ...)` — broadcasts each write to all writers.
- `io.TeeReader(r, w)` — every byte read from `r` is also written to `w`.
- `bufio.NewWriter(w)` — adds write buffering.
- `io.Copy(dst, src)` — reads from `src` and writes to `dst`, returns bytes copied.
- `fmt.Fprintf(w, format, args...)` — formatted write to any `io.Writer`.
- `io.WriteString(w, s)` — efficient string write if the writer has a `WriteString` method.

## Go example

```go
package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var buf strings.Builder
	hash := sha256.New()
	multi := io.MultiWriter(os.Stdout, &buf, hash)

	writer := bufio.NewWriter(multi)
	writer.WriteString("Hello, io.Writer!\n")
	fmt.Fprintf(writer, "Line %d\n", 2)
	writer.Flush()

	fmt.Printf("\nSHA-256: %s\n", hex.EncodeToString(hash.Sum(nil)))
}
```

## Step-by-step execution

1. `strings.Builder`, SHA-256 hash, and `os.Stdout` are combined with `io.MultiWriter`.
2. `bufio.NewWriter(multi)` creates a buffered writer around the multi-writer.
3. `writer.WriteString("Hello, io.Writer!\n")` buffers 17 bytes.
4. `fmt.Fprintf(writer, "Line %d\n", 2)` buffers 7 more bytes.
5. `writer.Flush()` writes all 24 buffered bytes through the multi-writer chain.
6. Each byte goes to stdout (visible), the string builder, and the hash.
7. `hash.Sum(nil)` returns the SHA-256 digest; `hex.EncodeToString` formats it.

## Common mistakes

- **Forgetting `Flush()` on `bufio.Writer`**: Data stays in the buffer and is lost on program exit. Always flush before closing.
- **Ignoring partial writes**: If `Write` returns `n < len(p)`, you must retry with `p[n:]`. Most standard writers are reliable, but network writers may not be.
- **Writing to a closed writer**: Panics or returns errors depending on implementation. Track writer lifecycle.
- **Assuming `Write` is safe for concurrent use**: Most writers are not. Use a mutex or dedicated goroutine for concurrent writes to the same writer.

## Debugging walkthrough

Buggy code:

```go
func main() {
	w := bufio.NewWriter(os.Stdout)
	w.WriteString("Hello\n")
	// No Flush!
}
```

**Symptom**: Nothing prints to stdout.

**Investigation**: Add `fmt.Println("after write")` — it prints. The string was buffered but never flushed.

**Fix**: Either call `w.Flush()` after writing, or use `os.Stdout` directly (or defer flush).

## Production notes

- Use `bufio.Writer` for all high-frequency writes to files and network connections.
- `io.MultiWriter` is useful for logging to both stdout and a file: `io.MultiWriter(os.Stdout, logFile)`.
- For structured logging, write directly to the `io.Writer` from your log library rather than going through `fmt.Print`.
- In HTTP handlers, `http.ResponseWriter` implements `io.Writer`, so you can use `io.Copy` to serve files.
- The `Write` method should never be called after `Close` returns.

## Performance implications

- Unbuffered small writes to `os.File` are expensive (one syscall per write). `bufio.Writer` batches them.
- `io.Copy` uses a 32KB buffer internally and is optimized with `io.CopyBuffer` for custom buffer sizes.
- `io.MultiWriter` writes to each writer sequentially — the slowest writer determines throughput.
- Hash writers compute checksums without allocation — they update internal state without copying.

## Practice task

Implement a `PrefixWriter` that wraps an `io.Writer` and prepends a prefix string to every write. For example, writing `"hello\nworld\n"` with prefix `"> "` should produce `"> hello\n> world\n"`. (Hint: detect newlines in the written data.) Save as `main.go`.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/06-io-writer
go test ./curriculum/modules/07-cli-files-json-config/lessons/06-io-writer
```

## Review questions

1. What does `io.Writer.Write` return when it writes all bytes successfully?
2. What happens if you call `bufio.Writer.Write` but never call `Flush`?
3. How does `io.MultiWriter` handle an error from one of its underlying writers?
4. What interface does `crypto/sha256.New()` implement that lets it be used as a writer?
5. Why would you use `io.Copy(dst, src)` instead of manually reading and writing in a loop?

## NEXT UP

Files — opening, creating, reading, and writing files with the `os` package.
