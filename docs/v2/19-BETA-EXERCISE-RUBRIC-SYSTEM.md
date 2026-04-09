# V2 Beta Exercise And Rubric System

Status: Draft 0.1
Audience: Maintainers and curriculum designers
Related issues: `#171`, `#174`

## Purpose

This document freezes the beta exercise system as a product requirement.

Alpha proved that migrated sections can carry milestone exercises and starter folders.
Beta must go further: exercises, starter-solution flow, and evaluation criteria must become a
coherent training system rather than a section-by-section convention.

This document is the planning authority for `#174`.

## Core Decision

Beta exercise design uses separate axes instead of one overloaded label:

- `exercise_type`
- `difficulty_band`
- `starter_mode`
- `verification_mode`

This separation matters because a guided exercise can be easy or hard, a mini-project can be core
or expert, and a checkpoint can be rubric-based even when it is not test-heavy.

## Canonical Exercise Types

### 1. Drill

Purpose:

- immediate reinforcement after a lesson

Expected size:

- one concept
- one file or one tightly bounded package

Starter expectation:

- optional

### 2. Guided Exercise

Purpose:

- combine a small cluster of lessons into one supported implementation task

Expected size:

- two to four concepts
- learner still benefits from scaffolding

Starter expectation:

- required

### 3. Checkpoint

Purpose:

- prove readiness at the end of a stage slice or meaningful unit

Expected size:

- independent synthesis
- no hidden concept jumps

Starter expectation:

- not allowed by default
- only permit limited skeleton setup when infrastructure burden would otherwise hide the learning goal

### 4. Mini-Project

Purpose:

- turn a stage into a runnable engineering artifact

Expected size:

- multiple files or packages when useful
- explicit scope boundary so it stays smaller than a flagship project

Starter expectation:

- optional partial scaffold
- never a nearly-complete implementation

### 5. Flagship Milestone

Purpose:

- prove integrated engineering growth inside the flagship project

Expected size:

- multi-stage
- architecture, testing, operations, and review pressure all matter

Starter expectation:

- stage scaffolds allowed
- full starter not allowed

## Canonical Difficulty Bands

Beta uses four public difficulty bands:

- `easy`
- `medium`
- `hard`
- `expert`

Use them like this:

- `easy`: direct reinforcement with low ambiguity
- `medium`: standard synthesis with one or two design choices
- `hard`: realistic ambiguity, trade-offs, and stronger boundary handling
- `expert`: review, failure analysis, or design pressure that assumes broad stage mastery

These difficulty bands do not replace exercise type.
They describe challenge, not format.

## Starter Modes

Every beta exercise item must declare one starter mode:

- `none`
- `full_starter`
- `partial_scaffold`

Rules:

- drills may use `none` or `full_starter`
- guided exercises must use `full_starter`
- checkpoints default to `none`
- mini-projects may use `partial_scaffold`
- flagship milestones may use `partial_scaffold`, but not `full_starter`

## Verification Modes

Every beta exercise item must declare one verification mode:

- `run`
- `test`
- `mixed`
- `rubric`

Meaning:

- `run`: learner verifies through explicit runtime behavior
- `test`: learner verifies mainly through tests
- `mixed`: both runnable behavior and tests matter
- `rubric`: human-visible criteria are the main proof surface

Practical rule:

- beta should prefer `test` or `mixed` whenever honest automation improves feedback
- beta should prefer `rubric` for architecture, production, expert, and flagship work where quality
  cannot be reduced to a green test suite

## Starter-Solution Contract

Beta uses one consistent starter-solution pattern.

### Repository rule

The canonical surface is:

```text
item/
  README.md
  main.go            # reference solution or canonical implementation
  _starter/          # learner scaffold when allowed
  main_test.go       # when automated verification is part of the contract
  testdata/          # optional
```

### Required rules

- if `_starter/` exists, it must prepare the same task that the main solution demonstrates
- starter comments may guide the learner, but must not hide the real design decisions inside
  copy-paste blocks
- the solution must prove the README claims directly
- the solution must not introduce hidden requirements absent from the README
- the starter must not describe a smaller or different exercise than the solution

### Anti-patterns

Do not allow:

- starter files that solve half the real exercise already
- README claims that the solution never demonstrates
- solutions that quietly depend on later-stage concepts
- checkpoints disguised as guided exercises

## Canonical Rubric Surface

Every checkpoint, mini-project, flagship milestone, and rubric-verified item must ship with a
visible rubric section in its README.

The canonical rubric dimensions are:

1. Correctness
2. Completeness
3. Boundary handling
4. Code quality
5. Verification discipline
6. Explanation or reasoning quality when the item asks for trade-offs or review thinking

### Rubric rule

Not every item needs deep scoring.
But every rubric-verified item must answer:

- what must work
- what must be handled carefully
- what quality expectations matter
- what proof counts as "done"

## Minimum Beta Exercise Coverage

Beta does not need perfect coverage everywhere, but it does need an honest minimum floor.

### Every beta stage must include

- at least one guided exercise or checkpoint
- at least one clearly verified milestone item
- explicit linkage from lessons into practice

### Early stages must include

- more drills and guided exercises
- lower ambiguity
- stronger starter support

### Middle stages must include

- checkpoints
- mini-projects
- more `mixed` verification

### Later stages must include

- rubric-heavy evaluation
- trade-off discussion
- failure and production thinking

## Beta Must Ship Versus Beta May Defer

### Beta must ship

- the canonical exercise type system
- the four difficulty bands
- the starter-mode contract
- the verification-mode contract
- visible rubric sections for rubric-verified items
- starter plus solution consistency on all newly redesigned beta practice items

### Beta may defer

- full easy-medium-hard-expert coverage for every stage
- auto-scoring infrastructure
- hidden test orchestration
- progress tracking beyond repo-native proof

## Metadata Implications

Beta practice metadata should eventually support these fields:

- `exercise_type`
- `difficulty_band`
- `starter_mode`
- `verification_mode`
- `prerequisites`
- `skills_validated`
- `estimated_time`
- `requires_rubric`

When the schema is updated, these fields should be treated as first-class, not ad hoc tags.

## Definition Of Done For The Exercise-System Phase

`#174` is ready to close when:

- beta uses one canonical exercise-type model
- starter-solution expectations are explicit
- rubric expectations are explicit
- beta implementation issues can describe practice work without inventing local rules each time
