# govulncheck

## Learning objective

Use the `govulncheck` tool to scan Go projects for known vulnerabilities, interpret its output, integrate it into CI/CD pipelines, and apply fixes for vulnerable dependencies.

## Why this matters

The Go vulnerability database tracks over 10,000 published CVEs affecting Go modules. Without automated scanning, vulnerabilities can go unnoticed for months or years. The average time to patch a known vulnerability in open-source dependencies is over 100 days. `govulncheck` is the official Go tool for this task, developed by the Go security team. It is integrated into gopls (the Go language server), VS Code, GitHub Actions, and CI pipelines. Go engineers must add vulnerability scanning to their development workflow to deploy secure software.

## Mental model

`govulncheck` is a security radar for your Go project. It continuously scans the dependency horizon for known threats (CVEs) and alerts you before they become problems. Unlike simple version checkers, `govulncheck` filters out noise: if a vulnerability exists in a module but your code never calls the vulnerable function, `govulncheck` does not report it. This is the difference between a weather radar that shows every storm in the hemisphere and one that only shows storms heading toward your location.

## Core idea

`govulncheck` works in three phases:

1. **Import graph analysis**: Scans your module and its dependencies to build the full import graph.
2. **Vulnerability matching**: Compares each module version against the Go vulnerability database (OSV format) to find known vulnerabilities affecting those versions.
3. **Call graph analysis**: For each vulnerable module, determines whether the vulnerable function or symbol is actually reachable from your application code. This eliminates false positives.

| Phase | What it checks | Output if problematic |
|---|---|---|
| Import graph | Which modules are used | List of modules |
| Vulnerability matching | Version against known CVEs | Modules with known vulns |
| Call graph analysis | Whether vulnerable code is reachable | Only reachable vulns reported |

## Under the hood

