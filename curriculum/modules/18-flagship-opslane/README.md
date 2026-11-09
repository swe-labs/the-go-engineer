# Flagship Opslane

> **Module 18 — Flagship Capstone.** Final graduation proof through a portfolio-grade backend system.

## Flagship mission

Opslane is the capstone of the core path: a multi-tenant incident-management SaaS backend that integrates every skill from Modules 00–15 — Go language mechanics, data modeling, persistence, authentication and authorization, concurrency and lifecycle, observability, distributed architecture, and containerized deployment — into one production-shaped system you can run, operate, and defend.

The mission is not to ship a toy demo. It is to prove you can design a real backend from first principles, secure it across tenant boundaries, observe it under load, deploy it reproducibly, and explain your engineering decisions to a panel.

> **Full Architectural Specification:** For complete system architecture, data models, event bus topology, and API endpoint contracts, see [`docs/flagship/OPSLANE_SAAS_BACKEND.md`](../../../docs/flagship/OPSLANE_SAAS_BACKEND.md).

## What Opslane proves

Every milestone maps a core-path capability onto a concrete production artifact:

| What it proves | Learned in | Demonstrated by |
|---|---|---|
| Layered backend design (transport → domain → repo) | Modules 05, 14 | HTTP API layer and repository separation |
| Relational persistence and migrations | Module 09 | Schema, models, and migrations |
| Authentication and authorization | Module 10 | Session/token auth and role checks |
| Tenant isolation | Modules 10, 14 | Row-scoped queries and enforced tenant boundaries |
| Concurrency and lifecycle | Module 11 | Event bus, worker pools, caching, graceful shutdown |
| Observability | Module 12 | Structured logs, metrics, and traces |
| Reproducible deployment | Module 15 | Dockerfile, CI pipeline, and deployment runbook |
| Engineering communication | Module 16 | Architecture write-up and live portfolio defense |

## System overview

Opslane is a single Go binary organized into three clean domain tiers, wrapped by a small amount of process wiring:

```
cmd/opslane         -> process entry, configuration loading, dependency wiring
transport/http      -> routers, middleware, request/response DTOs
domain/             -> business rules, entities, service interfaces
repo/postgres       -> SQL queries, transactions, migrations
```

Supporting systems:

- **Event bus** — an in-process, buffered publish/subscribe hub that decouples order processing from downstream work.
- **Worker pools** — bounded goroutine pools that drain the event bus with backpressure and context cancellation.
- **Cache layer** — a read-through cache in front of hot tenant and catalog lookups.
- **Observability** — `slog` structured logs, Prometheus metrics, and OpenTelemetry spans on every request.
- **Graceful shutdown** — a `SIGTERM` handler drains in-flight work, closes the pools, and flushes the bus before exit.

## Required proof

Graduation requires the following artifacts, each verified by its assessment:

1. **The Opslane codebase** — the running system, built and tested: [`projects/opslane-flagship/README.md`](./projects/opslane-flagship/README.md).
2. **Fourteen checkpoint assessments** — `opslane-01` through `opslane-14`, one per capability, each backed by a physical artifact and a passing rubric.
3. **The module checkpoint** — [`assessments/checkpoint/README.md`](./assessments/checkpoint/README.md): the integrated engineering review.

## Portfolio defense

The final capstone is a live defense of the whole system: [`projects/15-opslane-portfolio-defense/README.md`](./projects/15-opslane-portfolio-defense/README.md). You will walk a reviewer through the architecture, justify each tier boundary and dependency choice, demonstrate the security and tenant-isolation guarantees, present the observability and deployment evidence, and reason about the top production failure modes.

## Prerequisites

- [Module 15: Docker, CI/CD, and Deployment](../15-docker-cicd-deployment/README.md) — the deployment runbook (`core-15-19`) is the technical prerequisite.
- The full core path (Modules 00–15) supplies the engineering foundation the milestones assume.
- [Module 16: Portfolio and Interview Readiness](../16-portfolio-interview-readiness/README.md) is intentionally **non-blocking**: you can begin Opslane immediately after Module 15, then feed the finished system into your portfolio packaging. [Module 17: Advanced Electives](../../electives/advanced-electives/README.md) is an optional companion track.

## Lessons

Each lesson is a checkpoint milestone with its own assessment and physical artifact:

- `opslane-01` [Opslane foundation and configuration](./assessments/01-opslane-foundation-and-configuration/README.md)
- `opslane-02` [Opslane database and models](./assessments/02-opslane-database-and-models/README.md)
- `opslane-03` [Opslane authentication](./assessments/03-opslane-authentication/README.md)
- `opslane-04` [Opslane authorization](./assessments/04-opslane-authorization/README.md)
- `opslane-05` [Opslane tenant isolation](./assessments/05-opslane-tenant-isolation/README.md)
- `opslane-06` [Opslane HTTP API layer](./assessments/06-opslane-http-api-layer/README.md)
- `opslane-07` [Opslane order processing](./assessments/07-opslane-order-processing/README.md)
- `opslane-08` [Opslane payment pipeline](./assessments/08-opslane-payment-pipeline/README.md)
- `opslane-09` [Opslane event bus](./assessments/09-opslane-event-bus/README.md)
- `opslane-10` [Opslane worker pools](./assessments/10-opslane-worker-pools/README.md)
- `opslane-11` [Opslane caching layer](./assessments/11-opslane-caching-layer/README.md)
- `opslane-12` [Opslane observability](./assessments/12-opslane-observability/README.md)
- `opslane-13` [Opslane graceful shutdown](./assessments/13-opslane-graceful-shutdown/README.md)
- `opslane-14` [Opslane deployment](./assessments/14-opslane-deployment/README.md)
- `opslane-15` [Opslane portfolio defense](./projects/15-opslane-portfolio-defense/README.md)

## Completion

Complete every required lesson, lab, project, and assessment for this module.
