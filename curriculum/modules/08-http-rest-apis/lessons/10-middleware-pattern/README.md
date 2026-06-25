# Middleware pattern

## Learning objective

Implement HTTP middleware using `func(http.Handler) http.Handler`, build a middleware chain with logging, recovery, CORS, and request ID injection, and understand the before/after execution model.

## Why this matters

Cross-cutting concerns — logging, authentication, rate limiting, request tracing, CORS — appear in every HTTP handler. Without middleware, you duplicate this code in every handler. With middleware, you write it once and wrap any handler. This is the single most important design pattern in Go web programming. Every production Go service uses middleware, from the simplest API to the largest microservice mesh.

## Mental model

Think of middleware as layers of an onion. The core handler is at the center. Each middleware layer wraps the previous one. A request enters the outermost layer, passes through each layer (doing pre-processing), reaches the core handler, then passes back through each layer (doing post-processing). Each layer can short-circuit: if auth middleware rejects the request, it never reaches the core. If recovery middleware catches a panic, it never reaches the client as a crash.

```
Request → [Logger] → [Auth] → [CORS] → [Handler] → Response
                ←        ←        ←        ←
```

## Core idea

In Go, middleware is a function that takes an `http.Handler` and returns an `http.Handler`:

```go
type Middleware func(http.Handler) http.Handler
```

A middleware function calls the inner handler's `ServeHTTP` optionally before and after doing its own work:

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        // Pre-processing
        log.Printf("start %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r) // Call the inner handler
        // Post-processing
        log.Printf("end %s %s (%s)", r.Method, r.URL.Path, time.Since(start))
    })
}
```

Middleware is applied by wrapping:

```go
handler := loggingMiddleware(authMiddleware(corsMiddleware(coreHandler)))
```

Or with a helper:

```go
chain := applyMiddleware(coreHandler, loggingMiddleware, authMiddleware, corsMiddleware)
```

Common middleware types:
- **Logging**: log method, path, duration, status code.
- **Recovery**: recover from panics, return 500.
- **CORS**: set `Access-Control-Allow-*` headers.
- **Request ID**: inject a unique ID into context.
- **Auth**: validate JWT or API key, reject with 401.
- **Rate limiting**: check token bucket, reject with 429.

## Under the hood

Each middleware creates a closure that captures the inner handler and any configuration. When `ServeHTTP` is called on the returned `HandlerFunc`, it runs the pre-logic, calls `next.ServeHTTP(w, r)` (which may be another middleware or the final handler), then runs the post-logic. The `next.ServeHTTP` call is synchronous — the middleware waits for the inner handler to complete before continuing. This means post-processing (like logging duration) can observe the final status code and response size.

## How Go uses it

```go
func withCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

Middleware is applied at server startup or route registration, not per-request:

```go
mux := http.NewServeMux()
mux.Handle("/api/", withCORS(withLogger(coreHandler)))
```

## Go example

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Middleware type alias.
type Middleware func(http.Handler) http.Handler

func applyMiddleware(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// Logging middleware.
func withLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}

// CORS middleware.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Recovery middleware.
func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("PANIC recovered: %v", rec)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from middleware chain!")
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("something went wrong")
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/hello", applyMiddleware(
		http.HandlerFunc(helloHandler),
		withLogger, withCORS, withRecovery,
	))
	mux.Handle("/panic", applyMiddleware(
		http.HandlerFunc(panicHandler),
		withLogger, withCORS, withRecovery,
	))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

1. Client sends `OPTIONS /hello` (preflight CORS request).
2. `withRecovery` sets up defer/recover, calls `withCORS.ServeHTTP`.
3. `withCORS` sets CORS headers, sees `OPTIONS` method → returns 204 without calling next.
4. Client sends `GET /hello`.
5. `withRecovery` → `withCORS` → `withLogger` → `helloHandler`.
6. `withLogger` records start time, calls `helloHandler`.
7. `helloHandler` writes "Hello from middleware chain!".
8. `withLogger` logs the duration after `helloHandler` returns.
9. `withCORS` had already set headers (pre-processing only).
10. `withRecovery` defer checks: no panic, returns normally.
11. Response sent to client.

