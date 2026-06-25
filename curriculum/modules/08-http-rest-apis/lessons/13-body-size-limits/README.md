# Body size limits

## Learning objective

Limit request body sizes using `http.MaxBytesReader`, configure per-handler and global limits, and handle oversized requests with appropriate error responses to protect against DoS attacks.

## Why this matters

Without body size limits, an attacker can send a multi-gigabyte POST to your server, exhausting memory and crashing the process. Even non-malicious clients can accidentally send oversized payloads. Every production HTTP server must enforce request body limits as a basic security and reliability measure.

## Mental model

Think of `http.MaxBytesReader` as a fuse on an electrical circuit. The wire (your handler) can only carry so much current (data). The fuse limits the flow. When the current exceeds the rating, the fuse blows -- in Go terms, the reader returns an error and stops reading. Your handler checks for that error and returns `413 Request Entity Too Large` instead of trying to process the impossible.

## Core idea

`http.MaxBytesReader` wraps an `io.ReadCloser` (the request body) with a limit. When more than `n` bytes are read, subsequent reads return an error whose `Error()` method contains `"http: request body too large"`.

```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit
```

The function takes three arguments:
- `w http.ResponseWriter` -- used to write the error response if the limit is exceeded (it calls `w.WriteHeader(413)` before returning the error).
- `r io.ReadCloser` -- the original request body.
- `n int64` -- the maximum number of bytes to read.

Once triggered, the response code is `413 Request Entity Too Large` (status 413). The connection may also be closed to prevent the client from sending more data.

## Under the hood

`http.MaxBytesReader` returns a `*maxBytesReader` that implements `io.ReadCloser`. It tracks the number of bytes read. On each `Read` call, it checks if the cumulative count exceeds the limit. If it does, it calls `w.WriteHeader(http.StatusRequestEntityTooLarge)` and returns an `errors.New("http: request body too large")`.

The reader does NOT buffer the entire body. It limits as it reads, so streaming handlers remain memory-efficient. The limit applies to the `io.Copy` or `json.Decoder` calls, not just the raw byte count -- JSON decoding reads the body incrementally through the limited reader.

```go
type maxBytesReader struct {
    w       ResponseWriter
    r       io.ReadCloser
    n       int64
    read    int64
    err     error
}
```

## How Go uses it

The Go standard library uses `MaxBytesReader` in several places:
- `net/http` tests validate the behavior of oversized bodies.
- `net/http/cgi` and `net/http/fcgi` use it to protect FastCGI backends.
- Google's own production Go services use similar patterns for body limits.

Frameworks like `gin` and `echo` provide convenience methods (`c.Request.Body = http.MaxBytesReader(...)`) and middleware for body limits. The pattern is universal.

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Comment struct {
	Text string `json:"text"`
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Limit body to 100 bytes
	r.Body = http.MaxBytesReader(w, r.Body, 100)

	var c Comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "comment received: %s", c.Text)
}

