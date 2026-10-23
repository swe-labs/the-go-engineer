# Opslane graceful shutdown

## Mission

Implement orchestrator-aligned graceful shutdown and zero-downtime draining for Opslane, combining OS signal capture (`SIGTERM`, `SIGINT`), traffic shedding middleware (HTTP 503 with `Connection: close`), bounded context timeouts, and in-flight request completion via `server.Shutdown(ctx)`.

## Prerequisites

- `opslane-12`

## Mental Model

When cloud orchestrators (Kubernetes, AWS ECS) redeploy pods or scale down replicas, they send a `SIGTERM` signal to running containers. If an application exits immediately (`os.Exit(0)`), in-flight customer checkout requests are dropped mid-flight, resulting in HTTP 502 Bad Gateway errors. Graceful shutdown implements a deterministic sequence:
1. **Signal Interception**: Capture `SIGTERM` without terminating the process.
2. **Readiness Traffic Shedding**: Flip readiness checks to return HTTP 503 Service Unavailable, signaling to load balancers to stop routing new traffic to this instance.
3. **Connection Draining**: Call `server.Shutdown(ctx)` with a deadline context (e.g. 5–10 seconds), closing listeners while allowing active requests to complete.
4. **Resource Cleanup**: Close database connection pools, background worker channels, and cache flushers before final exit.

## Visual Model

```text
Orchestrator (SIGTERM / SIGINT)
  │
  ▼
┌─────────────────────────────────────────────────────────────┐
│ signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)      │
└──────────────────────────────┬──────────────────────────────┘
                               │ Signal Received
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ Step 1: Traffic Shedding (sc.shuttingDown.Store(true))       │
│   - healthMiddleware returns HTTP 503                       │
│   - Sets "Connection: close" header                         │
│   - Load balancers remove pod endpoint from iptables        │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ Step 2: sc.server.Shutdown(shutdownCtx)                     │
│   - Closes port 8080 listener (no new TCP accepts)          │
│   - Waits for active in-flight HTTP requests to complete    │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ Step 3: Cleanup Background Workers (close(sc.workerQueue))  │
│ Clean Exit (code 0) within timeout                          │
└─────────────────────────────────────────────────────────────┘
```

## Machine View

Under the hood:
1. `signal.Notify` registers a buffered channel with the Go runtime signal delivery mechanism. When the operating system delivers `SIGINT` (Ctrl+C) or `SIGTERM` (Kubernetes pod eviction), the signal is written to `stopChan` without terminating the process.
2. `server.Shutdown(ctx)` first closes all active listener sockets, ensuring no new TCP SYN packets are accepted. It then closes idle keep-alive connections and waits for active connections to finish their responses.
3. If active requests do not complete before `shutdownCtx` expires (e.g. within 5 seconds), `server.Shutdown` returns `ctx.Err()`, allowing the operator to force termination rather than hanging indefinitely.
4. In `main.go`, `http.ErrServerClosed` is ignored because it indicates that `Shutdown()` closed the listener as commanded.

## Run Instructions

```bash
# Run graceful shutdown test suite
go test -v ./curriculum/modules/18-flagship-opslane/assessments/13-opslane-graceful-shutdown

# Launch server and test SIGINT
go run ./curriculum/modules/18-flagship-opslane/assessments/13-opslane-graceful-shutdown
```

## Code Walkthrough

1. **`ShutdownCoordinator` Struct (lines 17–23)**: Holds `http.Server`, atomic `shuttingDown` flag, worker channel, and synchronization primitives.
2. **`healthMiddleware` (lines 43–52)**: If `shuttingDown.Load()` is true, returns HTTP 503 Service Unavailable and `Connection: close` header.
3. **Test Routes (lines 54–64)**: `/ping` returns 200 immediately; `/slow-task` simulates a 50ms in-flight request.
4. **`GracefulShutdown()` (lines 67–80)**:
   - Sets `shuttingDown` flag.
   - Invokes `server.Shutdown(ctx)` to drain HTTP requests.
   - Closes background worker queue.
5. **`main()` Signal Loop (lines 82–110)**: Starts server in a goroutine, waits on `SIGTERM`/`SIGINT`, and initiates `GracefulShutdown` with a 5-second timeout.

## Try It

1. Run the test suite: verify that in-flight `/slow-task` requests complete successfully even while `GracefulShutdown` is initiated.
2. Verify that requests sent *after* shutdown begins receive HTTP 503 Service Unavailable.
3. Run the application in a terminal, press `Ctrl+C`, and observe the clean 3-step draining sequence.

Common mistakes:
- Using `os.Exit(0)` directly in signal handlers, dropping active in-flight requests and severing database connections.
- Forgetting `Connection: close` during shutdown, allowing HTTP/1.1 keep-alive connections to keep sending requests to a terminating instance.
- Omitting a deadline context on `server.Shutdown`, allowing hung client connections to prevent the process from ever exiting.

## In Production

In Kubernetes, pod termination involves a race condition between the kubelet sending `SIGTERM` to the container and kube-proxy removing the pod's IP from node `iptables`. Adding a 3–5 second readiness failure delay before calling `server.Shutdown` ensures that external load balancers remove the pod from service endpoints before the server stops accepting TCP connections.

## Thinking Questions

1. Why does an abrupt process exit during rolling deployment cause HTTP 502 Bad Gateway errors for clients?
2. What is the race condition between Kubernetes `iptables` endpoint removal and pod `SIGTERM` reception?
3. Why does `server.Shutdown` require a `context.Context` parameter?

## Next Step

Proceed to [Milestone 14 — Deployment](../14-opslane-deployment/README.md) to containerize Opslane with Docker multi-stage builds and docker-compose orchestration.
