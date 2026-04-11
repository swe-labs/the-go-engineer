# 11 Code Generation

## Purpose

`11 Code Generation` teaches leverage after understanding.

## Who This Is For

- learners who already understand the code and systems they are working with
- developers who want to use generation responsibly instead of blindly

## Mental Model

Generation is multiplication, not understanding.
It becomes useful only after you know what the generated code is doing and how it fits into the
system.

## Why This Stage Exists

This stage exists so learners use generation as leverage after understanding, not as a shortcut
around understanding.

The goal is to make generation workflows reviewable, intentional, and connected to real system
needs.

## What You Should Learn Here

- where generation helps
- where generation creates hidden cost
- how generated code fits into normal review and maintenance workflows
- how to keep generated workflows understandable

## Stage Shape

This stage currently has one live public source section with one focused generation path:

1. `15-code-generation`
   - the live beta path for `go generate`, generated mocks, and schema-driven SQL generation

That makes this stage intentionally narrow.
It is about responsible leverage, not about turning the repo into a catalog of tooling for its own
sake.

## Current Source Content

- [15-code-generation](../../15-code-generation/)

## Stage Support Docs

Use these support docs when you want the beta-stage view of code generation:

- [Code Generation support index](./code-generation/README.md)
- [Stage map](./code-generation/stage-map.md)
- [Milestone guidance](./code-generation/milestone-guidance.md)

## Where This Stage Starts

Start with [Section 15: Code Generation](../../15-code-generation/).

That section is already the full public source surface for this stage.

## Recommended Order

Use this order for the current beta-facing path:

1. `CG.1` for the `go generate` mental model
2. `CG.2` for generated mocks and test workflow leverage
3. `CG.3` for schema-driven SQL generation in production-style code

## Path Guidance

### Full Path

Enter this stage after the flagship and expert-pressure work are already meaningful.
That keeps generation connected to real engineering needs instead of abstract tool enthusiasm.

### Bridge Path

If you already use generation tools professionally, you can move faster, but do not skip `CG.1`.
That lesson establishes the "build-time workflow, not runtime magic" contract for the whole stage.

### Targeted Path

If your immediate goal is tooling leverage:

- start with `CG.1` if you need the basic generation model
- start with `CG.2` if your gap is test double generation
- start with `CG.3` if your gap is query/code generation around SQL

## Stage Milestones

The current live milestone backbone is:

- `CG.3` sqlc workflow

## Finish This Stage When

- you can explain why a generation workflow exists
- you can read and review the code that generation produces
- you know when not to use generation
- you understand generation as leverage, not magic

More concretely:

- you can explain why generated code still belongs inside normal review and maintenance flows
- you can connect generation workflows back to testability, contracts, or data access needs
- you can explain when generation reduces risk and when it increases hidden cost

## Next Stage

Use this stage as a late specialization layer after the flagship and expert-pressure work are
already meaningful to you.
