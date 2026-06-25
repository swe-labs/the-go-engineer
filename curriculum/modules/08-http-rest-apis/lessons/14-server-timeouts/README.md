# Server timeouts

## Learning objective

Configure `ReadTimeout`, `WriteTimeout`, `IdleTimeout`, and `ReadHeaderTimeout` on `http.Server`, set client-side timeouts with `context.WithTimeout`, and understand how each timeout protects against resource exhaustion.

## Why this matters

Without timeouts, a single slow client can hold a goroutine (and its associated memory) indefinitely. A malicious client can open thousands of connections, send headers one byte per minute, and exhaust your server's goroutine pool. Timeouts are the server's immune system -- they kill sick connections before they spread.

## Mental model

Think of each timeout as a different kind of gate with a timer:

- **ReadHeaderTimeout**: The bouncer checking IDs at the door. If the guest doesn't present their ID (headers) quickly enough, they're turned away.
- **ReadTimeout**: The time allowed for the guest to fully enter (read the entire request, including body).
- **WriteTimeout**: The time allowed for the host to respond (write the response).
- **IdleTimeout**: For keep-alive connections, how long the server waits for the next request before closing the door.

Each gate trips independently. No single timeout covers all gates -- you must configure each one.

## Core idea

An `http.Server` struct has four timeout fields, all of type `time.Duration`:

```go
srv := &http.Server{
    Addr:              ":8080",
    ReadTimeout:       5 * time.Second,
    WriteTimeout:      10 * time.Second,
    IdleTimeout:       60 * time.Second,
    ReadHeaderTimeout: 2 * time.Second,
}
```

- `ReadHeaderTimeout` starts after the connection is accepted and the request is read. It limits the time to read the request headers. If exceeded, the connection is closed and a `408 Request Timeout` is returned.
- `ReadTimeout` covers reading the entire request, including the body. It includes `ReadHeaderTimeout` -- if you set both, `ReadTimeout` must be >= `ReadHeaderTimeout`.
- `WriteTimeout` covers writing the response. It starts after the request headers are parsed. If the handler takes too long to write, the connection is closed.
- `IdleTimeout` applies only to keep-alive connections. It's the maximum time between requests on the same connection. If not set, `ReadTimeout` is used.

## Under the hood

The `net/http` server creates a goroutine per connection. Inside that goroutine, it calls `c.serve(ctx)` which loops over requests on the connection (for HTTP/1.1 keep-alive). Each loop iteration:

1. Sets a read deadline on the `net.Conn` via `c.rwc.SetReadDeadline(time.Now().Add(srv.ReadTimeout))`.
2. Calls `c.readRequest(ctx)` which first sets a header-specific deadline if `ReadHeaderTimeout` is set.
3. After headers are read, the read deadline is updated to `ReadTimeout` (if not already elapsed).
4. The handler executes. When the handler calls `w.Write`, the server sets a write deadline: `c.rwc.SetWriteDeadline(time.Now().Add(srv.WriteTimeout))`.
5. After the response is written, `IdleTimeout` is set as the next read deadline for keep-alive.

Deadlines are absolute timestamps on the `net.Conn`. If a deadline passes, the next read or write operation returns a timeout error, and the connection is closed.

## How Go uses it

The Go standard library sets conservative defaults: zero for all timeouts, meaning no timeout at all. Production deployments always set non-zero values. Popular Go web frameworks like `gin`, `echo`, and `chi` document timeout configuration in their deployment guides.

Cloud providers often enforce timeouts at the load balancer level (e.g., AWS ALB has a 60-second idle timeout). Your `http.Server` timeouts should be shorter than the load balancer's timeout to avoid 502 errors.

## Go example

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func slowHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(5 * time.Second):
		fmt.Fprintln(w, "done")
	case <-r.Context().Done():
		return
	}
}

func fastHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "fast response")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/fast", fastHandler)
	mux.HandleFunc("/slow", slowHandler)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Println("Server starting on :8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// Client-side timeout using context
func fetchWithTimeout(url string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}
```

## Step-by-step execution

Client sends GET /fast to server with ReadTimeout=5s, WriteTimeout=10s:

1. TCP connection accepted. Read deadline set to `now + ReadTimeout`.
2. A subset deadline for headers is set to `now + ReadHeaderTimeout` (2s).
3. Client sends headers within 2s. Server parses them.
4. Read deadline shifted to `now + ReadTimeout` (5s from connection start) for the body read. Since GET has no body, this completes immediately.
5. `fastHandler` executes. It writes the response.
6. Write deadline is set to `now + WriteTimeout` (10s) for the first write.
7. Response is written within the deadline.
8. Idle deadline set to `now + IdleTimeout` (60s). If no further request arrives, connection closes.

If the client were a slowloris attack sending one header byte per minute:

1. ReadHeaderTimeout fires after 2 seconds.
2. Server closes connection with 408 status. Goroutine is freed.

## Common mistakes

- Setting `ReadTimeout` but not `WriteTimeout`. Long-running handlers (e.g., streaming, large responses) will still block.
- Not setting `ReadHeaderTimeout`. A slow header attack will succeed within `ReadTimeout` if the body is never sent.
- Setting `IdleTimeout` without enabling keep-alive. `IdleTimeout` only applies to keep-alive connections.
- Confusing `ReadTimeout` with `ReadHeaderTimeout`. `ReadTimeout` includes the entire request read; `ReadHeaderTimeout` is only for headers.
- Forgetting that `WriteTimeout` applies to the entire response write. If your handler uses `http.Flusher` to stream, `WriteTimeout` will terminate the connection.

## Debugging walkthrough

A server that hangs and never responds:

```go
func noTimeout() {
    srv := &http.Server{Addr: ":8080", Handler: mux}
    srv.ListenAndServe()
}
```

**Symptom**: Some requests appear to hang indefinitely. Clients eventually timeout on their end but the server goroutine stays alive.

**Investigation**: Check `srv` struct -- all timeout fields are zero (no timeout). A handler that blocks (e.g., waiting on a channel) never gets killed.

**Fix**: Add explicit timeouts:

```go
srv := &http.Server{
    Addr:              ":8080",
    Handler:           mux,
    ReadTimeout:       5 * time.Second,
    WriteTimeout:      10 * time.Second,
    IdleTimeout:       60 * time.Second,
    ReadHeaderTimeout: 2 * time.Second,
}
```

## Production notes

- Set `ReadTimeout` and `ReadHeaderTimeout` to short values (1-5s) for public-facing APIs.
- Set `WriteTimeout` based on your slowest expected handler. Batch endpoints that return large datasets may need 30-60s.
- Set `IdleTimeout` to match your load balancer's idle timeout or less.
- Client-side timeouts are equally important. Use `context.WithTimeout` or `http.Client.Timeout` to prevent your service from waiting on a dead server.
- Log timeout errors at WARN level with client IP and path for monitoring.

## Performance implications

- Each timeout is a single `SetReadDeadline` / `SetWriteDeadline` syscall per request. Negligible overhead.
- Proper timeouts prevent goroutine and memory leaks from hanging connections.
- Aggressive timeouts (too short) cause premature disconnections. Monitor your p99 response time and set `WriteTimeout` to at least 2x that value.
- Without timeouts, a slow client can hold resources indefinitely, effectively a DoS vector.

## Practice task

1. Create a server with `ReadTimeout=3s`, `WriteTimeout=5s`, `IdleTimeout=30s`, `ReadHeaderTimeout=1s`.
2. Write a handler that reads the request body and echoes it back. Verify timeouts are set correctly via test.
3. Write a handler that simulates a long computation (2s sleep). Confirm it completes within the `WriteTimeout`.
4. Write a client function `fetchWithTimeout` and test it with a timeout that is both sufficient and insufficient for the handler.

## Tests / verification

```bash
go test ./curriculum/modules/08-http-rest-apis/lessons/14-server-timeouts -v
```

Tests verify handler behavior, server configuration, and client-side timeout behavior with `context.WithTimeout`.

## Review questions

1. What is the difference between `ReadTimeout` and `ReadHeaderTimeout`? Can `ReadHeaderTimeout` be longer than `ReadTimeout`?
2. Which timeout would protect against a slowloris attack that sends headers slowly?
3. If `WriteTimeout` is set to 5 seconds and a handler uses `http.Flusher` to stream a response over 10 seconds, what happens?
4. How does `IdleTimeout` interact with HTTP keep-alive? What happens if `IdleTimeout` is 0?
5. How would you set a client-side timeout when making an outgoing HTTP request from a Go service?

## NEXT UP

Graceful shutdown -- shutting down an HTTP server without dropping in-flight requests using `http.Server.Shutdown` and signal handling.
