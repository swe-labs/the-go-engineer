# V2 Lesson Template

## Purpose

This document defines the canonical reusable lesson template for v2.

It turns the approved lesson planning and prototype work into a maintainable authoring surface that
contributors can copy without improvising:

- lesson metadata
- code header structure
- lesson anatomy
- README-first explanation rules
- validator expectations

This template is derived primarily from:

- `docs/v2/04-LESSON-SPEC.md`
- `docs/v2/11-CURRICULUM-SCHEMA.md`
- `docs/v2/prototype/LESSON-FEP-4-CUSTOM-ERRORS-WRAPPING.md`
- `CODE-STANDARDS.md`

## When To Use This Template

Use this template when the content item is:

- `type: lesson`

Allowed lesson subtypes in the current v2 draft:

- `concept`
- `pattern`
- `integration`

Do not use this template for:

- drills
- exercises
- checkpoints
- mini-projects
- reference items

## Canonical Lesson Directory Shape

Default lesson layout:

```text
N-lesson-name/
  main.go
  README.md
  main_test.go     # optional
  testdata/        # optional
```

The canonical beta lesson template assumes a learner-facing README exists.
The explanation depth may change by stage, but the README-first teaching contract stays the same.

Use a larger multi-file layout only when the teaching goal genuinely requires it.
Lesson complexity should not grow just because the file layout allows it.

## Metadata Stub

Every v2 lesson should start with a metadata draft like this:

```json
{
  "id": "SX.Y",
  "section_id": "sNN",
  "slug": "lesson-slug",
  "title": "Lesson Title",
  "type": "lesson",
  "subtype": "concept",
  "level": "core",
  "verification_mode": "run",
  "estimated_time": 30,
  "summary": "One-sentence description of the lesson's teaching goal.",
  "objectives": [
    "Primary teaching objective",
    "Optional supporting objective"
  ],
  "prerequisites": ["SX.Z"],
  "production_relevance": "One concrete sentence about where this lesson matters in real Go work.",
  "path": "NN-section-name/N-lesson-name",
  "run_command": "go run ./NN-section-name/N-lesson-name",
  "test_command": "",
  "starter_path": "",
  "next_items": ["SX.Y+1"],
  "tags": ["topic-a", "topic-b"]
}
```

## Metadata Field Notes

- `type` must always be `lesson`
- `subtype` should be one of `concept`, `pattern`, or `integration`
- `level` should usually be `foundation`, `core`, `stretch`, or `production`
- `verification_mode` should be `run` unless a test-driven lesson shape is clearly better
- `objectives` should stay small and focused
- `next_items` should point to the most meaningful immediate next step, not every possible path

## Lesson Header Template

Every lesson entry file should begin with a header shaped like this:

```go
// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.
//
// Section NN: Section Name - Lesson Title
//
// Mental model:
// One short sentence naming the learner's framing idea.
//
// Run: go run ./NN-section-name/N-lesson-name
```

Keep the header short.
Do not move the README's deep explanation into the file header.

## Canonical Lesson Skeleton

Use this as the default beta lesson code shape:

```go
// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package main

import "fmt"

// Section NN: Section Name - Lesson Title
//
// Mental model:
// One short framing line.
//
// RUN: go run ./NN-section-name/N-lesson-name

func main() {
	// Use local comments only where they genuinely help the learner understand
	// a non-obvious mechanic, boundary, or reasoning step.

	fmt.Println("replace with lesson output")
}
```

## Required Lesson Anatomy Checklist

Every v2 lesson should clearly cover these five layers:

1. Framing
2. Core example
3. Explanation
4. Production relevance
5. Exit ramp

If one of those layers is missing, the lesson is probably not ready.

## Subtype Guidance

### Concept Lesson

Use when:

- introducing a new language or design idea
- teaching one clear mental model

Default shape:

- one runnable example
- tight explanation
- low architecture overhead

### Pattern Lesson

Use when:

- teaching a reusable engineering habit
- comparing good and bad approaches
- showing one bounded layered example

Default shape:

- one strong pattern example
- explicit tradeoffs
- clear production boundary

### Integration Lesson

Use when:

- combining several prior ideas in a controlled example
- preparing the learner for a checkpoint or mini-project

Default shape:

- one synthesis example
- explicit relationship to prior lessons
- clear warning against overexpansion

## Canonical Lesson README

The beta lesson template requires a learner-facing README.

Use the README to carry the explanation load that would otherwise make the code noisy.

Every lesson README should usually include:

1. mission
2. prerequisites
3. mental model
4. what problem this lesson solves now
5. run instructions
6. code walkthrough
7. common questions or failure points
8. production relevance
9. next step

The walkthrough should explain the code line by line or in small logical chunks.
Group lines only when they belong to one inseparable step.

See `LESSON-README-TEMPLATE.md` for the default shape.

## README Decision Rules

Use a full README for learner-facing beta lessons.

README-first does not remove the code requirement.
The learner should read the explanation first, then run and inspect the code itself.

## Test Decision Rules

Add `main_test.go` when:

- the behavior being taught is naturally testable
- the lesson also teaches validation or correctness checks
- a test clarifies the lesson better than more prose would

Do not add tests to every lesson by default.

## Scope Control Rules

Split the lesson when:

- more than one major new concept is being introduced
- the example is large enough that the core idea is hard to see
- the explanation depends on several unrelated digressions

Keep the lesson compact by preferring:

- one strong example
- one primary objective
- one explicit next step

## Copy-Ready Authoring Checklist

Before calling a lesson draft complete, confirm:

- metadata matches the actual lesson
- the lesson has one clear primary objective
- the run or test command is real
- the production relevance is explicit
- the next step is named clearly in `next_items`
- the README explains the code at the right depth for the intended learner
- the source file stays readable instead of absorbing the whole explanation burden
- the code remains a required runnable proof surface after the README

## Validator Notes

The validator should eventually enforce these lesson-specific checks:

- `type` is `lesson`
- `subtype` uses an allowed lesson subtype
- `path` exists
- either `run_command` or `test_command` points to a real target
- `next_items` resolves
- declared README support is reflected by real files once the schema carries that signal

The validator should not try to judge:

- whether the example is pedagogically strong
- whether the explanation is vivid enough
- whether the lesson should have been split earlier

Those remain reviewer judgments.

## Relationship To Existing Public Standards

This template does not replace the public standards yet.

For now:

- `CODE-STANDARDS.md` remains the public baseline in `main`
- this template is the v2 planning/prototype-derived maintainer template on `planning/v2`

When learner-facing v2 implementation begins, the repo can decide how to merge the two cleanly.

## Success Signal

This template is working when a maintainer can copy it and produce a new v2 lesson without
guessing:

- the metadata shape
- the file header
- the lesson anatomy
- how the README and code file divide the teaching work
- what validator expectations the lesson should satisfy

