# Section 10: Web and Database

## Mission

This section introduces the first production-shaped data and HTTP surfaces in the curriculum.

By the end of the live v2 slice, you should be comfortable reading and writing:

- `database/sql` code that opens databases, executes queries, and manages transactions
- repository-style boundaries that keep application logic away from raw SQL details
- context-aware database workflows that are safer under real service load
- small data-access code you can test and reason about without hiding the underlying mechanics

Section 10 also contains additional HTTP and web-application material. Those directories remain
available as legacy reference surfaces, but the current live v2 slice focuses on the databases
track first.

## Who Should Start Here

### Full Path

Start here after completing Section 09 in order.

### Bridge Path

You can move faster if you already understand:

- explicit error handling
- structs and interfaces
- basic file I/O and config work
- why context and cleanup matter in operational code

Even on the bridge path, do not skip `DB.1` and `DB.3`.
They establish the `database/sql` mental model the rest of the slice depends on.

### Targeted Path

If you are here mainly for data access and repository design, the live path for this section is:

- `DB.1` connecting to SQLite
- `DB.2` executing parameterized inserts
- `DB.3` scanning query results safely
- `DB.4` prepared statements with context
- `DB.5` transactions
- `DB.6` repository pattern project

## Current Section Map

| Surface | Status | Entry | Milestone | Focus |
| --- | --- | --- | --- | --- |
| Databases | Live v2 slice | `DB.1` | `DB.6` | `database/sql`, transactions, and repository boundaries |
| HTTP Client | Legacy reference | `http-client/` | `HC.2` | outbound HTTP calls and dependency injection |
| Database Migrations | Legacy reference | `database-migrations/` | `1-embedded-migrations` | schema evolution and embedded migrations |
| Web Masterclass | Legacy reference | `web-masterclass/` | `11-websockets` | routing, templates, auth, CRUD, and live updates |

## Suggested Order

1. Work through the databases track from `DB.1` to `DB.6`.
2. Use the other Section 10 directories as optional reference material while they remain outside
   the live v2 track.
3. After `DB.6`, continue to Section 11 for concurrency.

## Section Milestone

`DB.6` is the current live milestone for this pilot section.

If you can complete it and explain:

- why `database/sql` should stay behind a repository boundary in application code
- why transactions are the safe way to keep multi-step writes consistent
- why parameterized queries and context-aware execution matter in real services

then you are ready to move into concurrent workflows in Section 11.

## Pilot Role In V2

This live v2 slice keeps the current `10-web-and-database` filesystem layout intact while
promoting the databases track into the public v2 map:

- `DB.1` through `DB.5` are the core lessons
- `DB.6` is the milestone project
- `http-client`, `database-migrations`, and `web-masterclass` remain legacy reference surfaces
  until later alpha waves

That keeps Section 10 useful now without pretending the entire mega-section is fully migrated in
one pass.

## Legacy To Pilot Mapping

- `DB.1` through `DB.6` stay in `10-web-and-database/databases/*`
- the databases track is the only live v2-mapped surface in this wave
- the other Section 10 directories remain available, but they are not yet promoted into
  `curriculum.v2.json`

## References

1. [Accessing relational databases](https://go.dev/doc/database/)
2. [Executing transactions](https://go.dev/doc/database/execute-transactions)
3. [Canceling in-progress operations](https://go.dev/doc/database/cancel-operations)

## Next Step

After `DB.6`, continue to [Section 11: Concurrency](../11-concurrency).
