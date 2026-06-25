# Dependency security

## Learning objective

Assess and mitigate software supply chain risks in Go projects by verifying dependency integrity with `go.sum`, minimizing dependency surface area, reviewing dependencies for malicious code, and maintaining a software bill of materials.

## Why this matters

Software supply chain attacks are one of the fastest-growing security threats. The 2024 XZ Utils backdoor (CVE-2024-3094) and the 2018 event-stream incident demonstrated that a compromised dependency can give attackers access to thousands of downstream projects. Go's module system includes strong supply chain security features: `go.sum` for tamper-proof dependency verification, the checksum database for global transparency, and reproducible builds. Go engineers must understand these features and the broader supply chain risk landscape to protect their applications.

## Mental model

Dependencies are like building materials in construction. You need to verify that each brick comes from a trusted supplier and is not a counterfeit. `go.sum` is your material manifest: it lists every module and its cryptographic fingerprint. The checksum database is the industry-wide ledger that lets you verify no one has tampered with the module since it was published. When you add a dependency, you are vetted by the supplier and accepting the risk that they might ship defective or malicious material.

A software bill of materials (SBOM) is a complete inventory of every building material used in your project, including supplier, version, and license.

## Core idea

Go's dependency security model rests on three pillars:

1. **`go.sum` file**: Contains cryptographic hashes (SHA-256) of every module version your project uses. `go mod verify` checks that the cached module content matches `go.sum`.
2. **Checksum database** (`sum.golang.org`): A global, append-only, transparent log of module hashes. When you run `go mod download`, Go queries the checksum database to ensure the module hash matches the known good hash. This prevents a compromised module proxy from serving tampered code.
3. **Reproducible builds**: `go mod vendor` and `go.sum` ensure that building the same module version always produces the same code.

| Security feature | What it prevents | How it works |
|---|---|---|
| `go.sum` | Module tampering after first use | Stores expected SHA-256 hash of each module |
| Checksum database | Module proxy serving tampered code | Global log of hashes, tamper-evident |
| `go mod verify` | Local cache tampering | Recomputes hashes and compares with `go.sum` |
| Module proxy | Module disappearance or modification | Caches immutable module versions |
| `GONOSUMCHECK` / `GONOSUMDB` | Internal modules (opt-out) | Skips checksum database for private modules |

## Under the hood

When `go mod download` fetches a module:

1. Go queries the module proxy (default: `proxy.golang.org`).
2. The proxy returns the module content and its hash.
3. Go looks up the hash in `go.sum`. If present, it verifies the downloaded content matches.
4. If not present in `go.sum`, Go queries the checksum database `sum.golang.org` to get the global hash.
5. Go appends the new hash to `go.sum`.
6. If the checksum database is unavailable and the module is not in `go.sum`, Go blocks the build (unless `GONOSUMCHECK` is set).

The `go.sum` file contains entries like:

```
github.com/gorilla/mux v1.8.1 h1:TuMF1mMtQ/b/NgMm0w==
github.com/gorilla/mux v1.8.1/go.mod h1:1lud6UwP+6orDFRuTfBEV8e9/aOMD4/
```

Each entry includes: module path, version, hash algorithm (h1: = SHA-256), and the hex-encoded hash. Go modules record both the source hash and the `go.mod` hash separately.

## How Go uses it

Go uses `go.sum` and the checksum database automatically. Key commands:

```bash
go mod download        # Downloads dependencies and verifies go.sum
go mod verify          # Verifies cached modules match go.sum
go mod tidy            # Adds missing go.sum entries, removes unused deps
GONOSUMCHECK=* go mod download  # Skips checksum verification (insecure)
```

For internal/private modules not in the checksum database, set `GONOSUMDB` and `GONOSUMCHECK` environment variables:

```bash
export GONOSUMDB=*.internal.example.com
export GONOSUMCHECK=*.internal.example.com
```

Minimizing dependencies is a manual but critical practice. Audit imports regularly and remove unused dependencies.

## Go example

```go
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

type ModuleEntry struct {
	Path    string
	Version string
	Hash    string
}

func parseGoSum(sumFile string) ([]ModuleEntry, error) {
	data, err := os.ReadFile(sumFile)
	if err != nil {
		return nil, err
	}
	var entries []ModuleEntry
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			entries = append(entries, ModuleEntry{
				Path:    parts[0],
				Version: parts[1],
				Hash:    parts[2],
			})
		}
	}
	return entries, nil
}

func main() {
	fmt.Println("=== Go Dependency Security ===")

	entries := []ModuleEntry{
		{Path: "github.com/gorilla/mux", Version: "v1.8.1", Hash: "h1:TuMF1mMtQ=="},
		{Path: "golang.org/x/crypto", Version: "v0.17.0", Hash: "h1:abc123=="},
	}
	for _, e := range entries {
		fmt.Printf("Module: %-35s %-12s %s\n", e.Path, e.Version, e.Hash)
	}
	fmt.Println("\nRun 'go mod verify' to verify cache integrity.")
}
```

## Step-by-step execution

When you add a new dependency to your Go project:

