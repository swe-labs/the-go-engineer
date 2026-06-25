# Multi-stage builds

## Learning objective

Design multi-stage Dockerfiles that separate build and runtime environments, produce tiny images using `scratch` or `distroless` base images, and configure `CGO_ENABLED=0` for fully static Go binaries.

## Why this matters

A multi-stage Dockerfile is the standard way to build Go containers in production. It produces a final image containing only the compiled binary — no compiler, no package manager, no shell, no source code. This reduces the attack surface, decreases pull time, and eliminates variant analysis work for security scanners. Every major Go project (Docker CLI, Kubernetes, Prometheus, Grafana, Caddy) uses multi-stage builds.

## Mental model

Multi-stage build is like renting a fully equipped kitchen (builder stage) to prepare a meal, then transferring only the finished dish to a clean plate (runtime stage). The kitchen is discarded after use. The builder stage has the Go compiler, module cache, and all build tools. The runtime stage has nothing but the binary. Nothing from the builder leaks into the final image unless you explicitly `COPY --from=builder`.

## Core idea

A multi-stage Dockerfile uses multiple `FROM` instructions. Each `FROM` starts a new stage. You can name stages with `AS` and copy artifacts between them with `COPY --from=<stage>`. Only the last stage is used for the final image (unless you specify a target with `--target`).

| Stage | Base image | Contents | Purpose |
|---|---|---|---|
| Builder | `golang:1.25-alpine` | Compiler, modules, source | Compile the binary |
| Runtime | `scratch` | Only the binary | Run the service |

## Under the hood

Each `FROM` in a multi-stage Dockerfile creates a separate build stage with its own layer graph. The `COPY --from=builder /app/server /server` instruction extracts the file from the builder stage's filesystem and places it in the current stage's layer. This does not carry over any of the builder's other layers — only the specified file is copied. Docker's build cache treats each stage independently. If the builder stage is invalidated, only the final COPY layer in the runtime stage needs rebuilding.

The `scratch` image is an empty filesystem — literally no files, no directories, no shell. A Go binary compiled with `CGO_ENABLED=0` needs nothing from the OS. It contains its own DNS resolver, TLS stack, and timezone data baked in. The `gcr.io/distroless/static` image adds `/etc/ssl/certs`, `/etc/passwd` (for `nobody` user), and `/tmp`.

## How Go uses it

The canonical multi-stage Dockerfile for Go:

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

FROM scratch
COPY --from=builder /app/server /server
EXPOSE 8080
USER 1001
CMD ["/server"]
```

Go-specific flags in the builder stage:
- `CGO_ENABLED=0` — produces a fully static binary with no libc dependency
- `GOOS=linux` — ensures Linux binary even if building on macOS/Windows
- `-ldflags="-s -w"` — strips debug info, reducing binary size by ~30%

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

type BuildInfo struct {
	Version     string `json:"version"`
	GoVersion   string `json:"go_version"`
	GoArch      string `json:"go_arch"`
	GoOS        string `json:"go_os"`
	CgoEnabled  bool   `json:"cgo_enabled"`
	StaticBuild bool   `json:"static_build"`
}

func buildInfoHandler(w http.ResponseWriter, r *http.Request) {
	info := BuildInfo{
		Version:     "1.0.0",
		GoVersion:   runtime.Version(),
		GoArch:      runtime.GOARCH,
		GoOS:        runtime.GOOS,
		CgoEnabled:  false,
		StaticBuild: true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func main() {
	http.HandleFunc("/build", buildInfoHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server starting on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
```

## Step-by-step execution

1. `docker build -t myapp .` sends the context to dockerd.
2. Docker pulls `golang:1.25-alpine` for the builder stage.
3. `WORKDIR /app` creates the working directory.
4. `COPY go.mod go.sum ./` copies dependency files.
5. `RUN go mod download` caches dependencies.
6. `COPY . .` copies source code.
7. `RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server .` produces a static binary.
8. Docker switches to the `scratch` stage.
9. `COPY --from=builder /app/server /server` copies only the binary.
10. `EXPOSE 8080` documents the port.
11. `USER 1001` sets the runtime user.
12. `CMD ["/server"]` sets the entrypoint.

