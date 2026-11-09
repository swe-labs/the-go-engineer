# Opslane Milestone 09 Assessment — Rubric

## Assessment Dimensions

### 1. Pub-Sub Architecture & Dispatch Concurrency (30%)
- **Excellent (5 pts)**: Non-blocking asynchronous dispatch; topic filtering and wildcard (`"*"`) support; lock released before invoking callbacks; handlers isolated in goroutines.
- **Proficient (3 pts)**: Pub-sub functions, but callbacks executed while holding mutex lock or wildcard matching omitted.
- **Developing (1 pt)**: Callbacks executed synchronously, blocking the caller thread.
- **Unacceptable (0 pts)**: Event bus fails to dispatch events.

### 2. Panic Isolation & Resilience (25%)
- **Excellent (5 pts)**: Every spawned handler goroutine contains a deferred panic recovery block; panics logged without crashing process; remaining listeners execute unaffected.
- **Proficient (3 pts)**: Panic recovery implemented, but panics are silenced without diagnostic logging.
- **Developing (1 pt)**: Panic in one listener crashes the entire application.
- **Unacceptable (0 pts)**: No panic handling.

### 3. Graceful Draining & Lifecycle Management (25%)
- **Excellent (5 pts)**: `Close` rejects future publishes; `sync.WaitGroup` tracks in-flight handlers; context timeout prevents hanging shutdown; subscription removal cleanly filters slices.
- **Proficient (3 pts)**: `Close` waits for handlers, but lacks context cancellation timeout.
- **Developing (1 pt)**: `Close` abruptly aborts running handlers without waiting.
- **Unacceptable (0 pts)**: No draining mechanism.

### 4. Conceptual & Distributed Messaging Defense (20%)
- **Excellent (5 pts)**: Accurately explains lock contention avoidance, goroutine leak prevention, `sync.WaitGroup` mechanics, and Transactional Outbox pattern.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in outbox or waitgroup internals.
- **Developing (1 pt)**: Cannot articulate why calling handlers under lock causes deadlocks.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
