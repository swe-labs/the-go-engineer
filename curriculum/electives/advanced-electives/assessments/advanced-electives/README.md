# Project Assessment — Advanced Electives Capstone

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `project-advanced-electives`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Visual Model: Specialization Stream Verification Pipeline

```
+-----------------------------------------------------------------------------------------+
|                          ADVANCED ELECTIVES CAPSTONE PIPELINE                           |
+-----------------------------------------------------------------------------------------+
|                                                                                         |
|   Specialization Streams                                                                |
|   +-----------------------+     Verification Gate         Assessment Deliverables       |
|   | Stream A: Generics    |    +----------------------+   +-----------------------+     |
|   | Stream B: gRPC / Proto|--->| - go test -v -race   |-->| Specialized Codebase  |     |
|   | Stream C: sync.Pool   |    | - Benchmark profiles |   | Flame Graph Profiles  |     |
|   | Stream D: CQRS / Raft |    | - Zero Allocation    |   | Architecture Decision |     |
|   +-----------------------+    +----------------------+   +-----------------------+     |
|                                                                                         |
+-----------------------------------------------------------------------------------------+
```

---

## 2. Assessment Requirements

1. **Production-Grade Implementation**:
   - Implement one or more specialization streams under `curriculum/electives/advanced-electives/`.
   - All tests must pass cleanly under race detection: `go test -race ./...`.
2. **Benchmark Proof**:
   - Deliver benchmark profiles comparing the naive implementation vs the optimized pattern (e.g., measuring latency and allocations).
3. **Architecture Decision Record (ADR)**:
   - Provide an ADR documenting why the elective pattern was chosen and the operational trade-offs accepted.

---

## 3. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Inadequate** | Code fails to compile or fails tests; no benchmarks; cannot explain trade-offs. |
| **Level 2: Developing** | Compiles and passes basic unit tests; lacks race detection verification; benchmarks without memory analysis. |
| **Level 3: Proficient** | All 20 elective tests compile and pass; benchmark profiles provided; ADR clearly documents architectural rationale. |
| **Level 4: Exemplary** | Exceptional mastery. Compares CPU cache line mechanics, monomorphization assembly output, or gRPC HTTP/2 framing; demonstrates measurable zero-allocation performance gains. |

---

## 4. Verification

```bash
# Run all 20 advanced elective test suites
go test -v ./curriculum/electives/advanced-electives/...
```
