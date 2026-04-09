# V2 Beta Stage Mapping

Status: Draft 0.1
Audience: Maintainers and curriculum designers
Related issues: `#171`, `#173`, `#175`

## Purpose

This document turns the beta architecture decision into an executable mapping plan.

It answers three practical questions:

1. Which alpha sections and tracks feed each beta stage?
2. Which alpha sections must split across more than one beta stage?
3. Which breaking migration groups should land together so beta stays coherent for learners?

This document is the planning authority for `#173`.

## Governing Rules

The mapping work follows these rules:

1. Beta organizes the public learner experience by stage, not by alpha section numbering.
2. Alpha content remains the source inventory unless beta explicitly adds new material.
3. Alpha sections may split across more than one beta stage when the learner-facing boundary is
   clearer than the current physical section boundary.
4. Public stage navigation may change before the repo's physical directory layout fully changes.
5. `0 Foundation`, `Production Engineering`, `Expert Layer`, and `Flagship Project` are real beta
   layers, not placeholder labels.

## Canonical Stage Map

### 0. Foundation

Primary source:

- `01-core-foundations/getting-started`

Beta additions:

- pre-Go mental models for execution, files, terminals, and tooling
- first-run environment validation and workflow confidence checks

Role:

- give true beginners a stable on-ramp before language-heavy lessons begin

### 1. Language Fundamentals

Primary sources:

- `01-core-foundations/language-basics`
- `02-control-flow`
- `03-data-structures`
- `04-functions-and-errors`

Role:

- build first real programming fluency in Go
- cover values, branching, repetition, collections, function boundaries, and error handling

### 2. Types And Design

Primary sources:

- `05-types-and-interfaces`
- `06-composition`
- `07-strings-and-text`

Role:

- move learners from syntax into modeling and design choices

### 3. Modules And IO

Primary sources:

- `08-modules-and-packages`
- `09-io-and-cli/cli-tools`
- `09-io-and-cli/encoding`
- `09-io-and-cli/filesystem`

Role:

- teach packaging, file and process boundaries, serialization, and practical input/output work

### 4. Backend Engineering

Primary sources:

- `10-web-and-database/http-client`
- `10-web-and-database/web-masterclass`
- `10-web-and-database/databases`
- `10-web-and-database/database-migrations`

Role:

- teach HTTP, persistence, handlers, request flow, and application data boundaries

### 5. Concurrency System

Primary sources:

- `11-concurrency/context`
- `11-concurrency/goroutines`
- `11-concurrency/time-and-scheduling`
- `12-concurrency-patterns`

Role:

- move from concurrency mental models into bounded real-world patterns and failure handling

### 6. Quality And Performance

Primary sources:

- `13-quality-and-performance/testing`
- `13-quality-and-performance/http-client-testing`
- `13-quality-and-performance/profiling`
- `13-quality-and-performance/4-testcontainers`

Role:

- teach how to verify, measure, and trust engineering work

### 7. Architecture

Primary sources:

- `14-application-architecture/package-design`
- `14-application-architecture/grpc`

Conditional source:

- architectural slices from `14-application-architecture/enterprise-capstone`

Role:

- teach package boundaries, service structure, and system-level code organization

### 8. Production Engineering

Primary sources:

- `14-application-architecture/structured-logging`
- `14-application-architecture/graceful-shutdown`
- `14-application-architecture/docker-and-deployment`

Beta additions:

- deployment, config, tracing, monitoring, and operating playbooks

Role:

- teach what it takes to run and support software, not just write it

### 9. Expert Layer

Primary source:

- new beta-only material

Expected seed material:

- review and anti-pattern slices derived from alpha lessons and projects

Role:

- create engineering pressure through critique, diagnosis, and trade-off analysis

### 10. Flagship Project

Primary sources:

- `14-application-architecture/enterprise-capstone`
- milestone mini-project patterns already proven in alpha

Beta additions:

- a clearer long-running product skeleton with staged project checkpoints

Role:

- give learners one large, integrated proof-of-skill path

