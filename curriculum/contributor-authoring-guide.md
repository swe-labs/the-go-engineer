# Contributor Authoring Guide

**How to write zero-magic lessons using golden lesson exemplars.**

---

## Prerequisites

Before authoring, read (in this order):

1. `zero-magic-authoring-spec.md` — pedagogical contract and field definitions
2. `quality-rubric.md` — scoring dimensions and minimums
3. This guide — how to apply the standards using golden lesson examples

---

## Quick Reference

| Spec                 | Location                                                 |
| -------------------- | -------------------------------------------------------- |
| Pedagogical contract | `zero-magic-authoring-spec.md`                           |
| Quality rubric       | `quality-rubric.md`                                      |
| Golden lessons       | `path.core.json` items in `goldenLessons` set            |
| Validator            | `go run ./internal/tools/curriculum validate-zero-magic` |

---

## Authoring Workflow

```
1. Identify the confusion: what do learners get wrong?
2. Write the mental model first (the analogy)
3. Write under_the_hood (the mechanics)
4. Write execution_timeline or memory_timeline (what happens over time/in memory — check the spec's trigger tables for when each is required)
5. Write step_by_step_execution (the prose walkthrough)
6. Write failure_modes (what breaks and how)
7. Write beginner_mistakes (with mechanical root cause)
8. Write hidden_magic_checks (questions that force mechanical reasoning)
9. Write proof_of_understanding (how mastery is demonstrated)
10. Write problem_solved and why_it_exists (framing)
11. Self-review against the Quality Rubric
12. Run validate-zero-magic
```

---

## Field-by-Field Guide with Golden Examples

### 1. `mental_model`

**Purpose:** A durable analogy that maps cleanly to runtime behavior and survives edge cases.

**Bad (placeholder):**

```
"Think of Concept X as one step in the learner's path from raw code to professional engineering judgment."
```

**Good — slice len/cap (core-03-15):**

```
"A slice is a window into a backing array. Len is how many elements you can see through the window. Cap is
how many elements exist in the backing array starting from the window's position. Append moves the window —
and sometimes builds a bigger house behind it."
```

**Good — request body (core-08-05):**

```
"The request body is a stream, not a buffer. It can be read exactly once, from beginning to end. Reading the
body is like drinking from a straw — once consumed, it is gone."
```

**Good — goroutines (core-11-11):**

```
"A goroutine is a lightweight thread managed by the Go runtime. Unlike OS threads (~1MB stack, expensive
context switch), a goroutine starts with a tiny stack (~4KB, now ~2KB minimum) and is multiplexed onto OS
threads by the Go scheduler. The go keyword forks the current execution path: the caller continues
immediately, and the new goroutine starts executing the function on its own stack."
```

**Pattern:**

- Start with a concrete analogy ("window," "straw," "fire alarm")
- Map the analogy to runtime mechanics explicitly
- Address the most common misconception directly

---

### 2. `under_the_hood`

**Purpose:** Mechanical explanation of what happens at runtime. Not API documentation — explain the machinery.

**Good — escape analysis (core-12-14):**

```
"Escape analysis is the compiler pass that decides whether a value lives on the stack or the heap...
The compiler builds SSA (Static Single Assignment) intermediate representation. Escape analysis traces
every allocation through function calls, assignments, and returns. If a value's address is assigned to
a global, returned, stored in an interface, captured by a closure... -> HEAP. If the address never
leaves the function scope -> STACK."
```

**Good — goroutines (core-11-11):**

```
"When go func() is called, runtime.newproc allocates a g struct (~600 bytes on heap) and a stack (~4KB
from stackalloc). The goroutine is placed on the P's local run queue. The scheduler (schedule() in
runtime/proc.go) picks runnable goroutines and executes them on M (OS thread). When a goroutine blocks
(channel, syscall, mutex), the runtime parks it (gopark) and schedules another goroutine on the same M."
```

---

### 3. `execution_timeline`

**Purpose:** Show what happens over time in `→`-prefixed sequential steps. Required for request handling, concurrency, and I/O lessons.

