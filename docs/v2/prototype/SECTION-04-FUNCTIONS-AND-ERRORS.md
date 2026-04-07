# Prototype Section Outline: Section 04 Functions And Errors

## Purpose

This document defines the canonical prototype section outline for the v2 structural prototype.

It is the container for the next prototype issues:

- `#86` canonical lesson
- `#87` guided exercise and checkpoint
- `#88` mini-project and metadata example

## Prototype Decision

The prototype section for the first v2 structural slice is:

- `Section 04: Functions and Errors`

This remains a planning/prototype artifact on `planning/v2`.
It is not yet a learner-facing migration on `main`.

## Why Section 04 Was Chosen

Section 04 is the strongest prototype candidate because it sits at the first real curriculum hinge.

It is where the learner moves from:

- basic syntax and data handling

to:

- explicit control of behavior
- explicit control of failure
- the first project-worthy engineering habits

It is the right prototype because it can prove:

- concept lessons
- a pattern lesson
- an integration lesson
- a drill
- a guided exercise
- a checkpoint
- the first mini-project bridge

without requiring:

- databases
- HTTP infrastructure
- concurrency orchestration
- package-heavy service architecture

## Why Other Sections Were Not Chosen First

### Section 01

Too early and too simple for the first full system proof.
It can prove lesson flow, but it is weaker at proving checkpoint, error-handling, and mini-project
behavior.

### Section 05

Strong candidate, but more concept-dense and easier to overload before the v2 section contract is
proven once in a simpler section.

### Section 09 Or Later

Those sections add infrastructure complexity too early.
That would blur whether the prototype is proving the curriculum system or just wrestling with larger
application setup.

## Section Mission

Teach learners to treat functions and failure handling as explicit engineering tools, not as
incidental syntax.

By the end of this section, learners should understand how Go organizes behavior through functions,
how Go models failure as values, and how cleanup and controlled failure shape real programs.

## Section Role In The Curriculum

This section is the final foundations section and the bridge into the first mini-project.

It should feel like:

- the close of the basic language arc
- the start of production-minded control flow
- the readiness gate for the first real milestone artifact

## Learner Outcomes

After this prototype section, the learner should be able to:

- define functions with clear input and output contracts
- use closures and variadic inputs without losing readability
- apply Go's multi-return and error-as-value convention
- model and inspect errors with enough context to be useful
- use `defer` intentionally for cleanup and failure safety
- explain when `panic` is inappropriate and when `recover` is a boundary tool
- complete a small program that combines function design, explicit errors, and cleanup behavior

## Prerequisite Expectations

### Full Path

Expected prior completion:

- Sections 01-03 in order

That means the learner should already be comfortable with:

- variables and constants
- control flow
- slices, maps, and pointers

### Bridge Path

The learner may skim some earlier repetition, but should already be comfortable with:

- basic Go syntax
- control flow
- reading and writing small functions
- understanding value versus reference effects at a basic level

### Targeted Path

A learner entering here directly should review:

- basic control flow
- slice and map usage
- pointer mutation basics
- the idea that Go prefers explicit returns over hidden control flow

For the prototype, this prerequisite review should live in the section README entry guidance rather
than requiring a separate artifact.

## Prototype Scope Decision

The prototype outline should use:

- five lessons
- one drill
- one guided exercise
- one checkpoint
- one mini-project linkage

That is enough to prove the section contract without turning the prototype into a full migration of
every existing v1 item.

## Canonical Prototype Flow

The prototype section should follow this sequence.

| Order | Prototype ID | Type | Working Title | Purpose |
| ----- | ------------ | ---- | ------------- | ------- |
| 1 | `FEP.1` | lesson (`concept`) | Function Shape and Return Contracts | establish function signatures, parameters, return values, and readable contracts |
| 2 | `FEP.2` | lesson (`concept`) | Function Values, Closures, and Variadic Input | show how Go functions stay flexible without becoming magical |
| 3 | `FEP.D1` | drill | Refactor Repetitive Call Sites | reinforce signature reading and variadic usage immediately |
| 4 | `FEP.3` | lesson (`concept`) | Errors as Values and Multiple Returns | establish the error-return convention as the section hinge |
| 5 | `FEP.4` | lesson (`pattern`) | Custom Errors, Wrapping, and Inspection | teach reusable error-shaping patterns and context propagation |
| 6 | `FEP.E1` | exercise | Safe Operations Runner | guided synthesis of functions, multi-return, custom errors, and wrapping |
| 7 | `FEP.5` | lesson (`integration`) | Defer, Panic, and Failure Boundaries | integrate cleanup and controlled failure into the existing error model |
| 8 | `FEP.C1` | checkpoint | Reliable Batch Processor Checkpoint | prove the learner can combine the section's main ideas without hand-holding |
| 9 | `FEP.P1` | mini_project | Task Journal CLI Linkage | bridge the section into the first foundations mini-project |

