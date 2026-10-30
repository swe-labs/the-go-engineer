# Zero-Magic Authoring Specification (v3 Curriculum Standard)

> **Pedagogical Authority:** The Go Engineer Curriculum Architecture  
> **Status:** Active Authoritative Standard (`v3.1.0` Schema)  
> **Core Mission:** Eliminate "magic" — ensuring learners never write code they cannot explain mechanically, down to memory layout, temporal execution, and runtime failure modes.

## Purpose

Every lesson must eliminate "magic" — concepts that learners can use but cannot explain
mechanically. This spec defines the pedagogical contract for all curriculum content.

## Core Principle

A learner should never write code they cannot explain step by step.

If a learner can complete an exercise but cannot describe _what happens in physical memory_,
_what happens over time_, or _what happens on failure_, the lesson has failed.

---

## The Universal Physical Footprint Standard (The 6 Physical Pillars)

Every lesson introducing Go primitives, data structures, concurrency, collections, or I/O MUST ground the concept in its physical machine reality:

1. **Pillar 1: Primitives & Zero-Value Hardware Guarantee**:
   - State size in bytes and bits (e.g. `uint8` = 1 byte/8 bits; `int64` = 8 bytes/64 bits) and word boundary alignment.
   - Explain the physical **Zero-Value Guarantee**: Contiguous memory initialized by the Go runtime to all binary zeros (`00000000`). It is not "empty" or "null"; it is a concrete bit pattern (`0`, `false`, `nil`).
   - Explain **Type Projection**: Bits have no intrinsic meaning; the type system instructs the compiler how many bytes to read and how to project those bits into meaningful operations.
2. **Pillar 2: Composite Descriptors & Machine Headers**:
   - Slices: 24-byte header (`*data`, `len`, `cap`) pointing to a contiguous backing array.
   - Strings: 16-byte header (`*data`, `len`) pointing to read-only byte arrays.
   - Interfaces: 16-byte header (`*itab`, `*data`). Passing a value to an interface forces heap allocation so `*data` holds a valid pointer.
3. **Pillar 3: Concurrency Engine Physical Footprint**:
   - Goroutines: 2,048-byte (2 KB) contiguous initial stack frame + runtime `g` struct descriptor (~376–400B). Depth triggers contiguous stack doubling on the heap. Provide capacity math ($N \times \text{bytes}$).
   - Channels (`hchan`): ~96B struct with circular ring buffer `buf`, 8B spinlock `lock mutex`, and `waitq` doubly linked lists with `sudog` (~152B each).
   - Mutexes (`sync.Mutex`): 8 bytes (4B `state` bitmask + 4B `sema`).
4. **Pillar 4: Hash Table Mechanics & Bucket Topology (`map`)**:
   - `hmap` pointer (48B) pointing to array of `bmap` buckets (8 key-value pairs per bucket).
   - Deleting keys (`delete(m, key)`) sets tombstones but **never frees bucket memory back to the OS**.
5. **Pillar 5: Hardware Physics & Cache Lines**:
   - 64-byte L1/L2/L3 CPU cache lines.
   - Multi-core false sharing: insert 64-byte padding (`_ [56]byte` or `cpu.CacheLinePad`) between hot multi-core fields.
   - Struct field ordering: order descending by size (8B, 4B, 2B, 1B) to eliminate compiler zero-padding bloat.
6. **Pillar 6: I/O, Network & OS Buffers**:
   - `bufio` 4,096-byte (4 KB) page buffers.
   - OS kernel socket buffers (`SO_RCVBUF`, `SO_SNDBUF`, 128KB–2MB).
   - Draining `resp.Body` to `io.Discard` before closing prevents TCP resets (`RST`), keeping keep-alive connection pools alive.

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

### 5. Spec-Driven Authoring Protocol

All instructional authoring must follow the 4-step engineering pipeline:
1. **Specification Contract**: Define measurable behavioral objectives and exact physical preconditions.
2. **Table-Driven TDD Test Scaffold**: Construct table-driven test cases covering nominal, boundary, and failure conditions.
3. **Hardware-Grounded Implementation & Prose**: Write the mechanical narrative and complete runnable Go code.
4. **Quality Gate Audit**: Verify against all 17 automated quality gates (`go run ./scripts/validate_curriculum.go`).

### 6. Cognitive Anchor Quotes Standard

To anchor difficult systems mental models, integrate authoritative aphorisms from computing pioneers:
- *"Simplicity is prerequisite for reliability."* — Edsger W. Dijkstra (clean minimal interfaces, clean domain boundaries)
- *"There are two ways of constructing a software design: One way is to make it so simple that there are obviously no deficiencies..."* — C.A.R. Hoare (zero-magic mechanics)
- *"Clear is better than clever."* — Rob Pike (explicit error handling `if err != nil`)
- *"Debugging is twice as hard as writing the code in the first place..."* — Brian Kernighan (Delve inspection over clever code)

