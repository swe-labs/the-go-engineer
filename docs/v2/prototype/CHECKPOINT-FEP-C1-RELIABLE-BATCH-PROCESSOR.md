# Prototype Checkpoint: FEP.C1 Reliable Batch Processor Checkpoint

## Purpose

This document defines the canonical section checkpoint for the first v2 structural prototype.

It is the checkpoint half of issue `#87` and serves as the reference shape for:

- section readiness validation
- checkpoint metadata
- explicit pass criteria
- checkpoint verification without mini-project sprawl

## Prototype Decision

The canonical section checkpoint for the first prototype slice is:

- `FEP.C1 Reliable Batch Processor Checkpoint`

Checkpoint shape:

- `type`: `checkpoint`
- `level`: `core`
- `verification_mode`: `mixed`

## Why `FEP.C1` Is The Right Checkpoint

`FEP.C1` is the right checkpoint because it validates the full section arc rather than just the
error-pattern portion.

It proves the learner can combine:

- function contracts
- batch coordination
- explicit error returns
- wrapped errors and inspection
- `defer` for completion or cleanup behavior
- `panic` and `recover` only at a controlled boundary

It is also intentionally smaller than `FEP.P1`.
That keeps the checkpoint from turning into a disguised mini-project.

## Checkpoint Mission

Validate that the learner can build a small, reliable batch processor that handles both expected
failures and unexpected failure boundaries without collapsing the whole run.

The learner should demonstrate that they can:

- separate per-item work from batch coordination
- preserve useful error context
- use `defer` intentionally
- recover from a bounded panic and convert it into controlled failure output
- finish with a trustworthy readiness summary

## Placement In The Section

This checkpoint belongs after:

- `FEP.5 Defer, Panic, and Failure Boundaries`

And before:

- `FEP.P1 Task Journal CLI Linkage`

Its job is to answer:

- "Can this learner now handle controlled function and failure flow well enough to take on the
  first mini-project?"

## Canonical Metadata

The checkpoint should use this metadata model in the prototype:

```json
{
  "id": "FEP.C1",
  "section_id": "s04-prototype",
  "slug": "reliable-batch-processor-checkpoint",
  "title": "Reliable Batch Processor Checkpoint",
  "type": "checkpoint",
  "level": "core",
  "verification_mode": "mixed",
  "estimated_time": 75,
  "summary": "Validate that the learner can build a small batch processor with explicit function contracts, wrapped errors, defer-based cleanup, and panic recovery at a clear boundary.",
  "objectives": [
    "Design a small batch processor that separates per-item work from batch coordination",
    "Use defer for cleanup or completion tracking and convert unexpected panic into controlled error output",
    "Produce a readiness artifact with per-item outcomes and a trustworthy final summary"
  ],
  "prerequisites": ["FEP.1", "FEP.2", "FEP.3", "FEP.4", "FEP.5"],
  "production_relevance": "Reliable worker and batch-style code needs bounded failure behavior, explicit cleanup, and useful summaries so one bad item does not destroy the whole run.",
  "path": "04-functions-and-errors/9-reliable-batch-processor-checkpoint",
  "run_command": "go run ./04-functions-and-errors/9-reliable-batch-processor-checkpoint",
  "test_command": "",
  "starter_path": "",
  "next_items": ["FEP.P1"],
  "tags": ["checkpoint", "functions", "errors", "defer", "panic", "recover", "batch-processing"]
}
```

## Canonical Prompt

The learner should build a small batch processor that executes a sequence of work items and
produces a final report.

The checkpoint should require the learner to:

1. define a work-item shape and a processing boundary
2. execute several items and record per-item success or failure
3. wrap ordinary failures with item-level context
4. use `defer` for completion tracking, cleanup, or result finalization
5. recover from a deliberately bounded panic inside the processing boundary and convert it into an
   ordinary error or failure record
6. print a final summary that distinguishes success, expected failure, and recovered failure

The checkpoint should not dictate every helper function or every file name.
Some design choice should remain with the learner.

## Scope Rules

To stay a checkpoint instead of a mini-project, this item should stay within these boundaries:

- one small domain only
- one package is enough
- no external dependencies
- no file or network infrastructure
- no persistence layer
- no need for multiple binaries

The complexity should come from failure boundaries and control flow discipline, not setup.

## `_starter/` Decision

This checkpoint should not include:

- `_starter/`

That is the right choice because a checkpoint is a readiness signal.
Giving starter scaffolding here would blur whether the learner can shape the implementation on
their own.

If setup overhead ever becomes the main difficulty, that is a sign the checkpoint is scoped
poorly, not a sign it needs more scaffolding.

## Verification And Pass Criteria

This checkpoint should use:

- `verification_mode: mixed`

It should be verified by running the artifact and comparing it against explicit pass criteria in the
README.

Minimum pass criteria:

- the processor completes a full run even when some items fail
- ordinary failures are wrapped with context instead of printed as flat strings only
- at least one recovered panic is converted into controlled output instead of crashing the program
- `defer` is used intentionally and visibly in the solution
- the final summary is accurate and distinguishes the important failure classes
- the code shape is still small enough to read and reason about in one sitting

## README Expectations

This checkpoint should have a learner-facing README with the standard rubric-oriented shape:

1. checkpoint mission
2. required behavior
3. pass criteria
4. verification steps
5. common failure modes
6. next step

The README should make "ready" mean something concrete.

## Why This Is A Checkpoint Instead Of An Exercise

`FEP.C1` differs from the exercise because it:

- removes starter scaffolding
- validates the whole section arc instead of only the error-pattern arc
- includes a readiness boundary around `defer` and panic recovery
- leaves more implementation structure to the learner

Its job is not to coach the learner through every move.
Its job is to confirm readiness for the mini-project.

## Why This Is Not Yet A Mini-Project

`FEP.C1` still stops short of project scope because it:

- stays in one tightly bounded domain
- avoids multi-package application design
- does not ask for persistent state or CLI polish
- focuses on readiness proof, not on building a milestone artifact

That keeps the mini-project meaningful later.

## Validator And Schema Implications

This checkpoint shape implies that the v2 validator should eventually confirm:

- the declared checkpoint path exists
- the run command points to a real target
- the prerequisites resolve
- the next item resolves to `FEP.P1`
- a learner-facing README exists
- rubric-style pass criteria are visible in the README

This also reinforces a useful prototype rule:

- checkpoints do not need `_starter/` by default, but they do need explicit pass criteria

## Success Signal For Issue #87

The checkpoint half of issue `#87` should be considered successful when maintainers can say:

- this clearly validates readiness rather than teaching through scaffolding
- the pass criteria are concrete
- the checkpoint proves the section's full control-flow arc
- the scope is still smaller than a mini-project

