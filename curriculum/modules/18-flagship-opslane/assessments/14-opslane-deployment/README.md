# Opslane deployment

## Mission

Configure cloud-native deployment health probes and container orchestration contracts for Opslane, implementing independent `/health/live` and `/health/ready` endpoints, parallel dependency checking, 50ms deadline enforcement, and degraded state signaling.

## Prerequisites

- `opslane-13`

## Mental Model

Cloud orchestrators (such as Kubernetes, AWS ECS, or Nomad) rely on distinct health probe contracts to manage container lifecycles:
1. **Liveness Probe (`/health/live`)**: Answers "Is this container process alive and un-deadlocked?" It must NEVER check downstream dependencies (databases, external APIs). If a liveness probe fails, the orchestrator immediately restarts the container.
2. **Readiness Probe (`/health/ready`)**: Answers "Is this instance capable of serving customer traffic right now?" It checks downstream dependencies (PostgreSQL, Redis cache, background worker pools) in parallel. If any dependency is degraded or unreachable, it returns HTTP 503 Service Unavailable. The orchestrator responds by temporarily removing the pod from load balancer endpoints, without killing or restarting the container process.

## Visual Model

```text
Kubernetes Kubelet Health Probes
  │
  ├──► GET /health/live (Every 10s)
  │    └── Return HTTP 200 {"status":"ALIVE"}
  │        (Never restarts unless Go runtime deadlocks)
  │
  └──► GET /health/ready (Every 5s)
       │
       ▼
┌─────────────────────────────────────────────────────────────┐
│ DeploymentHealthService.handleReadiness                     │
│   Timeout Context: 50ms                                     │
│   Execute Component Checks in Parallel (sync.WaitGroup)     │
├──────────────────────────────┬──────────────────────────────┤
│ PostgreSQL Check             │ Redis Cache Check            │
│ Worker Pool Check            │ External Gateway Check       │
└──────────────────────────────┴──────────────────────────────┘
       │
       ├── All Healthy ────────► HTTP 200 OK {"status":"UP"}
       │                        (Load balancer routes traffic)
       └── Any Component Fails ─► HTTP 503 {"status":"DEGRADED"}
                                (Load balancer stops routing traffic)
```

## Machine View

Under the hood:
1. `/health/live` is lightweight and non-blocking, returning HTTP 200 immediately to prevent accidental restart loops caused by temporary network blips.
2. `/health/ready` bounds all component health evaluations to a 50ms window via `context.WithTimeout(r.Context(), 50*time.Millisecond)`. Slow dependencies cannot cause readiness probes to hang.
3. Component checks execute concurrently using goroutines and `sync.WaitGroup`, collecting `ComponentHealth` structs into a preallocated slice `results[idx]` indexed by check number, avoiding slice append race conditions.
4. If `allHealthy` evaluates to `false`, `w.WriteHeader(http.StatusServiceUnavailable)` is emitted alongside latency measurements in milliseconds.

## Run Instructions

```bash
# Run deployment probe test suite
go test -v ./curriculum/modules/18-flagship-opslane/assessments/14-opslane-deployment

# Launch deployment health server
go run ./curriculum/modules/18-flagship-opslane/assessments/14-opslane-deployment
```

## Code Walkthrough

1. **Health Models (lines 12–25)**: Defines `ComponentHealth` (name, healthy flag, latency) and `ReadinessReport` (status, duration, component array).
2. **`HealthChecker` & Generic Checker (lines 27–55)**: Interface defining `CheckHealth(ctx)`; generic checker supports error and latency simulation.
3. **`DeploymentHealthService` (lines 58–73)**: Aggregates checkers and mounts `/health/live` and `/health/ready`.
4. **`handleLiveness` (lines 75–80)**: Emits instantaneous HTTP 200 `{"status":"ALIVE"}`.
5. **`handleReadiness` (lines 83–129)**:
   - Sets 50ms context timeout.
   - Spawns concurrent check goroutines with `sync.WaitGroup`.
   - Gathers results; determines overall health status.
   - Returns HTTP 200 `"UP"` if all pass, or HTTP 503 `"DEGRADED"` if any fail.

## Try It

1. Run the test suite: verify that all healthy components yield HTTP 200 with latency measurements.
2. Inject a failing database checker: verify that `/health/ready` returns HTTP 503 `"DEGRADED"`, while `/health/live` continues returning HTTP 200 `"ALIVE"`.
3. Simulate a slow checker (>50ms): verify that the context timeout halts evaluation and reports the failure.

Common mistakes:
- Checking database connectivity inside the liveness probe, causing Kubernetes to enter cascading restart loops during transient database connection spikes.
- Executing component checks sequentially, causing readiness probes to exceed orchestrator timeouts.
- Not specifying timeout contexts on health checks, allowing hanging sockets to freeze the health server.

## In Production

In Kubernetes manifests, liveness and readiness probes are declared with `periodSeconds: 10`, `timeoutSeconds: 2`, and `failureThreshold: 3`. Setting `initialDelaySeconds: 15` allows application startup and migration phases to complete before the kubelet begins polling probes.

## Thinking Questions

1. Why does checking a third-party database inside `/health/live` cause a cascading cluster outage when the database slows down?
2. Why must component health checks be executed in parallel rather than sequentially?
3. How does Kubernetes treat a pod whose readiness probe returns HTTP 503 versus one whose liveness probe returns HTTP 503?

## Next Step

Proceed to [Milestone 15 — Portfolio Defense](../projects/15-opslane-portfolio-defense/README.md) to rehearse technical interview defense questions and articulate Opslane's production architectural guarantees.
