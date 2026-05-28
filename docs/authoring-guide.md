# Authoring Guide

This guide explains how to create or update modules, lessons, labs, projects, and assessments.

## Authoring order

Always work in this order:

1. Update metadata.
2. Generate or create learner-facing files.
3. Write complete content.
4. Add code/tests/assets.
5. Run validation.
6. Audit quality.
7. Review and fix findings.

Do not write isolated files that are not represented in metadata.

## Creating a module

1. Add module metadata to `metadata/path.core.json` or `metadata/path.electives.json`.
2. Use a canonical module path:

```text
curriculum/modules/{module}/
curriculum/electives/{elective}/
```

3. Add entry and terminal item IDs.
4. Add prerequisites.
5. Add concepts and failure coverage.
6. Create module README.
7. Add lessons, labs, projects, and assessments.

Module README must include:

```text
Mission
Who this module is for
What you will build
Prerequisites
Concept map
Lessons
Labs
Projects
Assessments
Common failure modes
Completion checklist
Next module
```

## Creating a lesson

1. Add item metadata.
2. Assign a module.
3. Define prerequisites and next item.
4. Add zero-magic fields.
5. Add proof fields.
6. Add canonical files paths.
7. Create the learner-facing folder:

```text
curriculum/modules/{module}/lessons/{lesson}/
```

8. Write README.
9. Add code/test/starter/solution/assets as required.
10. Validate.

Lesson README required sections:

```text
Mission
Prerequisites
Mental Model
Visual Model
Machine View
Run Instructions
Code Walkthrough
Try It
Common Mistakes
Debugging Signals
In Production
Performance Notes
Security Notes
Thinking Questions
Proof of Understanding
Next Step
```

## Creating a lab

Labs are guided practice.

Use:

```text
curriculum/modules/{module}/labs/{lab}/README.md
```

A lab must have:

- clear scenario
- tasks
- hints
- verification
- reflection
- extension

Labs should reduce ambiguity, not increase it.

## Creating a project

Projects prove integrated ability.

Use:

```text
curriculum/modules/{module}/projects/{project}/
```

A project must have:

- README
- starter path when implementation is required
- solution path when learners need reference
- tests
- rubric
- assessment binding
- portfolio guidance

Project metadata must include an `assessment_id`.

## Creating an assessment

Use:

```text
curriculum/modules/{module}/assessments/{assessment}/
```

Assessment files:

```text
README.md
questions.md
answer-key.md
rubric.md
assets/
```

Assessments must include:

- purpose
- scope
- allowed resources
- tasks/questions
- evidence required
- scoring rubric
- passing standard
- retake policy

## Updating content

When updating a lesson, check all connected surfaces:

- metadata item
- README
- main.go
- main_test.go
- `_starter/`
- `_solution/`
- assets
- crossrefs
- concepts
- assessments
- projects
- module README

A change is incomplete if one surface disagrees with another.

## Naming rules

Use lowercase kebab-case for folders:

```text
03-handler-lifecycle
01-build-first-handler
http-api-project
```

Use stable IDs in metadata:

```text
core-08-03
project-http-api
assessment-module-08
```

Do not rename stable IDs unless migration metadata is updated.

## Content depth rule

Do not ship a lesson that only says what to type.

A complete lesson explains:

- why the concept exists
- what problem it solves
- how Go represents it
- how it behaves at runtime
- where it fails
- how to debug it
- how it appears in production
- how the learner proves understanding

## Author checklist

Before opening a PR:

- [ ] Metadata updated.
- [ ] Canonical files paths used.
- [ ] README complete.
- [ ] Code compiles.
- [ ] Tests pass.
- [ ] Starter/solution present if required.
- [ ] Assets present if referenced.
- [ ] Concepts updated.
- [ ] Crossrefs meaningful.
- [ ] Assessment/project proof updated.
- [ ] Validation passes.
