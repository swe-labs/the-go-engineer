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

This section belongs to [1 Language Fundamentals](../../docs/stages/01-language-fundamentals.md).

Inside the new repo architecture, it is the fifth foundations section:

1. `01-getting-started`
2. `02-language-basics`
3. `03-control-flow`
4. `04-data-structures`
5. `05-functions-and-errors`

## Planned Section Shape

The canonical foundations order for this section is:

| ID | Planned Surface | Core Job |
| --- | --- | --- |
| `FE.1` | functions basics | teach what a function boundary is |
| `FE.2` | parameters and returns | teach input/output contracts |
| `FE.3` | multiple return values | introduce Go's multi-result style |
| `FE.4` | errors as values | teach explicit failure handling |
| `FE.5` | validation | teach early input checks and honest failure |
| `FE.6` | orchestration | combine small functions into one readable flow |
| `FE.7` | milestone project | prove reusable logic plus explicit failure handling |

## Current Migration Goal

This section is being rebuilt so that:

- learner-facing explanation stays in lesson `README.md` files
- `main.go` stays runnable and clean
- the milestone proves earned foundations concepts only
- old Section `04` material can be kept as legacy reference where helpful without confusing the new
  primary path

## Next Step

The next implementation move is to rebuild the first lesson group under this section and then align
the milestone around the new foundations scope.
