# Task API with Postgres

## Learning objective

Build a robust, transactional task management service backed by SQL persistence using Go's `database/sql` package. Master connection pool tuning, automated schema migrations, parameterized queries immune to SQL injection, and multi-table transactional status transitions with audit logging.

## Why this matters

Almost every scalable backend relies on relational persistence. Poor connection management exhausts database connections, while lack of transactions causes inconsistent partial state when errors occur mid-operation. Mastering `database/sql` pooling, error checking, and transactional atomicity (`tx.BeginTx`, `tx.Commit`, `tx.Rollback`) is essential for any production backend engineer.

## Prerequisites

Complete all Module 09 lessons (01–26) covering SQL fundamentals, database drivers, connection pools, CRUD operations, transactions, and migration strategies.

## Available concepts

This project draws on `database/sql`, PostgreSQL / SQLite drivers, connection pool configuration (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`), prepared queries, parameterized arguments, ACID transactions, disaster recovery, and automated schema migration.

## Project overview

The `_starter/` directory contains HTTP handlers, error models, and interface skeletons. Your job is to implement the database repository methods, connection pool tuning, and transactional update logic so that:

- Database connections are properly pooled and bounded to prevent exhaustion.
- Automated migrations initialize `tasks` and `task_audit_logs` tables.
- All SQL operations use parameter binding to eliminate SQL injection vulnerabilities.
- Status changes execute atomically inside a single transaction with audit log records.
- Failures during status transitions trigger rollbacks with zero residual state changes.

A runnable `docker-compose.yml` provides a local PostgreSQL instance, while SQLite in-memory is supported out-of-the-box for instant local testing.

## Database Schema

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(64) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status VARCHAR(32) NOT NULL DEFAULT 'TODO',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS task_audit_logs (
    id VARCHAR(64) PRIMARY KEY,
    task_id VARCHAR(64) NOT NULL,
    from_status VARCHAR(32) NOT NULL,
    to_status VARCHAR(32) NOT NULL,
    changed_at TIMESTAMP NOT NULL,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);
```

## API Endpoints

| Method | Path | Status | Description |
|---|---|---|---|
| `GET` | `/healthz` | `200 OK` / `503 Service Unavailable` | Liveness and database ping check |
| `POST` | `/tasks` | `201 Created` | Create a new task in `TODO` status |
| `GET` | `/tasks` | `200 OK` | List tasks with optional `status` filter and pagination |
| `GET` | `/tasks/{id}` | `200 OK` / `404 Not Found` | Retrieve a task by ID |
| `PATCH` | `/tasks/{id}/status` | `200 OK` / `422 Unprocessable` | Atomically update status and log audit record |
| `GET` | `/tasks/{id}/audit` | `200 OK` | Retrieve history of status transitions |
| `DELETE` | `/tasks/{id}` | `204 No Content` / `404 Not Found` | Delete a task |

## Tasks

### Task 1 — Connection pool & migration runner

In `_starter/main.go`:
- Implement `ConfigurePool(db *sql.DB)`: set max open connections to 25, max idle connections to 10, and connection max lifetime to 15 minutes.
- Implement `RunMigrations(ctx context.Context, db *sql.DB) error`: apply the DDL schema idempotently.

### Task 2 — Implement parameterized repository CRUD

Implement:
- `Create(ctx context.Context, title, description string) (Task, error)`
- `GetByID(ctx context.Context, id string) (Task, error)`
- `List(ctx context.Context, status TaskStatus, limit, offset int) ([]Task, int, error)`
- `Delete(ctx context.Context, id string) error`

*Requirement:* Use `Rebind(query, driver)` and query arguments (`$1` or `?`). Never concatenate raw strings into queries.

### Task 3 — Transactional status update with audit logging

Implement `UpdateStatusTx(ctx context.Context, id string, newStatus TaskStatus) (Task, error)`:
1. Begin a transaction using `db.BeginTx(ctx, nil)`.
2. Ensure `defer tx.Rollback()` is registered immediately.
3. Query the current task status. Return `ErrNotFound` if task does not exist.
4. Validate `newStatus` against allowed values (`TODO`, `IN_PROGRESS`, `DONE`, `ARCHIVED`).
5. Update `status` and `updated_at` in `tasks`.
6. Insert an audit record into `task_audit_logs` tracking `from_status`, `to_status`, and `changed_at`.
7. Commit the transaction with `tx.Commit()`.

## Verification

Run all automated unit and database tests:

```bash
go test -v ./curriculum/modules/09-sql-postgres-persistence/projects/task-api-with-postgres/_starter
```

Expected output:
```
=== RUN   TestPoolConfiguration
--- PASS: TestPoolConfiguration (0.00s)
=== RUN   TestTaskRepository_CreateAndGet
--- PASS: TestTaskRepository_CreateAndGet (0.00s)
=== RUN   TestSQLInjectionResistance
--- PASS: TestSQLInjectionResistance (0.00s)
=== RUN   TestUpdateStatusTx_SuccessAndAudit
--- PASS: TestUpdateStatusTx_SuccessAndAudit (0.00s)
=== RUN   TestUpdateStatusTx_RollbackOnInvalidStatus
--- PASS: TestUpdateStatusTx_RollbackOnInvalidStatus (0.00s)
=== RUN   TestHTTPAPI_EndToEnd
--- PASS: TestHTTPAPI_EndToEnd (0.00s)
PASS
```

Optionally launch local PostgreSQL via Docker:
```bash
docker compose -f ./curriculum/modules/09-sql-postgres-persistence/projects/task-api-with-postgres/_solution/docker-compose.yml up -d
DATABASE_URL="postgres://postgres:postgrespassword@localhost:5432/taskdb?sslmode=disable" go run ./curriculum/modules/09-sql-postgres-persistence/projects/task-api-with-postgres/_solution
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Transactional Integrity** | Status update and audit logging are strictly atomic; rollback cleanly isolates errors | Transactions used but rollback omitted or partial commits possible | Non-transactional updates resulting in orphaned state |
| **SQL Safety & Injection Defense** | All queries parameterized; SQL injection attempts proven harmless in tests | Most queries parameterized with occasional raw formatting | String concatenation in SQL clauses |
| **Connection Pooling** | Pool explicitly tuned (MaxOpenConns, MaxIdleConns, Lifetime) with health checks | Pool default limits retained without tuning | Connections leaked by unclosed rows or statements |
| **Code Structure & Architecture** | Clean separation of HTTP handlers, repository abstractions, and migration logic | Basic separation with logic leaks into handlers | Handlers directly executing raw SQL |

## NEXT UP

[Secure Authenticated API](../10-auth-security/projects/secure-authenticated-api/README.md)
