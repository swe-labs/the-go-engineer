# File Organizer

## Learning objective

Build a production-grade directory organization CLI tool in Go. Master filesystem traversal (`filepath.WalkDir`), file attribute inspection, cross-device file move semantics, naming collision resolution, and command-line flag parsing with the standard library `flag` package.

## Why this matters

Backend services and data pipelines constantly ingest unorganized files from uploads, batch jobs, or external object stores. Writing safe filesystem automation requires atomic directory creation, handling filename collisions without data loss, supporting dry-run previews, and graceful fallback when renaming across filesystem mount boundaries.

## Prerequisites

Complete all Module 07 lessons (01–21) on CLI flags, file reading/writing, `os` and `path/filepath` APIs, directory walking, and error inspection.

## Available concepts

This project draws on `flag`, `os`, `path/filepath`, `io`, `io/fs`, directory walking, file collision algorithms, and `context.Context` cancellation.

## Project overview

The `_starter/` directory provides type signatures, mode constants, and CLI flag bindings. Your job is to implement directory walking, semantic categorisation, collision resolution, and file moving so that:

- Files are organized either by extension category (`images/`, `documents/`, `audio/`, `videos/`, `archives/`) or by modification date (`YYYY-MM`).
- Existing files with identical names are handled via user-specified collision policy (`rename`, `skip`, or `overwrite`).
- `--dry-run` simulates moves without modifying disk state.
- Cross-device file moves fall back to stream copying (`io.Copy`) and source removal when `os.Rename` fails.

## CLI Flags

| Flag | Default | Description |
|---|---|---|
| `-src` | `.` | Source directory to scan |
| `-target` | `(same as src)` | Destination directory for organized folders |
| `-mode` | `extension` | Organization strategy: `extension` or `date` |
| `-policy` | `rename` | Collision strategy: `rename`, `skip`, or `overwrite` |
| `-dry-run`| `false` | Preview operations without moving files |

## Tasks

### Task 1 — Implement extension categorization

In `_starter/main.go`:
- Implement `CategoryFromExtension(ext string) string`: map `.png`, `.jpg` to `images`, `.pdf`, `.txt`, `.md` to `documents`, `.zip`, `.tar` to `archives`, etc.

### Task 2 — Implement collision resolution

- Implement `resolveCollisionPath(dir, filename string) string`:
  - When a target file exists, iteratively append `_1`, `_2`, ... before the extension until an unused filepath is found.

### Task 3 — Implement cross-device file move

- Implement `moveFile(src, dst string) error`:
  - Attempt `os.Rename(src, dst)`.
  - If `os.Rename` fails (e.g. cross-device link error), fall back to `io.Copy(dst, src)` followed by `os.Remove(src)`.

### Task 4 — Implement directory traversal

- Implement `Organize(ctx context.Context, cfg OrganizerConfig) (OrganizeReport, error)`:
  - Walk `SourceDir` using `filepath.WalkDir`.
  - Skip subdirectories and hidden files (`.` prefix).
  - Compute target path, handle collisions, apply `DryRun` check, and execute `moveFile`.

## Verification

Run all automated unit tests:

```bash
go test -v ./curriculum/modules/07-cli-files-json-config/projects/file-organizer/_starter
```

Expected output:
```
=== RUN   TestOrganize_ExtensionMode
--- PASS: TestOrganize_ExtensionMode (0.01s)
=== RUN   TestOrganize_CollisionRename
--- PASS: TestOrganize_CollisionRename (0.02s)
=== RUN   TestOrganize_CollisionSkip
--- PASS: TestOrganize_CollisionSkip (0.00s)
=== RUN   TestOrganize_DryRun
--- PASS: TestOrganize_DryRun (0.00s)
=== RUN   TestOrganize_InvalidDirectory
--- PASS: TestOrganize_InvalidDirectory (0.00s)
PASS
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Filesystem Safety** | Cross-device moves handled; collisions resolved without silent overwrites; dry-run mode works | Rename works but fails across filesystem boundaries | Unhandled collisions overwrite user data silently |
| **Traversal & Filtering** | `filepath.WalkDir` used efficiently; skips hidden files and deep recursions into target folders | Directory walk works but traverses recursively into organized folders | Crashes or infinite loops on nested subdirectories |
| **CLI & Configuration** | CLI flags properly mapped; clean defaults; informative execution summary | CLI works with minor flag parsing limitations | Flags unparsed or panic on invalid input |
| **Error Handling** | Detailed error collection per file; invalid source directories caught immediately | Errors reported but aborts on single file error | Unchecked I/O errors |

## NEXT UP

[Log Search](../log-search/README.md)
