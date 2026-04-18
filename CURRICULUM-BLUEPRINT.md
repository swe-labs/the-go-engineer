# The Go Engineer Curriculum Blueprint

> This document defines how the 12-section v2.1 architecture should behave as a learning system.
> It is the curriculum contract for how we teach. The source of truth for section layout is
> `ARCHITECTURE.md`. If a conflict exists, `ARCHITECTURE.md` wins.

## Core Promise

The Go Engineer should help a learner move from:

- "I can copy code"

to:

- "I understand what this line does"
- "I can explain why this design exists"
- "I can predict what breaks"
- "I can build and operate a real system"

---

## Non-Negotiable Teaching Rules

### 1. README first, code second

Every learner-facing lesson teaches through `README.md` first.
The learner should understand the mission before opening `main.go`.

### 2. Code is never skipped

We do not replace code with prose.
We explain the code, then run the code, then modify the code.

### 3. Zero magic

Each section teaches only what has been earned.
If a concept depends on later ideas, it belongs later or must include a clearly labelled forward reference: _(more on this in Lesson X.N)_

### 4. Explanation must answer how, why, and what changes

Good teaching surfaces explain:

- what this line or block does
- why it exists
- what would change if we changed it
- what mistake a learner is likely to make here

### 5. Engineering depth must be stage-aware

We do want:

- design thinking
- failure thinking
- production relevance
- debugging instincts

But we add them at the right layer.
We do not dump senior-level pressure framing into beginner lessons just to sound impressive.

---

## Phase-Level Blueprint

The curriculum is split into 5 phases across 12 sections (s00–s11). See `ARCHITECTURE.md` for the exact breakdown.

### Phase 0: Machine Foundation (s00)

This phase explains WHY code works at all before writing any Go. Safe, explicit, visual.

Required elements:

- mission and mental model
- visual diagrams (Mermaid, not ASCII art)
- plain-language analogies
- runnable demonstrations

### Phase 1: Language Foundation (s01–s04)

These sections must feel safe, explicit, and zero-magic. The learner is building Go fluency.

Required elements:

- mission
- mental model
- literal or near-literal walkthroughs
- clean runnable code

Avoid:

- premature scale pressure
- advanced security catalogues
- abstract design jargon before the learner has concrete examples

### Phase 2: Engineering Core (s05–s08)

These sections start increasing engineering judgement. The learner is building systems.

Add more of:

- trade-off explanations
- failure cases
- safer defaults
- tests and verification surfaces
- performance and maintainability reasoning
- ⚠️ In Production notes with real-world consequences

### Phase 3: Systems Engineering (s09–s10)

These sections carry full engineering weight. Architecture, security, and deployment.

Add more of:

- architecture trade-offs
- production notes
- security implications
- deployment patterns

### Phase 4: Flagship Project (s11 — GoScale)

This phase carries the heaviest engineering pressure:

- integrated project proof
- production deployment
- operational pressure
- all prior concepts applied together

---

## Canonical Lesson Contract

For learner-facing lessons, the default shape is:

```text
lesson-name/
├── README.md          ← Written first. Learner reads this before opening main.go.
├── main.go            ← Primary lesson code
├── main_test.go       ← Tests (required for exercises)
└── _starter/
    └── main.go        ← TODO stubs (exercises only)
```

### Required README sections

Each lesson README must include these sections **in this order** (see `ARCHITECTURE.md` § Lesson README Contract):

1. `## Mission`
2. `## Prerequisites` — list lesson IDs (e.g., "- CF.1, CF.2")
3. `## Mental Model`
4. `## Visual Model` — Mermaid diagram (NOT ASCII art)
5. `## Machine View`
6. `## Run Instructions` — exact `go run` or `go test` command
7. `## Code Walkthrough`
8. `## Try It` — numbered steps the learner actually does
9. `## ⚠️ In Production` — required in EVERY lesson
10. `## 🤔 Thinking Questions` — minimum 3 per lesson
11. `## Next Step`

For exercises: replace `## Code Walkthrough` with `## Solution Walkthrough`, add `## Verification Surface`.

### Required source-file behaviour

Source files should stay readable and runnable.
They should not become giant essays.

Use inline comments for:

- non-obvious behaviour
- mutation or boundary traps
- subtle runtime implications

Do not use code headers as the main teaching surface.

Every `main()` must end with a KEY TAKEAWAY comment and a NEXT UP footer (see `CODE-STANDARDS.md`).

---

## Canonical Milestone Contract

Every section needs proof, not just lessons.

A milestone should usually provide:

- clear README instructions
- a runnable completed solution
- a starter surface when the learner is expected to implement
- tests when the behaviour should be provable

---

## Cross-Reference Rules

When a lesson uses a concept not yet formally taught:

**Forward reference:** _(We use `defer` here. Defer is taught in Lesson CF.5. For now, read it as "run this line when the function returns.")_

**Backward reference:** _(This builds on error wrapping from FE.4.)_

**Sibling reference:** _(Context cancellation is covered in CT.1–CT.5. Here we receive and respect it; there we learn to create it.)_

---

## How To Add New Lessons Without Breaking The Architecture

If the curriculum needs more depth:

1. Add the lesson inside an existing section.
2. Keep the learner-facing section count at exactly 12 (s00–s11).
3. Update `ARCHITECTURE.md` if the scope of the section expands.
4. Register the lesson in `curriculum.v2.json`.
5. Ensure the validator passes: `go run ./scripts/validate_curriculum.go`

Do not solve content growth by inventing a new root-level section unless the entire architecture is being reworked intentionally.

---

## Bottom Line

The Go Engineer should feel like one coherent engineering learning system.

The 12 sections give us the public spine.
The README-first teaching contract gives us the delivery standard.
Future expansion should deepen the sections we have, not fragment the learner path again.
