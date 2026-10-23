# Opslane Milestone 14 Assessment — Rubric

## Assessment Dimensions

### 1. Probe Separation & Orchestration Semantics (30%)
- **Excellent (5 pts)**: Complete decoupling of `/health/live` (lightweight, zero external dependencies) from `/health/ready` (comprehensive dependency evaluation); HTTP 200 vs 503 status codes conform to orchestrator specs.
- **Proficient (3 pts)**: Probes separated, but liveness executes light dependency queries or readiness returns 500 instead of 503.
- **Developing (1 pt)**: Single health endpoint combines liveness and readiness indiscriminately.
- **Unacceptable (0 pts)**: No health endpoints implemented.

### 2. Concurrent Execution & Timeout Hardening (25%)
- **Excellent (5 pts)**: Component checks executed concurrently via `sync.WaitGroup`; 50ms context deadline enforced; preallocated slice index assignment avoids data races.
- **Proficient (3 pts)**: Checks executed concurrently, but timeout is missing or uses sleep.
- **Developing (1 pt)**: Checks executed sequentially, accumulating latency and exceeding probe timeouts.
- **Unacceptable (0 pts)**: Health checks deadlock or timeout under test.

### 3. Diagnostic Reporting & Breakdown Details (25%)
- **Excellent (5 pts)**: Clean JSON report with overall status (`UP` / `DEGRADED`), total duration in milliseconds, and individual component health breakdown with latency measurements.
- **Proficient (3 pts)**: Status returned, but latency measurements or component breakdown omitted.
- **Developing (1 pt)**: Raw string responses without component breakdown.
- **Unacceptable (0 pts)**: Unparseable responses.

### 4. Conceptual & Cloud-Native Defense (20%)
- **Excellent (5 pts)**: Accurately explains cascading failure mechanics, rolling deployment zero-downtime guarantees, Kubernetes probe tuning parameters, and thread-safe index writes.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in probe timing parameters.
- **Developing (1 pt)**: Confuses liveness with readiness or cannot explain why databases must not be checked in liveness.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
