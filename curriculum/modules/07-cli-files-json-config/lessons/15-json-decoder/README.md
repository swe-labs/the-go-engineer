# JSON decoder

## Learning objective

Stream JSON input from any `io.Reader` using `json.NewDecoder`, decode tokens with `Token` and `More`, and handle large JSON payloads without loading the entire blob into memory.

## Why this matters

Real-world JSON payloads can be gigabytes — logs, data exports, API responses. Loading the entire payload into memory with `json.Unmarshal` is impractical or impossible. `json.NewDecoder` reads tokens one at a time from a stream, consuming a fixed amount of memory regardless of input size. Every professional Go engineer needs streaming JSON for high-throughput and large-scale systems.

## Mental model

`json.Unmarshal` is like photocopying an entire book, then reading it. `json.NewDecoder` is like reading the book page by page as it feeds through a scanner — at any moment, only the current page is in memory. The decoder maintains a cursor into the input stream and produces tokens (`Delim`, `String`, `Number`, `Bool`, `nil`) one at a time, or decodes directly into a Go value with `Decode`.

## Core idea

`json.NewDecoder(r io.Reader) *json.Decoder` creates a decoder that reads from `r`. It provides two modes:

1. **Value decoding**: `dec.Decode(v any) error` — reads the next complete JSON value from the stream and stores it in `v`. Works for any JSON value (object, array, string, number, bool, null).

2. **Token streaming**: `dec.Token() (json.Token, error)` — reads one JSON token at a time (`json.Delim` for `{}[]`, `string`, `float64`, `bool`, `nil`). Combined with `dec.More()` to check if there are more elements in an array or object.

Additional useful methods:
- `dec.Buffered() io.Reader` — returns remaining unread data.
- `dec.DisallowUnknownFields()` — returns error on unknown JSON keys.
- `dec.UseNumber()` — stores numbers as `json.Number` (string) instead of `float64`.

## Under the hood

The decoder reads from the reader in chunks (default 512 bytes). It tokenizes JSON incrementally using a state machine:

```
'{' → expect object start → look for string keys or '}'
'"key"' → read string → expect ':'
Value → read value (recursively for objects/arrays)
',' or '}' → continue or close
```

When `Decode` is called, it consumes one complete JSON value (which may contain nested objects/arrays). The stream position advances past the consumed value, so subsequent `Decode` calls read the next JSON value — enabling NDJSON parsing.

## How Go uses it

- HTTP request body parsing: `json.NewDecoder(r.Body).Decode(&data)`
- NDJSON (newline-delimited JSON) log processing
- Streaming large JSON arrays with `Token` and `More`
- API gateway response streaming
- JSON data migration pipelines

