# Module 17: Advanced Electives

> **Phase:** elective | **Status:** draft | **Required:** False | **Pacing:** Self-directed / post-core

## Goal

Provide focused, production-grade deep dives into advanced Go type mechanics, binary RPC protocols, specialized concurrency primitives, distributed systems architectures, and Kubernetes operations without blocking or artificially inflating the core engineering pathway.

## Prerequisites

- [Module 16: Portfolio and Interview Readiness](../../modules/16-portfolio-interview-readiness/README.md) — complete that module's completion checkpoint before starting a specialization stream. The full core path (Modules 00–15) supplies the engineering foundations these electives build on.
- Electives are **optional and non-blocking**: they can be taken alongside or after the flagship capstone.

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

## Lessons

The track holds 26 self-contained lessons across five specialization streams. Each link opens the lesson README:

- **Stream A — Advanced Go typing & metaprogramming:** [`01-generics`](./lessons/01-generics/README.md), [`02-complex-generic-constraints`](./lessons/02-complex-generic-constraints/README.md), [`03-generic-data-structures`](./lessons/03-generic-data-structures/README.md), [`04-functional-options`](./lessons/04-functional-options/README.md), [`05-method-values`](./lessons/05-method-values/README.md), [`06-build-tags`](./lessons/06-build-tags/README.md), [`07-go-generate`](./lessons/07-go-generate/README.md)
- **Stream B — Modern RPC & serialization protocols:** [`08-mockery`](./lessons/08-mockery/README.md), [`09-protobuf`](./lessons/09-protobuf/README.md), [`10-grpc-fundamentals`](./lessons/10-grpc-fundamentals/README.md), [`11-grpc-streaming`](./lessons/11-grpc-streaming/README.md), [`12-grpc-interceptors`](./lessons/12-grpc-interceptors/README.md), [`13-rest-vs-grpc-tradeoffs`](./lessons/13-rest-vs-grpc-tradeoffs/README.md)
- **Stream C — High-throughput concurrency primitives:** [`14-sync-map`](./lessons/14-sync-map/README.md), [`15-sync-pool`](./lessons/15-sync-pool/README.md)
- **Stream D — Distributed systems patterns & progressive delivery:** [`16-cqrs`](./lessons/16-cqrs/README.md), [`17-feature-flags`](./lessons/17-feature-flags/README.md), [`18-blue-green-deployments`](./lessons/18-blue-green-deployments/README.md), [`19-microservice-splitting`](./lessons/19-microservice-splitting/README.md), [`20-deeper-distributed-systems`](./lessons/20-deeper-distributed-systems/README.md)
- **Stream E — Kubernetes operations:** [`21-why-orchestration`](./lessons/21-why-orchestration/README.md), [`22-kubernetes-object-model`](./lessons/22-kubernetes-object-model/README.md), [`23-config-and-security-context`](./lessons/23-config-and-security-context/README.md), [`24-probes-and-graceful-shutdown`](./lessons/24-probes-and-graceful-shutdown/README.md), [`25-resources-and-go-runtime`](./lessons/25-resources-and-go-runtime/README.md), [`26-rollouts-and-troubleshooting`](./lessons/26-rollouts-and-troubleshooting/README.md)

The section below expands each stream with its focus and mechanics.

---

## Specialization Streams & Reference Topics

Module 17 is organized into five cohesive streams totaling 26 reference topics:

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

### Stream E: Kubernetes Operations
Focuses on running containerized Go services on a container orchestrator, from the desired-state model through to production-grade operation on a local cluster.

| ID | Slug | Title | Focus & Mechanics |
|---|---|---|---|
| `elective-21` | `why-orchestration` | **Why orchestration** | Desired state, observed state, the reconciliation loop, and when to adopt an orchestrator. |
| `elective-22` | `kubernetes-object-model` | **Pods, Deployments, Services, and Ingress** | The object graph and the request path from client to container. |
| `elective-23` | `config-and-security-context` | **Configuration, secrets, and security context** | ConfigMap and Secret delivery, fail-fast loading, and a hardened security context. |
| `elective-24` | `probes-and-graceful-shutdown` | **Probes, lifecycle, and graceful shutdown** | Startup, readiness, and liveness probes plus a bounded drain on `SIGTERM`. |
| `elective-25` | `resources-and-go-runtime` | **Resources, the Go runtime, and autoscaling** | Requests and limits, `GOMAXPROCS`, `GOMEMLIMIT`, CFS throttling, and OOM kills. |
| `elective-26` | `rollouts-and-troubleshooting` | **Rollouts, rollback, and troubleshooting** | Rolling updates, rollback, PodDisruptionBudget, and stalled-rollout diagnosis. |