**Good — context cancellation (core-11-06):**

```
→ handler calls context.WithTimeout(ctx, 5*time.Second) creating child ctx with timer
→ child ctx and cancel func returned; goroutine stores child ctx, defers cancel()
→ DB query goroutine starts with child ctx: ctx.Done() channel created, goroutine selects on it
→ client disconnects: server cancels the base request context
→ parent ctx done channel closes → child ctx inherits cancellation
→ child ctx.done channel closes
→ DB query goroutine receives from ctx.Done(), aborts the query
→ defers fire: cancel() called (safe to call multiple times), resources cleaned up
→ goroutine exits, handler returns
```

**Format rules:**

- Each step starts with `→ ` (U+2192 + space)
- The first line may be an un-prefixed "initial state" description (e.g., `Initial state: 0 leaked goroutines`)
- Steps are ordered chronologically
- Branching paths use indented sub-steps
- End with resolution or termination

---

### 4. `memory_timeline`

**Purpose:** Show what happens in memory over time. Required for slices, pointers, interfaces, escape analysis, goroutines, and goroutine leaks.

**Good — slice len/cap (core-03-15):**

```
make([]int, 3, 5) → allocates backing array of 5 ints; len=3, cap=5; elements 0-2 are zero-valued
→ append(s, 1) → writes to index 3; len becomes 4; cap unchanged (still 5); no allocation
→ append(s, 2) → writes to index 4; len becomes 5; cap unchanged; no allocation
→ append(s, 3) → exceeds cap 5; runtime allocates new array (cap doubles to 10); copies existing 5
  elements; writes new element at index 5; len=6, cap=10; old array unreachable
```

**Good — goroutine leaks (core-11-23):**

```
Initial state: 0 leaked goroutines, heap baseline = 50MB
→ HTTP request arrives: handler spawns go logMetrics() — goroutine stack allocated (~4KB)
→ logMetrics blocks on ch <- metric (no reader): goroutine parked (Gwaiting), stack preserved
→ After 1000 requests: 1000 leaked goroutines × 4KB stack = 4MB minimum
→ Over 24 hours at 100 req/s: 8.64M goroutines → 34GB stack → process OOMs
```

---

### 5. `failure_modes`

**Purpose:** Describe each failure with trigger → symptom → root cause → detection.

**Good — double WriteHeader (core-08-07):**

```
"Double WriteHeader panic: handler calls w.WriteHeader(200), then encounters an error and calls
w.WriteHeader(500) — the second call panics because the status was already sent; the client already
received 200 OK."
```

**Good — goroutine leak OOM (core-11-23):**

```
"Unbounded goroutine leak in a loop: an HTTP handler spawns a goroutine on every request without tracking
it — under load, thousands of leaked goroutines consume all available memory and the process OOMs."
```

**Pattern:**

- Name the failure descriptively
- Describe the trigger condition
- Describe the observable symptom
- Explain the mechanical root cause
- Implicitly or explicitly state the detection strategy (metric, log, profile)

---

### 6. `beginner_mistakes`

**Purpose:** Common mistakes with mechanical root cause, not just "don't do this."

**Good — double body read (core-08-05):**

```
"Not limiting the request body size — a client sends a 1GB payload, and ioutil.ReadAll allocates 1GB
of memory, causing OOM. Always use http.MaxBytesReader or io.LimitReader."
```

**Good — closing from receiver (core-11-16):**

```
"Closing a channel from the receiver side — the receiver does not know if all senders are done; closing
from the receiver causes the sender to panic on the next send. The owner (sender) should always close
the channel."
```

**Pattern:**

- State the mistake concretely
- Explain the mechanical consequence
- Provide the correct pattern

---

### 7. `hidden_magic_checks`

**Purpose:** Questions that force mechanical reasoning — if a learner can answer these, they truly understand.

**Good — escape analysis (core-12-14):**

```
"Learner must run go build -gcflags=-m on a given function and identify which allocations escape and why,
comparing with their own prediction."
```

**Good — goroutines (core-11-11):**

