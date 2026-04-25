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
| `internal/auth` | password hashing, signed tokens, auth context, and tenant identity middleware |
| `internal/config` | validated configuration loading and startup contract |
| `internal/db` | PostgreSQL setup, migrations, repository contracts, and transaction seams |
| `internal/models` | tenant-aware domain records for users, orders, and payments |
| `Dockerfile` | packages the application for containerized execution |
| `docker-compose.yml` | runs the application with PostgreSQL and the flagship API |
| `internal/` | holds the application boundaries that later modules will deepen |

## Module 3 Focus

Module 3 establishes:

- password policy and bcrypt hashing
- signed token issuance and verification
- tenant-aware authenticated request identity
- middleware that rejects anonymous requests before protected handlers run

This slice turns the persistence foundation into a request-level security boundary.
The important rule is simple: protected application code should not guess tenant identity from
untrusted request data. Auth middleware verifies the token once, then stores the tenant-scoped
identity in the request context for downstream handlers and services.

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
$env:OPSLANE_AUTH_TOKEN_SECRET="development-only-opslane-secret-change-me"
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
curl http://localhost:8080/me
```

The health endpoint now checks the live PostgreSQL connection before reporting the service as ready.
The `/me` endpoint is intentionally protected. It returns `401 Unauthorized` without a bearer token
and returns the tenant-scoped identity when a valid Opslane token is supplied.

## Next Step

After Module 3 is in place, the next build slice is the HTTP API layer so Opslane can expose user,
order, and payment operations through stable request and error contracts.
