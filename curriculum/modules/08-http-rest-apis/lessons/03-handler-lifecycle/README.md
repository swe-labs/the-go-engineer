# Handler lifecycle

## Learning objective

Implement the `http.Handler` interface with struct-based handlers, understand request-scoped state, and chain handlers together using the `ServeHTTP` method.

## Why this matters

Every HTTP framework in Go — from the standard library to Gin, Chi, and Echo — is built on the `http.Handler` interface. It is the single most important abstraction in Go web programming. When you understand that a handler is any object with a `ServeHTTP` method, you unlock the ability to write middleware, compose handlers, and build testable, request-scoped services without framework lock-in.

## Mental model

Think of HTTP request processing as a pipeline. The server accepts a TCP connection, parses the request into an `*http.Request` and `http.ResponseWriter`, then calls `handler.ServeHTTP(w, r)`. The handler does its work and returns. That's it. A handler can be a function (via `http.HandlerFunc` adapter), a struct with a `ServeHTTP` method, or a chain of middleware wrapping a final handler. Every piece of Go HTTP middleware is just a handler that calls another handler's `ServeHTTP`.

## Core idea

The `http.Handler` interface is one method:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Any type that implements this method can serve HTTP requests. Go provides `http.HandlerFunc` — an adapter that turns a function `func(w ResponseWriter, r *Request)` into a `Handler`. This is what `http.HandleFunc` registers.

A handler struct can hold dependencies (database, logger, config) that are available for every request:

```go
type UserHandler struct {
    DB     *sql.DB
    Logger *log.Logger
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // h.DB and h.Logger are available here
}
```

The request lifecycle is:
1. Server accepts connection → spawns goroutine.
2. Goroutine parses request → builds `Request` and `ResponseWriter`.
3. Calls `mux.ServeHTTP(w, r)` → mux matches pattern → calls handler.
4. Handler runs → writes response.
5. Handler returns → response is flushed.
6. Connection reused or closed.

## Under the hood

`http.ListenAndServe` creates a `net.Listener` and calls `server.Serve(listener)`. The `Serve` method has an infinite loop calling `listener.Accept()`. Each accepted connection creates a new `http.conn` object and spawns `go c.serve(ctx)`. Inside `c.serve`, the server reads the request, creates the `ResponseWriter`, looks up the handler via `mux.Handler(r)`, and calls `handler.ServeHTTP(w, r)`. After the handler returns, `c.serve` checks if the connection should be kept alive (HTTP/1.1 default) and loops to read the next request on the same connection.

Each request runs in its own goroutine, so handlers must synchronize access to shared state. Request-scoped state lives on `r.Context()`, not in global variables.

## How Go uses it

The standard library uses `http.Handler` everywhere:

- `http.ServeMux` implements `http.Handler`. Its `ServeHTTP` method matches the request path and delegates to the registered handler.
- `http.HandlerFunc` is a type `func(ResponseWriter, *Request)` with a `ServeHTTP` method that calls itself. This is how `http.HandleFunc` works.
- `http.NotFoundHandler()`, `http.RedirectHandler()`, `http.TimeoutHandler()` all return `http.Handler`.
- Middleware functions take an `http.Handler` and return an `http.Handler`.

## Go example

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// requestTimer is a struct-based handler that logs how long each request takes.
type requestTimer struct {
	handler http.Handler
}

func (rt *requestTimer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rt.handler.ServeHTTP(w, r)
	log.Printf("%s %s took %v", r.Method, r.URL.Path, time.Since(start))
}

// GreetingHandler holds a dependency (greeting template).
type GreetingHandler struct {
	Greeting string
}

func (gh *GreetingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%s, %s!", gh.Greeting, r.URL.Query().Get("name"))
}

