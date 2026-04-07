# Prototype Guided Exercise: FEP.E1 Safe Operations Runner

## Purpose

This document defines the canonical guided exercise for the first v2 structural prototype.

It is the exercise half of issue `#87` and serves as the reference shape for:

- guided synthesis after a pattern lesson
- `_starter/` expectations
- exercise metadata
- supported verification and success criteria

## Prototype Decision

The canonical guided exercise for the first prototype slice is:

- `FEP.E1 Safe Operations Runner`

Exercise shape:

- `type`: `exercise`
- `level`: `core`
- `verification_mode`: `mixed`

## Why `FEP.E1` Is The Right Exercise

`FEP.E1` grows naturally from the section flow that issue `#85` established and the error-pattern
lesson that issue `#86` selected.

It is the right prototype exercise because it proves:

- synthesis across several lessons rather than one lesson only
- guided practice with visible scaffolding
- structured errors that remain inspectable after wrapping
- a batch-style workflow without requiring larger infrastructure

It also stays intentionally smaller than a mini-project.
That keeps the practice layer honest.

## Exercise Mission

Teach learners to combine function design, multi-return error flow, custom errors, wrapping, and
inspection in one small batch runner.

The learner should leave this exercise knowing how to:

- model a small request and result shape
- implement several operations with explicit `(value, error)` contracts
- wrap failures with operation-level context
- continue processing a batch without hiding which items failed and why

## Placement In The Section

This exercise belongs after:

- `FEP.4 Custom Errors, Wrapping, and Inspection`

And before:

- `FEP.5 Defer, Panic, and Failure Boundaries`

Its job is to convert the learner from:

- "I can read and inspect an error chain"

to:

- "I can use that pattern in a small implementation of my own"

## Canonical Metadata

The exercise should use this metadata model in the prototype:

```json
{
  "id": "FEP.E1",
  "section_id": "s04-prototype",
  "slug": "safe-operations-runner",
  "title": "Safe Operations Runner",
  "type": "exercise",
  "level": "core",
  "verification_mode": "mixed",
  "estimated_time": 60,
  "summary": "Guide learners through building a small runner that applies safe operations, wraps context, and reports structured success and failure results across a batch.",
  "objectives": [
    "Implement small operations with explicit value and error contracts",
    "Wrap operation failures with contextual information while preserving inspectable error identity",
    "Process a batch of requests without stopping at the first failure"
  ],
  "prerequisites": ["FEP.1", "FEP.2", "FEP.3", "FEP.4"],
  "production_relevance": "CLI and service code often process small batches, preserve per-item failure context, and continue work while still producing useful summaries.",
  "path": "04-functions-and-errors/8-safe-operations-runner",
  "run_command": "go run ./04-functions-and-errors/8-safe-operations-runner",
  "test_command": "",
  "starter_path": "04-functions-and-errors/8-safe-operations-runner/_starter",
  "next_items": ["FEP.5"],
  "tags": ["exercise", "functions", "errors", "wrapping", "batch-processing"]
}
```

## Canonical Prompt

The learner should build a small program that processes a batch of operation requests.

Suggested request shape:

- `OperationRequest`
  - operation name
  - one or two integer inputs

Suggested result shape:

- `OperationResult`
  - request identity or operation name
  - computed value when successful
  - final error when unsuccessful

The exercise should ask the learner to:

1. implement at least two safe operations with clear return contracts
2. define one inspectable custom error type and at least one stable sentinel error
3. route each request through a `runOperation`-style boundary that wraps lower-level failures with
   operation context
4. process a small slice of requests without aborting the entire batch on first failure
5. print a final summary that distinguishes successful and failed requests
6. inspect errors with `errors.Is` or `errors.As` rather than string matching

## Scope Rules

To keep the exercise guided and teachable, it should stay within these boundaries:

- no file IO
- no network calls
- no concurrency
- no `panic` or `recover`
- one package is enough
- one batch runner is enough

The learner should be practicing explicit error flow, not wrestling with infrastructure.

## `_starter/` Decision

This guided exercise should include:

- `_starter/`

That is the right choice because this is the first true synthesis item in the section.
The learner should spend effort on function contracts and error shaping, not on guessing the entire
program skeleton.

The starter should provide:

- the request and result structs
- a small sample batch in `main`
- TODO markers for the missing operations and batch runner
- output shape expectations at a high level

The starter should not provide:

- completed custom error types
- completed routing logic
- completed wrapping and inspection behavior

## Verification And Success Criteria

This exercise should use:

- `verification_mode: mixed`

The learner should verify completion by both running the program and checking explicit success
criteria in the README.

Minimum success criteria:

- valid requests produce correct values
- divide-by-zero or similar invalid requests return inspectable errors
- unknown operations are classified through a stable error identity
- the batch continues after a failure instead of aborting immediately
- the final output makes success and failure counts clear
- the learner uses `errors.Is` or `errors.As` at the caller boundary

## README Expectations

This guided exercise should have a learner-facing README with:

1. exercise mission
2. requirements
3. starter guidance
4. verification steps
5. success criteria
6. next step

That README should stay practical and specific.
It should not become a second lesson.

## Why This Is An Exercise Instead Of A Checkpoint

`FEP.E1` is still a supported practice surface.

It differs from the checkpoint because it:

- includes starter scaffolding
- narrows the design space deliberately
- tells the learner what shape to build
- focuses on supported synthesis before readiness validation

Its job is to help the learner succeed.
It is not yet the section's proof gate.

## Validator And Schema Implications

This exercise shape implies that the v2 validator should eventually confirm:

- the declared exercise path exists
- the declared `_starter/` path exists
- the run command points to a real target
- the prerequisites resolve
- the next item resolves to `FEP.5`
- a learner-facing README exists for the exercise

The schema does not need a separate exercise subtype for this prototype.
The top-level `type: exercise` plus `verification_mode: mixed` is enough.

## Success Signal For Issue #87

The exercise half of issue `#87` should be considered successful when maintainers can say:

- this is clearly more than a drill
- this is clearly more supported than a checkpoint
- the starter gives enough help without solving the work
- the verification criteria are concrete and honest