```
"Learner must predict the output of a program that starts a goroutine that prints 'hello' and then the
main goroutine prints 'world', explaining why the order is non-deterministic."
```

---

### 8. `proof_of_understanding`

**Purpose:** A concrete task that proves mastery. Should require applying the mechanics, not repeating them.

**Good — escape analysis (core-12-14):**

```
"The learner must write three versions of a function that creates a struct and: (1) returns the struct
by value, (2) returns a pointer to a local struct, (3) returns the value via interface{}. They must run
-gcflags=-m on each version, record which allocations escape and why, and predict the relative GC
overhead of each version under high load."
```

**Good — request parsing (core-08-05):**

```
"The learner must write an HTTP handler that: (1) accepts a JSON payload with a name field, (2) limits
the body size to 1MB using http.MaxBytesReader, (3) decodes the JSON, (4) logs the name, (5) returns
200. They must then write a test that sends a 2MB payload and verifies the handler returns 413, and
another test that sends a valid payload and verifies 200."
```

---

### 9. `problem_solved` and `why_it_exists`

**Purpose:** Frame the lesson before the mechanics.

**Good — `problem_solved` (core-12-14):**

```
"New Go developers do not know whether their allocations go to stack or heap. They use pointers everywhere
(fearing copies), accidentally cause escapes, and wonder why GC overhead is high. They optimize by
guessing — 'use pointers to avoid copies' — which paradoxically increases heap allocation."
```

**Good — `why_it_exists` (core-11-11):**

```
"Goroutines exist because OS threads are too expensive for the concurrency model Go wanted to support.
Creating an OS thread takes ~1µs and ~1MB of stack — creating 10,000 threads would use 10GB of memory.
Goroutines cost ~200ns to create and ~4KB of stack — 10,000 goroutines use ~40MB. By multiplexing
goroutines onto a small number of OS threads, Go makes concurrency cheap enough to use per-connection,
per-request, and per-task — enabling the 'goroutine per connection' pattern that makes Go HTTP servers
handle tens of thousands of concurrent connections."
```

---

## Recommended Fields

The 9 required fields above are the minimum. Add these recommended fields when the lesson's topic matches the spec's trigger tables:

| Field                        | When to include                                                                    |
| ---------------------------- | ---------------------------------------------------------------------------------- |
| `execution_timeline`         | Request handling, concurrency, I/O (see spec trigger table)                        |
| `memory_timeline`            | Slices, pointers, escape analysis, goroutines, interfaces (see spec trigger table) |
| `performance_implications`   | Any lesson with allocation, concurrency, serialization, or query patterns          |
| `debugging_walkthroughs`     | Debugging or incident-response lessons                                             |
| `production_examples`        | Any lesson teaching a pattern used differently in production vs tutorials          |
| `operational_considerations` | Lessons about databases, networking, deployment, or observability                  |

Check the golden lesson index below for lessons that exemplify each recommended field.

---

## Quality Self-Review Checklist

Before submitting, score your lesson against each dimension from the full rubric (see `quality-rubric.md` for detailed criteria):

| Dimension            | Target (golden) | Check                                                                 |
| -------------------- | --------------- | --------------------------------------------------------------------- |
| Clarity              | 4+              | Can a motivated learner understand without external help?             |
| Operational realism  | 4+              | Do examples use realistic data and production patterns?               |
| Magic elimination    | 4+              | Can the learner predict behavior before running code?                 |
| Debuggability        | 4+              | Does the lesson teach how to diagnose failures?                       |
| Production relevance | 4+              | Are performance implications and operational concerns addressed?      |
| Mental-model quality | 4+              | Is the analogy accurate, durable, and non-generic?                    |
| Cognitive Pacing     | 4+              | Are prerequisites respected and complexity appropriate?               |
| Diagram Usefulness   | 4+              | Do diagrams clarify rather than decorate?                             |
| Code-Reading Quality | 4+              | Does the lesson teach code reading, not just writing?                 |
| Failure realism      | 4+              | Do failure modes describe trigger → symptom → root cause → detection? |

Passing threshold for golden standard: **4/5 per dimension, 45/50 overall.**

