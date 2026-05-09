# The Go Engineer — Architecture v2.1

## Teaching Software Engineering with Go: Zero to Senior

> This is the single source of truth for curriculum structure.
> All other documents (ROADMAP.md, LEARNING-PATH.md, PROGRESSION.md, curriculum.v2.json)
> must align with what is written here. When a conflict exists, this document wins.

---

## The Teaching Contract (Non-Negotiable)

1. **README first, code second** — every lesson teaches through prose before showing code
2. **Zero magic** — never use a concept before teaching it; cross-reference forward dependencies explicitly: _(more on this in Lesson X.N)_
3. **Every lesson runs** — `go run ./path` or `go test ./path` produces visible, meaningful output
4. **Earn each concept** — a tool appears only after the learner has seen the problem it solves
5. **Production notes** — every lesson includes `⚠️ In Production` with real-world consequences
6. **Thinking questions** — every lesson ends with `🤔 Thinking Questions` (minimum 3)
7. **Failure-based learning** — diagnose real bugs; do not only show the correct solution
8. **Stage-aware depth** — beginner lessons stay safe and explicit; engineering pressure enters Phase 2 onward

---

## Progress Model

```text
Phase 0       Phase 1              Phase 2                Phase 3            Phase 4
(0-5%)        (5-52%)              (52-87%)               (87-96%)           (96-100%)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
MACHINE       LANGUAGE             ENGINEERING            SYSTEMS            FLAGSHIP
FOUNDATION    FOUNDATION           CORE                   ENGINEERING        PROJECT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
"I understand "I write Go          "I build systems       "I design and      "I build and
 the machine"  fluently"            that work at scale"    operate systems"   ship it"
```

---

## 12-Section Curriculum Map

```text
s00  How Computers Work          Phase 0    0%  →  5%    HC
s01  Getting Started             Phase 1    5%  → 12%    GT
s02  Language Basics             Phase 1   12%  → 28%    LB, CF, DS
s03  Functions & Errors          Phase 1   28%  → 38%    FE
s04  Types & Design              Phase 1   38%  → 52%    TI, CO, ST
─────────────────────────────────────────────────────────────────
s05  Packages, I/O & CLI         Phase 2   52%  → 62%    MP, CL, EN, FS
s06  Backend, APIs & Databases   Phase 2   62%  → 75%    HS, API, DB
s07  Concurrency                 Phase 2   75%  → 83%    GC, SY, CT, TM, CP
s08  Quality & Testing           Phase 2   83%  → 87%    TE, PR
─────────────────────────────────────────────────────────────────
s09  Architecture & Security     Phase 3   87%  → 92%    PD, ARCH, SEC
s10  Production Operations       Phase 3   92%  → 96%    SL, GS, CFG, OPS, DOCKER, CG
─────────────────────────────────────────────────────────────────
s11  GoScale Flagship            Phase 4   96%  → 100%   12 modules
```

---

## Phase 0: Machine Foundation

### Section 00 — How Computers Work

**Folder:** `00-how-computers-work/`
**ID prefix:** `HC`

Before writing code, learners need a mental model of the machine. This section explains why code works at all. It is not optional and not a detour.

| ID   | Lesson                        | Type   |
| ---- | ----------------------------- | ------ |
| HC.1 | What is a program?            | lesson |
| HC.2 | How code becomes execution    | lesson |
| HC.3 | Memory basics — stack vs heap | lesson |
| HC.4 | Terminal confidence           | lesson |
| HC.5 | How the OS manages processes  | lesson |
| HC.6 | CPU cache and performance     | lesson |
| HC.7 | Syscalls                      | lesson |
| HC.8 | Blocking vs non-blocking I/O  | lesson |

**Section checkpoint:** Learner can explain in plain language: what the CPU does, how memory is divided, what a process is, why cache matters, what a syscall is, and why waiting on I/O changes performance.

**HC.1-HC.8 must remain registered in curriculum.v2.json and backed by runnable lesson surfaces.**

---

## Phase 1: Language Foundation

### Section 01 — Getting Started

**Folder:** `01-getting-started/`
**ID prefix:** `GT`

Builds confidence that Go is installed, running, and the toolchain is understood before teaching the language.

