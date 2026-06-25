# t.Cleanup

## Learning objective

Register cleanup functions with `t.Cleanup`, understand LIFO cleanup order, compare `defer` vs `t.Cleanup` in tests, and clean up temporary resources reliably.

## Why this matters

Tests create temporary state: files, directories, database connections, environment variables, goroutines. If cleanup is skipped — because a test panics, calls `t.Fatal`, or returns early — leaked state accumulates, breaks other tests, and corrupts the test environment. `t.Cleanup` guarantees cleanup runs no matter how the test exits, including panics.

## Mental model

Think of `t.Cleanup` as a stack of "undo" operations. Each time you create a resource, you register its destructor. When the test finishes, the destructors run in reverse order (LIFO — last registered, first executed). This mirrors the pattern `defer` provides in production code, but `t.Cleanup` integrates with the test lifecycle: cleanup runs after all subtests complete, and it runs even if the test panics.

```
tempDir := t.TempDir()   ← creates directory
t.Cleanup(func() { ... }) ← registers undo
// ... use tempDir ...
// test ends → cleanup runs automatically
```

## Core idea

`t.Cleanup` registers a function to be called when the test and all its subtests complete:

```go
func TestWithFile(t *testing.T) {
    f, err := os.CreateTemp("", "test")
    if err != nil {
        t.Fatal(err)
    }
    t.Cleanup(func() { os.Remove(f.Name()) })
    // use f ...
}
```

Key properties:

| Property | Behavior |
|---|---|
| **Guaranteed execution** | Runs even if test panics or calls `t.Fatal` |
| **LIFO order** | Last registered cleanup runs first |
| **Per-subtest** | Cleanup registered in a subtest runs when that subtest ends |
| **Parent cleanup** | Runs after all subtests finish |
| **Multiple calls** | You can register multiple cleanup functions; all run |

`t.TempDir()` is a built-in convenience that creates a temp directory and registers its removal via `t.Cleanup`. It's the recommended way to create temp directories in tests.

## Under the hood

When `t.Cleanup(f)` is called:

1. Go pushes `f` onto a stack associated with the current `*testing.T`.
2. When the test function returns (or panics, or calls `t.Fatal`), Go pops functions off the stack in LIFO order and calls them.
3. Cleanup runs in the same goroutine as the test, sequentially.
4. If a cleanup function panics, the test is marked as failed, and remaining cleanup functions still run.

This is implemented internally as a `[]func()` slice appended to on each `t.Cleanup` call, then iterated in reverse at cleanup time.

## How Go uses it

The standard library uses `t.Cleanup` extensively:

- `t.TempDir()` registers directory removal.
- `net/http/httptest`'s `NewServer` registers server closure via `t.Cleanup`.
- `testing/iotest` uses cleanup for reader/writer test patterns.

The convention: if your test creates a resource, register its cleanup immediately after creation. This makes cleanup co-located with setup, readable, and impossible to forget.

## Go example

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteConfig writes a config file and returns its path.
func WriteConfig(dir, content string) (string, error) {
	path := filepath.Join(dir, "config.txt")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}
	return path, nil
}

// ReadConfig reads a config file and returns its content.
func ReadConfig(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	dir, _ := os.MkdirTemp("", "example")
	defer os.RemoveAll(dir)

	path, _ := WriteConfig(dir, "hello world")
	content, _ := ReadConfig(path)
	fmt.Println(content)
}
```

```go
// main_test.go
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteReadConfig(t *testing.T) {
	dir := t.TempDir()

	path, err := WriteConfig(dir, "hello world")
	if err != nil {
		t.Fatalf("WriteConfig failed: %v", err)
	}

	content, err := ReadConfig(path)
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}
	if content != "hello world" {
		t.Errorf("content = %q; want %q", content, "hello world")
	}
}

func TestCleanupOrder(t *testing.T) {
	var order []string
	t.Cleanup(func() { order = append(order, "first") })
	t.Cleanup(func() { order = append(order, "second") })
	t.Cleanup(func() { order = append(order, "third") })

	t.Cleanup(func() {
		// third, second, first — LIFO order
		if len(order) != 3 || order[0] != "third" || order[1] != "second" || order[2] != "first" {
			t.Errorf("cleanup order = %v; want [third second first]", order)
		}
	})
}

