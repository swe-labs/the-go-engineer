# 05 Functions and Errors

## Mission

This section teaches learners how to turn inline logic into reusable functions and how to treat
failure as an explicit part of a function contract.

By the end of this section, a learner should be able to:

- write small functions with clear parameters and return values
- combine data-structure and control-flow knowledge inside reusable logic
- return `(value, error)` when a result can fail
- validate inputs before doing work
- sequence a few small functions into one honest milestone flow

## Why This Section Exists Now

The learner already knows:

- values and expressions
- branching and loops
- slices, maps, and pointer-aware mutation

That is enough to ask the next engineering questions:

- when should logic stop living directly inside `main()`?
- how should one piece of code give a result back to another?
- how should a program describe failure without hiding it?

Those are function-and-error questions.

## Zero-Magic Boundary

This section intentionally stays inside foundations-ready ideas:

- plain functions
- parameters and return values
- multiple return values
- errors as values
- validation and orchestration

It does **not** make later concepts part of the foundations-critical path:

- custom error types that depend on richer type modeling
- panic / recover as ordinary program flow
- functional options
- package architecture and cross-file system design

Those may remain in the repo as legacy or later-stage material, but they are not the primary
foundations route.

## Beta Stage Ownership

This section belongs to [03 Functions & Errors](../../docs/stages/03-functions-errors.md).

In the current repo architecture, it lives at `03-functions-errors`.

## Section Map

The canonical foundations order for this section is:

| ID | Type | Surface | Core Job |
| --- | --- | --- | --- |
| `FE.1` | Lesson | [functions basics](./1-functions-basics) | teach what a function boundary is |
| `FE.2` | Lesson | [parameters and returns](./2-parameters-and-returns) | teach input/output contracts |
| `FE.3` | Lesson | [multiple return values](./3-multiple-return-values) | introduce Go's multi-result style |
| `FE.4` | Lesson | [errors as values](./4-errors-as-values) | teach explicit failure handling |
| `FE.5` | Lesson | [validation](./5-validation) | teach early input checks and explicit rejection |
| `FE.6` | Lesson | [orchestration](./6-orchestration) | teach how one function coordinates smaller helpers |
| `FE.7` | Exercise | [order summary](./7-order-summary) | prove reusable logic plus explicit failure handling in one small milestone |

## Current Rebuild Goal

This section is being rebuilt so that:

- learner-facing explanation stays in lesson `README.md` files
- `main.go` stays runnable and clean
- the milestone proves earned foundations concepts only
- old Section `04` material can be kept as legacy reference where helpful without confusing the new
  primary path

## Suggested Learning Flow

1. Read each lesson `README.md` first.
2. Open `main.go` only after the lesson mission and machine view are clear.
3. Run the code and compare the output with the walkthrough.
4. Change one thing at a time using the `Try It` prompts.
5. Move to the milestone only after the early lessons stop feeling mechanical.

## Next Step

After `FE.7`, the learner should be ready to move into
[04 Types & Design](../../docs/stages/04-types-design.md).
That is where richer type modeling and design boundaries become the main topic.
