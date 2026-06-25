# go.mod

## Learning objective

Read and write a `go.mod` file, understand every directive (`module`, `go`, `require`, `exclude`, `replace`, `retract`), and use `go mod init`, `go mod tidy`, and `go mod why` effectively.

## Why this matters

The `go.mod` file is the identity card of every Go module. It declares what your module is, what language version it targets, and exactly which versions of every dependency it needs. Misconfigured `go.mod` files cause build failures, unexpected upgrades, and broken imports. Every Go developer reads and edits `go.mod` — either directly or through `go mod` commands — so understanding its structure is essential.

## Mental model

`go.mod` is a recipe card:
- **`module`** — the name of the dish (your module path).
- **`go`** — the kitchen version (Go language version).
- **`require`** — the ingredient list (dependencies and their versions).
- **`exclude`** — ingredients you refuse to use (forbidden versions).
- **`replace`** — ingredient substitutions (use this instead of that).
- **`retract`** — "ignore that version, it was a bad batch."

The recipe is authored by the module maintainer but can be tweaked by consumers via `replace` and `exclude`.

## Core idea

### `go.mod` directives

**`module`**: The module's import path. All packages in the module are imported relative to this path.

```
module github.com/rasel9t6/the-go-engineer
```

**`go`**: The Go language version the module was written for. This sets the language features available and the toolchain version.

```
go 1.25.0
```

**`require`**: Lists dependencies with exact versions.

```
require (
    github.com/google/uuid v1.6.0
    golang.org/x/sync v0.7.0 // indirect
)
```

- `// indirect` marks dependencies not directly imported by the module but needed by a direct dependency.

**`exclude`**: Prevents a specific version from being used, even if a dependency requests it.

```
exclude github.com/broken/module v0.5.0
```

**`replace`**: Substitutes a dependency with a different module path or local directory.

```
replace github.com/old/module => ./local/fork
replace github.com/old/module => github.com/new/module v2.0.0
```

**`retract`** (Go 1.16+): Marks versions as retracted — they should not be used.

```
retract v1.0.0 // Accidental publish
retract [v1.1.0, v1.5.0] // Range of bad versions
```

### `go mod init`

Creates a new `go.mod`:

```bash
go mod init github.com/user/project
```

### `go mod tidy`

The most important `go mod` command:
- Adds missing `require` entries for packages used in imports.
- Removes `require` entries for packages no longer imported.
- Adds `// indirect` comments where appropriate.
- Updates `go.sum`.

### `go mod why`

Explains why a dependency is needed:

```bash
go mod why -m github.com/google/uuid
```

## Under the hood

The `go.mod` file is parsed by the `cmd/go` tool. The `go` directive sets the language version, which controls language features (e.g., `go 1.18` enables generics). The toolchain version (`go 1.25.0` in the `go` directive) indicates which Go toolchain to use.

Module resolution uses **Minimum Version Selection (MVS)**: given a directed acyclic graph of module requirements, Go selects the minimum version of each module that satisfies all `require` statements. This ensures deterministic, reproducible builds.

The `exclude` directive is applied during MVS — excluded versions are removed from consideration before selection begins. The `replace` directive overrides the module path and version entirely.

## How Go uses it

- **`go mod init`**: Creates `go.mod` from scratch.
- **`go build`**: Reads `go.mod` to resolve all imports.
- **`go mod tidy`**: The most commonly used maintenance command.
- **`go get module@version`**: Updates `go.mod` with the specified version.
- **`go list -m all`**: Displays the final resolved versions after MVS.
- **`go mod edit`**: Programmatically edits `go.mod` (used by scripts and tools).
- **`go mod verify`**: Verifies dependencies against `go.sum` hashes.

## Go example

```go
package main

import (
	"fmt"
	"runtime/debug"
)

func main() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("no build info (run with 'go run' or 'go build')")
		return
	}
	fmt.Println("Module path:", bi.Main.Path)
	fmt.Println("Go version:", bi.GoVersion)
	for _, dep := range bi.Deps {
		fmt.Printf("  dep: %s@%s\n", dep.Path, dep.Version)
	}
}
```

A corresponding `go.mod` file:

```
module github.com/example/myapp

go 1.25.0

require (
    github.com/gorilla/mux v1.8.1
    github.com/lib/pq v1.10.9
    golang.org/x/crypto v0.28.0
)

exclude github.com/lib/pq v1.10.8

replace github.com/gorilla/mux => github.com/gorilla/mux v1.8.0

retract v0.5.0 // Contains a security vulnerability
```

## Step-by-step execution

Creating and verifying a `go.mod`:

1. `mkdir myapp && cd myapp`
2. `go mod init github.com/example/myapp`
   - Creates `go.mod` with only `module` and `go` lines.
