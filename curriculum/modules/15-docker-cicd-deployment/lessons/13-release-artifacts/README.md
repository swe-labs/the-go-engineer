# Release artifacts

## Learning objective

Version, build, checksum, and verify Go binaries as release artifacts, using ldflags for version injection and SHA256 checksums for integrity verification.

## Why this matters

A Go binary is the deliverable of a CI/CD pipeline. Without versioning, you cannot tell which commit produced a given binary in production. Without checksums, you cannot verify that a downloaded artifact has not been tampered with. Professional deployments rely on versioned, checksummed artifacts stored in a release registry so that every running binary is traceable to a specific commit, build time, and build environment.

## Mental model

A release artifact is a snapshot of your compiled Go program, tagged with metadata that answers three questions:

1. **What version?** (git tag or semver)
2. **From what commit?** (git SHA)
3. **When was it built?** (timestamp)

Think of a release artifact like a medicine bottle: the label (version, commit, build time) tells you what is inside, and the tamper-evident seal (SHA256 checksum) tells you nobody opened it after it left the factory. Without these, you cannot trust the artifact in production.

The checksum file (`SHA256SUMS`) is the manifest that maps binary names to their hashes. Before deploying, the deployment tool recomputes the hash and compares it against the manifest. A mismatch means the artifact was corrupted or replaced.

## Core idea

Release artifacts for Go are self-contained, statically compiled binaries. The artifact pipeline:

| Step | Action | Purpose |
|---|---|---|
| Build | `go build -ldflags "-X main.version=$VERSION"` | Embed version into binary |
| Checksum | `sha256sum <binary>` | Generate integrity hash |
| Sign | `gpg --sign SHA256SUMS` | Prove authenticity |
| Upload | Upload to artifact store (GitHub Releases, S3, Artifactory) | Make available for download |
| Verify | Compute checksum on download and compare | Ensure integrity |

ldflags (linker flags) use `-X` to set string variables in the Go binary at link time. This lets you inject `version`, `commit`, and `buildTime` without modifying source code. The values are baked into the `.data` section of the binary and are readable with `go version -m <binary>`.

## Under the hood

When building with ldflags, the Go linker (`go tool link`) replaces the value of a package-level string variable at link time:

```
go build -ldflags "-X main.version=v1.2.3 -X main.commit=abc1234 -X main.buildTime=2026-06-25T12:00:00Z" -o myapp .
```

The linker can only set `string` variables (not `int`, `bool`, or `const`). The variables must be declared at the package level (`var Version string`), not inside functions.

The resulting binary contains these strings in its read-only data section. You can inspect them with:

```
go version -m myapp
```

This prints all ldflags-injected values along with Go version and module information.

Checksums use SHA256, a cryptographic hash function producing a 32-byte (64 hex character) digest. The standard checksum file format (used by GNU coreutils) is:

```
<sha256 hash>  <binary name>
```

Each line contains the hex hash, two spaces, and the filename. Two spaces allow filenames containing single spaces to be parsed correctly.

## How Go uses it

Every Go release on GitHub follows this pattern. The official Go release workflow:

1. Tag the repository: `git tag v1.2.3 && git push --tags`.
2. Build for all target platforms using `GOOS=x GOARCH=y go build`.
3. For each binary, compute SHA256: `sha256sum myapp_v1.2.3_linux_amd64.tar.gz`.
4. Sign the checksum file with GPG: `gpg --detach-sign --armor SHA256SUMS`.
5. Create a GitHub Release with the tag, attach all binaries and `SHA256SUMS` and `SHA256SUMS.sig`.

Users verify downloaded artifacts:

```bash
sha256sum -c SHA256SUMS
gpg --verify SHA256SUMS.sig SHA256SUMS
```

## Go example

