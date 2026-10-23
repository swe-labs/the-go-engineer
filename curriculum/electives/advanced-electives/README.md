# Module 17: Advanced Electives

> **Phase:** elective | **Status:** draft | **Required:** False | **Pacing:** Self-directed / post-core

## Goal

Provide focused, production-grade deep dives into advanced Go type mechanics, binary RPC protocols, specialized concurrency primitives, and distributed systems architectures without blocking or artificially inflating the core engineering pathway.

---

## Architectural Context & Learning Path Flow

Module 17 is intentionally structured as a **companion elective track**. It sits outside the critical path required for the Flagship Opslane Capstone (Module 18) and initial career portfolio deployment (Module 16):

```
Core Path: [Mod 00 - 15] ──► [Mod 18: Flagship Opslane] ──► [Mod 16: Portfolio Packaging] ──► Graduation
                                        │
                                        └──► [Mod 17: Advanced Electives (Optional Specialization)]
```

- **Core Prerequisite Baseline:** Completion of Modules 00 through 15 (foundations through containerized production deployment).
- **Decoupling Rationale:** Topics such as gRPC, Protocol Buffers, CQRS, and `sync.Pool` represent critical production specializations, but forcing them into the first 40 weeks creates cognitive overload and distracts from core software engineering habits (clean interfaces, error handling, relational persistence, testing, and lifecycle concurrency).
- **Specialization Tracks:** Learners choose focused tracks relevant to their immediate engineering domain or target team stack.

---

## Specialization Streams & Reference Topics

Module 17 is organized into four cohesive streams totaling 20 reference topics:

### Stream A: Advanced Go Typing & Metaprogramming
Focuses on type parametrization, static safety, and compile-time code generation.

| ID | Slug | Title | Focus & Mechanics |
|---|---|---|---|
| `elective-01` | `generics` | **Generics** | Type parameters, constraints, and compile-time monomorphization. |
| `elective-02` | `complex-generic-constraints` | **Complex generic constraints** | Constraint interfaces, union elements (`~T`), and comparable semantics. |
| `elective-03` | `generic-data-structures` | **Generic data structures** | Type-safe linked lists, generic rings, and thread-safe generic caches. |
| `elective-04` | `functional-options` | **Functional options** | Extensible API design, configuration mutation, and builder patterns. |
| `elective-05` | `method-values` | **Method values** | Bound method values, currying closures, and callback passing. |
| `elective-06` | `build-tags` | **Build tags** | Conditional compilation (`//go:build`), platform and environment isolation. |
| `elective-07` | `go-generate` | **go generate** | Automating code generation, `stringer`, and custom generator binaries. |

### Stream B: Modern RPC & Serialization Protocols
Focuses on contract-first binary communications and service-to-service micro-protocols.

| ID | Slug | Title | Focus & Mechanics |
|---|---|---|---|
| `elective-08` | `mockery` | **Mockery** | Automated interface mock generation and mocking lifecycle strategies. |
| `elective-09` | `protobuf` | **Protobuf** | Schema-first IDL, binary wire format, and field serialization semantics. |
| `elective-10` | `grpc-fundamentals` | **gRPC fundamentals** | HTTP/2 multiplexing, unary RPC lifecycles, and code generation. |
| `elective-11` | `grpc-streaming` | **gRPC streaming** | Client, server, and bidirectional streaming data pipelines. |
| `elective-12` | `grpc-interceptors` | **gRPC interceptors** | Middleware chains, authentication propagation, and distributed tracing metadata. |
| `elective-13` | `rest-vs-grpc-tradeoffs` | **REST vs gRPC tradeoffs** | Wire size, latency, tooling, browser compatibility, and architectural fitness. |

### Stream C: High-Throughput Concurrency Primitives
Focuses on low-level memory reuse and lock contention reduction in extreme-throughput scenarios.

| ID | Slug | Title | Focus & Mechanics |
|---|---|---|---|
| `elective-14` | `sync-map` | **sync.Map** | Read-heavy concurrent access patterns, dirty/read maps, and lock avoidance. |
| `elective-15` | `sync-pool` | **sync.Pool** | GC allocation mitigation, ring buffers, and byte buffer recycling. |

### Stream D: Distributed Systems Patterns & Progressive Delivery
Focuses on advanced distributed consistency, operational decoupling, and release engineering.

| ID | Slug | Title | Focus & Mechanics |
|---|---|---|---|
| `elective-16` | `cqrs` | **CQRS** | Command Query Responsibility Segregation, read/write models, and projection pipelines. |
| `elective-17` | `feature-flags` | **Feature flags** | Runtime toggling, canary targeting, and trunk-based deployment decoupling. |
| `elective-18` | `blue-green-deployments` | **Blue/green deployments** | Zero-downtime routing, traffic flipping, and rollback validation. |
| `elective-19` | `microservice-splitting` | **Microservice splitting** | Strangler fig pattern, bounded context isolation, and transactional cutovers. |
| `elective-20` | `deeper-distributed-systems` | **Deeper distributed systems** | Consensus protocols (Raft), CAP tradeoffs, quorum writes, and vector clocks. |

---

## Applied Projects & Milestone Assessment

### Capstone Project: `project-advanced-electives`
- **Title:** Advanced Electives Capstone (`advanced-electives-capstone`)
- **Deliverables:**
  1. Production-grade Go codebase implementing a selected specialization stream (e.g., custom gRPC micro-service with interceptor telemetry, or a high-throughput lock-free cache using `sync.Pool`).
  2. Benchmark suite comparing naive vs. optimized implementations with flame graph profiles.
  3. Architecture Decision Record (ADR) justifying why the elective pattern was chosen over standard library idioms.
- **Verification:** Automated benchmark comparisons, race detector verification, and peer architectural review.

### Milestone Assessment: `assessment-module-17`
- **Title:** Module 17 Assessment -- Electives and Advanced Engineering
- **Format:** Technical project defense and architectural inquiry.
- **Passing Threshold:** 80% across conceptual clarity (30%), implementation quality (35%), empirical proof (25%), and engineering communication (10%).

---

## Completion Checklist

To complete an elective specialization:
- [ ] Implement and test all lessons in the chosen specialization stream.
- [ ] Deliver a functional implementation in `curriculum/electives/advanced-electives/projects/advanced-electives-capstone`.
- [ ] Provide automated verification (unit tests, benchmarks, and race analysis).
- [ ] Pass the defense rubric for `assessment-module-17`.
