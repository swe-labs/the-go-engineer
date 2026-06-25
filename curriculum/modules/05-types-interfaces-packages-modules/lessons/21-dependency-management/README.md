# Dependency management

## Learning objective

Add, upgrade, downgrade, and remove Go module dependencies using `go get`, understand Minimum Version Selection (MVS), use `go mod vendor` to create a vendor directory, and configure `GONOSUMCHECK`, `GONOSUMDB`, and `GOPRIVATE` for private modules.

## Why this matters

Real Go projects depend on dozens or hundreds of external packages. Managing these dependencies correctly — adding the right version, upgrading safely, removing unused packages, and handling private modules — is a daily task for professional Go engineers. Mistakes here cause build failures, security vulnerabilities from outdated deps, and broken CI pipelines. Mastery of `go get` and related tools keeps your project healthy.

## Mental model

Your module's dependency graph is a tree. `go get` is your gardening tool:

- `go get module@version` — plant a specific seed (add or pin a dependency).
- `go get -u` — water and fertilise everything (upgrade all direct and indirect deps).
- `go get module@none` — pull a weed (remove a dependency).
- `go mod tidy` — prune dead branches (remove unused imports).
- `go mod vendor` — take cuttings (copy dependencies into your workspace for offline builds).

**Minimum Version Selection (MVS)** is the rule: when multiple requirements for the same module conflict, Go chooses the minimum version that satisfies all requirements. This is simpler, more predictable, and more stable than the dependency-resolution algorithms used by npm, pip, or RubyGems.

## Core idea

### `go get` — the primary dependency tool

```bash
go get github.com/gorilla/mux              # latest version
go get github.com/gorilla/mux@v1.8.1       # specific version
go get github.com/gorilla/mux@v1.8.0       # downgrade
go get github.com/gorilla/mux@latest        # latest
go get github.com/gorilla/mux@none          # remove
go get -u ./...                             # upgrade all direct and indirect deps
go get -u=patch ./...                       # upgrade to latest patch version only
```

Each `go get` command updates `go.mod` and `go.sum` automatically.

### Minimum Version Selection (MVS)

MVS works like this:

1. Start with the `require` directives in the root module.
2. For each dependency, find the version that satisfies all its dependents.
3. Select the **minimum** such version.

Example: Module A requires `C v1.0.0`. Module B requires `C v1.2.0` and `A`. MVS selects `C v1.2.0` because that satisfies both A (which accepts any `v1.x`) and B.

MVS is deterministic: given the same `go.mod`, every machine selects the same versions. There are no lock files or resolution order dependencies.

### Pruning unused dependencies

`go mod tidy` removes `require` entries for packages that are no longer imported by any package in the module. This keeps `go.mod` minimal.

### `vendor` directory

```bash
go mod vendor
```

Creates a `vendor/` directory containing copies of all dependency source files. Build with `-mod=vendor` to use them:

```bash
go build -mod=vendor ./...
```

Vendoring is useful for:
- Air-gapped or offline builds.
- Ensuring dependencies are available even if the upstream disappears.
- Code review (reviewers can see exactly which dependency code is included).

### Private module configuration

For modules hosted on private servers, configure these environment variables:

| Variable | Purpose |
|---|---|
| `GOPRIVATE` | Shorthand that sets both `GONOSUMDB` and `GONOSUMCHECK` for the specified patterns. |
| `GONOSUMDB` | Skip checksum database lookup for matching modules. |
| `GONOSUMCHECK` | Skip hash verification for matching modules. |
| `GOINSECURE` | Allow HTTP (not HTTPS) and skip certificate verification for matching modules. |

```bash
export GOPRIVATE=github.com/myorg/*
```

Go reads these from environment variables or the `.netrc` / `gitconfig` for authentication.

## Under the hood

The module resolution algorithm:

1. **Parse**: Read `go.mod` for direct `require` entries.
2. **Build graph**: For each dependency, read its `go.mod` to find its own requirements.
3. **Select**: Apply MVS — pick the minimum version that satisfies each constraint.
4. **Replace**: Apply `replace` directives, if any.
5. **Exclude**: Remove `exclude`-ed versions from consideration.
6. **Download**: Fetch selected versions from proxy or direct source.
7. **Verify**: Check hashes against `go.sum` and the checksum database.
8. **Build**: Compile all packages.

The entire dependency graph is computed at build time, not at `go get` time. `go get` just modifies `go.mod`; `go build` resolves and verifies.

## How Go uses it

- **Adding a new dependency**: `go get github.com/new/lib`
- **Upgrading all dependencies**: `go get -u ./...`
- **Upgrading to latest patch**: `go get -u=patch ./...`
- **Removing an unused dependency**: `go mod tidy` (removes unused entries) or `go get module@none`
- **Inspecting why a dependency is needed**: `go mod why -m module/path`
- **Listing all module versions**: `go list -m -versions module/path`
- **Seeing the resolved dependency graph**: `go list -m all`

## Go example

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	tmpDir, err := os.MkdirTemp("", "dep-demo")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	// 1. Init module.
	runCmd(tmpDir, "go", "mod", "init", "example.com/dep-demo")

	// 2. Create main.go.
	mainSrc := `package main
import "fmt"
import "rsc.io/quote"
func main() { fmt.Println(quote.Hello()) }
`
	os.WriteFile(tmpDir+"/main.go", []byte(mainSrc), 0644)

	// 3. go get a specific version.
	runCmd(tmpDir, "go", "get", "rsc.io/quote@v1.5.2")

	// 4. Show go.mod.
	data, _ := os.ReadFile(tmpDir + "/go.mod")
	fmt.Printf("go.mod after go get:\n%s\n", data)

	// 5. List all modules.
	runCmd(tmpDir, "go", "list", "-m", "all")

	// 6. Upgrade to patch.
	runCmd(tmpDir, "go", "get", "-u=patch", "./...")

	// 7. Show final go.mod.
	data, _ = os.ReadFile(tmpDir + "/go.mod")
	fmt.Printf("After patch upgrade:\n%s\n", data)
}

