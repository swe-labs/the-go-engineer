# Metadata Contract

This document defines the contract for all files under `metadata/`.

Metadata is complete only when it can drive generation, validation, audit, migration, and release without guessing.

## Required files

```text
metadata/
├── README.md
├── workspace.json
├── schema.v3.json
├── path.core.json
├── path.electives.json
├── concepts.json
├── projects.json
├── assessments.json
├── crossrefs.json
├── failures.json
├── readme.contracts.json
├── migration.v2-to-v3.json
├── VALIDATION.metadata.json
└── legacy/
    ├── curriculum.v2.json
    ├── curriculum.v2.lock.json
    ├── unmapped-v2-report.json
    └── migration-notes.md
```

## File responsibilities

### `workspace.json`

Defines the active workspace policy:

- source-of-truth files
- validation commands
- repository layout
- forbidden folder names
- root cleanup rules
- legacy inputs

### `schema.v3.json`

Defines allowed document types and common metadata constraints.

It should validate split metadata documents without forcing all curriculum information into one monolithic file.

### `path.core.json`

Defines the required core curriculum graph:

- modules
- items
- ordering
- prerequisites
- next item links
- zero-magic metadata
- proof metadata
- content contracts
- verification commands
- canonical files paths
- source legacy IDs

### `path.electives.json`

Defines optional advanced content.

Electives must not block the core path unless a module explicitly declares an elective prerequisite.

### `concepts.json`

Defines canonical concept ownership.

Each concept should include:

- `concept`
- `canonical_owner`
- `preview_locations`
- `reinforcement_locations`

A concept is incomplete if it is introduced but never reinforced.

### `projects.json`

Defines portfolio and integration proof.

Each project must include:

- ID
- module ownership
- title
- purpose
- deliverables
- rubric
- verification
- assessment binding
- canonical learner-facing file paths

### `assessments.json`

Defines module and project assessments.

Each assessment must include:

- ID
- target IDs
- criteria
- evidence requirements
- scoring or rubric policy
- retake policy
- canonical learner-facing file paths

### `crossrefs.json`

Defines semantic relationships.

Every cross-reference reason must explain the actual relationship. Generic reasons are invalid.

Good:

```text
core-08-03 builds on core-07-04 because request parsing reuses JSON decoding and explicit error boundaries.
```

Bad:

```text
Both lessons build understanding in this module.
```

### `failures.json`

Defines failure engineering coverage.

Every technical module should identify relevant failures, such as:

- configuration failures
- timeout failures
- partial failures
- resource exhaustion
- security failures
- deployment failures
- data consistency failures

### `readme.contracts.json`

Defines required README sections and quality rules.

Contracts must match the actual learner-facing content structure under `curriculum/`.

### `migration.v2-to-v3.json`

Defines how v2 content maps into v3.

Every legacy item must have one outcome:

- `direct`
- `rewrite`
- `split`
- `merge`
- `move-to-elective`
- `replace`
- `archive`

### `VALIDATION.metadata.json`

Stores the latest metadata validation result.

It should be generated, not hand-written.

### `legacy/`

Stores frozen v2 inputs and migration reports. Legacy files are not active v3 source of truth.

## Required item fields

Each core/elective item should define:

```text
id
module_id
slug
title
type
subtype
status
difficulty
phase
order
estimated_minutes
learning_objective
required_prior_knowledge
prerequisites
next_item_ids
zero_magic
crossrefs
proof
content_contract
verification
files
source_legacy_ids
tags
```

## Required zero-magic fields

A lesson is not world-class unless it has:

```text
problem_solved
why_it_exists
mental_model
under_the_hood
how_go_uses_it
real_world_usage
beginner_mistakes
debugging_signals
performance_implications
security_implications
execution_timeline
proof_of_understanding
```

Some fields may be intentionally short for non-code orientation lessons, but they must not be missing.

## Canonical files contract

Lessons:

```json
{
  "files": {
    "readme_path": "curriculum/modules/{module}/lessons/{lesson}/README.md",
    "main_path": "curriculum/modules/{module}/lessons/{lesson}/main.go",
    "test_path": "curriculum/modules/{module}/lessons/{lesson}/main_test.go",
    "starter_path": "curriculum/modules/{module}/lessons/{lesson}/_starter",
    "solution_path": "curriculum/modules/{module}/lessons/{lesson}/_solution",
    "assets_dir": "curriculum/modules/{module}/lessons/{lesson}/assets"
  }
}
```

Labs:

```json
{
  "files": {
    "readme_path": "curriculum/modules/{module}/labs/{lab}/README.md"
  }
}
```

Projects:

```json
{
  "files": {
    "readme_path": "curriculum/modules/{module}/projects/{project}/README.md",
    "starter_path": "curriculum/modules/{module}/projects/{project}/_starter",
    "solution_path": "curriculum/modules/{module}/projects/{project}/_solution",
    "test_path": "curriculum/modules/{module}/projects/{project}/tests",
    "assets_dir": "curriculum/modules/{module}/projects/{project}/assets"
  }
}
```

Assessments:

```json
{
  "files": {
    "readme_path": "curriculum/modules/{module}/assessments/{assessment}/README.md",
    "questions_path": "curriculum/modules/{module}/assessments/{assessment}/questions.md",
    "answer_key_path": "curriculum/modules/{module}/assessments/{assessment}/answer-key.md",
    "rubric_path": "curriculum/modules/{module}/assessments/{assessment}/rubric.md"
  }
}
```

## Graph rules

- Module IDs must be unique.
- Item IDs must be unique across core and electives.
- Item prerequisites must resolve to valid item or module IDs.
- `next_item_ids` must resolve to valid item or project IDs.
- Module entry and terminal items must exist.
- Assessment targets must exist.
- Project assessment IDs must exist.
- Cross-reference targets must exist.
- Concept owners must exist.
- Legacy source IDs should be traceable when content came from v2.

## Stability rules

Before release:

- no `draft` status in required metadata
- no `placeholder` zero-magic status
- no `scaffolded` README status
- no generic cross-reference reasons
- no unmapped required v2 items
- no noncanonical learner-facing paths