## Go example

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func main() {
	// NDJSON: multiple JSON objects, one per line
	ndjson := `{"id":1,"name":"Alice"}
{"id":2,"name":"Bob"}
{"id":3,"name":"Carol"}`

	dec := json.NewDecoder(strings.NewReader(ndjson))
	type Person struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	for {
		var p Person
		err := dec.Decode(&p)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Decode error:", err)
			break
		}
		fmt.Printf("Person: ID=%d, Name=%s\n", p.ID, p.Name)
	}

	// Token-based streaming over a JSON array
	fmt.Println("\nToken-based array parsing:")
	array := `["go", "rust", "zig"]`
	dec2 := json.NewDecoder(bytes.NewReader([]byte(array)))
	for {
		t, err := dec2.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Token error:", err)
			break
		}
		if delim, ok := t.(json.Delim); ok {
			fmt.Printf("Delimiter: %c\n", delim)
		} else {
			fmt.Printf("Token: %v (type: %T)\n", t, t)
		}
	}
}
```

## Step-by-step execution

For `dec.Decode(&p)` on first line `{"id":1,"name":"Alice"}`:

1. Decoder reads bytes from the reader into internal buffer.
2. Scans `{` — object start. Begins object parsing.
3. Reads key `"id"` → expects `:` → reads value `1` → stores as `int` in `p.ID`.
4. Reads `,` — more keys. Reads `"name"` → `:` → `"Alice"` → stores in `p.Name`.
5. Reads `}` — object end.
6. Returns nil. Stream now points to `\n` (next line).
7. Next `Decode` call reads the next line's JSON object.
8. When no more JSON values remain, returns `io.EOF`.

For token-based parsing of `["go", "rust", "zig"]`:
1. `Token()` returns `[` (json.Delim).
2. `Token()` returns `"go"`.
3. `Token()` returns `"rust"`.
4. `Token()` returns `"zig"`.
5. `Token()` returns `]` (json.Delim).
6. `Token()` returns `io.EOF`.

## Common mistakes

- Mistake: Calling `Decode` after `Token` without tracking position — `Decode` expects a complete JSON value starting from the current position.
  - Fix: Use one approach consistently, or carefully track the stream state.

- Mistake: Ignoring `io.EOF` vs other errors — `io.EOF` means clean end of input; other errors mean malformed JSON.
  - Fix: Check `err == io.EOF` for loop termination, report other errors.

- Mistake: Using `More()` at the wrong nesting level — `More()` tells if there are more items in the current array/object being tokenized.
  - Fix: Only call `More()` between tokens when reading inside a container.

## Debugging walkthrough

This code loops forever:

```go
data := `[1, 2, 3]`
dec := json.NewDecoder(strings.NewReader(data))
for {
	var v int
	err := dec.Decode(&v)
	if err != nil {
		break
	}
	fmt.Println(v)
}
```

**Symptom**: Only prints `1` (or panics) or loops infinitely.

**Root cause**: `Decode` reads one complete JSON value. The first call reads the entire array `[1,2,3]` and tries to store it in `int` — this fails because a JSON array can't go into an int. The error breaks the loop.

**Fix**: Decode into `[]int` or use `Token` to iterate array elements:

```go
var values []int
json.NewDecoder(strings.NewReader(data)).Decode(&values)
```

## Production notes

- **Always check for `io.EOF`** to distinguish clean completion from errors.
- **Use `DisallowUnknownFields()`** in strict API servers to reject unexpected fields.
- **Use `UseNumber()`** when you need exact precision for large integers (IDs, financial amounts) — avoids float64 rounding.
- **Close the underlying reader** when done (e.g., `r.Body.Close()` in HTTP handlers).
- **For huge JSON arrays** (millions of elements), use `Token` + `More` to iterate without allocating a giant slice.

## Performance implications

- `json.NewDecoder` uses a small internal buffer (default growing as needed). It reads more efficiently than repeatedly calling `Read` on the underlying reader.
- Token-based parsing is more CPU-efficient than `Decode` when you only need specific fields — you can skip unwanted tokens without allocating Go objects for them.
- `Decode` allocates memory for the decoded Go value. For high-throughput, reuse structs by resetting them between decodes.
- Using `UseNumber()` adds a string allocation per number but avoids float64 precision loss.

## Practice task

Write a function `SumJSON(r io.Reader) (int, error)` that reads a stream of JSON numbers (one per line, NDJSON) and returns their sum. Each line is a single JSON number (e.g., `42\n`, `100\n`). Stop on `io.EOF`. Return an error if any line is not a valid JSON number.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/15-json-decoder
go test ./curriculum/modules/07-cli-files-json-config/lessons/15-json-decoder
```

## Review questions

1. What does `json.NewDecoder` return when there are no more JSON values in the stream?
2. How does `dec.Token()` differ from `dec.Decode(&v)`?
3. What does `dec.More()` tell you? When would you use it?
4. What problem does `dec.UseNumber()` solve?
5. Why is streaming JSON decoding important for large payloads?

## NEXT UP

Config files — patterns for loading application configuration from JSON and YAML files.
