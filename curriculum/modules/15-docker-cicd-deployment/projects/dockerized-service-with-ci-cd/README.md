# Dockerized Service with CI/CD

## Learning objective

Containerize a production-grade Go HTTP microservice using an optimized multi-stage `Dockerfile`. Implement orchestrator-ready liveness (`/healthz`) and readiness (`/readyz`) probes, static compilation with `CGO_ENABLED=0`, layer-caching optimization, non-root user sandboxing (`USER appuser`), resource constraints via `docker-compose.yml`, and an automated GitHub Actions CI/CD pipeline with testing, linting, race detection, and container build caching.

## Why this matters

Deploying raw Go binaries directly onto virtual machines is rarely how modern cloud services operate. Kubernetes, ECS, and Nomad expect containerized services with non-root security contexts, minimal image footprints, zero extraneous build tooling in runtime images, and responsive health/readiness signals during rolling deployments. Mastering containerization and CI/CD automation ensures your Go services deploy reliably and securely at enterprise scale.

## Prerequisites

Complete Module 15 lessons on Linux containers, Docker multi-stage builds, static compilation flags (`-ldflags`), image security, non-root privileges, GitHub Actions workflows, and Docker Compose orchestration.

## Available concepts

This project draws on multi-stage Docker builds, `scratch` and `alpine` base images, `CGO_ENABLED=0`, `go.mod` layer caching, non-root user sandboxing, Linux signals (SIGTERM/SIGINT) with graceful shutdown, HTTP health probes (`/healthz`, `/readyz`), and GitHub Actions workflow syntax.

## Project overview

The `_starter/` directory provides a baseline Go microservice and an unoptimized single-stage `Dockerfile`.
The starter exhibits classic containerization anti-patterns:
- **Massive Image Footprint**: A single-stage build that ships the full Go toolchain, compiler, and package cache in the production image (~800MB instead of ~15MB).
- **Security Vulnerability**: Runs as the default root user inside the container, granting unnecessary privileges.
- **Busted Layer Caching**: Copies application source files before `go.mod`, invalidating the dependency download cache on every code change.
- **Missing Orchestration Probes**: Inadequate readiness probe during graceful drain and no container-native health check binary.

Your mission is to build an enterprise-ready containerized service:
1. `Dockerfile`: Optimized multi-stage build (`builder` -> `scratch` / minimal runner) with static CGO-free binary, CA certificates, and non-root user.
2. `main.go`: Microservice with `/healthz`, `/readyz` (drainable state), `/api/v1/info`, and internal `-healthcheck` CLI flag.
3. `docker-compose.yml`: Local orchestrator declaration with CPU/memory limits, port forwarding, and periodic health checks.
4. `.github/workflows/ci.yml`: Production CI pipeline running unit tests with `-race`, formatting checks, and cached container image builds.

## Multi-Stage Container Architecture

```
┌──────────────────────────────────────────────────────────┐
│ STAGE 1: BUILDER (golang:1.24-alpine)                    │
│   1. Install ca-certificates & create appuser (UID 10001)│
│   2. COPY go.mod go.sum -> RUN go mod download (CACHED)  │
│   3. COPY source code                                    │
│   4. RUN CGO_ENABLED=0 go build -ldflags="-s -w"         │
└────────────────────────────┬─────────────────────────────┘
                             │ Copies static binary,
                             │ ca-certs & /etc/passwd ONLY
                             ▼
┌──────────────────────────────────────────────────────────┐
│ STAGE 2: RUNTIME (scratch)                               │
│   - Size: ~15 MB                                         │
│   - Attack surface: Zero shell, zero package manager     │
│   - Security: USER appuser:appuser (UID 10001)           │
│   - Entrypoint: /bin/server                              │
└──────────────────────────────────────────────────────────┘
```

## Tasks

### Task 1 — Multi-Stage Dockerfile Optimization

