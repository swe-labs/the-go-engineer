# Bytes

## Learning objective

Work with `[]byte` as the universal data buffer in Go, convert between strings and byte slices, and use the `bytes` package for searching, splitting, joining, and replacing byte data.

## Why this matters

Everything in a computer is bytes — network packets, file contents, encryption keys, images, protocol buffers, JSON, and UTF-8 text. Go's `[]byte` is the standard type for raw data. Every I/O operation, every wire format, and every serialization library uses `[]byte`. If you cannot work fluently with bytes, you cannot write real Go programs.

## Mental model

A `byte` is just a `uint8` — an 8-bit value from 0 to 255. A `[]byte` is a slice of these unsigned integers. You can index it, slice it, append to it, range over it, and pass it everywhere you use `[]int`, with one extra trick: `string` and `[]byte` interconvert efficiently because Go's strings are UTF-8 encoded byte sequences under the hood.

Think of `[]byte` as the lingua franca of data exchange. A string is the human-readable interpretation of bytes; `[]byte` is the raw machine representation.

## Core idea

- `byte` is an alias for `uint8`. The two are interchangeable; `byte` signals "raw data" vs `uint8` which signals "small number".
- `[]byte` is a slice of bytes — the standard Go type for I/O buffers, encoding output, and protocol data.
- String ↔ byte conversion: `[]byte("hello")` and `string([]byte{104, 101, 108, 108, 111})`.
- The `bytes` package provides utilities: `Contains`, `Split`, `Join`, `Replace`, `HasPrefix`, `HasSuffix`, `Trim`, and `Buffer`.

## Under the hood

A Go string is an immutable `{ptr, len}` header pointing to read-only memory. Converting `string` to `[]byte` allocates a new byte slice and copies the data, because `[]byte` must be mutable. The compiler optimizes this when it can prove the byte slice is never modified (escape analysis), but in general:

- `[]byte(s)` — O(n) copy. Produces a mutable copy.
- `string(b)` — O(n) copy. Produces an immutable copy.

The `bytes.Buffer` type wraps a `[]byte` with a read offset, implementing `io.Reader` and `io.Writer`. It grows the internal buffer as needed, same as append.

## How Go uses it

- **I/O**: `io.Reader` and `io.Writer` both operate on `[]byte`. `os.File.Read(buf)`, `net.Conn.Read(buf)`.
- **Encoding**: `encoding/json`, `encoding/gob`, `encoding/xml` marshal to / unmarshal from `[]byte`.
- **Crypto**: `crypto/sha256.Sum256(data)` takes `[]byte`. All hash and cipher functions use `[]byte`.
- **Networking**: `net` package reads and writes `[]byte` on connections.
- **Text processing**: `bytes` package mirrors `strings` package but for `[]byte` operands.

## Go example

```go
package main

import (
	"bytes"
	"fmt"
)

func main() {
	// byte is uint8
	var b byte = 65
	fmt.Println("byte as char:", string(b), "as int:", b)

	// []byte from string
	data := []byte("hello world")
	fmt.Println("data:", data)
	fmt.Println("data as string:", string(data))

	// bytes package utilities
	fmt.Println("Contains 'world':", bytes.Contains(data, []byte("world")))
	fmt.Println("HasPrefix 'hello':", bytes.HasPrefix(data, []byte("hello")))

	// Split and Join
	parts := bytes.Split(data, []byte(" "))
	fmt.Println("Split by space:", parts)
	joined := bytes.Join(parts, []byte(","))
	fmt.Println("Joined with comma:", string(joined))

	// Replace
	replaced := bytes.Replace(data, []byte("world"), []byte("go"), 1)
	fmt.Println("Replaced:", string(replaced))

	// bytes.Buffer
	var buf bytes.Buffer
	buf.WriteString("count: ")
	buf.Write([]byte{49, 50, 51})
	fmt.Println("Buffer:", buf.String())
}
```

## Step-by-step execution

For `data := []byte("hello")`:

1. Compiler sees `string` → `[]byte` conversion.
2. Runtime allocates a new backing array of length 5 on the heap (or stack if small and non-escaping).
3. Copies bytes `0x68 0x65 0x6c 0x6c 0x6f` into the new array.
4. Slice header: `{Data=&arr[0], Len=5, Cap=5}`.

For `buf.WriteString("go")` on a `bytes.Buffer`:

1. `buf`'s internal `[]byte` has length `L` and capacity `C`.
2. `WriteString` checks: `L + 2 > C`? If yes, grow capacity (same algorithm as append).
3. Copy "go" bytes into `buf.buf[L:L+2]`.
4. Set `buf.buf = buf.buf[:L+2]`.

