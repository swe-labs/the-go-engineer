# Error middleware

## Learning objective

Implement centralized error handling in middleware, build a panic recovery middleware that returns structured JSON error responses, and propagate errors from handlers to middleware for consistent error formatting.

## Why this matters

Without centralized error handling, every handler must write its own error response logic. This leads to inconsistent error shapes, missing status codes, and handlers that crash the server on panic. A single error middleware layer catches every unhandled error and panic, formats it consistently, and logs it. This is the difference between a server that returns a cryptic HTML stack trace on crash and one that calmly returns `{"code":"INTERNAL_ERROR","message":"something went wrong"}` with a correlation ID.

## Mental model

Think of error middleware as a safety net stretched under a tightrope walker. The walker (handler) might slip (panic) or fail (return an error). The net catches everything. It doesn't matter how the walker fails — the net ensures the outcome is always the same: a structured, safe response. The net also logs the details so the crew knows what went wrong, but the audience (client) only sees a sanitized message.

## Core idea

Error middleware combines two patterns:

1. **Panic recovery**: a `defer recover()` in middleware catches any panic from the handler chain. Instead of crashing the process, it logs the stack trace and returns a 500 JSON response.

2. **Error propagation**: handlers communicate errors to middleware by writing to `ResponseWriter` or by using a custom error type that middleware can inspect. The middleware can wrap `ResponseWriter` to detect the status code and ensure every error response uses the same JSON shape.

```go
func withErrorHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                log.Printf("panic recovered: %v", rec)
                writeJSONError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

For non-panic errors, handlers should still use a shared `writeError` function, but middleware can also wrap `ResponseWriter` to ensure that any 5xx response has the correct JSON shape.

## Under the hood

`recover()` only works inside a deferred function. When a panic occurs (e.g., nil pointer dereference, out-of-bounds access, explicit `panic()`), Go unwinds the stack, running deferred functions. If any deferred function calls `recover()`, the panic is caught and the unwinding stops. After `recover()`, the function continues normally. The middleware can then write an error response. If no `recover()` is called, the panic propagates to the goroutine's top level, crashing the server.

For non-panic errors, the middleware can wrap `ResponseWriter` to track whether an error status was written and ensure the body contains a proper JSON error:

```go
type errorWriter struct {
    http.ResponseWriter
    wroteError bool
}
```

## How Go uses it

```go
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // Log the full stack trace
                log.Printf("panic: %v\n%s", err, debug.Stack())
                // Return structured error
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusInternalServerError)
                json.NewEncoder(w).Encode(APIError{
                    Code:    "INTERNAL_ERROR",
                    Message: "An internal error occurred",
                })
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

