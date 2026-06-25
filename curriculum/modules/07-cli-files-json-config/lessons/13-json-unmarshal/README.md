# JSON unmarshal

## Learning objective

Parse JSON byte slices into Go structs, maps, and slices using `json.Unmarshal`, handle unknown fields and type mismatches, and use `json.RawMessage` for deferred decoding.

## Why this matters

Reading JSON is the mirror of writing it — every API client, config loader, and data pipeline needs unmarshaling. But unmarshaling is more error-prone than marshaling: incoming JSON may have unexpected fields, wrong types, or missing values. Understanding how Go maps JSON to Go types prevents runtime panics, data corruption, and silent failures.

## Mental model

Unmarshaling is like filling a prefabricated form. You have a Go struct (the form with labeled slots) and a JSON blob (the data to fill in). Go reads the JSON key by key, matches each key to a struct field tag (or name), converts the JSON value to the Go type, and stores it. If a key has no matching slot, it's discarded by default. If a value can't be converted to the slot's type, the whole operation fails.

## Core idea

`json.Unmarshal(data []byte, v any) error` parses JSON `data` and stores the result in the value pointed to by `v`. The target must be a pointer (to a struct, slice, map, or concrete type).

JSON-to-Go type mapping:

| JSON type | Go target type |
|-----------|---------------|
| string | `string` |
| number | `float64` (or `int`/`int64` if tagged appropriately) |
| boolean | `bool` |
| null | nil pointer, or zero value |
| object | `struct`, `map[string]any` |
| array | `[]any`, typed slice |

## Under the hood

`json.Unmarshal` works through a state machine that tokenizes JSON bytes. It reads one token at a time (`{`, `}`, `[`, `]`, `"string"`, `123`, `true`, etc.), matches it against the target Go type using reflection, and writes the decoded value.

The decoder handles:
- **Disappearing fields**: JSON keys without matching struct fields are silently ignored (unless `DisallowUnknownFields` is set on a decoder).
- **Null**: JSON `null` does nothing to a pointer field (leaves it nil) or zeroes a non-pointer field.
- **Type coercion**: JSON numbers can be stored in `int`, `float64`, or `string` fields (with tags like `json:"fieldname,string"`).
- **Time parsing**: `time.Time` fields with the `json:"field"` tag parse RFC 3339 strings automatically.

## How Go uses it

- HTTP request body parsing: `json.NewDecoder(r.Body).Decode(&data)`
- Reading config files: `json.Unmarshal(fileBytes, &config)`
- API client response parsing
- Database JSON column scanning
- Testing: parse expected output for comparison

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Tags      []string  `json:"tags,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
}

func main() {
	input := `{
		"id": 42,
		"name": "deploy",
		"timestamp": "2025-06-01T10:00:00Z",
		"tags": ["prod", "us-east"],
		"metadata": {"env": "production", "version": "2.1"}
	}`

	var e Event
	if err := json.Unmarshal([]byte(input), &e); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("ID: %d\n", e.ID)
	fmt.Printf("Name: %s\n", e.Name)
	fmt.Printf("Timestamp: %s\n", e.Timestamp.Format(time.RFC3339))
	fmt.Printf("Tags: %v\n", e.Tags)
	fmt.Printf("Metadata (raw): %s\n", string(e.Metadata))

	// json.RawMessage can be unmarshaled later
	var meta map[string]string
	json.Unmarshal(e.Metadata, &meta)
	fmt.Printf("Parsed metadata: %v\n", meta)
}
```

## Step-by-step execution

