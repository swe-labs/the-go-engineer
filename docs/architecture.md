# Architecture

This document defines the active Go Engineer repository architecture.

The architecture optimizes for:

- clear source-of-truth boundaries
- beginner-friendly learner navigation
- strict validation
- repeatable generation
- safe v2 to v3 migration
- section-by-section and lesson-by-lesson completion

## Top-level layout

```text
go-engineer/
├── README.md
├── go.mod
├── go.sum
├── .gitignore
├── LICENSE
├── Makefile
├── .github/
├── metadata/
├── curriculum/
├── tools/
├── docs/
└── dist/
```

## Directory responsibilities

### `metadata/`

`metadata/` is the active curriculum source of truth.

It defines:

- modules
- lessons
- labs
- projects
- assessments
- concepts
- cross-references
- failure coverage
- README contracts
- migration coverage
- validation status

It must not contain learner-facing long-form lesson prose, runnable lesson code, or diagrams except as metadata references.

### `curriculum/`

`curriculum/` is the learner-facing product.

It contains:

- module overviews
- lesson READMEs
- lab READMEs
- project READMEs
- assessment READMEs
- starter code
- solution code
- tests
- diagrams
- shared guides and reference assets

### `tools/`

`tools/` contains all automation:

```text
tools/validate/      Validation gates.
tools/generate/      File generation from metadata.
tools/audit/         Quality and completion audits.
tools/migrate/       v2 to v3 migration tooling.
tools/authoring/     Human-readable authoring workflows.
```

Use neutral responsibility names. Do not use tool-branded folder names such as `tool-specific`, `ai`, `agent`, `bot`, or `llm`.

### `docs/`

`docs/` contains maintainer-facing policy and process documents. These documents explain the standard but do not override metadata.

### `dist/`

`dist/` contains generated release artifacts only. Never hand-edit files in `dist/`.

## Learner-facing curriculum layout

Use typed folders inside each module. Do not place all lessons, labs, projects, and assessments in one flat module directory.

```text
curriculum/modules/{module}/README.md
curriculum/modules/{module}/lessons/{lesson}/README.md
curriculum/modules/{module}/labs/{lab}/README.md
curriculum/modules/{module}/projects/{project}/README.md
curriculum/modules/{module}/assessments/{assessment}/README.md
curriculum/modules/{module}/assets/
```

Electives follow the same typed pattern:

```text
curriculum/electives/{elective}/README.md
curriculum/electives/{elective}/lessons/{lesson}/README.md
curriculum/electives/{elective}/labs/{lab}/README.md
curriculum/electives/{elective}/projects/{project}/README.md
curriculum/electives/{elective}/assessments/{assessment}/README.md
curriculum/electives/{elective}/assets/
```

Shared learner resources live here:

```text
curriculum/shared/
├── README.md
├── glossary.md
├── zero-magic-guide.md
├── debugging-guide.md
├── command-cheatsheet.md
├── go-style-guide.md
├── testing-guide.md
├── security-guide.md
├── deployment-guide.md
├── portfolio-rubric.md
├── assessment-rubric.md
├── project-quality-standard.md
└── assets/
```

## Canonical module list

The core path is organized as:

```text
00-orientation
01-computers-terminal-git-web
02-go-setup-tooling
03-programming-fundamentals
04-functions-errors-data-semantics
05-types-interfaces-packages-modules
06-testing-debugging-refactoring
07-cli-files-json-config
08-http-rest-apis
09-sql-postgres-persistence
10-auth-security
11-lifecycle-context-concurrency
12-observability-diagnostics
13-performance-memory-engineering
14-architecture-distributed-systems
15-docker-cicd-deployment
16-portfolio-interview-readiness
18-flagship-opslane
```

Electives:

```text
advanced-electives
```

## Path policy

Metadata paths must be explicit and rooted at `curriculum/`.

Good:

```json
{
  "readme_path": "curriculum/modules/08-http-rest-apis/lessons/03-handler-lifecycle/README.md"
}
```

Bad:

```json
{
  "readme_path": "08-http-rest-apis/03-handler-lifecycle/README.md"
}
```

Bad:

```json
{
  "readme_path": "modules/08-http-rest-apis/03-handler-lifecycle/README.md"
}
```

The validator must not have to guess where learner-facing content lives.

## Root cleanup policy

The root stays small.

Allowed root files:

```text
README.md
go.mod
go.sum
.gitignore
LICENSE
Makefile
```

Allowed root directories:

```text
.github/
metadata/
curriculum/
tools/
docs/
dist/
```

Forbidden root files:

```text
.env
AGENTS.md
```

`.env` files are secrets risk. `AGENTS.md` is redundant because reusable assistant behavior belongs in the packaged Skill and neutral authoring workflows belong in `tools/authoring/`.

## Empty folder policy

Typed folders may be omitted until content exists.

For example, a module may start as:

```text
curriculum/modules/00-orientation/
├── README.md
└── lessons/
```

A full module may later become:

```text
curriculum/modules/08-http-rest-apis/
├── README.md
├── lessons/
├── labs/
├── projects/
├── assessments/
└── assets/
```

Metadata paths must still use typed folders even before the files are generated.
