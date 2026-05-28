# Repository structure standard

Use neutral, durable folder names. Name folders by responsibility, not by the tool or agent that may operate on them.

## Top-level responsibilities

```text
go-engineer/
├── metadata/     # source of truth: graph, concepts, projects, assessments, contracts
├── curriculum/   # learner-facing readmes, code, tests, assets, labs, projects
├── tools/        # validation, generation, audit, migration, authoring automation
├── docs/         # maintainer documentation, standards, governance, release process
└── dist/         # generated release artifacts only; never hand-edit
```

Do not create folders named `codex`, `ai`, `chatgpt`, or model-specific names. Use `tools/authoring/` for reusable authoring prompts/procedures.

## Metadata layout

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

## Learner-facing curriculum layout

Every module has a module overview and typed content folders. Omit empty typed folders until needed if the repository policy allows it.

```text
curriculum/modules/{module}/
├── README.md
├── lessons/{lesson}/
├── labs/{lab}/
├── projects/{project}/
├── assessments/{assessment}/
└── assets/
```

Electives use the same typed shape:

```text
curriculum/electives/{elective}/
├── README.md
├── lessons/{lesson}/
├── labs/{lab}/
├── projects/{project}/
├── assessments/{assessment}/
└── assets/
```

A complete lesson folder normally contains:

```text
curriculum/modules/{module}/lessons/{lesson}/
├── README.md
├── main.go
├── main_test.go
├── _starter/
├── _solution/
└── assets/
    └── diagrams/
```

## Tools layout

```text
tools/
├── validate/
│   ├── curriculum/   # strict go validation backbone
│   └── repository/   # repository/file/content validators
├── generate/         # generate readmes, lessons, assessments, snapshots
├── audit/            # audit quality, completion, projects, assessments
├── migrate/          # v2 to v3 migration helpers
└── authoring/        # reusable authoring workflows and prompt files
```

## Canonical file paths in metadata

Use repo-root relative paths. Prefer explicit paths; do not make validators infer where content lives.

```json
{
  "files": {
    "readme_path": "curriculum/modules/08-http-rest-apis/lessons/03-handler-lifecycle/README.md",
    "main_path": "curriculum/modules/08-http-rest-apis/lessons/03-handler-lifecycle/main.go",
    "test_path": "curriculum/modules/08-http-rest-apis/lessons/03-handler-lifecycle/main_test.go",
    "starter_path": "curriculum/modules/08-http-rest-apis/lessons/03-handler-lifecycle/_starter",
    "solution_path": "curriculum/modules/08-http-rest-apis/lessons/03-handler-lifecycle/_solution",
    "assets_dir": "curriculum/modules/08-http-rest-apis/lessons/03-handler-lifecycle/assets"
  }
}
```

## Canonical item type mapping

Use metadata `type` to choose typed folder:

| item type | folder |
|---|---|
| `lesson` | `lessons/` |
| `lab` | `labs/` |
| `project` | `projects/` |
| `assessment` | `assessments/` |
| `checkpoint` | `assessments/` |

When type is unknown, default to `lessons/` only for scaffolding and emit a warning in audits.
