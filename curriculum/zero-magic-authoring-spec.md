# Zero-Magic Authoring Spec v2

> Changes from v1:
>
> - Clarified `step_by_step_execution` vs `execution_timeline` boundary
> - Strengthened memory timeline trigger criteria with specific examples
> - Added failure mode completeness criteria (detection strategy)
> - Added `→` format compliance note
> - Updated Validator Integration with prioritized roadmap
> - Added Quality Rubric reference
> - Added format compliance examples

## Purpose

Every lesson must eliminate "magic" — concepts that learners can use but cannot explain
mechanically. This spec defines the pedagogical contract for all curriculum content.

## Core Principle

A learner should never write code they cannot explain step by step.

If a learner can complete an exercise but cannot describe _what happens in memory_,
_what happens over time_, or _what happens on failure_, the lesson has failed.

---

## Pedagogical Contract

Every lesson MUST contain these fields in its `zero_magic` block:

### Required Fields

| Field                    | Type          | Purpose                                                    |
| ------------------------ | ------------- | ---------------------------------------------------------- |
| `problem_solved`         | string        | What specific confusion or gap does this lesson eliminate? |
| `why_it_exists`          | string        | Why does this concept exist in Go/the system?              |
| `mental_model`           | string        | A durable analogy or mental model                          |
| `under_the_hood`         | string        | Mechanical explanation of what happens at runtime          |
| `step_by_step_execution` | array[string] | Ordered sequence of events during execution (prose format) |
| `hidden_magic_checks`    | array[string] | Checks that prevent magic-reliant explanations             |
| `failure_modes`          | array[string] | What breaks and how it fails                               |
| `beginner_mistakes`      | array[string] | Common mistakes with mechanical root cause                 |
| `proof_of_understanding` | string        | How the learner proves mastery                             |

### Recommended Fields

| Field                        | Type          | Purpose                                                                                       |
| ---------------------------- | ------------- | --------------------------------------------------------------------------------------------- |
| `execution_timeline`         | array[string] | Time-ordered event chain with `→` prefix (for I/O, concurrency, request handling)             |
| `memory_timeline`            | array[string] | Memory-ordered event chain with `→` prefix (for memory layout, allocation, pointer mechanics) |
| `debugging_walkthroughs`     | array[string] | Step-by-step debugging scenarios                                                              |
| `production_examples`        | array[string] | Real production code patterns                                                                 |
| `performance_implications`   | array[string] | Performance characteristics and tradeoffs                                                     |
| `operational_considerations` | array[string] | Production operations impact                                                                  |
| `code_reading_tasks`         | array[string] | Exercises in reading existing code                                                            |
| `refactoring_tasks`          | array[string] | Exercises in improving existing code                                                          |
| `review_questions`           | array[string] | Questions that test mechanical understanding                                                  |

---

## Explanation Philosophy

### 1. Execution Timeline Explanations

Do NOT only explain _what_ a concept is.

Explain _exactly what happens over time_.

**Bad**: "Context cancellation stops goroutines."

**Good**:

```
request arrives
→ handler creates context.WithTimeout
→ DB query starts with context
→ client disconnects
→ context canceled
→ done channel closes
→ query goroutine receives cancellation
→ cleanup executes
→ goroutine exits
→ handler returns
```

### 2. Memory Timeline Explanations

For memory-related concepts, explain _exactly what happens in memory over time_.

**Bad**: "Slices share the underlying array."

**Good**:

```
slice header copied to function
→ len=3 cap=5, both point to same array
→ append called, len becomes 4
→ still within cap, no new allocation
→ different goroutine appends, exceeds cap
→ new array allocated, references diverge
```

### 3. Failure Mode Explanations

Every production concept must answer: "What breaks and how?"

Each failure mode should describe:

- The trigger condition
- The observable symptom
- The mechanical root cause
- The detection strategy (e.g., which metric, log line, or profile artifact reveals this failure)

### 4. Magic Elimination

A lesson contains "magic" if a learner can:

- Use the concept correctly
- But cannot explain the runtime mechanics

