# Opslane Milestone 08 Assessment — Rubric

## Assessment Dimensions

### 1. Idempotency State Machine & Deduplication (30%)
- **Excellent (5 pts)**: Strict state tracking (`IN_PROGRESS`, `COMPLETED`); concurrent in-flight duplicates rejected with `ErrDuplicateInFlight`; completed requests replay cached receipt without touching gateway; reservations cleaned up on abort.
- **Proficient (3 pts)**: Idempotency implemented, but concurrent in-flight requests are not blocked or receipts are un-cached.
- **Developing (1 pt)**: Replaying requests executes second charge on gateway.
- **Unacceptable (0 pts)**: No idempotency logic.

### 2. Fault-Tolerant Retry & Backoff Mechanics (25%)
- **Excellent (5 pts)**: Exponential backoff retries implemented; strictly filters transient vs permanent errors; respects context cancellation (`ctx.Done()`) during backoff; terminates on permanent failure.
- **Proficient (3 pts)**: Retries implemented, but retries permanent errors (e.g. card declines) or lacks exponential backoff.
- **Developing (1 pt)**: Retries infinitely or ignores context timeouts.
- **Unacceptable (0 pts)**: No retry capability on transient network failure.

### 3. Thread Safety & Resource Cleanup (25%)
- **Excellent (5 pts)**: All map mutations protected via `sync.Mutex`; keys cleanly deleted on context cancellation or terminal error; atomic testing simulator accurately tracks call counts.
- **Proficient (3 pts)**: Mutex used, but key deletion omitted on some error return paths.
- **Developing (1 pt)**: Map accessed without lock protection, failing under `-race`.
- **Unacceptable (0 pts)**: Concurrency race crashes.

### 4. Conceptual & Financial Systems Defense (20%)
- **Excellent (5 pts)**: Accurately explains transient vs permanent errors, retry storms, randomized jitter, distributed idempotency persistence, and financial race condition hazards.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in jitter mechanics or distributed stores.
- **Developing (1 pt)**: Cannot articulate why card declines should not be retried.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
