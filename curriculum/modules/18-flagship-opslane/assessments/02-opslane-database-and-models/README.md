# Opslane database and models

## Mission

Establish the relational data persistence architecture for Opslane, implementing multi-tenant domain models, versioned transactional schema migrations, database connection pool tuning, and tenant-scoped repository access.

## Prerequisites

- `opslane-01`

## Mental Model

In a multi-tenant SaaS backend, database design must enforce tenant boundary isolation at every query and schema level. Rather than separate databases per tenant, Opslane employs a shared-database, shared-schema pattern where every domain entity (`User`, `Order`, `Payment`, `AuditLog`) carries a mandatory `tenant_id` foreign key. Schema evolution is managed deterministically via transactional migrations recorded in a `schema_migrations` history table, and repository methods strictly include `tenant_id` in all predicate clauses (`WHERE tenant_id = ? AND id = ?`).

## Visual Model

```text
┌─────────────────────────────────────────────────────────────┐
│ Application Repository Layer (OrderRepository)             │
│   GetByID(ctx, tenantID, orderID)                           │
└──────────────────────────────┬──────────────────────────────┘
                               │ SELECT ... WHERE tenant_id = ? AND id = ?
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ Connection Pool (*sql.DB)                                   │
│   - MaxOpenConns: 25 | MaxIdleConns: 10                     │
│   - ConnMaxLifetime: 15m | ConnMaxIdleTime: 5m              │
└──────────────────────────────┬──────────────────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ PostgreSQL / SQLite Relational Schema                       │
│   ├─► schema_migrations (version, name, applied_at)         │
│   ├─► tenants (id PK, name, tier, created_at)               │
│   ├─► users (id PK, tenant_id FK -> tenants, email, role)   │
│   ├─► orders (id PK, tenant_id FK, user_id FK, total...)    │
│   ├─► payments (id PK, tenant_id FK, idempotency_key UNIQUE)│
│   └─► audit_logs (id PK, tenant_id FK, actor_id, details)   │
└─────────────────────────────────────────────────────────────┘
```

## Machine View

Under the hood:
1. `ConnectDatabase` sets up Go's generic `*sql.DB` connection pool. Setting `MaxOpenConns(25)` protects the database from connection exhaustion under high concurrency. Setting `ConnMaxLifetime(15m)` ensures idle sockets are recycled before cloud firewalls or AWS RDS close them silently.
2. `RunMigrations` executes pending migrations inside an explicit database transaction (`db.Begin()`). If any SQL statement in a migration fails, `tx.Rollback()` is triggered, guaranteeing that partial schemas are never left in an inconsistent state.
3. Foreign keys enforce `ON DELETE CASCADE` for tenant data cleanup.
4. `OrderRepository.GetByID` uses `QueryRowContext` parameterized queries, preventing SQL injection and preventing tenant data leakage by requiring both `tenant_id` and `order_id` in the `WHERE` clause.

## Run Instructions

```bash
# Run database and model tests
go test -v ./curriculum/modules/18-flagship-opslane/assessments/02-opslane-database-and-models

# Run standalone database initialization
go run ./curriculum/modules/18-flagship-opslane/assessments/02-opslane-database-and-models
```

## Code Walkthrough

1. **Domain Models (lines 16–62)**: Declares structs for `Tenant`, `User`, `Order`, `Payment`, and `AuditLog` with JSON mappings and tenant scoping. Note `PasswordHash` is tagged `json:"-"` to prevent accidental credential leakage in API responses.
2. **Schema Migration System (lines 65–152)**: Defines versioned DDL scripts with up/down logic across 3 versions:
   - Version 1: `tenants` and `users` tables with index on `tenant_id`.
   - Version 2: `orders` and `payments` tables with unique constraint on `idempotency_key`.
   - Version 3: `audit_logs` table with composite index `(tenant_id, created_at)`.
3. **Database Connection Tuning (lines 156–177)**: Configures connection pool bounds (`MaxOpenConns`, `MaxIdleConns`, `ConnMaxLifetime`) and validates connectivity via `db.PingContext(ctx)` with timeout.
4. **`RunMigrations()` (lines 180–223)**: Queries `schema_migrations` table; skips applied versions; executes unapplied migrations within an atomic `sql.Tx`.
5. **`OrderRepository` (lines 226–266)**: Implements `Create` and `GetByID` using tenant-scoped SQL queries and `sql.NullString` scanning for nullable fields (`notes`).

## Try It

1. Run the test suite: verify that migrations run cleanly, insert valid orders, and enforce tenant isolation.
2. Attempt to fetch an existing order using an mismatched `tenant_id`: verify `GetByID` returns `nil, nil` (not found).
3. Inspect `schema_migrations`: verify that all 3 versions are recorded.

Common mistakes:
- Forgetting `tenant_id` in repository `WHERE` clauses, creating severe cross-tenant data leakage vulnerabilities.
- Executing migrations outside a transaction, leaving the database broken if a DDL statement fails halfway.
- Scanning nullable SQL columns into non-pointer Go strings, causing runtime scan errors.

## In Production

At enterprise scale, multi-tenant relational persistence requires defensive query design. Composite indices such as `(tenant_id, status)` and `(tenant_id, created_at)` ensure queries remain sub-millisecond even as tables grow to tens of millions of rows. Connection pool limits must be sized appropriately according to database server memory and CPU core count.

## Thinking Questions

1. Why is `tenant_id` included in the unique index or `WHERE` clause of every query even if `id` is a UUID?
2. What are the operational risks of running database schema migrations automatically at application startup in multi-replica deployments?
3. How does `sql.NullString` handle SQL `NULL` values compared to standard Go `string` variables?

## Next Step

Proceed to [Milestone 03 — Authentication](../03-opslane-authentication/README.md) to implement password hashing with bcrypt, JWT token generation, and secure session management.
