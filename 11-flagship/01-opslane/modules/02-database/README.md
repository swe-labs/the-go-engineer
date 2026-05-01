# OPSL.2 Database and Models

## Mission

Define the multi-tenant data model and persistence seams before adding more public behavior.

## What This Module Builds

- tenant-aware domain models
- PostgreSQL schema creation
- repository boundaries
- transaction seams for later business workflows

## You Are Here If

- you can explain why tenant ownership is stored directly in the schema
- you understand where repository methods stop and higher-level workflows should begin
- you can read the current migrations without guessing how orders and payments relate

## Proof Surface

```bash
go test ./11-flagship/01-opslane/internal/db/...
```

## Required Files and Boundaries

- `internal/db/migrations.go`
- `internal/db/repository.go`
- `internal/models/tenant.go`
- `internal/models/user.go`
- `internal/models/order.go`
- `internal/models/payment.go`

Implemented code surface: [SURFACE.md](./SURFACE.md)

The repository layer owns persistence concerns. It should not absorb auth rules or HTTP parsing.

## Engineering Questions

- Which relationships must be tenant-scoped in the database itself, not just in handlers?
- Where do transactions belong when multiple rows must change together?
- What query patterns deserve indexes now versus later?

## Next Step

Next: `OPSL.3` -> `11-flagship/01-opslane/modules/03-auth`

Open `11-flagship/01-opslane/modules/03-auth/README.md` to continue.
