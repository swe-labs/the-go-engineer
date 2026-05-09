# Roadmap

> This document tracks what is built, in progress, and planned.
> Section IDs and names here must match `ARCHITECTURE.md` exactly.
> If they conflict, `ARCHITECTURE.md` wins.

---

## Branch Model

- `main` — active v2 development and prereleases
- `release/v1` — stable v1 maintenance for current users
- `release/v2` — cut from `main` when v2 reaches feature freeze

---

## Current Status (as of April 2026)

The 12-section v2.1 architecture is approved. Active structural migration is underway.

- `release/v1`: stable for current learners
- `main`: primary curriculum build target
- Architecture: `ARCHITECTURE.md` v2.1

---

## Status Legend

| Symbol | Meaning                                              |
| ------ | ---------------------------------------------------- |
| ✅     | Content exists, registered in JSON, validator passes |
| 🚧     | Content partially exists or in progress              |
| 📋     | Planned, not started                                 |
| ❌     | Critical gap — blocking learner path                 |

---

## Phase 0: Machine Foundation

| ID  | Section            | Status | Notes                                                                              |
| --- | ------------------ | ------ | ---------------------------------------------------------------------------------- |
| s00 | How Computers Work | ✅     | Lessons HC.1–HC.8 are registered, runnable, and covered by foundations validation. |

### Lessons (s00)

| ID   | Lesson                        | Status |
| ---- | ----------------------------- | ------ |
| HC.1 | What is a program?            | ✅     |
| HC.2 | How code becomes execution    | ✅     |
| HC.3 | Memory basics — stack vs heap | ✅     |
| HC.4 | Terminal confidence           | ✅     |
| HC.5 | How the OS manages processes  | ✅     |
| HC.6 | CPU cache and performance     | ✅     |
| HC.7 | Syscalls                      | ✅     |
| HC.8 | Blocking vs non-blocking I/O  | ✅     |

---

## Phase 1: Language Foundation

| ID  | Section            | Status | Notes                                                                                   |
| --- | ------------------ | ------ | --------------------------------------------------------------------------------------- |
| s01 | Getting Started    | 🚧     | GT.1–GT.4 exist. GT.5–GT.6 not in JSON.                                                 |
| s02 | Language Basics    | 🚧     | LB, DS exist. CF missing CF.5–CF.6 (defer).                                             |
| s03 | Functions & Errors | 🚧     | FE.1–FE.7 exist. FE.8–FE.10 not created.                                                |
| s04 | Types & Design     | ✅     | TI, CO, and ST are registered, and validation now covers their alternate path families. |

### Lessons (s01)

| ID   | Lesson                    | Status                          |
| ---- | ------------------------- | ------------------------------- |
| GT.1 | Installation verification | ✅                              |
| GT.2 | Hello World               | ✅                              |
| GT.3 | How Go works              | ✅                              |
| GT.4 | Development environment   | ✅                              |
| GT.5 | go fmt, go vet, go doc    | 🚧 content written, not in JSON |
| GT.6 | Reading compiler errors   | 🚧 content written, not in JSON |

### Lessons (s02) — new items only

| ID   | Lesson                    | Status                               |
| ---- | ------------------------- | ------------------------------------ |
| CF.5 | Defer — mechanics & order | ❌ not created; blocks DB.3 and FS.1 |
| CF.6 | Defer in real use cases   | ❌ not created                       |

### Lessons (s03) — new items only

| ID    | Lesson                | Status                                                 |
| ----- | --------------------- | ------------------------------------------------------ |
| FE.8  | First-class functions | ❌ not created; blocks GC.1 understanding              |
| FE.9  | Closures — mechanics  | ❌ not created; blocks goroutine capture understanding |
| FE.10 | panic & recover       | ❌ not created; blocks HTTP middleware understanding   |

---

## Phase 2: Engineering Core

