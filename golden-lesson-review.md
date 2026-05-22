# Golden Lesson Zero-Magic Review

**Date:** 2026-05-22
**Scope:** All 27 golden lessons from `curriculum/path.core.json`
**Spec:** `curriculum/zero-magic-authoring-spec.md` v1
**Reviewer:** OpenCode agent

---

## 1. Summary Table

| # | Lesson ID | Title | Required (9/9) | Magic Elim (1-5) | Mental Model (1-5) | Failure Realism (1-5) | Exec Timeline | Mem Timeline |
|---|---|---|---|---|---|---|---|---|
| 1 | core-03-15 | Slice length and capacity | ✅ | 5 | 5 | 5 | ❌ missing | ✅ present |
| 2 | core-03-16 | Slice sharing and aliasing | ✅ | 5 | 5 | 5 | ❌ missing | ✅ present |
| 3 | core-04-08 | Pointer and value mutation | ✅ | 5 | 5 | 5 | ✅ → format | ✅ present |
| 4 | core-04-09 | Errors as values | ✅ | 5 | 5 | 5 | ❌ missing | ❌ missing |
| 5 | core-04-15 | defer mechanics | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 6 | core-04-17 | panic and recover | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 7 | core-05-05 | Receiver sets | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 8 | core-05-08 | Interfaces | ✅ | 5 | 4 | 4 | ❌ missing | ✅ present |
| 9 | core-05-14 | Nil interfaces | ✅ | 5 | 5 | 5 | ❌ missing | ✅ present |
| 10 | core-08-03 | Handler lifecycle | ✅ | 5 | 4 | 5 | ✅ → format | ❌ missing |
| 11 | core-08-05 | Request parsing | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 12 | core-08-07 | Response writing | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 13 | core-09-11 | sql.DB as a pool | ✅ | 5 | 4 | 5 | ❌ missing† | ❌ missing |
| 14 | core-09-17 | Transactions | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 15 | core-11-06 | Cancellation | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 16 | core-11-11 | Goroutines | ✅ | 5 | 4 | 5 | ✅ → format | ❌ missing‡ |
| 17 | core-11-16 | Channel ownership | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 18 | core-11-23 | Goroutine leaks | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing‡ |
| 19 | core-12-01 | Why observability exists | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 20 | core-12-04 | Correlation IDs | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 21 | core-12-14 | Escape analysis | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing‡ |
| 22 | core-12-19 | Incident debugging | ✅ | 5 | 4 | 5 | ✅ → format | ❌ missing |
| 23 | core-13-08 | Invariants | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 24 | core-13-11 | Retries and idempotency | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 25 | opslane-05 | Tenant isolation | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 26 | opslane-07 | Order processing | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |
| 27 | opslane-13 | Graceful shutdown | ✅ | 5 | 5 | 5 | ✅ → format | ❌ missing |

**Notes:**
† `core-09-11` involves I/O (database queries) and concurrency (pool sharing across goroutines) — should have `execution_timeline`.
‡ `core-11-11` (goroutines), `core-11-23` (goroutine leaks), and `core-12-14` (escape analysis) involve memory allocation directly — should have `memory_timeline`.

---

## 2. Lessons Missing Required Fields

**None.** All 27 golden lessons contain all 9 required fields:
- `problem_solved`
- `why_it_exists`
- `mental_model`
- `under_the_hood`
- `step_by_step_execution`
- `hidden_magic_checks`
- `failure_modes`
- `beginner_mistakes`
- `proof_of_understanding`

---

## 3. Lessons with Weak mental_model or under_the_hood (Score ≤ 2)

**None.** All 27 golden lessons have substantive, mechanically-grounded mental models and under_the_hood sections. Scores range 4–5.

The placeholder patterns (e.g., "Think of X as one step in the learner's path from raw code to professional engineering judgment") appear only in non-golden lessons (e.g., core-05-06, core-05-07, core-05-09, etc.) — these are correctly excluded from the golden set.

---

## 4. Lessons That Should Have execution_timeline But Don't

Per spec: *"Every lesson involving request handling, concurrency, or I/O MUST include an execution timeline"*