| ID   | Lesson                       | Status  | Milestone              |
| ---- | ---------------------------- | ------- | ---------------------- |
| GT.1 | Installation verification    | exists  | ✓ First Go install     |
| GT.2 | Hello World                  | exists  | ✓ First program runs   |
| GT.3 | How Go works                 | exists  | —                      |
| GT.4 | Development environment      | exists  | ✓ Dev loop established |
| GT.5 | `go fmt`, `go vet`, `go doc` | **new** | —                      |
| GT.6 | Reading compiler errors      | **new** | ✓ Toolchain confident  |

**GT.5 content:** The daily toolchain workflow. Formatting is enforced, not optional. `go vet` before every commit. `go doc` for self-service documentation.

**GT.6 content:** Error anatomy (file:line:column:message). The six most common error types. The blank identifier `_`. Fix errors top-down, one at a time.

---

### Section 02 — Language Basics

**Folder:** `02-language-basics/`
**ID prefixes:** `LB`, `CF`, `DS`

Teaches the entire Go syntax surface — values, flow, and data — as one zero-magic path.

#### Values (LB)

| ID   | Lesson             | Status |
| ---- | ------------------ | ------ | ---------- |
| LB.1 | Variables & types  | exists |
| LB.2 | Constants          | exists |
| LB.3 | Enums with iota    | exists |
| LB.4 | Application Logger | exists | ✓ Exercise |

#### Control Flow (CF)

| ID   | Lesson                    | Status  | Note                                       |
| ---- | ------------------------- | ------- | ------------------------------------------ |
| CF.1 | If / else                 | exists  | —                                          |
| CF.2 | For basics                | exists  | —                                          |
| CF.3 | Break / continue          | exists  | —                                          |
| CF.4 | Switch                    | exists  | —                                          |
| CF.5 | Defer — mechanics & order | **new** | Must precede any lesson using `defer`      |
| CF.6 | Defer in real use cases   | **new** | File close, mutex unlock, cleanup patterns |
| CF.7 | Pricing Checkout          | exists  | ✓ Exercise (formerly CF.5)                 |

**Critical ordering note:** `defer rows.Close()` appears in DB.3 and `defer file.Close()` appears in FS.1. CF.5 must be completed before a learner reaches either lesson. Until CF.5 is created, FS.1 and DB.3 must carry a cross-reference note: _(Defer is taught in Lesson CF.5. For now, read `defer` as "run this line when the current function returns, no matter what.")_

#### Data Structures (DS)

| ID   | Lesson                   | Status |
| ---- | ------------------------ | ------ | ---------- |
| DS.1 | Arrays                   | exists |
| DS.2 | Slices                   | exists |
| DS.3 | Maps                     | exists |
| DS.4 | Pointers                 | exists |
| DS.5 | Slice sharing & capacity | exists |
| DS.6 | Contact Directory        | exists | ✓ Exercise |

---

### Section 03 — Functions & Errors

**Folder:** `03-functions-errors/`
**ID prefix:** `FE`

Teaches functions as design boundaries and errors as explicit values — the two practices that separate Go from most other languages.

| ID    | Lesson                  | Status  | Engineering Detail            |
| ----- | ----------------------- | ------- | ----------------------------- |
| FE.1  | Functions as boundaries | exists  | Why decompose?                |
| FE.2  | Parameters & returns    | exists  | Pass-by-value cost            |
| FE.3  | Multiple return values  | exists  | —                             |
| FE.4  | Errors as values        | exists  | The (value, error) contract   |
| FE.5  | Validation patterns     | exists  | Fail-fast, guard clauses      |
| FE.6  | Orchestration           | exists  | Failure propagation           |
| FE.7  | Order Summary           | exists  | ✓ Exercise                    |
| FE.8  | First-class functions   | **new** | Callbacks, higher-order       |
| FE.9  | Closures — mechanics    | **new** | Capture, variable lifetime    |
| FE.10 | panic & recover         | **new** | When to use each; when NOT to |

**Why FE.8–FE.10 are required:**

- FE.8–FE.9: COMMON-MISTAKES mistake #1 ("Capturing loop variables incorrectly in closures") cannot be diagnosed without understanding closures. These lessons must precede GC.1 (goroutines).
- FE.10: Required before any HTTP middleware (HS.3), which uses recover internally.

---

### Section 04 — Types & Design

**Folders:** `04-types-design/`, `05-composition/`, `06-strings-and-text/`
**ID prefixes:** `TI`, `CO`, `ST`

Teaches Go's type system, composition model, and text handling. Three subsystems within one section because each builds on the last.

#### Types & Interfaces (TI)

