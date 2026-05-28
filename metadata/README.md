# Metadata

This directory is the active source of truth for the Go Engineer v3 curriculum graph.

Metadata defines what exists, how it is ordered, what each item proves, which concepts are owned where, which assessments and projects validate mastery, and how legacy v2 content maps into v3.

Learner-facing prose, source code, tests, diagrams, and project files belong under `curriculum/`, not here.

## Files

| File | Role |
|---|---|
| `workspace.json` | Workspace manifest, source-of-truth list, validation commands, repository layout policy. |
| `schema.v3.json` | JSON schema for split v3 metadata documents. |
| `path.core.json` | Core module and lesson graph. |
| `path.electives.json` | Elective module and lesson graph. |
| `concepts.json` | Canonical concept ownership, previews, and reinforcement locations. |
| `projects.json` | Project definitions, deliverables, rubrics, verification, and canonical learner-facing paths. |
| `assessments.json` | Module and project assessment definitions, criteria, evidence, retake policy, and canonical learner-facing paths. |
| `crossrefs.json` | Global cross-reference graph and semantic relationship reasons. |
| `failures.json` | Operational failure taxonomy and module-level failure coverage requirements. |
| `readme.contracts.json` | README contracts for modules, lessons, projects, assessments, electives, and flagship checkpoints. |
| `migration.v2-to-v3.json` | v2-to-v3 migration policy and item coverage summary. |
| `VALIDATION.metadata.json` | Latest metadata validation report. |
| `legacy/` | Frozen v2 migration inputs and migration coverage reports. |

## Canonical learner-facing paths

Metadata must reference explicit typed curriculum paths:

```text
curriculum/modules/{module}/lessons/{lesson}/README.md
curriculum/modules/{module}/labs/{lab}/README.md
curriculum/modules/{module}/projects/{project}/README.md
curriculum/modules/{module}/assessments/{assessment}/README.md

curriculum/electives/{elective}/lessons/{lesson}/README.md
curriculum/electives/{elective}/labs/{lab}/README.md
curriculum/electives/{elective}/projects/{project}/README.md
curriculum/electives/{elective}/assessments/{assessment}/README.md
```

## Completion standard

Metadata is complete only when:

- all JSON files parse
- every item has a valid module
- every prerequisite and `next_item_id` resolves
- every cross-reference target resolves
- every concept owner and reinforcement target resolves
- every project target and assessment target resolves
- every project and assessment has a canonical learner-facing path
- every required v2 item has a migration outcome
- no placeholder status remains
- no generic cross-reference reason remains
- canonical typed curriculum paths are used

Run:

```bash
go run ./tools/validate/curriculum validate-metadata
```
