# Project Assessment — Dockerized Service with CI/CD

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `project-dockerized-service-cicd`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Assessment Overview & Industrial Context

Modern cloud platforms (Kubernetes, AWS ECS, Nomad) deploy services as lightweight, secure container images. Deploying raw Go binaries directly or relying on bloated single-stage container builds exposes severe operational problems:
- **Massive Image Footprints**: A single-stage container that includes the entire Go compiler, package toolchain, and build cache can exceed 800MB, increasing deployment time, network bandwidth consumption, and cluster autoscaling latency.
- **Security & Privilege Escalation**: Running container workloads as the default `root` user violates security compliance (e.g. SOC2, CIS benchmarks) and enables privilege escalation if a container escape vulnerability is discovered.
- **Vulnerability Surface Area**: Shipping package managers (`apt`, `apk`), shells (`/bin/sh`, `/bin/bash`), and development utilities inside production images grants attackers ready-made tools upon code injection.
- **Unreliable Orchestration Lifecycle**: Containers lacking orchestrator-ready liveness (`/healthz`) and readiness (`/readyz`) probes suffer dropped requests and 502 Bad Gateway errors during rolling rollouts and autoscaling events.

In this assessment, you will containerize a production-grade Go HTTP microservice. You will craft an optimized multi-stage `Dockerfile` producing a minimal (~15MB) static binary image running as an unprivileged user (`USER appuser`), configure local orchestration with `docker-compose.yml`, implement orchestrator probes and graceful draining, and build an automated GitHub Actions CI/CD pipeline.

---

## 2. Core Competencies Tested

- **Multi-Stage Dockerfile Optimization**: Separating the build environment (`golang:1.24-alpine`) from the minimal runtime environment (`scratch` or minimal `alpine`), leveraging layer caching for `go.mod` dependencies.
- **Static Compilation & Binary Hardening**: Compiling CGO-free static binaries using `CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w"`.
- **Container Security & Least Privilege**: Creating and enforcing an unprivileged non-root user (`USER appuser:appuser`, UID `10001`), bundling essential root CA certificates (`ca-certificates`), and eliminating shells from production images.
- **Orchestration Probes & Graceful Drain**: Implementing independent `/healthz` (liveness) and `/readyz` (readiness) HTTP endpoints, a CLI healthcheck command flag (`-healthcheck`), and signal-driven graceful shutdown (`SIGTERM`/`SIGINT`).
- **Production CI/CD Automation**: Building a GitHub Actions workflow that executes testing with race detection (`go test -race`), linting, formatting checks, and cached container image builds.

---

## 3. Zero-Magic Architectural Breakdown

```
┌──────────────────────────────────────────────────────────┐
│ STAGE 1: BUILDER (golang:1.24-alpine)                    │
│   1. Install ca-certificates & tzdata                    │
│   2. Create unprivileged appuser (UID 10001)             │
│   3. COPY go.mod go.sum -> RUN go mod download (CACHED)  │
│   4. COPY source code -> RUN CGO_ENABLED=0 go build      │
│      -trimpath -ldflags="-s -w" -o /bin/server           │
└────────────────────────────┬─────────────────────────────┘
                             │ Copies ONLY static binary,
                             │ ca-certificates, and /etc/passwd
                             ▼
┌──────────────────────────────────────────────────────────┐
│ STAGE 2: PRODUCTION RUNTIME (scratch or minimal alpine)   │
│   - Image Footprint: ~15 MB                              │
│   - Zero shell, zero compiler, zero package manager      │
│   - Security Context: USER appuser:appuser (UID 10001)   │
│   - Entrypoint: ["/bin/server"]                          │
└──────────────────────────────────────────────────────────┘
```

### Why Layer Caching Order Matters
Docker builds layers incrementally and reuses cached layers if preceding instructions and input files are unchanged. Copying `go.mod` and `go.sum` and running `go mod download` *before* copying application source code ensures that dependency downloads are cached as long as dependencies remain identical, reducing CI build times from minutes to seconds.

### The Significance of `CGO_ENABLED=0`
By default, Go dynamically links to `libc` if certain standard library features (such as system DNS resolution or user/group lookups) are invoked. In a `scratch` or `alpine` runtime container, `glibc` is absent, causing dynamically linked binaries to fail at launch with cryptic `file not found` errors. Forcing `CGO_ENABLED=0` produces a truly standalone, statically linked executable containing all necessary Go runtime routines.

---

## 4. Implementation Requirements & Task Breakdown

### Task 1: Multi-Stage `Dockerfile`
In `_starter/Dockerfile`:
- Builder Stage:
  - Base: `golang:1.24-alpine AS builder`.
  - Install `ca-certificates` and `tzdata`.
  - Create non-root user: `RUN adduser -D -g '' -u 10001 appuser`.
  - Copy `go.mod` and `go.sum`; execute `go mod download`.
  - Copy application source files.
  - Compile static binary: `RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/server .`.
