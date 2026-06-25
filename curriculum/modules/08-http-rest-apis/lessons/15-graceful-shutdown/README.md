# Graceful shutdown

## Learning objective

Implement graceful HTTP server shutdown using `http.Server.Shutdown`, handle OS signals (SIGINT/SIGTERM), use context deadlines to bound drain time, and register shutdown hooks for cleanup.

## Why this matters

When a server is killed (`Ctrl+C`, deployment rollout, OS reboot), in-flight requests are dropped, connections are severed, and clients see `connection reset by peer`. Graceful shutdown tells the server: "stop accepting new requests, finish the ones in progress, then exit". This is table stakes for production services -- zero-downtime deployments depend on it.

## Mental model

Think of graceful shutdown like closing a store at the end of the day:

1. You lock the door (stop accepting new customers / requests).
2. You serve the customers already inside (finish in-flight requests).
3. You give remaining customers a deadline ("we close in 10 minutes" -- context timeout).
4. After the deadline, anyone still inside is politely escorted out (context cancellation kills lingering handlers).
5. You turn off the lights (exit the process).

`http.Server.Shutdown` orchestrates steps 1-4. You provide the "10 minutes" as a context deadline.

## Core idea

`http.Server.Shutdown(ctx)` does two things:

1. Immediately closes the listener; no new connections are accepted.
2. Waits for all active connections to finish. A connection is "active" from the moment `ServeHTTP` is called until the handler returns. If a handler is still running when the context expires, `Shutdown` returns the context error and the remaining connections are terminated.

The `RegisterOnShutdown` method adds callbacks that run after the listener is closed but before `Shutdown` returns. Use it for cleanup like closing database connection pools or flushing logs.

## Under the hood

Internally, `Shutdown` works by:

1. Setting a `closed` flag on the server, protected by a mutex.
2. Closing the `net.Listener` so `Serve` returns `http.ErrServerClosed`.
3. Iterating over a map of active connections and calling `close` on each, which triggers `SetReadDeadline` to the shutdown time.
4. Waiting on a `sync.WaitGroup` for each connection's goroutine to finish.
5. If the context expires, it closes all remaining connections forcefully via `SetReadDeadline(time.Now())`, causing pending reads to fail immediately.

The server does NOT interrupt handlers by default. Handlers must respect `r.Context().Done()` to be truly cooperative. If your handler ignores context cancellation, `Shutdown` will wait until the context deadline or the handler finishes, whichever comes first.

## How Go uses it

The Go standard library uses graceful shutdown in:
- `net/http/pprof` -- the profiling server can be shut down gracefully.
- `net/http/httptest` -- `Server.Close` waits for active requests.
- Google's production Go infrastructure relies on `Shutdown` for Kubernetes pod termination.

Nearly every production Go HTTP server follows the signal → shutdown → exit pattern shown in the example.

## Go example

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func fastHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(3 * time.Second):
		fmt.Fprintln(w, "done")
	case <-r.Context().Done():
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/fast", fastHandler)
	mux.HandleFunc("/slow", slowHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Run server in background
	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Wait for signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("received signal %s, shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}
```

## Step-by-step execution

Server is running with 3 slow handlers in progress when SIGTERM arrives:

1. OS sends SIGTERM to the process.
2. `signal.Notify` forwards it to the `quit` channel.
3. Main goroutine receives the signal, logs it, and calls `srv.Shutdown(ctx)` with a 10-second deadline.
4. `Shutdown` closes the listener. New connections are rejected with `connection refused`.
5. `Shutdown` marks the 3 active connections as "closing" and waits.
6. Each handler's context (`r.Context()`) is cancelled.
7. Handlers that respect `r.Context().Done()` return immediately. Handlers that ignore it keep running.
8. After 10 seconds, the context expires. Any remaining connections are force-closed via `SetReadDeadline`.
9. `Shutdown` returns. The program exits.

If all handlers finish before the deadline, step 8 is skipped and shutdown is clean.

## Common mistakes

- Forgetting to run `ListenAndServe` in a goroutine. `Shutdown` can only be called after `Serve` starts, and `ListenAndServe` blocks.
- Not checking for `http.ErrServerClosed` in the `ListenAndServe` error handler. After shutdown, `Serve` returns this sentinel error. Treating it as fatal causes a spurious crash.
- Setting the shutdown timeout too short. Monitor your p99 handler duration and set the timeout to at least 2x that value.
- Not calling `signal.Notify` with a buffered channel. If the signal arrives before the channel is read, it's dropped.
- Forgetting `defer cancel()` on the shutdown context, causing a context leak.
- Handlers that ignore context cancellation. Always check `r.Context().Done()` in long-running handlers.

## Debugging walkthrough

A server that never finishes shutting down:

```go
func hangingHandler(w http.ResponseWriter, r *http.Request) {
    time.Sleep(30 * time.Second) // ignores context
    fmt.Fprintln(w, "done")
}
```

**Symptom**: `Shutdown` waits the full context timeout before force-killing the connection.

**Investigation**: The handler does not select on `r.Context().Done()`. `Shutdown` cancels contexts but has no way to preempt `time.Sleep`.

**Fix**: Make the handler context-aware:

```go
func goodHandler(w http.ResponseWriter, r *http.Request) {
    select {
    case <-time.After(30 * time.Second):
        fmt.Fprintln(w, "done")
    case <-r.Context().Done():
        return
    }
}
```

## Production notes

- Always set a shutdown timeout (10-30 seconds is typical). An infinite wait is a hidden availability risk.
- Register cleanup callbacks with `RegisterOnShutdown` for database migrations, cache flushes, and connection pool draining.
- Use `shutdownWithWG` (WaitGroup pattern) to ensure background goroutines have completed before the process exits.
- In Kubernetes, set `terminationGracePeriodSeconds` in the pod spec to slightly more than your shutdown timeout.
- Monitor `http.Server` shutdown duration as a metric. Unexpectedly long shutdowns indicate handlers ignoring context.

## Performance implications

- `Shutdown` itself is cheap -- it sets a flag and closes a listener.
- The cost is in waiting. Handlers that finish quickly make shutdown nearly instant (milliseconds).
- The shutdown context timeout bounds the worst-case shutdown duration. Set it as low as your handlers allow.
- Each active connection during shutdown consumes a goroutine. If you have thousands of connections, shutdown can take seconds even with cooperative handlers.

## Practice task

1. Write a server with two handlers: `/fast` (responds immediately) and `/stream` (sends a chunk every second for 10 seconds).
2. Implement graceful shutdown with a 5-second timeout.
3. Test that `/fast` completes during shutdown but `/stream` is terminated by the timeout.
4. Add a `RegisterOnShutdown` callback that logs "cleanup complete".
5. Write a test that starts the server, sends a request to `/stream`, triggers shutdown, and verifies the handler is cancelled.

## Tests / verification

```bash
go test ./curriculum/modules/08-http-rest-apis/lessons/15-graceful-shutdown -v
```

Tests validate handler behavior, server shutdown mechanics, and context cancellation on handlers.

## Review questions

1. What does `http.Server.Shutdown` do with the listener? With active connections?
2. How does `Shutdown` interact with handler context cancellation? What must handlers do to be "graceful"?
3. Why do you need to run `ListenAndServe` in a goroutine when using graceful shutdown?
4. What is the purpose of `RegisterOnShutdown`? Give two real-world examples.
5. What happens if the context passed to `Shutdown` expires before all handlers finish?

## NEXT UP

Health and readiness probes -- implementing `/healthz` and `/readyz` endpoints for Kubernetes and load balancer health checking.
