# V2 Folder Structure

## Purpose

The v2 folder structure must improve consistency without causing unnecessary migration pain.
This means the structure should be conservative, predictable, and validator-friendly.

## Governing Rule

The first v2 migration waves should preserve the current top-level section directories.

That means keeping:

- `01-core-foundations`
- `02-control-flow`
- ...
- `15-code-generation`

The goal is to improve the contract inside the repo before renaming large trees.

## Root Layout

The recommended root structure remains:

```text
/
  01-core-foundations/
  02-control-flow/
  ...
  15-code-generation/
  docs/
    curriculum/
    v2/
  scripts/
  curriculum.json
```

If a richer v2 schema is adopted later, it should live alongside current metadata during migration.

## Section Layout

Each section should follow this structure:

```text
NN-section-name/
  README.md
  [optional-subgraph/]N-lesson-name/
  [optional-subgraph/]N-drill-name/
  [optional-subgraph/]N-exercise-name/
  [optional-subgraph/]N-checkpoint-name/
  [optional-subgraph/]N-mini-project-name/
  _shared/                    # optional support code or fixtures for the section
```

Subgraphs are allowed when a section naturally contains clear internal tracks such as:

- `getting-started`
- `language-basics`
- `filesystem`
- `encoding`

They should exist to help navigation, not to create extra nesting for its own sake.

## Lesson Directory Layout

Default lesson layout:

```text
N-lesson-name/
  main.go
  README.md        # optional
  main_test.go     # optional
  testdata/        # optional
```

## Drill Directory Layout

Default drill layout:

```text
N-drill-name/
  README.md        # optional
  main.go          # optional when the drill is prompt-first
  main_test.go     # optional when quick verification helps
```

## Exercise Directory Layout

Default exercise layout:

```text
N-exercise-name/
  README.md
  main.go
  _starter/
    main.go
  main_test.go     # optional
  testdata/        # optional
```

## Project Directory Layout

For projects that need multiple packages, the allowed layout is:

```text
N-project-name/
  README.md
  cmd/
  internal/
  pkg/             # only when a real public package boundary exists
  testdata/
  go.mod           # optional, only if isolation is genuinely valuable
```

Projects should not default to their own module unless that isolation helps teaching or tooling.

## Reference Placement

Reference content usually belongs in:

- `README.md`
- learner support docs
- setup or tooling docs

Only place reference content in numbered lesson-style directories when it is explicitly part of the
canonical learning path.

## Naming Rules

Use these naming rules consistently:

- root sections stay two-digit prefixed
- local lesson and exercise directories stay numerically ordered
- slugs are lowercase and hyphen-separated
- do not rename v1 directories only for aesthetic reasons during early migration

## Placement Rules

Use these placement rules:

- keep a lesson next to the practice that depends on it
- keep checkpoints near the end of a local track
- keep mini-projects inside the section or phase they validate
- do not create a distant top-level `projects/` tree during early migration

## Shared Assets

`_shared/` is allowed at the section level when:

- fixtures are reused across several lessons
- a helper package supports multiple exercises
- duplication would obscure the teaching goal

Do not create `_shared/` just to simulate enterprise layering.

## Migration Rule

During the first migration waves:

- preserve current links where possible
- add structure incrementally
- move only when the new structure clearly improves navigation, validation, or teaching clarity

## Validator Implications

The folder structure should be simple enough for validator rules to enforce:

- section directory exists
- lesson path exists
- run and test commands resolve
- starter directories exist when declared
- optional files are only required when metadata says they should exist
