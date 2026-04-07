# V2 Schema Transition Plan

## Purpose

This document turns the v2 schema direction into an implementation-ready transition plan.

It is the working answer for issue `#96`.
Its job is to tell maintainers:

- what the smallest safe next metadata step is
- what must remain in `curriculum.json` during transition
- what should move into a v2 metadata surface first
- which schema concerns are intentionally deferred

## Current Reality

The repo already has one live metadata system:

- `curriculum.json`

That file is still useful.
It powers the current lesson graph well enough for v1 and still supports existing validation and
navigation work.

The problem is not that `curriculum.json` is bad.
The problem is that it cannot describe the richer v2 system cleanly enough on its own.

It does not model first-class:

- exercises
- checkpoints
- mini-projects
- verification mode
- next-item navigation
- learning-path route logic

## Transition Goal

The first schema transition should make v2 metadata possible without forcing a full-repo cutover.

That means the first step should be:

- additive
- section-scoped
- validator-friendly
- reversible if the first migration wave exposes a better shape

## Recommendation

The smallest safe next step is:

1. keep `curriculum.json` in place as the live legacy lesson map
2. add a new additive file named `curriculum.v2.json`
3. populate `curriculum.v2.json` only for migrated v2 sections and items
4. keep learning-path metadata out of the first implementation step

This is the conservative path.
It gives the validator and migration waves enough structure without turning the first schema PR
into a full metadata migration.

## What Stays In `curriculum.json`

During the first v2 migration waves, `curriculum.json` should keep owning:

- legacy section ordering
- legacy lesson presence
- legacy lesson prerequisites
- legacy lesson paths needed by current repo tooling

Compatibility rule:

- if a v2 lesson is live on `main`, its legacy lesson record should remain representable in
  `curriculum.json` until the repo officially stops depending on that file

This keeps the current validator and learner-facing compatibility surfaces from collapsing
mid-migration.

## What Moves First To `curriculum.v2.json`

The first additive v2 file should own only the data the old file cannot model well.

Start with:

- section records for migrated v2 sections
- content-item records for migrated lessons
- content-item records for exercises
- content-item records for checkpoints
- content-item records for mini-projects

The first implementation should not require full metadata coverage for untouched legacy sections.

## First File Shape

The first `curriculum.v2.json` should stay intentionally small.

Recommended top-level shape:

```json
{
  "schema_version": 1,
  "sections": [],
  "items": []
}
```

Why this shape is the right first step:

- it covers the validator needs from issue `#95`
- it matches the Section 04 prototype metadata surface
- it avoids opening the learning-path storage question too early
- it stays easy to diff and reason about during the first migration wave

## First Implementation Scope

The first schema implementation should do only the following:

1. introduce `curriculum.v2.json`
2. add records for the first migrated v2 section only
3. extend `scripts/validate_curriculum.go` to read the v2 file when present
4. keep `curriculum.json` validation intact
5. document which file owns which truth during transition

The first implementation should not:

- backfill the whole repo into `curriculum.v2.json`
- generate `curriculum.v2.json` from `curriculum.json`
- generate `curriculum.json` from v2 data yet
- split metadata across many new files
- require learning paths to move out of docs immediately

## Compatibility Rules

These rules should govern the transition window.

### 1. `curriculum.json` Remains Required

Do not remove or hollow out `curriculum.json` during the first migration waves.

### 2. Dual Authorship Is Allowed Only For Lessons

During transition, a migrated lesson may appear in both files:

- `curriculum.json` for legacy compatibility
- `curriculum.v2.json` for rich v2 behavior

Exercises, checkpoints, and mini-projects should live only in `curriculum.v2.json`.

### 3. Shared Lesson Truth Must Agree

When a lesson exists in both files, these fields must stay aligned:

- stable id or approved mapped id
- path
- prerequisite truth at the lesson graph level

`curriculum.v2.json` owns the richer fields such as:

- `type`
- `subtype`
- `verification_mode`
- `next_items`
- `production_relevance`

### 4. Preserve Stable Lesson IDs When Reasonable

If a v1 lesson maps cleanly to one v2 lesson, preserve the lesson id when practical.

Assign new ids when:

- a lesson is split
- multiple lessons are merged
- a new exercise, checkpoint, or project item is introduced

This keeps migration mapping sane without forcing fake one-to-one identity.

### 5. Do Not Mirror Non-Lesson Items Into The Legacy File

The old file should not grow ad hoc support for:

- exercises
- checkpoints
- mini-projects
- capstones

That would blur the boundary instead of clarifying it.

### 6. No Full Learning-Path Metadata Yet

Keep learning-path logic in docs during the first schema step.

The schema can add explicit path records later, once the first migrated sections prove:

- what path granularity is actually useful
- whether one file or two files reads better in practice

## Recommended Section States

Maintainers should think about sections in three transition states:

- `legacy`: exists only in `curriculum.json`
- `mixed`: lesson compatibility remains in `curriculum.json`, rich v2 items exist in
  `curriculum.v2.json`
- `v2`: later-stage state where the repo no longer needs legacy compatibility for that section

These states do not need to become schema fields in the first implementation.
They are transition-planning language.

## Validator Alignment

This plan is intentionally aligned with issue `#95`.

The first v2 validator pass needs a metadata surface that can express:

- sections
- typed content items
- prerequisites
- next-item links
- path and command truth
- starter paths

`curriculum.v2.json` with `sections` and `items` is enough for that.

Nothing in the first validator expansion requires learning-path metadata yet.

## Deferred Schema Work

The following work should stay deferred beyond the first schema transition:

- deciding whether learning paths live in the same file or a second file
- deciding whether phases need their own top-level entity in the live file
- deciding whether `curriculum.v2.json` later generates legacy views
- removing `curriculum.json`
- introducing previous-item links in addition to `next_items`
- full metadata coverage for all sections before the first migration wave proves the model

## Recommended Next Step After This Plan

After issue `#96`, the next schema-related implementation should be:

1. create the first real `curriculum.v2.json`
2. populate it with one migrated section
3. extend the validator to consume it
4. keep the legacy file working in parallel

That is the smallest real proof that the schema transition works.

## Success Signal

Issue `#96` should be considered successful when maintainers can say:

- we know what the first v2 metadata file is
- we know what remains in `curriculum.json`
- we know which compatibility rules must hold during migration
- we can start implementation without reopening the whole schema debate
