# Directories

## Learning objective

Create, remove, and list directories using `os.Mkdir`, `os.MkdirAll`, `os.Remove`, `os.RemoveAll`, and `os.ReadDir`, and traverse directory trees with `filepath.WalkDir`.

## Why this matters

Programs that manage files must manage directories too. Creating build output folders, cleaning temp caches, organizing logs by date, and recursively scanning project trees are daily tasks in production software. These operations look simple but have sharp edges: non-empty directories cannot be removed with `os.Remove`, and `os.Mkdir` fails if the parent does not exist.

## Mental model

Directories are containers that map names to files or subdirectories. The filesystem is a tree rooted at `/` (or a drive letter on Windows). Operations on directories are either structural (create, remove) or navigational (list, walk).

```
        /root
       /     \
    docs/    src/
    /          \
README.md     main.go
```

## Core idea

| Operation | Function | Behavior |
|---|---|---|
| Create one dir | `os.Mkdir(path, perm)` | Fails if parent doesn't exist or path exists |
| Create all parents | `os.MkdirAll(path, perm)` | Creates parents as needed; no error if exists |
| Remove empty dir | `os.Remove(path)` | Fails if dir is not empty |
| Remove dir tree | `os.RemoveAll(path)` | Removes everything recursively |
| List entries | `os.ReadDir(path)` | Returns `[]os.DirEntry` (names + types) |
| Walk tree | `filepath.WalkDir(root, fn)` | Calls fn for every file/dir recursively |

## Under the hood

`os.Mkdir` calls the `mkdir` system call on Unix or `CreateDirectory` on Windows. `os.MkdirAll` is a loop that walks the path prefix, calling `os.Mkdir` for each missing component. `os.RemoveAll` recursively deletes directory entries using `ReadDir` and `Remove` — it's depth-first, deleting children before parents.

`os.ReadDir` returns `[]os.DirEntry`, an interface with `Name()`, `IsDir()`, `Type()`, and `Info()`. It is more efficient than the older `ioutil.ReadDir` because it does not call `Stat` on each entry unless you call `Info()`.

## How Go uses it

- `os.Mkdir("uploads", 0755)` — create one directory.
- `os.MkdirAll("backups/2024/01/15", 0755)` — create nested directories, no-op if they exist.
- `os.Remove("tempdir")` — remove empty directory, returns error if not empty.
- `os.RemoveAll("build")` — remove entire build tree.
- `os.ReadDir(".")` — list current directory.
- `filepath.WalkDir("root", walkFn)` — traverse recursively with `WalkDir` (Go 1.16+).
- `os.Chdir(path)` — change working directory (use sparingly).
- `os.Getwd()` — get current working directory.

## Go example

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	os.Mkdir("mydir", 0755)
	os.MkdirAll("a/b/c/d", 0755)

	entries, _ := os.ReadDir(".")
	for _, e := range entries {
		info, _ := e.Info()
		fmt.Printf("  %s %s (%d bytes)\n", e.Type(), e.Name(), info.Size())
	}

	os.RemoveAll("a")
	os.Remove("mydir")

	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil { return err }
		fmt.Println(path)
		return nil
	})
}
```

## Step-by-step execution

1. `os.Mkdir("mydir", 0755)` creates `/cwd/mydir` with permissions `drwxr-xr-x`.
2. `os.MkdirAll("a/b/c/d", 0755)` creates four levels: `a/`, `a/b/`, `a/b/c/`, `a/b/c/d/`.
3. `os.ReadDir(".")` returns all files and dirs in the current directory.
4. `e.Info()` returns `FileInfo` for each entry; `info.Size()` gives the byte count.
5. `os.RemoveAll("a")` deletes the entire `a/` tree.
6. `os.Remove("mydir")` deletes `mydir/` (must be empty).
7. `filepath.WalkDir` traverses remaining files recursively.

## Common mistakes

- **Calling `os.Remove` on a non-empty directory**: Returns error "directory not empty". Use `os.RemoveAll` instead.
- **Calling `os.Mkdir` when the directory already exists**: Returns error. Use `os.MkdirAll` or check `os.IsExist(err)`.
- **Not checking `err` from `os.RemoveAll`**: The function may partially fail, leaving some files behind.
- **Using `ioutil.ReadDir` (deprecated)**: Use `os.ReadDir` (Go 1.16+).
- **Creating directories with wrong permissions**: A directory needs execute (`x`) permission to be traversable. `0644` on a directory prevents listing.

## Debugging walkthrough

Buggy code:

```go
os.Mkdir("logs/2024/01", 0755)
logFile, _ := os.Create("logs/2024/01/events.log")
```

**Symptom**: `Mkdir` fails with "no such file or directory" — `logs/2024` doesn't exist.

**Investigation**: Print the specific error: `fmt.Println(err)`. The error says the parent is missing.

**Fix**: Use `os.MkdirAll("logs/2024/01", 0755)` which creates intermediate directories.

## Production notes

- Always use `os.MkdirAll` for nested directory creation — it's safe to call even if the directory exists.
- For atomic directory creation (ensuring only one goroutine creates it), check `os.IsExist(err)` after `os.Mkdir`.
- Clean up with `os.RemoveAll` in test teardown and temp file management.
- Set restrictive permissions: `0700` for private data, `0755` for public directories.
- Avoid `os.Chdir` in libraries — it changes the global process state. Pass absolute paths instead.

## Performance implications

- `os.Mkdir` is a single syscall. `os.MkdirAll` makes one `mkdir` per missing component.
- `os.ReadDir` reads the directory's inode entries from the filesystem — fast for small directories (microseconds), linear in the number of entries.
- `os.RemoveAll` issues one syscall per file and directory. For large trees, this may take seconds.
- `filepath.WalkDir` is faster than the older `filepath.Walk` because it avoids a `Stat` call per entry.

## Practice task

Write a program that creates a directory structure: `projects/myapp/{src,bin,docs}` with a `README.md` file in the project root. Then walk the structure and print every file. Clean up afterward. Save as `main.go`.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/09-directories
go test ./curriculum/modules/07-cli-files-json-config/lessons/09-directories
```

## Review questions

1. What is the difference between `os.Mkdir` and `os.MkdirAll`?
2. Why does `os.Remove` fail on a non-empty directory, and what should you use instead?
3. What does `os.ReadDir` return, and how is it different from the older `ioutil.ReadDir`?
4. What directory permission is needed to allow entering (but not listing) a directory?
5. How does `filepath.WalkDir` differ from `filepath.Walk` in terms of performance?

## NEXT UP

Temp files — creating temporary files and directories safely.