For `bytes.Contains(data, []byte("world"))`:

1. Calls `bytes.Index(data, []byte("world"))`.
2. Uses optimized Rabin-Karp or two-way string matching for pattern search.
3. Returns `true` if Index >= 0.

## Common mistakes

- Mistake: Assuming `string` to `[]byte` conversion is free.
  - Why it happens: Strings and byte slices look similar, but strings are immutable and slices are mutable. The conversion must copy.
  - Fix: Use `unsafe.Slice(unsafe.StringData(s), len(s))` in performance-critical code where you can guarantee immutability. But prefer the safe conversion.

- Mistake: Modifying `[]byte` returned by `bytes.Buffer.Bytes()`.
  - Why it happens: `Bytes()` returns a slice sharing the buffer's internal array. Writing to it corrupts the buffer.
  - Fix: Copy immediately: `cp := append([]byte(nil), buf.Bytes()...)`.

- Mistake: Confusing `[]byte` with `string` in map keys.
  - Why it happens: `[]byte` is not comparable, `string` is. You cannot use `[]byte` as a map key.
  - Fix: Convert to string: `m[string(key)]`.

- Mistake: Not resetting `bytes.Buffer` in a loop, causing unbounded growth.
  - Why it happens: Each `Write` appends. Without `buf.Reset()`, the buffer keeps growing.
  - Fix: Call `buf.Reset()` between iterations, or use `buf = new(bytes.Buffer)`.

## Debugging walkthrough

```go
package main

import (
	"bytes"
	"fmt"
)

func main() {
	var buf bytes.Buffer
	for i := 0; i < 3; i++ {
		buf.WriteString(fmt.Sprintf("line %d\n", i))
	}
	data := buf.Bytes()
	data[0] = 'X' // intentional corruption to show the bug
	fmt.Println(buf.String())
}
```

**Symptom**: Prints `Xine 0\nline 1\nline 2\n` because `data[0] = 'X'` modified the buffer's internal array.

**Investigation**: The `data` slice returned by `buf.Bytes()` shares the backing array with `buf`. Print addresses:

```go
fmt.Printf("buf internal: %p\ndata: %p\n", buf.Bytes(), data)
```

Both are the same pointer.

**Root cause**: `buf.Bytes()` documents that the returned slice shares memory. The write corrupts the buffer.

**Fix**: Copy the data: `data := append([]byte(nil), buf.Bytes()...)`.

## Production notes

- **Buffer reuse**: Use `sync.Pool` with `bytes.Buffer` in high-throughput servers to reduce allocation.
- **Preallocate buffers**: `buf := make([]byte, 0, 4096)` for I/O reduces growth overhead.
- **`bytes.Buffer` vs `strings.Builder`**: Use `strings.Builder` when building strings (avoids `[]byte` to `string` conversion on `String()`). Use `bytes.Buffer` when you need `[]byte` output or implement `io.Writer`.
- **Zero-copy string to bytes**: In Go 1.20+, `unsafe.String(ptr, len)` and `unsafe.Slice(ptr, len)` allow zero-copy conversion but require care. Never modify the result.
- **`bytes.Cut` (Go 1.18+)**: Efficiently splits before and after a separator: `before, after, found := bytes.Cut(data, sep)`.

## Performance implications

- `[]byte(s)` allocates. For large strings, this dominates CPU profiles. Avoid in hot loops.
- `bytes.Buffer.Write` grows using the same algorithm as slice append: amortized O(1).
- `bytes.Contains` and `bytes.Index` are highly optimized (assembly on most platforms).
- `bytes.Replace` allocates a new slice. In-place replacement isn't possible because the result may be longer or shorter.
- Zero-copy conversion with `unsafe` is ~10x faster but sacrifices safety. Profile before using it.

## Practice task

Write a function `wordCount(text []byte) map[string]int` that splits `text` by whitespace and returns a map of word → count (case-sensitive). Then write a function `toUpperASCII(b []byte) []byte` that returns a new byte slice where each lowercase ASCII letter (a-z) is converted to uppercase. In `main()`, test on `[]byte("hello world hello Go")` and print results.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/20-bytes
go test ./curriculum/modules/03-programming-fundamentals/lessons/20-bytes
```

## Review questions

1. What is the difference between `byte` and `uint8` in Go?
2. Does `string([]byte{104, 105})` allocate? Why?
3. What does `bytes.Buffer.Bytes()` return, and why must you copy it?
4. How does `bytes.Split` differ from `bytes.SplitN`?
5. Can you use `[]byte` as a map key? How would you work around this?

## NEXT UP

Runes and Unicode — working with individual Unicode code points and Go's UTF-8 support.
