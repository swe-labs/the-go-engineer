# V2 Content Type System

## Purpose

v2 needs a clear content type system so the curriculum can be designed, navigated, validated, and
migrated consistently.

This document defines the canonical content roles used across lessons, exercises, checkpoints,
projects, metadata, and learning paths.

## Core Rule

Every content item in v2 should have one primary job.

If an item is trying to introduce a concept, provide guided practice, validate mastery, and act as a
project at the same time, it is usually the wrong shape.

## The Four Classification Dimensions

v2 content should be classified across four separate dimensions:

### 1. Top-Level Type

The item's curriculum role:

- `lesson`
- `drill`
- `exercise`
- `checkpoint`
- `mini_project`
- `capstone`
- `reference`

### 2. Subtype Or Shape

The internal form of the item.
This is most useful for lessons.

For the first v2 draft:

- lesson subtypes: `concept`, `pattern`, `integration`

### 3. Difficulty Band

The pacing and ambiguity level:

- `foundation`
- `core`
- `stretch`
- `production`

### 4. Verification Mode

How the learner or repo verifies the item:

- runnable example
- testable behavior
- prompt/rubric verification
- mixed verification

These dimensions should not be collapsed into one overloaded field.

## Canonical Top-Level Types

| Type | Primary Job | Typical Scale | Default Verification | Typical Position |
| ---- | ----------- | ------------- | -------------------- | ---------------- |
| `lesson` | teach a core idea clearly | small to medium | `go run`, inspection, or focused test | throughout a section |
| `drill` | reinforce the lesson immediately | small | quick run, prompt, or small assertion | directly after a lesson |
| `exercise` | guided synthesis with scaffolding | medium | run/test + explicit success criteria | after several lessons |
| `checkpoint` | validate readiness to move forward | medium | test, rubric, or explicit completion rules | late in a section or track |
| `mini_project` | build a meaningful workflow | medium to large | runnable artifact + project criteria | section or phase milestone |
| `capstone` | synthesize a larger curriculum span | large | multi-surface verification | phase or program milestone |
| `reference` | support tooling, setup, or lookup | small to medium | command or setup verification when possible | entry/setup or targeted support |

## Type Contracts

### Lesson

A `lesson` is the primary teaching unit.

Use it when:

- introducing a new concept
- teaching a reusable pattern
- combining a small set of prior ideas in a controlled example

Do not use it when:

- the main goal is to test readiness
- the learner should already know most of the content and now needs guided synthesis
- the scope has grown into a real workflow or project

Required characteristics:

- one primary learning objective
- a clear core example
- explanation and tradeoffs
- explicit production relevance
- a next step

Allowed lesson subtypes:

- `concept`
- `pattern`
- `integration`

### Drill

A `drill` is immediate reinforcement.

Use it when:

- learners need quick repetition right after a lesson
- the practice target is narrow and focused
- the goal is confidence, not broad synthesis

Do not use it when:

- the item needs a large prompt or scaffold
- several lessons must be combined
- the item is actually testing readiness to move on

Required characteristics:

- short completion time
- one narrow skill target
- explicit link back to the lesson it reinforces

### Exercise

An `exercise` is guided synthesis.

Use it when:

- two to four prior lessons need to be combined
- the learner benefits from a prompt, constraints, and starter scaffolding
- the goal is applied practice with support

Do not use it when:

- the item should act as the formal readiness gate for a section
- the work has become a meaningful application workflow better treated as a mini-project

Required characteristics:

- clear requirements
- success criteria
- starter scaffolding when appropriate
- explicit prerequisites
- verification instructions

### Checkpoint

A `checkpoint` validates progression.

Use it when:

- the learner needs a clear readiness signal
- the section has introduced enough material that a gate is helpful
- the repo needs an explicit boundary before larger work

Do not use it when:

- the learner mainly needs supported guided practice
- the activity has become a large build artifact

Required characteristics:

- clear pass criteria
- explicit scope
- visible relation to the section outcomes
- minimal ambiguity about what "ready" means

