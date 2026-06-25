# File permissions

## Learning objective

Read and set Unix file permissions using Go's `os.FileMode`, understand octal permission notation, and avoid common permission pitfalls when writing files in Docker containers.

## Why this matters

When your Go service runs as PID 1 in a Docker container, it creates files — logs, configs, uploads, database files. If the permissions are wrong, other processes (sidecars, init containers, backup jobs) cannot read them. If they are too permissive (0777), secrets are exposed. Understanding file permissions is critical for secure containerized applications and for writing CI/CD tooling that manipulates file systems.

## Mental model

Each file has three permission rings: owner, group, and other. Each ring has three bits: read (r=4), write (w=2), execute (x=1). Octal notation encodes all three rings in one number: owner-group-other. Mode `0644` means owner can read+write (6), group can read (4), others can read (4). Go's `os.FileMode` is the language's direct port of this Unix model.

## Core idea

`os.FileMode` is a `uint32` that mirrors the Unix permission bitmap. The lower 12 bits encode:

| Bits | Field | Octal |
|---|---|---|
| 0-2 | Other rwx | 0007 |
| 3-5 | Group rwx | 0070 |
| 6-8 | Owner rwx | 0700 |
| 9 | Sticky bit | 1000 |
| 10 | Setgid | 2000 |
| 11 | Setuid | 4000 |

Common modes: `0644` (regular file), `0755` (executable), `0600` (private), `0777` (wide open, almost never correct).

## Under the hood

When `os.Stat` returns a `FileInfo`, the `Mode()` method returns an `os.FileMode`. The `os.Chmod` syscall on Linux calls `syscall.Fchmodat` or `syscall.Chmod`, which updates the inode permission bits in the filesystem. On Windows, Go translates permission bits to Windows ACLs but the mapping is approximate: regular files get 0666 and directories get 0777 by default. Docker layers each add their own overlay metadata layer — when you `RUN chmod` in a Dockerfile, the overlay filesystem records the new mode without rewriting the underlying file data.

## How Go uses it

Every file-related operation in Go involves permissions:
- `os.WriteFile(name, data, perm)` — creates files with the given mode (subject to umask)
- `os.OpenFile(name, flag, perm)` — opens with flags, optionally creates with perm
- `os.MkdirAll(name, perm)` — creates directories
- `os.Chmod(name, mode)` — changes permissions on existing files
- `os.Stat(name).Mode()` — reads current permissions
- `os.CreateTemp(dir, pattern)` — creates temporary files with 0600

## Go example

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.CreateTemp("", "perms-*.txt")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	info, _ := os.Stat(file.Name())
	fmt.Printf("Default mode: %#o\n", info.Mode().Perm())

	os.Chmod(file.Name(), 0600)
	info, _ = os.Stat(file.Name())
	fmt.Printf("After chmod 0600: %#o\n", info.Mode().Perm())

	os.Chmod(file.Name(), 0755)
	info, _ = os.Stat(file.Name())
	fmt.Printf("After chmod 0755: %#o\n", info.Mode().Perm())

	fmt.Printf("IsRegular: %v\n", info.Mode().IsRegular())
	fmt.Printf("Perm string: %s\n", info.Mode().Perm())
}
```

## Step-by-step execution

1. `os.CreateTemp` creates a temp file with default permissions (on Unix, 0600 after umask).
2. `os.Stat` reads the file's metadata from the filesystem and populates a `FileInfo`.
3. `info.Mode()` returns an `os.FileMode` value.
4. `info.Mode().Perm()` extracts the lower 12 bits (the permission octal).
5. `os.Chmod("/tmp/test.txt", 0600)` calls `syscall.Chmod` on Unix, which updates the inode.
6. After chmod, the file's permission bits are `-rw-------` (owner read+write).
7. `%#o` prints the mode as a zero-prefixed octal number.

## Common mistakes

- **Using 0777 as a default to "make things work".** This bypasses all permission controls and creates security vulnerabilities. Always set the most restrictive mode that still works.
- **Confusing octal and decimal.** Writing `os.WriteFile(name, data, 644)` instead of `os.WriteFile(name, data, 0644)`. Go interprets `644` as decimal (`0644` octal = `420` decimal). Always use the `0` prefix for octal.
- **Running containers as root.** When the Dockerfile does not specify `USER`, the container runs as root. Files created by the Go process are owned by root and inaccessible to non-root CI/CD or host processes.
- **Assuming os.WriteFile defaults to secure permissions.** It uses `0666` (before umask). On systems with umask 022, the result is `0644` — world-readable. For secrets, explicitly pass `0600`.
- **Forgetting umask.** `os.WriteFile("f", data, 0666)` creates a file, but the kernel applies the process umask. If umask is 027, the result is `0640`, not `0666`. Call `syscall.Umask(0)` to check or set umask.

## Debugging walkthrough

Consider a Go service that writes an API key file:

```go
os.WriteFile("/secrets/api.key", []byte(key), 0644)
```

**Symptom:** The API key file is world-readable in the container. Other processes can read it.

**Investigation:** Run `ls -la /secrets/api.key` inside the container. The output shows `-rw-r--r--` (0644), confirming world-readable.

**Root cause:** The developer used 0644 (owner write, everyone read). The file contains a secret and should be owner-only.

**Fix:** Change to 0600:
```go
os.WriteFile("/secrets/api.key", []byte(key), 0600)
```
Also ensure the container runs as a non-root user so only that user can read the file.

## Production notes

- TLS private keys must be 0600. Go's `tls.LoadX509KeyPair` reads the file, and most orchestrators refuse to start if key files have group or world permissions.
- SSH keys, database passwords, and API tokens should always be stored with 0600.
- Docker secrets mounted into containers have permissions set by the orchestrator (typically 0444). Do not chmod them.
- In Kubernetes, use `securityContext.fsGroup` and `securityContext.runAsUser` to control file ownership.
- Use `gosec` (G301, G302, G306) to catch overly permissive FileMode arguments in CI.
- When using `os.OpenFile` with `os.O_CREATE`, always specify the permission mode explicitly. The zero value (0) creates unusable files.

## Performance implications

- Permission checks via `os.Stat` are filesystem metadata lookups — they do not read file content. Cost is typically microseconds.
- `os.Chmod` is a syscall that modifies inode metadata. On overlay filesystems in Docker, it triggers a copy-up of the metadata, which is fast (no data copy).
- Setting permissions on creation (via `os.WriteFile` or `os.OpenFile`) has zero additional cost — the kernel sets the bits atomically during inode creation.
- Avoid calling `os.Chmod` on the hot path; set the correct mode at file creation time.

## Practice task

Write a function `WriteSecureFile(path string, data []byte) error` that writes a file with 0600 permissions. Then write a function `ReadPermissions(path string) (os.FileMode, error)` that reads and returns the permissions of a file. Write a `main()` that creates a secure file, reads its permissions, and prints them in octal format.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/04-file-permissions
go test -v ./curriculum/modules/15-docker-cicd-deployment/lessons/04-file-permissions
```

## Review questions

1. What is the difference between `0644` and `0755`? When would you use each?
2. Why is `os.WriteFile(name, data, 644)` a bug (without the `0` prefix)?
3. What does the umask do, and why does it affect the effective permission of files created by `os.WriteFile`?
4. Why should Docker containers run as non-root, and how does this relate to file permissions?
5. What permission mode should you use for a TLS private key file?

## NEXT UP

Networking basics — TCP/UDP, ports, listening, dialing, and the difference between localhost and 0.0.0.0.