`govulncheck` uses the Go vulnerability database (https://vuln.go.dev), which publishes entries in the Open Source Vulnerability (OSV) format. Each OSV entry includes:

- `id`: Unique identifier (e.g., `GO-2024-0001`)
- `affected`: List of affected modules and version ranges
- `fixed`: Version where the vulnerability is fixed
- `description`: Human-readable description
- `severity`: CVSS score and vector
- `references`: Links to advisories, patches, and CVE entries
- `ecosystem_specific`: Go-specific fields including the vulnerable function/symbol name

`govulncheck` builds the call graph by analyzing your binary or source code. It uses the same static analysis infrastructure as the Go compiler. If the vulnerable function is never called (directly or transitively) from your `main` package or exported API, the vulnerability is not reported.

```text
Module A v1.0.0 -> has CVE-2024-XXXX -> function Foo() is vulnerable
  -> Your code imports module A but never calls Foo()
  -> govulncheck does NOT report it (unreachable)

Module B v2.0.0 -> has CVE-2024-YYYY -> function Bar() is vulnerable
  -> Your code calls Bar()
  -> govulncheck DOES report it (reachable)
```

## How Go uses it

Installation:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
```

Basic usage:

```bash
govulncheck ./...                  # Scan entire module
govulncheck -json ./...            # JSON output for CI
govulncheck -show verbose ./...    # Show full vulnerability details
```

Integration in CI (GitHub Actions):

```yaml
- name: govulncheck
  run: |
    go install golang.org/x/vuln/cmd/govulncheck@latest
    govulncheck ./...
```

Fixing vulnerabilities:

```bash
# Update the specific vulnerable module
go get golang.org/x/crypto@v0.17.0
# Or bump all dependencies (careful with breaking changes)
go get -u ./...
# Run govulncheck again to verify
govulncheck ./...
```

## Go example

```go
package main

import (
	"fmt"
	"strings"
)

type VulnEntry struct {
	ID          string
	Module      string
	Version     string
	Description string
	FixedIn     string
	Severity    string
}

type ModuleVersion struct {
	Path    string
	Version string
}

func checkVulnerabilities(modules []ModuleVersion, vulns []VulnEntry) []VulnEntry {
	var found []VulnEntry
	for _, mod := range modules {
		for _, v := range vulns {
			if v.Module == mod.Path && strings.HasPrefix(mod.Version, strings.TrimSuffix(v.Version, "+")) {
				found = append(found, v)
			}
		}
	}
	return found
}

func main() {
	vulns := []VulnEntry{
		{ID: "GO-2024-0001", Module: "golang.org/x/crypto", Version: "<0.17.0",
			Description: "DoS in RSA key generation", FixedIn: "0.17.0", Severity: "HIGH"},
	}
	modules := []ModuleVersion{
		{Path: "golang.org/x/crypto", Version: "v0.16.0"},
	}
	results := checkVulnerabilities(modules, vulns)
	for _, v := range results {
		fmt.Printf("[%s] %s: %s (fix: %s)\n", v.Severity, v.ID, v.Description, v.FixedIn)
	}
}
```

## Step-by-step execution

Running `govulncheck ./...` on a Go project:

1. `govulncheck` reads `go.mod` to find the module path and dependencies.
2. It downloads the latest vulnerability database from `vuln.go.dev` (cached locally).
3. For each module in `go.mod`, it checks if the version falls within any affected range in the database.
4. For matching vulnerabilities, it performs call graph analysis: can the vulnerable function be reached from your code?
5. It prints results:
   - **Vulnerabilities** (reachable and unfixed): requires immediate action
   - **Vulnerable modules** (exist in deps but unreachable): informational, still worth fixing
   - **Fixed versions available**: upgrade recommendations

Example output:

```
=== Informational ===
Found 1 vulnerability in module golang.org/x/crypto
  GO-2024-0001: Denial of service in RSA key generation
  More: https://pkg.go.dev/vuln/GO-2024-0001
  Found in: golang.org/x/crypto@v0.16.0
  Fixed in: golang.org/x/crypto@v0.17.0

Your code is affected by 1 vulnerability.
Run 'go get golang.org/x/crypto@v0.17.0' to fix.
```

## Common mistakes

- Mistake: Running `govulncheck` only once during initial setup instead of integrating it into CI/CD.
  - Why it happens: Developers check once, see no vulns, and never run it again.
  - Fix: Add `govulncheck` to CI pipeline. Run on every commit, or at least daily.

- Mistake: Ignoring `govulncheck` findings because they are "only in test dependencies."
  - Why it happens: Test code imports production packages as transitive deps. A vuln in a test dep can still affect production if the vulnerable package is also a transitive dep of production code.
  - Fix: Treat all findings seriously. Update whatever module is flagged.

- Mistake: Not updating the vulnerability database regularly.
  - Why it happens: `govulncheck` caches the database. A stale cache misses newly published CVEs.
  - Fix: Run `govulncheck` without cache or clear the cache periodically: `go clean -cache`.

- Mistake: Fixing the wrong package: updating the direct dependency when the CVE is in a transitive dependency.
  - Why it happens: Developers see `github.com/gorilla/mux` in `go.mod` and run `go get github.com/gorilla/mux@latest`, but the vuln is in `golang.org/x/net` which mux depends on.
  - Fix: Read the full `govulncheck` output to identify the actual vulnerable module. Update that module directly.

- Mistake: Assuming `govulncheck` catches all vulnerabilities.
  - Why it happens: Developers think scanning is a silver bullet.
  - Fix: `govulncheck` only catches known, published vulnerabilities. It does not catch zero-days, logic flaws, or misconfigurations. Use it as one layer of a defense-in-depth strategy.

## Debugging walkthrough

Consider this CI failure:

```
$ govulncheck ./...
Found 1 vulnerability in module golang.org/x/crypto
  GO-2024-0001: Denial of service in RSA key generation
  Found in: golang.org/x/crypto@v0.16.0
  Fixed in: golang.org/x/crypto@v0.17.0
```

Symptom: CI pipeline fails due to a vulnerability in a transitive dependency.

Investigation: `go mod graph | grep golang.org/x/crypto` reveals which direct dependency pulls it in.

Root cause: `github.com/gorilla/mux v1.8.0` depends on `golang.org/x/crypto@v0.16.0` (as a transitive dep through another module).

Fix: Update the vulnerable module directly, even if it is transitive:

```bash
go get golang.org/x/crypto@v0.17.0
go mod tidy
govulncheck ./...  # Verify fix
```

If the direct dependency does not support the fixed version, update the direct dependency too:

```bash
go get github.com/gorilla/mux@latest
go mod tidy
```

## Production notes

- Add `govulncheck` to CI as a required check. Block merges on new vulnerabilities.
- Run `govulncheck` in a scheduled job (daily) in addition to on-commit, because new vulnerabilities are published daily.
- Use `govulncheck -json` in CI for machine-parseable output that can be posted to Slack or a dashboard.
- Monitor the Go vulnerability database RSS feed (https://groups.google.com/g/golang-announce) for new advisories.
- For enterprise compliance, generate regular vulnerability reports using `govulncheck -json` and archive them.
- Consider `golangci-lint` which includes govulncheck as one of its linters for a unified linting pipeline.

## Performance implications

- `govulncheck` completes in seconds for most projects (typically 2-10 seconds).
- The first run downloads the vulnerability database (~10MB). Subsequent runs use a cached version.
- Call graph analysis adds some time (5-30 seconds for large projects) but significantly reduces noise.
- CI integration adds negligible overhead to build times compared to the security benefit.
- Running `govulncheck` on every commit is practical for projects of any size.

## Practice task

Write a function `runGovulncheckReport(modulePath string, vulnDBPath string) ([]VulnEntry, error)` that:

1. Loads a vulnerability database from a JSON file (OSV format, simplified).
2. Parses `go.mod` of the given module path to extract dependency versions.
3. Checks each dependency against the vulnerability database.
4. Returns a list of matching vulnerabilities with severity, description, and fix version.

Then write a `main()` that creates a mock `go.mod`, a mock vulnerability database, runs the report, and prints results including a summary line like "Found 2 vulnerabilities. Run 'go get ...' to fix."

## Tests / verification

```bash
go run ./curriculum/modules/10-auth-security/lessons/22-govulncheck
go test ./curriculum/modules/10-auth-security/lessons/22-govulncheck
```

The existing tests verify vulnerability database loading, module version matching, detection of vulnerable modules, and that patched versions are not flagged.

## Review questions

1. How does `govulncheck` differ from a simple version-matching vulnerability scanner?
2. What is the Go vulnerability database, and what format does it use?
3. Why would `govulncheck` report a vulnerability in a transitive dependency but not in a direct dependency?
4. How do you fix a vulnerability reported by `govulncheck` if the fix requires a breaking change in a dependency?
5. What are the limitations of `govulncheck` and what other security tools should complement it?

## NEXT UP

OWASP for Go APIs -- applying the OWASP API Security Top 10 to Go API development, with a security checklist and common Go-specific vulnerabilities.