### Mini-Project

A `mini_project` is a meaningful workflow milestone.

Use it when:

- the learner must integrate several skills into one artifact
- package boundaries, interfaces, IO, or workflow design matter
- the project represents a real curriculum jump

Do not use it when:

- the work is too small and would be better as an exercise
- the scope spans too much of the curriculum and has become a capstone

Required characteristics:

- explicit scope
- concrete deliverable
- verification rules
- section or phase relevance

### Capstone

A `capstone` is wide synthesis.

Use it when:

- the learner is consolidating a whole phase or the final arc
- quality, architecture, or deployment concerns become part of the teaching goal

Do not use it when:

- the work is only a larger mini-project with no broader curriculum role

Required characteristics:

- broad but bounded scope
- visible connection to a larger curriculum span
- stronger quality expectations than smaller items

### Reference

A `reference` item supports the path without carrying the same teaching role as a lesson.

Use it when:

- setup, tooling, or command guidance is needed
- the item helps learners enter or troubleshoot the main path
- the information is valuable but should not masquerade as a core concept lesson

Do not use it when:

- the item teaches a core conceptual idea that belongs in a lesson
- the item is only a dumping ground for details that should be documented more clearly elsewhere

Required characteristics:

- narrow purpose
- clear relation to the main path
- verification where possible

## Boundary Rules

These are the most important distinctions to keep stable.

### Lesson vs Reference

- `lesson`: core teaching object
- `reference`: support object

If skipping the item would leave the learner without a needed concept, it is probably a lesson.

### Drill vs Exercise

- `drill`: immediate reinforcement of one narrow skill
- `exercise`: guided synthesis across multiple lessons

### Exercise vs Checkpoint

- `exercise`: supportive practice
- `checkpoint`: readiness validation

An exercise can help a learner get ready.
A checkpoint should tell whether the learner is ready.

### Mini-Project vs Capstone

- `mini_project`: local milestone
- `capstone`: broader phase or program synthesis

### Integration Lesson vs Exercise

An `integration` lesson still teaches through a guided example.
An `exercise` asks the learner to perform the synthesis with their own implementation decisions.

## Type-Specific Expectations

| Type | README | `_starter/` | Tests | Next Step | Notes |
| ---- | ------ | ----------- | ----- | --------- | ----- |
| `lesson` | optional | no | optional | required | code-first by default |
| `drill` | optional | optional | optional | required | keep lightweight |
| `exercise` | preferred | usually yes | optional but useful | required | must have success criteria |
| `checkpoint` | preferred | optional | often useful | required | readiness signal matters |
| `mini_project` | required | optional | optional but encouraged | required | scope must be explicit |
| `capstone` | required | optional | optional but encouraged | required | stronger quality bar |
| `reference` | preferred | no | optional | optional | only when it is part of a route |

## Schema Implications

The schema should treat:

- `type` as the top-level role
- `subtype` as the optional internal shape
- `level` as the difficulty band
- `verification_mode` as the primary verification strategy

First draft rule:

- every content item must declare a `type`
- lessons may also declare a `subtype`
- other item types may add subtype conventions later if needed
- every content item should eventually declare a `verification_mode`

## Folder And Navigation Implications

This content type system implies:

- drills are first-class content objects even when they stay lightweight
- checkpoints should be visible in section flow, not buried as unnamed exercises
- reference content usually belongs in docs or support surfaces unless it is part of the canonical path

## Validator Implications

The validator should eventually enforce:

- all items use an allowed top-level `type`
- lesson `subtype` values are valid
- required paths and commands exist for the declared item shape
- items that declare `_starter/` actually have it
- checkpoints and projects declare enough verification information to be meaningful

## First Prototype Guidance

For the first prototype:

- use one `lesson`
- use one `exercise`
- use one `checkpoint`
- use one `mini_project`
- optionally include a `drill` if it helps prove the flow

That is enough to validate the type system without overbuilding the prototype.