## Common mistakes

- Applying middleware in the wrong order. The first middleware in the list is the outermost. Logging should be outermost (to capture the full duration including other middleware). Recovery should be outermost (to catch panics from all layers).
- Not calling `next.ServeHTTP(w, r)` in a middleware. If you forget, the chain breaks silently and the client gets no response (or a 200 with no body if WriteHeader was called).
- Setting response headers after `next.ServeHTTP`. Once the inner handler writes the status, headers are locked. Set or modify headers before calling `next.ServeHTTP`.
- Writing to `ResponseWriter` both before and after `next.ServeHTTP`. If the inner handler writes the status code, your pre-write may have set a different status that gets overridden.
- Not handling `OPTIONS` preflight in CORS middleware. Browsers send `OPTIONS` first; if it's not handled, cross-origin requests fail.

## Debugging walkthrough

A middleware chain logs requests but the logged status code is always 200, even when handlers return 404:

```go
func withLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w, r)
        log.Printf("Status: %d", ???) // Can't read the status!
    })
}
```

**Root cause**: `http.ResponseWriter` doesn't expose the status code. The middleware can't see what the inner handler wrote.

**Fix**: Wrap `ResponseWriter` with a custom type that captures the status code:

```go
type responseWriter struct {
    http.ResponseWriter
    status int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
}

func withLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rw := &responseWriter{ResponseWriter: w, status: 200}
        next.ServeHTTP(rw, r)
        log.Printf("%s %s -> %d", r.Method, r.URL.Path, rw.status)
    })
}
```

## Production notes

- Keep middleware stateless where possible. Configuration (log level, CORS origins) should be set at construction, not per-request.
- Order matters: Recovery outermost, then RequestID/Logging, then Auth, then CORS, then handler.
- Use `http.Handler` wrapping, not `http.HandleFunc` adapter, when building middleware chains. The explicit type makes the wrapping order visible.
- Third-party middleware libraries (Chi, Alice) provide `chi.Use()` and `alice.New()` for cleaner chain construction. For most projects, a simple `applyMiddleware` helper is sufficient.
- Middleware that modifies the request (adding values to context) should use `r.WithContext(ctx)` and pass the new request to `next.ServeHTTP`.

## Performance implications

- Each middleware layer adds two function calls (one to enter the middleware, one to call `next`). With 3-7 middleware layers, this is <1µs overhead.
- Wrapping `ResponseWriter` to capture status allocates one additional object per request. In high-throughput services (>10K req/s), consider pooling or using an atomic field instead.
- Logging middleware that formats and writes logs synchronously can add 10-100µs per request. Use async logging or structured loggers that minimize allocations.
- Recovery middleware's `defer` has a small cost (~30ns). This is negligible.

## Practice task

Build a middleware chain with:
1. `withRequestID` — generates a UUID (or simple random hex) and injects it into `r.Context()`.
2. `withLogger` — logs method, path, status, and the request ID from context.
3. `withAuth` — checks `Authorization: Bearer <token>` header, rejects with 401 if missing or token != `"secret"`.
4. An inner handler that returns `{"user":"data"}` as JSON.
5. Apply middleware in the correct order. Test with httptest for both valid auth and invalid auth.

## Tests / verification

```bash
go run ./curriculum/modules/08-http-rest-apis/lessons/10-middleware-pattern
go test ./curriculum/modules/08-http-rest-apis/lessons/10-middleware-pattern
```

## Review questions

1. What is the signature of a Go HTTP middleware function?
2. Why should recovery middleware be the outermost layer?
3. How can a middleware read the HTTP status code written by the inner handler?
4. What happens if a middleware never calls `next.ServeHTTP`?
5. How does CORS middleware handle the OPTIONS preflight request?

## NEXT UP

Error middleware — centralized error handling and panic recovery in middleware, returning structured errors.
