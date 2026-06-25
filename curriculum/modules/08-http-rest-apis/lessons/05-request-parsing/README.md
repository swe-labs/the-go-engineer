# Request parsing

## Learning objective

Extract data from HTTP requests using query strings (`r.URL.Query`), path values (`r.PathValue`), JSON bodies (`json.NewDecoder(r.Body)`), form data, and headers (`r.Header`).

## Why this matters

Every HTTP handler needs to read input. The input can come from the URL path, query parameters, headers, or the request body. Parsing these correctly is the difference between a robust API that gracefully handles edge cases and one that panics on missing data or crashes on malformed input. Go's standard library provides direct, zero-allocation access to each of these sources — but each has its own gotchas around encoding, escaping, and buffering.

## Mental model

Imagine a request as a filing envelope. On the outside (the URL), you have the delivery address (path) and routing notes (query string). On the flap, you have stamps and labels (headers). Inside (the body), you have the actual document. Go lets you open the envelope and read each part separately. The path tells you *what* resource, the query string tells you *how* (filter, sort, paginate), headers tell you *who* and *what format*, and the body tells you *what to change*.

## Core idea

There are four main places to read input from an HTTP request:

| Source | Access method | Use case |
|--------|--------------|----------|
| Path value | `r.PathValue("key")` | Resource identifier (`/users/{id}`) |
| Query string | `r.URL.Query().Get("key")` | Filters, pagination, options |
| Headers | `r.Header.Get("Key")` | Auth tokens, content negotiation, metadata |
| Body (JSON) | `json.NewDecoder(r.Body)` | Structured data for create/update |
| Body (form) | `r.ParseForm()` / `r.FormValue("key")` | HTML form submissions |

Key rules:
- `r.Body` is an `io.ReadCloser`. It must be closed (Go server does this automatically) and can only be read once. To read it twice, buffer it first.
- `r.URL.Query()` parses the query string into a `url.Values` (a `map[string][]string`). `Get` returns the first value.
- `r.FormValue` calls `r.ParseForm` internally if not already called, but only reads URL query and form body — not JSON.
- `r.PathValue` is only available in Go 1.22+.

## Under the hood

`json.NewDecoder(r.Body)` creates a streaming JSON decoder. It reads from `r.Body` incrementally, which means it can process very large payloads without buffering the entire body in memory. The decoder stops after reading one JSON value — subsequent data in the body is ignored. For forms, `r.ParseForm` reads the entire body into memory and populates `r.Form`. The body is consumed and cannot be read again unless buffered. `r.URL.Query` parses the raw query string from `r.URL.RawQuery` using `url.ParseQuery`, which splits on `&` and `=` and decodes percent-encoded characters.

