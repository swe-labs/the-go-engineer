# V2 Curriculum Schema

## Purpose

v2 needs a richer metadata model than the current `curriculum.json` if we want the curriculum,
exercise bank, project ladder, and validation rules to stay aligned.

This document defines the target schema direction.

## Current Reality

The current `curriculum.json` already gives the repo a valuable foundation:

- sections
- lessons
- ids
- paths
- prerequisites
- entry and exercise flags

That is strong enough for v1, but too small for the full v2 training system.

## Recommendation

The first v2 implementation should use a staged approach:

### Stage 1

- keep `curriculum.json` as the active v1-compatible source
- design the richer v2 schema in parallel
- validate the new schema through the prototype first

### Stage 2

- add a dedicated v2 metadata file, likely `curriculum.v2.json`
- keep compatibility tooling while migration is in progress

### Stage 3

- decide whether the v2 schema replaces the old file or becomes the canonical source that generates
  legacy views

## Schema Goals

The v2 schema should be able to answer:

- what exists
- where it lives
- what it depends on
- what type of content it is
- what learner outcome it serves
- how it is verified

## Core Entities

The v2 schema should model at least these entities:

- phase
- section
- content item
- learning path

`content item` is the important general type because lessons, exercises, checkpoints, and projects
should all be first-class items.

## Section Contract

Each section record should include:

- `id`
- `number`
- `slug`
- `title`
- `phase`
- `summary`
- `status`
- `prerequisites`
- `path_prefix`
- `entry_points`
- `outputs`

## Content Item Contract

Each content item should include:

- `id`
- `section_id`
- `slug`
- `title`
- `type`
- `subtype`
- `level`
- `verification_mode`
- `estimated_time`
- `summary`
- `objectives`
- `prerequisites`
- `production_relevance`
- `path`
- `run_command`
- `test_command`
- `starter_path`
- `next_items`
- `tags`

## Content Types

At minimum, `type` should support:

- `lesson`
- `drill`
- `exercise`
- `checkpoint`
- `mini_project`
- `capstone`
- `reference`

`type` is the top-level curriculum role.
`subtype` is optional and is most useful for lessons.

First prototype convention:

- lesson subtypes: `concept`, `pattern`, `integration`
- level values: `foundation`, `core`, `stretch`, `production`
- verification modes: `run`, `test`, `rubric`, `mixed`

## Learning Path Contract

Learning paths should not duplicate content.
They should reference canonical items and provide route logic.

Each learning path should include:

- `id`
- `title`
- `audience`
- `entry_assumptions`
- `recommended_start_items`
- `required_items`
- `optional_items`
- `milestones`

## Example Section Record

```json
{
  "id": "s01",
  "number": "01",
  "slug": "core-foundations",
  "title": "Core Foundations",
  "phase": "foundations",
  "summary": "Introduce the Go toolchain, syntax confidence, and the basic mental model of a Go program.",
  "status": "planned",
  "prerequisites": [],
  "path_prefix": "01-core-foundations",
  "entry_points": ["GS.1", "LB.1"],
  "outputs": ["beginner-checkpoint"]
}
```

## Example Content Item Record

```json
{
  "id": "FE.4",
  "section_id": "s04",
  "slug": "multiple-returns",
  "title": "Multiple Returns",
  "type": "lesson",
  "subtype": "concept",
  "level": "core",
  "verification_mode": "run",
  "estimated_time": 30,
  "summary": "Introduce Go's multi-return pattern and how it leads into error-as-value design.",
  "objectives": [
    "Return multiple values from a function",
    "Recognize why Go pairs data returns with error returns"
  ],
  "prerequisites": ["FE.1", "FE.2", "FE.3"],
  "production_relevance": "Multi-return is the foundation for idiomatic error handling throughout Go codebases.",
  "path": "04-functions-and-errors/4-multiple-returns",
  "run_command": "go run ./04-functions-and-errors/4-multiple-returns",
  "test_command": "",
  "starter_path": "",
  "next_items": ["FE.5"],
  "tags": ["errors", "functions", "go-idioms"]
}
```

## Validator Responsibilities

The v2 validator should eventually verify:

- ids are unique
- referenced sections exist
- prerequisite ids resolve
- paths exist
- declared starter paths exist
- run and test commands point to valid targets
- content types use allowed values
- next-item links resolve

## Migration Constraints

The schema should stay realistic.
Do not add metadata simply because it looks complete.

Every field should meet one of these tests:

- it helps learners navigate
- it helps contributors author correctly
- it helps validators catch drift
- it helps release planning and migration

## Open Schema Decisions

These points still need review:

- whether `estimated_time` should be integer minutes or rough bands
- whether path navigation should be represented only by ids or also by explicit previous/next links
- whether learning paths belong in the same file as curriculum items or in a separate metadata file
