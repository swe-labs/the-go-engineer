# Opslane foundation and configuration

## Mission

Initialize the Opslane backend foundation, configuring environment variable parsing, strict parameter validation, structured logging with `log/slog`, panic recovery middleware, and orchestrator health endpoints.

## Prerequisites

- `core-15-19`

## Mental Model

The foundation and configuration layer provides a single, immutable, and strictly validated source of truth for all runtime dependencies. Rather than scattered `os.Getenv` calls throughout handlers and services, configuration is evaluated once at process startup. If any parameter is missing, out-of-bounds, or violates production security guarantees (such as using default development secrets in production), the application crashes immediately (fail-fast) before opening network ports or accepting client traffic.

## Visual Model

```text
┌─────────────────────────────────────────────────────────────┐
│ OS Environment Variables (PORT, ENV, DB_URL, JWT_SECRET...) │
└──────────────────────────────┬──────────────────────────────┘
                               │ LoadConfig() & Strict Validation
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ Validated Config Struct                                     │
│   - Port (1024-65535)                                       │
│   - Environment ("development", "staging", "production")    │
│   - Security Guard (production rejects default secret)      │
└──────────────────────────────┬──────────────────────────────┘
                               │ NewApplication(cfg)
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ Opslane Application Runtime                                 │
│   ├─► slog.Logger (JSON in Prod / Text in Dev)              │
│   ├─► Panic Recovery Middleware (recoverMiddleware)         │
│   ├─► GET /healthz (status, uptime, service)                │
│   └─► GET /readyz (workers, readiness state)                │
└─────────────────────────────────────────────────────────────┘
```

## Machine View

Under the hood:
1. `LoadConfig` evaluates `os.Getenv` against sensible local development defaults. Numeric conversions (`strconv.Atoi`) ensure ports lie between 1024 and 65535, preventing binding to privileged root ports or overflow values.
2. In production (`ENV=production`), an explicit security invariant check runs: if `JWTSecret` matches the dev fallback, startup aborts with `os.Exit(1)`.
3. `InitLogger` binds a typed `slog.Handler`. Production uses `slog.NewJSONHandler` targeting stdout for collector daemon ingestion (FluentBit, Vector). Development uses `slog.NewTextHandler` for terminal ergonomics.
4. `recoverMiddleware` intercepts goroutine panics in HTTP request chains using `defer recover()`, preventing process termination while logging structured error details and returning HTTP 500 JSON payloads.

## Run Instructions

```bash
# Run milestone tests
go test -v ./curriculum/modules/18-flagship-opslane/assessments/01-opslane-foundation-and-configuration

# Execute with environment overrides
PORT=9090 ENV=staging go run ./curriculum/modules/18-flagship-opslane/assessments/01-opslane-foundation-and-configuration
```

## Code Walkthrough

1. **`Config` Struct (lines 16–24)**: Models all operational parameters (`Port`, `Environment`, `DatabaseURL`, `JWTSecret`, `WorkerCount`, `LogLevel`, `RequestTimeout`).
2. **`LoadConfig()` (lines 27–93)**:
   - Sets safe local defaults (port `8080`, dev database URL, dev secret).
   - Overrides from environment variables with bounds checking (`port >= 1024 && port <= 65535`, `WorkerCount >= 1 && WorkerCount <= 64`).
   - Enforces cryptographic secret length: `JWT_SECRET` must be $\ge 16$ characters.
   - Enforces production invariant: rejects default secret in `production`.
3. **`InitLogger()` (lines 96–118)**: Sets up structured logging level and selects `NewJSONHandler` for production or `NewTextHandler` for dev.
4. **`Application` & Composition Root (lines 121–137)**: Initializes the dependency container and records process `StartTime`.
5. **`Routes()` & Health Probes (lines 140–186)**: Mounts `GET /healthz` (service name, uptime, env) and `GET /readyz` (worker count, readiness), wrapped with `recoverMiddleware`.
6. **`main()` (lines 188–213)**: Validates configuration, prints startup banner, and starts `http.Server` with explicit read/write timeouts.

## Try It

1. Run the test suite: verify that defaults pass and production insecure secret triggers an error.
2. Set `PORT=80` and run: observe the validation error rejecting ports below 1024.
3. Set `ENV=production` without setting `JWT_SECRET`: verify fail-fast termination before server binds.

Common mistakes:
- Using `log.Fatal` inside library packages rather than returning errors to `main()`.
- Allowing default fallback secrets to run in production containers.
- Skipping read and write timeouts on `http.Server`, exposing the service to Slowloris resource exhaustion.

## In Production

Production microservices at scale require strictly validated, immutable configuration. Systems like Kubernetes rely on `/healthz` and `/readyz` probes to orchestrate rolling updates and traffic routing. Structuring logging to stdout via `log/slog` adheres to Twelve-Factor App methodologies, allowing centralized log routers to ingest and index records without disk I/O bottlenecks.

## Thinking Questions

1. Why should configuration parsing always happen in `main` or an explicit init function rather than via `init()` packages?
2. How does separating `/healthz` (liveness) from `/readyz` (readiness) prevent cascading pod restarts during transient downstream database saturation?
3. Why are HTTP timeouts (`ReadTimeout`, `WriteTimeout`) critical on internet-facing Go HTTP servers?

## Next Step

Proceed to [Milestone 02 — Database and Models](../02-opslane-database-and-models/README.md) to set up transactional PostgreSQL persistence and tenant-isolated domain models.
