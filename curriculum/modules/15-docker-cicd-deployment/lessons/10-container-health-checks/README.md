# Container health checks

## Learning objective

Implement `/healthz` and `/readyz` HTTP endpoints in Go, configure Docker `HEALTHCHECK` and Kubernetes readiness/liveness probes, and distinguish between application startup, liveness, and readiness states.

## Why this matters

Without health checks, Kubernetes thinks your pod is healthy even when your Go service is deadlocked, has exhausted its connection pool, or is still starting up. Traffic is routed to a broken pod. Requests fail. Users see errors. Health checks are the mechanism by which the orchestrator knows whether your service can handle traffic. Implementing them correctly is essential for zero-downtime deployments.

## Mental model

Health checks are the orchestrator asking your service: "Are you alive?" (liveness), "Can you handle traffic?" (readiness), and "Have you finished starting up?" (startup). Think of them as three distinct questions:
- **Liveness**: Is the Go process still running and not deadlocked? If no, restart the container.
- **Readiness**: Are all dependencies connected and warm? If no, remove from load balancer.
- **Startup**: Has the service finished initializing? If no, delay liveness checks to avoid premature restarts.

## Core idea

| Probe type | Purpose | Failure action |
|---|---|---|
| Startup | Was initialization successful? | Restart container |
| Liveness | Is the process alive? | Restart container |
| Readiness | Can the process serve traffic? | Remove from service |

In Docker, only `HEALTHCHECK` exists — a single command that Docker runs periodically. It distinguishes `starting`, `healthy`, and `unhealthy`. In Kubernetes, all three probe types exist as separate configurable checks.

## Under the hood

Docker `HEALTHCHECK` runs a command inside the container every N seconds. The exit code determines the health state: 0 = healthy, 1 = unhealthy. After a configurable number of consecutive failures, the container is marked unhealthy and Docker stops sending traffic to it (in Swarm mode) or restarts it (with `--restart=always`).

Kubernetes probes work differently. The kubelet on each node executes the probe against the container's IP:port. For HTTP probes, it sends GET requests to the specified path and port. A 2xx or 3xx response is a success. Anything else (including connection failure or timeout) is a failure. The kubelet updates the pod status, which the service controller uses to decide whether to route traffic.

## How Go uses it

All production Go HTTP servers should expose at least `/healthz` and `/readyz` endpoints:

```go
http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
})

http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
    if dbPool.Ping(r.Context()) != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
})
```

The `/healthz` endpoint is a light check that the process is alive — it should be extremely fast and never block. The `/readyz` endpoint checks dependencies: database connectivity, cache connectivity, and that any pre-warming is complete.

## Go example

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

type HealthStatus struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
}

var (
	healthy   atomic.Bool
	ready     atomic.Bool
	startTime time.Time
)

func init() {
	healthy.Store(false)
	ready.Store(false)
	startTime = time.Now()
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if !healthy.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(HealthStatus{
			Status:    "unhealthy",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Uptime:    time.Since(startTime).String(),
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    time.Since(startTime).String(),
	})
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	if !ready.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "not ready"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

func main() {
	healthy.Store(true)

	go func() {
		time.Sleep(2 * time.Second)
		ready.Store(true)
		slog.Info("service is ready")
	}()

	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/readyz", readyzHandler)

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
```

## Step-by-step execution

1. The Go service starts. `healthy` is set to true immediately.
2. `ready` starts as false — the service is alive but not yet ready to serve traffic.
3. Kubernetes sends a readiness probe to `/readyz`. The service returns 503 (Service Unavailable).
4. Kubernetes removes the pod from the service endpoints — no traffic is routed to it.
5. After 2 seconds, the goroutine sets `ready = true`.
6. The next readiness probe returns 200 (OK). Kubernetes adds the pod to the service.
7. Traffic is routed to the pod.
8. If the service later becomes unhealthy, the health check returns 503 and Kubernetes restarts the pod.

## Common mistakes

- **Using the same check for health and readiness.** Health should check process sanity. Readiness should check dependency connectivity. If you check the database in both, a transient DB failure causes a restart instead of a graceful traffic drain.
- **Making health checks depend on external services.** A database outage should mark the pod as not ready, not trigger a restart. Health checks should be fast, local, and cheap.
- **Not handling the case where the service is still starting.** Without a startup probe, Kubernetes starts sending traffic before the service is ready. Use a startup probe or an initial delay on the readiness probe.
- **Returning 200 even when dependencies are down.** The health endpoint always returns 200. The orchestrator sees "healthy" and keeps routing traffic to a broken pod.
- **Applying the same check to every instance.** In a multi-instance deployment, one instance with a slow goroutine could pass health checks while rejecting traffic.

## Debugging walkthrough

Consider a Go service that keeps getting restarted by Kubernetes even though the process is running:

```go
func healthzHandler(w http.ResponseWriter, r *http.Request) {
    if err := dbPool.Ping(r.Context()); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
}
```

**Symptom:** The pod restarts every few minutes. `kubectl describe pod` shows "unhealthy" probe failures.

**Investigation:** Check the kubelet logs and the service logs. The database had a brief network hiccup. The health check failed because the DB ping timed out. Kubernetes interpreted this as the service being dead and restarted the pod.

**Root cause:** The liveness probe checks the database. A transient DB failure triggers a restart instead of a traffic drain.

**Fix:** Move the DB check to the readiness probe, not the liveness probe:
```go
func readyzHandler(w http.ResponseWriter, r *http.Request) {
    if err := dbPool.Ping(r.Context()); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
}
```
The liveness probe should only check process-level health.

## Production notes

- Always set `timeoutSeconds` shorter than `periodSeconds` in Kubernetes probe config.
- Use a startup probe for services with long initialization (e.g., loading ML models, warming caches). Value: `failureThreshold * periodSeconds` should exceed the expected startup time.
- Monitor health check failures in your observability system. A pattern of readiness failures indicates a dependency problem.
- Log health check requests at DEBUG level only — logging every health check at INFO level creates noise and cost.
- Use separate ports for health checks in high-throughput services to avoid mutex contention with the main request handler.
- In Docker, set `HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD ["/healthcheck"]`.

## Performance implications

- Health check handlers should be extremely fast — typically under 1ms. A slow health check blocks the probe and triggers unnecessary restarts.
- Each health check creates a new HTTP request. At 10-second intervals across 100 pods, that is 10 health checks per second — negligible.
- Avoid mutex contention in health check handlers. Use `atomic` variables for simple state tracking.
- Database ping in readiness probes adds latency but is acceptable if the pool is configured with a short timeout. Use a 1-second timeout for health check pings.

## Practice task

Write a Go HTTP server with `/healthz` and `/readyz` endpoints. The server starts in "not ready" state and becomes ready after 3 seconds (simulating dependency initialization). Write a test that verifies the health endpoint returns 200 immediately but the ready endpoint returns 503 in the first 3 seconds.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/10-container-health-checks
go test -v ./curriculum/modules/15-docker-cicd-deployment/lessons/10-container-health-checks
```

## Review questions

1. What is the difference between a liveness probe and a readiness probe in Kubernetes?
2. Why should you avoid checking external dependencies in a liveness probe?
3. What does Docker do when a container's `HEALTHCHECK` fails three consecutive times?
4. How would you implement a startup probe equivalent in a Go HTTP server?
5. Why might you want health check endpoints on a separate port from the main application?

## NEXT UP

GitHub Actions — automating build, test, and deploy pipelines for Go projects with GitHub's CI/CD platform.
