# Docker Compose

## Learning objective

Define multi-service Go applications using `docker-compose.yml` with services, environment variables, volumes, ports, and inter-service dependencies managed by Compose.

## Why this matters

A Go HTTP service never runs alone in production. It connects to PostgreSQL, Redis, a message broker, or other internal services. Docker Compose lets you define the entire stack in a YAML file and spin it up with one command. For local development, Compose replaces the need to install and configure databases, caches, and queues manually. For CI, Compose provides a reproducible test environment.

## Mental model

Compose is a multi-container orchestrator for a single host. You declare services (containers), networks (communication channels), and volumes (persistent storage) in a YAML file. `docker compose up` creates everything in the right order based on `depends_on`. Each service gets its own hostname (the service name) on an internal DNS, so your Go app connects to `postgres:5432` instead of `localhost:5432`.

## Core idea

A `docker-compose.yml` file defines:

| Key | Purpose | Example |
|---|---|---|
| `services` | List of containers | `api`, `db`, `redis` |
| `build` | Build context for the image | `.` or `./api` |
| `image` | Pre-built image | `postgres:16-alpine` |
| `ports` | Host-to-container port mapping | `"8080:8080"` |
| `environment` | Environment variables | `DATABASE_URL=postgres://...` |
| `volumes` | Persistent data/configuration | `pgdata:/var/lib/postgresql/data` |
| `depends_on` | Service startup order | `- db` |
| `networks` | Custom network definitions | `appnet` |

## Under the hood

`docker compose up` reads the YAML file, creates a default network (bridge), and starts each service. The `depends_on` option controls startup order but NOT readiness — it waits for the container to start, not for the process inside to be ready. Compose assigns DNS records: the service name resolves to the container's IP on the Compose network. Go's `net.LookupHost("db")` resolves to the PostgreSQL container. Compose manages the lifecycle of all containers: `docker compose down` stops and removes everything.

## How Go uses it

A typical Go development environment with Compose:

```yaml
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: "postgres://user:pass@db:5432/myapp?sslmode=disable"
      REDIS_URL: "redis://redis:6379"
    depends_on:
      - db
      - redis

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: myapp
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  pgdata:
```

Your Go app connects to `db:5432` (not `localhost:5432`). The Go `database/sql` connection string references the service name as the host.

## Go example

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Config struct {
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
	DatabaseURL string `json:"database_url"`
	RedisURL    string `json:"redis_url"`
	ListenAddr  string `json:"listen_addr"`
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	cfg := Config{
		ServiceName: getEnv("SERVICE_NAME", "api"),
		Version:     getEnv("VERSION", "1.0.0"),
		DatabaseURL: getEnv("DATABASE_URL", "not-configured"),
		RedisURL:    getEnv("REDIS_URL", "not-configured"),
		ListenAddr:  getEnv("LISTEN_ADDR", ":8080"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/health", healthHandler)

	addr := getEnv("LISTEN_ADDR", ":8080")
	fmt.Printf("Server starting on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
```

## Step-by-step execution

1. `docker compose up` reads `docker-compose.yml`.
2. Compose creates a default network called `myapp_default`.
3. The `db` service starts first (PostgreSQL container).
4. Compose does not wait for PostgreSQL to be ready — it just waits for the container to start.
5. The `redis` service starts.
6. The `api` service starts. The Go binary reads `DATABASE_URL` from the environment: `postgres://user:pass@db:5432/myapp`.
7. The Go HTTP server starts on `:8080`.
8. The Go app connects to `db:5432` via the Compose internal DNS.

## Common mistakes

- **Connecting to `localhost` instead of the service name.** Inside a Compose network, `localhost` refers to the container itself, not the host or other containers. Use the service name as the hostname (`db:5432`).
- **Not handling container readiness.** Compose `depends_on` only waits for container start, not for the service inside to be ready. The Go app crashes on startup because PostgreSQL is still initializing. Fix: add a retry loop in Go or use `healthcheck` on the dependency.
- **Exposing database ports in production.** Mapping `5432:5432` in Compose is useful for local debugging but unnecessary and dangerous in production. Remove `ports` for internal services in production Compose files.
- **Hardcoding URLs.** Writing `DATABASE_URL=postgres://user:pass@localhost:5432/myapp` in the Dockerfile makes the image environment-specific. Use environment variables set in Compose.
- **Using `latest` tags.** `postgres:latest` changes unpredictably. Pin to `postgres:16-alpine`.

## Debugging walkthrough

Consider a Go service that cannot connect to PostgreSQL in Compose:

```go
db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/myapp")
```

**Symptom:** The Go app starts, logs "connection refused", and exits.

**Investigation:** Inside the container, `localhost` is the container itself, not the host or the `db` container. Run `docker compose exec api ping db` to verify network connectivity.

**Root cause:** The connection string uses `localhost` instead of `db` (the Compose service name). Inside the container, `localhost` resolves to 127.0.0.1 — the container's own loopback interface, which has no PostgreSQL running.

**Fix:** Change the connection string to use the service name:
```go
db, err := sql.Open("postgres", "postgres://user:pass@db:5432/myapp")
```
Set it via environment variable:
```yaml
environment:
  DATABASE_URL: "postgres://user:pass@db:5432/myapp?sslmode=disable"
```

## Production notes

- Compose is for development and CI, not production. Use Kubernetes, Docker Swarm, or ECS for production.
- Use `docker compose config` to validate and view the resolved configuration.
- Use `.env` files for per-developer overrides. Compose reads `.env` automatically.
- Set `restart: unless-stopped` on services in CI to handle crashes gracefully.
- Use `docker compose --profile` for optional services (e.g., `--profile monitoring` for Prometheus and Grafana).
- For production-like CI, use `docker compose up --wait` with health checks to ensure services are ready before tests run.

## Performance implications

- Compose adds negligible overhead — it is a thin Python/Go CLI over the Docker API.
- Running multiple databases in Compose for development uses significant RAM. Use resource limits in `deploy.resources` for CI.
- The internal DNS in Compose is backed by Docker's embedded DNS server — resolution is fast (<1ms).
- Network performance between Compose services is native bridge networking, near line rate.

## Practice task

Write a Go program that reads `DATABASE_URL` and `REDIS_URL` from environment variables and exposes them at a `/config` JSON endpoint. Then write a `docker-compose.yml` that runs the Go service alongside PostgreSQL and Redis. The Compose file should go in the lesson directory.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/09-docker-compose
go test -v ./curriculum/modules/15-docker-cicd-deployment/lessons/09-docker-compose
```

## Review questions

1. Why does `localhost` in a Go app inside a Compose service NOT refer to the host machine?
2. How does service discovery work between Compose services?
3. What is the difference between `depends_on` and a health check for startup ordering?
4. Why should database credentials be set via environment variables rather than hardcoded in the Dockerfile?
5. When would you remove the `ports` section from a service in `docker-compose.yml`?

## NEXT UP

Container health checks — implementing `/healthz` endpoints in Go, using Docker's HEALTHCHECK instruction, and distinguishing readiness from liveness.
