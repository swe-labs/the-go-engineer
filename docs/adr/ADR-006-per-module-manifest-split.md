# ADR-006: Per-Module Manifest Split

- Status: accepted
- Date: 2026-11-14

## Context

Through v3.1.0 the authoritative curriculum metadata lived in two large manifests: `curriculum/path.core.json` (68,730 lines at peak — all 421 items, 18 module records, and flagship milestones in one array) and `curriculum/path.electives.json`. Three structural redundancies compounded the maintenance pain:

- `crossrefs.json` (18,861 lines) byte-mirrored crossref data already embedded in every item (2,086/2,086 entries identical), while electives had always carried embedded-only data.
- `projects.json` and `assessments.json` duplicated records the path manifests also carried, with divergent copies of `verification` and `order` fields.
- Cross-file boilerplate (`architectural_decision`, `quality_gates`/`validation_commands` with two conflicting gate-name vocabularies) was duplicated across manifests.

Editing, reordering, including, or excluding a lesson required scripted edits into a file too large to review, and any two parallel work streams touching manifests produced merge conflicts. A dead-weight audit also found ~9.7k lines of provably unread fields, dangling references, and double-encoded UTF-8 prose (Phase A, PR #588).

## Decision

1. **Per-module manifests.** Each module directory carries `manifest.json` beside its `README.md`, holding that module's `items`, folded `projects`, folded `assessments`, and the items' embedded crossrefs. `curriculum/electives/advanced-electives/manifest.json` serves module-17.
2. **Thin indexes.** `path.core.json` and `path.electives.json` remain (same filenames, so directory discovery and existing citations stay valid) but hold only the header and the module registry: ordering, prerequisites, entry/terminal item ids, and module metadata.
3. **Crossrefs embed-only.** The central `crossrefs.json` registry is deleted. Embedded per-item crossrefs are the single source; `validate-crossrefs` checks embedded targets directly and the reconciliation-against-registry step is retired. The `builds_on` backward-pointing DAG in `validate-graph` derives from embedded data for all items.
4. **Projects/assessments folded.** `projects.json` and `assessments.json` are deleted; records live in the module manifests keyed by `module_id`. Cross-module references (placement anchors, prerequisites into other modules) keep working because item ids are global.
5. **Centralized loading.** `internal/tools/curriculum/loader.go` assembles the core/elective bundles once (cached) and exposes `moduleManifestPaths()`; gates consume the assembled bundles or iterate manifests with their own local raw shapes.

Phase A (PR #588) preceded the split to remove dead weight first, so the split migrated the final schema exactly once.

## Consequences

- Reordering, adding, or removing a lesson is a single-module-file edit; parallel work on different modules can no longer conflict on a shared manifest.
- `validate-schema` header checks and the raw-byte key-order check run per module manifest; `expected.json` inventory pins are unchanged and computed over assembled bundles.
- The catalog's repository-structure listing and manifest counts changed (reviewed deltas); `generate-catalog` remains byte-pinned by test.
- `concepts.json`, `failures.json`, `readme.contracts.json`, `expected.json`, and `workspace.json` remain global registries — their entries are inherently cross-module.
- Adding a module = new directory + `manifest.json` + one index row; deleting the module manifests of a module removes its records atomically.