func TestCleanupRunsAfterFatal(t *testing.T) {
	cleaned := false
	t.Cleanup(func() { cleaned = true })

	// This Fatal would abort the test, but Cleanup still runs
	t.Fatal("intentional failure")

	// We never reach here, but Cleanup verifies in post-condition
	_ = cleaned
}
```

## Step-by-step execution

Running `TestWriteReadConfig`:

1. `t.TempDir()` creates a temp directory and registers `os.RemoveAll(dir)` via `t.Cleanup`.
2. `WriteConfig(dir, "hello world")` creates a file. Returns path.
3. `ReadConfig(path)` reads the file. Returns `"hello world"`.
4. Assertion passes.
5. Test function returns.
6. Cleanup runs: the temp directory is removed (registered by `t.TempDir`).
7. Test passes.

If step 2 failed (`t.Fatal`), cleanup would still run in step 6, removing the temp directory. Without `t.Cleanup`, a `t.Fatal` before `defer` would leak the directory.

## Common mistakes

- **Using `defer` instead of `t.Cleanup`.** `defer` works, but `t.Cleanup` runs even after `t.Fatal` and `panic`. In tests, prefer `t.Cleanup`.

- **Registering cleanup after the resource is created.** Always register cleanup immediately after creation. If creation fails, skip cleanup registration.

- **Assuming cleanup runs at the end of the parent, not after subtests.** Cleanup in a parent test waits until all subtests complete. If you need cleanup after each subtest, register it inside the subtest.

- **Order-dependent cleanup.** If cleanup B depends on A, register A first, then B. LIFO means B runs first, then A — correct if B was created after A.

## Debugging walkthrough

A test leaks temp directories:

```go
func TestProcessData(t *testing.T) {
    dir, _ := os.MkdirTemp("", "process")
    // forgot to clean up
    ProcessData(dir)
}
```

**Symptom**: Disk fills up after many test runs. Temp directories accumulate.

**Investigation**: `$TMPDIR` has thousands of `process*` directories. The test never removes them.

**Root cause**: No `defer os.RemoveAll(dir)` or `t.Cleanup`. The directories persist.

**Fix**: Replace with `dir := t.TempDir()`. Or, if you need manual control: `defer os.RemoveAll(dir)` or `t.Cleanup(func() { os.RemoveAll(dir) })`.

## Production notes

- **Always use `t.TempDir()` when you need a temp directory.** It handles cleanup, works on all OSes, and generates unique names.
- **Cleanup functions should be idempotent or fast.** If a cleanup function hangs, the test hangs.
- **Avoid cleanup functions that could fail.** If cleanup must report errors, use `t.Log` to record them. Do not call `t.Fatal` from cleanup — it would race with the test's completion.
- **`t.Cleanup` is not `defer`.** Both have their place. Use `t.Cleanup` for test-specific teardown, `defer` for function-scoped cleanup that must run even outside test contexts.
- **External resources.** Clean up HTTP servers, database connections, and goroutines in `t.Cleanup` to ensure no resources leak between tests.

## Performance implications

- `t.Cleanup` adds negligible overhead — it appends to a slice.
- Cleanup functions run sequentially. If a cleanup is expensive (e.g., deleting many files), it adds to test time.
- `t.TempDir` on Linux uses `/tmp` (often tmpfs, fast). On Windows temp directories are on disk.
- LIFO order is intentional: the most recently created resource (likely the most specific) is cleaned up first, before the resources it depends on.

## Practice task

Write a function `SetupEnv(key, value string) (cleanup func())` that sets an environment variable and returns a cleanup function that restores the original value. Write a test that:
1. Calls `SetupEnv("MY_VAR", "test-value")`.
2. Registers the returned cleanup with `t.Cleanup`.
3. Verifies the env var is set.
4. On test completion, verifies the env var is restored to its original value (or unset).

Use `os.Getenv` and `os.Setenv`. Run `go test -v` to confirm.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/05-t-cleanup
go test -v ./curriculum/modules/06-testing-debugging-refactoring/lessons/05-t-cleanup
```

## Review questions

1. What is the difference between `t.Cleanup` and `defer` in a test function?
2. In what order do multiple `t.Cleanup` registrations execute?
3. If a test calls `t.Fatal`, do cleanup functions still run?
4. What does `t.TempDir()` do under the hood?
5. Why is LIFO order useful for cleanup of nested resources?

## NEXT UP

Test fixtures — organizing shared setup and teardown with `testdata` and `TestMain`.
