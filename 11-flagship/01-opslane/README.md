# Opslane

## Mission

This surface is the active flagship implementation project for the stable v2.1 curriculum.

Opslane is where the curriculum stops being a collection of isolated lessons and becomes one
production-shaped backend that we can build, test, and harden in public.

## Stage Ownership

This project belongs to [11 Flagship](../README.md).

## Why This Project Matters

This project exists so the curriculum has one integrated system where earlier stage skills meet:

- backend request and data flow
- repository-style persistence boundaries
- package and handler structure
- runtime and deployment-oriented behavior
- longer feedback loops than a single exercise can provide

## Current Project Shape

| Area | Role |
| --- | --- |
| `cmd/server` | runnable flagship server entrypoint |
| `internal/config` | validated configuration loading and startup contract |
| `Dockerfile` | packages the application for containerized execution |
| `docker-compose.yml` | runs the application with supporting services |
| `internal/` | holds the application boundaries that later modules will deepen |

## Module 1 Focus

Module 1 establishes:

- a real server shell
- environment-aware configuration
- safe defaults and explicit overrides
- fail-fast startup validation

This is intentionally the smallest useful slice of the flagship.
We are proving the project can boot cleanly before layering in database, auth, and workflow logic.

## Run the Project

For the local shell:

```bash
go run ./11-flagship/01-opslane/cmd/server
```

Optional environment overrides:

```bash
$env:OPSLANE_ENV="development"
$env:OPSLANE_HTTP_ADDR=":8080"
$env:OPSLANE_LOG_LEVEL="debug"
go run ./11-flagship/01-opslane/cmd/server
```

Make sure Docker is available, then run from this directory:

```bash
docker-compose up -d --build
```

## Current HTTP Surface

```bash
curl http://localhost:8080/
curl http://localhost:8080/health
```

## Next Step

After Module 1 is in place, the next build slice is database and model foundations so Opslane can
move from a bootable shell into a real multi-tenant service.
