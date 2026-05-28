# Migration Guide

This guide defines how to migrate legacy v2 content into the v3 curriculum.

## Migration principle

Do not blindly preserve the old structure.

v3 uses a new source-of-truth model:

```text
metadata/     active v3 graph
curriculum/   learner-facing v3 content
tools/        migration, generation, audit, validation
```

Legacy v2 content is preserved for traceability under:

```text
metadata/legacy/
docs/legacy-v2/
```

## Legacy files

```text
metadata/legacy/curriculum.v2.json
metadata/legacy/curriculum.v2.lock.json
metadata/legacy/unmapped-v2-report.json
metadata/legacy/migration-notes.md
```

`curriculum.v2.lock.json` is the frozen migration input.

Do not edit it.

## Migration outcomes

Every v2 item must receive one outcome:

| Outcome | Meaning |
|---|---|
| `direct` | v2 item maps directly to a v3 item or project. |
| `rewrite` | v2 item is preserved conceptually but rewritten for v3. |
| `split` | v2 item becomes multiple v3 items. |
| `merge` | multiple v2 items become one stronger v3 item. |
| `move-to-elective` | content is preserved but removed from the required core path. |
| `replace` | v2 content is replaced by a better modern equivalent. |
| `archive` | v2 content is intentionally kept only for history. |

No v2 item may be left unmapped without an explicit report entry.

## Migration workflow

1. Freeze the v2 input.
2. Generate or update `unmapped-v2-report.json`.
3. Decide outcome for each v2 item.
4. Update v3 metadata.
5. Preserve source legacy IDs where relevant.
6. Move advanced content to electives when it should not block core progression.
7. Write or update learner-facing v3 content.
8. Run migration report.
9. Run metadata validation.
10. Run repository validation.

## What should move to electives

Move content to electives when it is:

- advanced but not required for zero-to-engineer core
- specialized to a particular architecture style
- useful but not foundational
- likely to distract beginners
- better taught after the learner has production context

Examples:

- advanced generics
- complex generic constraints
- gRPC internals
- service mesh
- Helm/Kubernetes operational depth
- CQRS/event sourcing patterns

## What should stay core

Keep content core when it is required for professional Go backend work:

- basic Go syntax
- errors
- testing
- packages/modules
- CLI and files
- HTTP
- SQL/Postgres
- authentication/security basics
- context and concurrency
- observability
- performance basics
- architecture tradeoffs
- Docker/CI/CD/deployment basics
- portfolio and interview readiness
- flagship integration

## Source legacy IDs

When a v3 item derives from v2, preserve traceability:

```json
{
  "source_legacy_ids": ["FE.4"]
}
```

For rewritten or merged content, include all source IDs that materially influenced the new item.

## Migration validation

A complete migration has:

- `legacy_v2_unmapped: 0`
- source SHA recorded
- every v2 item has an outcome
- every direct mapping points to a valid v3 target
- every policy mapping has a reason
- v2 docs moved out of active root
- v3 docs no longer claim v2 is active architecture

## Anti-patterns

Do not:

- keep two active architectures
- leave v2 root docs as active contracts
- mark content complete because it existed in v2
- copy old prose without updating path, pedagogy, and validation contract
- move hard content earlier just to preserve old order