```go
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"runtime"
	"strings"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

type Artifact struct {
	Name     string
	Version  string
	GOOS     string
	GOARCH   string
	Checksum string
	Data     []byte
}

func NewArtifact(name, version, goos, goarch string, data []byte) *Artifact {
	return &Artifact{
		Name:    name,
		Version: version,
		GOOS:    goos,
		GOARCH:  goarch,
		Data:    data,
	}
}

func (a *Artifact) BinaryName() string {
	ext := ""
	if a.GOOS == "windows" {
		ext = ".exe"
	}
	if a.Name == "" {
		a.Name = "app"
	}
	return fmt.Sprintf("%s-%s-%s-%s%s", a.Name, a.Version, a.GOOS, a.GOARCH, ext)
}

func (a *Artifact) ComputeChecksum() string {
	h := sha256.Sum256(a.Data)
	return hex.EncodeToString(h[:])
}

func (a *Artifact) ChecksumLine() string {
	if a.Checksum == "" {
		a.Checksum = a.ComputeChecksum()
	}
	return fmt.Sprintf("%s  %s", a.Checksum, a.BinaryName())
}

type ChecksumEntry struct {
	Checksum   string
	BinaryName string
}

type ChecksumFile struct {
	Entries []ChecksumEntry
}

func ParseChecksums(content string) (*ChecksumFile, error) {
	cf := &ChecksumFile{}
	for _, line := range strings.Split(strings.TrimSpace(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid checksum line: %s", line)
		}
		cf.Entries = append(cf.Entries, ChecksumEntry{
			Checksum:   parts[0],
			BinaryName: strings.Join(parts[1:], " "),
		})
	}
	return cf, nil
}

func (cf *ChecksumFile) Verify(artifact *Artifact) bool {
	artifact.Checksum = artifact.ComputeChecksum()
	wantName := artifact.BinaryName()
	for _, entry := range cf.Entries {
		if entry.BinaryName == wantName && entry.Checksum == artifact.Checksum {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Release artifacts")
	fmt.Printf("Version: %s, Commit: %s, Built: %s\n", Version, Commit, BuildTime)

	app := NewArtifact("myapp", "v1.2.3", runtime.GOOS, runtime.GOARCH, []byte("binary-content"))
	app.Checksum = app.ComputeChecksum()
	fmt.Printf("Binary: %s\n", app.BinaryName())
	fmt.Printf("SHA256: %s\n", app.Checksum)

	cf := &ChecksumFile{
		Entries: []ChecksumEntry{
			{Checksum: app.Checksum, BinaryName: app.BinaryName()},
		},
	}
	if cf.Verify(app) {
		fmt.Println("Integrity: OK")
	}
}
```

## Step-by-step execution

For building and verifying a release artifact:

1. Developer tags commit: `git tag v1.2.3`.
2. CI pipeline detects the tag and runs `go build -ldflags "-X main.version=v1.2.3 -X main.commit=$(git rev-parse --short HEAD) -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"`.
3. The linker overwrites the default values of `Version`, `Commit`, `BuildTime` in the compiled binary.
4. The pipeline generates SHA256 checksum: `sha256sum myapp-v1.2.3-linux-amd64 >> SHA256SUMS`.
5. The artifact and checksum file are uploaded to GitHub Releases.
6. On deployment, the tool downloads `SHA256SUMS` and the binary, recomputes the SHA256, and compares.
7. If hashes match, the binary is deployed. If not, deployment is aborted.

For the `ParseChecksums` function:

1. Input: `"abc...  myapp-v1.2.3-linux-amd64\n"`.
2. Split by newline, trim whitespace.
3. Split line by whitespace: `["abc...", "myapp-v1.2.3-linux-amd64"]`.
4. Create `ChecksumEntry` with the hash and filename.
5. Append to `ChecksumFile.Entries`.

## Common mistakes

- **Using `const` instead of `var` for ldflags**: The linker cannot set constants. ldflags only work with mutable `var` declarations.
- **Wrong package path in ldflags**: The `-X` flag requires the full package path: `-X main.version=x` works when `version` is in package `main`. For other packages, use `-X github.com/user/repo/pkg.version=x`.
- **Not including `.exe` extension for Windows**: Windows users expect `.exe` on binaries. Without it, Windows may not recognize the file as executable.
- **Checksum mismatch from extra whitespace**: Checksum files use two spaces between hash and filename. Using a tab or one space breaks `sha256sum -c`.
- **Not signing the checksum file**: A checksum alone prevents accidental corruption but not malicious tampering. An attacker who replaces the binary can also replace the checksum file. Signing with GPG prevents this.

