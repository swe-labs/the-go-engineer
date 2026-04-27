# OPSL.6 Implemented Code Surface

This module is now implemented in the current tree.

## Primary Code Files

- [`internal/payment/gateway.go`](../../internal/payment/gateway.go)
- [`internal/payment/worker.go`](../../internal/payment/worker.go)
- [`internal/services/payment.go`](../../internal/services/payment.go)
- [`internal/handlers/api.go`](../../internal/handlers/api.go)
- [`internal/handlers/handlers.go`](../../internal/handlers/handlers.go)
- [`internal/db/repository.go`](../../internal/db/repository.go)

## Proof Commands

```bash
go test ./11-flagship/01-opslane/internal/payment/...
go test ./11-flagship/01-opslane/internal/services/...
go test ./11-flagship/01-opslane/internal/handlers/...
go run ./11-flagship/01-opslane/scripts/progress.go
```

## What This Proves

- payment gateway behavior is behind an explicit interface
- worker processing has a bounded queue boundary
- duplicate `provider_reference` values do not double-charge
- timeout outcomes stay pending for reconciliation
- settled and failed payments move the order workflow intentionally

## What To Read Next

If you want the current learner map, go back to [MODULES.md](../../MODULES.md).
