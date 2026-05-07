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
- `OPSL.6` complete: payment pipeline
- `OPSL.7` complete: event bus and worker pools
- `OPSL.8` complete: caching layer
- `OPSL.9` complete: observability
- `OPSL.10` complete: graceful shutdown and deployment

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
- [OPSL.6 module spec](./modules/06-payment-pipeline/README.md)
- [OPSL.7 module spec](./modules/07-event-workers/README.md)
- [OPSL.8 module spec](./modules/08-caching/README.md)

## Module 5 Snapshot

`OPSL.5` establishes:

- an explicit order service layer instead of direct handler-to-store writes
- state transition rules for `pending`, `processing`, `paid`, `failed`, and `cancelled`
- idempotent order entry that returns the existing order on retry
- an inventory coordination seam that future payment and worker modules can deepen

This slice turns order creation into workflow code.
The HTTP handler now parses the request and trusted identity, then hands business rules to the
order service instead of writing straight through to persistence.

## Module 6 Snapshot

`OPSL.6` establishes:

- an explicit payment gateway interface
- a payment worker boundary for queued payment jobs
- payment service logic for duplicate provider references, retries, and timeouts
- reconciliation-safe pending payments when the gateway outcome is unknown
- order workflow integration so settled payments mark orders paid and failed payments mark orders failed

This slice turns payment creation into reliability-focused workflow code.
The HTTP handler now creates a tenant-scoped payment job and lets the payment service decide how the
database, gateway, and order state machine should cooperate.

## Module 7 Snapshot

`OPSL.7` establishes:

- a bounded event bus with explicit queue-full behavior
- a fixed-size worker pool that drains submitted work
- order, payment, and notification processors behind workflow interfaces
- error callbacks so background failures do not disappear

This slice introduces asynchronous building blocks without hiding goroutines inside handlers.
The system can now talk about queue capacity, backpressure, and safe draining before later modules
wire those primitives into caching, observability, and shutdown behavior.

## Module 8 Snapshot

`OPSL.8` establishes:

- a bounded in-memory cache with TTL and insert-order eviction
- explicit invalidation after order transitions and payment settlements
- singleflight stampede prevention when hot cache keys expire
- HTTP Cache-Control middleware for public and authenticated endpoints
- copy-on-read/write to prevent callers from mutating cached data

This slice introduces caching as an additive optimization, not a hidden source of truth.
PostgreSQL remains the system of record. The cache sits between the service layer and
the repository layer, and invalidation always follows successful writes.

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

## Database Migrations

The project uses formal SQL migrations in the `migrations/` directory:

```bash
# Run migrations
go run ./scripts/migrate.go

# Rollback last migration
go run ./scripts/migrate.go -direction down

# Check migration status
go run ./scripts/migrate.go -direction status
```

Migrations are numbered (001-006) and include:
- `001_create_tenants` - tenant registry
- `002_create_users` - tenant-scoped users
- `003_create_orders` - order workflow
- `004_create_payments` - payment tracking
- `005_seed_data` - development demo data
- `006_create_rate_limits` - distributed rate limiting

Migration governance details:

- policy: [`./docs/MIGRATION_POLICY.md`](./docs/MIGRATION_POLICY.md)
- production semantics: [`./docs/PRODUCTION_SEMANTICS.md`](./docs/PRODUCTION_SEMANTICS.md)
- OpenAPI spec: [`./docs/openapi.yaml`](./docs/openapi.yaml)

## Database Backup & Restore

```bash
# Create backup
./scripts/backup.sh

# Restore from backup (interactive)
./scripts/restore.sh <backup_file>
```

## Seed Data

Load demo data for development:

```bash
go run ./scripts/seed.go

# Reset and re-seed
go run ./scripts/seed.go -reset
```

Demo credentials: `admin@demo.com` / `password123`

## Next Step

All modules are now complete. The repository provides a fully integrated flagship backend demonstrating configuration, database, auth, workflow, async workers, cache, observability, and graceful shutdown.

## Security Notes

Opslane is production-shaped and intentionally educational. For security boundaries and reporting, see the docs folder for threat model, security policy, and known limitations.
