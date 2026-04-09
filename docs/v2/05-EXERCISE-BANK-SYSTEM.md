# V2 Exercise Bank System

## Purpose

Exercises need to become a first-class system in v2.
Right now, exercises are one of the strongest parts of the repo, but they are not yet governed by a
single architecture.

This document defines the first draft of that architecture.

Within the broader v2 content type system, this document focuses on the practice-oriented items:
drills, exercises, checkpoints, mini-projects, and capstones.

## Exercise Types

v2 should use five exercise classes:

### 1. Lesson Drill

- short and focused
- reinforces the lesson that just finished
- usually solved in one file

### 2. Guided Exercise

- combines two to four lessons
- includes clear requirements and constraints
- usually ships with `_starter/`

### 3. Section Checkpoint

- validates the main ideas of the section
- should feel like a confidence test, not a trick
- often includes a test file or explicit verification steps

### 4. Mini-Project

- spans a meaningful workflow
- uses multiple packages or interfaces when appropriate
- produces an artifact the learner can run, inspect, or extend

### 5. Capstone

- synthesizes a whole phase of the curriculum
- should model real engineering concerns such as validation, IO boundaries, testing, or deployment

## Canonical Repo Layout

The default exercise layout should be:

```text
NN-section-name/
  N-exercise-name/
    README.md
    main.go
    _starter/
      main.go
    internal/            # optional for larger exercises
    testdata/            # optional when input fixtures matter
    main_test.go         # optional but preferred when exercise behavior is testable
```

For mini-projects and capstones, package-sized layouts are allowed when they improve realism.

## Required Exercise Components

Every non-trivial v2 exercise should include:

- a short learner-facing README
- explicit requirements
- success criteria
- starter code when self-implementation is expected
- verification instructions
- clear linkage to prerequisite lessons

## Metadata Proposal

Each exercise should eventually carry metadata for:

- `id`
- `title`
- `section`
- `level`
- `exercise_type`
- `prerequisites`
- `estimated_time`
- `starter_path`
- `solution_path`
- `test_command`
- `skills_validated`

This can live in an extended curriculum schema later.

## Difficulty Bands

Use the same bands across the repo:

- `foundation`: direct reinforcement
- `core`: standard synthesis
- `stretch`: more design choice and ambiguity
- `production`: realistic boundaries and tradeoffs

Beta note:

- `docs/v2/19-BETA-EXERCISE-RUBRIC-SYSTEM.md` is the stricter beta-phase authority for starter
  mode, verification mode, and rubric expectations, but it keeps these same difficulty bands

## Definition Of Done

An exercise is ready for v2 only when:

- the prompt is clear
- the solution matches the prompt
- the starter scaffolding is usable
- the verification instructions actually work
- the required concepts match the declared prerequisites

## Migration Rules From v1

When migrating v1 exercises into v2:

- keep the strongest existing exercise ideas
- normalize naming and scaffolding
- add README and verification instructions where missing
- add tests where they improve feedback
- split oversized exercises into checkpoint plus mini-project when needed

## First Draft Recommendation

The first prototype should include:

- one lesson drill
- one guided exercise with `_starter/`
- one section checkpoint
- one mini-project

That set is enough to validate whether the exercise bank system feels coherent before large-scale
migration begins.
