# Advanced Thinking Sections

This document is a reference for later-stage lessons that need deeper design and trade-off prompts.

It is not a foundations template.

## 1. Error Design

Use in `02-engineering-core` and above, after learners already understand ordinary `(value, error)` flow.

```markdown
## Thinking Questions

### Error Taxonomy
1. When should you create custom error types instead of plain errors?
2. What is the difference between wrapping with `%w` and formatting with `%v`?
3. How do you keep error codes stable across versions?

### Recovery Strategy
1. When should a caller retry versus return immediately?
2. When is backoff required?
3. Where does a circuit breaker belong?

### Layer Boundaries
1. How should repository, service, and handler layers shape errors differently?
2. When should one layer translate an error from a lower layer?
3. What error information should stay internal?
```

## 2. Concurrent Systems

Use in `04-concurrency`.

```markdown
## Thinking Questions

### Goroutine Management
1. What is the difference between bounded and unbounded goroutine creation?
2. How do worker pools change shutdown behavior?
3. What makes a goroutine leak easy to miss?

### Channel Design
1. When should you use buffered versus unbuffered channels?
2. How do channels create backpressure?
3. When is fan-out/fan-in a good fit?

### Synchronization
1. When is a mutex simpler than a channel?
2. When is `RWMutex` worth the complexity?
3. What state should never be shared directly?
```

## 3. Performance and Memory

Use in `07-advanced`.

```markdown
## Thinking Questions

### Memory Behavior
1. What is escape analysis and when does it matter?
2. What changes when data moves from stack to heap?
3. What allocation patterns increase GC pressure?

### Measurement
1. When should you profile instead of guessing?
2. How do CPU and memory profiles answer different questions?
3. What is the hot path in this code?

### Trade-offs
1. When is micro-optimization worth the readability cost?
2. How do you compare clarity against allocation savings?
3. Which optimization would you reject even if it benchmarks faster?
```

## 4. System Design

Use in `07-advanced` and `08-projects`.

```markdown
## Thinking Questions

### Architecture
1. What boundary is this design trying to protect?
2. When is a monolith better than a distributed design?
3. Which dependency direction keeps the system easiest to change?

### Reliability
1. What breaks first when a dependency slows down?
2. Where should rate limiting or load shedding happen?
3. How would you stop one subsystem from cascading failure into another?

### Observability
1. What signal would tell you the system is degraded before users complain?
2. When is logging enough, and when do you need metrics or tracing?
3. What would make a 3 AM incident easier to debug?
```

## Reusable Shape

```markdown
## Thinking Questions

### Design Level
- What decision is this lesson really about?
- What trade-off does the learner need to justify?

### Boundary Level
- What inputs, outputs, or dependencies make this tricky?
- What should stay simple even if other parts grow?

### Failure Level
- What breaks first?
- What evidence would prove the diagnosis?

### Evolution Level
- How would this design change in a larger system?
- What extension would you reject for now?
```

## Rule

Use these prompts when the learner has already earned:

- the core syntax
- the local mental model
- the basic happy path

Foundations should stay with smaller, local questions.
