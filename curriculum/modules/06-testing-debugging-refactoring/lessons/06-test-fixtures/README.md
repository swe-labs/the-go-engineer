# Test fixtures

## Learning objective

Organize shared test data using the `testdata` directory, implement setup and teardown patterns, and use `TestMain(m)` for package-level fixture management.

## Why this matters

As test suites grow, many tests need shared data: input files, expected outputs, database schemas, configuration. Duplicating data across tests is wasteful and inconsistent. Fixtures provide a single source of truth for test data. `TestMain` lets you run package-level setup once, not once per test. These patterns are essential for any test suite that deals with files, databases, or external configuration.

## Mental model

A fixture is pre-arranged state that a test needs to run. Think of it as a stage with props: before the play (test) begins, the stagehands (setup) place the props. After the play, the stagehands clear them (teardown). The `testdata` directory is the prop closet — a well-known location where all test data lives.

```
testdata/            ← shared data directory
├── input.txt         ← test input
├── expected.json     ← expected output
└── schema.sql        ← database schema
```

## Core idea

Go treats the `testdata` directory specially: `go test` changes the working directory to the package directory, so tests can reference `testdata/` with relative paths. Any file or directory named `testdata` at the package root is ignored by the build — it is never compiled into the test binary.

Key fixture patterns:

| Pattern | When to use |
|---|---|
| **`testdata/` files** | Input/output files shared across tests |
| **Setup/teardown in `TestMain`** | Package-level setup once (database, server, env) |
| **Setup/teardown in test functions** | Per-test setup with `t.Cleanup` |
| **Helper functions** | Reusable setup logic with `t.Helper()` |

`TestMain(m *testing.M)` is the entry point for a package's tests. If defined, it runs instead of `go test`'s default runner. You call `m.Run()` to execute all tests, and do setup before / teardown after:

```go
func TestMain(m *testing.M) {
    // setup
    code := m.Run()
    // teardown
    os.Exit(code)
}
```

## Under the hood

When `go test` runs a package:

1. It changes the working directory to the package directory.
2. If `TestMain` exists, it is called instead of the default test runner.
3. `m.Run()` discovers and runs all tests in the package.
4. `m.Run()` returns an exit code. The package's `TestMain` should call `os.Exit` with this code.
5. If no `TestMain` exists, `go test` handles discovery and execution directly.

Files in `testdata` are not compiled. They are copied into the test's temporary cache when `go test` runs. You can read them with `os.ReadFile("testdata/input.txt")`.

## How Go uses it

The standard library's own tests use `testdata/` extensively:

- `archive/tar/testdata/` — tar archives for reader/writer tests.
- `compress/gzip/testdata/` — compressed streams.
- `encoding/json/testdata/` — JSON input files.
- `net/http/testdata/` — HTTP response bodies and certificates.

The pattern is consistent: fixtures live in `testdata/`, tests reference them by relative path, and `TestMain` handles expensive shared setup.

