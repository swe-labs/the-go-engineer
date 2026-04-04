# Docker & Deployment (§14 · Docker)

## Overview

This section teaches containerization with Docker through practical examples. You'll learn to write efficient Dockerfiles, leverage multi-stage builds for minimal images, and optimize layer caching for faster rebuilds.

## Why Docker for Go?

Go has an **unfair advantage** in the container world:

- **Statically linked binaries**: No runtime dependencies
- **Single binary output**: Unlike Python/Node which need large runtimes
- **Fast startup**: Milliseconds instead of seconds
- **Minimal Docker images**: 5-20MB instead of 500MB+

## Learning Path

### Lesson 1: Single-Stage Dockerfile (DO.1)
**Concepts**: Basic dockerfile structure, FROM, WORKDIR, COPY, RUN, ENTRYPOINT

Learn the fundamentals:
```dockerfile
FROM alpine:3.19
WORKDIR /app
COPY . .
RUN go build -o main main.go
ENTRYPOINT ["./main"]
```

### Lesson 2: Multi-Stage Builds (DO.2)
**Concepts**: Build stage, runtime stage, COPY --from, minimal final image

The Go superpower - separate compile from runtime:
```dockerfile
# Build stage (heavy, compile only)
FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o app main.go

# Runtime stage (minimal, run only)
FROM alpine:3.19
COPY --from=builder /build/app .
ENTRYPOINT ["./app"]
```

**Impact**: Reduces final image from 350MB → 15MB!

### Lesson 3: Layer Caching & Optimization (DO.3)
**Concepts**: Docker layer caching, go.mod caching, .dockerignore, build optimization

Cache your dependencies separately:
```dockerfile
# Cache dependencies first (rarely change)
COPY go.mod go.sum ./
RUN go mod download

# Then copy source (changes frequently)
COPY . .

# Build
RUN go build -o app main.go
```

**Key Insight**: Order matters! Put stable layers earlier to leverage build cache.

---

## Directory Structure

```
docker-and-deployment/
├── README.md                  (This file)
├── 1-dockerfile/              (Basic Dockerfile example)
├── 2-multi-stage/             (Optimized multi-stage build)
└── 3-layer-caching/           (Caching best practices)
```

---

## Best Practices Summary

✅ **Always use multi-stage builds for Go**
✅ **Leverage layer caching** - dependencies before source
✅ **Use Alpine Linux** - minimal, secure, fast
✅ **Create .dockerignore** - exclude `.git`, test files, etc.
✅ **COPY before RUN** - cache dependencies separately
✅ **Tag images with versions** - avoid relying on `latest`
✅ **Run as non-root user** - security best practice
✅ **Keep base images updated** - security patches

---

## Real-World Application

This Docker setup is used in the **Enterprise Capstone Project** where:
- Build stage compiles the Go binary
- Runtime stage serves the REST API
- Docker Compose orchestrates PostgreSQL + Go app
- Migrations run automatically on startup

See [Enterprise Capstone](../enterprise-capstone) for the full application.

## Next Steps

🚀 **After this section**: [Enterprise Capstone](../enterprise-capstone) — Build a complete REST API application with Docker Compose and PostgreSQL deployment
