# Opslane Milestone 13 Assessment — Rubric

## Assessment Dimensions

### 1. Signal Handling & Phased Coordination (30%)
- **Excellent (5 pts)**: `signal.Notify` cleanly captures `SIGINT` and `SIGTERM`; buffered channel prevents missed signals; 3-step sequence (traffic shedding -> server shutdown -> worker queue closure) strictly ordered.
- **Proficient (3 pts)**: Signals captured, but shutdown sequence executes out of order (e.g. closing worker queues before draining HTTP server).
- **Developing (1 pt)**: Signal handling exits immediately (`os.Exit`), dropping active requests.
- **Unacceptable (0 pts)**: Signals not captured; process terminates abruptly.

### 2. Traffic Shedding & In-Flight Draining (25%)
- **Excellent (5 pts)**: Incoming requests shed with HTTP 503 and `Connection: close`; in-flight requests complete cleanly without errors; `server.Shutdown` executes within deadline context.
- **Proficient (3 pts)**: Server drains requests, but fails to set `Connection: close` header on 503 responses.
- **Developing (1 pt)**: Server refuses new requests but aborts in-flight requests immediately.
- **Unacceptable (0 pts)**: In-flight requests severed or dropped.

### 3. Timeout Bounds & Error Semantics (25%)
- **Excellent (5 pts)**: `context.WithTimeout` enforces hard termination bound; `http.ErrServerClosed` correctly ignored in `main()`; shutdown errors bubble up to logging.
- **Proficient (3 pts)**: Timeout configured, but `http.ErrServerClosed` is logged as a fatal error.
- **Developing (1 pt)**: Uses unbounded context, risking hanging indefinitely on stuck client connections.
- **Unacceptable (0 pts)**: Server hangs indefinitely on shutdown.

### 4. Conceptual & Orchestration Defense (20%)
- **Excellent (5 pts)**: Accurately explains Kubernetes `iptables` race conditions, HTTP 502 root causes, SIGKILL vs SIGTERM mechanics, and proper resource teardown ordering.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in `iptables` or signal buffer mechanics.
- **Developing (1 pt)**: Cannot articulate why database connections must be closed last.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
