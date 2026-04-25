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
| `internal/db` | PostgreSQL setup, migrations, repository contracts, and transaction seams |
| `internal/models` | tenant-aware domain records for users, orders, and payments |
| `Dockerfile` | packages the application for containerized execution |
| `docker-compose.yml` | runs the application with PostgreSQL and the flagship API |
| `internal/` | holds the application boundaries that later modules will deepen |

## Module 2 Focus

Module 2 establishes:

- tenant-aware core models
- PostgreSQL startup and migration flow
- explicit repository contracts and transaction seams
- conservative pool settings and local persistence defaults

We now move beyond a bootable shell.
This slice gives Opslane a real persistence boundary so later modules can add auth, order flow, and
async processing on top of a coherent data model instead of placeholders.

## Run the Project

For the local shell:

```bash
docker-compose up -d db
go run ./11-flagship/01-opslane/cmd/server
```

Optional environment overrides:

```bash
$env:OPSLANE_ENV="development"
$env:OPSLANE_HTTP_ADDR=":8080"
$env:OPSLANE_LOG_LEVEL="debug"
$env:OPSLANE_DB_DSN="postgres://opslane:secretpassword@localhost:5432/opslane?sslmode=disable"
go run ./11-flagship/01-opslane/cmd/server
```

Make sure Docker is available, then run from this directory:

```bash
docker-compose up -d --build
```

## Local Data Volume

The Compose stack uses the versioned volume `opslane_pg16_data` for PostgreSQL 16 data.

That name is intentional. It avoids reusing older local `pgdata` volumes that may contain a
previous PostgreSQL major version or older bootstrap credentials. If you previously ran an older
Opslane compose stack, the new stack starts from its own database volume instead of silently
reusing incompatible data.

## Current HTTP Surface

```bash
curl http://localhost:8080/
curl http://localhost:8080/health
```

The health endpoint now checks the live PostgreSQL connection before reporting the service as ready.

## Next Step

After Module 2 is in place, the next build slice is authentication and tenant isolation so Opslane
can turn this persistence foundation into an actual multi-tenant request flow.