The builder stage and its ~350 MB of layers are discarded. The final image contains one ~15 MB layer.

## Common mistakes

- **Forgetting `CGO_ENABLED=0`.** The resulting binary links against glibc and crashes on `scratch` with "no such file or directory" when trying to load shared libraries.
- **Not stripping debug info.** Without `-ldflags="-s -w"`, the binary includes DWARF debug info and symbol tables, increasing size by 30-50%.
- **Running the runtime stage as root.** The `scratch` image has no `/etc/passwd` by default, so `USER 1001` creates an anonymous user. For distroless, `USER nonroot` is already defined.
- **Including unnecessary files from the builder.** `COPY --from=builder /go/bin/server` is better than `COPY --from=builder /go/ .` which copies the entire Go toolchain.
- **Not pinning the Go version.** Using `golang:latest` in CI breaks builds when a new Go version is released. Pin to `golang:1.25-alpine`.

## Debugging walkthrough

Consider a Dockerfile that builds a Go binary with CGO:

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server .

FROM scratch
COPY --from=builder /app/server /server
CMD ["/server"]
```

**Symptom:** `docker run myapp` immediately exits with "exec /server: no such file or directory".

**Investigation:** Run `file server` on the built binary. Output shows "dynamically linked" with references to `libc.so.6`. The binary was compiled with CGO_ENABLED=1 (default).

**Root cause:** `go build` without `CGO_ENABLED=0` produces a dynamically linked binary that needs glibc. The `scratch` image has no libc.

**Fix:** Add the environment variable:
```dockerfile
RUN CGO_ENABLED=0 go build -o server .
```
Verify with `file server` — output should say "statically linked".

## Production notes

- Use `gcr.io/distroless/static` instead of `scratch` when you need TLS certificates or timezone data. It adds only ~2 MB.
- Use `-ldflags="-s -w -X main.version=$(gitsem)"` to embed version info at build time.
- Scan the final image with `trivy` or `docker scout` — a multi-stage build should have zero vulnerabilities because it contains almost nothing.
- For Alpine-based runtime images, use `alpine:3.20` (not `golang:alpine`) which is ~5 MB.
- Set `LABEL` in the final stage for metadata: `LABEL org.opencontainers.image.source="https://github.com/myorg/myapp"`.
- Use `docker build --target builder` to produce a development image with all tools for local debugging.

## Performance implications

- Multi-stage builds produce images 10-50x smaller than single-stage builds. A 15 MB image pulls in <1 second vs 30+ seconds for an 800 MB image.
- The build itself takes longer because there are more layers, but the CI cache mitigates this. Only modified stages are rebuilt.
- Stripping debug info (`-ldflags="-s -w"`) reduces binary size by ~30% at the cost of losing line numbers in stack traces. For production, strip symbols but keep the Go binary's own crash reporting.
- `CGO_ENABLED=0` binaries have slightly different DNS resolution behavior (pure Go resolver vs glibc's nsswitch) but are faster in containers because they skip `/etc/nsswitch.conf` parsing.

## Practice task

Write a Go program that prints its build mode (static or dynamic) and Go version. Then write a multi-stage Dockerfile that builds the program with `CGO_ENABLED=0` and runs it on `scratch`. Verify the image size with `docker images`. The Dockerfile should go in the lesson directory.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/08-multi-stage-builds
go test -v ./curriculum/modules/15-docker-cicd-deployment/lessons/08-multi-stage-builds
```

## Review questions

1. How does `COPY --from=builder` work across stages in a multi-stage Dockerfile?
2. Why does a Go binary compiled with the default `CGO_ENABLED=1` fail on a `scratch` image?
3. What do the `-ldflags="-s -w"` flags do, and why are they recommended for production?
4. When would you choose `gcr.io/distroless/static` over `scratch` as the runtime base image?
5. What is the final image size difference between a single-stage build with `golang:1.25` and a multi-stage build with `scratch`?

## NEXT UP

Docker Compose — defining multi-service environments with `docker-compose.yml` for Go apps with databases and caches.
