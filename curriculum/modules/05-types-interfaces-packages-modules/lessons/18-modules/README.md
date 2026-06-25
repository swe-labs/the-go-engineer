# Modules

## Learning objective

Create and use Go modules with proper module paths, understand versioning semantics, use the `replace` directive for local development, and navigate multi-module repositories.

## Why this matters

Before Go modules (pre-Go 1.11), dependency management was a pain point — `GOPATH` mode required all code to live under a single workspace, and there was no built-in versioning. Modules solved this. Today, every Go project uses modules. Understanding module paths, versioning, and the `replace` directive is essential for building reproducible, maintainable Go software.

## Mental model

A Go module is a self-contained unit of versioned code. Think of it as a library on a shelf:

- The **module path** is the label on the spine: `github.com/gin-gonic/gin`.
- The **version** is the edition number: `v1.9.1`.
- The **`go.mod`** is the card catalog entry listing the module's own dependencies.
- The **`replace` directive** is a sticky note saying "for this book, use my local draft instead of the published edition."

Multiple modules in one repository is like a box set — each book has its own ISBN (module path) but they ship together in one box (the repo).

## Core idea

A module is defined by a `go.mod` file at its root. The module path is the import path prefix for all packages within the module.

```go
module github.com/rasel9t6/the-go-engineer

go 1.25.0
```

- Every package inside this module (e.g., `curriculum/modules/...`) is imported relative to `github.com/rasel9t6/the-go-engineer`.
- The `go` directive specifies the Go language version expected.

### Module versioning

Modules are versioned according to semantic versioning (`vMAJOR.MINOR.PATH`):
- `v0.x.x` — unstable, no compatibility guarantees.
- `v1.x.x` — stable, backward compatible within major version.
- `v2+` — breaking changes, module path includes `/v2`, `/v3`, etc.

Versions are tagged in Git: `v1.0.0`, `v1.2.3`, etc.

### The `replace` directive

`replace` tells the module system to use a different location for a dependency:

```go
replace github.com/old/module => ../local/path
replace github.com/old/module => github.com/fork/module v1.2.3
```

Use cases:
- Local development: test changes to a dependency without publishing.
- Fork replacement: use a patched fork instead of the original.
- Broken upstream: pin to a specific fork that fixes a bug.

### Multi-module repositories

A single Git repository can contain multiple modules, each with its own `go.mod`:

```
monorepo/
  go.mod           → module github.com/myorg/monorepo
  pkg/
    api/
      go.mod       → module github.com/myorg/monorepo/pkg/api
  cmd/
    server/
      main.go      → imports github.com/myorg/monorepo/pkg/api
```

Each module has independent version tags. A tag for `pkg/api` would be `pkg/api/v1.0.0`.

## Under the hood

The Go module system (introduced in Go 1.11, default since Go 1.16) works in two phases:

1. **Resolution**: `go build` reads `go.mod`, resolves each `require` to a specific version (using the module proxy or direct source), downloads if needed, and records the content hash in `go.sum`.
2. **Linking**: The compiler treats each module's packages as a namespace. All packages are compiled together; module boundaries do not affect the final binary size, only the import path resolution.

The `replace` directive is applied during resolution — before any version comparison, the module system substitutes the replacement path. This means `replace` works even for indirect dependencies.

## How Go uses it

- **Standard library**: The standard library is not a module but is always available.
- **`go mod init`**: Creates a new module with `go.mod`.
- **`go mod tidy`**: Adds missing dependencies, removes unused ones.
- **`go get`**: Adds, upgrades, or downgrades dependencies.
- **`go mod download`**: Downloads all dependencies into the module cache.
- **`go mod vendor`**: Copies dependencies into a `vendor/` directory.
- **`go list -m all`**: Lists all module versions in the build.
- **`go mod why`**: Explains why a dependency is needed.

## Go example

```go
// File: go.mod
module github.com/rasel9t6/the-go-engineer/curriculum/modules/05-types-interfaces-packages-modules/lessons/18-modules/example

go 1.25.0

require github.com/google/uuid v1.6.0

replace github.com/google/uuid => ./local-uuid-fork
```

```go
// File: main.go
package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	id := uuid.New()
	fmt.Printf("Generated UUID: %s\n", id)
}
```

## Step-by-step execution

Creating and using a module:

