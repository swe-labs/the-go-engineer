# Temp files

## Learning objective

Create temporary files and directories with `os.CreateTemp` and `os.MkdirTemp`, ensure automatic cleanup, and understand security implications of temporary storage.

## Why this matters

Every production program needs temporary storage at some point: downloading a file before processing, writing intermediate results, extracting an archive, or creating lock files. Using fixed paths like `/tmp/mylock` creates security vulnerabilities (symlink attacks, race conditions) and portability issues. Go's temp file functions create unique, collision-resistant names in the OS temp directory.

## Mental model

Think of temp files as short-term scratch space. You ask the OS "give me a safe place to put this temporarily" and it returns a unique path. You use it, then clean it up. The system temp directory is like a public whiteboard — anyone can write to it, so your files must have unpredictable names.

```
os.CreateTemp("", "prefix-*")  -->  /tmp/prefix-123456789
os.MkdirTemp("", "build-*")    -->  /tmp/build-987654321/
```

## Core idea

- `os.CreateTemp(dir, pattern) (*os.File, error)` — creates a temp file in `dir` (or `os.TempDir()` if `dir == ""`). The `pattern` may contain `*` which is replaced by random characters.
- `os.MkdirTemp(dir, pattern) (string, error)` — creates a temp directory with the same naming rules.
- Both functions use a cryptographically random suffix to prevent name collisions.
- Cleanup is your responsibility: `defer os.Remove(f.Name())` and `defer os.RemoveAll(dir)`.
- `os.TempDir()` returns the system temp directory (`/tmp` on Unix, `%TMP%` on Windows).

## Under the hood

`os.CreateTemp` generates a random suffix using `rand.Reader` (cryptographic randomness) and tries to create the file with `O_RDWR | O_CREATE | O_EXCL`. The `O_EXCL` flag ensures the file does not already exist — if it does (astronomically unlikely), Go retries with a new random string. The file is created with mode `0600` (owner read/write only) regardless of the umask.

The `*` in the pattern is replaced with the random string. If the pattern has no `*`, the random string is appended to the end. For example, `"log-*.txt"` → `"log-a1b2c3.txt"`, `"build"` → `"build123456789"`.

## How Go uses it

- `os.CreateTemp("", "downloader-*.part")` — temp file for partial downloads.
- `os.MkdirTemp("", "test-*")` — test-specific temp directory.
- `defer os.RemoveAll(tmpDir)` — ensure cleanup even on panic.
- `os.TempDir()` — get system temp directory for information.
- Temp files are written and then atomically renamed to their final location for safe writes.

## Go example

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	tmpFile, err := os.CreateTemp("", "example-*.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create temp file error: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	fmt.Printf("Temp file: %s\n", tmpFile.Name())
	tmpFile.Write([]byte("temporary data"))

	tmpDir, err := os.MkdirTemp("", "example-dir-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create temp dir error: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	fmt.Printf("Temp dir: %s\n", tmpDir)
	subFile := filepath.Join(tmpDir, "data.txt")
	os.WriteFile(subFile, []byte("data in temp dir"), 0644)
}
```

## Step-by-step execution

1. `os.CreateTemp("", "example-*.txt")` calls the OS to create a file in the system temp dir with pattern `example-*.txt`.
2. The OS returns a file e.g., `/tmp/example-x7k2m9.txt` and an `*os.File`.
3. `defer os.Remove(tmpFile.Name())` schedules cleanup.
4. `tmpFile.Write([]byte("temporary data"))` writes 14 bytes.
5. `os.MkdirTemp("", "example-dir-*")` creates e.g., `/tmp/example-dir-3p8q1/`.
6. `os.RemoveAll(tmpDir)` at cleanup removes the directory tree.

## Common mistakes

- **Not cleaning up temp files**: Leftover files fill up `/tmp`. Always `defer os.Remove(f.Name())` or `defer os.RemoveAll(dir)`.
- **Using a fixed path in `/tmp`**: Race condition and security risk. Another user or process could create a symlink at that path. Always use `os.CreateTemp`.
- **Calling `defer os.Remove(f.Name())` before checking error**: If `CreateTemp` fails, `f.Name()` panics on nil. Check error first, then defer.
- **Assuming `os.CreateTemp` file is world-readable**: The file is created with `0600` — only the owner can read it.
- **Closing the temp file before remove on Windows**: On Windows, an open file cannot be removed. Close before removing.

## Debugging walkthrough

Buggy code:

```go
f, err := os.CreateTemp("", "data-*")
defer os.Remove(f.Name()) // panic if CreateTemp fails!
// use f...
```

**Symptom**: Panic with "nil pointer dereference" when CreateTemp fails.

**Investigation**: The error is from disk full or permission denied; `f` is nil.

**Fix**: Check `err` before the defer:

```go
f, err := os.CreateTemp("", "data-*")
if err != nil {
	panic(err) // or handle gracefully
}
defer os.Remove(f.Name())
```

## Production notes

- Always clean up temp files with `defer`. A crashed process leaks temp files — consider a startup cleanup routine.
- Use `os.CreateTemp` in the same directory as the final file for atomic renames (the rename syscall works within the same filesystem).
- Set `$TMPDIR` environment variable on Unix to control where temp files go (CI systems often set this).
- For large temp files (>1GB), consider writing to a dedicated temp partition or using `/var/tmp` for persistent temp files.
- Temp files are not a substitute for proper cache management. Use `os.CreateTemp` for truly temporary data only.

## Performance implications

- Creating a temp file is one `open` syscall with `O_EXCL` — very fast (~1µs).
- The random suffix generation uses crypto/rand which is slower than math/rand but security-critical.
- Writing to a temp file on an SSD is fast; on a tmpfs (RAM-backed /tmp), it's extremely fast.
- Temp directories in tests should use `t.TempDir()` (Go 1.15+) which auto-cleans at test end.

## Practice task

Write a program that simulates a safe file save: write content to a temp file in the same directory as the target, then rename the temp file to the target name. This prevents readers from seeing a half-written file. Save as `main.go`.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/10-temp-files
go test ./curriculum/modules/07-cli-files-json-config/lessons/10-temp-files
```

## Review questions

1. What does the `*` in `os.CreateTemp("", "prefix-*.txt")` do?
2. Why should you use `os.CreateTemp` instead of writing to a hard-coded path like `/tmp/data.txt`?
3. What permissions are set on a file created by `os.CreateTemp`?
4. How do you ensure a temp file is cleaned up even if the function panics?
5. Why might you want to create a temp file in the same directory as the final output file?

## NEXT UP

`fs.FS` — the virtual filesystem interface for portable, testable file access.
