# V2 QA And Validation Model

## Purpose

v2 needs a quality model that keeps the curriculum, docs, metadata, and runnable code aligned as
one system.

This document defines what should be checked automatically, what should be reviewed by humans, and
what should block migration work.

## Core Rule

Validation should enforce reality, not aspiration.

If the repo cannot actually prove a rule through tooling or consistent review, that rule should not
be treated like a hard guarantee.

v2 QA should prefer:

- a smaller set of trusted checks
- explicit ownership of each check
- type-aware validation tied to the curriculum model

over giant generic checklists that drift from the real repo.

## Current Baseline

The repo already has a useful quality baseline.

Today, the main guardrails are:

- `scripts/validate_curriculum.go` for curriculum paths and runnable path references
- CI checks for build, vet, formatting, module tidy, tests, race detection, coverage, benchmarks,
  and curriculum validation
- standards docs such as `CODE-STANDARDS.md` and `TESTING-STANDARDS.md`

This is a strong starting point.
The main v2 need is not "invent QA from scratch."
The need is to align QA with the richer content model, learning paths, and navigation system.

## QA Goals

The v2 QA system should answer:

- is the repo internally consistent
- can learners actually run or verify what the docs promise
- do the docs and metadata agree
- does each content type meet its contract
- are migration waves safe enough to release incrementally

## Quality Layers

v2 should treat QA as several distinct layers.

### 1. Structural Integrity

Checks that the repo layout and referenced paths exist.

Examples:

- section paths exist
- lesson and exercise paths exist
- `_starter/` paths exist when declared
- commands point to valid targets

### 2. Metadata Integrity

Checks that the curriculum metadata is internally coherent.

Examples:

- ids are unique
- section references resolve
- prerequisite ids resolve
- next-item links resolve
- allowed `type`, `subtype`, `level`, and `verification_mode` values are used

### 3. Content Contract Integrity

Checks that the declared content type matches the repo reality closely enough.

Examples:

- a `lesson` exposes a run or test path
- an `exercise` with a declared starter actually has `_starter/`
- a `checkpoint` declares explicit verification
- a `mini_project` has a required README

### 4. Navigation Integrity

Checks that learners can actually navigate the system.

Examples:

- section README exists
- required milestone docs exist
- learner-facing docs link to real sections or items
- navigation references do not point to dead paths

### 5. Runtime And Regression Integrity

Checks that runnable content and repo code stay healthy.

Examples:

- `go build ./...`
- `go vet ./...`
- `gofmt -l .`
- `go test ./...`
- `go test -race ./...`

### 6. Human Review Integrity

Checks that should not be delegated entirely to tooling.

Examples:

- whether a lesson is actually understandable
- whether a project scope is appropriate
- whether path-aware docs are honest
- whether a checkpoint meaningfully proves readiness

## Ownership Model

Different quality problems should be owned by different surfaces.

| Layer | Primary Owner | Typical Enforcement |
| ----- | ------------- | ------------------- |
| structural integrity | validator | blocking script/CI |
| metadata integrity | validator | blocking script/CI |
| content contract integrity | validator + reviewer | mixed |
| navigation integrity | validator + reviewer | mixed |
| runtime and regression integrity | CI | blocking checks |
| teaching quality and scope | reviewer/maintainer | manual review |

This keeps the system honest about what can be automated and what cannot.

## Validator Responsibilities

The v2 validator should gradually evolve from the current curriculum checker into a broader system
validator.

It should eventually verify at least:

- section ids and content ids are unique
- declared paths exist
- run and test commands point to real targets
- starter paths exist when declared
- prerequisite links resolve
- next-item links resolve
- type values are allowed
- lesson subtype values are allowed
- verification modes are allowed
- section README paths exist when required
- learner-facing navigation links do not point to missing locations

The validator should not attempt to judge:

- whether teaching quality is good
- whether a lesson explanation is clear
- whether project scope is pedagogically ideal

Those are review concerns, not parser concerns.

## CI Responsibilities

CI should enforce repo-health checks that are cheap, trustworthy, and blocking.

The current baseline is close to the right model.

The default blocking CI checks should remain:

- build
- vet
- format check
- module tidy check
- test
- race detector
- curriculum validator

Coverage and benchmarks can stay informative unless the team later decides to set hard thresholds.

## Type-Aware Validation Rules

The QA model should validate items differently based on content type.

### Lesson

Minimum expectations:

- valid metadata
- valid run or test command
- next step is declared
- production relevance exists in the authored surface

