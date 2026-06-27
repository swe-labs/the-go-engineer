# Files, bytes, directories, and paths

## Learning objective

Understand and apply Files, bytes, directories, and paths in the context of professional Go software engineering.

## Why this matters

Files are the fundamental storage abstraction in every OS. Without understanding how paths resolve, how bytes relate to content, and how directories organize data, every file operation feels fragile and unpredictable.

## Mental model

Go compiles to a native binary with no VM. The compiler type-checks, inlines, escapes-analyzes, and optimises before producing a static binary. This gives Go fast startup and simple deployment.

## Core idea

Without a filesystem, every program would need to manage raw disk sectors — an impossible burden. Files, directories, and paths give developers a simple tree-based abstraction that works across programming languages and operating systems.

## Under the hood

On Linux, ext4 stores files in blocks (default 4 KB). Each file's inode (index node) stores 15 block pointers: 12 direct, 1 indirect, 1 doubly indirect, 1 triply indirect — supporting files up to 16 TB. NTFS uses MFT (Master File Table) entries instead of inodes. All filesystems use B-trees or similar structures for directory lookups.

## How Go uses it

Go's `os` and `path/filepath` packages provide cross-platform filesystem access. `os.Open`, `os.ReadFile`, `os.WriteFile`, and `filepath.Join` abstract over OS differences. Go's `embed` package even allows embedding files directly into the binary at compile time.

## Go example

The example defines a `FileInfo` struct to hold path, size, directory flag, and content. Two functions — `CreateFile` and `ReadFile` — demonstrate writing bytes to disk and reading them back. `main()` creates a file, prints its metadata and content, then cleans up. Path manipulation functions like `filepath.Join`, `filepath.Abs`, `filepath.Base`, `filepath.Dir`, and `filepath.Ext` are used to inspect the path.

## Step-by-step execution

1. The OS mounts a filesystem (NTFS, ext4, APFS) that organizes data as a tree of directories and files.
2. A file is a named sequence of bytes stored in fixed-size blocks on disk, tracked by an inode containing metadata (size, permissions, timestamps).
3. A directory is a special file that maps names to inode numbers — `ls`/`dir` reads this mapping.
4. A path is a string of directory names separated by separators (`/` or `\`) leading to a target file or directory.
5. When a program opens a path, the OS filesystem driver traverses the tree, resolving each component until it finds the target inode.
6. `CreateFile` uses `os.MkdirAll` to create parent directories, then `os.WriteFile` to write bytes.
7. `ReadFile` uses `os.ReadFile` and `os.Stat` to retrieve content and metadata.
8. Path functions decompose the path into its absolute form, base name, directory, and extension.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Confusing a file's name with its contents | Renaming `.txt` to `.pdf` does not make it a PDF | The extension is just part of the name; content determines the format |
| Assuming paths like `C:\Users\name\file.txt` work on every OS | Linux uses forward slashes and no drive letters | Use `filepath.Join` for cross-platform path construction |
| Thinking directories store files hierarchically | On disk, directories are special files mapping names to inode numbers | Understand that the tree structure is a human-friendly abstraction over flat inode tables |

## Debugging walkthrough

**Scenario: A program fails with "file not found" even though the file exists in the same folder as the executable.**
- **Cause:** The program's working directory is different from the executable's directory; relative paths resolve from the working directory, not the binary location.
- **Diagnosis:** Print `os.Getwd()` and compare it to the expected file location.
- **Resolution:** Use absolute paths or resolve paths relative to the executable using `os.Executable()` or `filepath.Abs()`.

**Scenario: A developer stores sensitive data in a world-readable file and other users can read it.**
- **Cause:** Files have permission bits controlling read/write/execute access by owner, group, and others.
- **Diagnosis:** Run `ls -la` (Unix) or check Properties > Security (Windows) to inspect permissions.
- **Resolution:** Set appropriate permissions (0600 for secrets) and use OS-level access controls.

## Production notes

Every application reads configuration files, writes logs, processes data files, and manages assets. Understanding files, bytes, directories, and paths is essential for building reliable software that works across platforms.

## Performance implications

- Reading a file byte-by-byte is 100-1000x slower than reading in large buffers due to syscall overhead.
- SSDs provide ~500 MB/s sequential reads; HDDs provide ~150 MB/s. Random access is 10-100x slower on HDDs.
- Directory traversal for deeply nested paths adds latency — each component requires a separate directory lookup.

## Practice task

Write a Go program that creates a directory, writes a file with specific bytes, reads it back, and confirms the content matches using `os` and `path/filepath` packages.

## Tests / verification

```bash
go test ./curriculum/modules/01-computers-terminal-git-web/03-files-bytes-directories-and-paths/
```

## Review questions

1. Predict whether a given path is absolute or relative and explain what the OS uses as the base for resolution.
2. Why does a directory listing show file sizes in bytes and what does a "byte" actually represent?
3. How can you determine whether two paths refer to the same file on disk?
4. What does `filepath.Join` do differently from simple string concatenation?
5. Why is reading a file byte-by-byte much slower than reading in buffers?

## NEXT UP

[Lesson 04: Terminal basics](../04-terminal-basics/README.md)
