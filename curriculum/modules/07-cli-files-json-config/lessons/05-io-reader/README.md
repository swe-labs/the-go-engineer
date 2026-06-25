# io.Reader

## Learning objective

Implement and compose `io.Reader` — Go's universal interface for reading data — and chain readers through standard library combinators.

## Why this matters

`io.Reader` is the single most important interface in Go. It represents anything you can read from: a file, a network socket, an HTTP response body, a compressed stream, an encrypted blob, or a string. Every I/O library in Go speaks `io.Reader`. Mastering it unlocks all of Go's streaming data processing.

## Mental model

An `io.Reader` is a tap. You attach a bucket (`[]byte` slice), turn the handle (`Read`), and water flows in. The tap tells you how much water you got (n) and whether there's more coming (nil error) or the pipe is dry (io.EOF).

```
         +--------+
data --> | Reader | --> p []byte
         +--------+
         (n int, err error)
```

## Core idea

The `io.Reader` interface is one method:

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

- `Read` fills `p` with up to `len(p)` bytes.
- Returns the number of bytes read (`n`).
- Returns `io.EOF` when there is no more data (after reading the last bytes, or on a subsequent call).
- May return `n > 0` with an error (e.g., EOF after reading the last bytes).
- Must not return `n > 0` with a nil `p`; may return `0, nil` (nothing to read, no error).

## Under the hood

`Read` is a blocking syscall under the hood for file and network readers. The Go runtime uses non-blocking I/O internally with goroutine scheduling, but the `Read` method appears synchronous. For buffered readers (`bufio.Reader`), the actual `Read` call on the underlying reader is deferred until the buffer is empty.

The `io.EOF` sentinel error is a special value created with `errors.New("EOF")`. Callers should compare with `== io.EOF` (not `!= nil`) to distinguish end-of-stream from actual errors.

## How Go uses it

- `strings.NewReader(s)` — wraps a string as an `io.Reader`.
- `bytes.NewReader(b)` — wraps a `[]byte` as an `io.Reader`.
- `os.File` implements `io.Reader` (and `io.Writer`).
- `bufio.NewReader(r)` — adds buffering to any reader.
- `compress/gzip.NewReader(r)` — wraps a reader to decompress gzip data.
- `io.MultiReader(r1, r2, ...)` — concatenates multiple readers sequentially.
- `io.LimitReader(r, n)` — restricts reading to `n` bytes.
- `io.TeeReader(r, w)` — reads from `r` and writes everything read to `w`.
- `io.ReadAll(r)` — reads all remaining data from `r` into a `[]byte`.

## Go example

```go
package main

import (
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

func main() {
	data := "Hello, io.Reader!"

	var b strings.Builder
	gzWriter := gzip.NewWriter(base64.NewEncoder(base64.StdEncoding, &b))
	gzWriter.Write([]byte(data))
	gzWriter.Close()

	encoded := b.String()

	base64Reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encoded))
	gzReader, _ := gzip.NewReader(base64Reader)
	defer gzReader.Close()

	decoded, _ := io.ReadAll(gzReader)
	fmt.Println("Decoded:", string(decoded))
}
```

## Step-by-step execution

For `r := strings.NewReader("hi"); buf := make([]byte, 4); n, err := r.Read(buf)`:

1. `strings.Reader` has 2 bytes remaining ("hi").
2. `Read(buf)` copies "hi" into `buf[0:2]`, sets `n = 2`, returns `err = nil`.
3. Call `r.Read(buf)` again: no data left, copies 0 bytes, returns `n = 0, err = io.EOF`.

For `io.ReadAll(r)`:

1. Creates a `bytes.Buffer`.
2. Loops: `Read` into a 512-byte temp buffer, `Write` the bytes into the buffer.
3. Stops when `Read` returns `io.EOF`.
4. Returns the buffer's bytes.

## Common mistakes

- **Assuming `Read` fills the buffer**: `Read` may return fewer bytes than `len(p)` even without error. Always use `n`, do not assume `n == len(p)`.
- **Ignoring the return value `n`**: Using `buf` directly after `Read` without checking `n` processes garbage bytes.
- **Using `err != nil` to detect EOF**: `Read` may return `n > 0` with `io.EOF` — process the bytes first, then check for EOF.
- **Reusing the same buffer incorrectly**: Reset or re-slice the buffer between reads to avoid stale data.
- **Forgetting to close the reader**: Some readers (gzip, network) hold OS resources. Always `defer closer.Close()`.

## Debugging walkthrough

Buggy code:

```go
func readAll(r io.Reader) []byte {
	var buf []byte
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		buf = append(buf, buf[:n]...)
	}
	return buf
}
```

**Symptom**: Infinite loop or zero bytes returned.

**Investigation**: `buf` starts as nil (len 0). `Read(buf)` with a nil/empty slice returns `0, nil` immediately — `buf` never grows.

**Fix**: Allocate a fixed-size buffer: `buf := make([]byte, 1024)`, read into it, then append `buf[:n]` to the result.

## Production notes

- Prefer `io.ReadAll` or `io.Copy` over manual read loops for most use cases.
- Always check the `n` return value before using read data.
- Use `bufio.Reader` for small reads from network sources to reduce syscalls.
- Implement `io.ReaderFrom` on your types if you can optimize bulk reads.
- For streaming protocols, consider using `io.LimitReader` to prevent unbounded memory allocation.

## Performance implications

- One `Read` call on an unbuffered `os.File` issues one system call. Use `bufio.Reader` to batch reads.
- `io.ReadAll` grows a buffer exponentially (doubling each time), so amortized cost is O(n).
- Wrapping readers (gzip on top of TLS on top of TCP) adds CPU overhead per layer but enables clean composability.
- Zero-copy reads are possible with `io.ReaderAt` and `io.WriterTo` but not with the basic `Reader` interface.

## Practice task

Implement a `Rot13Reader` that wraps an `io.Reader` and applies ROT13 substitution to every byte (add 13 to A-Z/a-z, wrapping). Import it from a local package or define it in main. Use it to read and transform: `"Hello, World!"` → `"Uryyb, Jbeyq!"`. Save as `main.go`.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/05-io-reader
go test ./curriculum/modules/07-cli-files-json-config/lessons/05-io-reader
```

## Review questions

1. What does `io.Reader.Read` return when the source has 3 bytes available and the buffer is 10 bytes?
2. What is the difference between `io.EOF` and a non-nil error from `Read`?
3. Why is the return value `n` important even when `err == nil`?
4. How would you compose two readers so they appear as one continuous stream?
5. What happens if you call `Read` on an `io.Reader` with a zero-length buffer?

## NEXT UP

`io.Writer` — the counterpart interface for writing data.
