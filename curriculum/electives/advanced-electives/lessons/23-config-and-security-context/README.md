# Configuration, secrets, and security context

## Learning objective

By the end of this lesson, you will be able to inject configuration and secrets into a Go service with ConfigMaps and Secrets, load them with a fail-fast validator, and constrain the container with a non-root, read-only security context.

## Why this matters

A twelve-factor service reads its configuration from the environment, not from files baked into the image. Kubernetes gives you two objects for that: a ConfigMap for non-sensitive values and a Secret for sensitive ones. The distinction is operational, not cryptographic: a Secret is base64-encoded, not encrypted, unless the cluster is configured for encryption at rest. Treating a Secret as automatically safe is a common and expensive mistake.

The second half is the security context. Running as root inside a container is the default, and root in a container is root in a user namespace on the host. A container escape or a writable root filesystem turns a small bug into a host compromise. A restrictive security context is not a hardening afterthought; it is part of the deployment contract, and it must be set before the service reaches production.

## Mental model

Configuration is delivered in two shapes, and the shape decides the update behavior.

| Shape | Delivery | Update behavior | Fits |
|---|---|---|---|
| Environment variable | Rendered at container start | Fixed until the Pod restarts | Small, stable, non-sensitive values |
| Mounted file | Projected volume into the filesystem | The kubelet updates the mounted file in place | Rotatable values, certificates, large config |

The analogy is a sealed envelope versus a mailbox. An environment variable is a sealed envelope: the value is fixed when the container starts and cannot change until a new container is created. A mounted ConfigMap or Secret is a mailbox: the kubelet overwrites the file when the source changes, so a process that re-reads the file sees the new value without a restart. The service's loader must match the delivery: a value read once at startup cannot be rotated in place, no matter how the volume updates.

## Visual model

```mermaid
flowchart TD
    CM["ConfigMap<br/>(non-sensitive)"] --> ENV["Env vars<br/>(fixed at start)"]
    SEC["Secret<br/>(sensitive)"] --> FILE["Mounted file<br/>(rotates in place)"]
    ENV --> LOADER["LoadConfig<br/>(fail fast)"]
    FILE --> LOADER
    LOADER -->|"valid"| RUN["Serve"]
    LOADER -->|"invalid"| EXIT["Exit 1<br/>with a clear error"]

    classDef config fill:#1a233a,stroke:#4c8bf5,stroke-width:1px,color:#ffffff;
    classDef secret fill:#2a1e3a,stroke:#a371f7,stroke-width:1px,color:#ffffff;
    classDef run fill:#1e3a1e,stroke:#3fb950,stroke-width:1px,color:#ffffff;
    classDef fail fill:#3a1e1e,stroke:#f85149,stroke-width:1px,color:#ffffff;
    class CM,ENV config;
    class SEC,FILE secret;
    class LOADER,RUN run;
    class EXIT fail;
```

Caption: configuration and secrets enter through different paths, and the loader validates both before the service accepts traffic. An invalid value exits the process instead of serving a misconfigured service.

## Core idea

The delivery rules:

- **ConfigMap** holds non-sensitive key-value data. It can be projected as environment variables or as files in a volume.
- **Secret** holds sensitive data. It is base64-encoded in the API, stored as `Opaque` by default, and projected as environment variables or files. Use file mounts for anything you intend to rotate.
- **Environment variables are immutable for the container's lifetime.** Changing the ConfigMap or Secret does not change the running container's environment. Environment-sourced secrets also leak into crash dumps, child processes, and `kubectl describe`.
- **File mounts update in place.** The kubelet syncs projected volumes periodically (up to about a minute). The application must re-read the file to benefit. Reading once at startup gives you validation but not rotation.
- **Fail fast.** Validate every required value at startup and exit non-zero with a clear message. A service that starts with a missing database URL and fails on the first request is harder to diagnose than one that refuses to start.

The security context rules:

| Setting | Effect |
|---|---|
| `runAsNonRoot: true` | The kubelet refuses to start the container as root |
| `runAsUser: 10001` | Runs the process as an unprivileged UID |
| `readOnlyRootFilesystem: true` | The container filesystem is immutable; temp files need an `emptyDir` |
| `allowPrivilegeEscalation: false` | Blocks setuid binaries from gaining more privilege |
| `capabilities.drop: ["ALL"]` | Removes all Linux capabilities from the process |
| `seccompProfile.type: RuntimeDefault` | Applies the runtime's default syscall filter |
| `automountServiceAccountToken: false` | No API token is mounted unless the Pod needs one |

## Go example

