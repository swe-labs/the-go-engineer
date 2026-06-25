# net/http basics

## Learning objective

Start an HTTP server using `http.ListenAndServe`, register handlers with `http.HandleFunc`, read the default `ServeMux` behavior, and write responses with proper status codes and content types.

## Why this matters

Every Go web service — from a tiny health-check endpoint to a multi-service production API — starts with the same three functions: `http.HandleFunc`, `http.ListenAndServe`, and `w.Write`/`w.WriteHeader`. There is no Rails or Django equivalent in Go; the standard library is the framework. Mastering these primitives means you can build, deploy, and debug any Go HTTP service without reaching for a third-party dependency.

## Mental model

Think of `net/http` as a dispatcher. You give it a map (the `ServeMux`): "when someone visits `/hello`, call this function." Then you tell it to listen on a port. Every incoming TCP connection is accepted, parsed into an `http.Request`, and handed to the matching handler function. The handler receives a `ResponseWriter` — a pipe to the client — and writes headers, status, and body into it. When the handler returns, the response is flushed and the connection is available for the next request.

## Core idea

The `net/http` server is built on three core pieces:

1. **Handler functions**: functions matching `func(w http.ResponseWriter, r *http.Request)`. Every request is dispatched to one of these.
2. **ServeMux** (multiplexer): a request router. The default one is `http.DefaultServeMux`. It matches the request path against registered patterns using longest-prefix matching.
3. **Server**: `http.ListenAndServe(addr, handler)` starts listening on a TCP address. If `handler` is `nil`, it uses `DefaultServeMux`.

The `ResponseWriter` interface has three methods:
- `Write([]byte) (int, error)` — writes the body (implicitly writes StatusOK if not already set)
- `WriteHeader(int)` — sets the status code (must be called before `Write`)
- `Header() http.Header` — returns the header map for setting response headers

## Under the hood

When `http.ListenAndServe(":8080", nil)` is called, Go creates a `net.Listener` on port 8080 and enters an accept loop. For each accepted TCP connection, a new goroutine is spawned. The goroutine reads the request from the connection, parses it into an `http.Request`, looks up the matching handler in `DefaultServeMux`, and calls `handler.ServeHTTP(w, r)`. The handler writes to the `ResponseWriter`, which buffers headers and body. After the handler returns, the server sends the HTTP response line, headers, and body over the connection. The connection is then returned to a keep-alive pool for reuse, or closed if `Connection: close` was requested.

## How Go uses it

Every Go HTTP server starts the same way:

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello"))
})
http.ListenAndServe(":8080", nil)
```

The `HandleFunc` call registers a pattern and function on `DefaultServeMux`. The server loop uses `DefaultServeMux` as its handler because the second argument to `ListenAndServe` is `nil`.

## Go example

```go
package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, Go HTTP!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok"}`)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/health", healthHandler)
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}
```

## Step-by-step execution

1. `main` calls `http.HandleFunc("/hello", helloHandler)` registering the pattern `/hello` with `helloHandler`.
2. `http.HandleFunc("/health", healthHandler)` registers `/health`.
3. `http.ListenAndServe(":8080", nil)` creates a `net.Listener` on port 8080 and starts the accept loop.
4. A client sends `GET /hello HTTP/1.1`.
5. The server parses the request, creates `http.ResponseWriter` and `*http.Request`.
6. `DefaultServeMux` matches `/hello` → calls `helloHandler`.
7. `helloHandler` sets `Content-Type: text/plain; charset=utf-8`.
8. `w.WriteHeader(http.StatusOK)` sends `HTTP/1.1 200 OK`.
9. `fmt.Fprintln(w, "Hello, Go HTTP!")` writes the body.
10. Handler returns, server flushes the response to the client.
11. If the client sends `Keep-Alive`, the connection waits for the next request.

## Common mistakes

- Calling `w.WriteHeader` after `w.Write`. `Write` implicitly sets the status to 200 if not already set. Any subsequent `WriteHeader` call is silently ignored.
- Not setting `Content-Type`. Go's `net/http` will try to detect it via `http.DetectContentType`, which may guess wrong for JSON (it might return `text/plain; charset=utf-8`).
- Using `http.ListenAndServe` without error handling. It returns an error on failure (e.g., port in use). Always check it: `log.Fatal(http.ListenAndServe(...))`.
- Registering handlers with conflicting patterns. `"/"` matches everything. If you register `"/"` and `"/api"`, a request to `"/api/users"` matches `"/api"` first (longest prefix wins).

## Debugging walkthrough

A server compiles but returns 404 for a registered route:

```go
http.HandleFunc("/users", userHandler)
http.ListenAndServe(":8080", nil)
```

**Symptom**: `curl http://localhost:8080/users/` returns 404, but `curl http://localhost:8080/users` works.

**Root cause**: The pattern `/users` matches exactly `/users`. The request `/users/` has a trailing slash and does not match. Pattern `/users/` would match `/users/` and any path under it.

**Fix**: Register both patterns or use the trailing-slash form:

```go
http.HandleFunc("/users", userHandler)
http.HandleFunc("/users/", userHandler)
```

## Production notes

- Never use the default `http.ServeMux` in production without wrapping it in middleware (logging, recovery, timeouts). An unhandled panic kills the process.
- Always set `ReadTimeout` and `WriteTimeout` on the `http.Server` struct to prevent slow-client attacks.
- Listen on a configurable address (from env var or flag), never hardcode `:8080`.
- Use `http.Server` with `ErrorLog` set to a structured logger instead of the default `log` package.

## Performance implications

- Each request gets a new goroutine. Under high concurrency, this can lead to high memory usage. Use `http.Server` with `MaxHeaderBytes` and timeouts to limit resource consumption.
- Go's HTTP server is production-grade: the standard library powers services handling millions of requests per second (Docker, Kubernetes, InfluxDB all use it).
- `ResponseWriter` writes are buffered internally. Large responses are chunked automatically when the content length is unknown.

## Practice task

Write a Go program that:
- Registers three handlers: `/greet` (returns `"Hello, Gopher!"`), `/time` (returns current time as JSON `{"now":"..."}`), and `/version` (returns `"Go 1.24"` as plain text).
- Each handler sets the appropriate `Content-Type`.
- Uses `http.ListenAndServe` with error handling (`log.Fatal`).
- Run with `go run .` and test with `curl` or the test file.

## Tests / verification

```bash
go run ./curriculum/modules/08-http-rest-apis/lessons/02-net-http-basics
go test ./curriculum/modules/08-http-rest-apis/lessons/02-net-http-basics
```

The tests start the server, make HTTP requests to each route, and verify status codes and body content.

## Review questions

1. What happens if you call `w.Write` before `w.WriteHeader`? What status code is sent?
2. Why is it important to set `Content-Type` explicitly rather than relying on Go's auto-detection?
3. What is the difference between the patterns `"/api"` and `"/api/"` in `DefaultServeMux`?
4. How does `http.ListenAndServe(":0", nil)` behave differently from `http.ListenAndServe(":8080", nil)`?
5. What three methods does the `http.ResponseWriter` interface require?

## NEXT UP

Handler lifecycle — understanding the `http.Handler` interface, `ServeHTTP`, and how Go processes requests from connection to response.