func runCmd(dir, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	fmt.Printf("> %s %v\n%s", name, args, out)
	if err != nil {
		panic(fmt.Sprintf("%v: %v\n%s", name, err, out))
	}
}
```

## Step-by-step execution

Adding and upgrading a dependency:

1. `go get github.com/gorilla/mux` — resolves the latest version, adds to `go.mod`, downloads, adds to `go.sum`.
2. Write code importing `github.com/gorilla/mux`.
3. `go build ./...` — builds using the resolved version.
4. `go list -m -versions github.com/gorilla/mux` — see available versions.
5. `go get github.com/gorilla/mux@v1.7.0` — downgrade to test compatibility.
6. `go get -u=patch ./...` — upgrade all deps to the latest patch versions within their major/minor.
7. `go mod tidy` — remove any now-unused dependencies.

Removing a dependency:

1. Remove all imports of the dependency from your Go source files.
2. `go mod tidy` — removes the `require` entry and the corresponding `go.sum` entries.
3. Alternatively: `go get module@none`.

## Common mistakes

- **Running `go get -u` without testing**: `-u` upgrades all direct and indirect dependencies. Newer versions may introduce breaking changes or API incompatibilities. Always run tests after `go get -u`.
- **Using `go get` outside a module**: In Go 1.17+, `go get` outside a module warns that it only works in module mode. Use `go install` for installing tools.
- **Not running `go mod tidy` after removing imports**: The `require` entry remains in `go.mod` until you run `go mod tidy`.
- **Vendoring without `-mod=vendor`**: After running `go mod vendor`, you must add `-mod=vendor` to build commands, or set `GOFLAGS=-mod=vendor`. Otherwise, Go uses the module cache.
- **Forgetting to set `GOPRIVATE` for private modules**: Builds fail because the checksum database cannot access private repos.
- **Expecting `go get -u` to respect semver ranges**: `-u` upgrades to the latest version, not the latest within a range. Use `-u=patch` for conservative upgrades.

## Debugging walkthrough

Symptom: `go get` fails with:

```
go: github.com/private/module@v1.0.0: verifying module: checksum mismatch
```

**Root cause**: The module is private but `GONOSUMDB` is not set. Go tried to verify the checksum against `sum.golang.org`, which does not have the private module's hash.

**Fix**:
```bash
export GOPRIVATE=github.com/myorg/*
go get github.com/private/module@v1.0.0
```

Symptom: `go build` uses an older version than expected.

**Root cause**: A transitive dependency requires the older version, and MVS selects it. Run `go list -m all | grep module` and `go mod why -m module` to trace the requirement.

**Fix**: Either upgrade the transitive dependency (by upgrading the direct dependency that requires it) or add a `require` directive in your module to force a newer version (MVS will select the maximum).

## Production notes

- **Pin direct dependencies, let indirect float**: Only `require` the versions you directly import. Let transitive dependencies be handled by MVS. Run `go mod tidy` to keep indirect entries accurate.
- **Weekly dependency updates**: Schedule a weekly CI job to run `go get -u=patch ./... && go mod tidy && go test ./...` to keep dependencies patched.
- **Security updates**: Subscribe to Go vulnerability database (`vuln.go.dev`). Use `govulncheck` to scan for known vulnerabilities.
- **Vendor in regulated environments**: Finance, healthcare, and government projects often require vendoring for auditability. Use `go mod vendor` and commit the `vendor/` directory.
- **Dependency review**: Treat `go.mod` and `go.sum` changes with the same scrutiny as source changes. Review every new or updated dependency.
- **`GONOSUMCHECK` and `GONOSUMDB`**: Set `GOPRIVATE=*corp.com` for all internal modules. This avoids checksum database lookups and hash verification for private code.

## Performance implications

- `go get` downloads the dependency and its transitive dependencies. Initial download time depends on dependency count and network speed.
- Subsequent builds use the module cache (no network access), so dependency count has minimal impact on build time.
- The `vendor/` directory adds ~1-10MB to the repository size. CI checkout times increase proportionally.
- MVS resolution is O(|graph|) and completes in milliseconds for typical projects (tens of microseconds per dependency).

## Practice task

1. Create a module: `mkdir /tmp/dep-lab && cd /tmp/dep-lab && go mod init example.com/dep-lab`
2. `go get rsc.io/quote@v1.5.2`
3. Write `main.go` that uses `quote.Hello()` and prints it.
4. `go build && go list -m all` — observe the dependency list.
5. `go mod why -m rsc.io/sampler` — trace why `sampler` is included.
6. `go get rsc.io/quote@v1.5.1` — downgrade and build.
7. `go mod tidy` — remove unused modules.
8. `go mod vendor` — create vendor directory.
9. Build with vendored dependencies: `go build -mod=vendor ./...`

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/21-dependency-management
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/21-dependency-management
```

## Review questions

1. What is the difference between `go get -u` and `go get -u=patch`?
2. How does Minimum Version Selection (MVS) differ from npm's dependency resolution?
3. When would you use `go mod vendor` instead of relying on the module cache?
4. What environment variables control access to private modules?
5. Why should you run `go mod tidy` after removing imports from your code?

## NEXT UP

Workspaces — developing multiple modules locally with `go.work` for seamless multi-module workflows.