## Section Flow Rationale

This order is intentional.

### Lessons 1-2

Build confidence with functions as behavior-shaping tools before failure handling enters the
picture.

### Drill

The drill sits early because this is the first point where learners can practice function shape and
variadic thinking without yet carrying the full error model.

### Lessons 3-4

Turn functions into real engineering tools by introducing multi-return and error modeling.

### Guided Exercise

The exercise belongs after the error-model lessons so the learner can synthesize function contracts,
multi-return, and error inspection before the section adds cleanup and controlled failure.

### Lesson 5

`defer`, `panic`, and `recover` should appear after the learner already understands explicit error
returns.
That keeps the section from accidentally teaching panic as a substitute for normal Go error
handling.

### Checkpoint

The checkpoint belongs after the integration lesson because it should test:

- explicit error flow
- meaningful cleanup
- disciplined failure boundaries

### Mini-Project Linkage

The mini-project sits after the checkpoint as the section-exit milestone.
The checkpoint proves readiness; the mini-project proves the learner can use that readiness in a
slightly larger artifact.

## Section README Shape

The prototype section should require a top-level section README with these parts:

1. section mission
2. who should start here
3. prerequisite guidance for Full Path, Bridge Path, and Targeted Path
4. ordered section map by content type
5. checkpoint and mini-project meaning
6. what the learner should do next after the section

Because this section is single-track, the prototype does not require local sub-track README files.

## Expected Item-Level Docs

For this prototype section:

- lessons may rely mainly on code-first teaching surfaces
- the drill may stay lightweight
- the guided exercise must have a learner-facing README
- the checkpoint must have a README with explicit pass criteria
- the mini-project must have a README with requirements, verification guidance, and milestone
  framing

## v1 To Prototype Mapping

The prototype should reuse the strongest v1 material while normalizing the section shape.

| v1 Item | Prototype Direction |
| ------- | ------------------- |
| `FE.1 functions` | feeds `FEP.1` |
| `FE.2 closures & recursion` | feeds `FEP.2` |
| `FE.3 variadic functions` | feeds `FEP.2` and `FEP.D1` |
| `FE.4 multiple return values` | feeds `FEP.3` |
| `FE.5 custom errors` | feeds `FEP.4` |
| `FE.6 error wrapping` | feeds `FEP.4` |
| `FE.7 defer` | feeds `FEP.5` |
| `FE.8 panic & recover` | feeds `FEP.5` |
| `FE.9 error handling project` | source material for `FEP.E1` and `FEP.C1` |
| `FE.10 functional options pattern` | defer beyond the first prototype slice |

## Why `FE.10` Is Deferred

The functional options pattern is valuable, but it is not required to prove the first v2 section
contract.

Deferring it keeps the prototype focused on the section's core hinge:

- function design
- explicit error handling
- cleanup and failure boundaries

It can return later as:

- an added advanced lesson in the final v2 section
- a pattern extension after the prototype proves the core section flow

## Prototype Metadata Expectations

The prototype section should be able to express, at minimum:

- section mission
- prerequisite expectations
- ordered content items
- item types and lesson subtypes
- checkpoint location
- mini-project linkage

This section outline should feed the metadata example in `#88`.

## Success Signal For Issue #85

Issue `#85` should be considered complete when the team agrees that:

- Section 04 is the right prototype section
- the section flow is coherent
- the practice layers are placed intentionally
- the README/navigation shape is clear
- `#86`, `#87`, and `#88` can now proceed without guessing the container shape

## Recommended Follow-On Work

After this outline:

1. `#86` should define `FEP.4` as the canonical prototype lesson candidate
2. `#87` should define `FEP.E1` and `FEP.C1`
3. `#88` should define `FEP.P1` and the metadata example that ties the section together
