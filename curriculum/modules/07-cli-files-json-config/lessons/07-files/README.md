# Files

## Learning objective

Open, create, read, and write files using the `os` package, handle file permissions and modes, and properly close files with `defer`.

## Why this matters

Files are the most persistent storage mechanism in any operating system. Configuration files, log files, data exports, cached assets — every production Go program interacts with the filesystem. Knowing the difference between `os.ReadFile`, `os.Create`, and `os.OpenFile` lets you choose the right tool for each task.

## Mental model

Think of the filesystem as a key-value store where the key is a path (string) and the value is a sequence of bytes. Read operations retrieve bytes; write operations store bytes. File descriptors are handles to open files — like a cursor that tracks your position in the byte sequence.

```
path --> [file bytes]     os.ReadFile(path)  --> []byte
       [file bytes]       os.WriteFile(path, data, perm)
```

## Core idea

Three main patterns for file access:

1. **Whole file at once** — `os.ReadFile(path)` and `os.WriteFile(path, data, perm)`.
2. **Open + stream** — `os.Open(path)` returns `*os.File` for reading; `os.Create(path)` for writing.
3. **Open with flags** — `os.OpenFile(path, mode, perm)` for append, read-write, exclusive create, etc.

Every opened file must be closed. Use `defer f.Close()` right after opening.

## Under the hood

The operating system manages file descriptors (FDs). Each process has a limited number (typically 256–1024 per process, configurable). An `*os.File` wraps an FD. When you call `os.Open`, the kernel translates the path to an inode, checks permissions, and returns an FD. `Read` and `Write` move the file offset. `Close` releases the FD back to the OS.

`os.ReadFile` reads the entire file into memory — fine for config files and small payloads, dangerous for multi-GB files. `os.WriteFile` truncates the file before writing — it's an atomic replacement from the user's perspective.

## How Go uses it

- `os.ReadFile(path) ([]byte, error)` — reads entire file, suitable for files under ~100MB.
- `os.WriteFile(path, data, perm)` — writes data, creating or truncating the file.
- `os.Open(path) (*os.File, error)` — opens existing file for reading only.
- `os.Create(path) (*os.File, error)` — creates or truncates file for read-write.
- `os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)` — opens with custom flags.
- `f.Stat() (FileInfo, error)` — returns file metadata (size, mode, mod time).
- `f.Close() error` — releases the file descriptor.

## Go example

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	data := []byte("Hello, file!\nLine 2\n")
	err := os.WriteFile("test_output.txt", data, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		os.Exit(1)
	}

	content, err := os.ReadFile("test_output.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Read %d bytes:\n%s", len(content), content)

	f, _ := os.Open("test_output.txt")
	defer f.Close()
	info, _ := f.Stat()
	fmt.Printf("Size: %d bytes\n", info.Size())

	os.Remove("test_output.txt")
}
```

## Step-by-step execution

1. `os.WriteFile("test_output.txt", []byte("Hello..."), 0644)` creates the file with read/write for owner, read for others.
2. `os.ReadFile("test_output.txt")` reads the entire content into memory.
3. `os.Open("test_output.txt")` returns a file handle; `f.Stat()` reads metadata.
4. `info.Size()` returns the byte count.
5. `f.Close()` releases the FD (deferred).
6. `os.Remove("test_output.txt")` deletes the file.

## Common mistakes

- **Forgetting to close files**: Leaks file descriptors. Under heavy load, you hit the OS limit and all file operations fail.
- **Not checking write errors**: Partial writes happen (disk full, quota exceeded). Check the error from `WriteFile` and `Write`.
- **Using `os.ReadFile` for huge files**: Reads the entire file into memory. Use `os.Open` + `bufio.Scanner` or `io.Copy` for streaming.
- **Wrong permission bits**: `0644` (octal) is owner read-write, others read. Use the `os.FileMode` type, not a decimal number.
- **Skipping `defer` placement**: `defer f.Close()` should come immediately after `os.Open`, not 20 lines later where you may forget.

## Debugging walkthrough

Buggy code:

```go
f, err := os.Open("config.json")
if err != nil {
	panic(err)
}
// use f...
f.Close()

// 50 lines later...
g, err := os.Open("data.txt")  // May fail: FD leak from f.Close() error ignored
```

**Symptom**: Random "too many open files" errors under load.

**Investigation**: Count FDs with `lsof -p PID` on Linux. The count grows.

**Fix**: Always `defer f.Close()` immediately after open, regardless of error handling path.

## Production notes

- Use `os.ReadFile` for config files (< 1MB). Use `os.Open` + streaming for data files.
- File permissions are modified by `umask` on Unix. `os.WriteFile(path, data, 0644)` may produce `0640` if umask is `0024`.
- For atomic file writes: write to a temp file in the same directory, then `os.Rename`. This prevents readers from seeing a half-written file.
- On Windows, file locking semantics differ. Files opened for writing may be locked by the OS.
- Use `path/filepath` for path manipulation, not string concatenation.

## Performance implications

- `os.ReadFile` allocates a single `[]byte` — efficient for small files.
- Opening a file is a syscall (~1µs on modern hardware). Reuse open file handles for repeated reads/writes.
- `os.WriteFile` truncates and writes in one call — generally faster than open/write/close.
- Sequential reads are heavily optimized by the OS page cache; random reads within a large file are slower.
- SSD vs HDD: SSDs handle random I/O well, but syscall overhead is the same.

## Practice task

Write a program that copies a file from `src.txt` to `dst.txt` using `os.ReadFile` and `os.WriteFile`. Then implement the same with `os.Open`, `os.Create`, and `io.Copy`. Compare the two approaches. Save as `main.go`.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/07-files
go test ./curriculum/modules/07-cli-files-json-config/lessons/07-files
```

## Review questions

1. What is the difference between `os.Open` and `os.OpenFile`?
2. Why should you use `defer f.Close()` immediately after `os.Open`?
3. What happens to the existing content when you call `os.WriteFile` on an existing file?
4. What does `f.Stat()` return, and how do you get the file size from it?
5. What is the octal permission value for owner read-write, group read, others read?

## NEXT UP

Paths — manipulating file paths portably with `path/filepath`.
