# Opslane Module Progress Map

This file is the learner-facing map for the Opslane flagship.

Use it together with the checker:

```bash
go run ./11-flagship/01-opslane/scripts/progress.go
```

The checker is authoritative about what is complete in the current tree.
This file explains what each module means, what proof looks like, and what comes next.

## How To Use This Map

1. Run the progress checker from the repository root.
2. Open the matching module README in `11-flagship/01-opslane/modules/`.
3. Finish the proof surface for that module before moving on.
4. Re-run the checker to confirm the next module unlocked.

## Required Path

| ID | Module | Current repository state |
| --- | --- | --- |
| `OPSL.1` | Foundation and Configuration | complete |
| `OPSL.2` | Database and Models | complete |
| `OPSL.3` | Authentication and Tenant Isolation | complete |
| `OPSL.4` | HTTP API Layer | complete |
| `OPSL.5` | Order Processing | complete |
| `OPSL.6` | Payment Pipeline | complete |
| `OPSL.7` | Event Bus and Worker Pools | complete |
| `OPSL.8` | Caching Layer | complete |
| `OPSL.9` | Observability | next |
| `OPSL.10` | Graceful Shutdown and Deployment | locked |

## OPSL.1 Foundation and Configuration

What you build: the runnable server shell and validated startup configuration.

Proof surface:

```bash
go test ./11-flagship/01-opslane/internal/config/...
go run ./11-flagship/01-opslane/cmd/server
```

Required files:

- `cmd/server/main.go`
- `internal/config/config.go`
- `internal/config/environment.go`
- `.env.example`
- `Dockerfile`
- `docker-compose.yml`

Module spec: [modules/01-foundation/README.md](./modules/01-foundation/README.md)

## OPSL.2 Database and Models

What you build: the multi-tenant PostgreSQL schema, models, and repository seams.

Proof surface:

```bash
go test ./11-flagship/01-opslane/internal/db/...
```

Required files:

- `internal/db/migrations.go`
- `internal/db/repository.go`
- `internal/models/tenant.go`
- `internal/models/user.go`
- `internal/models/order.go`
- `internal/models/payment.go`

Module spec: [modules/02-database/README.md](./modules/02-database/README.md)

## OPSL.3 Authentication and Tenant Isolation

What you build: token issuance, password hashing, auth middleware, and trusted tenant identity flow.

Proof surface:

```bash
go test ./11-flagship/01-opslane/internal/auth/...
```

Required files:

- `internal/auth/token.go`
- `internal/auth/password.go`
- `internal/auth/service.go`
- `internal/auth/middleware.go`
- `internal/auth/context.go`

Module spec: [modules/03-auth/README.md](./modules/03-auth/README.md)

## OPSL.4 HTTP API Layer

What you build: the public JSON contract, protected routes, rate limiting, and CORS behavior.

Proof surface:

```bash
go test ./11-flagship/01-opslane/internal/handlers/...
go test ./11-flagship/01-opslane/internal/middleware/...
```

Manual spot checks:

```bash
curl http://localhost:8080/health
curl -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d "{\"tenant_id\":1,\"email\":\"admin@example.com\",\"password\":\"CorrectHorse7Battery\"}"
```

Required files:

- `internal/handlers/handlers.go`
- `internal/handlers/api.go`
- `internal/middleware/middleware.go`

Module spec: [modules/04-http-api/README.md](./modules/04-http-api/README.md)

## OPSL.5 Order Processing

What you build: the order state machine, validation rules, inventory coordination, and idempotent workflow entry.

Proof surface:

```bash
go test ./11-flagship/01-opslane/internal/services/...
go test ./11-flagship/01-opslane/internal/handlers/...
```

Behavior proof:

- order state transition tests prove valid and invalid transitions
- idempotent order creation returns the original order instead of creating duplicates

Required files:

- `internal/services/order.go`
- `internal/services/inventory.go`
- `internal/services/validation.go`

Read this module:

- [modules/05-order-processing/README.md](./modules/05-order-processing/README.md)
- [modules/05-order-processing/SURFACE.md](./modules/05-order-processing/SURFACE.md)

## OPSL.6 Payment Pipeline

What you build: payment retries, duplicate protection, and reconciliation-safe gateway flow.

Proof surface:

