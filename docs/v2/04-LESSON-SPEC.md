# V2 Lesson Spec

## Purpose

This document defines the canonical contract for a v2 lesson.
It should align curriculum design, repo structure, contributor behavior, and future validator rules.

## Definition

A lesson is the smallest canonical teaching unit in the curriculum.
It is larger than a note and smaller than a project.

A good lesson:

- teaches one primary concept
- uses a small number of supporting concepts
- is runnable or directly verifiable
- points to prerequisites and the next step
- includes production relevance

## Lesson Subtypes

Inside the broader `lesson` content type, v2 should support three lesson subtypes:

### 1. Concept Lesson

- introduces a core language or design idea
- usually centers one runnable example

### 2. Pattern Lesson

- teaches a reusable engineering pattern
- usually compares good and bad approaches

### 3. Integration Lesson

- combines several prior concepts in a controlled example
- often sits immediately before an exercise or checkpoint

`reference` is a separate top-level content type in the content type system, not a standard lesson
subtype.

## Canonical Lesson Contract

Every v2 lesson should define:

- `id`
- `title`
- `slug`
- `section`
- `type`
- `level`
- `subtype`
- `estimated_time`
- `prerequisites`
- `primary_objective`
- `supporting_objectives`
- `production_relevance`
- `verification_mode`
- `run_command` or `test_command`
- `next_items`

For lesson records, `type` should always be `lesson`.
`subtype` then carries the internal lesson shape such as `concept`, `pattern`, or `integration`.

## Repo Layout

The default lesson layout should be:

```text
NN-section-name/
  [optional-subgraph/]N-lesson-name/
    main.go
    README.md        # required for learner-facing lessons
    main_test.go     # optional when testable behavior adds value
    testdata/        # optional when fixtures matter
```

For larger architectural lessons, a small package layout is allowed if the teaching goal requires
it. That exception should be used intentionally, not casually.

## Required Teaching Anatomy

Every lesson should include the following layers.

### 1. Framing

State:

- what the learner is about to learn
- why it matters
- where it sits in the wider curriculum

### 2. Core Example

Show a runnable or inspectable implementation that expresses the main idea clearly.

### 3. Explanation

Explain behavior, tradeoffs, gotchas, and the reason the example is shaped this way.

### 4. Production Relevance

Answer:

- where this shows up in real Go work
- what common failure mode it prevents
- what habit the learner should keep

### 5. Exit Ramp

Tell the learner what to do next:

- next lesson
- drill
- exercise
- section checkpoint

## Teaching Surface Rule

For beta, the teaching surface should change with learner needs.

Use this rule:

- `main.go` is the runnable and inspectable code surface
- `README.md` is the primary explanation surface for learner-facing lessons
- inline comments should explain local mechanics, boundaries, and non-obvious reasoning
- deep walkthroughs, line-by-line explanation, runtime-model framing, and "why this line exists"
  guidance should live in the README rather than bloating the source header

This keeps the code readable while still giving beginners enough explanation to answer:

- what is happening
- why it is happening
- how the code changes behavior
- what would break if a line or step changed

## Code Header Contract

Every lesson entry file should stay lean and include:

- copyright line
- license pointer
- section and lesson title
- one short mental-model or focus note
- run command

Do not use the code header as the main home for:

- long objective lists
- level labels
- engineering-depth prose
- repeated section-level background

Those explanations belong in the learner-facing README and metadata.

## When A Lesson Needs A README

A lesson must include `README.md` when:

- the lesson is part of the learner-facing curriculum
- the learner benefits from line-by-line or small-chunk explanation to follow the example
- the lesson introduces a new mental model, runtime behavior, boundary, or design move that needs
  more than short inline comments

A lesson should also include `README.md` when:

- runtime setup is non-trivial
- the lesson has multiple files or commands
- diagrams, tables, or structured comparisons help understanding
- the learner needs troubleshooting help beyond inline comments

For lessons with a README, the README should usually own:

- mission and prerequisites
- mental model
- the walkthrough of the code line by line or in small logical chunks
- common questions and failure points
- production relevance and next step

The code file should remain mandatory and runnable after the docs.
README-first does not mean prose-only. It means the explanation comes first and the code remains the
required proof surface.

## When A Lesson Needs Tests

A lesson should include tests when:

- the behavior being taught is naturally testable
- the lesson teaches testing or validation itself
- the example would benefit from showing how correctness is asserted

Do not add tests only to satisfy a checkbox if they do not improve the lesson.

## Size Guidelines

A lesson likely needs to be split when:

- it introduces more than one major new concept
- the explanation becomes dependent on several unrelated digressions
- the runnable example is large enough that the core idea is hard to see

A lesson likely needs to be merged when:

- it is too small to justify navigation weight
- it cannot stand on its own objective
- it exists only because of naming or folder convenience

## Authoring Rules

Maintain these rules when writing a lesson:

- explain why, not only what
- do not assume hidden prerequisites
- prefer one strong example over many shallow ones
- show mistakes and tradeoffs when they teach something real
- do not overbuild architecture in early sections

## Definition Of Done

A lesson is ready for v2 when:

- the objective is clear
- the run or test command works
- the metadata matches the lesson reality
- the next step is obvious
- the production relevance is explicit
- the explanation depth matches the intended learner, with README support where needed
- the lesson is small enough to be teachable and large enough to matter
