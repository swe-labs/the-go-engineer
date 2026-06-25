# Health and readiness probes

## Learning objective

Implement `/healthz` (liveness) and `/readyz` (readiness) HTTP endpoints, distinguish between liveness and readiness semantics, check dependency status, and configure Kubernetes probes.

## Why this matters

Kubernetes and cloud orchestrators use HTTP probes to decide whether a pod is alive and ready to serve traffic. A misconfigured probe causes cascading failures: the orchestrator kills healthy pods or sends traffic to broken ones. Every production service must expose correct liveness and readiness endpoints.

## Mental model

Think of liveness vs readiness like a chef in a kitchen:

- **Liveness** (`/healthz`): Is the chef alive? Are they breathing? If not, replace them (restart the pod).
- **Readiness** (`/readyz`): Is the chef ready to cook? Do they have ingredients (database), a clean station (cache connection), and sharp knives (dependencies)? If not, don't send orders (traffic) yet.

Liveness says "don't kill me". Readiness says "send me traffic". A pod can be alive but not ready (e.g., warming cache after startup). A pod can be ready but not alive (shouldn't happen -- readiness implies liveness).

## Core idea

Expose two endpoints:

| Endpoint | Purpose | Return 200 when | Return 503 when |
|---|---|---|---|
| `/healthz` | Liveness | Process is running, basic health check passes | Deadlocked, OOM, panic loop |
| `/readyz` | Readiness | All dependencies (DB, cache, etc.) are available | A dependency is down, cache is cold |

Implementation pattern:

```go
func healthzHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if !isHealthy() {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy"})
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
```

## Under the hood

Kubernetes probes work by polling:

1. The kubelet sends an HTTP GET to the probe endpoint at the configured interval (default: 10s).
2. If the endpoint returns 2xx-3xx, the probe is successful.
3. If the endpoint returns 4xx-5xx (or fails to connect), the probe fails.
4. After `failureThreshold` consecutive failures (default: 3), the action is taken:
   - Liveness failure: kubelet restarts the container.
   - Readiness failure: kubelet removes the pod from Service endpoints.
5. `initialDelaySeconds` skips probing for a period after startup.
6. `periodSeconds` controls the polling interval.

At the Go level, each probe endpoint should be lightweight (no complex computation) and should not depend on the very resource it's checking causing cascading failures.

## How Go uses it

The Go community has standard patterns for probes. `chi` and `gin` both have middleware for health checks. The `health` package in the Go ecosystem provides composable checkers. Kubernetes itself is written in Go and uses HTTP probes internally.

Many production Go services expose three endpoints:
- `/healthz` -- simple, always returns 200 if the process is running.
- `/readyz` -- checks database, cache, and upstream dependencies.
- `/livez` -- sometimes separate from healthz for custom liveness logic.

## Go example

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type probeState struct {
	mu         sync.RWMutex
	dbReady    bool
	cacheReady bool
}

var state = &probeState{dbReady: true, cacheReady: true}

type healthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	state.mu.RLock()
	healthy := state.dbReady && state.cacheReady
	state.mu.RUnlock()

	if !healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(healthResponse{
			Status: "unhealthy",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthResponse{
		Status: "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	state.mu.RLock()
	ready := state.dbReady
	state.mu.RUnlock()

	if !ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(healthResponse{
			Status: "not ready",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthResponse{
		Status: "ready",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/readyz", readyzHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Step-by-step execution

When Kubernetes starts a pod:

1. Pod is scheduled, container starts, Go binary runs.
2. `main()` configures probes and starts the HTTP server.
3. After `initialDelaySeconds`, kubelet sends GET `/healthz`. The handler checks dependency flags, returns 200, `{"status": "ok"}`.
4. kubelet sends GET `/readyz`. The handler checks `dbReady` flag (true) and returns 200, `{"status": "ready"}`.
5. Pod is added to Service endpoints. Traffic begins flowing.
6. The database goes down. Application sets `state.dbReady = false`.
7. Next `/readyz` probe returns 503. kubelet removes pod from Service endpoints. No new traffic arrives.
8. `/healthz` still returns 200 (the app process is alive). kubelet does NOT restart the pod.
9. Database recovers. Application sets `state.dbReady = true`.
10. Next `/readyz` returns 200. Pod is re-added to Service endpoints.

## Common mistakes

- Using the same logic for liveness and readiness. A DB failure should make readiness fail but NOT liveness -- otherwise Kubernetes restarts a perfectly healthy process that just can't reach the DB.
- Making probes expensive. A probe that queries every row of a database table will slow down under load, causing cascade failures.
- Not implementing probes at all. Without probes, Kubernetes cannot detect or recover from failures.
- Probes depending on the probe infrastructure. A probe that logs to a full disk will fail because the disk is full, creating a self-fulfilling restart loop.
- Not using `ReadHeaderTimeout` or `ReadTimeout` on the probe server. A stuck probe handler can block the entire server.

## Debugging walkthrough

A pod that keeps restarting:

```
Events:
  Warning Unhealthy  3s  kubelet  Liveness probe failed: HTTP probe failed with statuscode: 503
```

**Symptom**: Pod is stuck in `CrashLoopBackOff`.

**Investigation**: Check the liveness endpoint directly: `kubectl exec <pod> -- curl localhost:8080/healthz`. It returns 503. The handler checks `isHealthy()` which checks `dbReady && cacheReady`. The database is down, but the liveness probe shouldn't depend on the database.

**Fix**: Separate liveness from readiness:

```go
// Liveness: only checks process health
func healthzHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Readiness: checks dependencies
func readyzHandler(w http.ResponseWriter, r *http.Request) {
    if !dbReady {
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
}
```

## Production notes

- Liveness should be a simple `200 OK` response, optionally checking process-level health (goroutine count, memory usage).
- Readiness should check all external dependencies: database, cache, message queue, upstream APIs.
- Use `initialDelaySeconds: 5` to allow the application to start before probing.
- Use `failureThreshold: 3` to tolerate transient failures.
- For liveness, set `periodSeconds` relatively high (15-30s) since process death is rare.
- For readiness, set `periodSeconds` lower (5-10s) to respond quickly to dependency failures.
- Add a `/livez` endpoint for custom liveness checks if needed.

## Performance implications

- Probes are polled every few seconds by the kubelet. Each probe is a TCP connection + HTTP request/response cycle. Minimal overhead.
- Expensive probes (full DB queries) can overload dependencies during rolling restarts when many pods probe simultaneously.
- Liveness probes should complete in <1ms. Readiness probes should complete in <100ms.
- Consider caching readiness check results for a short period (a few seconds) to avoid hammering dependencies.

## Practice task

1. Write `/healthz` and `/readyz` handlers for a service with two dependencies: `db` and `cache`.
2. Use a `sync.RWMutex`-protected state that tracks each dependency's status.
3. Write table-driven tests for: both dependencies healthy, DB down, cache down, both down.
4. Ensure liveness returns 200 even when dependencies are down.
5. Write a test that simulates a dependency going down and verifies the readiness probe transitions correctly.

## Tests / verification

```bash
go test ./curriculum/modules/08-http-rest-apis/lessons/16-health-and-readiness-probes -v
```

Tests cover healthy/unhealthy states for both liveness and readiness, method validation, and the `livenessHandler` factory.

## Review questions

1. What is the difference between liveness and readiness probes in Kubernetes?
2. Why should a liveness probe NOT check database connectivity?
3. What HTTP status code should a probe return when the application is not ready?
4. How does Kubernetes use the readiness probe to manage traffic routing?
5. What are `initialDelaySeconds`, `periodSeconds`, and `failureThreshold` in the Kubernetes probe configuration?

## NEXT UP

REST design principles -- resource-oriented URLs, statelessness, HATEOAS, idempotency, and the REST maturity model.
