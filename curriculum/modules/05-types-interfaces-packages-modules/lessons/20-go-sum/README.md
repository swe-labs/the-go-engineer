# go.sum

## Learning objective

Explain the purpose of `go.sum`, how it stores content hashes for integrity verification, and how to use `go mod verify` and navigate the checksum database with appropriate privacy considerations.

## Why this matters

When you `go build`, the Go toolchain downloads modules from the internet. How do you know the downloaded code is exactly what the author published — and not a tampered copy injected with malware? `go.sum` is the answer. It ensures that every build uses byte-for-byte identical source code. Without it, your builds are not reproducible and are vulnerable to supply-chain attacks. This is not theoretical — the checksum database has prevented real-world attacks.

## Mental model

`go.sum` is a lockbox of cryptographic fingerprints:

- Each line is a hash of a specific module version's `go.mod` file and its full source tree.
- When you first download a module version, its hash is recorded.
- On subsequent builds (or other developers' machines), the toolchain recomputes the hash and compares it against `go.sum`.
- If the hashes match, the source is authentic. If they differ, the build fails with a security error.
- The checksum database (`sum.golang.org`) acts as a notary — it publishes the hash on a public, append-only log, so you can verify that everyone sees the same source.

## Core idea

### `go.sum` file structure

Each line in `go.sum` has one of two formats:

```
<module> <version> h1:<hash>
<module> <version>/go.mod h1:<hash>
```

- The first form is the hash of the entire module source tree
- The second form is the hash of just the `go.mod` file
- The hash algorithm is SHA-256, encoded as base64 with an `h1:` prefix

Example:
```
rsc.io/quote v1.5.2 h1:3fEykkD9k7lYzXq5aQAYN4sQFYqMvB8tmr5+TSowIc=
rsc.io/quote v1.5.2/go.mod h1:2KgnzaLcBqiXEg+BYhBDx+Hrm0oVk3GqS0Q2OQAvOKs=
```

### Integrity verification

When building:
1. The compiler resolves a module version from `go.mod`.
2. It downloads the module and computes SHA-256 hashes of the `go.mod` and the source tree.
3. It checks these against the entries in `go.sum`.
4. If no entry exists, it queries the checksum database (`sum.golang.org`) to get the authoritative hash.
5. It records the hash in `go.sum` for future verification.
6. If any hash does not match, the build fails.

### `go mod verify`

Verifies that the cached module files match `go.sum`:

```bash
go mod verify
```

This checks every module in the cache against the recorded hashes. If any file has been modified or corrupted, it reports the mismatch.

### Checksum database and privacy

- By default, `go get` queries `sum.golang.org` for checksums of public modules.
- This leaks your dependency list to Google.
- For private modules, set `GONOSUMDB` or `GOPRIVATE` to avoid sending their paths.
- The checksum database uses a **transparent log** (Trillian) — all entries are public and append-only, preventing undetectable tampering.

## Under the hood

The Go team operates `sum.golang.org`, a Go Checksum Database. It implements the `checksum` protocol defined in the Go modules specification. The database returns signed tree heads, and clients verify inclusion proofs. This means:

- Even if the database is compromised, existing entries cannot be retroactively modified.
- Clients can detect if the database shows different results to different requesters (by gossiping tree heads).

The `go.sum` file is append-only: when a module version is removed from `go.mod`, the corresponding `go.sum` entries are NOT automatically removed. Only `go mod tidy` removes unused entries. This is intentional — you might want to keep entries for modules that were previously used but are now removed, to verify that they were authentic at the time.

## How Go uses it

- **Every `go build`**: Verifies downloaded modules match `go.sum`.
- **`go mod download`**: Downloads modules and records hashes in `go.sum`.
- **`go mod tidy`**: Adds missing `go.sum` entries, removes stale ones.
- **`go mod verify`**: Checks cached module integrity against `go.sum`.
- **`go get -u`**: Updates dependencies and updates `go.sum`.
- **`GONOSUMCHECK`**: Environment variable to skip hash verification for specific modules. Use for private modules.
- **`GONOSUMDB`**: Environment variable to skip checksum database lookup for specific modules.
- **`GOFLAGS=-mod=mod`**: Forces `go build` to update `go.mod` and `go.sum` for missing modules.

## Go example

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Demonstrate go.sum verification by running go mod verify.
	cmd := exec.Command("go", "mod", "verify")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("go mod verify: %v\n%s", err, out)
		return
	}
	fmt.Printf("All modules verified:\n%s", out)
}
```

But the real "code" for this lesson is not a Go source file — it is the process of inspecting `go.sum`:

```bash
# View the go.sum file
cat go.sum

# Verify cached modules
go mod verify

