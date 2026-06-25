# Workspaces

## Learning objective

Create and use Go workspaces with `go.work` files for local multi-module development, understand the `use` directive, and compare workspaces to `replace` directives.

## Why this matters

When you work on multiple Go modules simultaneously — for example, your application and a helper library you are developing in parallel — you need to test changes across modules without publishing new versions. Before Go 1.18, the only option was the `replace` directive in `go.mod`. Workspaces (`go.work`) provide a cleaner, more scalable solution. They are essential for monorepo development, multi-module microservices, and library authors testing changes against consumer code.

## Mental model

A workspace is a workbench where you lay out multiple modules side by side.

- The **`go.work`** file is the workbench blueprint — it says "these modules are on my bench."
- Each **`use`** directive points to a module on the bench.
- When you build, Go **replaces** the published versions of those modules with your local copies — without modifying any `go.mod` files.
- Modules not listed in `go.work` are fetched normally from the module proxy.

Think of `go.work` as a **temporary, developer-local overlay** on top of the module graph. It does not affect published builds or other developers unless they also use workspaces.

## Core idea

### The `go.work` file

```
go 1.25.0

use (
    ./myapp
    ./mylib
)
```

- **`go` directive**: The Go version used to interpret the workspace.
- **`use`**: Paths to local modules (relative to the `go.work` file's directory).
- **`replace`**: Workspace-level `replace` directives (optional).

### Creating a workspace

```bash
go work init ./myapp ./mylib
```

This creates a `go.work` file with `use` directives for the specified modules.

### How it differs from `replace`

| Aspect | `replace` in `go.mod` | `go.work` workspace |
|---|---|---|
| Scope | Single module | All modules in the workspace |
| Persistence | Committed to VCS | Local-only (add to `.gitignore`) |
| Overrides | Applied only to that module | Applied globally during workspace builds |
| Merging | Requires editing each `go.mod` | Single file for all modules |
| Team use | Must be in VCS | Each developer creates their own |

### Using workspace mode

When a `go.work` file is present in the current directory or any parent, Go automatically enters workspace mode. Explicit commands:

```bash
go work init       # create go.work
go work use ./mod  # add a module to the workspace
go work sync       # sync workspace dependencies across modules
```

## Under the hood

When Go builds in workspace mode:

1. It reads `go.work` to find all `use` directives.
2. For each `use`d module, it reads the module's `go.mod`.
3. It creates a synthetic module graph: the workspace acts as the root module, with all `use`d modules as direct dependencies.
4. Any `go.work` `replace` directives are applied to the entire graph.
5. MVS runs on this combined graph — local versions may satisfy dependencies that would otherwise resolve to published versions.
6. The `go.work` file does **not** affect the module cache or published builds — it is purely a local development overlay.

The `go.work` file is intentionally not part of the module's public API. It is not used by `go build` when building a module for publishing.

## How Go uses it

- **Monorepo development**: A repository with multiple Go modules (e.g., `api/`, `services/`, `libs/`) uses a root `go.work` to link them all for local development.
- **Library authoring**: A library author maintains a separate test consumer project and uses a workspace to test changes before releasing.
- **Microservice development**: Each microservice is its own module; a workspace links them for integration testing.
- **Replacing `replace`**: Teams that previously used `replace` directives in every module's `go.mod` now use a single `go.work` file.
- **`go work sync`**: Ensures that all modules in the workspace use consistent dependency versions.

## Go example

```
// Directory structure:
//   workspace-demo/
//     go.work
//     app/
//       go.mod          → module example.com/app
//       main.go
//     lib/
//       go.mod          → module example.com/lib
//       lib.go
```

```go
// go.work
go 1.25.0

use (
    ./app
    ./lib
)
```

```go
// app/main.go
package main

import (
	"fmt"

	"example.com/lib"
)

func main() {
	fmt.Println(lib.Greet("Workspace"))
}
```

```go
// lib/lib.go
package lib

func Greet(name string) string {
	return "Hello, " + name + "!"
}
```

## Step-by-step execution

Creating and using a workspace:

1. Create `app/` and `lib/` directories, each with their own `go.mod` and source files.
2. From the parent directory, run: `go work init ./app ./lib`
3. This creates `go.work` with `use (./app ./lib)`.
4. `go run ./app` — builds and runs, using the local `lib` module.
5. Modify `lib/lib.go` — changes are visible immediately in `app` without `go get` or version bumps.
6. `go work use ./newmod` — adds another module.
7. `go work sync` — syncs dependencies across all workspace modules.

To exit workspace mode: delete or rename `go.work`. Builds then use published versions.

## Common mistakes

- **Committing `go.work` to version control**: `go.work` is a local development file. Add it to `.gitignore` so it does not pollute CI or other developers' environments. Each developer creates their own.
- **Forgetting `use` directives**: Adding a dependency that exists in a local module but is not listed in `use` causes Go to fetch the published version instead.
- **Mixing `replace` in `go.mod` with workspace**: If both are present, workspace `use` takes precedence for modules listed in `use`. But `replace` directives in `go.mod` still apply for modules not in the workspace.
- **Building without `go.work` present**: If you `cd` into a submodule and build, Go does not look for `go.work` in parent directories past the module root. Keep `go.work` at a common ancestor.
- **Expecting `go mod tidy` to use workspace modules**: `go mod tidy` runs in the context of a module, not the workspace. After `go mod tidy`, the `go.mod` may still reference published versions, but workspace mode overrides them at build time.

## Debugging walkthrough

Symptom: `go run ./app` uses a published version of `lib` instead of the local one.

**Root cause**: The `go.work` file exists but does not include `./lib` in `use`.

**Fix**: Add `./lib` to the `use` directive:
```bash
go work use ./lib
```

Symptom: `go build ./...` in a workspace complains about inconsistent versions.

**Root cause**: Two modules in the workspace require different versions of the same dependency.

**Fix**: Use `go work sync` to align versions:
```bash
go work sync
```

Symptom: CI fails because it does not have `go.work` and cannot find local modules.

**Root cause**: A module's `go.mod` has a `replace` directive pointing to a local path that only exists on the developer's machine. The `replace` was intended for workspace development but was committed.

**Fix**: Remove `replace` directives from `go.mod` and use `go.work` locally. Add `go.work` to `.gitignore`.

## Production notes

- **`go.work` is local-only**: The Go team's official recommendation is to keep `go.work` out of version control. Each developer generates their own with `go work init`.
- **CI/CD**: In CI, do not use `go.work`. Build each module independently with `go build ./...` from the module root. Use `go mod download` to cache dependencies.
- **Multi-module CI**: Set up separate CI jobs for each module, or use a matrix build. A workspace is not needed for CI — it is a development convenience.
- **VS Code / gopls**: The Go language server supports workspace mode. When `go.work` is present, `gopls` resolves imports across modules seamlessly.
- **Migration from `replace`**: If your team uses `replace` directives extensively, migrate to a workspace. The workspace is cleaner, does not require editing each `go.mod`, and prevents accidental commits of local paths.

## Performance implications

- Workspace mode adds minimal build-time overhead — the workspace file is parsed once, and module resolution runs on the combined graph.
- The module cache is shared between workspace and non-workspace builds. Local modules in `use` bypass the cache (they are read from disk) but compilation speed is unaffected.
- `go work sync` runs MVS across all workspace modules, which can be slightly slower than single-module resolution but still completes in milliseconds.
- There is zero runtime cost — workspace mode is a build-time concept only.

## Practice task

Create a two-module workspace:

1. Create directories: `mkdir -p /tmp/workspace-lab/{app,lib}`
2. In `lib/`: `go mod init example.com/lib` and create `lib.go` with an exported `Help() string` function.
3. In `app/`: `go mod init example.com/app` and create `main.go` that imports `example.com/lib` and calls `lib.Help()`.
4. From `/tmp/workspace-lab`, run `go work init ./app ./lib`.
5. Run `go run ./app` and verify it works.
6. Modify `lib.Help()` to return a different string and re-run — observe the local change is picked up.
7. Remove `go.work` and try building `./app` — it fails because `example.com/lib` is not available.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/22-workspaces
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/22-workspaces
```

## Review questions

1. What command creates a `go.work` file with two modules?
2. How does a workspace differ from a `replace` directive in `go.mod`?
3. Should you commit `go.work` to version control? Why or why not?
4. What happens if a module is imported but not listed in the workspace's `use` directive?
5. What is the purpose of `go work sync`?

## NEXT UP

Documentation comments — how to write godoc-compatible documentation for packages, types, functions, and examples.
