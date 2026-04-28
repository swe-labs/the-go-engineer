# DB.6 Repository Pattern Project

## Mission

Build a small SQLite-backed user repository that proves application logic can depend on behavior
contracts instead of raw database calls.

This exercise is the live Stage 06 milestone.
It is where connections, queries, scans, prepared statements, transactions, and interfaces come
together in one runnable artifact.

## Prerequisites

Complete these first:

- `DB.1` connecting to SQLite
- `DB.2` queries and parameterization
- `DB.3` select queries
- `DB.4` prepared statements with context
- `DB.5` transactions

## What You Will Build

Implement a repository-driven user flow that:

1. defines a `UserRepository` interface
2. backs that interface with a SQLite implementation
3. creates a user and profile inside one transaction
4. queries users back into domain models
5. keeps SQL details away from the calling application code
6. proves core behavior with tests

## Files

- [main.go](./main.go): complete runnable solution
- [main_test.go](./main_test.go): tests for the repository contract
- [repository/user.go](./repository/user.go): concrete repository implementation
- [models/user.go](./models/user.go): domain models used by the repository
- [_starter/main.go](./_starter/main.go): starter surface with requirements and TODOs

## Run Instructions

Run the completed solution:

```bash
go run ./06-backend-db/01-web-and-database/databases/6-repository
```

Run the tests:

```bash
go test ./06-backend-db/01-web-and-database/databases/6-repository
```

Run the starter:

```bash
go run ./06-backend-db/01-web-and-database/databases/6-repository/_starter
```

Note: the current repository lesson uses `github.com/mattn/go-sqlite3`, so local run/test commands
need a CGO-enabled Go toolchain.

## Success Criteria

Your finished solution should:

- keep write operations transactional when multiple tables are involved
- use parameterized SQL instead of string-built queries
- map rows into Go models predictably
- expose a repository interface instead of leaking `*sql.DB` into higher-level logic
- pass the provided tests

## Common Failure Modes

- letting application code talk directly to `*sql.DB` instead of the repository contract
- forgetting that `rows.Close()` is part of the resource-management contract
- splitting multi-table writes into separate non-transactional calls
- treating the repository pattern as an excuse to hide all SQL behavior instead of explaining it


## Mental Model

Think of this as the conceptual blueprint. The components interact by exchanging state, defining clear boundaries between what is requested and what is provided.

## Visual Model

Visualizing this process involves tracing the execution path from the input entry point, through the processing layers, and out to the final output or side effect.

## Machine View

At the hardware level, this translates into specific memory allocations, CPU instruction cycles, and OS-level system calls to manage resources efficiently.

## Solution Walkthrough

The solution demonstrates a complete implementation, proving the concept by bridging the individual requirements into a single, cohesive executable.

## Try It

Run the code locally. Modify the inputs, toggle the conditions, and observe how the output shifts. Experimentation is the fastest way to cement your understanding.

## Verification Surface

The correctness of this component is proven by its associated test suite. We verify boundaries, handle edge cases, and ensure performance constraints are met.

## In Production

The repository pattern is the standard database access boundary in production Go services. Without it, SQL queries leak into HTTP handlers, business logic becomes untestable without a running database, and schema changes ripple through every layer of the application. In production, the repository interface enables teams to swap database implementations (SQLite in tests, PostgreSQL in production), add caching layers transparently, and mock database behavior in unit tests without touching real infrastructure. The transaction pattern shown in this exercise — grouping related writes into a single atomic operation — prevents data inconsistencies that are nearly impossible to debug in production. A common production failure mode is forgetting to call `rows.Close()`, which leaks database connections from the pool until the service runs out and every request starts timing out. Teams that treat the repository as a first-class boundary — not just a convenience wrapper — gain the ability to reason about database behavior independently from business logic, which is critical when debugging production incidents under pressure.

## Thinking Questions

1. Why is it important to keep SQL inside the repository layer instead of letting handlers build queries directly?
2. If a transaction commits the user but fails on the profile insert, what state does the database end up in, and why does wrapping both in one transaction prevent this?
3. How would you test repository behavior without requiring a real database connection in CI?
4. What happens to in-flight queries when the database connection pool is exhausted, and how would you detect this in production?

## Next Step

After you complete this milestone, continue to [Stage 07: Concurrency](../../../../07-concurrency/01-concurrency) or
explore the other Stage 06 legacy reference surfaces if you need more web/database examples first.