| ID    | Lesson                      | Status |
| ----- | --------------------------- | ------ | ---------- |
| TI.1  | Structs                     | exists |
| TI.2  | Methods                     | exists |
| TI.3  | Interfaces                  | exists |
| TI.4  | Interface embedding         | exists |
| TI.5  | Stringer                    | exists |
| TI.6  | Type switch                 | exists |
| TI.7  | Receiver sets               | exists |
| TI.8  | Custom error types          | exists |
| TI.9  | Generics introduction       | exists |
| TI.10 | Payroll Processor           | exists | ✓ Exercise |
| TI.11 | Empty interface & `any`     | exists |
| TI.12 | Type assertions             | exists |
| TI.13 | Nil interfaces              | exists |
| TI.14 | Functional options          | exists |
| TI.15 | Method values               | exists |
| TI.16 | Complex generic constraints | exists |
| TI.17 | Generic data structures     | exists |

#### Composition (CO)

| ID   | Lesson       | Status |
| ---- | ------------ | ------ | ---------- |
| CO.1 | Composition  | exists |
| CO.2 | Embedding    | exists |
| CO.3 | Bank Account | exists | ✓ Exercise |

#### Strings & Text (ST)

| ID   | Lesson          | Status |
| ---- | --------------- | ------ | ---------- |
| ST.1 | Strings         | exists |
| ST.2 | Formatting      | exists |
| ST.3 | Unicode & runes | exists |
| ST.4 | Regex           | exists |
| ST.5 | Text templates  | exists |
| ST.6 | Config Parser   | exists | ✓ Exercise |

**Note on section size:** Section 04 contains 26 lessons across 3 subsystems. This is appropriate — the type system is the most important design tool in Go and must be thorough. The subsystem labels (TI/CO/ST) provide internal organisation. A learner's progress path flows: TI → CO → ST, then into Phase 2.

---

## Phase 2: Engineering Core

### Section 05 — Packages, I/O & CLI

**Folder:** `05-packages-io/`
**ID prefixes:** `MP`, `CL`, `EN`, `FS`

Teaches how Go programs are organised, how they interact with disk and streams, and how to build CLI tools.

#### Modules & Packages (MP)

| ID   | Lesson                | Status |
| ---- | --------------------- | ------ | ---------- |
| MP.1 | Module basics         | exists |
| MP.2 | Managing dependencies | exists |
| MP.3 | Versioning Workshop   | exists | ✓ Exercise |
| MP.4 | Build tags            | exists |

#### CLI Tools (CL)

| ID   | Lesson         | Status |
| ---- | -------------- | ------ | ---------- |
| CL.1 | Args           | exists |
| CL.2 | Flags          | exists |
| CL.3 | Subcommands    | exists |
| CL.4 | File Organizer | exists | ✓ Exercise |

#### Encoding (EN)

| ID   | Lesson             | Status |
| ---- | ------------------ | ------ | ---------- |
| EN.1 | JSON marshalling   | exists |
| EN.2 | JSON unmarshalling | exists |
| EN.3 | JSON encoder       | exists |
| EN.4 | JSON decoder       | exists |
| EN.5 | Base64             | exists |
| EN.6 | Config Parser      | exists | ✓ Exercise |

#### Filesystem (FS)

| ID   | Lesson             | Status |
| ---- | ------------------ | ------ | ---------- |
| FS.1 | Files              | exists |
| FS.2 | Paths              | exists |
| FS.3 | Directories        | exists |
| FS.4 | Temp files         | exists |
| FS.5 | Embed              | exists |
| FS.6 | I/O patterns       | exists |
| FS.7 | Log Search         | exists | ✓ Exercise |
| FS.8 | fs.FS testing seam | exists |

**Cross-reference note for FS.1:** Uses `defer file.Close()`. Until CF.5–CF.6 exist, FS.1 must include: _(We use `defer` here. Defer ensures `file.Close()` runs when the function returns, even if an error occurs earlier. Defer is covered in full in Lesson CF.5.)_

---

### Section 06 — Backend, APIs & Databases

**Folder:** `06-backend-db/`
**ID prefixes:** `HS`, `API`, `DB`

**This section currently has NO HTTP or gRPC content. That must change.** This is the section that teaches learners to build production-grade backend services.

#### HTTP Servers (HS) — ENTIRELY NEW

