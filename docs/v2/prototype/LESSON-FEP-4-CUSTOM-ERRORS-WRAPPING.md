# Prototype Canonical Lesson: FEP.4 Custom Errors, Wrapping, And Inspection

## Purpose

This document defines the canonical lesson for the first v2 structural prototype.

It is the answer to issue `#86` and serves as the reference lesson for:

- lesson anatomy
- lesson metadata
- production framing
- code-first versus README support decisions
- exit-ramp behavior inside a section

## Prototype Decision

The canonical lesson for the first prototype slice is:

- `FEP.4 Custom Errors, Wrapping, and Inspection`

Lesson shape:

- `type`: `lesson`
- `subtype`: `pattern`

## Why `FEP.4` Was Chosen

`FEP.4` is the strongest canonical lesson candidate because it proves more of the v2 lesson
contract than `FEP.3` while still staying contained.

It proves:

- a clear lesson objective
- a reusable engineering pattern
- visible production relevance
- a compact but meaningful example chain
- an obvious next step into guided synthesis

It is also a better canonical lesson than `FEP.3` because it shows how a v2 lesson can move beyond
"here is a syntax feature" and into "here is an engineering pattern you will use repeatedly."

## Why `FEP.3` Was Not Chosen As The Canonical Lesson

`FEP.3 Errors as Values and Multiple Returns` still matters and stays in the section.
It is the section hinge.

But as the canonical prototype lesson it is less complete because:

- it is more concept-heavy than pattern-heavy
- it proves fewer decisions about production framing
- it is more foundational and less representative of how richer lessons should feel later

`FEP.3` should still exist in the section as the prerequisite lesson that makes `FEP.4` possible.

## Lesson Mission

Teach learners how to model errors as meaningful values that carry identity, context, and
inspectable structure across function boundaries.

The learner should leave this lesson understanding that Go error handling is not only:

- returning `err`

but also:

- creating stable error identities
- wrapping failures with context
- extracting structured failure information safely

## Placement In The Section

This lesson belongs after:

- `FEP.3 Errors as Values and Multiple Returns`

And before:

- `FEP.E1 Safe Operations Runner`

Its job is to convert the learner from "I know `(value, error)` exists" to "I can shape errors as a
real engineering boundary."

## Canonical Metadata

The lesson should use this metadata model in the prototype:

```json
{
  "id": "FEP.4",
  "section_id": "s04-prototype",
  "slug": "custom-errors-wrapping-inspection",
  "title": "Custom Errors, Wrapping, and Inspection",
  "type": "lesson",
  "subtype": "pattern",
  "level": "core",
  "verification_mode": "run",
  "estimated_time": 35,
  "summary": "Teach learners to create stable error identities, wrap failures with context, and inspect error chains safely.",
  "objectives": [
    "Define sentinel and typed errors for meaningful failure states",
    "Wrap lower-level errors with contextual information",
    "Use errors.Is and errors.As to inspect an error chain safely"
  ],
  "prerequisites": ["FEP.3"],
  "production_relevance": "Go services rely on wrapped and inspectable errors to preserve control-flow clarity while still exposing enough context for logs, retries, and user-safe handling.",
  "path": "04-functions-and-errors/5-custom-errors-wrapping",
  "run_command": "go run ./04-functions-and-errors/5-custom-errors-wrapping",
  "test_command": "",
  "starter_path": "",
  "next_items": ["FEP.E1"],
  "tags": ["errors", "wrapping", "inspection", "go-idioms"]
}
```

## Primary Objective

The primary objective is:

- teach the error-pattern stack of sentinel errors, typed errors, wrapping, and inspection through
  one coherent example

This is one main objective, not four separate lesson objectives, because all four pieces belong to
the same engineering pattern.

## Supporting Objectives

Supporting objectives:

- show why string matching on errors is fragile
- show how context is added without losing original error identity
- show when a typed error is more useful than a plain sentinel

## Verification Mode

The prototype lesson should use:

- `verification_mode: run`

That is the right choice because:

- the lesson is teaching an inspectable runtime pattern
- the example can demonstrate error chains clearly through output
- tests are helpful later, but they are not necessary to prove the lesson contract itself

## README Decision

This canonical lesson should stay:

- code-first
- no separate README in the first prototype

Why:

- the runtime setup is simple
- the example can stay in one file
- the pattern is best learned by reading a compact error chain and running it
- adding a README here would prove less than proving that a code-first lesson can still be
  complete, readable, and production-minded

This does **not** mean pattern lessons never need READMEs.
It means this particular prototype lesson proves that the default v2 lesson can still live in code
when the scope is well bounded.

## Canonical Anatomy

This lesson should prove the five required anatomy layers.

### 1. Framing

The opening should answer:

- what the learner is about to learn
- why plain `if err != nil` is not the whole story
- why this pattern matters before the exercise and checkpoint

### 2. Core Example

The lesson should use one compact layered example, such as:

- a lower-level operation returning a sentinel or typed error
- a middle layer wrapping with `%w`
- a top layer inspecting with `errors.Is` and `errors.As`

The example should be realistic enough to feel like application code, but small enough to read in
one sitting.

### 3. Explanation

The explanation should cover:

- when to use a sentinel error
- when to use a typed error
- what wrapping preserves
- why `errors.Is` and `errors.As` are safer than string matching

### 4. Production Relevance

The production framing should connect the lesson to:

- service-layer context propagation
- retry or classification decisions
- cleaner logs and safer handling boundaries

It should stay concrete and not drift into generic architecture talk.

### 5. Exit Ramp

The exit ramp should explicitly point to:

- `FEP.E1 Safe Operations Runner`

The learner should understand that the exercise is where this lesson stops being demonstration and
becomes implementation.

## Canonical Example Shape

The lesson should use one bounded error chain with three layers:

1. origin layer
   - returns a sentinel or typed error
2. service layer
   - wraps with contextual information
3. caller layer
   - inspects with `errors.Is` or `errors.As`

This is enough structure to prove the pattern without growing into a multi-package application.

## What The Lesson Should Not Try To Do

To stay canonical and maintainable, this lesson should not:

- introduce `defer`, `panic`, or `recover`
- become the guided exercise itself
- teach multiple unrelated error domains
- add concurrency, IO, or HTTP infrastructure

That scope discipline is part of what makes it a good canonical lesson.

## Lesson Size Decision

This lesson should be medium-sized.

Too small:

- it would fail to prove the pattern lesson shape

Too large:

- it would collapse into an exercise

The right prototype size is:

- one strong example
- one clear engineering pattern
- enough explanation to teach tradeoffs
- no extra architecture noise

## Code Comments Versus README Support

For this canonical lesson, use:

- inline code comments for local mechanics and reasoning
- header framing for lesson intent and engineering depth
- no extra README for the prototype

This establishes an important v2 default:

- if the code can teach clearly on its own, keep the lesson lightweight

## Lesson-Spec Validation Result

This prototype lesson confirms that these lesson-spec expectations are realistic:

- one primary objective is enough
- `type`, `subtype`, `level`, and `verification_mode` are useful and not redundant
- `next_items` is essential
- production relevance should be explicit even in earlier sections
- README support should be optional, not automatic

## Small Lesson-Spec Adjustment Discovered

The prototype suggests one practical clarification for v2:

- pattern lessons should usually center on one bounded engineering pattern with a single layered
  example, not a catalog of related techniques

That is already compatible with the lesson spec, but the prototype makes the boundary clearer.

## Success Signal For Issue #86

Issue `#86` should be considered complete when maintainers can point to this lesson and say:

- this is the right size
- this is the right amount of framing
- this is enough production depth without over-teaching
- this is how a code-first pattern lesson should feel in v2

## Recommended Follow-On Work

Now that the canonical lesson is chosen:

1. `#87` should define `FEP.E1 Safe Operations Runner` and `FEP.C1 Reliable Batch Processor Checkpoint`
   around the error-pattern flow established here
2. `#88` should define the `FEP.P1` mini-project linkage and the metadata example that includes
   this lesson as a prerequisite and predecessor
