# Opslane Milestone 10 Assessment — Rubric

## Assessment Dimensions

### 1. Bounded Concurrency & Queue Backpressure (30%)
- **Excellent (5 pts)**: Fixed worker ceiling strictly maintained; buffer capacity bounded; submit operations enforce backpressure; peak concurrency never exceeds configured workers.
- **Proficient (3 pts)**: Concurrency bounded, but queue uses unbounded capacity or peak concurrency fluctuates slightly above limit.
- **Developing (1 pt)**: Spawns unconstrained goroutines per job, failing under heavy load.
- **Unacceptable (0 pts)**: Concurrency model deadlocks or crashes.

### 2. Lock-Free Telemetry & Panic Isolation (25%)
- **Excellent (5 pts)**: Active, max, and completed jobs tracked using lock-free atomics (`sync/atomic`); CAS loop updates `maxConcurrencySeen` safely; panic recovery preserves telemetry counters.
- **Proficient (3 pts)**: Atomic counters used, but CAS loop missing or panic recovery omits counter decrements.
- **Developing (1 pt)**: Counters updated without atomic primitives, causing race conditions under `go test -race`.
- **Unacceptable (0 pts)**: Panicking job crashes entire worker pool.

### 3. Graceful Draining & Shutdown Lifecycle (25%)
- **Excellent (5 pts)**: `Shutdown` prevents new submissions via `closed.Swap(true)`; channel cleanly closed; workers drain remaining buffered jobs; `wg.Wait()` bounded by context timeout.
- **Proficient (3 pts)**: Workers drain jobs, but `Shutdown` lacks timeout protection.
- **Developing (1 pt)**: Closing channel while producers are active causes panics (`send on closed channel`).
- **Unacceptable (0 pts)**: Shutdown hangs indefinitely or drops in-flight jobs.

### 4. Conceptual & Runtime Systems Defense (20%)
- **Excellent (5 pts)**: Accurately explains backpressure mechanics, CAS loop semantics, channel closing rules, and OOM prevention strategies.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in CAS or channel close panics.
- **Developing (1 pt)**: Cannot articulate what backpressure means or why unbounded goroutines cause crashes.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
