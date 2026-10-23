# Opslane Milestone 11 Assessment — Rubric

## Assessment Dimensions

### 1. Singleflight & Stampede Defense (30%)
- **Excellent (5 pts)**: `singleflight.Group` correctly wraps cache miss path; concurrent requests coalesced to single database read; double-check pattern verified; non-blocking broadcast to all callers.
- **Proficient (3 pts)**: Singleflight used, but double-check inside barrier omitted or key partitioning is coarse.
- **Developing (1 pt)**: Singleflight missing, allowing concurrent cache stampedes to hit database.
- **Unacceptable (0 pts)**: Concurrency causes race crashes.

### 2. TTL Expiration & Thread-Safe Cache Design (25%)
- **Excellent (5 pts)**: `TTLCache` uses `sync.RWMutex` with concurrent read scaling; expiration checked against `time.Now()`; safe serialization into byte slices prevents pointer aliasing.
- **Proficient (3 pts)**: Cache implemented with mutex, but stores raw mutable pointers instead of serialized bytes.
- **Developing (1 pt)**: Cache lacks read lock optimization or expiration logic is flawed.
- **Unacceptable (0 pts)**: Map accessed without lock protection.

### 3. Invalidation & Mutation Consistency (25%)
- **Excellent (5 pts)**: Cache-Aside pattern strictly followed; database updated before cache invalidation (`cache.Delete`); zero stale read window; atomic counters verify DB query counts.
- **Proficient (3 pts)**: Cache invalidated on update, but done before DB write, risking stale cache repopulation.
- **Developing (1 pt)**: Cache never invalidated on update, serving stale data indefinitely.
- **Unacceptable (0 pts)**: Invalidation logic broken.

### 4. Conceptual & High-Scale Caching Defense (20%)
- **Excellent (5 pts)**: Accurately explains cache stampedes, cache penetration, passive vs active TTL cleanup, pointer aliasing risks, and singleflight internals.
- **Proficient (3 pts)**: Answers conceptual questions accurately with minor gaps in LRU eviction or pointer aliasing.
- **Developing (1 pt)**: Cannot articulate why singleflight is needed or how stampedes overwhelm databases.
- **Unacceptable (0 pts)**: No explanation provided.

## Grading Thresholds
- **12–14 pts**: Exceeds expectations (Gold standard)
- **10–11 pts**: Meets expectations (Passing)
- **7–9 pts**: Needs revision
- **0–6 pts**: Unsatisfactory — review prerequisites and retake
