# Go Module Root Basics

## Learning objective

By the end of this lesson, you will understand what a Go module is, how `go.mod` defines the module root, and how Go uses the module path to resolve imports.

## Why this matters

Every Go project starts with `go mod init`. The `go.mod` file defines your module's identity, its dependencies, and the Go version it targets. Without understanding modules, you cannot add dependencies, build reproducible binaries, or organize code across multiple packages.

## Mental model

A Go module is a collection of packages with a single `go.mod` file at the root. The module path in `go.mod` is the import path prefix for all packages in the module. `go mod tidy` ensures the `go.mod` file matches the source code exactly.

## Core idea

The `go.mod` file defines the module path, the Go version, and the dependencies required by the module. Every directory with `.go` files inside the module is a package, and its import path is `module-path/directory-path`.

## Under the hood

The `go.mod` file is a line-oriented text file with directives:

- `module example.com/myapp` — the module's import path prefix
- `go 1.22` — the Go version used for compilation
- `require example.com/lib v1.2.3` — a dependency and its version
- `replace example.com/lib => ../local-lib` — override a dependency's location
- `exclude example.com/lib v1.2.4` — prevent a specific version from being used

The `go.sum` file contains cryptographic hashes of each dependency version, ensuring reproducible builds.

## How Go uses it

Key commands:

- `go mod init example.com/myapp` — creates a new module
- `go mod tidy` — adds missing dependencies and removes unused ones
- `go mod download` — downloads all dependencies to the module cache
- `go mod verify` — checks that cached dependencies match `go.sum`
- `go list -m all` — lists all dependencies with their versions

Import path convention:

```
module github.com/user/repo
├── .                   → github.com/user/repo
├── cmd/server/         → github.com/user/repo/cmd/server
├── pkg/api/            → github.com/user/repo/pkg/api
└── internal/db/        → github.com/user/repo/internal/db
```

## Go example

```go
package main

import "fmt"

func main() {
	fmt.Println("This package is part of a Go module.")
	fmt.Println("The module root has a go.mod file.")
}
```

## Step-by-step execution

1. Look at the root `go.mod` file in this curriculum repository — it defines the module path and Go version.
2. Run `go list -m all` from the module root to see all dependencies.
3. Create a new directory and run `go mod init example.com/test-module` to create a new module.
4. Add a dependency in your test module: `go get rsc.io/quote`
5. Run `go mod tidy` and examine the updated `go.mod` and `go.sum` files.

## Common mistakes

- **Forgetting `go mod tidy` after adding imports**: The code imports a package but `go.mod` does not list it. Run `go mod tidy` to synchronize.
- **Using `go get` for everything**: `go get` is for getting a specific dependency. `go mod tidy` is the standard way to update `go.mod` after editing imports in source code.
- **Committing without `go mod tidy`**: CI will fail when `go mod tidy` changes `go.mod` or `go.sum`. Run it before every commit.
- **Editing `go.mod` by hand for complex changes**: Use `go get`, `go mod tidy`, and `go mod edit` commands instead of manual edits.

## Debugging walkthrough

**Scenario**: `go run .` fails with `missing go.sum entry for module providing package rsc.io/quote`.

1. The package is imported in your code but not recorded in `go.sum`.
2. Run `go mod tidy` to add the missing entry.
3. `go run .` should now succeed.

**Scenario**: `go mod tidy` removes a dependency you thought you needed.

1. The dependency is imported in a file that is conditionally compiled (build tags) or only imported in test files.
2. Check if the import uses a build tag. If so, run `go mod tidy -tags=yourtag` to include it.
3. If the dependency is needed for tests only, ensure it is imported in a `_test.go` file.

## Production notes

- Commit both `go.mod` and `go.sum`. The `go.sum` file ensures reproducible builds.
- Pin your Go version in `go.mod`: the `go 1.22` directive tells Go which language version to use.
- Use `replace` directives for local development but remove them before merging.
- In CI, run `go mod verify` to ensure no dependency tampering.

## Performance implications

- `go mod tidy` parses all Go files in the module. For large modules, this takes a few seconds.
- The module cache (in `$GOPATH/pkg/mod`) stores downloaded dependencies. Once cached, builds are offline.
- A module with many dependencies takes longer to build on first CI run. Use Docker layer caching for faster builds.

## Practice task

1. Look at the `go.mod` file in the module root. Identify the module path, Go version, and at least one dependency.
2. Create a new temporary directory, run `go mod init example.com/myapp`, create a `main.go` that imports `rsc.io/quote`, run `go mod tidy`, and run the program.
3. Examine the generated `go.mod` and `go.sum` files.

## Tests / verification

```bash
go list -m all | head -10
```

Expected output: shows the current module and its top dependencies.

```bash
go run .
```

Expected output:
```
This package is part of a Go module.
The module root has a go.mod file.
```

## Review questions

1. What file defines the root of a Go module?
2. What does `go mod tidy` do?
3. What is the purpose of `go.sum`?
4. How does Go determine the import path for a package inside a module?
5. What happens if you import a package without running `go mod tidy`?

## NEXT UP

[Tooling checklist](../14-tooling-checklist/README.md) — verify your entire Go development environment.
