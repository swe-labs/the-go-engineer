# Advanced Failure Scenarios

This document is a reference for later-stage lessons that intentionally teach diagnosis and recovery.

It is not a foundations template.

## 1. Concurrency Failures

Use in `04-concurrency`.

```markdown
### Scenario: Goroutine Leak in Background Worker

**Setup:**
```go
for {
    task := getTask()
    go process(task)
}
```

**Expected:** work continues normally

**Actual:** goroutine count grows until memory pressure or collapse

**Root Cause:** there is no cancellation or bound on worker creation

**Solution Direction:**
- propagate context
- stop on cancellation
- use bounded workers

**Prevention:**
- track goroutine count
- review all long-lived loops for exit paths
```

```markdown
### Scenario: Data Race in Shared Counter

**Setup:**
```go
counter++
```

**Expected:** deterministic count

**Actual:** non-deterministic results under concurrent access

**Solution Direction:**
- mutex
- atomic
- message passing

**Prevention:**
- run with `-race`
- review all shared mutable state
```

## 2. Backend Failures

Use in `03-backend-systems` and `06-production`.

```markdown
### Scenario: Slow Request Body Exhausts Memory

**Setup:** request body is read without limits

**Expected:** request completes normally

**Actual:** server memory spikes or process crashes

**Solution Direction:**
- bound body size
- enforce timeouts
- reject oversized payloads early
```

```markdown
### Scenario: Connection Pool Exhaustion

**Setup:** too many concurrent database users without tuned pool limits

**Expected:** steady throughput

**Actual:** blocked requests, timeouts, or dependency collapse

**Solution Direction:**
- configure connection limits
- add timeouts
- monitor pool metrics
- reduce unnecessary concurrent work
```

## 3. Performance and Memory Failures

Use in `07-advanced`.

```markdown
### Scenario: Retained Backing Array

**Setup:** a tiny slice view keeps a large backing array alive

**Expected:** memory drops after processing

**Actual:** heap stays large

**Solution Direction:**
- copy required data out
- redesign buffer lifecycle
- measure before and after with heap profiling
```

```markdown
### Scenario: Allocation Regression in Hot Path

**Setup:** a loop allocates every iteration

**Expected:** stable throughput

**Actual:** latency and GC pressure climb

**Solution Direction:**
- profile first
- identify avoidable allocations
- compare clarity against optimization cost
```

## Reusable Shape

```markdown
### Scenario: Name

**Setup:**
code or operating context

**Expected:**
healthy behavior

**Actual:**
broken behavior

**Root Cause:**
the real mechanism

**Solution Direction:**
- fix 1
- fix 2

**Prevention:**
- habit 1
- signal 2
```

## Rule

Only add these scenarios when the learner can already:

- read the code surface confidently
- explain the happy path
- benefit from a diagnosis task instead of just copying the solution
