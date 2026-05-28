# Go Engineer Curriculum

Go Engineer is a zero-magic, metadata-driven curriculum for taking a learner from beginner to production-ready Go software engineer.

The repository is intentionally split into five responsibilities:

```text
metadata/     Source of truth for graph, concepts, projects, assessments, crossrefs, contracts, and migration.
curriculum/   Learner-facing content: READMEs, code, tests, labs, projects, assessments, diagrams, and assets.
tools/        Validation, generation, audit, migration, and authoring automation.
docs/         Maintainer documentation, standards, governance, and release process.
dist/         Generated release artifacts only. Never hand-edit.
```

## Current architecture

Use typed learner-facing folders. Do not mix lessons, labs, projects, and assessments in one flat module folder.

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

## Canonical metadata path convention

Metadata must point to explicit learner-facing files:

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

## Quality gates

A section, lesson, lab, project, or assessment is complete only when its metadata and learner-facing files both pass validation.

Metadata must prove:

- every module, item, project, assessment, concept, and cross-reference has a stable ID
- every prerequisite and next item resolves
- every project and assessment targets real curriculum objects
- every concept has a canonical owner and reinforcement plan
- no placeholder status remains
- no generic cross-reference explanation remains
- module phases and cognitive-load flags are accurate
- migration from v2 is traceable

Learner-facing content must provide:

- a complete README for each module and lesson
- runnable Go code where required
- tests for behavior-heavy lessons and projects
- starter and solution files where learners implement code
- diagrams/assets when required by the contract
- mental models, analogies, machine views, common mistakes, debugging guidance, production context, practice, and proof of understanding

## Commands

Run metadata validation:

```bash
make validate-metadata
```

Run strict learner-facing repository validation:

```bash
make validate-repository
```

Run the full release gate:

```bash
make validate-release
```

Generate release artifacts:

```bash
make release-artifacts
```

## Root cleanup policy

Keep the repository root small. Root should contain only project entry points and release-critical files:

```text
README.md
go.mod
go.sum
.gitignore
LICENSE
Makefile
.github/
metadata/
curriculum/
tools/
docs/
dist/
```

Older v2 root documents belong under `docs/legacy-v2/` or `metadata/legacy/`, not at the root.

Do not keep tool-branded maintainer prompt files in root. Authoring workflows belong in `tools/authoring/`, and reusable assistant behavior belongs in the packaged Skill.