| ID  | Section                   | Status | Notes                                                       |
| --- | ------------------------- | ------ | ----------------------------------------------------------- |
| s05 | Packages, I/O & CLI       | ✅     | MP, CL, EN, FS all exist and registered.                    |
| s06 | Backend, APIs & Databases | ❌     | **DB exists. HS (HTTP) and API (gRPC) completely missing.** |
| s07 | Concurrency               | 🚧     | GC, CT, TM, CP exist. SY (sync primitives) missing.         |
| s08 | Quality & Testing         | 🚧     | TE.1–TE.4 exist. TE.5–TE.10 and PR.3–PR.5 missing.          |

### Lessons (s06) — new items

| ID    | Lesson                         | Status |
| ----- | ------------------------------ | ------ |
| HS.1  | net/http basics                | ❌     |
| HS.2  | Routing patterns               | ❌     |
| HS.3  | Middleware — the pattern       | ❌     |
| HS.4  | Request parsing & validation   | ❌     |
| HS.5  | Response writing patterns      | ❌     |
| HS.6  | Error handling middleware      | ❌     |
| HS.7  | Server timeouts                | ❌     |
| HS.8  | Graceful HTTP shutdown         | ❌     |
| HS.9  | Health & readiness probes      | ❌     |
| HS.10 | REST API                       | ❌     |
| API.1 | REST design principles         | ❌     |
| API.2 | API versioning                 | ❌     |
| API.3 | Pagination & filtering         | ❌     |
| API.4 | Protobuf basics                | ❌     |
| API.5 | gRPC fundamentals              | ❌     |
| API.6 | gRPC streaming                 | ❌     |
| API.7 | gRPC interceptors              | ❌     |
| API.8 | REST vs gRPC trade-off         | ❌     |
| API.9 | gRPC Service                   | ❌     |
| DB.7  | N+1 query detection            | ❌     |
| DB.8  | Query timeouts via context     | ❌     |

### Lessons (s07) — new items only

| ID   | Lesson                    | Status |
| ---- | ------------------------- | ------ |
| SY.1 | sync.Mutex & RWMutex      | ❌     |
| SY.2 | sync.Once & sync.Map      | ❌     |
| SY.3 | Atomic operations         | ❌     |
| SY.4 | Race conditions           | ❌     |
| SY.5 | Goroutine leaks           | ❌     |
| SY.6 | Deadlocks                 | ❌     |

### Lessons (s08) — new items only

| ID    | Lesson                      | Status |
| ----- | --------------------------- | ------ |
| TE.5  | Sub-tests & t.Cleanup       | ❌     |
| TE.6  | Fuzz testing                | ❌     |
| TE.7  | Interfaces for testability  | ❌     |
| TE.8  | Mocking with interfaces     | ❌     |
| TE.9  | Integration tests           | ❌     |
| TE.10 | Golden files                | ❌     |
| PR.3  | Memory profiling            | ❌     |
| PR.4  | Escape analysis             | ❌     |
| PR.5  | Benchmark-driven dev        | ❌     |

---

## Phase 3: Systems Engineering

| ID  | Section                 | Status | Notes                                                                |
| --- | ----------------------- | ------ | -------------------------------------------------------------------- |
| s09 | Architecture & Security | ❌     | PD.1–PD.3 exist. ARCH (9 lessons) and SEC (11 lessons) missing.      |
| s10 | Production Operations   | 🚧     | SL and GS exist. CFG, OPS, DOCKER/DEPLOY, CG to be added/registered. |

### Lessons (s09) — new items

| ID      | Lesson                                        | Status |
| ------- | --------------------------------------------- | ------ |
| ARCH.1  | Monolith vs Modular Monolith vs Microservices | ❌     |
| ARCH.2  | Domain-Driven Design basics                   | ❌     |
| ARCH.3  | Hexagonal architecture in Go                  | ❌     |
| ARCH.4  | Repository pattern deep dive                  | ❌     |
| ARCH.5  | Service layer pattern                         | ❌     |
| ARCH.6  | Event-driven architecture                     | ❌     |
| ARCH.7  | CQRS basics                                   | ❌     |
| ARCH.8  | When to split services                        | ❌     |
| ARCH.9  | Modular Refactor                              | ❌     |
| SEC.1   | Input validation patterns                     | ❌     |
| SEC.2   | SQL injection prevention                      | ❌     |
| SEC.3   | XSS & CSRF                                    | ❌     |
| SEC.4   | Authentication basics                         | ❌     |
| SEC.5   | JWT — implementation & risks                  | ❌     |
| SEC.6   | Password hashing                              | ❌     |
| SEC.7   | Rate limiting patterns                        | ❌     |
| SEC.8   | TLS & HTTPS in Go                             | ❌     |
| SEC.9   | Secrets management                            | ❌     |
| SEC.10  | OWASP Top 10 for Go                           | ❌     |
| SEC.11  | Secure API                                    | ❌     |

