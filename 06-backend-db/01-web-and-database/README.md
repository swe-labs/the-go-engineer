# Stage 06: Web and Database

## Mission

This section teaches the full backend shape for Stage 06: HTTP servers, API contracts, and database-backed application code.

By the end of this section, you should be comfortable reading and writing:

- `net/http` handlers, middleware, request parsing, timeouts, and readiness surfaces
- REST and gRPC APIs with explicit transport and contract choices
- `database/sql` code that opens databases, executes queries, manages transactions, and respects context
- repository-style boundaries that keep application logic away from raw SQL details

## Stage Ownership

This section belongs to [06 Backend, APIs & Databases](../README.md).

## Stage Map

| Track | Entry | Milestone | Focus |
| --- | --- | --- | --- |
| HTTP Servers | [HS.1 net/http basics](./http-servers/1-net-http-basics) | `HS.10` | handler flow, middleware, validation, timeouts, shutdown, probes |
| APIs | [API.1 REST design principles](./apis/1-rest-design-principles) | `API.9` | REST design, versioning, protobuf, gRPC, interceptors, trade-offs |
| Databases | [DB.1 connecting to SQLite](./databases/1-connecting-to-db) | `DB.8` | SQL access, transactions, repositories, query safety, context timeouts |

## Supporting Reference Surfaces

The stage also includes supporting reference directories that deepen the backend story:

- [`http-client`](./http-client) for outbound HTTP dependency patterns
- [`database-migrations`](./database-migrations) for schema evolution workflows
- [`web-masterclass`](./web-masterclass) for larger application-shaped references

These are part of the stage inventory, but the canonical promoted path is the `HS`, `API`, and `DB` track set in `curriculum.v2.json`.

## Suggested Order

1. Start with the HTTP Servers track so request boundaries feel concrete.
2. Move through the Databases track so persistence and transactions stay explicit.
3. Use the APIs track to compare REST and gRPC after the HTTP and data layers are clear.
4. Use the supporting reference surfaces when you want broader backend context or extra drills.

## Stage Milestones

This section has three promoted proof surfaces:

- `HS.10` REST API
- `API.9` gRPC Service
- `DB.8` query timeouts via context, with `DB.6` remaining the repository milestone that proves data-boundary discipline

You are ready for Stage 07 when you can explain:

- why handlers, middleware, and timeouts belong together in real services
- why REST and gRPC are transport choices with trade-offs, not competing fashions
- why transactions, prepared statements, and context-aware database access matter under real load

## Next Step

After Stage 06, continue to [07 Concurrency](../../07-concurrency).
