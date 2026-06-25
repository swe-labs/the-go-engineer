# JSON encoder

## Learning objective

Stream JSON output directly to any `io.Writer` using `json.NewEncoder`, control formatting with `SetIndent` and `SetEscapeHTML`, and understand when to use the encoder over `json.Marshal`.

## Why this matters

Every HTTP response, log line, and file write that produces JSON is an encoding operation. Using `json.Marshal` to build a `[]byte` and then writing it to a writer duplicates memory and forces allocation of the full output before sending a single byte. `json.NewEncoder` streams JSON directly to the destination, reducing latency and memory pressure — critical for large responses or high-throughput servers.

## Mental model

Think of `json.Marshal` as taking a photo and printing it; the photo exists fully in your hand before you hand it over. `json.NewEncoder` is a live TV feed — the encoder starts sending bytes to the writer immediately as it encodes, without storing the complete picture in memory. The encoder wraps a writer and writes JSON tokens (`{`, `"key"`, `:`, `value`, `,`, `}`) directly as it traverses the Go value.

## Core idea

`json.NewEncoder(w io.Writer) *json.Encoder` creates an encoder that writes JSON to `w`. Call `enc.Encode(v any) error` to encode a single value.

Key differences from `json.Marshal`:
- Writes directly to `w` instead of returning `[]byte`.
- Appends a trailing newline (`\n`) after each encoded value — making it suitable for newline-delimited JSON (NDJSON).
- Supports `SetIndent(prefix, indent string)` for pretty printing (equivalent to `MarshalIndent`).
- Supports `SetEscapeHTML(on bool)` to disable HTML escaping of `<`, `>`, `&` (default: on).

## Under the hood

`json.Encoder` contains an internal `bytes.Buffer`-like state and a reference to the target writer. On `Encode`, it:
1. Resets internal state.
2. Calls the same encoding functions as `json.Marshal`, but writes tokens to the buffer then flushes to `w`.
3. Appends `\n` after the complete JSON value.
4. Returns any write error from the underlying writer.

The encoder is not safe for concurrent use — if multiple goroutines call `Encode` simultaneously on the same encoder, output will be interleaved.

## How Go uses it

- HTTP response encoding: `json.NewEncoder(w).Encode(response)`
- Logging structured data: `json.NewEncoder(os.Stdout).Encode(logEntry)`
- Writing NDJSON (newline-delimited JSON) files for data pipelines
- Writing JSON config files with pretty printing (`SetIndent`)
- Streaming large datasets to a file without buffering the entire set in memory

