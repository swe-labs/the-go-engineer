# V2 Beta Architecture Decision

Status: Draft 0.1
Audience: Maintainers and curriculum designers
Related issues: `#171`, `#172`, `#173`, `#174`, `#175`

## Purpose

This document freezes the architectural direction for the v2 beta phase.

Alpha proved that the repo can migrate section by section without losing validation discipline.
Beta must do something bigger: turn the repo from a migrated topic inventory into a real
software engineering learning system.

This is not an "alpha plus polish" document.
It is the decision that beta is a learner-facing architecture redesign.

## Decision Summary

The beta line will follow these rules:

1. The current alpha `01` through `15` sections are a strong content inventory, but they are not the
   final learner-facing curriculum architecture.
2. The public beta curriculum will be organized by engineering maturity and learning purpose rather
   than by the current alpha section numbering alone.
3. Beta must add four layers that alpha does not yet represent well enough:
   - `0 Foundation`
   - `Production Engineering`
   - `Expert Layer`
   - `Flagship Project`
4. Beta must treat exercises, starter-solution pairs, and evaluation criteria as part of the product,
   not as optional extras.
5. Breaking changes are allowed when they improve the learning system materially.

## Why Alpha Is Not Enough

Alpha succeeded at:

- migrating the current curriculum into `curriculum.v2.json`
- giving each migrated section a cleaner README and milestone surface
- proving validator-backed curriculum maintenance
- proving that the repo can evolve safely in public

Alpha did not yet solve the bigger product problem:

- the learner-facing system is still mostly section-first
- the beginner on-ramp is still not deep enough
- the exercise bank is still partial and uneven
- production and expert layers are still under-defined
- there is no flagship project spine

In short:

Alpha proved migration discipline.
Beta must prove learning-system architecture.

## Beta Target Architecture

The beta target is a staged curriculum with these public layers:

### 0. Foundation

Purpose:

- introduce computers, execution, memory, CLI, and tool basics before Go syntax takes over

Current alpha source:

- the `getting-started` track inside Section `01`
- new beta-only mental-model and pre-Go support content

### 1. Language Fundamentals

Purpose:

- move learners from first syntax to values, control flow, data structures, functions, and errors

Current alpha source:

- the `language-basics` track inside Section `01`
- Sections `02` to `04`

### 2. Types And Design

Purpose:

- teach structs, methods, interfaces, composition, and text/data modeling as design tools

Current alpha source:

- Sections `05` to `07`

### 3. Modules And IO

Purpose:

- teach packages, modules, CLI boundaries, encoding, and filesystem work as practical systems input/output

Current alpha source:

- Sections `08` to `09`

### 4. Backend Engineering

Purpose:

- teach HTTP, persistence, transactions, handlers, and application boundaries

Current alpha source:

- Section `10`

### 5. Concurrency System

Purpose:

- teach concurrency from mental model through patterns, failure modes, and bounded production use

Current alpha source:

- Sections `11` to `12`

### 6. Quality And Performance

Purpose:

- teach testing, benchmarks, profiling, and disciplined measurement

Current alpha source:

- Section `13`

### 7. Architecture

Purpose:

- teach package boundaries, logging, shutdown, and service-level structural thinking

Current alpha source:

- the package and service-structure portions of Section `14`

### 8. Production Engineering

Purpose:

- add deployment, configuration, tracing, monitoring, and operating concerns that alpha does not yet
  represent as a first-class layer

Current alpha source:

- the operational tracks already present inside Section `14`, especially structured logging and
  graceful shutdown
- new beta-only production topics such as deployment, config, tracing, and monitoring

### 9. Expert Layer

Purpose:

- add code review, anti-pattern analysis, failure analysis, and system-design-style engineering pressure

Current state:

- new beta layer

### 10. Flagship Project

Purpose:

- give the learner one large proof-of-skill project that pulls the curriculum together

Current state:

- new beta layer

### 11. Code Generation

Purpose:

- teach generation and workflow leverage after the learner already understands the code being generated

Current alpha source:

- Section `15`

## Alpha-To-Beta Mapping