1. `go get github.com/example/lib@v1.0.0` resolves the module version from the proxy.
2. Go downloads the module content and computes its SHA-256 hash: `h1:abc123def456...`.
3. Go queries `sum.golang.org` for the known hash of `github.com/example/lib@v1.0.0`.
4. If the hash matches, Go appends to `go.sum`:
   ```
   github.com/example/lib v1.0.0 h1:abc123def456...=
   ```
5. If the hash does not match the checksum database, Go prints a security error and refuses to build.
6. On subsequent builds, Go verifies the cached module content matches `go.sum` before using it.

For an attacker trying to serve a tampered version:

1. Attacker compromises the module proxy and replaces `example/lib` with malicious code.
2. The malicious code has a different hash: `h1:xyz789...`.
3. Go queries the checksum database and finds the expected hash is `abc123...`.
4. Go rejects the download with: `security: checksum mismatch`.
5. The attacker would need to compromise the checksum database, which is append-only and tamper-evident.

## Common mistakes

- Mistake: Running `go mod vendor` without reviewing the vendored code.
  - Why it happens: Teams vendor dependencies and trust them blindly.
  - Fix: Review vendored changes on every dependency update. Use `git diff vendor/` during code review.

- Mistake: Using `@latest` version of a dependency without pinning.
  - Why it happens: Developers want to always get the latest version.
  - Fix: Pin to a specific version: `@v1.2.3`. Use Dependabot or Renovate for automated updates with review.

- Mistake: Adding direct dependencies for functionality that Go's standard library already provides.
  - Why it happens: Developers reach for third-party packages out of habit.
  - Fix: Prefer the standard library. For example, use `net/http` instead of `gorilla/mux` for simple routing, or `encoding/json` instead of `json-iterator/go`.

- Mistake: Not running `go mod verify` after a security incident.
  - Why it happens: Teams do not know about the command.
  - Fix: Add `go mod verify` to CI pipeline and run it after any security-relevant incident.

- Mistake: Ignoring `go.sum` merge conflicts.
  - Why it happens: `go.sum` is auto-generated, so teams often blindly accept the merged version.
  - Fix: Run `go mod tidy` after resolving `go.sum` merge conflicts to regenerate the correct hashes.

## Debugging walkthrough

Consider this CI failure:

```
$ go mod verify
github.com/example/lib@v1.0.0: hash mismatch
    expected: h1:abc123def456...
    got:      h1:xyz789abc123...
```

Symptom: CI build fails on a previously working project without any dependency changes.

Investigation: Check if the `go.sum` file was modified by a merge conflict, or if the module proxy is returning a different version:

```bash
go mod download -json github.com/example/lib@v1.0.0 | jq .
```

Root cause: A team member resolved a `go.sum` merge conflict by keeping both conflicting lines, creating an ambiguous hash entry. Go chooses the first matching line, which is the wrong hash.

Fix: Run `go mod tidy` to regenerate `go.sum` from the correct dependency graph:

```bash
go mod tidy
git add go.sum
git commit -m "fix: regenerate go.sum after merge conflict"
```

## Production notes

- Run `go mod verify` in CI to detect dependency tampering.
- Use a private module proxy (Athens, Artifactory, GoCenter) for internal modules and to cache public modules against availability issues.
- Set `GONOSUMDB` and `GONOSUMCHECK` for private modules. Never set them globally -- that disables supply chain security for all modules.
- Generate an SBOM (Software Bill of Materials) using `github.com/spdx/tools-golang` or `syft` for compliance and vulnerability tracking.
- Review dependency updates in a dedicated PR. Do not merge dependency updates without review.
- Periodically audit dependencies with `go mod graph` to identify unnecessary transitive dependencies.
- Consider using `golang.org/x/tools/cmd/deadcode` to find unused code that may pull in unnecessary dependencies.

## Performance implications

- `go.sum` parsing adds milliseconds to build time.
- Module download time depends on network speed and module size. Using a local proxy cache reduces this significantly.
- The number of dependencies has a minor impact on compile time (more packages to compile). The larger impact is on binary size.
- Minimizing dependencies reduces both build time and attack surface. Each dependency is a potential supply chain risk.
- Using the standard library is almost always faster at compile time and produces smaller binaries than third-party alternatives.

## Practice task

Write a function `auditDependencies(goSumPath string) error` that:

1. Parses the `go.sum` file into module entries.
2. Checks if any dependency has a known vulnerability (simulate by checking a hardcoded list of "known bad" module paths).
3. Prints a report of all dependencies, marking known bad ones.
4. Returns an error if any known bad dependency is found.

Then write a `main()` that creates a test `go.sum` file with several entries (one being a "known bad" module), runs the audit, and prints results.

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/21-dependency-security
go test ./curriculum/modules/10-auth-security/lessons/21-dependency-security
```

The existing tests verify `go.sum` parsing, module hash computation, dependency minification, and empty-sum handling.

## Review questions

1. How does `go.sum` protect against module tampering, and what is the role of the checksum database?
2. What is the difference between `GONOSUMCHECK` and `GONOSUMDB`, and why should neither be set globally?
3. How would you detect a compromised dependency in a Go project using the standard toolchain?
4. Why is minimizing dependencies a security practice, not just a convenience?
5. What information does a Software Bill of Materials contain, and why is it useful for incident response?

## NEXT UP

govulncheck -- scanning Go projects for known vulnerabilities using the official Go vulnerability tool, integrating it into CI, and understanding the vulnerability database format.