1. `json.Unmarshal` receives the JSON string and a pointer to `Event`.
2. It reads `{` — expects an object. Target is a struct, so it proceeds.
3. Key `"id"` → matches `ID` field (tag `json:"id"`). Value `42` is a JSON number → stored as `int` in `e.ID`.
4. Key `"name"` → matches `Name`. Value `"deploy"` is a JSON string → stored as `string`.
5. Key `"timestamp"` → matches `Timestamp`. Value `"2025-06-01T10:00:00Z"` is a JSON string. Go checks if `Timestamp` type (`time.Time`) implements `json.Unmarshaler` — it does. Calls `UnmarshalJSON` which parses RFC 3339.
6. Key `"tags"` → matches `Tags`. Value is a JSON array → each element decoded as `string` and appended to slice.
7. Key `"metadata"` → matches `Metadata` (type `json.RawMessage`). Value is stored as raw bytes without parsing.
8. Unknown key `"extra"` (if present) would be silently discarded.
9. All keys consumed, `}` closes the object. Unmarshal returns nil.

## Common mistakes

- Mistake: Passing a non-pointer to `Unmarshal` — `json.Unmarshal(data, v)` where `v` is a struct value, not a pointer.
  - Fix: Always pass a pointer: `json.Unmarshal(data, &v)`.

- Mistake: Expecting JSON numbers to unmarshal into `int` fields by default — they become `float64` in `any` targets.
  - Fix: Use `json.Number` or a typed struct field with `json.Decoder` (UseNumber).

- Mistake: Forgetting `time.Time` expects RFC 3339 — other formats cause parse errors.
  - Fix: Use a custom type with `UnmarshalJSON` for non-standard formats.

- Mistake: Assuming unknown fields cause an error — they are silently ignored.
  - Fix: Use `json.NewDecoder` with `DisallowUnknownFields()`.

## Debugging walkthrough

This code fails silently:

```go
type Config struct {
	Port int `json:"port"`
}
data := []byte(`{"port": "8080"}`)
var cfg Config
err := json.Unmarshal(data, &cfg)
fmt.Println(cfg.Port, err) // 0 json: cannot unmarshal string into Go struct field Config.port of type int
```

**Symptom**: `Port` is 0, error is returned.

**Root cause**: JSON value `"8080"` is a string (note the quotes), but Go expects a number for `int`.

**Fix**: Change the JSON to `{"port": 8080}` (unquoted), or accept the error and handle it. To accept string-encoded numbers, use `json:"port,string"` tag.

## Production notes

- **Always check the error** from `Unmarshal`. A single type mismatch fails the entire parse.
- **Use `json.RawMessage`** for fields with dynamic schemas — decode the known fields first, then inspect and decode the raw payload based on a discriminator field.
- **Validate after unmarshal**: JSON may be syntactically valid but semantically wrong (e.g., negative port). Always validate business logic after parsing.
- **Prefer `json.Decoder`** for streaming large inputs (files, HTTP bodies) — it doesn't need the entire payload in memory.

## Performance implications

- `json.Unmarshal` allocates memory for each decoded value. Unmarshaling into `map[string]any` allocates more than structs because of interface boxing.
- Using `json.RawMessage` avoids parsing a subtree until needed — saves CPU if the subtree is conditionally used.
- `json.Number` avoids float64 precision loss but adds string-to-number conversion cost.
- For large JSON (>1 MB), prefer `json.NewDecoder` which streams tokens instead of loading the entire blob.

## Practice task

Write a function `UnmarshalEvent(data []byte) (Event, error)` that parses a JSON event with fields `id` (int), `title` (string), `timestamp` (RFC 3339 string), and an optional `payload` (json.RawMessage). Test it with valid JSON, JSON with wrong types, and JSON with extra unknown fields. Print results from `main`.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/13-json-unmarshal
go test ./curriculum/modules/07-cli-files-json-config/lessons/13-json-unmarshal
```

## Review questions

1. What happens if you pass a struct value (not a pointer) to `json.Unmarshal`?
2. What Go type does a JSON number unmarshal into when the target is `any`?
3. How does `time.Time` know how to parse the timestamp string from JSON?
4. When would you use `json.RawMessage` instead of a concrete struct type?
5. What does `json.Unmarshal` do with a JSON key that has no matching struct field?

## NEXT UP

JSON encoder — streaming JSON output to writers instead of building byte slices in memory.