| Alpha source | Beta destination | Decision |
| --- | --- | --- |
| `01-core-foundations/getting-started` | `0 Foundation` | split and extend |
| `01-core-foundations/language-basics` | `1 Language Fundamentals` | regroup |
| `02-control-flow` | `1 Language Fundamentals` | regroup |
| `03-data-structures` | `1 Language Fundamentals` | regroup |
| `04-functions-and-errors` | `1 Language Fundamentals` | regroup |
| `05-types-and-interfaces` | `2 Types And Design` | regroup |
| `06-composition` | `2 Types And Design` | regroup |
| `07-strings-and-text` | `2 Types And Design` | regroup |
| `08-modules-and-packages` | `3 Modules And IO` | regroup |
| `09-io-and-cli` | `3 Modules And IO` | regroup |
| `10-web-and-database` | `4 Backend Engineering` | regroup |
| `11-concurrency` | `5 Concurrency System` | regroup |
| `12-concurrency-patterns` | `5 Concurrency System` | regroup |
| `13-quality-and-performance` | `6 Quality And Performance` | regroup |
| `14-application-architecture/package-design` | `7 Architecture` | split and regroup |
| `14-application-architecture/structured-logging` | `8 Production Engineering` | split and regroup |
| `14-application-architecture/graceful-shutdown` | `8 Production Engineering` | split and regroup |
| other beta-only production topics | `8 Production Engineering` | add |
| none | `9 Expert Layer` | add |
| none | `10 Flagship Project` | add |
| `15-code-generation` | `11 Code Generation` | keep as late-stage specialization |

## Non-Negotiable Beta Requirements

Beta must include:

- a real `0 Foundation` layer
- a public stage-based learner architecture
- a mandatory mental-model surface in every section or stage
- a canonical lesson standard
- a canonical exercise standard
- starter and solution expectations
- a reusable evaluation and rubric pattern
- a mini-project ladder
- a flagship project plan and initial scaffold
- a production-engineering layer

## Mental-Model Rule

Every beta section or stage must include a mental-model explanation.

That does not mean every section needs a separate `mental-model/` folder.

Use this rule instead:

- simple sections may satisfy the requirement through the stage README, section README, and first
  lesson
- complex sections should have a dedicated mental-model sublayer

Complex sections usually include:

- concurrency
- backend engineering
- architecture
- performance and profiling
- production engineering
- expert-layer analysis

Practical takeaway:

- mental model is mandatory everywhere
- a separate mental-model folder is mandatory only where abstraction, failure modes, or systems
  intuition make it necessary

## What Beta Does Not Need To Ship Immediately

Beta does not need:

- a separate LMS or platform
- automatic scoring infrastructure
- complete easy-medium-hard-expert exercise coverage for every stage on day one
- a perfect physical folder rename for the entire repo in one massive cutover

Those items are useful, but they are not the gating definition of beta.

## Source Content Versus Public Architecture

The current alpha sections remain valuable source content.

That means:

- beta should reuse and regroup alpha content rather than discard it
- the public learner-facing architecture should change even if some internal folder movement remains phased
- the stage model becomes the main navigation truth during beta work
- some alpha sections will intentionally split across more than one beta stage when that reflects the
  real learner-facing boundary better than keeping whole-section moves

Practical rule:

- logical learner architecture first
- physical directory migration second, only when it clearly improves the system

## Beta Implementation Rules

1. Do not treat the alpha section structure as sacred.
2. Do not break repo validation discipline while redesigning the curriculum.
3. Do not defer `0 Foundation`, the exercise system, or the flagship project until "later."
4. Do not let platform ideas delay curriculum architecture work.
5. Prefer one coherent beta system over many isolated local improvements.

## First Beta Implementation Order

1. freeze the beta stage model and mapping
2. define the global exercise, starter-solution, and rubric system
3. define `0 Foundation`
4. define `Production Engineering`
5. define `Expert Layer`
6. define the `Flagship Project`
7. update root docs and learning-path navigation around the beta architecture
8. begin regrouping alpha content into beta-facing stage surfaces

## Immediate Follow-Up Decisions

The next beta planning questions to resolve are:

1. how much of the beta architecture should be reflected in physical folders immediately
2. what the first version of `0 Foundation` must contain
3. what the flagship project should be
4. what exercise depth is mandatory for beta versus post-beta

## Bottom Line

Beta is the point where v2 stops being "a cleaner set of migrated sections" and becomes a real
software engineering learning system.

That means the beta success bar is architectural, not cosmetic.
