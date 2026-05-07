# Opslane Migration Policy

## Purpose

Keep schema evolution deterministic across local runs, CI, and release branches.

## Rules

1. Migrations are append-only and versioned sequentially (`001`, `002`, ...).
2. Every `*.up.sql` migration must have a matching `*.down.sql`.
3. Checksums are recorded in `schema_migrations` and validated before apply.
4. Dirty migration state blocks new runs until resolved by an operator.
5. Startup migration path and standalone migration runner must remain equivalent in effect.

## Operational Workflow

```bash
go run ./11-flagship/01-opslane/scripts/migrate.go -direction status
go run ./11-flagship/01-opslane/scripts/migrate.go -direction up
go run ./11-flagship/01-opslane/scripts/migrate.go -direction down
```

## Failure Handling

- **Checksum mismatch**: stop rollout, verify migration file integrity, and reconcile environment state.
- **Dirty state**: resolve manually (rollback/repair) before applying new migrations.
- **Order gap**: add missing migration version or renumber before merge.

## Review Requirements

Migration PRs should include:

- SQL change intent and rollback plan
- proof that `up` and `down` both execute
- validator and test evidence
