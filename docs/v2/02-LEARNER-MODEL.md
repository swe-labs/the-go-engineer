# V2 Learner Model

## Purpose

The learner model exists to resolve curriculum tradeoffs.
When maintainers disagree about pacing, depth, scaffolding, or project size, the answer should come
from the learner model instead of personal preference.

## Primary Learner Promise

The default promise of The Go Engineer is:

"Take a motivated learner from first-contact Go to production-minded engineering habits through
clear code, repeated practice, and realistic system-building."

This promise is wider than "learn syntax" and narrower than "replace a full degree program."

## Core Learner Segments

v2 should explicitly design for three canonical learner types.

### 1. The Beginner Builder

Profile:

- little or no prior programming experience
- willing to learn slowly
- benefits from explicit explanation and small wins

Needs:

- careful pacing
- lesson READMEs that explain code line by line or in very small chunks, with depth adjusted by stage
- readable inline comments for local mechanics inside the code itself
- repeated reinforcement
- small but visible exercises
- fewer hidden assumptions

Risks:

- pointer and memory confusion
- discouragement from oversized projects
- being overwhelmed by advanced terminology too early

### 2. The Transfer Learner

Profile:

- already programs in Python, JavaScript, Java, C#, Rust, or similar
- does not need generic programming instruction
- does need Go idioms, tooling, and mental model rewiring

Needs:

- clear statements of what is uniquely Go
- bridges from familiar concepts to Go patterns
- permission to skim beginner mechanics without losing the path
- fast access to exercises that surface idiomatic mistakes

Risks:

- treating Go like their previous language
- dismissing error handling and composition as "simple"
- skipping foundational sections and later hitting architectural confusion

### 3. The Working Go Improver

Profile:

- already writes some Go in production or at work
- wants better concurrency, testing, performance, or architecture habits

Needs:

- targeted entry points
- higher signal-to-noise ratio
- exercises that surface tradeoffs instead of only mechanics
- direct production framing

Risks:

- underestimating foundational gaps
- jumping only to advanced topics without reviewing assumptions
- wanting patterns without understanding prerequisites

## Default Authoring Target

When maintainers write a lesson, they should target:

- clarity for the Beginner Builder
- honesty and speed for the Transfer Learner
- depth and relevance for the Working Go Improver

If a lesson can only satisfy one audience, it should still remain understandable to the Beginner
Builder while optionally pointing faster learners to skim paths.

## Learner States

The same learner can move through several states:

- orientation
- concept acquisition
- guided practice
- synthesis
- self-directed extension

v2 should make those states visible instead of treating every directory as the same type of object.

## Pacing Modes

v2 should support three pacing modes without forking the curriculum.

### Full Path

For complete beginners.
Follow all lessons, drills, checkpoints, and projects in order.

### Bridge Path

For experienced programmers new to Go.
Skim selected foundation lessons, then do the canonical checkpoints and exercises that prove Go
specific understanding.

### Targeted Path

For working Go developers.
Jump into a section with a short prerequisite recap and a clear dependency statement.

## Where Learners Need Extra Support

The curriculum should assume extra support is needed at these transitions:

- zero values to control flow
- slices and maps to pointers and mutation
- functions to errors as values
- structs to interfaces and composition
- files and CLI to multi-package applications
- simple goroutines to cancellation and bounded concurrency
- "code that works" to "code that is testable, observable, and maintainable"

## Design Implications

This learner model implies the following product rules:

- early sections must stay concrete and confidence-building
- learner-facing lessons should use README-first explanation while keeping runnable code clean
- advanced sections must clearly state prerequisites
- projects should grow in complexity gradually
- fast-track learners should be able to skip repetition without skipping validation
- docs should explain where to enter and why

## Success Signals

The learner model is working when:

- beginners can tell where to start and what to do next
- transfer learners can move quickly without becoming lost
- working Go developers can enter advanced sections without the repo feeling too basic
- maintainers can explain why a lesson is shaped the way it is
