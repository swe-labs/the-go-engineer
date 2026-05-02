# The Go Engineer Curriculum Blueprint

This document defines how the locked v2.1 architecture behaves as a learning system.

If this file and [ARCHITECTURE.md](./ARCHITECTURE.md) disagree on public structure, `ARCHITECTURE.md` wins.

## Core Promise

The curriculum should help a learner move from copying code to understanding, changing, testing, and operating Go software with clear reasoning.

The expected progression is:

- understand what a line of Go does
- explain why that design exists
- predict what can break
- prove behavior with tests or validation
- connect isolated lessons into a production-shaped backend system

## Architecture Lock

The public curriculum spine is locked at 12 sections, `s00` through `s11`.

Do not solve content growth by adding a new public root section. Add depth inside existing sections unless a maintainer explicitly approves architecture work.

## Teaching Rules

### README first, code second

Every learner-facing lesson teaches through `README.md` before source code.

The learner should understand the mission, prerequisites, mental model, and machine view before reading `main.go`.

### Code is never skipped

Do not replace code with prose. Explain the code, run the code, then modify the code.

### Zero magic

Each section teaches only what has been earned.

If a concept depends on later ideas, either move it later or add a local source comment or README alert that names the future lesson or section.

### Explanations answer how, why, and what changes

Good teaching surfaces explain:

- what the line or block does
- why it exists
- what would change if the design changed
- what mistake a learner is likely to make

### Engineering depth is stage-aware

Beginner sections stay concrete. Later sections add design trade-offs, failure modes, security, operations, and reliability reasoning.

## Phase Blueprint

| Phase | Sections | Required emphasis |
| --- | --- | --- |
| 0 | s00 | machine model, terminal confidence, execution basics |
| 1 | s01-s04 | syntax, control flow, data structures, functions, errors, types |
| 2 | s05-s08 | packages, I/O, APIs, databases, concurrency, tests, profiling |
| 3 | s09-s10 | architecture, security, production operations |
| 4 | s11 | integrated system design through Opslane |

## Canonical Lesson Contract

Default lesson shape:

```text
lesson-name/
|-- README.md
|-- main.go
|-- main_test.go
`-- _starter/
    `-- main.go
```

Not every lesson needs tests or starter code. Exercises and behavior-heavy lessons do.

### Required README Sections

Each lesson README must include these sections in this order:

1. `## Mission`
2. `## Prerequisites`
3. `## Mental Model`
4. `## Visual Model`
5. `## Machine View`
6. `## Run Instructions`
7. `## Code Walkthrough`
8. `## Try It`
9. `## In Production`
10. `## Thinking Questions`
11. `## Next Step`

The `## Next Step` section must keep the next curriculum path visible and make it clickable.
Use the next item ID, show the next curriculum path as the link label, and point the Markdown link to that lesson's `README.md`.

For exercises:

- replace `## Code Walkthrough` with `## Solution Walkthrough`
- include `## Verification Surface`

### Required Source Behavior

Source files should stay readable and runnable.

Use inline comments for:

- non-obvious behavior
- mutation or boundary traps
- runtime implications
- local cross-references that explain why another lesson's concept appears here

Do not turn source headers into the main teaching surface. The README is primary.
Follow `CODE-STANDARDS.md` for Machine Role comments, source header fields, `RUN:` command format, `NEXT UP:` footers, and proof-surface consistency.

## Cross-Reference Rules

When a lesson uses a concept not yet formally taught:

- introduce the borrowed idea in one or two sentences
- name the future lesson or section where it is taught in detail
- keep the explanation local to the paragraph where the concept appears

When a lesson reuses a concept taught earlier:

- name the earlier lesson ID or section
- remind the learner why that earlier idea matters here
- do not repeat the full earlier lesson

When neighboring tracks use the same idea:

- point to the sibling lesson only when it improves navigation.
- Use GitHub-style alerts (`[!NOTE]` or `[!TIP]`) to integrate references without interrupting the narrative flow.
- Include the lesson ID and a clickable local `README.md` link when referencing a specific lesson.
- Avoid detached, standalone "Forward/Backward Reference" headlines.

## Milestone Contract

Every section needs proof, not only explanations.

A milestone should usually provide:

- clear README instructions
- a runnable completed solution
- starter code when the learner is expected to implement
- tests when behavior should be provable
- matching curriculum metadata, source `Level`/`RUN:` headers, source `NEXT UP:` footer, README run instructions, and README next-step link

## Revision Checklist

When adding or revising lessons:

1. work inside an existing section
2. keep the public section count at exactly 12
3. update `curriculum.v2.json`
4. update the section README
5. keep README, source, tests, starter code, and metadata aligned
6. run `go run ./scripts/validate_curriculum.go`

## Bottom Line

The Go Engineer should read as one coherent engineering learning system. The 12 sections are the public spine; README-first teaching, runnable code, strict validation, and clear cross-references are the delivery standard.
