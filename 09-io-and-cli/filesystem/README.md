# Track C: Filesystem

## Mission

This track teaches you how Go works with files, paths, directories, embedded assets, and I/O
composition without turning everything into giant in-memory blobs.

## Track Map

| ID | Type | Surface | Why It Matters | Requires |
| --- | --- | --- | --- | --- |
| `FS.1` | Lesson | [files](./1-files) | Introduces basic file I/O and buffered reading. | entry |
| `FS.2` | Lesson | [paths](./2-paths) | Shows path-safe cross-platform filename handling. | `FS.1` |
| `FS.3` | Lesson | [directories](./3-dir) | Adds directory creation, listing, and traversal. | `FS.1`, `FS.2` |
| `FS.4` | Lesson | [temp files](./4-temp) | Teaches temporary resources and cleanup habits. | `FS.1`, `FS.3` |
| `FS.5` | Lesson | [embed](./5-embed) | Shows how to compile static assets into a Go binary. | `FS.3`, `FS.4` |
| `FS.6` | Lesson | [I/O patterns](./6-io-patterns) | Explains `io.Reader` and `io.Writer` composition. | `FS.1`, `FS.2`, `FS.3` |
| `FS.7` | Exercise | [log search tool](./7-log-search) | Combines traversal, scanning, and filtering in one milestone. | `FS.1`, `FS.2`, `FS.3`, `FS.6` |
| `FS.8` | Stretch Lesson | [fs.FS testing seam](./8-fs-testing-seam) | Adds a testable filesystem abstraction without real-disk dependence. | `FS.5`, `FS.6` |

## Suggested Order

1. Work through `FS.1` to `FS.6` in order.
2. Complete `FS.7` as the main filesystem milestone.
3. Use `FS.8` as a stretch lesson after the milestone.

## Track Milestone

`FS.7` is the current filesystem track milestone.

If you can complete it and explain:

- why `filepath.WalkDir` is a better default than manually recursing path strings
- why `bufio.Scanner` is safer than loading whole files into memory for search tasks
- why `fs.FS` gives you a cleaner testing boundary than hard-coded disk paths

then the filesystem part of Section 09 is doing its job.

## Next Step

After `FS.7`, continue to `FS.8` as a stretch lesson or back to the
[Section 09 overview](../README.md).