- Runtime Stage:
  - Base: `scratch` or minimal `alpine:3.20`.
  - Copy `/etc/ssl/certs/ca-certificates.crt`.
  - Copy `/etc/passwd` to resolve user `appuser`.
  - Copy `/bin/server`.
  - Set `USER appuser:appuser`.
  - Declare `EXPOSE 8080`.
  - Set `ENTRYPOINT ["/bin/server"]`.

### Task 2: Probes, Healthcheck, and Graceful Shutdown
In `_starter/main.go`:
- Expose `/healthz` returning HTTP 200 while server is running.
- Expose `/readyz` returning HTTP 200 normally, but transitioning to HTTP 503 Service Unavailable when shutdown signal (`SIGTERM`) is received to drain traffic.
- Support CLI flag `-healthcheck`:
  - When invoked with `-healthcheck`, make an HTTP GET request to `http://localhost:8080/healthz`. Exit with code 0 on success, code 1 on failure.
- Listen for `syscall.SIGTERM` and `os.Interrupt`, trigger readiness failure, wait for in-flight requests with a 10-second timeout, and invoke `server.Shutdown(ctx)`.

### Task 3: Local Orchestration with `docker-compose.yml`
- Define service `web` building from the current context.
- Configure port mapping `8080:8080`.
- Specify resource limits: `cpus: '0.5'`, `memory: 128M`.
- Define healthcheck using the container's internal probe.

### Task 4: GitHub Actions CI/CD Pipeline
In `_starter/.github/workflows/ci.yml`:
- Trigger on `push` and `pull_request` against `main`.
- Setup Go environment with caching.
- Run `go test -v -race ./...`.
- Run container build step ensuring the multi-stage image compiles cleanly.

---

## 5. Verification & Acceptance Criteria

Run test verification suite:

```bash
go test -v ./curriculum/modules/15-docker-cicd-deployment/projects/dockerized-service-with-ci-cd/_solution
```

### Acceptance Checklist
- [ ] `TestDockerfile_MultiStageIntegrity`: verifies multi-stage build stages, `CGO_ENABLED=0`, non-root user `appuser`, and static linking flags.
- [ ] `TestServer_HealthAndReadinessEndpoints`: verifies `/healthz` returns 200, `/readyz` returns 200 and toggles to 503 on shutdown.
- [ ] `TestServer_HealthcheckCLIFlag`: confirms binary supports internal `-healthcheck` CLI flag for container environments lacking `curl` or `wget`.
- [ ] `TestDockerCompose_ResourceConstraints`: verifies CPU and memory limits are declared.
- [ ] `TestCIDeploymentWorkflow`: validates GitHub Actions workflow yaml structure and race detection testing.

---

## 6. Common Pitfalls & Anti-Patterns

1. **Missing CA Certificates in `scratch`**:
   *Anti-pattern*: Using `FROM scratch` without copying `/etc/ssl/certs/ca-certificates.crt`. Any outbound HTTPS call (e.g. third-party APIs, webhooks) will fail with `x509: certificate signed by unknown authority`.
   *Correction*: Always copy root CA bundles from the builder stage into the runtime container.

2. **Root User Execution**:
   *Anti-pattern*: Omitting `USER` directive, allowing the container process to execute as UID 0 (root).
   *Correction*: Create and switch to `USER appuser:appuser` (UID >= 10000).

3. **Treating Liveness Identical to Readiness**:
   *Anti-pattern*: Failing liveness probes when downstream dependencies (database, cache) are briefly saturated, causing Kubernetes to enter a restart loop.
   *Correction*: Liveness only monitors if the Go process is alive and un-deadlocked; readiness monitors if the service can currently accept customer traffic.

---

## 7. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Novice** | Uses a bloated single-stage image (>500MB); runs as root; binary fails to execute in scratch due to missing CGO flags; no health endpoints. |
| **Level 2: Advanced Beginner** | Multi-stage build created, but layer caching is broken by copying source before `go.mod`. Health endpoints exist but liveness and readiness are coupled. |
| **Level 3: Competent** | Minimal image (~15MB) running as non-root `appuser`. Static binary with `-trimpath -ldflags="-s -w"`. Independent `/healthz` and `/readyz` probes with graceful drain. Working CI workflow. |
| **Level 4: Proficient** | Exceptional container engineering. Incorporates CLI health checks for scratch images without shell. CA certificates and timezone bundles managed properly. Resource limits and graceful shutdown timeouts tuned for cloud orchestrators. |

---

## 8. Submission & Reflection

1. Submit the output from `docker build` showing stage transitions and final image size.
2. Submit test output from `go test -v ./...`.
3. Answer: *Why is `/bin/sh` or `/bin/bash` intentionally excluded from production runtime container images in high-security environments?*
