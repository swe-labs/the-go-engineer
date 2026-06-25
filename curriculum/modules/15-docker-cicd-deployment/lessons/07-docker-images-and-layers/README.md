# Docker images and layers

## Learning objective

Explain how Docker image layers work, optimize `COPY` ordering for cache efficiency, write `.dockerignore` files, and evaluate layer size impact on build and deployment speed.

## Why this matters

Every Dockerfile instruction creates a layer. Layers are cached and reused across builds. A poorly ordered Dockerfile invalidates the cache on every source change, forcing a full rebuild every time. A well-ordered Dockerfile rebuilds in seconds by reusing cached layers. Understanding layers is the difference between a CI pipeline that takes 30 seconds and one that takes 10 minutes.

## Mental model

Think of image layers as a stack of transparencies. Each transparency shows only the files that were added or changed by that instruction. When a process reads a file, Docker's union filesystem (overlayfs) walks the stack from top to bottom and returns the first version it finds. Deleting a file in a later layer does not erase it from earlier layers — it just places a "hidden" marker on top. This explains why a 800 MB build stage can yield an 8 MB final image, and why deleting intermediate files mid-build does nothing to shrink the image.

## Core idea

Each Dockerfile instruction creates a read-only layer:

| Instruction | Layer content | Cache key |
|---|---|---|
| `FROM golang:1.25` | Base image filesystem | Image digest |
| `WORKDIR /app` | Directory metadata | Command string |
| `COPY go.mod go.sum ./` | Files from build context | File checksums |
| `RUN go mod download` | Downloaded modules | Command output |
| `COPY . .` | All source files | File checksums |
| `RUN go build -o server .` | Compiled binary | Command output |

The cache is invalidated when the cache key changes. For `COPY`, the key is the file checksum. For `RUN`, it is the command string. When one layer is invalidated, all subsequent layers are rebuilt.

## Under the hood

Docker images are stored under `/var/lib/docker/overlay2` (Linux) as a set of read-only snapshots. Each layer is a tarball of filesystem diffs. The union mount (overlayfs) combines all layers via a `lowerdir` (image layers) and `upperdir` (container writable layer). When a container writes a file, the copy-on-write mechanism copies the file from a lower layer to the upper layer. `docker history` decompiles an image into its instruction list with per-layer sizes built from the diff metadata. The `--squash` flag merges all layers into one, trading cache granularity for a smaller image.

## How Go uses it

Go's large standard library and module graph make layer management critical:

```dockerfile
# INEFFICIENT: Every source change invalidates the entire build
FROM golang:1.25-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server .

# EFFICIENT: Dependencies cached separately from source
FROM golang:1.25-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .
```

The efficient version separates dependency resolution (slow, changes rarely) from source compilation (changes on every edit).

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type LayerInfo struct {
	Command   string `json:"command"`
	SizeBytes int64  `json:"size_bytes"`
	Cacheable bool   `json:"cacheable"`
}

var recommendedLayers = []LayerInfo{
	{Command: "FROM golang:1.25-alpine", SizeBytes: 150_000_000, Cacheable: true},
	{Command: "WORKDIR /app", SizeBytes: 0, Cacheable: true},
	{Command: "COPY go.mod go.sum ./", SizeBytes: 50_000, Cacheable: true},
	{Command: "RUN go mod download", SizeBytes: 200_000_000, Cacheable: true},
	{Command: "COPY . .", SizeBytes: 500_000, Cacheable: false},
	{Command: "RUN CGO_ENABLED=0 go build -o server .", SizeBytes: 15_000_000, Cacheable: false},
}

func layersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendedLayers)
}

