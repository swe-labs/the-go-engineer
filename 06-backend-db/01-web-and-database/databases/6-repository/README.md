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

## Next Step

After you complete this milestone, continue to [Stage 07: Concurrency](../../../../07-concurrency/01-concurrency) or
explore the other Stage 06 legacy reference surfaces if you need more web/database examples first.
