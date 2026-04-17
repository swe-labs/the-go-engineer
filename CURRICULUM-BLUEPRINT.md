# The Go Engineer Curriculum Blueprint

## Purpose

This document defines how the 11-stage architecture should behave as a learning system.

It is not just a list of topics.
It is the curriculum contract for how we teach:

- from zero programming knowledge
- through Go fluency
- into software-engineering judgment

## Core Promise

The Go Engineer should help a learner move from:

- "I can copy code"

to:

- "I understand what this line does"
- "I can explain why this design exists"
- "I can predict what breaks"
- "I can build and operate a real system"

## Non-Negotiable Teaching Rules

### 1. Doc first, then code

Every learner-facing lesson teaches through `README.md` first.
The learner should understand the mission before opening `main.go`.

### 2. Code is never skipped

We do not replace code with prose.
We explain the code, then run the code, then modify the code.

### 3. Zero magic

Each stage teaches only what has been earned.
If a concept depends on later ideas, it belongs later or must be reduced to a clearly labeled
preview.

### 4. Explanation should answer how, why, and what changes

Good teaching surfaces explain:

- what this line or block does
- why it exists
- what would change if we changed it
- what mistake a learner is likely to make here

### 5. Engineering depth must be stage-aware

We do want:

- design thinking
- failure thinking
- production relevance
- debugging instincts

But we add them at the right layer.
We do not dump senior-level pressure framing into beginner lessons just to sound impressive.

## Canonical 11 Stages

| Stage | Goal | Current Source Surface | Required Proof Style |
| --- | --- | --- | --- |
| `01 Getting Started` | build environment confidence and first execution trust | [01-getting-started](./01-getting-started/) | first-run execution |
| `02 Language Basics` | build value, control-flow, and data-structure fluency | [02-language-basics](./02-language-basics/) | small runnable milestones |
| `03 Functions & Errors` | turn inline logic into reusable, honest function boundaries | [03-functions-errors](./03-functions-errors/) | starter + runnable + test-backed milestone |
| `04 Types & Design` | shape programs with types, interfaces, composition, and text workflows | [04-types-design](./04-types-design/), [05-composition](./05-composition/), [06-strings-and-text](./06-strings-and-text/) | multiple section milestones inside one stage |
| `05 Packages & IO` | make code navigable across modules and useful across files, encodings, and CLI surfaces | [07-modules-and-packages](./07-modules-and-packages/), [08-io-and-cli](./08-io-and-cli/) | track-based proof surfaces |
| `06 Backend & DB` | build HTTP and database behavior that feels like real application code | [09-web-and-database](./09-web-and-database/) | integrated backend milestone |
| `07 Concurrency` | teach coordination, cancellation, and bounded concurrent work | [10-concurrency](./10-concurrency/), [11-concurrency-patterns](./11-concurrency-patterns/) | runnable concurrency capstones |
| `08 Quality & Test` | verify, benchmark, and profile code instead of guessing | [12-quality-and-performance](./12-quality-and-performance/) | test and benchmark proof |
| `09 Architecture` | reason about package, service, and boundary design | [13-application-architecture/package-design](./13-application-architecture/package-design/), [13-application-architecture/grpc](./13-application-architecture/grpc/) | design-oriented milestones |
| `10 Production` | think about runtime behavior, observability, shutdown, and deployment | [13-application-architecture/structured-logging](./13-application-architecture/structured-logging/), [13-application-architecture/graceful-shutdown](./13-application-architecture/graceful-shutdown/), [13-application-architecture/docker-and-deployment](./13-application-architecture/docker-and-deployment/) | operational proof surfaces |
| `11 Flagship` | integrate multiple engineering layers into one project | [13-application-architecture/enterprise-capstone](./13-application-architecture/enterprise-capstone/), [14-code-generation](./14-code-generation/) | flagship checkpoints and integrated system proof |

## Stage-Level Blueprint

### Stages 01-04: Beginner-to-Builder

These stages must feel safe, explicit, and zero-magic.

Required elements:

- mission
- mental model
- visual model where useful
- machine view where useful
- literal or near-literal walkthroughs
- `Try It` prompts
- clean runnable code

Avoid:

- premature scale pressure
- advanced security catalogs
- concurrency failure systems
- abstract design jargon before the learner has concrete examples

### Stages 05-08: Builder-to-Engineer

These stages should start increasing engineering judgment.

Add more of:

- trade-off explanations
- failure cases
- safer defaults
- tests and verification surfaces
- performance and maintainability reasoning

### Stages 09-11: Engineer-to-System Thinker

These stages should carry heavier:

- architecture trade-offs
- production notes
- failure scenarios
- judgment-heavy exercises
- integrated project proof

## Canonical Lesson Contract

For learner-facing lessons, the default shape is:

```text
lesson/
├── README.md
├── main.go
├── _bad/
│   └── main.go            (optional)
├── _starter/
│   └── main.go            (optional)
└── main_test.go           (optional)
```

### Required README sections

Each lesson README should include, as appropriate:

- mission
- prerequisites
- mental model
- visual model
- machine view
- run instructions
- code walkthrough
- `Try It`
- next step

### Required source-file behavior

Source files should stay readable and runnable.
They should not become giant essays.

Use inline comments for:

- non-obvious behavior
- mutation or boundary traps
- subtle runtime implications

Do not use code headers as the main teaching surface.

## Canonical Milestone Contract

Every stage needs proof, not just lessons.

A milestone should usually provide:

- clear README instructions
- a runnable completed solution
- a starter surface when the learner is expected to implement
- tests when the behavior should be provable

## Stage Support Docs

Stage support docs under [docs/stages](./docs/stages/README.md) should do three jobs only:

1. explain stage purpose
2. route learners to the current source surfaces
3. explain proof expectations

They should not become replacements for lesson teaching.

## Advanced Overlays

We do want advanced thinking, failure, and production overlays.
We just apply them by stage.

### Early stages

Use only the local README-first contract with small, grounded production relevance.

### Later stages

Use the advanced templates in [docs/templates](./docs/templates/README.md) when they genuinely help:

- thinking questions
- failure scenarios
- production notes

## How To Add New Lessons Without Breaking The Architecture

If the curriculum needs more depth:

1. add the lesson inside an existing stage
2. keep the learner-facing stage count at 11
3. update stage docs and metadata honestly
4. add or update the stage proof surface if the change shifts readiness expectations

Examples:

- add another language-basics lesson inside Stage `02`
- add another backend lesson inside Stage `06`
- add another production lesson inside Stage `10`

Do not solve content growth by inventing a new public stage unless the entire architecture is being
reworked intentionally.

## Beta Completion Standard

We should treat the beta migration as complete only when:

- the 11-stage architecture is consistent across the repo
- the major learner-facing sections are migrated to the current teaching contract
- stale 15-stage and 17-stage public narratives are retired
- the validator supports the current contract
- the active learner path no longer depends on hidden legacy assumptions

## Bottom Line

The Go Engineer should feel like one coherent engineering learning system.

The 11 stages give us the public spine.
The README-first teaching contract gives us the delivery standard.
Future expansion should deepen the stages we have, not fragment the learner path again.