```go
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config is the validated runtime configuration for the service.
type Config struct {
	Addr        string
	LogLevel    string
	MaxConns    int
	DatabaseURL string
}

// LoadConfig builds a Config from environment variables and a file-mounted
// secret. It returns the first validation error it finds so the process can
// fail fast at startup instead of serving with missing or invalid settings.
func LoadConfig(getenv func(string) string, readFile func(string) ([]byte, error)) (Config, error) {
	c := Config{
		Addr:     getenv("APP_ADDR"),
		LogLevel: getenv("APP_LOG_LEVEL"),
	}
	if c.Addr == "" {
		c.Addr = ":8080"
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return Config{}, fmt.Errorf("APP_LOG_LEVEL %q is not one of debug, info, warn, error", c.LogLevel)
	}

	raw := getenv("APP_MAX_CONNS")
	if raw == "" {
		return Config{}, errors.New("APP_MAX_CONNS is required")
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return Config{}, fmt.Errorf("APP_MAX_CONNS %q must be a positive integer", raw)
	}
	c.MaxConns = n

	path := getenv("APP_DATABASE_URL_FILE")
	if path == "" {
		return Config{}, errors.New("APP_DATABASE_URL_FILE is required")
	}
	data, err := readFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read secret %s: %w", path, err)
	}
	c.DatabaseURL = strings.TrimSpace(string(data))
	if c.DatabaseURL == "" {
		return Config{}, fmt.Errorf("secret %s is empty", path)
	}
	return c, nil
}

func main() {
	cfg, err := LoadConfig(os.Getenv, os.ReadFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("listening on %s log=%s max_conns=%d\n", cfg.Addr, cfg.LogLevel, cfg.MaxConns)
}
```

The matching ConfigMap, Secret, and Deployment:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: orders-api-config
data:
  APP_ADDR: ":8080"
  APP_LOG_LEVEL: "info"
  APP_MAX_CONNS: "25"
---
apiVersion: v1
kind: Secret
metadata:
  name: orders-api-secrets
type: Opaque
stringData:
  DATABASE_URL: "postgres://orders:REPLACE_ME@postgres:5432/orders"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: orders-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: orders-api
  template:
    metadata:
      labels:
        app: orders-api
    spec:
      automountServiceAccountToken: false
      securityContext:
        runAsNonRoot: true
        runAsUser: 10001
        runAsGroup: 10001
        seccompProfile:
          type: RuntimeDefault
      containers:
        - name: api
          image: registry.example.com/orders-api:v1.4.0
          envFrom:
            - configMapRef:
                name: orders-api-config
          env:
            - name: APP_DATABASE_URL_FILE
              value: /etc/orders/db-url
          securityContext:
            readOnlyRootFilesystem: true
            allowPrivilegeEscalation: false
            capabilities:
              drop: ["ALL"]
          volumeMounts:
            - name: secrets
              mountPath: /etc/orders
              readOnly: true
            - name: tmp
              mountPath: /tmp
      volumes:
        - name: secrets
          secret:
            secretName: orders-api-secrets
            items:
              - key: DATABASE_URL
                path: db-url
        - name: tmp
          emptyDir: {}
