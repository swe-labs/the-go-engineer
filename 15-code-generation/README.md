# Section 15: Code Generation

## Mission

This section teaches you how to use generation tools to remove repetitive work without hiding the
real code behind runtime magic.

The live v2 slice focuses on three habits:

- treating `go generate` as a build-time workflow instead of a runtime trick
- generating mocks from interfaces when hand-written test doubles stop scaling
- generating typed data-access code from SQL instead of relying on reflection-heavy ORMs

## Who Should Start Here

### Full Path

Start here after finishing Section 14 in order.

### Bridge Path

You can move faster if you already understand:

- interfaces and test doubles from the testing sections
- package boundaries and code ownership
- basic SQL and repository-style data access

Even on the bridge path, do not skip `CG.1`.
It explains the build-time mental model the rest of the section depends on.

### Targeted Path

This section is a single focused track.
Follow the lessons in order:

- `CG.1` for the `go generate` workflow
- `CG.2` for generated mocks in test-heavy codebases
- `CG.3` for schema-driven SQL code generation

## Section Map

| ID | Type | Surface | Why It Matters |
| --- | --- | --- | --- |
| `CG.1` | Lesson | [go generate primer](./1-go-generate) | Explains how generation tools fit into normal Go builds. |
| `CG.2` | Lesson | [mockery workflow](./2-mockery) | Replaces repetitive manual mocks with generated test doubles. |
| `CG.3` | Lesson | [sqlc workflow](./3-sqlc) | Generates typed query code from real SQL schemas and queries. |

## Suggested Order

1. Start with `CG.1` so the generation model is clear before any tool-specific workflow.
2. Continue to `CG.2` to see code generation improve testing ergonomics.
3. Finish with `CG.3` to see code generation applied to production data access.

## Section Milestone

`CG.3` is the current Section 15 output.

If you can explain:

- why code generation is a compile-time productivity tool instead of runtime magic
- why generated mocks and generated query layers are safer than repetitive handwritten glue
- why generated code should still be readable, reviewable, and committed intentionally

then the code-generation part of the curriculum is doing its job.

## Next Step

After `CG.3`, you have reached the end of the current public curriculum line.
The best next move is to revisit any weak milestone, rebuild a capstone from memory, or apply the
patterns in your own production-style Go project.
