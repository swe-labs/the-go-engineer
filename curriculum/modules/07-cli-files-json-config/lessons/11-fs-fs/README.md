# fs.FS

## Learning objective

Use the `io/fs` package to read files from virtual filesystems, embed static assets with `embed.FS`, and test filesystem-dependent code with `fstest.MapFS`.

## Why this matters

The `io/fs` package (Go 1.16) introduced a standard interface for filesystems. Any filesystem — real OS files, embedded assets, zip archives, in-memory test data — can be treated uniformly. Your code that reads configuration, templates, or static files no longer depends on the real OS. This makes testing trivial and enables powerful patterns like embedding assets in the binary.

## Mental model

An `fs.FS` is a virtual directory tree. You open it with `fs.Open(name)` and get a file (which can be a directory). The interface is small but composable: `fs.ReadFile`, `fs.ReadDir`, `fs.WalkDir`, and `fs.Glob` all work on any `fs.FS`.

```
                  +-------+
os.DirFS(".") --> | fs.FS | --> files from disk
                  +-------+
embed.FS       --> | fs.FS | --> files baked into binary
                  +-------+
fstest.MapFS   --> | fs.FS | --> files defined in-memory
                  +-------+
```

## Core idea

The `fs.FS` interface:

```go
type FS interface {
	Open(name string) (File, error)
}
```

An `fs.File` is a read-only file handle:

```go
type File interface {
	Stat() (FileInfo, error)
	Read([]byte) (int, error)
	Close() error
}
```

Package-level helpers that work on any `fs.FS`:

- `fs.ReadFile(fsys, name)` — reads entire file.
- `fs.ReadDir(fsys, name)` — lists directory entries.
- `fs.WalkDir(fsys, root, fn)` — walks the tree.
- `fs.Glob(fsys, pattern)` — matches files by glob pattern.
- `fs.Sub(fsys, dir)` — returns a sub-filesystem rooted at `dir`.

## Under the hood

`embed.FS` is a compile-time feature: the `//go:embed` directive reads files at build time and stores their content in the binary. The embedded filesystem is read-only, immutable, and requires no disk I/O at runtime. It is implemented as a compressed (or uncompressed) blob in the binary's data section.

`os.DirFS(path)` wraps a real OS directory as a `fs.FS`. Each `Open` call translates to `os.Open` with the path joined to the root. The resulting `fs.File` is the actual OS file descriptor.

`fstest.MapFS` is a `map[string]*fstest.MapFile` for testing. No disk I/O, no permissions, no cleanup.

## How Go uses it

- `//go:embed static/*` embeds all files in the `static/` directory.
- `os.DirFS(".")` creates a read-only view of the current directory.
- `fs.ReadFile(fsys, "config.json")` reads a config from any FS.
- `embed.FS` is used for web server templates, migrations, and static assets.
- `fstest.TestFS(fsys, "file1", "file2")` validates an FS implementation.
- `fs.WalkDir` replaces `filepath.Walk` for FS-agnostic traversal.

## Go example

```go
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
)

//go:embed static/*
var embeddedFiles embed.FS

func main() {
	data, _ := fs.ReadFile(embeddedFiles, "static/hello.txt")
	fmt.Printf("Embedded file: %s", data)

	fs.WalkDir(embeddedFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil { return err }
		info, _ := d.Info()
		fmt.Printf("  %s (%d bytes)\n", path, info.Size())
		return nil
	})

	diskFS := os.DirFS(".")
	data, _ = fs.ReadFile(diskFS, "main.go")
	fmt.Printf("\nmain.go on disk: %d bytes\n", len(data))
}
```

## Step-by-step execution

1. `//go:embed static/*` embeds everything in `static/` into `embeddedFiles` at compile time.
2. `fs.ReadFile(embeddedFiles, "static/hello.txt")` opens and reads the embedded file — no disk access at runtime.
3. `fs.WalkDir(embeddedFiles, ".")` lists all embedded files recursively.
4. `os.DirFS(".")` creates a virtual FS backed by the real current directory.
5. `fs.ReadFile(diskFS, "main.go")` opens and reads the actual disk file.

## Common mistakes

- **`//go:embed` path starts after the directory**: `//go:embed static/*` makes files available at `static/hello.txt`, not `hello.txt`.
- **Embedding outside the module**: `//go:embed ../foo` is rejected. Embedded paths must be within the module directory.
- **Writing to an embedded FS**: `embed.FS` is read-only. You cannot modify embedded files at runtime.
- **Forgetting `embed` import**: The `embed` package must be imported (use `_ "embed"` if not using any exported name).
- **Using `os.Open` instead of `fs.Open` with `os.DirFS`**: If you write code against `fs.FS`, it works with any FS. Hard-coding `os.Open` ties you to the disk.

## Debugging walkthrough

Buggy code:

```go
//go:embed static/config.json
var config embed.FS

func main() {
	data, err := config.ReadFile("config.json")
	// ...
}
```

**Symptom**: Compile error: `config.ReadFile undefined` (type `embed.FS` has no `ReadFile` method).

**Investigation**: `embed.FS` only implements `fs.FS`, which has `Open`, not `ReadFile`.

**Fix**: Use `fs.ReadFile(config, "static/config.json")` — the package-level helper calls `Open` internally.

## Production notes

- `embed.FS` is ideal for shipping templates, SQL migrations, and static web assets in a single binary.
- Use `os.DirFS` for local development (hot reload) and `embed.FS` for production deployment.
- `fstest.MapFS` makes filesystem tests fast, deterministic, and parallel-safe.
- The `fs` package is read-only by design. For write operations, use `os` package functions directly.
- Third-party FS implementations exist for zip, tar, S3, and FTP — all implementing `fs.FS`.

## Performance implications

- `embed.FS` is the fastest filesystem — no syscalls, no disk I/O. Data is read from memory.
- `os.DirFS` has the same performance as `os.Open` — each `Open` is a syscall.
- `fs.WalkDir` on `embed.FS` is O(n) with very low constant factor.
- `fstest.MapFS` uses a `map[string]*MapFile` — `Open` is O(1) amortized.

## Practice task

Create a `//go:embed static/*` directive with at least two files in a `static/` directory. Write a program that reads and prints all embedded file names and their sizes using `fs.WalkDir`. Then use `fstest.MapFS` to test a function that reads a greeting from a file and returns "Hello, <content>!". Save as `main.go`.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/11-fs-fs
go test ./curriculum/modules/07-cli-files-json-config/lessons/11-fs-fs
```

## Review questions

1. What is the only method required by the `fs.FS` interface?
2. How do you embed files in a Go binary at compile time?
3. What is the difference between `os.DirFS(".")` and using `os.Open` directly?
4. Why is `embed.FS` read-only?
5. What is `fstest.MapFS` used for, and how does it make testing easier?

## NEXT UP

JSON marshal — encoding and decoding JSON data in Go.