| Lesson | Topic | Reason It Should Have One |
|--------|-------|--------------------------|
| core-03-15 | Slice length and capacity | Uses append, allocation — sequential execution with branching paths (allocation vs no allocation). Borderline. |
| core-03-16 | Slice sharing and aliasing | Sub-slicing involves sequential operations with aliasing. Borderline. |
| core-04-09 | Errors as values | Sequential error-return path; step_by_step_execution exists but not in `→` format. Borderline. |
| core-05-08 | Interfaces | Interface method dispatch is a runtime operation with sequential steps. Borderline. |
| core-05-14 | Nil interfaces | Two-word interface representation has a clear sequential timeline. Borderline. |
| **core-09-11** | **sql.DB as a pool** | **DEFINITE: involves I/O (database queries) and concurrency (pool shared across goroutines). Has `step_by_step_execution` but no `execution_timeline` in the required `→` format.** |

**Primary gap:** core-09-11 (sql.DB as a pool) — this lesson describes a concurrent connection pool servicing I/O requests and should follow the Execution Timeline Standard.

---

## 5. Lessons That Should Have memory_timeline But Don't

Per spec: *"Every lesson involving data structures, pointers, or allocation MUST include a memory timeline"*

| Lesson | Topic | Reason It Should Have One |
|--------|-------|--------------------------|
| **core-11-11** | **Goroutines** | **DEFINITE: goroutine stack allocation (~4KB initial), M:P:G model, memory consumption at scale. The `under_the_hood` describes `runtime.newproc` and stack allocation — a memory timeline would show stack growth and GC interaction.** |
| **core-11-23** | **Goroutine leaks** | **DEFINITE: entire lesson is about memory consumption from leaked goroutines. The `under_the_hood` describes GC scanning leaked stacks — a memory timeline would show cumulative memory growth.** |
| **core-12-14** | **Escape analysis** | **DEFINITE: the entire lesson is about stack vs heap allocation decisions. A memory timeline showing value lifetime from creation to escape decision to allocation would directly serve the learning objective.** |
| core-04-09 | Errors as values | Interface values involve two-word (type, data) representation — could benefit from memory timeline. |
| core-04-15 | defer mechanics | Defer stack (linked list of `_defer` structs) involves runtime data structures — could benefit. |
| core-04-17 | panic and recover | Panic stack (`_panic` struct, per-goroutine) involves runtime data structures — could benefit. |
| core-05-05 | Receiver sets | Itab/itable construction involves compile-time data structures — borderline. |
| core-09-11 | sql.DB as a pool | Connection pool internals (channel-based semaphore, free list) involve data structures — could benefit. |
| core-11-16 | Channel ownership | `hchan` struct involves mutex, circular buffer, send/recv queues — could benefit from memory timeline. |
| core-13-08 | Invariants | Database isolation levels and locking — involves data integrity, not memory layout. Optional. |

**Primary gaps:** core-11-11 (goroutines), core-11-23 (goroutine leaks), core-12-14 (escape analysis) — all three directly involve memory allocation models and would be significantly strengthened by a `memory_timeline`.

---

## 6. Recommended Spec Refinements

### 6.1 Clarify `step_by_step_execution` vs `execution_timeline` Boundary

Currently, the spec defines `execution_timeline` as a recommended field with a specific `→` format, while `step_by_step_execution` is required. In practice, many golden lessons have rich `step_by_step_execution` arrays that already describe time-ordered execution but without the `execution_timeline`-specific `→` prefix format.

**Recommendation:** Add guidance on when `step_by_step_execution` (required) suffices vs when `execution_timeline` (recommended) is additionally required. Specifically:
- `step_by_step_execution` = the required mechanical walkthrough (can be prose)
- `execution_timeline` = the recommended `→`-prefixed format for lessons involving request handling, concurrency, or I/O

This would resolve the ambiguity for lessons like core-03-15, core-03-16, core-04-09, core-05-08, core-05-14, and core-09-11.

### 6.2 Strengthen Memory Timeline Trigger Criteria

The current spec says "data structures, pointers, or allocation" — this is broad enough to cover nearly every lesson. Consider narrowing to:

> "Every lesson whose primary learning objective involves understanding how values are laid out in memory, how memory is allocated/freed, or how pointers/references affect state MUST include a `memory_timeline`."

This would clearly flag:
- core-11-11 (goroutines: stack allocation)
- core-11-23 (goroutine leaks: memory leak mechanics)
- core-12-14 (escape analysis: allocation decisions)

While excluding lessons where memory is secondary context.

### 6.3 Add `failure_modes` Completeness Criteria

The spec says each failure mode should describe trigger → symptom → root cause → detection strategy. The detection strategy is implicit in most existing failure modes. Consider making it explicit:

