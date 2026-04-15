# Advanced Production Notes

This document is a reference for later-stage lessons that need production-hardening guidance.

It is not a foundations template.

## 1. Backend Systems

Use in `03-backend-systems`.

```markdown
## Production Notes

### HTTP Timeouts
- read timeout
- write timeout
- idle timeout
- header timeout

### Common Issues
1. Slow header or body attacks
   - Fix: set header/read timeouts

2. Large request bodies
   - Fix: bound body size before reading

3. Handler-spawned background work
   - Fix: use context, cancellation, and limits
```

## 2. Database and Stateful Services

Use in `03-backend-systems` and `06-production`.

```markdown
## Production Notes

### Connection Pool Management
- max open connections
- max idle connections
- connection lifetime
- idle lifetime

### Common Issues
1. Pool exhaustion
   - Symptom: requests block or time out
   - Fix: set limits, use timeouts, monitor pool metrics

2. Deadlocks
   - Fix: consistent lock order and shorter transactions

3. Slow queries
   - Fix: inspect plans and indexes
```

## 3. Concurrency

Use in `04-concurrency` and `06-production`.

```markdown
## Production Notes

### Common Issues
1. Goroutine leaks
   - Fix: cancellation, worker limits, shutdown discipline

2. Race conditions
   - Fix: `-race`, mutex, channels, or atomics where appropriate

3. Backpressure collapse
   - Fix: bounded queues, semaphores, admission control
```

## 4. Memory and Performance

Use in `07-advanced`.

```markdown
## Production Notes

### Common Issues
1. Unbounded maps and caches
2. Retained backing arrays
3. Allocation-heavy hot paths

### Debugging Tools
- pprof for CPU
- pprof for heap
- trace for scheduler and blocking behavior
```

## 5. Shutdown and Observability

Use in `06-production` and `08-projects`.

```markdown
## Production Notes

### Shutdown
1. Stop taking new work
2. Drain in-flight work
3. Close stateful dependencies
4. Flush logs and exit cleanly

### Observability
- structured logs
- request or correlation IDs
- latency percentiles
- error rate
- saturation metrics
```

## Reusable Shape

```markdown
## Production Notes

### Critical Configuration
- item 1
- item 2

### Common Failure Modes
1. Issue
   - Symptom
   - Fix

### Monitoring Checklist
- metric 1
- metric 2

### Debugging Commands
```bash
# command 1
# command 2
```
```

## Rule

Only add this layer when the lesson already has:

- a clear local mental model
- a working happy path
- enough earned context for operational consequences to make sense
