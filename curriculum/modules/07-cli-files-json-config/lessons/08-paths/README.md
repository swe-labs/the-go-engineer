# Paths

## Learning objective

Construct, decompose, clean, and traverse filesystem paths portably across Windows and Unix using the `path/filepath` package.

## Why this matters

Hard-coding path separators (`/` vs `\`) is the #1 cross-platform bug in Go programs. The `path/filepath` package uses the correct separator for the OS automatically. Whether you're building a CLI tool for Linux or a desktop app for Windows, `filepath.Join` and its siblings make your code portable without `// +build` tags.

## Mental model

A path is a sequence of directory names and a filename separated by `os.PathSeparator` (`/` on Unix, `\` on Windows). The `filepath` package treats paths as abstract sequences — you join, split, and clean them without worrying about the separator.

```
"docs/2024/report.txt"  -->  Split --> dir="docs/2024/", file="report.txt"
"a/b/../c/./d"         -->  Clean --> "a/c/d"
"a", "b", "c"          -->  Join  --> "a/b/c"
```

## Core idea

Seven essential `filepath` functions:

| Function | Purpose | Example |
|---|---|---|
| `Join(elem ...string)` | Build paths with correct separator | `Join("a", "b")` → `"a/b"` |
| `Clean(path)` | Remove `..`, `.`, double separators | `Clean("a/../b")` → `"b"` |
| `Split(path)` | Split into `(dir, file)` | `Split("a/b.txt")` → `"a/", "b.txt"` |
| `Dir(path)` | Return the directory portion | `Dir("a/b.txt")` → `"a"` |
| `Base(path)` | Return the last element | `Base("a/b.txt")` → `"b.txt"` |
| `Ext(path)` | Return the extension | `Ext("a/b.txt")` → `".txt"` |
| `IsAbs(path)` | Check if path is absolute | `IsAbs("/a")` → `true` |

## Under the hood

Go's `filepath` package detects the OS at compile time using build tags. On Unix, it uses `/` as separator and treats paths as UTF-8 byte sequences. On Windows, it handles drive letters (`C:`), backslashes, and UNC paths (`\\server\share`). The `Clean` function normalizes separators, removes trailing slashes (except root), and resolves `..` without accessing the filesystem — it is purely lexical.

`filepath.Walk` uses `os.Lstat` internally to read directory entries, recursing depth-first. It calls the provided `WalkFunc` for every file and directory.

## How Go uses it

- `filepath.Join` for building paths — never use string concatenation with `+ "/" +`.
- `filepath.Clean` to normalize user-provided paths before using them.
- `filepath.Dir` and `filepath.Base` to decompose paths without string manipulation.
- `filepath.Ext` to get file extensions for content-type detection.
- `filepath.Walk` for recursive directory processing.
- `filepath.Match` and `filepath.Glob` for pattern matching (e.g., `*.txt`).
- `filepath.Rel` to compute relative paths between two absolute paths.

## Go example

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	base := "docs"
	sub := "2024"
	file := "report.txt"
	fullPath := filepath.Join(base, sub, file)
	fmt.Println("Joined path:", fullPath)

	messy := "docs/../docs//2024/./report.txt"
	clean := filepath.Clean(messy)
	fmt.Println("Cleaned path:", clean)

	dir, name := filepath.Split(fullPath)
	fmt.Printf("Dir: %q, File: %q\n", dir, name)

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		fmt.Printf("  %s (%d bytes)\n", path, info.Size())
		return nil
	})
}
```

## Step-by-step execution

For `filepath.Join("docs", "2024", "report.txt")` on Unix:

1. First element `"docs"` starts the path.
2. Append separator `/` then `"2024"` → `"docs/2024"`.
3. Append separator `/` then `"report.txt"` → `"docs/2024/report.txt"`.

For `filepath.Clean("docs/../docs//2024/./report.txt")`:

1. Replace double `/` with single → `"docs/../docs/2024/./report.txt"`.
2. Resolve `..` with preceding `docs` → `"docs/2024/./report.txt"`.
3. Remove `.` → `"docs/2024/report.txt"`.

## Common mistakes

- **Using `+` for path construction**: `"dir" + "/" + "file"` breaks on Windows.
- **Assuming `/` is always root**: On Windows, `/` is relative to the current drive's root.
- **Forgetting that `filepath.Clean` is lexical**: `Clean("a/../../b")` → `"../b"`, it does not check if paths exist.
- **Modifying the path inside `filepath.Walk` walk function**: The walk function receives the path; do not assume you can `os.Chdir` or modify the walked directory during traversal.
- **Passing user input directly to `filepath.Join` without cleaning**: `Join` does not resolve `..`, so `Join("/safe", userInput)` may escape `/safe`.

## Debugging walkthrough

Buggy code:

```go
path := "/data/" + filename
f, err := os.Open(path)
```

**Symptom**: On Windows, `"/data/" + "config.json"` becomes `"/data/config.json"` which is not a valid absolute path.

**Investigation**: Print `path` — on Windows it's interpreted as relative to `C:`, not root.

**Fix**: Use `filepath.Join("\\data", filename)` or better, use platform-independent path logic.

## Production notes

- Always use `filepath.Join` to build paths, never string concatenation.
- For configurable base directories, use `filepath.Abs` to resolve relative paths to absolute before storing.
- Validate user-provided paths: reject paths containing `..` if you want to restrict to a directory tree.
- `filepath.Walk` walks in lexical order, not creation order. Use `filepath.WalkDir` (Go 1.16+) for better performance.
- For temporary paths, use `os.TempDir()` and `filepath.Join` — never hard-code `/tmp`.

## Performance implications

- `filepath` functions operate on strings without OS calls (except `Abs`, `EvalSymlinks`, and `Walk`).
- `Clean`, `Join`, `Split` are O(n) in path length — negligible cost.
- `Walk` issues a `stat` syscall per file/directory. For large trees (100K+ files), this is significant.
- `filepath.WalkDir` (Go 1.16+) uses `os.ReadDir` internally, which is faster than `os.Lstat` for directory listing.

## Practice task

Write a program that accepts a directory path as a CLI arg, walks it recursively, and prints every `.go` file with its line count (count `\n` in the file). Use `filepath.Walk` and `os.ReadFile`. Save as `main.go`.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/08-paths
go test ./curriculum/modules/07-cli-files-json-config/lessons/08-paths
```

## Review questions

1. What does `filepath.Join("a", "b", "..", "c")` return?
2. Why is `filepath.Ext("archive.tar.gz")` only `".gz"`, not `".tar.gz"`?
3. How does `filepath.Clean` handle a path containing `..`?
4. What is the difference between `filepath.Walk` and `filepath.WalkDir`?
5. On Windows, what does `filepath.IsAbs("/foo")` return?

## NEXT UP

Directories — creating, removing, and listing directories.