## Go example

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// withRecovery catches panics and returns structured JSON errors.
func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC: %v\n%s", rec, debug.Stack())
				writeJSON(w, http.StatusInternalServerError, APIError{
					Code:    "INTERNAL_ERROR",
					Message: "An unexpected error occurred",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// withErrorTracking wraps ResponseWriter to catch handlers that write
// error status codes without a proper JSON body.
type errorTrackingWriter struct {
	http.ResponseWriter
	status int
}

func (w *errorTrackingWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func withErrorTracking(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tw := &errorTrackingWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(tw, r)
		if tw.status >= 500 {
			log.Printf("ERROR: %s %s returned %d", r.Method, r.URL.Path, tw.status)
		}
	})
}

func crashHandler(w http.ResponseWriter, r *http.Request) {
	panic("database connection lost")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotFound, APIError{
		Code:    "NOT_FOUND",
		Message: "The requested resource was not found",
	})
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/ok", withRecovery(withErrorTracking(http.HandlerFunc(okHandler))))
	mux.Handle("/crash", withRecovery(withErrorTracking(http.HandlerFunc(crashHandler))))
	mux.Handle("/notfound", withRecovery(withErrorTracking(http.HandlerFunc(notFoundHandler))))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

1. Client sends `GET /crash`.
2. `withRecovery` sets up `defer recover()`.
3. `withErrorTracking` wraps `ResponseWriter`.
4. `crashHandler` calls `panic("database connection lost")`.
5. Stack unwinds: `errorTrackingWriter.WriteHeader` is not called.
6. `withRecovery` defer runs `recover()` → catches the panic.
7. Logs the panic and stack trace.
8. Writes `{"code":"INTERNAL_ERROR","message":"An unexpected error occurred"}` with 500.
9. Client receives a clean JSON error instead of a crashed server.
10. For `GET /notfound`, `notFoundHandler` calls `writeJSON` with 404 and a structured error.
11. `errorTrackingWriter` logs the 404 status.
12. Client receives `{"code":"NOT_FOUND","message":"The requested resource was not found"}`.

## Common mistakes

- Not calling `recover()` in a defer. If the defer doesn't call `recover()`, the panic continues unwinding. The defer must be the first thing in the middleware.
- Calling `recover()` outside of a deferred function. `recover()` returns `nil` if not called directly inside a `defer`.
- Writing to `ResponseWriter` after `recover()` without checking if headers were already written. If the handler had partially written a response, writing again may corrupt the output. Track whether headers were sent.
- Only handling panics, not regular errors. Panic recovery is necessary but not sufficient. Handlers can still return 500 with plain text via `http.Error`. Error tracking middleware ensures all 5xx responses are consistent.
- Logging the full stack trace in production responses. Always log it server-side, but return a sanitized message to the client.

## Debugging walkthrough

A server crashes silently with no log output:

**Symptom**: The process exits with a panic but no log message appears.

**Root cause**: The panic occurs in a goroutine that is not the main request goroutine. For example, a handler spawns a goroutine that panics:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    go func() {
        panic("in goroutine") // Not caught by recovery middleware!
    }()
}
```

**Fix**: Recovery middleware only catches panics in the request goroutine. Panics in other goroutines must be handled separately. Either use an error channel to propagate the panic back to the handler, or add recovery to every goroutine:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    go func() {
        defer func() {
            if rec := recover(); rec != nil {
                log.Printf("goroutine panic: %v", rec)
            }
        }()
        // dangerous code
    }()
}
```

## Production notes

- Always pair recovery middleware with a structured logger that captures the full stack trace (`debug.Stack()`). Include the request ID if available.
- Set `GOTRACEBACK=single` or `GOTRACEBACK=crash` in production to control behavior on unrecovered panics. `single` prints only the crashing goroutine's stack.
- Use `http.Server` with `ErrorLog` to capture server-level errors (TLS handshake failures, etc.) separately from application errors.
- Consider using `sentry-go` or similar for centralized error reporting. Integrate it into the recovery middleware.
- Test your recovery middleware with a handler that panics at different stages: before writing headers, after writing partial body, after writing full body.

## Performance implications

- `defer` has a small overhead (~30-50ns per defer). Recovery middleware typically has one defer, so the cost is negligible.
- `debug.Stack()` allocates memory proportional to the stack depth. For deep stacks (100+ frames), this can be expensive. In hot paths, consider `debug.Stack()` only on panics (which are rare).
- Writing a JSON error response in the recovery path adds one allocation for the error struct and one JSON encode. This is fine since panics are exceptional.
- Wrapping `ResponseWriter` to track status adds one allocation per request. For 10K req/s, this is ~80KB/s — acceptable.

## Practice task

Build an error middleware that:
1. Recovers from panics and returns `{"code":"PANIC","message":"internal error"}` with status 500.
2. Logs the panic error and stack trace.
3. Also tracks all 4xx and 5xx responses from handlers.
4. For 5xx responses that don't have a JSON body, writes a default `{"code":"SERVER_ERROR","message":"internal error"}`.

Create handlers that test panic, 404 via `writeJSON`, and a handler that calls `http.Error` (plain text 500). Verify that the middleware catches all cases and returns consistent JSON.

## Tests / verification

```bash
go run ./curriculum/modules/08-http-rest-apis/lessons/11-error-middleware
go test ./curriculum/modules/08-http-rest-apis/lessons/11-error-middleware
```

## Review questions

1. Why must `recover()` be called inside a deferred function?
2. What is the difference between catching a panic in middleware vs. catching it in each handler?
3. How can middleware ensure that every 5xx response has a consistent JSON shape?
4. What happens to panics in goroutines spawned by a handler? How would you catch them?
5. Why should the full stack trace be logged server-side but never sent to the client?

## NEXT UP

Handler testing with httptest — writing table-driven tests for HTTP handlers using `httptest.NewRecorder` and `httptest.NewRequest`.