Eliminate magic by requiring learners to:

- Draw memory diagrams
- Trace execution timelines
- Walk through failure scenarios
- Read and explain source code before writing it

---

## Quality Rubric

Every lesson is reviewed against the dimensions in `quality-rubric.md`.

Minimum scores for golden lessons: 4/5 per dimension, 45/50 overall.

---

## Execution Timeline Standard

Every lesson involving **request handling, concurrency, or I/O** MUST include an
`execution_timeline` in the following `→`-prefixed format:

```
start condition
→ step 1: what happens
→ step 2: what happens next
→ step 3: branching point
  → path A: what happens
  → path B: what happens
→ step 4: resolution
→ end condition
```

**Format compliance:** The `→ ` prefix (U+2192 + space) at the start of each step
is the minimum requirement. Explicit "step N:" numbering is optional.
`execution_timeline` differs from `step_by_step_execution` in format (arrow-prefixed)
and in scope (focused on cross-cutting event chains like request lifecycle,
connection pooling, or goroutine scheduling).

**Trigger table:**

| Lesson Topic                                         | Must have `execution_timeline`?              |
| ---------------------------------------------------- | -------------------------------------------- |
| HTTP request handling (handler, middleware, parsing) | Yes                                          |
| Concurrency (goroutines, channels, select)           | Yes                                          |
| I/O (database, network, filesystem)                  | Yes                                          |
| Data structures (slices, maps, interfaces)           | Optional — `step_by_step_execution` suffices |
| Error handling, panics, defer                        | Optional — `step_by_step_execution` suffices |

---

## Memory Timeline Standard

Every lesson whose primary learning objective involves **how values are laid out in
memory, how memory is allocated/freed, or how pointers/references affect state**
MUST include a `memory_timeline` showing state changes over time.

Format:

```
initial state: [variable] → [memory description]
→ operation 1: [change description]
→ operation 2: [change description]
→ final state: [variable] → [memory description]
```

**Trigger table:**

| Lesson Topic                                    | Must have `memory_timeline`? |
| ----------------------------------------------- | ---------------------------- |
| Slice internals (len/cap, sharing, aliasing)    | Yes                          |
| Pointer and value mutation behavior             | Yes                          |
| Escape analysis (stack vs heap decisions)       | Yes                          |
| Goroutines (stack allocation, growth)           | Yes                          |
| Goroutine leaks (cumulative memory consumption) | Yes                          |
| Interface representation (two-word itable)      | Yes                          |
| defer/panic runtime data structures             | Recommended                  |
| Channel internals (hchan struct)                | Recommended                  |
| Connection pools, transaction pinning           | Optional                     |
| Errors as values (interface wrapping)           | Optional                     |

---

## Timeline Format Compliance

Both `execution_timeline` and `memory_timeline` MUST use the `→` (U+2192) prefix
for each step. Acceptable variations:

```
→ step description
→ step description: detail
  → sub-step description
initial state: description
```

Do NOT use markdown lists (`-`, `*`, `1.`) or code fences inside individual steps.
Each step is a plain string in a JSON array; the `→` prefix is the visual indicator.

---

## Validator Integration

The curriculum tool checks golden lessons via `validate-zero-magic` for:

- All 9 required fields present
- No placeholder content (detected by phrase matching)
- Mental model and under_the_hood quality (non-generic)

Priority-ordered roadmap for additional validators:

1. **Memory timeline validator** — flag golden lessons with topics from the
   Memory Timeline Standard trigger table if `memory_timeline` is absent
2. **Execution timeline validator** — flag golden lessons with topics from the
   Execution Timeline Standard trigger table if `execution_timeline` is absent
3. **Failure mode completeness** — verify each failure mode describes trigger +
   symptom + root cause + detection strategy
4. **Beginner mistake root-cause quality** — verify each mistake explains the
   mechanical root cause, not just the surface behavior
5. **Timeline format validator** — verify `execution_timeline` and `memory_timeline`
   entries use the correct `→` prefix format