1. `mkdir mymodule && cd mymodule`
2. `go mod init github.com/user/mymodule` — creates `go.mod` with module path and Go version.
3. Write `main.go` importing `github.com/google/uuid`.
4. `go mod tidy` — resolves the import, adds `require github.com/google/uuid v1.6.0` to `go.mod`, downloads the module, creates `go.sum`.
5. `go build ./...` — compiles everything.
6. `go list -m all` — shows all module versions.

Using `replace`:

1. Clone `github.com/google/uuid` to `./local-uuid-fork`.
2. Add `replace github.com/google/uuid => ./local-uuid-fork` to `go.mod`.
3. Build with the local fork.

## Common mistakes

- **Module path does not match repo URL**: The module path should usually match the repository URL where others can find it. Mismatch causes import confusion. Use `go mod init <repo-url>`.
- **Forgetting `go mod tidy`**: After adding imports or removing dependencies, always run `go mod tidy`.
- **Missing `go.sum` file**: Commit `go.sum` to your repository. It ensures reproducible builds across machines. Deleting it makes builds non-reproducible.
- **Using `replace` across modules in a workspace**: In Go 1.18+, use `go.work` instead of `replace` for multi-module local development.
- **Not tagging versions**: Go module versions come from Git tags. If you publish a module and do not tag it, users cannot pin to a stable version. Run `git tag v1.0.0 && git push --tags`.
- **Using `replace` with absolute paths**: Absolute paths break on other machines. Use relative paths (relative to the module root).

## Debugging walkthrough

Symptom: `go build` fails with:

```
module github.com/user/mymodule@latest (v0.0.0-...)
  found (v0.0.0-...), but:
  could not import github.com/user/mymodule/foo (no required module provides package "...")
```

**Root cause**: The module path in `go.mod` does not match the import path. For example, `go.mod` says `module mymodule` but imports use `github.com/user/mymodule`.

**Fix**: Ensure `go.mod` starts with `module github.com/user/mymodule` (the full import path).

Symptom: `go mod tidy` adds an `// indirect` dependency unexpectedly.

**Root cause**: Your code imports a package that imports another package, but you do not use that transitive dependency directly. The `// indirect` comment is normal — it tells readers that this dependency is required but not directly imported by your code.

## Production notes

- **Always commit `go.mod` and `go.sum`**: These files define your reproducible build. Losing them means losing the exact dependency versions.
- **Use `GONOSUMCHECK` and `GONOSUMDB`**: For private modules, set these environment variables to bypass the checksum database. `GOPRIVATE` is a convenient shorthand that sets both.
- **Multi-module repos**: In a monorepo, each module should have a unique version tag prefix. Use `pkg/api/v1.0.0` for module `pkg/api`.
- **`replace` in production**: Avoid `replace` in published `go.mod` files — they break consumers who do not have the same local paths. Use `replace` only for local development, and remove before publishing.
- **Minimum version selection (MVS)**: Go uses MVS to select dependency versions — it picks the minimum version that satisfies all requirements. This ensures reproducibility and avoids "dependency hell."

## Performance implications

- Module resolution is a build-time cost, not a runtime cost.
- The module cache (`$GOPATH/pkg/mod` or `$GOCACHE`) stores downloaded modules. Subsequent builds on the same machine skip downloading.
- The number of dependencies affects build time linearly. Each additional module adds resolution and compilation overhead.
- Using `replace` with a local path does not affect build speed compared to using the cached module.
- CI/CD pipelines benefit from caching the module directory: `go mod download` then `go build`.

## Practice task

Create a new Go module in a temporary directory:

1. `mkdir -p /tmp/gomod-demo && cd /tmp/gomod-demo`
2. `go mod init example.com/gomod-demo`
3. Create `main.go` that imports `github.com/google/uuid` and prints a UUID.
4. Run `go mod tidy`.
5. Add a `replace` directive pointing `github.com/google/uuid` to a local directory you create (it can be empty with its own `go.mod`).
6. Run `go build` and verify the output.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/18-modules
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/18-modules
```

## Review questions

1. What file defines a Go module, and what information does it contain?
2. What is the difference between `require` and `replace` in `go.mod`?
3. How does `go mod tidy` help maintain a healthy dependency set?
4. Why must `go.sum` be committed to version control?
5. How does Go handle versioning for breaking changes (major version bumps)?

## NEXT UP

go.mod — the complete structure of the `go.mod` file and every directive it supports.