## Go example

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// LoadFixture reads a fixture file from the testdata directory.
func LoadFixture(name string) (string, error) {
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	content, err := LoadFixture("greeting.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Print(content)
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

func TestMain(m *testing.M) {
	// Package-level setup: ensure testdata exists
	dir := filepath.Join(".", "testdata")
	if err := os.MkdirAll(dir, 0755); err != nil {
		os.Exit(1)
	}
	// Create test fixture
	os.WriteFile(filepath.Join(dir, "greeting.txt"), []byte("Hello from testdata!\n"), 0644)

	code := m.Run()

	// Package-level teardown: remove testdata
	os.RemoveAll(dir)
	os.Exit(code)
}

func TestLoadFixture(t *testing.T) {
	got, err := LoadFixture("greeting.txt")
	if err != nil {
		t.Fatalf("LoadFixture failed: %v", err)
	}
	want := "Hello from testdata!\n"
	if got != want {
		t.Errorf("LoadFixture = %q; want %q", got, want)
	}
}

func TestLoadFixtureMissing(t *testing.T) {
	_, err := LoadFixture("nonexistent.txt")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestWithTempDirFixture(t *testing.T) {
	// Per-test fixture using t.TempDir
	dir := t.TempDir()
	path := filepath.Join(dir, "data.txt")
	content := "per-test content"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(got) != content {
		t.Errorf("content = %q; want %q", string(got), content)
	}
}
```

## Step-by-step execution

Running `go test -v`:

1. `TestMain` runs first. It creates `testdata/` directory and writes `greeting.txt`.
2. `m.Run()` discovers tests: `TestLoadFixture`, `TestLoadFixtureMissing`, `TestWithTempDirFixture`.
3. `TestLoadFixture` calls `LoadFixture("greeting.txt")`. The function reads `testdata/greeting.txt`, returns the content. Assertion passes.
4. `TestLoadFixtureMissing` calls `LoadFixture("nonexistent.txt")`. File not found, returns error. Pass.
5. `TestWithTempDirFixture` creates a temp dir, writes a file, reads it back. Pass.
6. All tests pass. `m.Run()` returns 0.
7. Teardown in `TestMain` removes the `testdata/` directory.
8. `os.Exit(0)`.

If step 1 failed (cannot create directory), `m.Run()` is never called and tests do not execute.

## Common mistakes

- **Hardcoding paths.** Tests that use absolute paths fail when run on another machine. Always use relative paths from the package directory or `t.TempDir()`.

- **Sharing mutable fixtures.** If a test modifies a shared fixture file, other tests see the modified state. Use immutable fixtures, or copy them per test.

- **TestMain without os.Exit.** If `TestMain` does not call `os.Exit(code)`, the test results are lost and `go test` reports a failure.

- **Putting Go code in testdata.** Files in `testdata` are not compiled. If you need Go source as test data, put it in a `testdata` subdirectory but do not expect it to be built.

- **Assuming testdata is writable.** In CI, the working directory may be read-only. Prefer `t.TempDir()` for writable test fixtures.

## Debugging walkthrough

A test fails with `open testdata/input.txt: The system cannot find the path specified`:

```go
func TestLoadData(t *testing.T) {
    data, err := os.ReadFile("testdata/input.txt")
    // ...
}
```

**Symptom**: Test passes locally but fails in CI.

**Investigation**: The test runner's working directory is not the package directory. This happens when using `go test ./pkg` from a parent directory: `go test` changes to the package directory, but certain build tools or IDEs may not.

**Root cause**: The test depends on the working directory being the package root. `go test` guarantees this, but some custom runners do not.

**Fix**: Use `filepath.Join` with the test's package path, or use `os.DirFS("testdata")`. Alternatively, embed the fixture with `//go:embed`:

```go
//go:embed testdata/input.txt
var testFixture string
```

## Production notes

- **Check in fixture files to version control.** Fixtures are part of the test suite. They should be reviewed alongside code changes.
- **Keep fixtures small.** Large binary fixtures bloat the repository. For large data, generate fixtures in `TestMain` or download them during CI.
- **Use `testdata` for inputs and golden files.** Golden files are expected outputs stored alongside inputs for comparison.
- **Avoid secrets in fixtures.** Never commit real passwords, tokens, or keys in `testdata`. Use environment variables or dummy values instead.
- **Document fixture formats.** If a fixture is complex, add a comment at the top of the test file explaining what it contains and how to regenerate it.

## Performance implications

- Reading fixture files from disk is slower than constructing data in memory. For hot tests, consider embedding fixtures with `//go:embed` or loading them once in `TestMain`.
- `TestMain` setup runs once per package, not once per test. Use it for costly shared setup like starting a database container.
- `t.TempDir()` is fast on most systems but involves a syscall. For tests that do not touch the filesystem, avoid it.

## Practice task

Create a `testdata/` directory with a file `users.json` containing:
```json
[{"name": "Alice", "age": 30}, {"name": "Bob", "age": 25}]
```

Write a function `LoadUsers() ([]User, error)` that reads and parses this file. Write tests for:
- Successful load and parse
- Missing file
- Malformed JSON

Use `TestMain` to create the `testdata` directory and write the fixture. Remove it after tests.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/06-test-fixtures
go test -v ./curriculum/modules/06-testing-debugging-refactoring/lessons/06-test-fixtures
```

## Review questions

1. What is special about the `testdata` directory in Go tests?
2. What is the purpose of `TestMain(m *testing.M)` and what must it always call?
3. How does `go test` handle the working directory when a test function reads a file?
4. What are the disadvantages of sharing mutable fixture state between tests?
5. When would you use `TestMain` setup versus per-test setup with `t.TempDir`?

## NEXT UP

Golden files — comparing test output against stored expected output with `-update` flag support.
