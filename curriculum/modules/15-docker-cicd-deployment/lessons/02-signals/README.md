# Signals

## Learning objective

Handle OS signals (SIGINT, SIGTERM, SIGHUP) in Go programs using `os/signal.Notify`, implement graceful shutdown in HTTP servers, and avoid the pitfalls of signal handling in containerized environments.

## Why this matters

When Kubernetes decides to stop your pod, it sends SIGTERM and waits for the grace period. If your Go service does not handle SIGTERM, it is killed forcefully (SIGKILL) after the timeout, potentially dropping in-flight requests, corrupting data, or leaving database connections open. Proper signal handling is the difference between a graceful shutdown and a silent disaster.

## Mental model

Signals are OS-level interrupts delivered to a process. Think of them as urgent kernel-to-process messages: "stop what you are doing immediately" (SIGTERM), "user pressed Ctrl+C" (SIGINT), "terminal disconnected" (SIGHUP), "illegal instruction" (SIGILL). A signal handler is like a fire alarm — it must be fast, non-blocking, and designed to run when the system is in an unpredictable state. Go's `signal.Notify` converts these OS interrupts into Go channel messages, allowing goroutines to handle them synchronously.

## Core idea

The `os/signal` package connects OS signals to Go channels. The key function is `signal.Notify(ch chan<- os.Signal, sig ...os.Signal)`, which registers the channel to receive the specified signals. Signals are delivered to the channel via a non-blocking send — if the channel is full, the signal is dropped. Always use a buffered channel of size 1 for signal delivery.

| Signal | Default action | Typical use |
|---|---|---|
| SIGINT (Ctrl+C) | Terminate | Interactive shutdown |
| SIGTERM | Terminate | Graceful shutdown (Kubernetes) |
| SIGKILL | Terminate (cannot be caught) | Force kill |
| SIGHUP | Terminate | Config reload, log rotation |
| SIGPIPE | Terminate | Broken pipe (writing to closed connection) |

## Under the hood

`signal.Notify` registers a channel in a global map of signal->channels inside the `os/signal` package. When a signal arrives, the Go runtime's signal handler (installed in `sigtramp`) looks up the registered channels and sends the signal value to each channel in a non-blocking `select`. If the channel is full, the send is skipped — the signal is dropped. `signal.NotifyContext` creates a cancel context and spawns a goroutine that listens on the signal channel, calling `cancel()` when a signal arrives. Go replaces the OS signal handler with its own `sigtramp` handler that queues signals in a global ring buffer and wakes the runtime's signal goroutine.

## How Go uses it

Every production Go HTTP server should handle SIGTERM/SIGINT for graceful shutdown:
- HTTP server drains active connections before closing
- Database connection pool is closed cleanly
- In-flight requests complete (within a deadline)
- Metrics are flushed
- Message queue consumers NACK their in-flight messages

Go's standard library `net/http` provides `Server.Shutdown(ctx)` which blocks until all connections are drained or the context expires — exactly what you need in a signal handler.

## Go example

```go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	srv := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	slog.Info("received signal, shutting down", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
	fmt.Println("server exited gracefully")
}
```

## Step-by-step execution

1. A buffered channel `quit` of size 1 is created.
2. `signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)` registers the channel.
3. The server starts in a goroutine.
4. The main goroutine blocks on `<-quit`, waiting for a signal.
5. When the user presses Ctrl+C, the kernel sends SIGINT to the process.
6. Go's runtime signal handler receives SIGINT and does a non-blocking send on `quit`.
7. The main goroutine receives the signal and logs it.
8. `srv.Shutdown(ctx)` is called. The server stops accepting new connections and waits for active requests to complete.
9. If the context times out (10s), `Shutdown` returns even if connections are still active.
10. The program exits.

## Common mistakes

- **Using an unbuffered channel for signals.** When the channel is not ready to receive, the signal is dropped. Always use a buffered channel of size 1.
- **Not handling SIGTERM in Docker containers.** Docker sends SIGTERM to PID 1 inside the container. If your Go binary uses `CMD ["/app/server"]` (exec form), it receives SIGTERM directly. If you use `CMD /app/server` (shell form), SIGTERM is sent to the shell, not your Go program, and your process is killed abruptly.
- **Calling `os.Exit(1)` in a signal handler.** This prevents deferred cleanup from running. Instead, let the main function return naturally after cleanup.
- **Forgetting that `signal.Notify` affects the whole process.** Signal handlers are process-wide, not goroutine-scoped. Calling `signal.Notify` in one goroutine registers the channel globally.

## Debugging walkthrough

Consider a Go service that is killed by Kubernetes with no logs after a SIGTERM:

```go
signal.Notify(quit, syscall.SIGINT)
// missing syscall.SIGTERM
```

**Symptom:** The service does not shut down gracefully when Kubernetes sends SIGTERM. It is killed after the grace period expires.

**Investigation:** Check which signals are registered. Kubernetes sends SIGTERM, but the code only handles SIGINT. Add logging when signals are received and inspect the signal type.

**Root cause:** SIGTERM was not included in the `signal.Notify` call. Without a registered handler, the default action for SIGTERM is to terminate the process immediately — no shutdown logic runs.

**Fix:** Add `syscall.SIGTERM` to the `signal.Notify` call. Also consider using `signal.NotifyContext` for simpler context propagation.

## Production notes

- Always test your signal handling with `docker kill -s SIGTERM <container>` or `kill -TERM <pid>`.
- Set `terminationGracePeriodSeconds` in Kubernetes Pod spec to at least 30 seconds for Go services.
- Log the signal type and the shutdown phase for debugging abrupt terminations.
- Use `signal.NotifyContext` to propagate cancellation through the entire application context tree.
- Register `syscall.SIGPIPE` in long-lived services that write to network connections — Go ignores SIGPIPE by default, but some C libraries called via cgo may not.
- For config reload on SIGHUP, re-read config files and apply them atomically.

## Performance implications

- Signal handling is inherently synchronous at the OS level. A slow signal handler blocks the signal delivery to all other threads.
- Go's signal goroutine runs at a high priority. Keep handlers minimal and avoid blocking operations.
- `signal.Notify` creates internal goroutines and maps; memory overhead is negligible (one goroutine per package, plus one map entry per signal).
- The non-blocking send on the signal channel means your handler goroutine must read from the channel promptly or signals will be dropped.

## Practice task

Write a function `GracefulShutdown(addr string, timeout time.Duration) error` that:
- Starts an HTTP server on `addr` with a root handler that writes "ok".
- Returns immediately (runs the server in a goroutine).
- The caller provides a channel; when a signal is received, the server shuts down with the given timeout.
- Returns any error from `Shutdown`.

Then write a `main()` that calls this function and sends itself SIGTERM after 100ms. Verify graceful shutdown completes.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/02-signals
go test -v ./curriculum/modules/15-docker-cicd-deployment/lessons/02-signals
```

## Review questions

1. Why must the channel passed to `signal.Notify` be buffered with size 1?
2. What is the difference between SIGTERM and SIGKILL? Why should applications prefer handling SIGTERM?
3. How does `signal.NotifyContext` simplify graceful shutdown compared to raw `signal.Notify`?
4. What happens if your Go binary is started with `CMD ["./app"]` vs `CMD ./app` in a Dockerfile?
5. Why does `http.Server.Shutdown()` exist instead of just calling `os.Exit`?

## NEXT UP

Logs — where to write them (stdout/stderr), structured logging with slog, and the 12-factor app approach.