---

## Before/After: Placeholder → Golden

### Before (detected by validator)

```json
"mental_model": "Think of Request parsing as one step in the learner's path..."
"problem_solved": "This lesson explains what problem Request parsing solves for a Go developer."
"under_the_hood": "The lesson must explain the underlying mechanics..."
```

### After

```json
"mental_model": "The request body is a stream, not a buffer. It can be read exactly once..."
"problem_solved": "New Go developers treat the request body as a reusable buffer..."
"under_the_hood": "http.Request.Body is an http.body wrapper around the connection's bufio.Reader..."
```

The validator flags the "before" patterns via these phrases:

- `"mechanically without understanding"`
- `"one step in the learner's path"`
- `"in the context of professional Go"`

---

## Common Pitfalls

| Pitfall               | Example                                   | Fix                                                                                                  |
| --------------------- | ----------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| API documentation     | "Handler.ServeHTTP handles HTTP requests" | Explain the execution timeline: TCP accept → request parse → handler call → response write           |
| Abstract mental model | "Think of it as a box"                    | Be specific: "The request body is a stream, not a buffer. Reading it is like drinking from a straw." |
| Missing root cause    | "Don't double-read the body"              | Explain why: "The body stream is consumed on first read; subsequent reads return io.EOF"             |
| Toy examples          | Use x, y, z as variable names             | Use realistic data: order_id, user_name, status_code                                                 |
| No detection strategy | "The pool can exhaust"                    | Add: "Monitor with db.Stats().OpenConnections; alert when OpenConnections >= MaxOpenConns \* 0.9"    |

---

## Reference: Golden Lesson Index

These 27 lessons exemplify the standards:

| ID         | Topic                      | Best exemplar fields                         |
| ---------- | -------------------------- | -------------------------------------------- |
| core-03-15 | Slice length and capacity  | mental_model, memory_timeline                |
| core-03-16 | Slice sharing and aliasing | memory_timeline, failure_modes               |
| core-04-08 | Pointer and value mutation | memory_timeline, step_by_step_execution      |
| core-04-09 | Errors as values           | performance_implications, beginner_mistakes  |
| core-04-15 | defer mechanics            | execution_timeline, failure_modes            |
| core-04-17 | panic and recover          | mental_model, under_the_hood                 |
| core-05-05 | Receiver sets              | execution_timeline, hidden_magic_checks      |
| core-05-08 | Interfaces                 | mental_model, memory_timeline                |
| core-05-14 | Nil interfaces             | mental_model, memory_timeline                |
| core-08-03 | Handler lifecycle          | execution_timeline, failure_modes            |
| core-08-05 | Request parsing            | mental_model, beginner_mistakes              |
| core-08-07 | Response writing           | failure_modes, execution_timeline            |
| core-09-11 | sql.DB as a pool           | execution_timeline, performance_implications |
| core-09-17 | Transactions               | failure_modes, under_the_hood                |
| core-11-06 | Context cancellation       | execution_timeline, mental_model             |
| core-11-11 | Goroutines                 | mental_model, memory_timeline                |
| core-11-16 | Channel ownership          | beginner_mistakes, failure_modes             |
| core-11-23 | Goroutine leaks            | memory_timeline, failure_modes               |
| core-12-04 | Correlation IDs            | execution_timeline, real_world_usage         |
| core-12-14 | Escape analysis            | memory_timeline, proof_of_understanding      |
| core-12-01 | Why observability exists   | problem_solved, why_it_exists                |
| core-12-19 | Incident debugging         | debugging_walkthroughs, failure_modes        |
| core-13-08 | Invariants                 | under_the_hood, failure_modes                |
| core-13-11 | Retries and idempotency    | step_by_step_execution, failure_modes        |
| opslane-05 | Tenant isolation           | execution_timeline, hidden_magic_checks      |
| opslane-07 | Order processing           | failure_modes, step_by_step_execution        |
| opslane-13 | Graceful shutdown          | execution_timeline, beginner_mistakes        |

Read any of these lessons in `path.core.json` to see the complete zero_magic block in context.
