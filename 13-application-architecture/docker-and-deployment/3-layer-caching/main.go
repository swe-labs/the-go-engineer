package main

/*
Docker Layer Caching Optimization

This lesson demonstrates how to optimize Docker builds using layer caching.

KEY INSIGHT: Docker caches layers. If a layer hasn't changed, it reuses the cached result.

WRONG ORDER (forces rebuild on every source change):
    COPY . .                    ← Copies everything, invalidates cache
    RUN go mod download         ← Dependencies downloaded again!
    RUN go build -o app main.go

RIGHT ORDER (caches dependencies):
    COPY go.mod go.sum ./       ← Only copies dependency files
    RUN go mod download         ← Cached if go.mod unchanged
    COPY . .                    ← Only invalidates if source changes
    RUN go build -o app main.go ← Uses cached dependencies

BENEFIT: 10x faster rebuilds during development!

.dockerignore helps too:
    - Exclude .git/ (can be 100MB+)
    - Exclude node_modules/
    - Exclude .env files
    - Exclude test fixtures

See the Dockerfile in this directory for the optimized approach.

BUILD WITH LAYER CACHING BENEFITS:
docker build -t myapp:latest .

First build: Takes ~30 seconds (downloads dependencies)
Second build (no source change): Takes ~5 seconds (uses cache)
Third build (fast rebuild after code change): Takes ~10 seconds (recompile, reuses dep cache)
*/

func main() {
	// This lesson demonstrates Docker layer caching optimization.
	// See the Dockerfile in this directory for the optimized approach.
	// Run: docker build -t myapp:latest .
	// First build takes ~30 seconds (downloads dependencies)
	// Subsequent builds are much faster due to layer caching.
	println("Docker layer caching lesson - see Dockerfile in this directory")
}
