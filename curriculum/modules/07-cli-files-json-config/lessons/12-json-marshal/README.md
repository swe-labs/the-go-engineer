# JSON marshal

## Learning objective

Serialize Go structs, maps, and slices into JSON byte slices using `json.Marshal` and `json.MarshalIndent`, control field names and inclusion with struct tags, and understand why only exported fields appear in JSON output.

## Why this matters

JSON is the universal wire format for web APIs, configuration files, and data interchange. Every Go service you write will marshal JSON — to send HTTP responses, write config files, or persist data. Misunderstanding marshaling leads to missing fields, silent errors, and serialization bugs that are hard to trace. Mastering `encoding/json` is non-negotiable.

## Mental model

Think of JSON marshaling as a camera that photographs Go values. The camera sees only exported fields (capital letter = visible). Struct tags are like sticky notes you put on each field telling the camera what name to use and how to treat special cases like empty values or string conversion. `json.Marshal` takes the picture into memory; `json.MarshalIndent` takes the same picture but arranges it neatly for human viewing.

## Core idea

`json.Marshal(v any) ([]byte, error)` converts a Go value into a JSON byte slice using reflection. It traverses the value at runtime, reads struct tags, and produces the corresponding JSON text.

Key rules:
- Only exported struct fields are marshaled (package-level visibility check via reflection).
- Each field's JSON name defaults to the field name; override with `json:"name"` tag.
- The tag option `omitempty` omits the field if it has a zero value.
- The tag option `string` forces the value to be quoted as a JSON string.
- `json.MarshalIndent` adds newlines and indentation for readability.

## Under the hood

`encoding/json` uses the `reflect` package to inspect values at runtime. For a struct, `reflect.Type.NumField()` and `Type.Field(i)` retrieve each field's name, type, and tag string. The encoder then dispatches to type-specific encoding functions (`encodeInt`, `encodeString`, `encodeSlice`, etc.).

The encoder writes into a `bytes.Buffer` internally. Each value is encoded recursively: slices loop over elements, structs loop over fields, maps iterate over key-value pairs. The zero-alloc path uses `sync.Pool` for reusing encoder buffers.

`time.Time` implements `json.Marshaler` — its `MarshalJSON()` method outputs RFC 3339 format (`"2006-01-02T15:04:05Z07:00"`).

## How Go uses it

- HTTP response bodies: `json.NewEncoder(w).Encode(data)`
- Configuration files: write JSON configs with `json.MarshalIndent`
- Database JSON columns: marshal structs to `[]byte` for storage
- Logging: serialize structured log fields
- Testing: generate expected JSON output for golden file tests

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Order struct {
	ID        int       `json:"id"`
	Product   string    `json:"product"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	Notes     string    `json:"notes,omitempty"`
	Secret    string    `json:"-"`
}

func main() {
	o := Order{
		ID:        1001,
		Product:   "Widget",
		Quantity:  3,
		Price:     14.99,
		CreatedAt: time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC),
		Notes:     "",
		Secret:    "admin-password",
	}

	compact, _ := json.Marshal(o)
	fmt.Println("Compact:", string(compact))

	pretty, _ := json.MarshalIndent(o, "", "  ")
	fmt.Println("Pretty:")
	fmt.Println(string(pretty))

	// MarshalIndent with prefix
	prefixed, _ := json.MarshalIndent(o, "> ", "  ")
	fmt.Println("Prefixed:")
	fmt.Println(string(prefixed))
}
```

## Step-by-step execution

1. `Order` struct is created with field values.
2. `json.Marshal(o)` is called. The encoder uses reflection to discover `Order` has 7 fields.
3. Field `ID` has tag `json:"id"` → JSON key `"id"`, value `1001` → JSON number.
4. Field `Notes` has tag `json:"notes,omitempty"` — value is `""` (zero string) → field omitted.
5. Field `Secret` has tag `json:"-"` — field is always omitted.
6. `CreatedAt` is `time.Time` — custom `MarshalJSON()` produces `"2025-06-01T10:00:00Z"`.
7. Compact output: single line with no extra whitespace.
8. `MarshalIndent` adds `"\n"` after `{`, `,`, and `}` plus `"  "` indentation per nesting level.
9. Prefix `"> "` is prepended to each new line inside the JSON body.

## Common mistakes

- Mistake: Forgetting to export fields — `product string` (lowercase) won't marshal.
  - Fix: Capitalize field names: `Product string`.

- Mistake: Using `json:"-"` accidentally — the field is silently dropped.
  - Fix: Remove the `"-"` tag or use a different tag name.

- Mistake: Expecting `time.Time` to output a custom format other than RFC 3339.
  - Fix: Implement `json.Marshaler` on a custom type, or use a wrapper.

- Mistake: Assuming `omitempty` omits `false` booleans or `0` integers.
  - Fix: It does omit them — they are zero values. Use a pointer if you need to distinguish "not set" from "zero".

## Debugging walkthrough

Consider this code that produces unexpected output:

```go
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	role string
}

u := User{Name: "Alice", Age: 30, role: "admin"}
data, _ := json.Marshal(u)
fmt.Println(string(data)) // {"name":"Alice","age":30}
```

**Symptom**: The `role` field is missing from JSON.

**Investigation**: Check field export status — `role` starts with lowercase `r`. Go's `json.Marshal` only sees exported fields.

**Root cause**: The `role` field is unexported (private). Reflection cannot access it from another package (`encoding/json`).

**Fix**: Capitalize to `Role string \`json:"role"\``.

## Production notes

- **Always handle marshal errors**: disk-full, nil pointer in a struct field, or cyclic data structures cause errors. Never swallow them with `_`.
- **Use `MarshalIndent` only for human-readable output** (config files, logs). For APIs, use compact `Marshal` to save bandwidth.
- **Omit sensitive fields** with `json:"-"` — never rely on it for security (other encodings may still expose it), but it prevents accidental leakage in JSON paths.
- **Pre-allocate buffers** with `json.NewEncoder(buf)` when marshaling many values in a loop — reuses encoder state.

## Performance implications

- `json.Marshal` uses reflection, which is slower than hand-written serialization but much safer.
- For high-throughput JSON (>10K objects/sec), consider `easyjson` or `ffjson` which generate static marshal code.
- `MarshalIndent` is ~20-30% slower than `Marshal` due to extra formatting work.
- Zero values with `omitempty` reduce output size but add tag parsing overhead at init time.
- `map[string]any` marshaling is slower than struct marshaling because of type assertions.

## Practice task

Write a function `MarshalUser(name string, age int, email string) (string, error)` that returns a pretty-printed JSON string of a struct with fields `name`, `age`, and `email` (using struct tags for lowercase JSON keys, `omitempty` on email). If `age` is negative, return an error. Then call it from `main` with three different inputs and print the results.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/12-json-marshal
go test ./curriculum/modules/07-cli-files-json-config/lessons/12-json-marshal
```

## Review questions

1. What happens if a struct field has a lowercase name? Will it appear in `json.Marshal` output?
2. What does the struct tag `json:"-"` do? What about `json:"name,omitempty"`?
3. How does `time.Time` get marshaled — what format does it use?
4. Why might `MarshalIndent` be the wrong choice for an HTTP API response?
5. Given `json:"total,string"`, what type will the JSON value `"42"` be in Go?

## NEXT UP

JSON unmarshal — the reverse operation: parsing JSON byte slices back into Go values.
