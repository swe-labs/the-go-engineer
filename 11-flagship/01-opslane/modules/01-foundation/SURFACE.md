# OPSL.1 Implemented Code Surface

This module is already implemented in the current tree.

## Primary Code Files

- [`cmd/server/main.go`](../../cmd/server/main.go)
- [`internal/config/config.go`](../../internal/config/config.go)
- [`internal/config/environment.go`](../../internal/config/environment.go)
- [`Dockerfile`](../../Dockerfile)
- [`docker-compose.yml`](../../docker-compose.yml)
- [`.env.example`](../../.env.example)

## Proof Commands

```bash
go test ./11-flagship/01-opslane/internal/config/...
go run ./11-flagship/01-opslane/cmd/server
```

## What To Read Next

If you want the current learner map, go back to [MODULES.md](../../MODULES.md).