## How Go uses it

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Path value (Go 1.22+)
    id := r.PathValue("id")

    // Query string
    name := r.URL.Query().Get("name")
    page := r.URL.Query().Get("page")

    // Headers
    auth := r.Header.Get("Authorization")

    // JSON body
    var input CreateUserInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "bad request", 400)
        return
    }
}
```

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// SearchParams holds parsed request data.
type SearchParams struct {
	Query string `json:"query"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query string.
	params := SearchParams{
		Query: r.URL.Query().Get("q"),
		Limit: 20,
	}
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			params.Page = v
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			params.Limit = v
		}
	}

	// Read auth header.
	auth := r.Header.Get("Authorization")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"params": params,
		"auth":   auth[:min(len(auth), 8)] + "...",
	})
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /search", searchHandler)
	mux.HandleFunc("POST /data", createHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

1. Client sends `GET /search?q=golang&page=1&limit=10` with header `Authorization: Bearer tok123`.
2. `searchHandler` calls `r.URL.Query().Get("q")` → `"golang"`.
3. `r.URL.Query().Get("page")` → `"1"`, parsed to `int 1` via `strconv.Atoi`.
4. `r.URL.Query().Get("limit")` → `"10"`, parsed to `int 10`.
5. `r.Header.Get("Authorization")` → `"Bearer tok123"`.
6. Response encoded as JSON with the parsed params and truncated auth header.
7. For `POST /data`, `json.NewDecoder(r.Body)` reads the stream, decodes into a `map[string]interface{}`.
8. The decoded map is sent back as JSON with 201 Created.

## Common mistakes

- Calling `json.NewDecoder(r.Body).Decode` multiple times. The body is a stream; after the first read, it's exhausted. Save the decoded value.
- Not handling `strconv.Atoi` errors. Query string values are strings; converting to int can fail. Always check the error and fall back to a default.
- Using `r.FormValue` expecting JSON to be parsed. `r.FormValue` only parses URL query and form-encoded bodies (or multipart forms). For JSON, use `json.Decode`.
- Reading `r.Body` directly without `json.Decode` and forgetting to close it. While Go closes it after the handler returns, a large unread body wastes memory.
- Assuming query parameter order. `r.URL.Query()` returns a `map[string][]string`; the order of values for the same key is preserved (first occurrence first), but iteration order over keys is random.

## Debugging walkthrough

A handler always sees empty path values even though the route is registered with `{id}`:

```go
mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    fmt.Println("id:", r.PathValue("id")) // always ""
})
```

**Symptom**: `r.PathValue("id")` returns empty string.

**Root cause**: Go version < 1.22. `r.PathValue` was introduced in Go 1.22. On older versions, the method exists (it was added to `*http.Request` in 1.22) but returns empty string because `ServeMux` doesn't set path values.

**Fix**: Update `go.mod` to Go 1.22+: `go mod edit -go=1.22`. Or use a third-party router that supports path parameters.

## Production notes

- Use `json.Decoder` instead of `io.ReadAll` + `json.Unmarshal` for large bodies. The streaming decoder avoids allocating a buffer for the entire body.
- Validate content type before parsing. Check `r.Header.Get("Content-Type")` before decoding JSON or form data. Return 415 Unsupported Media Type if the content type is wrong.
- Limit body size with `http.MaxBytesReader` to prevent memory exhaustion attacks.
- Consider using `r.Clone` or buffering `r.Body` with `io.NopCloser(bytes.NewBuffer(buf))` if multiple readers are needed (e.g., for logging).

## Performance implications

- `json.NewDecoder(r.Body).Decode` is streaming — it reads from the TCP connection directly without buffering the entire body. This reduces memory allocation for large payloads.
- `r.URL.Query()` parses the query string on every call. Cache the result if accessed multiple times.
- `r.FormValue` calls `r.ParseForm` on first invocation, which reads and parses the entire body for form data. This is a one-time cost.
- `r.PathValue` is O(1) — it reads from a pre-populated map inside the `Request` struct.

## Practice task

Write a handler `processHandler` that:
- Reads a path value `{id}` → required.
- Reads a query parameter `format` → optional, defaults to `"json"`.
- Reads header `X-Trace-ID` → optional, for logging.
- Reads a JSON body with fields `name` (string) and `count` (int) → both required.
- Returns a 400 error as JSON with field-level messages if any required data is missing or invalid.
- Returns 200 with the parsed data as JSON on success.

## Tests / verification

```bash
go run ./curriculum/modules/08-http-rest-apis/lessons/05-request-parsing
go test ./curriculum/modules/08-http-rest-apis/lessons/05-request-parsing
```

## Review questions

1. Why can you only read `r.Body` once? How would you read it twice?
2. What is the difference between `r.FormValue("key")` and `r.URL.Query().Get("key")`?
3. How do you ensure a JSON decoder doesn't consume more than one JSON value?
4. What happens if you call `r.PathValue("id")` on a request that was routed to a handler without a `{id}` wildcard?
5. How would you parse a multipart form upload in Go?

## NEXT UP

Input validation — validating parsed input for required fields, type correctness, and range constraints using manual checks and `go-playground/validator`.