| ID    | Lesson                       | Engineering Detail                     |
| ----- | ---------------------------- | -------------------------------------- |
| HS.1  | `net/http` basics            | Handler interface, mux, ListenAndServe |
| HS.2  | Routing patterns             | ServeMux, method checks                |
| HS.3  | Middleware — the pattern     | Chaining, wrapping, recovery           |
| HS.4  | Request parsing & validation | Body limits, header parsing            |
| HS.5  | Response writing patterns    | Status codes, JSON responses           |
| HS.6  | Error handling middleware    | Centralised error responses            |
| HS.7  | Server timeouts              | ✓ Critical — all four timeout types    |
| HS.8  | Graceful HTTP shutdown       | Drain pattern                          |
| HS.9  | Health & readiness probes    | Kubernetes-compatible endpoints        |
| HS.10 | REST API                     | ✓ Exercise — full middleware stack     |

**⚠️ In Production note for HS.7 (required):** The four timeout types are ALL required. Omitting any one creates a specific vulnerability. ReadHeaderTimeout prevents Slowloris. ReadTimeout prevents slow body attacks. WriteTimeout prevents slow response consumption. IdleTimeout prevents connection pool exhaustion. Show all four and explain each attack vector.

**Cross-reference note for HS.3:** _(Middleware uses closures to wrap handlers. Closures are covered in FE.9.)_
**Cross-reference note for HS.8:** _(This teaches HTTP-layer graceful shutdown. Full process-level shutdown is covered in GS.2 in Section 10. They work together.)_

#### APIs & Protocols (API) — ENTIRELY NEW

| ID    | Lesson                       | Engineering Detail            |
| ----- | ---------------------------- | ----------------------------- |
| API.1 | REST design principles       | Resource naming, HTTP verbs   |
| API.2 | API versioning strategies    | URL vs header versioning      |
| API.3 | Pagination & filtering       | Cursor vs offset              |
| API.4 | Protobuf basics              | Schema-first design, .proto   |
| API.5 | gRPC fundamentals            | Service definition, codegen   |
| API.6 | gRPC streaming               | Server, client, bidirectional |
| API.7 | gRPC interceptors            | Auth, logging, recovery       |
| API.8 | REST vs gRPC — the trade-off | When to choose each           |
| API.9 | gRPC Service                 | ✓ Exercise                    |

**Cross-reference note for API.4:** _(Protobuf uses code generation. The generated Go code uses `go generate`. Code generation is covered in depth in Lesson CG.1 in Section 10. Here we run `protoc` directly and use the output.)_

#### Databases (DB)

| ID   | Lesson                     | Status  |
| ---- | -------------------------- | ------- | ---------- |
| DB.1 | Connecting to SQLite       | exists  |
| DB.2 | Executing queries (INSERT) | exists  |
| DB.3 | Select queries             | exists  |
| DB.4 | Prepared statements        | exists  |
| DB.5 | Transactions               | exists  |
| DB.6 | Repository Pattern         | exists  | ✓ Exercise |
| DB.7 | N+1 query detection        | **new** |
| DB.8 | Query timeouts via context | **new** |

