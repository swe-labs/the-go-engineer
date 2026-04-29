# OPSL.7 Implemented Code Surface

This module is now implemented in the current tree.

## Primary Code Files

- [`internal/events/types.go`](../../internal/events/types.go)
- [`internal/events/bus.go`](../../internal/events/bus.go)
- [`internal/workers/pool.go`](../../internal/workers/pool.go)
- [`internal/workers/order_processor.go`](../../internal/workers/order_processor.go)
- [`internal/workers/payment_processor.go`](../../internal/workers/payment_processor.go)
- [`internal/workers/notification_worker.go`](../../internal/workers/notification_worker.go)

## Proof Commands

```bash
go test ./11-flagship/01-opslane/internal/events/...
go test ./11-flagship/01-opslane/internal/workers/...
go run ./11-flagship/01-opslane/scripts/progress.go
```

## What This Proves

- async work has bounded queue capacity
- queue saturation returns explicit errors
- worker pools drain submitted work on stop
- processor adapters call order, payment, and notification seams
- handler failures travel through an explicit error callback instead of disappearing

## What To Read Next

If you want the current learner map, go back to [MODULES.md](../../MODULES.md).
