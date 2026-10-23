# Project Assessment — File Organizer

## Purpose

Evaluate your ability to build a production-grade filesystem CLI tool using `filepath.WalkDir`, cross-device file move semantics, collision resolution, and the `flag` package. This assessment verifies that you can write safe filesystem automation that handles edge cases without data loss.

## Targets

- `project-file-organizer`

## Scope

This assessment covers:

- Directory traversal with `filepath.WalkDir` and `io/fs`
- File extension inspection and category mapping
- `os.Rename` behavior and cross-device link error fallback
- `io.Copy` for stream-based file copying
- Collision resolution algorithms (rename, skip, overwrite)
- CLI flag parsing with the `flag` package
- Dry-run mode for safe previewing
- Error collection per file (no abort on single failure)

## Evidence required

Submit:

1. Your completed `_starter/main.go` with `CategoryFromExtension`, `resolveCollisionPath`, `moveFile`, and `Organize` implemented
2. Output of `go test -v .` showing all tests passing
3. A brief explanation of why `os.Rename` can fail across filesystem mount boundaries and how your `moveFile` fallback handles it

## Time expectation

Plan for 90 to 120 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Filesystem safety** | 30% | Cross-device moves handled via `io.Copy` fallback; collisions resolved per policy without silent overwrites; dry-run mode emits previews without modifying disk | `os.Rename` works but fails across mount boundaries; collisions mostly handled | Unhandled collisions silently overwrite user data |
| **Traversal and filtering** | 25% | `filepath.WalkDir` used efficiently; skips hidden files (`.` prefix) and avoids recursive traversal into already-organized target folders | Walk works but traverses recursively into organized subfolders | Crashes or infinite loops on nested directories |
| **CLI and configuration** | 20% | CLI flags properly mapped with clean defaults; execution summary reports files moved, skipped, and errors | CLI works with minor flag parsing limitations | Flags unparsed or panic on invalid input |
| **Error handling** | 25% | Detailed error collection per file; invalid source directories caught with specific error message; continues processing after individual file errors | Errors reported but processing aborts on first file error | Unchecked I/O errors; missing files crash the tool |

**Passing score**: 75%

## Review questions

1. What error does `os.Rename` return when source and destination are on different filesystems? How does your code detect this specific condition?
2. Walk through your `resolveCollisionPath` function: if `report.txt` already exists in the target directory, what filename does it produce? What about if `report_1.txt` also exists?
3. Why does your organizer skip hidden files (files starting with `.`)? What would happen if it processed `.git/` contents?
4. Explain your dry-run implementation: how do you ensure that dry-run mode computes the same file paths as real mode without calling `os.Rename` or `io.Copy`?
5. What is the difference between `filepath.Walk` and `filepath.WalkDir`? Why does the project use `WalkDir`?

## Passing standard

You pass when all tests pass, cross-device moves work, collisions are handled safely, and you can explain filesystem edge cases mechanically.

## Retake policy

If you score below 75%, review Module 07 lessons on filesystem APIs, `filepath.WalkDir`, and CLI flags, then resubmit.

## NEXT UP

[Project Assessment: Config Parser](../config-parser/README.md)
