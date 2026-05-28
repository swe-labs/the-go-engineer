# Curriculum workflows

## Creation workflow

1. Locate the target metadata item by `item_id`, slug, module/order, or project/assessment ID.
2. Confirm the item has stable metadata: `module_id`, title, type, difficulty, learning objective, prerequisites, next items, content contract, verification, files, proof, and zero-magic fields.
3. Confirm paths use the final learner-facing layout:
   - `curriculum/modules/{module}/lessons/{lesson}/...`
   - `curriculum/modules/{module}/labs/{lab}/...`
   - `curriculum/modules/{module}/projects/{project}/...`
   - `curriculum/modules/{module}/assessments/{assessment}/...`
   - `curriculum/electives/{elective}/lessons/{lesson}/...`
4. Determine required artifacts from `content_contract` and `files`.
5. Create or update the lesson/lab/project/assessment folder.
6. Write the README using `lesson-template.md` and the quality standard.
7. Add runnable code and tests when required.
8. Add starter/solution folders and real assets/diagrams when required.
9. Update project/assessment/crossref/concept references if the item introduces or reinforces a concept.
10. Run metadata and repository validation.
11. Report exactly what changed and what remains.

## Audit workflow

1. Scope the audit: item, module, project, assessment, phase, migration, or whole repository.
2. Run `scripts/audit_curriculum.py --repo <repo> --strict-repository` when actual files should exist.
3. If a Go validator exists under `tools/validate/curriculum`, run the strongest relevant command, usually:
   `go run ./tools/validate/curriculum validate-all --strict-repository`.
4. Inspect failures and classify them:
   - **Blocker**: broken graph, missing required file, invalid target, placeholder content, incomplete assessment/project, missing tests for test-required item, legacy/non-canonical path, or missing migration trace.
   - **Major**: weak explanation, insufficient analogy, poor rubric, unclear proof task, shallow production notes, or weak debugging coverage.
   - **Minor**: naming inconsistency, weak footer, style issue, or missing optional diagram.
5. Produce a prioritized remediation plan with exact file paths.
6. Never report “100% complete” unless strict repository validation passes and the completion rubric passes.

## Completion workflow

1. Sort remaining gaps by dependency order: metadata → module README → lesson README → code → tests → assets → project → assessment → crossrefs → migration trace.
2. Complete one module at a time; within a module, complete entry-to-terminal order.
3. For every lesson/lab, ensure README, code, tests, and metadata agree.
4. For every project, verify it exercises multiple lessons and has rubric, deliverables, verification, and portfolio narrative.
5. For every assessment, verify criteria map to target lessons and require explanation, code reading, debugging, implementation, and tradeoff reasoning.
6. Re-run validation after each module.
7. Maintain a machine-readable completion report in `dist/completion-report.json` or the repo’s configured report path.

## Migration workflow

1. Load legacy v2 material from `metadata/legacy/` and current v3 metadata.
2. Classify each legacy item as `keep`, `rewrite`, `merge`, `move-to-elective`, or `remove`.
3. Preserve traceability via `source_legacy_ids` and migration mappings.
4. Do not import weak v2 content unchanged. Rewrite it to zero-magic standard.
5. Move advanced/non-blocking topics to electives unless required for core job readiness.
6. Map migrated content to final learner-facing paths under `curriculum/modules/` or `curriculum/electives/`.
7. Validate that migrated content has no duplicate lessons, broken prerequisites, hidden dependency jumps, or orphaned legacy IDs.
8. Produce a migration report listing every v2 item and its final disposition.
