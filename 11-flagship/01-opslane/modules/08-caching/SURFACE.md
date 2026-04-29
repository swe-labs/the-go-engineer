# OPSL.8 Implemented Code Surface

This module is now implemented in the current tree.

## Primary Code Files

- [`internal/cache/cache.go`](../../internal/cache/cache.go)
- [`internal/cache/store.go`](../../internal/cache/store.go)
- [`internal/middleware/cache.go`](../../internal/middleware/cache.go)

## Proof Commands

```bash
go test ./11-flagship/01-opslane/internal/cache/...
go run ./11-flagship/01-opslane/scripts/progress.go
```

## What This Proves

- in-memory cache respects bounded capacity with insert-order eviction
- TTL expiry is enforced on reads and cleaned by a background janitor
- cached values are copied on read and write to prevent mutation
- explicit invalidation removes stale data after order and payment writes
- prefix-based deletion clears tenant-scoped cache groups atomically
- singleflight deduplicates concurrent loads for the same key
- HTTP Cache-Control middleware sets correct headers for public and authenticated endpoints
- NoopCache provides a safe zero-allocation fallback

## What To Read Next

If you want the current learner map, go back to [MODULES.md](../../MODULES.md).
