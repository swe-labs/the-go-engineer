# v2 to v3 Migration Notes

This directory preserves the previous v2 curriculum registry for traceability while the active v3 curriculum uses split metadata files under `metadata/`.

## Files

| File | Purpose |
|---|---|
| `curriculum.v2.json` | Original v2 registry preserved for migration reference. |
| `curriculum.v2.lock.json` | Byte-identical frozen copy used as the stable migration input. |
| `unmapped-v2-report.json` | Machine-readable coverage report for every v2 item. |
| `migration-notes.md` | Human-readable explanation of the migration approach. |

## Frozen source

```text
metadata/legacy/curriculum.v2.lock.json
sha256: e429366ee3b24b88ae8bc8a3a650362fa6e2087495f0728152aaac66c0cdf2fb
```

## Migration policy

The v3 curriculum intentionally breaks the old v2 public section spine and rebuilds around the v3 zero-to-software-engineer path.

Each v2 item must have one of these outcomes:

- `direct`
- `rewrite`
- `split`
- `merge`
- `move-to-elective`
- `replace`
- `archive`

The current report covers all 215 v2 items.

```text
directly mapped: 83
covered by policy: 132
unmapped: 0
```

## Active v3 source of truth

Do not edit legacy files as the active curriculum source. Update these files instead:

```text
metadata/path.core.json
metadata/path.electives.json
metadata/concepts.json
metadata/projects.json
metadata/assessments.json
metadata/crossrefs.json
metadata/failures.json
metadata/readme.contracts.json
metadata/migration.v2-to-v3.json
```