func maxBytesMiddleware(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/comment", commentHandler)

	secured := maxBytesMiddleware(1 << 20)(http.HandlerFunc(commentHandler))
	mux.Handle("/secure-comment", secured)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

For a POST to `/comment` with body `{"text":"hello"}` (18 bytes):

1. `commentHandler` is called. It sets `r.Body = http.MaxBytesReader(w, r.Body, 100)`.
2. `json.NewDecoder(r.Body).Decode(&c)` starts reading. The underlying `maxBytesReader` wraps the real body.
3. Each `Read` call returns data from the real body, and `maxBytesReader` increments its internal counter.
4. After 18 bytes, the JSON decoder finishes. Total read: 18 bytes, well under the 100 byte limit.
5. `fmt.Fprintf` writes `"comment received: hello"` as the response.

For a body of 200 bytes:

1. Same setup. The maxBytesReader wraps the body with a 100 byte limit.
2. The JSON decoder reads. After `read` exceeds 100, the next `Read` returns `0, io.EOF` equivalent with an error.
3. `maxBytesReader` calls `w.WriteHeader(http.StatusRequestEntityTooLarge)` behind the scenes.
4. `json.Decode` returns an error containing `"http: request body too large"`.
5. The handler checks for this string and returns `413 Request Entity Too Large`.

## Common mistakes

- Setting the limit after already reading part of the body. `MaxBytesReader` only limits future reads, not past ones.
- Forgetting to check for the `"http: request body too large"` error string. A generic `http.Error(w, "bad request", 400)` loses the distinction between a bad payload and an oversized one.
- Setting limits too high (e.g., 10 MB for a field that is never more than 1 KB). This defeats the purpose. Tune limits to expected payload sizes.
- Not using `MaxBytesReader` in middleware, requiring every handler to repeat the setup. Apply it once at the mux or middleware level.
- Confusing `MaxBytesReader` with `io.LimitReader`. `io.LimitReader` silently stops reading at the limit -- the caller cannot distinguish between a complete body and a truncated one. `MaxBytesReader` returns an error, signaling truncation.

## Debugging walkthrough

A handler that silently truncates request data instead of rejecting it:

```go
func badHandler(w http.ResponseWriter, r *http.Request) {
    limited := io.LimitReader(r.Body, 100)
    body, _ := io.ReadAll(limited)
    fmt.Fprintf(w, "body: %s", body)
}
```

**Symptom**: Large requests appear to succeed but data is silently truncated. No error is returned.

**Investigation**: Check the response body. A 500-byte payload produces the same response as the first 100 bytes. The client has no way to know data was lost.

**Fix**: Replace `io.LimitReader` with `http.MaxBytesReader`, and check the error:

```go
func goodHandler(w http.ResponseWriter, r *http.Request) {
    r.Body = http.MaxBytesReader(w, r.Body, 100)
    body, err := io.ReadAll(r.Body)
    if err != nil {
        // MaxBytesReader already wrote 413 status
        return
    }
    fmt.Fprintf(w, "body: %s", body)
}
```

## Production notes

- Set a global body limit via middleware and per-route overrides for endpoints that need more (e.g., file upload).
- Log oversized requests at WARN level with client IP and content-length header for monitoring.
- Return `413 Request Entity Too Large` with a clear message like `"request body exceeds 1 MB limit"`.
- Consider returning a `Retry-After` header if the limit is part of rate limiting.
- Document limits in your API OpenAPI spec or README.

## Performance implications

- `MaxBytesReader` adds a single integer comparison per `Read` call -- negligible overhead (nanoseconds).
- Early rejection of oversized bodies avoids expensive JSON unmarshalling, database writes, or file processing.
- The body limit prevents memory exhaustion attacks. Without it, a single 1 GB POST can cause OOM on a server with limited RAM.
- Combined with `ReadTimeout` and `WriteTimeout`, body limits form a key part of your DoS defense strategy.

## Practice task

1. Write a handler `profileHandler` that accepts POST with a JSON body `{"bio": "..."}` and limits the body to 500 bytes.
2. Write a middleware `bodyLimitMiddleware` that wraps any handler with a configurable limit.
3. Write table-driven tests for: valid body under limit, body exactly at limit, body over limit by 1 byte, malformed JSON, and GET request (should be rejected as method not allowed).
4. Test that the middleware correctly wraps the handler and returns 413 for oversized payloads.

## Tests / verification

```bash
go test ./curriculum/modules/08-http-rest-apis/lessons/13-body-size-limits -v
```

All tests should pass, covering valid bodies, oversized bodies, malformed JSON, and method rejection.

## Review questions

1. What does `http.MaxBytesReader` do when the body exceeds the limit? Does it call `w.WriteHeader` automatically?
2. Why is `http.MaxBytesReader` preferred over `io.LimitReader` for HTTP handlers?
3. What HTTP status code should you return when a request body exceeds the limit?
4. Can `http.MaxBytesReader` be applied in middleware? If so, how would you access the `http.ResponseWriter` in the middleware?
5. What is the difference in behavior when `MaxBytesReader` triggers versus when the body is simply too large for `json.Decode` for other reasons?

## NEXT UP

Server timeouts -- configuring `ReadTimeout`, `WriteTimeout`, `IdleTimeout`, and `ReadHeaderTimeout` to protect your server from slow clients and hanging connections.
