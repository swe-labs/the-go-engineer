# Opslane graceful shutdown

## Mission

Understand and apply Opslane graceful shutdown in the context of professional Go software engineering.

## Prerequisites

- opslane-12

## Mental Model

Graceful shutdown is a coordinated sequence: stop accepting new work, wait for active work to complete, close shared resources, and exit. The order matters — you cannot close the database before the HTTP server finishes draining. Each goroutine must have a shutdown signal (context cancellation or channel close) and must complete within a deadline. The shutdown sequence is the inverse of the startup sequence: start workers last, stop them first.

## Visual Model

```text
input -> concept boundary -> behavior -> observable result
```

## Machine View

signal.NotifyContext creates a context that is canceled when any of the specified signals arrive. Internally, it calls signal.Notify with a buffered channel and spawns a goroutine that listens on the channel and cancels the context. This is simpler than managing a raw signal channel manually. http.Server.Shutdown(ctx) closes the net.Listener, sets a flag on the server to reject new connections, and waits for the Server.activeConnections counter to reach 0. Active connections are tracked via a sync.WaitGroup — each connection's goroutine calls Add(1) on accept and Done() on close. When all connections complete (or the context deadline expires), Shutdown returns. If the deadline expires before all connections drain, Shutdown returns the context's error and the remaining connections are terminated. db.Close() closes the connection pool: it closes all idle connections immediately and waits for in-use connections to return to the pool before closing them (up to a timeout).

## Run Instructions

```bash
Read the lesson and complete the practice task.
go test ./curriculum/modules/18-flagship-opslane/assessments/13-opslane-graceful-shutdown
```

## Code Walkthrough

Read the code from top to bottom. Connect each line to the mental model above. The important question is not only what the line does, but why the program needs that line.

## Try It

1. Run the command exactly as shown.
2. Change one input.
3. Predict the output before running again.
4. Explain the result in your own words.

Common mistakes:

- Closing the HTTP listener before draining active connections — in-flight requests are cut mid-execution, causing partial writes, corrupted responses, and client errors.
- Using signal.Notify with an unbuffered channel — the signal is sent to the channel synchronously; if the channel is not ready to receive (main goroutine is busy), the signal is dropped and the process never shuts down.
- Not setting a shutdown timeout — if a request hangs indefinitely, the server waits forever, and the rolling update never completes (minReadySeconds never satisfied).
- Ignoring background goroutines in shutdown — the HTTP server shuts down gracefully, but background workers (log flushers, metric exporters, queue consumers) keep running or crash mid-operation.
- Stopping the database connection pool before the HTTP server finishes draining — in-flight requests try to query the database, fail with 'pool closed', and produce 500 errors for the last few requests.

## In Production

Every production Go service in Kubernetes depends on graceful shutdown. Kubernetes sends SIGTERM during rolling updates, scale-downs, and node drains. If the process does not handle SIGTERM, connections are cut, requests fail, and clients see 502 Bad Gateway from the load balancer. Opslane's graceful shutdown must handle: HTTP drain (30s timeout), background worker drain (queue messages are NACKed back to the queue if not processed), database connection close (pending queries complete or are canceled), and metrics flush (last batch of metrics is pushed before exit).

## Thinking Questions

1. What problem does this concept solve?
2. What would break if this concept were removed?
3. How would you test the behavior?

## Next Step

Continue to the next item listed in metadata: `opslane-14`.
