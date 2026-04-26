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

## Current Progress

Opslane now uses an explicit module system instead of one static flagship label.

The current repository state is:

- `OPSL.1` complete: foundation and configuration
- `OPSL.2` complete: database and models
- `OPSL.3` complete: authentication and tenant isolation
- `OPSL.4` complete: HTTP API layer
- `OPSL.5` complete: order processing
- `OPSL.6` next: payment pipeline

Use the progress surface instead of guessing:

```bash
go run ./11-flagship/01-opslane/scripts/progress.go
```

Then use the learner map:

- [Opslane module map](./MODULES.md)
- [OPSL.1 module spec](./modules/01-foundation/README.md)
- [OPSL.2 module spec](./modules/02-database/README.md)
- [OPSL.3 module spec](./modules/03-auth/README.md)
- [OPSL.4 module spec](./modules/04-http-api/README.md)
- [OPSL.5 module spec](./modules/05-order-processing/README.md)

## Module 5 Snapshot

`OPSL.5` establishes:

- an explicit order service layer instead of direct handler-to-store writes
- state transition rules for `pending`, `processing`, `paid`, `failed`, and `cancelled`
- idempotent order entry that returns the existing order on retry
- an inventory coordination seam that future payment and worker modules can deepen

This slice turns order creation into workflow code.
The HTTP handler now parses the request and trusted identity, then hands business rules to the
order service instead of writing straight through to persistence.

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
curl http://localhost:8080/api/v1/me
```

The health endpoint now checks the live PostgreSQL connection before reporting the service as ready.
The `/me` and `/api/v1/me` endpoints are intentionally protected. They return `401 Unauthorized`
without a bearer token and return the tenant-scoped identity when a valid Opslane token is supplied.

### Setup and Login

```bash
curl -X POST http://localhost:8080/api/v1/tenants \
  -H "Content-Type: application/json" \
  -d '{"name":"Acme Inc","slug":"acme"}'

curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"tenant_id":1,"email":"admin@example.com","display_name":"Admin","password":"CorrectHorse7Battery","role":"admin"}'

curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"tenant_id":1,"email":"admin@example.com","password":"CorrectHorse7Battery"}'
```

### Protected API

```bash
curl http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer <token>"

curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"total_cents":2500,"currency":"USD","idempotency_key":"checkout-123"}'

curl -X POST http://localhost:8080/api/v1/payments \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"order_id":1,"provider_reference":"pay_123","amount_cents":2500}'

curl http://localhost:8080/api/v1/orders/1/payments \
  -H "Authorization: Bearer <token>"
```

## Next Step

After `OPSL.5`, continue to [OPSL.6](./modules/06-payment-pipeline/README.md).
That module moves payment handling behind its own reliability-focused workflow boundary.
