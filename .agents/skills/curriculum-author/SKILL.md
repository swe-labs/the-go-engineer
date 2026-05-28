---
name: curriculum-author
description: create, audit, migrate, validate, and complete a zero-magic software-engineering curriculum from metadata and repository files. use when asked to create or revise curriculum modules, lessons, labs, projects, assessments, readmes, starter/solution code, tests, diagrams, assets, or migration plans; audit quality and completion; identify gaps; enforce world-class teaching standards; or help an authoring workflow implement the full curriculum repository end to end.
---

# Curriculum Author

## Operating principles

Treat `metadata/` as the source of truth and `curriculum/` as the learner-facing implementation. A lesson, lab, project, or assessment is not complete until metadata, README, code, tests, assets, diagrams, project links, assessments, and cross-references all agree.

Use the final neutral repository architecture. Do not introduce folders named `codex/`, `ai/`, `chatgpt/`, or tool-branded workflow names. Put reusable authoring prompts and procedures under `tools/authoring/`.

Never mark a section, lesson, project, assessment, module, or repository complete unless it satisfies `references/completion-rubric.md` and passes the strongest available validation. Prefer strict failure over optimistic status.

Use this skill to perform four kinds of work:

1. **Create** a module, lesson, lab, project, assessment, README, code file, test, or asset scaffold from metadata.
2. **Audit** a module, lesson, lab, project, assessment, phase, migration, or whole repository.
3. **Complete** missing files section-by-section or lesson-by-lesson.
4. **Migrate** v2 content into the v3 metadata and learner-facing repository architecture.

## Canonical repository layout

Use this responsibility map:

```text
metadata/     source of truth: graph, concepts, projects, assessments, contracts
curriculum/   learner-facing content: readmes, code, tests, assets, labs, projects
tools/        validation, generation, audit, migration, and authoring automation
docs/         maintainer documentation, standards, governance, release process
dist/         generated release artifacts only; never hand-edit
```

Use typed learner-facing module folders:

```text
curriculum/modules/{module}/README.md
curriculum/modules/{module}/lessons/{lesson}/README.md
curriculum/modules/{module}/labs/{lab}/README.md
curriculum/modules/{module}/projects/{project}/README.md
curriculum/modules/{module}/assessments/{assessment}/README.md
curriculum/modules/{module}/assets/

curriculum/electives/{elective}/README.md
curriculum/electives/{elective}/lessons/{lesson}/README.md
curriculum/electives/{elective}/labs/{lab}/README.md
curriculum/electives/{elective}/projects/{project}/README.md
curriculum/electives/{elective}/assessments/{assessment}/README.md
curriculum/electives/{elective}/assets/
```

Metadata file paths must be explicit and repo-root relative, for example:

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

See `references/repository-structure.md` before changing paths or generating folders.

## Required workflow

1. **Discover the repo layout.** Locate `metadata/` first. Confirm learner-facing content lives under `curriculum/modules/` and `curriculum/electives/`. Confirm tools live under `tools/`, especially `tools/validate/`, `tools/generate/`, `tools/audit/`, `tools/migrate/`, and `tools/authoring/`.
2. **Load metadata.** Inspect `path.core.json`, `path.electives.json`, `projects.json`, `assessments.json`, `concepts.json`, `crossrefs.json`, `failures.json`, `readme.contracts.json`, `migration.v2-to-v3.json`, and `workspace.json` when available.
3. **Choose workflow.**
   - Creating content: follow `references/workflows.md#creation-workflow`.
   - Auditing content: follow `references/workflows.md#audit-workflow`.
   - Completing repository files: follow `references/workflows.md#completion-workflow`.
   - Migrating v2: follow `references/workflows.md#migration-workflow`.
4. **Use deterministic checks.** Run `scripts/audit_curriculum.py` before and after major changes. Use `scripts/scaffold_lesson.py` to create missing lesson skeletons when appropriate.
5. **Write world-class content.** Use `references/quality-standard.md`, `references/lesson-template.md`, and `references/completion-rubric.md`.
6. **Validate strictly.** If a repo validator exists under `tools/validate/curriculum`, run it. Prefer:
   - `go run ./tools/validate/curriculum validate-metadata`
   - `go run ./tools/validate/curriculum validate-repository --strict-repository`
   - `go run ./tools/validate/curriculum validate-all --strict-repository`
   If not available, run the bundled audit script and report remaining gaps.
7. **Report status honestly.** Separate metadata completeness, README completeness, code/test completeness, asset/diagram completeness, project/assessment completeness, and migration completeness.

## Output rules

For audits, provide:

- completion status by module and lesson
- blockers vs warnings
- exact files to create or change
- metadata/repository mismatches
- priority order for fixes
- validation commands run and results
- whether strict repository validation passes

For generated lessons, produce or edit actual files, not only prose in chat. Every generated lesson must include:

- README with concept explanation, analogy, mental model, under-the-hood behavior, Go/runtime details, common mistakes, debugging, production notes, exercises, review questions, and NEXT UP footer
- runnable example when `content_contract.runnable_required` is true
- tests when `content_contract.tests_required` is true
- starter and solution folders when metadata requires them
- real diagrams/assets when metadata requires them; never call an empty placeholder complete

For migration, always preserve traceability from legacy IDs to v3 modules/items and document keep/merge/rewrite/move-to-elective/remove decisions.

## Scripts

- `scripts/audit_curriculum.py` audits metadata and repository files. It can emit JSON or Markdown.
- `scripts/scaffold_lesson.py` creates README/code/test/starter/solution/assets skeletons from a lesson ID using the final typed `curriculum/modules/.../{lessons,labs,projects,assessments}/...` layout.

Example commands:

```bash
python scripts/audit_curriculum.py --repo . --format markdown --out AUDIT.curriculum.md
python scripts/audit_curriculum.py --repo . --strict-repository --format json --out AUDIT.curriculum.json
python scripts/scaffold_lesson.py --repo . --item-id core-08-03 --overwrite-missing-only
```

## Reference loading guide

Load only what is needed:

- `references/repository-structure.md` before editing paths or directory layout.
- `references/workflows.md` for process decisions.
- `references/quality-standard.md` before writing or judging content.
- `references/lesson-template.md` before creating README content.
- `references/completion-rubric.md` before declaring anything complete.
- `references/migration-guide.md` when migrating v2.
