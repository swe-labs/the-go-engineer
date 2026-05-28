# Curriculum Quality Validation Backbone

This repository package contains a replacement validation tool for `internal/tools/curriculum` plus metadata fixes needed for the current 100% metadata architecture.

## Validation profiles

### 1. Metadata profile

Use this when editing JSON metadata only.

```bash
go run ./internal/tools/curriculum validate-metadata --root .
```

This profile validates:

- JSON file presence and parseability.
- Core/elective separation.
- Module, item, project, assessment, concept, and cross-reference IDs.
- Module phases, cognitive load, and module-size guardrails.
- Item graph reachability and cycle safety.
- Concept ownership and reinforcement coverage.
- Assessment targets, criteria, evidence, and retake policies.
- Project bindings, anchors, deliverables, rubrics, verification, and portfolio metadata.
- Zero-magic required fields and placeholder literals.
- README contract references.

The current metadata package passes this profile with zero errors and zero warnings.

### 2. Repository content profile

Use this when implementing actual curriculum files.

```bash
go run ./internal/tools/curriculum validate-repository --root . --strict-repository
```

This profile validates the real repository contents declared by metadata:

- Every required `README.md` exists.
- Every required `main.go`, `main_test.go`, starter, solution, asset directory, and diagram exists.
- README quality includes problem framing, mental model, analogy, production context, Go code or commands, practice, review questions, and `NEXT UP` footer.
- Go code is formatted with `gofmt`.
- Test files contain real `Test*` functions.
- Runnable/test-required lessons declare verification commands.
- Visual-model-required lessons have at least one visual asset.

This profile is intentionally strict and should fail until the full lesson/code/asset repository is generated.

### 3. Full profile

Use this in final CI before release.

```bash
go run ./internal/tools/curriculum validate-all --root . --strict-repository
```

This runs both metadata and repository validation. The curriculum should not be released until this command passes.

## Removed/retired legacy behavior

The old tool assumed a `curriculum/` directory, mixed validators with transform scripts, duplicated `main` functions, and only checked a narrow set of content statuses. The new tool:

- Locates `metadata/`, `curriculum/`, or root-level metadata automatically.
- Keeps validation separate from migration/transformation.
- Supports JSON reports for CI.
- Uses explicit profiles for metadata vs full repository completion.
- Treats repository content validation as a release gate, not a metadata-only check.

## CI recommendation

Use two jobs:

```bash
# Fast metadata gate for every commit
go run ./internal/tools/curriculum validate-metadata --root . --json

# Full release gate for publish/release branches
go run ./internal/tools/curriculum validate-all --root . --strict-repository --json
```