---

## Applied Projects & Milestone Assessment

### Capstone Project: `project-advanced-electives`
- **Title:** Advanced Electives Capstone (`advanced-electives-capstone`)
- **Deliverables:**
  1. Production-grade Go codebase implementing a selected specialization stream (e.g., custom gRPC micro-service with interceptor telemetry, or a high-throughput lock-free cache using `sync.Pool`).
  2. Benchmark suite comparing naive vs. optimized implementations with flame graph profiles.
  3. Architecture Decision Record (ADR) justifying why the elective pattern was chosen over standard library idioms.
- **Verification:** Automated benchmark comparisons, race detector verification, and peer architectural review.

### Capstone Project: `project-kubernetes-go-service`
- **Title:** Kubernetes Go Service Capstone (`kubernetes-go-service-capstone`)
- **Prerequisites:** Stream E (`elective-21` through `elective-26`), anchored on `elective-26` and the core deployment runbook (`core-15-19`).
- **Deliverables:**
  1. The Module 15 Go service deployed to a local kind or k3d cluster with a Service and an Ingress (or a NodePort fallback).
  2. A Deployment with resource requests and limits, all three probes, a non-root read-only security context, and a ConfigMap/Secret split.
  3. A demonstrated rolling update, a rollback, and a PodDisruptionBudget held during a node drain.
  4. A runbook extension covering the rollout, verification, rollback, and the top three failure modes.
- **Verification:** Run the verification commands to confirm rollout status, the serving revision, the security context, resource settings, the PodDisruptionBudget, the rollback, and a node drain, then defend the result in writing.

### Milestone Assessment: `assessment-module-17`
- **Title:** Module 17 Assessment -- Electives and Advanced Engineering
- **Format:** Technical project defense and architectural inquiry.
- **Passing Threshold:** 80% across conceptual clarity (30%), implementation quality (35%), empirical proof (25%), and engineering communication (10%).

### Project Assessment: `assessment-project-kubernetes-go-service`
- **Title:** Assessment -- Kubernetes Go Service Capstone
- **Format:** Evidence review and written defense.
- **Passing Threshold:** 80% across correctness (35%), production readiness (25%), proof (25%), and explanation (15%).

## Assessments

- [`assessments/advanced-electives/README.md`](./assessments/advanced-electives/README.md) — **Module 17 Assessment — Electives and Advanced Engineering**
- [`assessments/kubernetes-go-service/README.md`](./assessments/kubernetes-go-service/README.md) — **Assessment — Kubernetes Go Service Capstone**

---

## Completion Checklist

To complete an elective specialization:
- [ ] Implement and test all lessons in the chosen specialization stream.
- [ ] Deliver a functional implementation in `curriculum/electives/advanced-electives/projects/advanced-electives-capstone`.
- [ ] Provide automated verification (unit tests, benchmarks, and race analysis).
- [ ] Pass the defense rubric for `assessment-module-17`.

The Kubernetes Operations stream (Stream E) is completed through its own capstone:
- [ ] Deploy a Go service to a local cluster and complete `curriculum/electives/advanced-electives/projects/kubernetes-go-service-capstone`.
- [ ] Demonstrate a rolling update, a rollback, and a node drain under a PodDisruptionBudget.
- [ ] Pass the defense rubric for `assessment-project-kubernetes-go-service`.

## Next module

Module 17 is an optional companion track, so it does not gate the core arc. When your chosen stream and its capstone are complete, rejoin the core path at [Module 18: Flagship Opslane](../../modules/18-flagship-opslane/README.md) to assemble the graduation capstone and portfolio defense.