### 7. 3-Tier Clean Domain Architecture Standard for Capstones

All capstones and backend service projects must strictly decouple:
- `transport/`: HTTP handlers, routing, JSON parsing, status codes. Enforces boundary sanitization; never passes raw `*http.Request` into domain entities.
- `domain/`: Pure business entities, state invariants, domain errors, repository interfaces. Never imports transport or database drivers.
- `repo/`: Implements domain storage interfaces using SQL/caching.

### 8. Delve (`dlv`) Runtime Inspection Standard

Debugging walkthroughs must feature real Delve (`dlv`) stack frame navigation commands (`frame`, `locals`, `goroutines`, `goroutine <id> stack`) rather than defaulting to scattershot `fmt.Println` print debugging.

### 9. Anti-Overengineering & Anti-Slop Directive

Do not force all 17 headings onto a lesson if the subject matter does not warrant them. Forcing an orientation lesson, shell introduction, or career strategy document to invent a fake `Performance implications` or fake `Under the hood` produces **artificial template slop**, violating zero-magic principles. Omit non-applicable sections rather than populating them with placeholder filler. Always respect the item's `content_contract` flags and `subtype_profiles` defined in `curriculum/readme.contracts.json`.

- **6 Core Invariant Headings**: Mandatory across ALL lessons (`Learning objective`, `Why this matters`, `Mental model`, `Common mistakes`, `Review questions`, `NEXT UP`).
- **11 Contextual Headings**: Enforced only when applicable based on `subtype` and `content_contract` flags.

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

## Visual Model & Mermaid Diagram Standard

Every curriculum README and lesson specifying a `## Visual model` must provide a concrete, fully-rendered Mermaid diagram (or mechanical ASCII layout if documenting raw memory word offsets). Mental models in Go systems engineering require spatial and temporal anchoring: describing concurrent interleavings or middleware recursion purely in prose invites dangerous misconceptions.

### Go Domain to Mermaid Syntax Mapping

| Domain / Concept Area | Preferred Mermaid Syntax | Mechanical Purpose | Example Use Case |
|---|---|---|---|
| **Concurrency & Synchronization** | `sequenceDiagram` | Channel send/recv handshakes, `select` fairness, mutex contention, worker rendezvous. | `goroutine-leaks`, `channel-buffering`, `bounded-worker-pool` |
| **HTTP Middleware & Request Lifecycle** | `sequenceDiagram` | Onion-style wrap (`ServeHTTP` recursion), context deadline propagation, response interception. | `middleware-chains`, `correlation-id`, `book-catalog-rest-api` |
| **System Architecture & Packaging** | `flowchart TD` or `flowchart LR` | Hexagonal ports and adapters, Clean Architecture boundaries, package dependency graph. | `clean-architecture`, `modular-refactor-review`, `opslane` |
| **Pipelines & Worker Pools** | `flowchart LR` | Data streaming, fan-out/fan-in worker topologies, backpressure queueing. | `bounded-worker-pool`, `concurrent-downloader`, `stream-processing` |
| **Lifecycle & State Transitions** | `stateDiagram-v2` | Entity lifecycle, connection pool state transitions (`Idle` -> `Active` -> `Closed`), circuit breakers. | `database-pooling`, `task-api-with-postgres`, `url-health-checker` |
| **Data Structures & Memory Offsets** | `flowchart TD` or ASCII box | Slice headers (`ptr`, `len`, `cap`), interface two-word `eface`/`iface` structures, ring buffers. | `slice-internals`, `interface-internals`, `bank-account` |
| **Database Schemas & Relations** | `erDiagram` | Multi-tenant isolation, foreign keys, cascade constraints. | `task-api-with-postgres`, `secure-authenticated-api` |
| **Version Control & Migrations** | `gitGraph` | Database schema migration forward/rollback chains, git trunk workflows. | `database-migrations`, `dockerized-service-with-ci-cd` |

### Authoring Rules & Quality Gates for Diagrams

1. **Zero Generic Placeholders**: Never use abstract placeholders (e.g. `input -> concept boundary -> behavior -> observable result`). Every node and edge must represent a concrete Go type, channel, goroutine, or system component.
2. **Syntax Safety**:
   - Always quote labels containing spaces, punctuation, or Go syntax (e.g. `G1["Goroutine 1: producer"] -->|"ch <- item"| C["Buffered Channel (cap=4)"]`).
   - All blocks (`subgraph`, `loop`, `alt`, `opt`) must have matching `end` delimiters.
   - Do not inject raw HTML tags (`<br>`, `<div >`) that break non-browser renderers; use Mermaid newline syntax `\n` or standard labels.
3. **Coupled Narrative**: A diagram cannot stand alone. It must be accompanied by a caption explaining what critical transition to observe, and immediately connect into the `## Machine view` or execution mechanics.

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
