# Stage 06 Track: Databases

## Mission

This track teaches you how Go talks to relational databases without hiding the core `database/sql`
contracts.

By the end of the live track, you should be comfortable with:

- blank imports and driver registration
- the difference between `sql.Open` and `db.Ping`
- parameterized queries and scanning rows safely
- prepared statements and context-aware execution
- transactions for multi-step writes
- repository boundaries that keep application code decoupled from raw SQL

## Engineering Depth

The `database/sql` package is not an ORM.
It is a connection-pool manager plus a set of low-level query contracts.

That means the real engineering hazards are not "how do I write SQL in Go?" but:

- leaking `rows` and exhausting the connection pool
- holding transactions open while doing slow non-database work
- scattering raw `*sql.DB` calls through handlers and business logic
- confusing a clean "not found" result with an operational failure

## Learning Path

| ID | Lesson | Concept | Requires |
| --- | --- | --- | --- |
| `DB.1` | [connecting](./1-connecting-to-db) | blank imports, `sql.Open`, `db.Ping`, connection pools | entry |
| `DB.2` | [query - INSERT](./2-query) | `db.Exec`, `?` parameters, `LastInsertId`, bcrypt | `DB.1` |
| `DB.3` | [query - SELECT](./3-select) | `QueryRow`, `Query`, `rows.Scan`, `rows.Close`, `rows.Err` | `DB.1`, `DB.2` |
| `DB.4` | [prepared statements](./4-prepare) | `db.Prepare`, `ExecContext`, statement reuse | `DB.2`, `DB.3` |
| `DB.5` | [transactions](./5-transactions) | `BeginTx`, `defer Rollback`, `Commit`, ACID consistency | `DB.1`, `DB.2`, `DB.3` |
| `DB.6` | [repository pattern project](./6-repository) | interface-driven data access, transactions, model mapping | `DB.1`, `DB.2`, `DB.3`, `DB.4`, `DB.5` |

## Suggested Order

1. Work through `DB.1` to `DB.5` in order.
2. Complete `DB.6` as the live milestone.
3. Return to the Stage 06 overview when you are ready to explore the legacy reference surfaces.

## Live Milestone

`DB.6` is the current live milestone for this track.

It is the point where the lower-level lessons stop being isolated mechanics and start looking like
service-ready data access code.

## References

1. [Accessing relational databases](https://go.dev/doc/database/)
2. [Executing transactions](https://go.dev/doc/database/execute-transactions)
3. [Canceling in-progress operations](https://go.dev/doc/database/cancel-operations)

## Next Step

After `DB.6`, continue to the [Stage 06 overview](../README.md) or move on to
[Stage 07: Concurrency](../../../07-concurrency/01-concurrency).
