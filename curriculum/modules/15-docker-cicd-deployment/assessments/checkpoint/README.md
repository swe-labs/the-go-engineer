# Module 15 Assessment — Containerization, CI/CD, and Deployment

## Purpose

Validate mastery for the target curriculum items.

## Targets

- `core-15-01`
- `core-15-02`
- `core-15-03`
- `core-15-04`
- `core-15-05`
- `core-15-06`
- `core-15-07`
- `core-15-08`
- `core-15-09`
- `core-15-10`
- `core-15-11`
- `core-15-12`
- `core-15-13`
- `core-15-14`
- `core-15-15`
- `core-15-16`
- `core-15-17`
- `core-15-18`
- `core-15-19`

## Evidence

Submit answers, code evidence, test output, and explanation where required.

---

## 1. Deep-Dive Cloud-Native & Deployment Scenarios

Provide rigorous, production-grade answers to each scenario, referencing Linux container specifications (OCI), cgroups v2, Docker engine internals, and orchestration platforms (Kubernetes).

### Scenario 1: Linux Namespaces, Cgroups, and Container Isolation (`core-15-01` through `core-15-04`)
A junior developer asks: "Is a Docker container a lightweight virtual machine?"
- Demystify this misconception. Explain the exact Linux kernel primitives that construct a container:
  - **Namespaces**: PID, NET, MNT, IPC, UTS, USER.
  - **Control Groups (cgroups v2)**: Memory limits, CPU quotas (`cpu.cfs_quota_us` vs `cpu.cfs_period_us`).
  - **OverlayFS (chroot/pivot_root)**.
- Why does the Go runtime prior to Go 1.19 allocate `GOMAXPROCS` based on the host CPU count rather than the container's cgroup CPU quota? How does `uber-go/automaxprocs` or modern Go runtime flags resolve CPU throttling?

### Scenario 2: Multi-Stage Builds, Security Hardening, and Base Images (`core-15-05` through `core-15-08`)
A security audit flags a production Go container image for running as UID 0 (root) and containing multiple critical CVEs originating from an unpatched Linux shell and package manager.
- Compare the trade-offs of runtime base images: `scratch`, `alpine`, and `distroless` (`gcr.io/distroless/static`).
- Why is `CGO_ENABLED=0` essential when targeting `scratch` or `alpine` from a standard Debian/Ubuntu build environment?
- Explain the security significance of `USER appuser:appuser` (UID 10001). What permissions are denied to a non-root user inside a Linux container?
- How do compiler flags `-trimpath` and `-ldflags="-s -w"` enhance security and reduce binary footprint?

### Scenario 3: Health Probes, Rolling Updates, and Zero-Downtime Deployments (`core-15-09` through `core-15-13`)
During a Kubernetes rolling deployment (`maxSurge: 25%`, `maxUnavailable: 0`), client API traffic experiences a sudden spike of HTTP 502 Bad Gateway and connection refused errors.
- Diagram the sequence of events during a Kubernetes Pod termination:
  1. Service controller removes pod endpoint from kube-proxy / iptables.
  2. Kubelet sends `SIGTERM` to the container process.
  3. Kubelet waits for `terminationGracePeriodSeconds`.
  4. Kubelet sends `SIGKILL` if process has not exited.
- Why does an immediate process exit on `SIGTERM` cause HTTP 502 errors? (Explain the race condition between iptables endpoint removal and pod shutdown).
- How does adding a 5-second sleep or readiness failure delay in the Go shutdown sequence guarantee zero-downtime traffic draining?

### Scenario 4: Automated CI/CD Pipelines & Continuous Delivery (`core-15-14` through `core-15-19`)
A software engineering team wants to implement an enterprise GitHub Actions CI/CD pipeline for their Go microservices.
- Design an end-to-end pipeline covering:
  - Commit validation: `golangci-lint`, formatting (`go fmt`), and static analysis (`govulncheck`).
  - Automated testing: Unit and integration tests with data race detection (`go test -race`).
  - Build & Caching: Multi-stage Docker build with GitHub Actions cache backend (`type=gha`).
  - Artifact signing and container registry publishing with Cosign.
- Contrast Canary Deployments, Blue-Green Deployments, and Rolling Updates. What are the observability requirements for automated rollbacks?

---

## 2. Practical Containerization Challenge

Write and verify the deployment files for a production Go microservice:
1. An optimized multi-stage `Dockerfile` producing a static non-root `scratch` image with CA certificates and timezone data.
2. A `docker-compose.yml` file declaring:
   - Resource limitations (0.5 CPUs, 128MB RAM).
   - Container healthcheck.
   - Non-root security context (`read_only: true`, `no-new-privileges: true`).
3. An HTTP server in Go with an internal `-healthcheck` probe flag and graceful shutdown with in-flight drain handling.

---

## 3. Invariant Deployment Rules

1. **Never Run as Root in Production**: Always specify `USER <non-root-uid>` in the production container stage.
2. **Deterministic Static Builds**: Always compile with `CGO_ENABLED=0 -trimpath -ldflags="-s -w"` for container deployments.
3. **Decouple Liveness from External Dependencies**: Never check third-party databases in liveness probes; use readiness probes for dependency health.
4. **Always Test with `-race` in CI**: Never deploy code to production without passing automated race detection in the CI pipeline.

---

## 4. 4-Tier Assessment Rubric

| Level | Criteria & Observable Evidence |
|---|---|
| **Level 1: Novice** | Ships compiler in production image; runs as root; dynamic CGO binary crashes on scratch; confuses liveness with readiness; no CI pipeline. |
| **Level 2: Advanced Beginner** | Uses multi-stage build, but breaks layer caching; copies source files before `go.mod`; shuts down immediately on SIGTERM dropping active connections. |
| **Level 3: Competent** | Produces minimal static container (~15MB); runs as non-root UID 10001; handles SIGTERM with graceful in-flight request draining; configures working CI/CD with race detection. |
| **Level 4: Proficient** | Exceptional cloud-native engineering. Explains cgroup CFS quotas and kernel namespaces; tunes container termination lifecycle to prevent network race conditions; implements zero-CVE distroless containers with automated vulnerability scanning. |

---

## 5. Submission Requirements

Submit your detailed answers for Scenarios 1–4 and the complete Dockerfile, docker-compose.yml, and Go source code for the practical containerization challenge.