### 11. Code Generation

Primary source:

- `15-code-generation`

Role:

- teach leverage only after learners understand the code and systems being generated

## Split Map For Mixed Alpha Sections

The following alpha sections must split during beta:

| Alpha source | Split destination | Why the split is required |
| --- | --- | --- |
| `01-core-foundations` | `0 Foundation` + `1 Language Fundamentals` | setup confidence and first syntax are different beginner problems |
| `09-io-and-cli` | `3 Modules And IO` only, but through three internal tracks | CLI, encoding, and filesystem should stay distinct even when grouped |
| `10-web-and-database` | `4 Backend Engineering` with track-aware internal grouping | HTTP, persistence, and migrations should stay visible as different concerns |
| `14-application-architecture` | `7 Architecture` + `8 Production Engineering` + `10 Flagship Project` | architecture, operations, and capstone work are different learner goals |

## Breaking Migration Groups

These groups are the units that should land together when beta work starts changing the public
learner architecture.

### Group A: Public Beta Shell

Must land together:

- root README beta positioning
- beta stage navigation entry points
- learning-path updates for stage-based routing
- clear "alpha section inventory vs beta stage architecture" messaging

Why:

- learners should never see a half-switched navigation model

### Group B: Foundation Extraction

Must land together:

- `0 Foundation` stage shell
- Section `01` split between `getting-started` and `language-basics`
- beginner mental-model docs and entry guidance

Why:

- beta cannot claim a real zero-to-engineer system without a beginner on-ramp

### Group C: Fundamentals Regroup

Must land together:

- stage surface for Language Fundamentals
- regrouped entry guidance for Sections `01` to `04`
- clear mental-model and milestone ownership across the stage

Why:

- this stage becomes the main early-learning backbone

### Group D: Design And IO Regroup

Must land together:

- `2 Types And Design` stage shell
- `3 Modules And IO` stage shell
- public routing from design topics into module and I/O work

Why:

- this is where learners shift from syntax to real software boundaries

### Group E: Backend And Concurrency Regroup

Must land together:

- `4 Backend Engineering` stage shell
- `5 Concurrency System` stage shell
- project and checkpoint expectations that bridge the two

Why:

- these are the first clearly application-level engineering stages

### Group F: Quality, Architecture, And Production Split

Must land together:

- `6 Quality And Performance`
- `7 Architecture`
- `8 Production Engineering`
- the explicit split of Section `14`

Why:

- this is the most breaking regroup in the beta design and needs one coherent public cutover

### Group G: Expert And Flagship Layers

Must land together:

- `9 Expert Layer`
- `10 Flagship Project`
- the project ladder handoff into flagship work

Why:

- flagship work without expert pressure becomes a capstone only, not a real seniority bridge

### Group H: Code Generation Re-entry

Must land together:

- `11 Code Generation`
- guidance that this stage sits after architectural and production understanding

Why:

- generation only makes sense after the learner understands what is being automated

## Recommended Beta Execution Order

1. Group A: Public Beta Shell
2. Group B: Foundation Extraction
3. Group C: Fundamentals Regroup
4. Group D: Design And IO Regroup
5. Group E: Backend And Concurrency Regroup
6. Group F: Quality, Architecture, And Production Split
7. Group G: Expert And Flagship Layers
8. Group H: Code Generation Re-entry

This order keeps the learner journey coherent from the front door outward.

## Do Not Do

During beta mapping work, do not:

- duplicate the same alpha content into multiple public beta stages without a clear ownership rule
- force physical directory moves before the public learner architecture is stable
- leave Section `14` unsplit just because the current folder is convenient
- claim `0 Foundation` exists before beginner tooling and mental-model support actually land

## Definition Of Done For The Mapping Phase

`#173` is ready to close when:

- every beta stage has an explicit alpha source or explicit beta-only addition plan
- every known mixed alpha section has an intentional split rule
- the breaking migration groups are named and ordered
- the beta implementation work can be broken into stage-aware issues without re-arguing the map