### Lessons (s10) — new items

| ID       | Lesson                        | Status                 |
| -------- | ----------------------------- | ---------------------- |
| CFG.1    | Environment variables         | ❌                     |
| CFG.2    | Configuration files           | ❌                     |
| CFG.3    | Flag parsing                  | ❌                     |
| CFG.4    | 12-Factor App principles      | ❌                     |
| CFG.5    | Config validation on boot     | ❌                     |
| OPS.1    | Metrics basics                | ❌                     |
| OPS.2    | Prometheus integration        | ❌                     |
| OPS.3    | Distributed tracing basics    | ❌                     |
| OPS.4    | Feature flags                 | ❌                     |
| OPS.5    | Alerting mindset              | ❌                     |
| DOCKER.1 | Docker basics                 | ❌                     |
| DOCKER.2 | Multi-stage builds            | ❌                     |
| DOCKER.3 | Docker Compose                | ❌                     |
| DEPLOY.1 | CI/CD pipelines               | ❌                     |
| DEPLOY.2 | Blue/green & rollback         | ❌                     |
| DEPLOY.3 | Dockerised Service            | ❌                     |
| CG.1     | go generate primer            | ✅ (move from old s11) |
| CG.2     | Mockery workflow              | ✅ (move from old s11) |
| CG.3     | sqlc workflow                 | ✅ (move from old s11) |

---

## Phase 4: Flagship Project

| ID  | Section          | Status | Notes                                        |
| --- | ---------------- | ------ | -------------------------------------------- |
| s11 | GoScale Flagship | 🚧     | Skeleton exists. Full modules not built out. |

---

## Immediate Priorities (ordered)

1. Register GT.5–GT.6 in curriculum.v2.json
2. Create CF.6 (defer) + add cross-reference notes to FS.1 and DB.3
3. Create FE.8–FE.10 (closures, first-class functions, panic/recover)
4. Create HS.1–HS.10 (HTTP servers) — largest content gap
5. Create API.1–API.9 (REST design + gRPC)
6. Create SY.1–SY.6 (sync primitives)
7. Create TE.5–TE.10, PR.3–PR.5
8. Create ARCH.1–ARCH.9 + SEC.1–SEC.11
9. Create CFG.1–CFG.5 + OPS.1–OPS.5 + DOCKER.1–DEPLOY.3

## Doc Fixes Required

- ~~`COMMON-MISTAKES.md` — all "Covered in: Section N" references are wrong~~ ✅ Fixed
- ~~`CODE-STANDARDS.md` — remove `HTTP_STATUS_OK = 200` example~~ ✅ Fixed
- ~~`ROADMAP.md` — phase/section boundaries must match new 12-section structure~~ ✅ Fixed
- ~~`LEARNING-PATH.md` — phase/section boundaries must match new 12-section structure~~ ✅ Fixed
- ~~`PROGRESSION.md` — milestone table must match milestones in `ARCHITECTURE.md`~~ ✅ Fixed
- ~~`CONTRIBUTING.md` — emoji in "NEXT UP" footer breaks validator regex~~ ✅ Fixed

## Version Plan

| Version    | Target       | Criteria                                  |
| ---------- | ------------ | ----------------------------------------- |
| v2.0-alpha | current      | s01–s11 with existing content registered  |
| v2.0-beta  | near-term    | HC registered + HS + API + SY added       |
| v2.0-rc    | mid-term     | ARCH + SEC + CFG + DOCKER complete        |
| v2.0       | release      | All 12 sections complete, validator green |
| v2.1       | post-release | GoScale modules fully implemented         |

---