## Debugging walkthrough

Consider a deployment that fails with:

```
sha256sum: WARNING: 1 computed checksum did NOT match
```

**Symptom**: The production deployment script reports a checksum mismatch and refuses to deploy.

**Investigation**:
1. Get the expected checksum from the release artifacts' `SHA256SUMS` file.
2. Compute the checksum of the downloaded binary: `sha256sum myapp-v1.2.3-linux-amd64`.
3. Compare manually. If they differ, the binary was corrupted during download (network issue, storage corruption).

**Root cause**: The CI pipeline uploaded a new binary without updating `SHA256SUMS`, or the download was interrupted.

**Fix**: Re-download the artifact from the release page. If the checksum still fails, re-run the CI pipeline to regenerate both the binary and checksum.

Another scenario:

```
$ go version -m ./myapp
./myapp: go1.22.0
...
    build	-main.version=v1.2.3
    build	-main.commit=abc1234
```

If `version` shows "dev" instead of "v1.2.3", the CI pipeline did not pass the ldflags correctly. Check the CI script for correct quoting and variable substitution.

## Production notes

In production release pipelines:

- **Semantic versioning**: Use `vMAJOR.MINOR.PATCH` tags. Auto-increment patch for bug fixes, minor for features, major for breaking changes.
- **Immutable release tags**: Never move a git tag after publishing a release. A moved tag can silently deploy a different binary to production.
- **Multi-architecture builds**: For each release, build for `linux/amd64`, `linux/arm64`, and `darwin/amd64` at minimum. Use a build matrix in CI.
- **SBOM generation**: Generate a Software Bill of Materials (SBOM) for each release using `go-sbom` or `syft`. Attach SBOM to the release for vulnerability scanning.
- **Artifact retention**: Keep all release artifacts indefinitely for auditability. GitHub Releases stores artifacts forever; for self-managed storage, implement a retention policy.

## Performance implications

- **ldflags do not affect runtime performance**: The linker simply writes string literals into the binary's data section. There is zero runtime overhead.
- **Binary size impact**: Including version, commit, and build time adds approximately 100-200 bytes to the binary. Negligible.
- **Checksum computation**: SHA256 of a 50MB binary takes under 100ms on modern hardware. This runs once during the release pipeline and optionally on deploy.
- **Cross-compilation**: Building for multiple platforms multiplies build time by the number of targets. Use parallel CI jobs to amortize this.

## Practice task

Write a function `GenerateReleaseArtifacts(name, version string, platforms [][2]string) ([]*Artifact, error)` that:

- Takes a list of `[goos, goarch]` pairs.
- For each platform, creates an `Artifact` with mock data (e.g., `[]byte(fmt.Sprintf("%s-%s-%s-%s", name, version, goos, goarch))`).
- Computes the SHA256 checksum for each artifact.
- Returns the list of artifacts.
- Then generate a `SHA256SUMS` formatted string from all artifacts.

Write a `main()` that builds artifacts for `linux/amd64`, `linux/arm64`, `windows/amd64`, and `darwin/amd64`, and prints the checksum manifest.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/13-release-artifacts
go test ./curriculum/modules/15-docker-cicd-deployment/lessons/13-release-artifacts
```

The existing tests verify artifact creation, binary naming (including .exe), checksum computation, checksum file parsing, and integrity verification of valid and tampered artifacts. After completing the practice task, add tests for `GenerateReleaseArtifacts` covering multiple platforms and empty platforms.

## Review questions

1. What Go types can ldflags `-X` set? Why can it not set constants?
2. Why does a checksum file need to be signed (e.g., with GPG) rather than just uploaded alongside the binary?
3. What is the format of a standard `SHA256SUMS` file? Why does the format use two spaces between the hash and filename?
4. How would you verify that a Go binary deployed in production was built from a specific git commit?
5. What happens if a git tag is moved after a release is published? Why is this dangerous?

## NEXT UP

Config in deployment: managing environment variables, config files, and the 12-factor app approach for Go services.