```diff
 Each failure mode should describe:
 - The trigger condition
 - The observable symptom
 - The mechanical root cause
 - The detection strategy
+++ (e.g., which metric, log line, or profile artifact reveals this failure)
```

### 6.4 Add Validator Integration for Golden Lesson Gaps

The spec mentions future validators for:
- Execution timeline presence in golden lessons
- Memory timeline presence for memory-related concepts

Based on this review, prioritize:
1. **Memory timeline validator**: flag golden lessons whose module involves allocation/pointers (modules 03, 04, 11, 12-14) if `memory_timeline` is absent
2. **Execution timeline validator**: flag golden lessons involving request handling (module 08), concurrency (module 11), or I/O (module 09) if `execution_timeline` is absent

### 6.5 Consider Making `memory_timeline` Required for Specific Modules

The data shows a stark pattern: only 5 of 27 golden lessons have `memory_timeline` (core-03-15, core-03-16, core-04-08, core-05-08, core-05-14). All 5 are from early modules (03–05) that deal with slice internals and interface representation. Later modules (08–13, opslane) have zero memory timelines despite involving significant memory concepts.

**Recommendation:** If `memory_timeline` is valued as a standard, consider making it required (not recommended) for any golden lesson whose `under_the_hood` section describes runtime data structures, allocations, or pointer mechanics. This would capture ~15 more lessons.

### 6.6 Execution Timeline `→` Format Compliance

The spec defines the execution_timeline format as:
```
start condition
→ step 1: what happens
→ step 2: what happens next
```

Reviewing existing `execution_timeline` fields:
- Core-04-08: Uses `→ ` prefix but not "→ step N:" numbering — close but not exact
- Core-04-15: Same pattern — uses `→ ` prefix with descriptive text
- Core-08-03: Same pattern
- Core-08-05: Same pattern
- All authored `execution_timeline` fields follow the `→ ` prefix convention consistently but omit explicit "step N:" labels

**Recommendation:** The spec format is being followed in spirit (arrow-prefixed sequential steps). Consider whether explicit "step N:" numbering is required or whether the `→` prefix alone satisfies the standard. If strict numbering is desired, add explicit examples in the spec.

---

## Appendix A: Recommended Fields Coverage

| Recommended Field | Lessons Present In | Lessons Absent From |
|---|---|---|
| `execution_timeline` | 19 lessons | 8 lessons (03-15, 03-16, 04-09, 05-08, 05-14, 09-11, plus others listed in §4) |
| `memory_timeline` | 5 lessons | 22 lessons |
| `debugging_walkthroughs` | 0 lessons | All 27 |
| `production_examples` | 3 lessons (03-15, 04-09, 11-06) | 24 lessons |
| `performance_implications` | 24 lessons | 3 lessons (none — covered well) |
| `operational_considerations` | 0 lessons | All 27 |
| `code_reading_tasks` | 0 lessons | All 27 |
| `refactoring_tasks` | 0 lessons | All 27 |
| `review_questions` | 0 lessons | All 27 |

**Note:** `debugging_walkthroughs`, `operational_considerations`, `code_reading_tasks`, `refactoring_tasks`, and `review_questions` are recommended but not implemented in any golden lesson. These may be intentionally deferred to README content or exercise files rather than the JSON metadata.

## Appendix B: Mental Model Quality Notes

All 27 golden lessons have non-generic, durable analogies. Standout examples:
- **core-03-15**: "A slice is a window into a backing array. Append moves the window — and sometimes builds a bigger house behind it."
- **core-05-14**: "An interface value is a box with two compartments: one holds the type tag, one holds the data."
- **core-08-05**: "The request body is a stream, not a buffer. Reading the body is like drinking from a straw — once consumed, it is gone."
- **core-11-06**: "The done channel is like a fire alarm: once it rings, it keeps ringing, and every goroutine in the building should have a plan for what to do when it hears it."
- **opslane-05**: "Tenant isolation is a data firewall between tenants."

No generic "one step in the learner's path" patterns found in any golden lesson.

## Appendix C: Failure Mode Completeness

Failure modes across all 27 lessons consistently describe trigger + observable symptom + mechanical root cause. Detection strategy is present implicitly in most cases (e.g., "the log aggregator drops events", "the client parses a truncated JSON document") but could be more explicit. Strongest failure realism scores (5/5) go to lessons where the failure mode describes a specific production scenario with concrete observability signals.
