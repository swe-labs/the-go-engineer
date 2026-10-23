# Module 09 Assessment — SQL and PostgreSQL Persistence

## Purpose

Verify that the learner can independently apply and explain `database/sql`, connection pooling, schema migrations, CRUD queries, ACID transactions, prepared statements, and PostgreSQL-specific persistence features.

## Targets

- `core-09-01`
- `core-09-02`
- `core-09-03`
- `core-09-04`
- `core-09-05`
- `core-09-06`
- `core-09-07`
- `core-09-08`
- `core-09-09`
- `core-09-10`
- `core-09-11`
- `core-09-12`
- `core-09-13`
- `core-09-14`
- `core-09-15`
- `core-09-16`
- `core-09-17`
- `core-09-18`
- `core-09-19`
- `core-09-20`
- `core-09-21`
- `core-09-22`
- `core-09-23`
- `core-09-24`
- `core-09-25`
- `core-09-26`
- `core-09-27`
- `core-09-28`

## Scope

This checkpoint covers:

- Relational data modeling: primary keys, foreign keys, unique constraints, and indexes
- Go database architecture: `database/sql` driver abstraction, connection lifecycle, and DSN configuration
- Connection pool tuning: configuring `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`, and `SetConnMaxIdleTime`
- Schema migration management: versioned DDL migrations, idempotent execution, and rollback strategies
- Safe query execution: parameter binding across dialects (`$1` in PostgreSQL, `?` in SQLite), preventing SQL injection
- Query patterns: `db.Exec`, `db.Query`, and `db.QueryRow` with context deadlines
- Result scanning: mapping database rows to Go structs, handling NULL values with `sql.Null*` types, and checking `rows.Err()`
- ACID transactions: atomic operations with `db.BeginTx`, deferred rollback pattern (`defer tx.Rollback()`), and isolation levels
- Performance optimizations: prepared statement caching and index utilization
- Disaster recovery operations: database dumps, restoration procedures, and data integrity verification

## Evidence required

Submit:

1. **Go database program**: Complete persistence layer implementing CRUD queries with parameterized arguments and context propagation
2. **Schema migration scripts**: Automated migration files with reversible up/down changes applied idempotently
3. **Transaction implementation**: Multi-table transactional operation verifying atomic commit and rollback on simulated failure
4. **Written explanations**: Answers to the five review questions below

## Time expectation

Plan for 2 to 3 hours.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Conceptual clarity** | 30% | Clearly explains connection pooling, query execution, transaction lifecycle, isolation levels, and migration safety without memorized magic | Explains basics of pooling, queries, and transactions with minor gaps | Cannot explain connection pool mechanics, isolation levels, or transaction rollback |
| **Implementation** | 35% | Complete, clean implementation with 100% parameterized queries, proper connection pool limits, deferred rollbacks, and no resource leaks (`rows.Close()`) | Working implementation with minor pooling or scanning flaws | Incomplete persistence code; raw string concatenation in SQL queries |
| **Proof** | 25% | Strong test suite verifying transactional rollback on error, parameter binding, and idempotent migrations | Basic tests covering happy path only | No test evidence or verification |
| **Communication** | 10% | Professional explanation linking SQL design and transaction boundaries to data integrity, disaster recovery, and system scalability | Basic explanation with some clarity | Unable to articulate trade-offs or design decisions |

**Passing score**: 80%

## Review questions

1. Explain the `database/sql` connection pool. What happens when all connections are in use and a new query arrives?
2. Walk me through a SQL transaction with a deferred rollback. What are the isolation levels and when do you need SERIALIZABLE?
3. What are SQL migrations and why do you run them at startup rather than by hand? What makes a migration reversible?
4. What is the difference between `db.Exec` and `db.QueryRow`, and what happens if you forget to call `rows.Close()` on a multi-row query?
5. Show me how SQL injection occurs with string concatenation, and explain how database drivers and prepared statements prevent it at the protocol level.

## Passing standard

Score of 80% or higher across the four rubric criteria, with all required evidence submitted.

## Retake policy

Learner may retry after fixing rubric gaps and writing a short correction note.

## NEXT UP

[Module 10: Authentication and Security](../../10-auth-security/README.md)