```bash
go test ./11-flagship/01-opslane/internal/payment/...
go test ./11-flagship/01-opslane/internal/services/...
go test ./11-flagship/01-opslane/internal/handlers/...
```

Behavior proof:

- gateway timeout tests prove pending payments stay reconciliation-safe
- duplicate provider references return existing payments instead of creating duplicate charges
- reconciliation tests prove a delayed success can settle an existing pending payment

Required files:

- `internal/payment/gateway.go`
- `internal/payment/worker.go`
- `internal/services/payment.go`

Read this module:

- [modules/06-payment-pipeline/README.md](./modules/06-payment-pipeline/README.md)
- [modules/06-payment-pipeline/SURFACE.md](./modules/06-payment-pipeline/SURFACE.md)

## OPSL.7 Event Bus and Worker Pools

What you build: bounded asynchronous work, observable worker lifecycle, and queue-pressure handling.

Proof surface:

```bash
go test ./11-flagship/01-opslane/internal/events/...
go test ./11-flagship/01-opslane/internal/workers/...
```

Behavior proof:

- event bus tests prove bounded publish behavior and closed-bus handling
- worker pool tests prove queue saturation, drain behavior, and handler error reporting
- processor tests prove order, payment, and notification adapters call explicit workflow seams

Required files:

- `internal/events/bus.go`
- `internal/events/types.go`
- `internal/workers/pool.go`
- `internal/workers/order_processor.go`
- `internal/workers/payment_processor.go`
- `internal/workers/notification_worker.go`

Read this module:

- [modules/07-event-workers/README.md](./modules/07-event-workers/README.md)
- [modules/07-event-workers/SURFACE.md](./modules/07-event-workers/SURFACE.md)

## OPSL.8 Caching Layer

What you build: cache-aside reads, invalidation boundaries, and bounded TTL behavior.

Proof surface:

```bash
go test ./11-flagship/01-opslane/internal/cache/...
go run ./11-flagship/01-opslane/scripts/progress.go
```

The proof surface covers:

- bounded in-memory cache with TTL and insert-order eviction
- lazy expiry on reads and background janitor sweep
- copy-on-read/write mutation safety
- explicit invalidation after order and payment writes
- prefix-based batch invalidation for tenant-scoped groups
- singleflight stampede prevention
- HTTP Cache-Control middleware

Implemented files:

- `internal/cache/cache.go`
- `internal/cache/store.go`
- `internal/middleware/cache.go`

Read the implementation details:

- [modules/08-caching/README.md](./modules/08-caching/README.md)
- [modules/08-caching/SURFACE.md](./modules/08-caching/SURFACE.md)

## OPSL.9 Observability

What you build: structured logs, correlation IDs, metrics, and trace-friendly request flow.

Proof surface:

```bash
go test ./11-flagship/01-opslane/internal/logging/...
go test ./11-flagship/01-opslane/internal/metrics/...
go test ./11-flagship/01-opslane/internal/tracing/...
go run ./11-flagship/01-opslane/scripts/progress.go
```

The proof surface covers:

- structured logger factory with JSON and text formats
- correlation ID context propagation from HTTP entry through all layers
- HTTP middleware with correlation ID, status capture, and structured request logging
- atomic counters and fixed-bucket histograms for application metrics
- pre-registered HTTP, cache, and worker metrics
- HTTP metrics middleware for request counting and latency
- span tracking with correlation ID linkage
- cross-service correlation header injection and extraction

Implemented files:

- `internal/logging/context.go`
- `internal/logging/logger.go`
- `internal/logging/middleware.go`
- `internal/metrics/metrics.go`
- `internal/metrics/middleware.go`
- `internal/tracing/tracing.go`

Read this module:

- [modules/09-observability/README.md](./modules/09-observability/README.md)
- [modules/09-observability/SURFACE.md](./modules/09-observability/SURFACE.md)

## OPSL.10 Graceful Shutdown and Deployment

What you build: safe drain behavior, deployment packaging, and final integrated system proof.

Target proof surface once this module is implemented:

- `go build` succeeds for the Opslane server
- the full Opslane test suite stays green
- drain behavior is verified under shutdown pressure

Required files:

- `cmd/server/shutdown.go`
- repository root `.github/workflows/ci.yml`

Read this before starting:

- [modules/10-shutdown-deploy/README.md](./modules/10-shutdown-deploy/README.md)