### Drill

Minimum expectations:

- valid metadata
- narrow verification rule
- explicit parent lesson or local reinforcement context

### Exercise

Minimum expectations:

- valid metadata
- explicit requirements
- verification instructions
- `_starter/` when the design says scaffolding is required

### Checkpoint

Minimum expectations:

- explicit pass criteria
- clear section linkage
- verification mode that matches the item shape

### Mini-Project

Minimum expectations:

- README present
- explicit deliverable
- verification or review rubric
- section or phase milestone meaning

### Capstone

Minimum expectations:

- README present
- wider curriculum linkage
- stronger quality expectations than a mini-project

### Reference

Minimum expectations:

- narrow purpose
- correct placement
- live links or commands when they are part of a learner workflow

## Verification Mode Rules

The schema's `verification_mode` should drive what QA expects from an item.

### `run`

The item should declare a runnable command and the target should exist.

### `test`

The item should declare a test command or naturally resolvable test target.

### `rubric`

The item should declare visible success criteria because there may not be a strict automated test.

### `mixed`

The item should combine runnable/testable behavior with human-readable completion rules.

This prevents v2 from pretending all teaching artifacts can be validated the same way.

## Docs And Navigation QA

The navigation system needs both automatic and manual checks.

### Automatic Checks

Eventually validate:

- required section README files exist
- root-level learner links resolve
- curriculum map links resolve
- learning path references resolve
- milestone docs referenced by metadata exist

### Manual Checks

Reviewers should still verify:

- the doc surface is doing the right job
- the root README is not bloated
- the section README actually helps entry and progression
- path-aware guidance is honest and not generic filler

## Authoring And Review Gates

Every content change should pass through three gates.

### Gate A: Authoring Gate

The author confirms:

- metadata reflects reality
- commands run
- docs references are current
- the content type contract is respected

### Gate B: Automated Gate

The repo confirms:

- validator passes
- CI checks pass

### Gate C: Review Gate

A maintainer or reviewer confirms:

- teaching scope is right
- navigation is clear
- the change fits the current migration wave

## Prototype QA Requirements

The prototype phase should prove that the QA model is practical.

At minimum, the prototype should demonstrate:

- one lesson that passes a lesson-type contract
- one exercise with scaffolding and verification
- one checkpoint with explicit pass criteria
- one mini-project with milestone documentation
- metadata that links all of the above coherently

If the validator cannot meaningfully check the prototype slice, the QA model is still incomplete.

## Migration Wave QA

Each migration wave should have a clear readiness bar.

Minimum wave readiness:

- migrated items use the approved content model
- required docs surfaces exist
- metadata is coherent
- validator passes
- CI passes
- reviewer agrees the wave is understandable to learners

This keeps alpha releases honest even before the full repo is migrated.

## Release-Stage QA

The QA bar should rise as the release line matures.

### Alpha

Focus on:

- structural correctness
- honest path and docs guidance
- prototype rules working at section scale

### Beta

Focus on:

- consistency across sections
- navigation completeness
- schema and validator coverage
- migration guide quality

### RC And Final

Focus on:

- no known blocking drift between docs, metadata, and repo paths
- stable release messaging
- clear learner migration guidance

## Standards Docs Relationship

The existing standards docs remain useful, but v2 should position them more clearly.

- `CODE-STANDARDS.md` should define style and authoring expectations
- `TESTING-STANDARDS.md` should define testing patterns and when tests add value
- this QA document should define the system-level enforcement model

That distinction matters because not every useful practice is a validator rule.

## Definition Of Done

The v2 QA and validation model is ready when:

- validator scope is clearly separated from human review scope
- CI responsibilities are explicit
- content types have distinct minimum validation expectations
- docs and navigation quality have an explicit home
- prototype and migration wave gates are defined
- maintainers can explain why a change failed QA without guessing

## Open Decisions

These questions still need review:

- how much of navigation integrity should be validator-enforced versus review-enforced
- whether rubric-style verification should require a standard README section template
- whether coverage should remain informational or become a release gate for specific sections
- whether the current single validator script should grow gradually or split into focused tools

## Working Recommendation

For the first v2 implementation:

- keep the validator focused on structural, metadata, command, and link truth
- keep CI blocking on build, vet, format, tidy, tests, race, and validator success
- treat coverage and benchmarks as informative by default
- require type-aware review of checkpoints and projects
- expand validator scope gradually as the prototype proves which checks are worth automating
