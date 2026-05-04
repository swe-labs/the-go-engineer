# OPSL.2 Implemented Code Surface

This module is already implemented in the current tree.

## Primary Code Files

- [`internal/db/migrations.go`](../../internal/db/migrations.go) (embedded)
- [`migrations/`](../../migrations/) (formal SQL migrations)
- [`internal/db/repository.go`](../../internal/db/repository.go)
- [`internal/models/tenant.go`](../../internal/models/tenant.go)
- [`internal/models/user.go`](../../internal/models/user.go)
- [`internal/models/order.go`](../../internal/models/order.go)
- [`internal/models/payment.go`](../../internal/models/payment.go)

## Proof Commands

```bash
go test ./11-flagship/01-opslane/internal/db/...
```

## What To Read Next

If you want the current learner map, go back to [MODULES.md](../../MODULES.md).
