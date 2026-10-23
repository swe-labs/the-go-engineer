# Project Assessment — Task API with Postgres

## Purpose

Evaluate your ability to build a transactional task management API backed by SQL persistence using Go's `database/sql` package. This assessment verifies that you can tune database connection pools, execute schema migrations idempotently, write parameterized SQL queries immune to injection, coordinate multi-table transactional status transitions with audit logging (`tx.BeginTx`, `tx.Commit`, `tx.Rollback`), and conduct database disaster recovery drills.

## Targets

- `project-task-api-postgres`

## Scope

This assessment covers:

- Database connection pool tuning: `SetMaxOpenConns`, `SetMaxIdleConns`, and `SetConnMaxLifetime`
- Automated schema migrations applying DDL for `tasks` and `task_audit_logs` idempotently
- Parameterized CRUD queries using placeholder binding (`$1` / `?`) to eliminate SQL injection
- ACID transactions: atomic state updates and audit log inserts with `defer tx.Rollback()` guarantees
- Error handling: mapping `sql.ErrNoRows` to `ErrNotFound` (404) and invalid status transitions to 422
- HTTP REST integration: routing and handlers exposing task lifecycle and audit log queries
- Database disaster recovery protocol: pg_dump backup, database destruction, and schema/data restoration

## Evidence required

Submit:

1. Your completed `_starter/main.go` with connection pooling, migrations, repository CRUD, and transactional status transitions implemented
2. Output of `go test -v .` showing all unit, repository, and transaction tests passing
3. A written explanation answering: Walk me through your database disaster recovery drill: what exact steps did you take to dump, destroy, restore, and verify database integrity?

## Time expectation

Plan for 90 to 120 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Correctness** | 35% | Connection pool configured; migrations run idempotently; CRUD endpoints and status transitions work; all automated tests pass | Most endpoints work; minor status code mismatch or migration edge case | Non-functional repository or migrations fail |
| **Transactional integrity** | 25% | Status update and audit log write execute atomically inside `BeginTx`; `defer tx.Rollback()` guarantees zero partial state on error; `tx.Commit()` executes cleanly | Transaction used but rollback not deferred or audit log committed outside transaction | No transaction used; partial state persists on error |
| **SQL security & quality** | 25% | 100% parameterized queries; zero string concatenation in queries; SQL injection test suite passes; `rows.Close()` and `rows.Err()` checked | Parameterized queries used but minor query formatting flaws | Raw string concatenation in SQL queries creating injection vulnerabilities |
| **Explanation & DR** | 15% | Learner demonstrates disaster recovery procedure and articulates transaction isolation levels and connection pool exhaustion risks | Explains transactions and pooling with minor gaps | Unable to explain transactions or disaster recovery protocol |

**Passing score**: 80%

## Review questions

1. Walk me through your database disaster recovery drill: what exact steps did you take to dump, destroy, restore, and verify database integrity?
2. How do database transactions (`tx.Commit()`, `tx.Rollback()`) ensure consistency when updating tasks with associated audit log entries?
3. Why did you use parameterized queries instead of string concatenation, and what are the specific SQL injection risks prevented?
4. What happens when your connection pool reaches `MaxOpenConns` and a new query arrives? How do you prevent request deadlocks?
5. Why is `defer tx.Rollback()` safe even when `tx.Commit()` succeeds at the end of the transaction?

## Passing standard

You pass when all tests pass, transactions guarantee atomicity without dangling locks, queries are strictly parameterized, and you can explain connection pooling mechanics and DR protocols.

## Retake policy

Learner may revise code and resubmit after addressing all rubric gaps.

## NEXT UP

[Module 09 Checkpoint](../checkpoint/README.md)