# Check which modules have checksum database entries
GONOSUMDB='' go list -m -json rsc.io/quote
```

## Step-by-step execution

Inspecting `go.sum` integrity:

1. `go mod tidy` — downloads dependencies and populates `go.sum`.
2. `cat go.sum` — see the hashes. Each module version has two lines (source tree and go.mod).
3. `go mod verify` — confirms cached files match the hashes. Output: "all modules verified."
4. Deliberately corrupt a module file in the cache and re-run `go mod verify` — it reports the mismatch.
5. `GONOSUMCHECK=* go build` — builds skip hash verification (not recommended for production).

To observe the checksum database in action:

1. Clear the module cache: `go clean -modcache`
2. `GONOSUMDB='' go get rsc.io/quote@v1.5.2` — queries the checksum database.
3. `GONOSUMDB=rsc.io/quote go get rsc.io/quote@v1.5.2` — skips the database (builds with a warning if no `go.sum` entry).

## Common mistakes

- **Deleting `go.sum`**: This breaks reproducibility. Every developer who checks out your code will get different hash expectations if `go.sum` is missing. Never delete it.
- **Not committing `go.sum`**: Keeping `go.sum` out of version control defeats its purpose. It must be shared across all developers and CI environments.
- **Adding `go.sum` to `.gitignore`**: This is a common mistake. `go.sum` must be committed.
- **Assuming `go.sum` is a lock file**: `go.sum` is not a lock file — it does not pin transitive versions. It only records hashes for verification. The `go.mod` file determines versions; `go.sum` confirms authenticity.
- **Using `GONOSUMCHECK=*` in production**: This disables all hash verification, making you vulnerable to supply-chain attacks. Only use it for private modules.
- **Merging `go.sum` conflicts incorrectly**: `go.sum` is append-only and usually merges cleanly. If there is a merge conflict, accept both entries and run `go mod tidy` to clean up.

## Debugging walkthrough

Symptom: Build fails with:

```
go: .../module@v1.2.3: verifying module: checksum mismatch
    downloaded: h1:abc...
    go.sum:     h1:xyz...
```

**Root cause**: The module source you downloaded does not match the hash in `go.sum`. This can happen if:
- A bad `replace` directive points to a different version of the module.
- The module cache is corrupted (run `go clean -modcache`).
- Someone tampered with `go.sum` or the module source.
- A proxy served different code than the original module.

**Investigation**:
```bash
go clean -modcache
go mod tidy  # re-download and re-verify
```

If the error persists, the module source has genuinely changed (unlikely for a tagged version) or `go.sum` has been tampered with.

**Fix**: Determine the correct hash. For a trusted module version, delete the module's cache entry and the incorrect `go.sum` line, then run `go mod tidy`. For private modules, ensure `GONOSUMDB` and `GONOSUMCHECK` are set appropriately.

Symptom: `go mod verify` reports "mismatched hash" for a module you did not modify.

**Root cause**: The module cache file was corrupted (disk error, interrupted download, antivirus quarantine).

**Fix**:
```bash
go clean -modcache
go mod download
```

## Production notes

- **Always commit `go.sum`**: This is non-negotiable for reproducible builds. Make it part of your CI linting: check that `go.sum` is not modified by `go mod tidy`.
- **Private modules**: Set `GOPRIVATE=*.corp.com,*.internal.com` to bypass the checksum database and hash verification for internal modules. This is the standard practice for enterprise Go development.
- **Module proxies**: If your team uses a private module proxy (like Athens, Artifactory, or GoCenter), the proxy handles hash verification for you. Your `go.sum` still verifies the proxy served the correct code.
- **Air-gapped environments**: Pre-populate the module cache and `go.sum` on a networked machine, then copy both to the air-gapped build server.
- **CI/CD**: Run `go mod verify` in CI to catch corrupted cache or tampered dependencies early.

## Performance implications

- Hash computation is fast (~50-100ms per module version on modern hardware) and cached after first download.
- `go mod verify` scans the entire module cache — for large projects with hundreds of modules, this takes 1-5 seconds.
- The checksum database lookup adds latency to the first `go get` for a new module version (typically 200-500ms).
- Hashes in `go.sum` are not used at runtime — they are a build-time verification mechanism with zero runtime cost.

## Practice task

1. Create a new module:`mkdir /tmp/gosum-demo && cd /tmp/gosum-demo && go mod init example.com/gosum-demo`.
2. Add a dependency: write `main.go` importing `rsc.io/quote` and run `go mod tidy`.
3. Inspect `go.sum` — identify the two lines per module version.
4. Run `go mod verify` and confirm "all modules verified."
5. Set `GONOSUMCHECK=*` and run `go build`. Observe the warning (or lack thereof).
6. Clean the module cache and re-download: `go clean -modcache && go mod download`.
7. Re-run `go mod verify`.

## Tests / verification

```bash
go run ./curriculum/modules/05-types-interfaces-packages-modules/lessons/20-go-sum
go test ./curriculum/modules/05-types-interfaces-packages-modules/lessons/20-go-sum
```

## Review questions

1. What hash algorithm does `go.sum` use to fingerprint module versions?
2. What two pieces of content are hashed per module version?
3. What is the purpose of the checksum database (`sum.golang.org`)?
4. Why can you not simply delete `go.sum` and regenerate it?
5. What environment variable controls whether a module's hash is verified?

## NEXT UP

Dependency management — how to add, upgrade, downgrade, and prune Go dependencies using `go get` and related tooling.