## Go example

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Product struct {
	SKU   string  `json:"sku"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {
	products := []Product{
		{"A100", "Widget", 9.99},
		{"A200", "Gadget", 24.99},
		{"A300", "Doohickey", 4.99},
	}

	// Encode to os.Stdout with indentation
	fmt.Println("=== Stdout with indentation ===")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	for _, p := range products {
		if err := enc.Encode(p); err != nil {
			fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		}
	}

	// Encode to a buffer with SetEscapeHTML disabled
	fmt.Println("\n=== Buffer (no HTML escaping) ===")
	var buf bytes.Buffer
	enc2 := json.NewEncoder(&buf)
	enc2.SetEscapeHTML(false)
	enc2.Encode(Product{"B500", "Chip < v2 > & Co", 149.99})
	fmt.Print(buf.String())

	// NDJSON example
	fmt.Println("=== NDJSON (compact, newline-delimited) ===")
	var ndjson bytes.Buffer
	enc3 := json.NewEncoder(&ndjson)
	// default: no indent, HTML escaping on
	for _, p := range products {
		enc3.Encode(p)
	}
	lines := strings.Split(strings.TrimSpace(ndjson.String()), "\n")
	fmt.Printf("Wrote %d lines\n", len(lines))
	for _, line := range lines {
		fmt.Println("  " + line)
	}
}
```

## Step-by-step execution

1. `json.NewEncoder(os.Stdout)` creates an encoder writing to stdout.
2. `SetIndent("", "  ")` configures pretty printing — stored in encoder state.
3. `enc.Encode(p)` is called for `Product{"A100", "Widget", 9.99}`:
   - Encoder writes `{` then `\n` to internal buffer.
   - Writes `  "sku": "A100"`, `  "name": "Widget"`, `  "price": 9.99`.
   - Writes `\n}` followed by `\n`.
   - Flushes buffer to `os.Stdout`.
4. Each subsequent `Encode` call repeats — output is independent JSON objects separated by newlines.
5. `SetEscapeHTML(false)` on the second encoder means `<`, `>`, `&` in strings are not escaped to `\u003c`, `\u003e`, `\u0026`.

## Common mistakes

- Mistake: Using the same encoder concurrently from multiple goroutines — output interleaves.
  - Fix: Use one encoder per goroutine, or protect with a mutex.

- Mistake: Forgetting `Encode` appends `\n`. Concatenating encoded values produces valid NDJSON but is not a valid single JSON array.
  - Fix: If you need a JSON array, encode a slice: `enc.Encode(allItems)`.

- Mistake: Using `SetIndent` in high-throughput HTTP handlers — it adds overhead for every response.
  - Fix: Only indent for human-readable endpoints (e.g., admin debug pages).

- Mistake: Assuming `Encode` returns before all bytes are written — for buffered writers, call `w.Flush()` if needed.

## Debugging walkthrough

This code writes nothing visible:

```go
var buf bytes.Buffer
enc := json.NewEncoder(&buf)
enc.Encode("hello")
fmt.Println("Buffer is empty?")
```

**Symptom**: The buffer seems empty when printed later (if printed before encode), or looks correct after encode. Actually `Encode` appends `\n`, so `buf.String()` will be `"hello\n"`.

**Investigation**: Print `buf.String()` after `Encode`.

**Root cause**: The encoder does write to the buffer, but the newline may be unexpected. If you trim it: `strings.TrimSpace(buf.String())`.

## Production notes

- **Use `SetEscapeHTML(false)`** when JSON values don't contain HTML (most API responses) — it's slightly faster and produces more readable output.
- **For HTTP responses**, use `json.NewEncoder(w).Encode(data)` directly without buffering into a `bytes.Buffer` first — this reduces memory allocation and latency.
- **Always check `Encode` errors** — a broken network connection or full disk causes write errors that you must handle.
- **NDJSON** is the standard format for streaming logs, events, and data pipeline records. Each line is a complete JSON object.

## Performance implications

- `json.NewEncoder` allocates less memory than `json.Marshal` + `w.Write()` because it avoids the intermediate `[]byte`.
- For small payloads (< 1 KB), the difference is negligible. For large payloads (> 100 KB), the encoder can be 2-3x more memory efficient.
- `SetIndent` adds ~20% CPU overhead.
- The encoder's internal buffer grows as needed — for repeated encodes of similar-sized objects, reuse the encoder and call `buf.Reset()` between uses.

## Practice task

Write a function `WriteProducts(products []Product, w io.Writer) error` that encodes each product as a separate JSON line (NDJSON) to the writer with indentation disabled. If the writer returns an error, stop immediately and return it. Test by writing to a `bytes.Buffer` and verifying the output has one JSON object per line.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/14-json-encoder
go test ./curriculum/modules/07-cli-files-json-config/lessons/14-json-encoder
```

## Review questions

1. How does `json.NewEncoder` differ from `json.Marshal` in terms of memory and output?
2. What character does `Encode` append after every JSON value? Why?
3. What does `SetEscapeHTML(false)` change in the output?
4. Can you safely call `Encode` on the same encoder from multiple goroutines?
5. When would you choose `json.MarshalIndent` over `json.Encoder` with `SetIndent`?

## NEXT UP

JSON decoder — streaming JSON input from readers instead of parsing byte slices in memory.
