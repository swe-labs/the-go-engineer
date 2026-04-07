# V2 Mini-Project Template

## Purpose

This document defines the canonical reusable mini-project template for v2.

It turns the approved mini-project planning and prototype work into a maintainable authoring
surface that contributors can copy without improvising:

- mini-project metadata
- milestone-sized project scope
- folder layout
- README structure
- verification expectations
- validator implications

This template is derived primarily from:

- `docs/v2/06-PROJECT-LADDER.md`
- `docs/v2/09-FOLDER-STRUCTURE.md`
- `docs/v2/11-CURRICULUM-SCHEMA.md`
- `docs/v2/prototype/MINI-PROJECT-FEP-P1-TASK-JOURNAL-CLI.md`

## When To Use This Template

Use this template when the content item is:

- `type: mini_project`

Use a mini-project when:

- the learner should produce a named runnable artifact
- the section or phase needs a meaningful milestone beyond a checkpoint
- package boundaries, entrypoints, or workflow shape now matter
- the work validates a real curriculum jump

Do not use this template for:

- guided exercises
- checkpoints
- phase capstones
- final capstones

## Canonical Mini-Project Directory Shape

Default mini-project layout:

```text
N-project-name/
  README.md
  cmd/
    project-name/
      main.go
  internal/
    domain/
    workflow/
  testdata/        # optional
  main_test.go     # optional
```

Optional additions:

- `pkg/` only when a real public package boundary exists
- `go.mod` only when project isolation genuinely improves teaching or tooling

Mini-projects should normally stay inside the section or phase they validate.
Do not move them into a top-level `projects/` tree during early migration waves.

## Metadata Stub

Every v2 mini-project should start with a metadata draft like this:

```json
{
  "id": "SX.P1",
  "section_id": "sNN",
  "slug": "project-slug",
  "title": "Project Title",
  "type": "mini_project",
  "level": "core",
  "verification_mode": "mixed",
  "estimated_time": 150,
  "summary": "One-sentence description of the milestone artifact.",
  "objectives": [
    "Primary milestone objective",
    "Optional supporting objective"
  ],
  "prerequisites": ["SX.C1"],
  "production_relevance": "One concrete sentence about why this project shape matters in real Go work.",
  "path": "NN-section-name/N-project-name",
  "run_command": "go run ./NN-section-name/N-project-name/cmd/project-name",
  "test_command": "",
  "starter_path": "",
  "next_items": ["sNN+1"],
  "tags": ["mini-project", "topic-a", "topic-b"]
}
```

## Metadata Field Notes

- `type` must always be `mini_project`
- `verification_mode` should usually be `mixed`
- `starter_path` should normally be empty
- `next_items` may point to a later local item or the next section id when the project is a
  section-exit milestone

## Canonical README Shape

Every v2 mini-project README should include these sections:

1. project mission
2. prerequisites
3. required features
4. project layout guidance
5. run instructions
6. success criteria
7. common failure modes
8. extension ideas
9. next step

This keeps the project grounded as a milestone artifact rather than a vague "build something"
prompt.

## Canonical README Skeleton

Use this as the default mini-project README shape:

~~~md
# Project Title

## Mission

What artifact the learner is building and why it marks a curriculum milestone.

## Prerequisites

- lesson, checkpoint, or prior project ids the learner should already understand

## Required Features

1. concrete feature requirement
2. concrete feature requirement
3. one scope boundary that keeps the project honest

## Project Layout Guidance

- expected `cmd/` entrypoint
- expected `internal/` ownership
- any intentionally deferred structure

## Run Instructions

~~~bash
go run ./NN-section-name/N-project-name/cmd/project-name
~~~

Optional:

~~~bash
go test ./NN-section-name/N-project-name/...
~~~

## Success Criteria

- the artifact runs from the declared entrypoint
- valid workflows behave correctly
- failure behavior is explicit and inspectable
- the package split is readable and not overbuilt

## Common Failure Modes

- one likely architecture mistake
- one likely error-handling mistake
- one likely scope-creep mistake

## Extension Ideas

- one safe extension for a later section
- one note about what later phases would add

## Next Step

Point to the next section, phase artifact, or migration wave surface.
~~~

## `_starter/` Rules

Mini-projects should normally not include:

- `_starter/`

That is the default because a mini-project is supposed to feel like a milestone artifact, not a
half-built guided exercise.

Only consider starter scaffolding when:

- the project introduces a teaching concern unrelated to the milestone itself
- setup burden would dominate the project unfairly
- the README alone cannot keep the scope honest

That should be uncommon.

## Scope Rules

A mini-project should:

- validate one meaningful curriculum jump
- stay in one bounded domain
- use enough structure to feel like real software
- avoid unrelated feature creep
- stop well short of capstone breadth

A mini-project should not:

- require infrastructure the learner has not studied yet
- simulate enterprise complexity for style points
- become a phase-spanning system
- ask for too many unrelated features at once

## Verification Rules

Mini-project verification should:

- confirm the artifact runs from the expected entrypoint
- confirm the declared workflow works end to end
- confirm expected failures are handled intentionally
- confirm the README's success criteria match the actual scope

When useful, verification may combine:

- runnable output
- focused tests
- rubric-style project checks

That is why `verification_mode: mixed` is the usual default.

## Failure Pressure Rules

Mini-projects should include meaningful non-happy-path behavior.

That can be:

- invalid commands or actions
- duplicate or conflicting domain operations
- not-found or state-transition errors
- summaries that must remain trustworthy after some failures

The goal is not to turn the project into a trap.
The goal is to prove the learner can build a small tool that behaves honestly under imperfect input.

## Package Layout Rules

Use package structure intentionally:

- `cmd/` owns the entrypoint
- `internal/` owns domain and workflow logic
- `pkg/` is optional and should be rare in early mini-projects

Prefer restrained package boundaries over "clean architecture theater."

## Validator Notes

The validator should eventually enforce these mini-project-specific checks:

- `type` is `mini_project`
- `path` exists
- `run_command` points to a real target
- `next_items` resolves
- a learner-facing README exists
- the project layout includes the declared entrypoint

The validator should not try to judge:

- whether the project is interesting enough
- whether the package design is perfectly calibrated
- whether the milestone should have been a project instead of a checkpoint

Those remain reviewer judgments.

## Success Signal

This template is working when a maintainer can copy it and produce a new v2 mini-project without
guessing:

- how big the project should be
- what layout is expected
- what belongs in the README
- how to keep the project milestone-sized instead of drifting into capstone scope
