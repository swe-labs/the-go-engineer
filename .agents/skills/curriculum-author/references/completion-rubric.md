# Completion rubric

## Status levels

- **0 — missing**: metadata or files absent.
- **1 — skeleton**: file exists but contains placeholders or shallow generic text.
- **2 — draft**: main ideas present but weak explanation, missing code/tests/assets, or inconsistent references.
- **3 — teachable**: complete explanation and working artifacts, but needs editorial/technical polish.
- **4 — production quality**: clear, accurate, validated, tested, and aligned with project/assessment.
- **5 — world-class**: beginner-friendly, technically deep, memorable, validated, and portfolio/interview useful.

## 100% complete lesson gate

A lesson is complete only when all are true:

- metadata status is stable
- `zero_magic_status` is golden
- `readme_status` is golden
- metadata file paths use `curriculum/modules/.../lessons/...` or `curriculum/electives/.../lessons/...`
- README exists and follows the required contract
- required code files exist
- required tests exist and pass
- required starter/solution folders exist
- required assets/diagrams exist
- verification commands work or manual steps are explicit
- crossrefs resolve and have specific reasons
- concepts are introduced/reinforced correctly
- no placeholder text remains
- review questions test explanation, application, debugging, and tradeoffs

## 100% complete lab gate

A lab is complete only when:

- it lives under `curriculum/modules/{module}/labs/{lab}/` or `curriculum/electives/{elective}/labs/{lab}/`
- it has a clear task, constraints, expected output, and verification steps
- starter and solution materials exist when required
- it reinforces a named lesson/concept
- it can be completed without hidden assumptions

## 100% complete module gate

A module is complete only when:

- every lesson passes the lesson gate
- every lab passes the lab gate
- entry and terminal item chains are valid
- module project(s) pass the project gate when applicable
- module assessment passes the assessment gate
- cognitive load and pacing are appropriate
- no hidden prerequisite jumps exist
- module README exists and explains goals, sequence, artifacts, and completion criteria
- module folders use typed children: `lessons/`, `labs/`, `projects/`, `assessments/`, and `assets/`

## 100% complete project gate

A project is complete only when:

- it lives under `curriculum/modules/{module}/projects/{project}/` or `curriculum/electives/{elective}/projects/{project}/`
- it has a realistic scenario
- deliverables are concrete
- milestones are ordered
- verification is runnable or manually precise
- rubric weights add to 100 or 1.0 according to repo convention
- it reinforces named concepts and lessons
- portfolio guidance is included when appropriate
- failure modes are included

## 100% complete assessment gate

An assessment is complete only when:

- it lives under `curriculum/modules/{module}/assessments/{assessment}/` or `curriculum/electives/{elective}/assessments/{assessment}/`
- target IDs resolve
- criteria map to target lessons/concepts
- criteria weights are valid
- it includes explanation, code reading, debugging, implementation, and tradeoff/prod questions
- answer key or review guidance exists when the repository policy requires it

## 100% complete repository gate

The repository is complete only when strict repository validation passes and every module, lesson, lab, project, assessment, README, code file, test file, diagram, asset, and migration mapping passes its gate.