```

## Step-by-step execution

Predict before you run: with `readOnlyRootFilesystem: true`, the process writes a log file to `/var/log/orders.log`. What happens, and where should the file go instead? Write your answer before reading the Debugging walkthrough.

1. `LoadConfig` reads `APP_ADDR` and `APP_LOG_LEVEL` from the environment. Both have defaults, so the service starts without them, but `APP_LOG_LEVEL` must be one of the four allowed values.
2. It reads `APP_MAX_CONNS`, which has no default. A missing or non-numeric value returns an error.
3. It reads the path in `APP_DATABASE_URL_FILE`, then reads that file. The path points into the mounted Secret volume, so the value never appears in the process environment.
4. The Secret value is trimmed of its trailing newline, checked for emptiness, and stored.
5. If every check passes, `main` prints the resolved configuration and starts serving. If any check fails, `main` prints the error and exits 1.
6. When the Secret is updated, the kubelet rewrites `/etc/orders/db-url`. A loader that re-reads the file picks up the new value; this one reads at startup, so rotation requires either a re-read or a restart.

## Runtime evidence

Run the tests with the race detector:

```bash
go test -race -count=1 .
```

Expected output (from `curriculum/electives/advanced-electives/lessons/23-config-and-security-context`):

```text
ok  	github.com/swe-labs/the-go-engineer/curriculum/electives/advanced-electives/lessons/23-config-and-security-context	1.7s
```

The tests cover the valid path, the defaults, a missing required value, a non-numeric value, an invalid log level, and secret rotation. Now demonstrate fail-fast behavior by running the binary with no configuration:

```bash
go run .
```

Expected output:

```text
config error: APP_MAX_CONNS is required
exit status 1
```

Interpretation: the process refuses to start instead of accepting traffic it cannot serve. In the cluster this surfaces as a `CrashLoopBackOff` with a clear message in `kubectl logs`, which is far cheaper to diagnose than a runtime failure on the first request.

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Baking secrets into the image | It is the simplest path in development | Use a Secret and mount it as a file; never commit real values |
| Reading a Secret from an environment variable | It is one line of code | Prefer a file mount; env values leak into dumps and child processes |
| Assuming a Secret is encrypted at rest | "Secret" implies safety | Confirm etcd encryption; use Sealed Secrets, SOPS, or a cloud manager |
| Running as root | It is the container default | Set `runAsNonRoot` and `runAsUser` |
| Enabling `readOnlyRootFilesystem` then writing to `/` | Logs and temp files often assume a writable root | Mount an `emptyDir` at `/tmp` and send logs to stdout |
| Reading config once and expecting rotation | Startup reading is simpler | Re-read the mounted file on an interval, or restart on change |
| Leaving the service account token mounted | It is mounted by default | Set `automountServiceAccountToken: false` unless the Pod needs the API |

## Debugging walkthrough

**Symptom:** after adding `readOnlyRootFilesystem: true`, the Pod starts and then exits with `read-only file system` in the logs. Before the change, the same image ran fine.

**Hypotheses, cheapest first:**

1. The image entrypoint writes to a path that is now read-only.
2. The Secret mount is misconfigured and the loader fails before serving.
3. The process needs a writable directory that does not exist under the new filesystem.

**Tool:** `kubectl logs orders-api-<hash>` shows the exact path in the error; `kubectl describe pod` shows the mount table and the exit code.

**Root cause:** the service opens `/var/log/orders.log` for append. With a read-only root filesystem, that open fails and the process exits. The volume mounts only cover `/etc/orders` and `/tmp`.

**Fix:** write logs to stdout and let the cluster collect them, or mount a dedicated `emptyDir` at the log directory. Sending logs to stdout is the correct fix: it removes the file dependency and matches the twelve-factor model.

**Prevention:** run the container locally with `--read-only` before deploying the change. A local `docker run --read-only --tmpfs /tmp image` reproduces the failure in seconds.

## Production notes

- **Fail fast at startup.** Validate required configuration and exit non-zero on the first error. A Pod that crashes with a clear message triggers `CrashLoopBackOff` and surfaces the reason in `kubectl logs`.
- **Separate ConfigMap from Secret.** Non-sensitive values belong in a ConfigMap so they are visible in review and diff cleanly. Put only sensitive values in a Secret.
- **Prefer file mounts for rotation.** A projected Secret volume updates in place. A service that periodically re-reads the file rotates without a restart.
- **Set the full security context.** `runAsNonRoot`, `runAsUser`, `readOnlyRootFilesystem`, `allowPrivilegeEscalation: false`, `capabilities.drop: ["ALL"]`, and a seccomp profile together.
- **Do not log secrets.** Redact connection strings and tokens. A Secret in a log line is a secret in your log aggregation system forever.
- **Use `stringData` for authoring, `data` for storage.** `kubectl apply` converts `stringData` to base64 in `data`. Never paste base64 by hand.
- **Manage secrets outside the repo.** External Secrets Operator, Sealed Secrets, SOPS, or a cloud secret manager keeps plaintext out of Git.

## Practice task

Harden the following in-memory Deployment and adapt the loader:

1. The container currently runs as root with a writable root filesystem and no capability drop. Add the full security context from this lesson.
2. The `DATABASE_URL` is currently an environment variable from a Secret. Move it to a file mount and update `LoadConfig` if needed.
3. Add an `APP_FEATURE_TIMEOUT` value read from the ConfigMap with a default of `5s`, validated with `time.ParseDuration`.
4. Write a table-driven test that covers a valid timeout, an empty timeout, and an invalid timeout.
5. Explain in one paragraph why the readiness probe should not fail when the database is briefly unreachable, and which probe should reflect that instead.

## Tests / verification

```bash
go run ./curriculum/electives/advanced-electives/lessons/23-config-and-security-context
go test -race -count=1 ./curriculum/electives/advanced-electives/lessons/23-config-and-security-context
```

The tests verify:

- A fully specified environment with a mounted secret loads correctly.
- Missing address and log level fall back to defaults.
- A missing `APP_MAX_CONNS` returns an error.
- A non-numeric `APP_MAX_CONNS` returns an error.
- An invalid log level returns an error.
- Re-reading a changed secret returns the new value.

## Review questions

1. What is the practical difference between a ConfigMap and a Secret?
2. Why does changing a ConfigMap not change a running container's environment variables?
3. When does a mounted Secret file change on disk, and what must the application do to observe it?
4. What does `readOnlyRootFilesystem: true` break, and how do you fix it without disabling the setting?
5. Why should the liveness probe not check the database, and which probe should?

## NEXT UP

[Probes, lifecycle, and graceful shutdown](../24-probes-and-graceful-shutdown/README.md)