3. Write `main.go` importing `github.com/gorilla/mux`.
4. `go mod tidy`
   - Downloads `gorilla/mux` and its dependencies.
   - Adds `require github.com/gorilla/mux v1.8.1` and any indirect deps.
   - Generates `go.sum`.
5. `go mod why -m golang.org/x/crypto`
   - Explains which package triggered `x/crypto` being a dependency.
6. `cat go.mod` — verify the file.

To add an `exclude`:

1. Manually edit `go.mod` or use: `go mod edit -exclude github.com/lib/pq@v1.10.8`

## Common mistakes

- **Editing `go.mod` manually without running `go mod tidy`**: Always run `go mod tidy` after manual edits to ensure consistency.
- **Removing `// indirect` comments**: The `// indirect` comment is significant. Tools like `go mod tidy` add it for a reason — removing it can cause confusion about why a dependency exists.
- **Using `replace` with an absolute path**: Use relative paths so the module works on other machines.
- **Forgetting to `go mod tidy` after changing dependencies**: Builds may use stale versions or fail with import errors.
- **Committing a `go.mod` with `replace` pointing to a local path**: This breaks other developers. Use `go.work` instead for local development, or remove `replace` before committing.
- **Setting `go` directive too high**: If you set `go 1.25` but use syntax only available in `go 1.26`, the code compiles but users on Go 1.25 get confusing errors. Match the directive to your actual Go version.

## Debugging walkthrough

Symptom: `go build` uses an unexpected version of a dependency.

```
go build ./...
go: finding module for package github.com/some/dep
go: github.com/some/dep@v1.5.0 found in go.mod, but...
```

**Root cause**: The `go.mod` requires `v1.5.0`, but a transitive dependency requires `v1.6.0`. MVS selects `v1.6.0`.

**Investigation**:
```bash
go list -m all | grep some/dep
go mod why -m github.com/some/dep
```

**Fix**: Either upgrade your direct dependency to use the newer API, or if the transitive dep's version requirement is an `exclude` candidate, add `exclude github.com/some/dep v1.6.0`.

Symptom: `go mod tidy` removes a dependency you think you need.

**Root cause**: The dependency is imported in a platform-specific file that is not being built on your current OS/arch.

**Fix**: Use build tags correctly. Run `go mod tidy` on the target platform, or set `GOOS` and `GOARCH` to the target.

## Production notes

- **Commit `go.mod` and `go.sum` to version control**: These define your reproducible build. Do not add them to `.gitignore`.
- **Use `go mod tidy` in CI**: Add `go mod tidy && git diff --exit-code go.mod go.sum` to verify that `go.mod` is clean.
- **Vendor mode**: If your team requires vendoring, use `go mod vendor` and check in the `vendor/` directory. Build with `-mod=vendor`.
- **`GONOSUMCHECK` / `GONOSUMDB` / `GOPRIVATE`**: For private modules, set `GOPRIVATE=*.corp.com` to bypass the checksum database and proxy.
- **Multi-module repos**: Use `go mod tidy` separately in each module. A shared CI script can iterate over all modules.
- **`retract` for bad releases**: If you accidentally publish a bad version, add a `retract` directive and tag a new version. The retraction is shown by `go list -m -versions` and `go get` avoids retracted versions.

## Performance implications

- `go mod tidy` reads the full import graph and can be slow on large modules (hundreds of dependencies). It runs in seconds for typical projects.
- The module file itself (a few KB) has zero runtime impact.
- A clean, minimal `go.mod` reduces `go mod tidy` time and makes code reviews easier.
- Excessive `replace` directives can slow down `go mod tidy` slightly as each replacement is resolved.

## Practice task

Create a new module manually (without `go mod init`):

1. Create a directory and write a `go.mod` file by hand with:
   - `module example.com/handcrafted`
   - `go 1.25.0`
   - `require rsc.io/quote v1.5.2`
2. Create `main.go` that imports `rsc.io/quote` and prints `quote.Hello()`.
3. Run `go mod tidy` to download the dependency and verify.
4. Add an `exclude` for `rsc.io/quote v1.5.2` and try building — observe the error.
5. Remove the `exclude` and add a `replace` pointing `rsc.io/quote` to `rsc.io/quote/v2 v2.0.0` (or any local directory).

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/19-go-mod
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/19-go-mod
```

## Review questions

1. What are the six main directives in a `go.mod` file?
2. What does the `// indirect` comment on a `require` line mean?
3. What is the purpose of `go mod tidy`?
4. How does `go mod why` help debug dependency issues?
5. What is Minimum Version Selection (MVS)?

## NEXT UP

go.sum — the checksum file that ensures reproducible, tamper-proof builds.
