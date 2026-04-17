# The Go Engineer Architecture

## Canonical Public Model

The Go Engineer now has one canonical learner-facing architecture:

1. `01 Getting Started`
2. `02 Language Basics`
3. `03 Functions & Errors`
4. `04 Types & Design`
5. `05 Packages & IO`
6. `06 Backend & DB`
7. `07 Concurrency`
8. `08 Quality & Test`
9. `09 Architecture`
10. `10 Production`
11. `11 Flagship`

This 11-stage model replaces the older 15-section and 17-section public narratives.

Those older layouts may still matter as internal history or source-inventory context, but they are
no longer the public navigation truth.

## Why The Stage Count Is 11

We want a curriculum that is:

- clear for beginners
- stable for maintainers
- extensible without constant public renumbering
- honest about how Go engineering skills compound over time

The old 15-section and 17-section models helped us reason about content inventory, but they created
too much public fragmentation.

The 11-stage model keeps the learner path simpler while still giving us room to grow depth inside
each stage.

## Stage Map

| Stage | Public Focus | Current Source Surface |
| --- | --- | --- |
| `01 Getting Started` | install Go, run code, understand files and terminal basics | [01-getting-started](./01-getting-started/) |
| `02 Language Basics` | values, control flow, collections, pointer-aware mutation | [02-language-basics](./02-language-basics/) |
| `03 Functions & Errors` | reusable logic, explicit failure handling, validation | [03-functions-errors](./03-functions-errors/) |
| `04 Types & Design` | structs, methods, interfaces, composition, text-heavy modeling | [04-types-design](./04-types-design/), [05-composition](./05-composition/), [06-strings-and-text](./06-strings-and-text/) |
| `05 Packages & IO` | modules, packages, filesystem work, encoding, CLI tooling | [07-modules-and-packages](./07-modules-and-packages/), [08-io-and-cli](./08-io-and-cli/) |
| `06 Backend & DB` | HTTP servers, handlers, clients, repositories, database workflows | [09-web-and-database](./09-web-and-database/) |
| `07 Concurrency` | goroutines, channels, context, concurrency patterns, scheduling | [10-concurrency](./10-concurrency/), [11-concurrency-patterns](./11-concurrency-patterns/) |
| `08 Quality & Test` | testing, mocks, benchmarks, profiling | [12-quality-and-performance](./12-quality-and-performance/) |
| `09 Architecture` | package design, service boundaries, gRPC, system structure | [13-application-architecture/package-design](./13-application-architecture/package-design/), [13-application-architecture/grpc](./13-application-architecture/grpc/) |
| `10 Production` | structured logging, graceful shutdown, deployment, runtime behavior | [13-application-architecture/structured-logging](./13-application-architecture/structured-logging/), [13-application-architecture/graceful-shutdown](./13-application-architecture/graceful-shutdown/), [13-application-architecture/docker-and-deployment](./13-application-architecture/docker-and-deployment/) |
| `11 Flagship` | integrated project work and late-stage system proof | [13-application-architecture/enterprise-capstone](./13-application-architecture/enterprise-capstone/), [14-code-generation](./14-code-generation/) |

## Progression Model

### Phase 1: Foundations

Stages `01` through `04` teach the beginner path:

- code can run
- values move
- control flow decides work
- collections hold data
- functions create boundaries
- types shape programs

### Phase 2: Engineering Core

Stages `05` through `08` turn working programs into engineering systems:

- module and package boundaries
- IO and CLI work
- backend and database flows
- concurrency
- verification and performance awareness

### Phase 3: System Design And Operations

Stages `09` through `11` teach:

- architecture and system boundaries
- runtime and operational behavior
- integrated flagship proof

## Rules For Future Growth

### 1. The public stage count stays at 11

If we need more depth, we add:

- lessons
- drills
- milestone variants
- optional stretch tracks

inside an existing stage.

We do not add a twelfth public stage just because a topic becomes richer.

### 2. Old section inventories are implementation detail, not public truth

The source tree may still contain historical lesson groupings.
That is acceptable as long as:

- the learner-facing docs stay aligned to the 11-stage model
- `curriculum.v2.json` routes learners correctly
- stage support docs explain the mapping honestly

### 3. Expanding a stage is allowed

If a stage needs more teaching surface to meet the repo goal, add it.

Examples:

- Stage `02` can gain another language-basics lesson if it removes hidden magic.
- Stage `05` can gain another IO lesson if it closes a real readiness gap.
- Stage `10` can gain another production lesson if operational behavior is still under-taught.

The constraint is:

- expand within the stage
- do not change the public stage architecture

### 4. Every stage must keep one honest proof surface

A stage is not complete because its lessons exist.
A stage is complete when it has:

- a clear learner mission
- a stable stage entry doc
- coherent lesson routing
- at least one meaningful proof surface

### 5. README-first teaching is part of the architecture

For learner-facing content, the teaching contract is:

1. `README.md` first
2. code second
3. explanation beside the code, not hidden elsewhere

That rule applies across the architecture, even though the density of explanation can change by
stage.

## What “Complete Beta Migration” Means

The beta migration is complete when all of these are true:

- the repo speaks one public 11-stage language
- learner routing matches the 11-stage model
- source paths and support docs map honestly to those stages
- the major learner-facing sections follow the current README-first contract
- stale public references to the old 15-stage or 17-stage models are removed or clearly demoted to history

## Bottom Line

The Go Engineer is no longer a repo with competing public architectures.

Its public curriculum architecture is the 11-stage model above.
Future lesson growth should deepen those stages, not replace them.
