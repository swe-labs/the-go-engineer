# Docker basics

## Learning objective

Write a Dockerfile for a Go application, build a minimal Docker image using `FROM`, `RUN`, `COPY`, `CMD`, and understand how Docker containers map to Linux processes.

## Why this matters

Docker is the standard way to package and distribute Go services. A single static Go binary is ideal for Docker — it has no dependencies, no interpreter, and starts in milliseconds. Every production Go service you encounter will be deployed as a container. Knowing how to write an efficient, secure Dockerfile is a prerequisite for shipping software.

## Mental model

A container is just a Linux process isolated by namespaces (PID, network, mount) and constrained by cgroups. The Dockerfile is a recipe that assembles a read-only filesystem snapshot (image). Each instruction adds a layer — an overlayfs diff. When you run a container, Docker overlays a writable layer on top and executes CMD/ENTRYPOINT as PID 1. For Go developers, the model simplifies: Go binary + minimal base image = a running process. Everything else is packaging detail.

## Core idea

A Dockerfile is a sequence of instructions that build an image:

| Instruction | Purpose | Go-specific advice |
|---|---|---|
| `FROM` | Base image | Use `golang:1.25-alpine` for build, `scratch` or `distroless` for runtime |
| `RUN` | Execute command | Use for `go build`, never for runtime config |
| `COPY` | Add files from context | Copy only the Go binary, not source |
| `CMD` | Default command | Use exec form: `CMD ["./server"]` |
| `ENTRYPOINT` | Fixed command | Prefer over CMD when arguments are fixed |
| `EXPOSE` | Document port | Informational only, does not publish the port |

## Under the hood

Docker is a client-server architecture. The `docker` CLI talks to `dockerd` via a REST API on `/var/run/docker.sock`. When you `docker build`, the daemon receives the build context (tarred), unpacks it, and parses the Dockerfile instruction-by-instruction. Each instruction becomes a layer: a commit of the filesystem diff using overlayfs. Layers are content-addressed and cached. When you `docker run`, the daemon creates a new mount namespace, overlays image layers plus a writable layer, sets up networking (bridge by default), then forks ENTRYPOINT/CMD as PID 1.

## How Go uses it

Go's static compilation model makes it the ideal language for Docker:
- No runtime dependency: a Go binary compiles to a static executable with no libc requirement (when `CGO_ENABLED=0`).
- Small binary size: 10-20 MB for a full HTTP server.
- Fast startup: typically under 100ms.
- Built-in signal handling: Go's runtime handles SIGTERM/SIGINT for graceful shutdown.

The canonical Go Dockerfile pattern is:
```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM scratch
COPY --from=builder /app/server /server
CMD ["/server"]
```

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
)

type Info struct {
	Service  string `json:"service"`
	Version  string `json:"version"`
	GoArch   string `json:"go_arch"`
	GoOS     string `json:"go_os"`
	Hostname string `json:"hostname"`
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	info := Info{
		Service:  "docker-basics",
		Version:  "1.0.0",
		GoArch:   runtime.GOARCH,
		GoOS:     runtime.GOOS,
		Hostname: hostname,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func main() {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	http.HandleFunc("/info", infoHandler)
	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
```

## Step-by-step execution

1. `docker build -t myapp .` sends the Dockerfile and directory to dockerd.
2. Docker pulls the `golang:1.25-alpine` base image if not cached.
3. `WORKDIR /app` creates the directory in the image.
4. `COPY go.mod go.sum ./` copies dependency files.
5. `RUN go mod download` downloads dependencies and caches them.
6. `COPY . .` copies the source code.
7. `RUN CGO_ENABLED=0 go build -o server .` compiles a static binary.
8. `FROM scratch` starts a fresh image.
9. `COPY --from=builder /app/server /server` copies only the binary.
10. `CMD ["/server"]` sets the default command.
11. `docker run myapp` starts the container, running the binary as PID 1.

## Common mistakes

- **Using a full OS image like `ubuntu` as the base.** This produces an 800 MB+ image for a 15 MB binary. Use `scratch` or `gcr.io/distroless/static`.
- **Using shell form CMD.** `CMD ./server` wraps the command in `/bin/sh -c`, which does not forward signals. Always use exec form: `CMD ["./server"]`.
- **Building with CGO_ENABLED=1.** The resulting binary links against glibc and crashes on `scratch` or `alpine`. Always set `CGO_ENABLED=0` in the builder stage.
- **Copying the entire source tree.** `COPY . .` includes `.git`, `node_modules`, `vendor`, and other bloat. Use `.dockerignore` to exclude them.
- **Running as root in the container.** This violates the principle of least privilege. Add `USER 1001` at the end of your Dockerfile.

## Debugging walkthrough

Consider a Dockerfile that produces a 1.2 GB image for a simple Go server:

```dockerfile
FROM golang:1.25
COPY . .
RUN go build -o server .
CMD ["./server"]
```

**Symptom:** `docker images` shows 1.2 GB. Deployment is slow.

**Investigation:** Check each layer's size with `docker history myapp`. The `golang:1.25` base image is ~800 MB. The `COPY . .` layer includes the host's Go module cache and `.git` history.

**Root cause:** Using the full `golang` image instead of `golang:alpine` or multi-stage build. Copying everything instead of using a `.dockerignore`.

**Fix:** Use multi-stage build with `golang:1.25-alpine` for the builder and `scratch` for the runtime. Add a `.dockerignore` with `.git`, `vendor/`, `README.md`. Final image: ~15 MB.

## Production notes

- Always tag your images with a specific version, not `latest`: `docker build -t myapp:v1.0.0 .`.
- Use `gcr.io/distroless/static` as the runtime base for Go apps — it has ca-certificates, /etc/passwd, and /tmp, but no shell.
- Scan images for vulnerabilities with `docker scout` or `trivy` before deploying.
- Set `USER 1001` in the Dockerfile to run as non-root.
- Use `HEALTHCHECK` instruction to enable Docker's container health monitoring.
- Pin the base image digest (`golang:1.25-alpine@sha256:...`) for reproducible builds.

## Performance implications

- Image size directly affects pull time on Kubernetes nodes. A 15 MB Go image pulls in <1 second on a fast connection. An 800 MB image takes minutes.
- The `scratch` image adds zero overhead — it has no files. The only runtime cost is the Go binary itself.
- Docker's copy-on-write layer has near-zero overhead for reads. Writes to the container layer are slower on some storage drivers (devicemapper) but fine on overlay2.
- Startup time for a Go container is dominated by binary loading (~50-100ms) and init logic, not Docker overhead.

## Practice task

Write a Go HTTP server with a `/health` endpoint that returns `{"status": "ok"}` and a `/version` endpoint that returns the version from an environment variable. Then write a Dockerfile that uses multi-stage build with `scratch` as the final base. The Dockerfile should go in the lesson directory alongside the Go source.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/06-docker-basics
go test -v ./curriculum/modules/15-docker-cicd-deployment/lessons/06-docker-basics
```

## Review questions

1. Why is `scratch` preferred over `ubuntu` as the base image for Go applications?
2. What is the difference between `CMD` and `ENTRYPOINT` in a Dockerfile?
3. Why must the Dockerfile use the exec form `CMD ["./server"]` instead of `CMD ./server`?
4. What does `CGO_ENABLED=0` do, and why is it important for Docker builds?
5. How does `.dockerignore` reduce image size and improve build cache performance?

## NEXT UP

Docker images and layers — how layers work, layer caching, COPY order optimization, and .dockerignore.