In `_starter/Dockerfile`:
- Define `FROM golang:1.24-alpine AS builder`.
- Install `ca-certificates` and `tzdata`.
- Create an unprivileged user `appuser` (UID `10001`).
- Copy `go.mod` and `go.sum` first; run `go mod download`.
- Copy remaining source files and compile with `CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/server .`.
- Define runtime stage `FROM scratch` (or `alpine`).
- Copy `/etc/ssl/certs/ca-certificates.crt`, `/etc/passwd`, and `/bin/server`.
- Switch user: `USER appuser:appuser`.
- Expose port `8080` and set `ENTRYPOINT ["/bin/server"]`.

### Task 2 — Implement Liveness, Readiness, and Healthcheck Probes

In `_starter/main.go`:
- Implement `/healthz`: returns HTTP 200 OK with `{"status":"ok"}`.
- Implement `/readyz`: returns HTTP 200 OK when service is active, and HTTP 503 Service Unavailable when the service receives a SIGTERM and begins graceful shutdown draining.
- Implement flag `-healthcheck`: queries local `http://127.0.0.1:8080/healthz` and exits with code 0 on success, or code 1 on failure.

### Task 3 — Compose Specification

In `_starter/docker-compose.yml`:
- Declare service `app` building from the local `Dockerfile`.
- Map port `8080:8080`.
- Declare deployment resource limits: `cpus: '0.5'`, `memory: 128M`.
- Configure healthcheck: `test: ["/bin/server", "-healthcheck"]` with interval `15s`, timeout `3s`, and retries `3`.

### Task 4 — GitHub Actions CI/CD Pipeline

In `_starter/.github/workflows/ci.yml`:
- Trigger on `push` and `pull_request` to `main`.
- Add test job executing `go test -v -race -coverprofile=coverage.out ./...`.
- Add Docker build job with `docker/setup-buildx-action` and GitHub Actions layer caching (`type=gha`).

## Verification

Run the automated test suite:

```bash
go test -v ./curriculum/modules/15-docker-cicd-deployment/projects/dockerized-service-with-ci-cd/_solution
```

Expected output:
```
=== RUN   TestHealthzEndpoint
--- PASS: TestHealthzEndpoint (0.00s)
=== RUN   TestReadyzEndpoint_StateTransition
--- PASS: TestReadyzEndpoint_StateTransition (0.00s)
=== RUN   TestServiceInfoEndpoint
--- PASS: TestServiceInfoEndpoint (0.00s)
=== RUN   TestInternalHealthCheckFunction
--- PASS: TestInternalHealthCheckFunction (0.00s)
=== RUN   TestDockerfileSpecification
--- PASS: TestDockerfileSpecification (0.00s)
=== RUN   TestCIDeploymentWorkflow
--- PASS: TestCIDeploymentWorkflow (0.00s)
PASS
ok  	github.com/rasel9t6/the-go-engineer/curriculum/modules/15-docker-cicd-deployment/projects/dockerized-service-with-ci-cd/_solution	0.965s
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Multi-Stage Dockerfile** | Minimal scratch/alpine image, non-root user, CA certs included, static CGO-free binary, and proper layer caching | Multi-stage build works but runs as root or invalidates dependency layer | Single-stage image shipping entire Go SDK |
| **Probes & Graceful Shutdown** | Distinct `/healthz` and `/readyz` endpoints; flips readiness to false on SIGTERM during drain; native healthcheck CLI flag | Health endpoints present but does not drain readiness on shutdown | Missing health probes |
| **Compose & Resource Bounds** | Declarative compose file with CPU/memory constraints, non-root ports, and healthcheck interval | Compose file present but missing resource constraints | No compose configuration |
| **CI/CD Automation** | GitHub Actions pipeline runs tests with `-race`, formatting checks, and buildx caching | CI pipeline runs tests but omits race detector or container build | No CI/CD workflow |

## NEXT UP

[Module 16: Portfolio & Interview Readiness](../../16-portfolio-interview-readiness/README.md)