func main() {
	greet := &GreetingHandler{Greeting: "Hello"}
	timed := &requestTimer{handler: greet}

	http.Handle("/greet", timed)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Step-by-step execution

1. Client sends `GET /greet?name=Alice`.
2. Server accepts connection, spawns goroutine, parses request.
3. `DefaultServeMux.ServeHTTP(w, r)` matches `/greet` → calls `timed.ServeHTTP(w, r)`.
4. `requestTimer.ServeHTTP` records `time.Now()`.
5. Calls `gh.ServeHTTP(w, r)` on the inner `GreetingHandler`.
6. `GreetingHandler.ServeHTTP` sets `Content-Type`, reads `r.URL.Query().Get("name")` → `"Alice"`.
7. Writes `"Hello, Alice!"` to response.
8. Returns to `requestTimer.ServeHTTP`, logs `"GET /greet took 1.2ms"`.
9. Returns. Response flushed. Connection kept alive.

## Common mistakes

- Storing per-request state in handler struct fields. Handler structs are shared across goroutines; fields must be read-only or properly synchronized. Use request context (`r.Context()`) for per-request state.
- Forgetting that `ServeHTTP` must be a pointer receiver if the handler struct has mutable fields. Value receivers create copies.
- Not calling the inner handler in middleware. If `mw.handler.ServeHTTP` is never called, the chain breaks and the client gets no response.
- Writing to `ResponseWriter` after the handler returns. Nothing happens (the writer is closed), but no error is reported either.

## Debugging walkthrough

A handler struct with a database field panics with nil pointer dereference:

```go
type UserHandler struct {
    DB *sql.DB
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    row := h.DB.QueryRow("SELECT name FROM users WHERE id = $1", ...)
}
```

**Symptom**: panic on `h.DB.QueryRow`.

**Investigation**: Check where `UserHandler` is created:

```go
http.Handle("/users", &UserHandler{}) // DB is nil!
```

**Root cause**: `DB` field is never initialized.

**Fix**: Ensure `DB` is set before registration:

```go
db, _ := sql.Open("postgres", dsn)
http.Handle("/users", &UserHandler{DB: db})
```

Or use a constructor function that guarantees initialization.

## Production notes

- Handler structs are the idiomatic way to inject dependencies in Go web services. Avoid global variables for DB connections and configs.
- The `http.Handler` interface makes every piece of HTTP logic testable with `httptest.NewRecorder` and `httptest.NewRequest`.
- In production, always wrap handlers with recovery middleware to prevent panics from crashing the server.
- Use `http.Handler` composition to build a middleware chain: each layer adds cross-cutting behavior (logging, auth, rate limiting) without modifying the core handler.

## Performance implications

- Each handler invocation in the chain adds function call overhead, but with Go's inlining and the small number of middleware layers (typically 3-7), the cost is negligible (sub-microsecond).
- Struct-based handlers with many fields increase memory per registration but not per-request allocation, since the struct is shared.
- The goroutine-per-request model means handlers must not block on I/O without releasing resources. Use context deadlines for all outbound calls.

## Practice task

Define a struct `Logger` that implements `http.Handler` and logs every request's method, path, and duration. Define a struct `AdminCheck` that implements `http.Handler` and returns 403 Forbidden unless the `X-Admin` header equals `"true"`. Chain them: create a final handler that returns `"Welcome, admin!"` and wrap it with `AdminCheck` then `Logger`. Test with both a valid and missing admin header.

## Tests / verification

```bash
go run ./curriculum/modules/08-http-rest-apis/lessons/03-handler-lifecycle
go test ./curriculum/modules/08-http-rest-apis/lessons/03-handler-lifecycle
```

## Review questions

1. What is the difference between `http.Handler` and `http.HandlerFunc`?
2. Why must handler struct fields be read-only or synchronized across requests?
3. What happens to the HTTP response if `ServeHTTP` returns without calling `w.Write` or `w.WriteHeader`?
4. How does `http.Handle("/path", handler)` differ from `http.HandleFunc("/path", fn)`?
5. What is the purpose of wrapping one handler inside another's `ServeHTTP` method?

## NEXT UP

Routing — using `ServeMux` patterns, method-based routing, path parameters, and third-party routers like Chi and gorilla/mux.