**Cross-reference note for DB.3:** _(This lesson uses `defer rows.Close()`. Defer is taught in CF.5. If you haven't completed CF.5 yet, read `defer` as "run this line when the function returns, regardless of errors.")_
**Cross-reference note for DB.8:** _(Uses `context.WithTimeout`. Context is covered fully in CT.1–CT.5 in Section 07.)_

---

### Section 07 — Concurrency

**Folder:** `07-concurrency/`
**ID prefixes:** `GC`, `SY`, `CT`, `TM`, `CP`

Teaches Go's concurrency model completely — goroutines, channels, synchronisation primitives, context, time, and advanced patterns.

#### Goroutines & Channels (GC)

| ID   | Lesson                | Status |
| ---- | --------------------- | ------ | ---------- |
| GC.1 | Goroutines            | exists |
| GC.2 | WaitGroups            | exists |
| GC.3 | Channels              | exists |
| GC.4 | Buffered channels     | exists |
| GC.5 | Closing channels      | exists |
| GC.6 | Pipeline project      | exists |
| GC.7 | Concurrent Downloader | exists | ✓ Exercise |

**Cross-reference note for GC.1:** _(Goroutines use closures to capture variables. The closure capture bug covered in COMMON-MISTAKES #1 is explained in Lesson FE.9.)_

#### Sync Primitives (SY) — MOSTLY NEW

| ID   | Lesson               | Status  | Note                                |
| ---- | -------------------- | ------- | ----------------------------------- |
| SY.1 | sync.Mutex & RWMutex | **new** | When to use Mutex vs channel        |
| SY.2 | sync.Once & sync.Map | **new** | Singleton patterns, concurrent maps |
| SY.3 | Atomic operations    | **new** | `sync/atomic`, CAS operations       |
| SY.4 | Race conditions      | **new** | `-race` detection, reading reports  |
| SY.5 | Goroutine leaks      | **new** | Detection, goroutine lifetime       |
| SY.6 | Deadlocks            | **new** | Diagnosis, prevention               |

**Why SY is new and required:** COMMON-MISTAKES mistake #14 ("concurrent map read/write") shows a Mutex fix — but there is currently no Mutex lesson. Learners cannot fix what they haven't learned. SY.1–SY.3 must precede the advanced patterns section (CP).

#### Context & Cancellation (CT)

| ID   | Lesson                     | Status |
| ---- | -------------------------- | ------ | ---------- |
| CT.1 | Background & TODO          | exists |
| CT.2 | WithCancel                 | exists |
| CT.3 | WithTimeout & WithDeadline | exists |
| CT.4 | WithValue                  | exists |
| CT.5 | Timeout-Aware API Client   | exists | ✓ Exercise |

#### Time & Scheduling (TM)

| ID   | Lesson           | Status |
| ---- | ---------------- | ------ | ---------- |
| TM.1 | Time basics      | exists |
| TM.2 | Time formatting  | exists |
| TM.3 | Timers & tickers | exists |
| TM.7 | Console Reminder | exists | ✓ Exercise |

**Note on TM placement:** Time belongs in Section 07 (not Section 05) because timers communicate via channels and tickers require goroutine lifecycle management. These are concurrency tools.

#### Advanced Concurrency Patterns (CP)

| ID   | Lesson                | Status |
| ---- | --------------------- | ------ | ---------- |
| CP.1 | errgroup basics       | exists |
| CP.2 | errgroup with context | exists |
| CP.3 | sync.Pool             | exists |
| CP.4 | Bounded Pipeline      | exists | ✓ Exercise |
| CP.5 | URL Health Checker    | exists | ✓ Exercise |

---

### Section 08 — Quality & Testing

**Folder:** `08-quality-test/`
**ID prefixes:** `TE`, `PR`

Teaches testing as a first-class engineering discipline. Tests are proof of behaviour, not an afterthought.

#### Testing (TE)

| ID    | Lesson                     | Status                         |
| ----- | -------------------------- | ------------------------------ |
| TE.1  | Unit testing basics        | exists                         |
| TE.2  | Table-driven tests         | exists                         |
| TE.3  | HTTP handler testing       | exists                         |
| TE.4  | Benchmarking               | exists                         |
| TE.5  | Sub-tests & t.Cleanup      | **new**                        |
| TE.6  | Fuzz testing               | **new** — Go 1.18+ standard    |
| TE.7  | Interfaces for testability | **new** — designing test seams |
| TE.8  | Mocking with interfaces    | **new**                        |
| TE.9  | Integration tests          | **new**                        |
| TE.10 | Golden files               | **new**                        |

#### Profiling (PR)

| ID   | Lesson               | Status                               |
| ---- | -------------------- | ------------------------------------ |
| PR.1 | CPU profiling        | exists                               |
| PR.2 | Live pprof endpoint  | exists                               |
| PR.3 | Memory profiling     | **new** — heap profiles, GC pressure |
| PR.4 | Escape analysis      | **new** — which data goes to heap?   |
| PR.5 | Benchmark-driven dev | **new** — before/after optimisation  |

---

## Phase 3: Systems Engineering

### Section 09 — Architecture & Security

**Folder:** `09-architecture/`
**ID prefixes:** `PD`, `ARCH`, `SEC`

**This section currently has only 3 lessons (PD.1–PD.3). It needs major expansion.** Architecture thinking and security engineering are what distinguish a senior engineer from a mid-level one.

#### Package Design (PD)

| ID   | Lesson              | Status |
| ---- | ------------------- | ------ |
| PD.1 | Naming conventions  | exists |
| PD.2 | Visibility & export | exists |
| PD.3 | Project layout      | exists |

#### Architecture Patterns (ARCH) — ENTIRELY NEW

| ID     | Lesson                                        | Engineering Detail                        |
| ------ | --------------------------------------------- | ----------------------------------------- |
| ARCH.1 | Monolith vs Modular Monolith vs Microservices | Trade-off matrix, team size               |
| ARCH.2 | Domain-Driven Design basics                   | Bounded contexts, ubiquitous language     |
| ARCH.3 | Hexagonal architecture in Go                  | Ports & adapters pattern                  |
| ARCH.4 | Repository pattern — deep dive                | Data access boundaries, interfaces        |
| ARCH.5 | Service layer pattern                         | Orchestration, thin handlers              |
| ARCH.6 | Event-driven architecture                     | Pub/sub, outbox pattern                   |
| ARCH.7 | CQRS basics                                   | Read/write separation, when it's worth it |
| ARCH.8 | When to split services                        | Warning signs, cost of splitting          |
| ARCH.9 | Modular Refactor                              | ✓ Capstone Exercise                       |

**Architecture verdict taught in ARCH.1:** The correct default is a Modular Monolith. Learners should finish ARCH.9 knowing how to start there, and how to recognise when splitting is justified — not that microservices are always the goal.

#### Security Engineering (SEC) — ENTIRELY NEW

| ID     | Lesson                       | Engineering Detail                          |
| ------ | ---------------------------- | ------------------------------------------- |
| SEC.1  | Input validation patterns    | Allow-list, fail-fast                       |
| SEC.2  | SQL injection prevention     | Parameterised queries only                  |
| SEC.3  | XSS & CSRF                   | Template escaping, CSRF tokens              |
| SEC.4  | Authentication basics        | Session vs token models                     |
| SEC.5  | JWT — implementation & risks | `alg:none` attack, signing                  |
| SEC.6  | Password hashing             | `bcrypt`, `argon2id`                        |
| SEC.7  | Rate limiting patterns       | Token bucket, sliding window                |
| SEC.8  | TLS & HTTPS in Go            | TLS config pitfalls                         |
| SEC.9  | Secrets management           | Env vars, vault patterns, never log secrets |
| SEC.10 | OWASP Top 10 for Go          | Applied checklist                           |
| SEC.11 | Secure API                   | ✓ Exercise                                  |

**Cross-reference note for SEC.5:** _(JWT is used in GoScale Module 3. This lesson gives you the foundation to implement it correctly. See also `docs/ENGINEERING_ERROR_FRAMEWORK.md` for how authentication failures should be categorised as UserError, not SystemError.)_

**Cross-reference note for SEC.6:** _(GoScale Module 3 uses `bcrypt` for password hashing. Complete this lesson before starting the flagship.)_

---

### Section 10 — Production Operations

**Folder:** `10-production/`
**ID prefixes:** `SL`, `GS`, `CFG`, `OPS`, `DOCKER`, `DEPLOY`, `CG`

Teaches everything needed to deploy, observe, and operate a Go service in production. This section is the bridge between "code that works" and "systems that run."

#### Structured Logging (SL)

| ID   | Lesson               | Status |
| ---- | -------------------- | ------ | ---------- |
| SL.1 | slog basics          | exists |
| SL.2 | Context-keyed logger | exists |
| SL.3 | Custom slog handler  | exists |
| SL.4 | zerolog comparison   | exists |
| SL.5 | PII Redactor         | exists | ✓ Exercise |

#### Graceful Shutdown (GS)

| ID   | Lesson               | Status |
| ---- | -------------------- | ------ | ---------- |
| GS.1 | signal.NotifyContext | exists |
| GS.2 | HTTP graceful drain  | exists |
| GS.3 | Shutdown Capstone    | exists | ✓ Capstone |

#### Configuration (CFG) — ENTIRELY NEW

| ID    | Lesson                    | Note                                 |
| ----- | ------------------------- | ------------------------------------ |
| CFG.1 | Environment variables     | `os.Getenv`, `.env` patterns         |
| CFG.2 | Configuration files       | JSON/TOML/YAML loading               |
| CFG.3 | Flag parsing              | `flag` package, precedence           |
| CFG.4 | 12-Factor App principles  | Production config discipline         |
| CFG.5 | Config validation on boot | Fail fast at startup, not at runtime |

**Cross-reference note for CFG.5:** _(GoScale Module 1 requires config validation on boot. This lesson teaches the pattern. Complete it before starting the flagship.)_

#### Observability (OPS) — ENTIRELY NEW

| ID    | Lesson                     | Note                              |
| ----- | -------------------------- | --------------------------------- |
| OPS.1 | Metrics basics             | What to measure, cardinality      |
| OPS.2 | Prometheus integration     | Counters, gauges, histograms      |
| OPS.3 | Distributed tracing basics | Correlation IDs, span propagation |
| OPS.4 | Feature flags              | Gradual rollout, kill switch      |
| OPS.5 | Alerting mindset           | SLOs, runbooks, on-call hygiene   |

**Cross-reference note for OPS.1–OPS.2:** _(GoScale Module 10 requires Prometheus metrics. Complete OPS.1–OPS.2 before starting Module 10.)_

#### Deployment (DOCKER/DEPLOY) — ENTIRELY NEW

| ID       | Lesson                | Note                      |
| -------- | --------------------- | ------------------------- |
| DOCKER.1 | Docker basics         | Dockerfile, image layers  |
| DOCKER.2 | Multi-stage builds    | Minimal production images |
| DOCKER.3 | Docker Compose        | Dev and prod profiles     |
| DEPLOY.1 | CI/CD pipelines       | GitHub Actions example    |
| DEPLOY.2 | Blue/green & rollback | Zero-downtime deployment  |
| DEPLOY.3 | Dockerised Service    | ✓ Exercise                |

**Cross-reference note for DOCKER.2:** _(Go produces self-contained binaries — this is what makes multi-stage Docker builds so effective. You learned why in Lesson GT.3.)_

#### Code Generation (CG)

| ID   | Lesson             | Status | Note               |
| ---- | ------------------ | ------ | ------------------ |
| CG.1 | go generate primer | exists | Moved from old s11 |
| CG.2 | Mockery workflow   | exists | Moved from old s11 |
| CG.3 | sqlc workflow      | exists | Moved from old s11 |

**Why CG is in Production Operations:** `go generate` is a build-time tool, not a language concept. It belongs alongside CI/CD and deployment tooling, not alongside language basics.

---

## Phase 4: Flagship Project

### Section 11 — GoScale SaaS Backend

**Folder:** `11-flagship/`

A production-grade, multi-tenant SaaS backend that integrates every concept from the curriculum.

**All prerequisites:** Complete all capstone exercises in s00–s10 before starting.

| Module | Focus                          | Applies From                         |
| ------ | ------------------------------ | ------------------------------------ |
| 1      | Foundation & Configuration     | CFG.1–CFG.5                          |
| 2      | Database & Models              | DB.1–DB.8, ARCH.4                    |
| 3      | Authentication & Authorization | SEC.4–SEC.6, ARCH.5                  |
| 4      | HTTP API                       | HS.1–HS.10, API.1–API.3              |
| 5      | Order Processing               | GC.1–GC.7, CT.1–CT.5                 |
| 6      | Payment Pipeline (mock)        | CP.1–CP.5, GS.1–GS.3                 |
| 7      | gRPC Internal API              | API.4–API.9                          |
| 8      | Event System & Workers         | CP.4, GS.3, SY.1–SY.3                |
| 9      | Caching                        | CT.3, SY.2                           |
| 10     | Observability                  | SL.1–SL.5, OPS.1–OPS.3               |
| 11     | Security Hardening             | SEC.1–SEC.10                         |
| 12     | Deployment & Operations        | DOCKER.1–DOCKER.3, DEPLOY.1–DEPLOY.2 |

---

## Milestone Progression

| %    | Milestone                 | Lesson   |
| ---- | ------------------------- | -------- |
| 5%   | Machine model checkpoint  | HC.8     |
| 10%  | First program             | GT.2     |
| 13%  | Toolchain confident       | GT.6     |
| 18%  | Pricing Checkout          | CF.7     |
| 24%  | Contact Directory         | DS.6     |
| 30%  | Order Summary             | FE.7     |
| 44%  | Payroll Processor         | TI.10    |
| 51%  | Config Parser (strings)   | ST.6     |
| 58%  | Log Search CLI            | FS.7     |
| 66%  | REST API with timeouts    | HS.10    |
| 70%  | gRPC Service              | API.9    |
| 74%  | Repository Pattern        | DB.6     |
| 77%  | Concurrent Downloader     | GC.7     |
| 81%  | URL Health Checker        | CP.5     |
| 85%  | Benchmark optimisation    | PR.5     |
| 88%  | Modular Refactor Capstone | ARCH.9   |
| 91%  | Secure API                | SEC.11   |
| 94%  | Dockerised Service        | DEPLOY.3 |
| 96%  | Shutdown Capstone         | GS.3     |
| 100% | GoScale Complete          | s11      |

---

## Lesson README Contract

Every learner-facing lesson must include these sections **in this order:**

```markdown
## Mission

## Prerequisites ← list lesson IDs (e.g., "- CF.1, CF.2")

## Mental Model

## Visual Model ← Mermaid diagram (NOT ASCII art)

## Machine View

## Run Instructions ← exact go run or go test command

## Code Walkthrough

## Try It ← numbered steps the learner actually does

## ⚠️ In Production ← required in EVERY lesson

## 🤔 Thinking Questions ← minimum 3 per lesson

## Next Step
```

For exercises: replace `## Code Walkthrough` with `## Solution Walkthrough`, add `## Verification Surface`.

**Validator enforcement:** The validator (`scripts/internal/curriculumvalidator/validator.go`) checks for these headings. Lessons that omit them will fail CI. New lessons must include all required headings.

---

## Cross-Reference Rules

When a lesson uses a concept not yet formally taught, it must include one of these forms:

**Forward reference (concept taught later):**

> _(We use `defer` here without fully explaining it. Defer is taught in Lesson CF.5. For now, read it as "run this line when the function returns, regardless of what happened.")_

**Backward reference (concept taught earlier):**

> _(This builds on error wrapping from FE.4. `fmt.Errorf("query user: %w", err)` preserves the full error chain.)_

**Sibling reference (same section, different subsystem):**

> _(Context cancellation via `ctx context.Context` is covered in CT.1–CT.5. Here we receive and respect it; there we learn to create it.)_

---

## Error Framework

Every section that touches errors must be consistent with `docs/ENGINEERING_ERROR_FRAMEWORK.md`.

| Layer      | Error Type  | HTTP Status | Action           |
| ---------- | ----------- | ----------- | ---------------- |
| Validation | UserError   | 400         | Return to caller |
| System     | SystemError | 500         | Wrap + propagate |
| Bug        | FatalError  | n/a         | Log + exit       |

Framework taught progressively:

- FE.4: errors as values (the (value, error) contract)
- TI.8: custom error types (implement the error interface)
- FE.5: validation as UserError
- HS.6: centralised error middleware
- ARCH.5: full three-tier framework in service layer
- GoScale Modules 3–4: applied in production

---

## Fixes Required Immediately

These are blocking issues that must be resolved before the curriculum is coherent:

| Priority | Fix                                                          | File(s) to change                | Status |
| -------- | ------------------------------------------------------------ | -------------------------------- | ------ |
| 1        | Backfill foundations README contracts across `s00`–`s04`     | lesson README files + validator  | ✅     |
| 2        | Add GT.5, GT.6 to curriculum.v2.json                         | `curriculum.v2.json`             | 📋     |
| 3        | Create CF.5, CF.6 lessons; add to curriculum.v2.json         | new files + `curriculum.v2.json` | 📋     |
| 4        | Add cross-reference note to FS.1 and DB.3 for `defer`        | lesson README files              | 📋     |
| 5        | Create HS.1–HS.10; add to curriculum.v2.json                 | new files + `curriculum.v2.json` | 📋     |
| 6        | Create API.1–API.9; add to curriculum.v2.json                | new files + `curriculum.v2.json` | 📋     |
| 7        | Create SY.1–SY.6; add to curriculum.v2.json                  | new files + `curriculum.v2.json` | 📋     |
| 8        | Create SEC.1–SEC.11; add to curriculum.v2.json               | new files + `curriculum.v2.json` | 📋     |
| 9        | Create ARCH.1–ARCH.9; add to curriculum.v2.json              | new files + `curriculum.v2.json` | 📋     |
| 10       | Create CFG.1–CFG.5; add to curriculum.v2.json                | new files + `curriculum.v2.json` | 📋     |
| 11       | Create OPS.1–OPS.5; add to curriculum.v2.json                | new files + `curriculum.v2.json` | 📋     |
| 12       | Create DOCKER.1–DEPLOY.3; add to curriculum.v2.json          | new files + `curriculum.v2.json` | 📋     |
| 13       | Fix all "Covered in:" references in COMMON-MISTAKES.md       | `COMMON-MISTAKES.md`             | ✅     |
| 14       | Update ROADMAP.md to match this document                     | `ROADMAP.md`                     | ✅     |
| 15       | Update PROGRESSION.md to match this document                 | `PROGRESSION.md`                 | ✅     |
| 16       | Fix CODE-STANDARDS.md constant naming (UPPER_SNAKE is wrong) | `CODE-STANDARDS.md`              | ✅     |

---

_The Go Engineer. Not a tutorial. An engineering transformation system._