func main() {
	http.HandleFunc("/layers", layersHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Listening on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
```

## Step-by-step execution

1. `docker build` sends the build context to dockerd.
2. Docker checks its layer cache for `FROM golang:1.25-alpine`. If the image is cached, layer 0 is reused.
3. `WORKDIR /app` creates a metadata layer. Always cached unless the command changes.
4. `COPY go.mod go.sum ./` computes checksums of the two files. If they match the cache, this layer is reused.
5. `RUN go mod download` runs the command. The cache key is the command string and the previous layer hash. Cached if nothing changed.
6. `COPY . .` computes checksums of all source files. Any edit invalidates this layer.
7. `RUN go build` must recompile because the previous layer changed.

## Common mistakes

- **Putting `COPY . .` before `RUN go mod download`.** Every source change forces re-downloading all modules. Solution: copy dependency files first, download, then copy source.
- **Copying the entire source tree without `.dockerignore`.** The `.git` directory alone can be hundreds of megabytes. Solution: create a `.dockerignore` with `.git`, `node_modules`, `README.md`, `Dockerfile`.
- **Running `apt-get install` in the same layer as `go build`.** This mixes build-time tools with the application, inflating the final image. Solution: multi-stage build.
- **Deleting files in the same RUN command that created them.** `RUN wget ... && rm ...` still leaves the downloaded file in the previous layer. Solution: use multi-stage build or chain commands in a single RUN.
- **Not considering layer count.** Each layer adds metadata overhead. Squash layers for production images shared across many deploys.

## Debugging walkthrough

Consider a Dockerfile that produces a 450 MB image:

```dockerfile
FROM golang:1.25
RUN apt-get update && apt-get install -y curl
COPY . .
RUN go build -o server .
RUN rm -rf /var/lib/apt/lists/*
CMD ["./server"]
```

**Symptom:** Image is 450 MB even though the Go binary is only 15 MB.

**Investigation:** Run `docker history myapp` and inspect per-layer sizes. The `apt-get install` layer is 200 MB. The `rm -rf` layer creates a new layer with a deletion marker but the apt files are still in the previous layer.

**Root cause:** Removing files in a separate RUN command does not shrink the image. The files exist in the lower layer and are only hidden by the upper layer's whiteout file.

**Fix:** Use multi-stage build. The builder stage has all the tools, the runtime stage has only the binary:
```dockerfile
FROM golang:1.25 AS builder
COPY . .
RUN go build -o server .

FROM scratch
COPY --from=builder /go/server /server
CMD ["/server"]
```
Final size: ~15 MB.

## Production notes

- Use `.dockerignore` to exclude `.git`, `Dockerfile`, `README.md`, `*.md`, `vendor/`, `node_modules/`, `tmp/`, and CI-specific files.
- Order `COPY` instructions from least-frequently-changing to most-frequently-changing: `COPY go.mod go.sum ./` before `COPY . .`.
- Monitor image size in CI. Alert if it exceeds a threshold (e.g., 50 MB for a Go service).
- Use `docker scout quickview` to analyze image efficiency.
- Pin base image digests for reproducible builds and to avoid surprise cache invalidations.
- The `--squash` flag (experimental) merges all layers into one, useful for minimal deployment images.

## Performance implications

- Each layer adds ~4 KB of metadata overhead. 50 layers = 200 KB overhead, negligible.
- Layer retrieval from a registry is parallelized across layers. More layers means more concurrent downloads. A 2-layer image downloads fewer but larger blobs; a 20-layer image downloads more but smaller blobs.
- Cache hits reduce build time from minutes to seconds. Optimize for cache hits in CI.
- Copy-on-write has a performance cost for write-heavy workloads inside containers. Write to bind-mounted volumes instead of the container layer.
- `scratch`-based images have zero layer overhead beyond the binary itself.

## Practice task

Write a function `LayerSizeBreakdown(imageName string) ([]LayerInfo, error)` that uses `os/exec` to run `docker history --format` and parse the output. If Docker is not available, return a hardcoded set of typical Go layers. Write a `main()` that prints the layer breakdown for a hypothetical optimized Go Docker image.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/07-docker-images-and-layers
go test -v ./curriculum/modules/15-docker-cicd-deployment/lessons/07-docker-images-and-layers
```

## Review questions

1. Why does `RUN rm -rf /var/lib/apt/lists/*` in a separate layer NOT reduce the final image size?
2. How does Docker determine whether a `COPY` instruction can be served from the cache?
3. What is the recommended order of Dockerfile instructions for optimal caching?
4. Why should `.dockerignore` include `.git` and `node_modules`?
5. What is the tradeoff of using `--squash` on a Docker image?

## NEXT UP

Multi-stage builds — separating build environment from runtime for tiny, secure Docker images.
