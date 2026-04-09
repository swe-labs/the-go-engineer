# V2 Beta Layer Definitions

Status: Draft 0.1
Audience: Maintainers and curriculum designers
Related issues: `#171`, `#175`

## Purpose

This document defines the four beta layers that most clearly separate beta from alpha:

- `0 Foundation`
- `Production Engineering`
- `Expert Layer`
- `Flagship Project`

These layers are not optional polish.
They are the parts of the curriculum that make the repo behave like a zero-to-engineer software
engineering learning system instead of a strong topic inventory.

This document is the planning authority for `#175`.

## Governing Rule

Each new beta layer must answer four questions clearly:

1. What learner problem does this layer solve?
2. What outcomes does it own?
3. What does it explicitly not try to become?
4. Where does it sit in the learner journey?

## 0 Foundation

### Problem it solves

Alpha still assumes too much beginner readiness.
Learners who are new to programming, terminals, files, editors, and execution flow need a real
on-ramp before heavy syntax and engineering structure begin.

### Core outcomes

`0 Foundation` should teach:

- what a program is
- how code becomes execution
- what files, folders, and the terminal are doing
- how to install, run, and verify Go tools without fear
- the first mental models for memory, state, and input-output boundaries

### What it is not

`0 Foundation` is not:

- a full CS-prep program
- a math curriculum
- a long detour before writing any Go

It should shorten confusion, not delay progress.

### Placement

`0 Foundation` sits before `1 Language Fundamentals`.

### Implementation guidance

- use `01-core-foundations/getting-started` as the primary alpha source
- add new mental-model and workflow support content where alpha is too thin
- keep the first successful run and first debugging loop inside this layer

## Production Engineering

### Problem it solves

Alpha teaches parts of operations and runtime thinking, but it does not yet present them as a
first-class engineering layer.

Learners need a stage that answers:

- how software is configured
- how software is observed
- how software is shut down safely
- how software is deployed and run

### Core outcomes

`Production Engineering` should teach:

- structured logging
- graceful shutdown
- configuration boundaries
- deployment packaging
- container-aware operation
- observability basics
- runtime and support checklists

### What it is not

`Production Engineering` is not:

- a full DevOps certification
- a Kubernetes-heavy platform track
- a replacement for architecture or backend design

It should stay application-centric and repo-native.

### Placement

`Production Engineering` follows `7 Architecture` and comes before the most advanced project and
expert-pressure work.

### Implementation guidance

- split `14-application-architecture` into architecture-owned and operations-owned tracks
- reuse `structured-logging`, `graceful-shutdown`, and `docker-and-deployment`
- add tracing, monitoring, and config surfaces only where they support the flagship system honestly

## Expert Layer

### Problem it solves

Alpha builds content breadth, but it does not yet create enough explicit engineering pressure for
learners who need to move from "can build" to "can reason, review, diagnose, and defend trade-offs."

### Core outcomes

`Expert Layer` should teach:

- review and critique
- anti-pattern detection
- failure analysis
- refactoring under constraints
- boundary and trade-off reasoning
- reading and judging unfamiliar code

### What it is not

`Expert Layer` is not:

- extra syntax content
- prestige labeling
- vague "senior tips" without evidence

It must create pressure through concrete analysis and better judgment.

### Placement

`Expert Layer` should sit late in the beta journey and feed directly into flagship work and final
proof surfaces.

### Implementation guidance

- seed it from review and debugging slices derived from earlier stages
- prefer rubric-heavy work over toy correctness-only tasks
- use code review, diagnosis, and redesign exercises as the main format

## Flagship Project

### Problem it solves

Alpha has meaningful milestone projects, but it does not yet have one product spine that proves
integrated engineering growth over time.

### Core outcomes

`Flagship Project` should prove:

- staged product growth
- architectural decisions
- testing and quality discipline
- operational readiness
- refactoring and extension under real constraints

### What it is not

`Flagship Project` is not:

- a random extra project folder
- a final dump of features after the real curriculum is over
- a second repo that drifts away from the learning system

It must be a curriculum spine.

### Placement

`Flagship Project` should begin once learners have enough backend and concurrency depth to build
meaningful features, then continue through production and expert work.

### Implementation guidance

- use `14-application-architecture/enterprise-capstone` as the strongest alpha seed
- break the flagship into staged milestones, not one giant final assignment
- connect flagship checkpoints to the exercise and rubric system directly

## Relationship Between The Four Layers

These layers are not independent extras.
They form the beta identity.

- `0 Foundation` makes honest beginner entry possible
- `Production Engineering` turns application code into operated software
- `Expert Layer` turns implementation into engineering judgment
- `Flagship Project` proves the learner can carry all of that in one system

If one of these layers is missing, beta is weaker in a way alpha polish cannot hide.

## Rollout Guidance

The recommended rollout order is:

1. define and ship `0 Foundation`
2. define and ship `Production Engineering`
3. define the `Flagship Project` skeleton together with the minimum `Expert Layer` pressure model
4. expand `Expert Layer` pressure across the later stages and flagship milestones

Why this order:

- learners need an honest entry before they need pressure
- production concepts need a real system context
- flagship work is strongest when it is not separated from the first real expert-pressure surfaces
- deeper expert pressure is strongest when learners already have something substantial to critique
  and improve

## Definition Of Done For The Layer-Definition Phase

`#175` is ready to close when:

- each new beta layer has a clear learner problem and owned outcomes
- the flagship project is positioned as a curriculum spine
- later implementation issues can use these layer definitions without re-arguing their purpose
