# OPSL.10 Implemented Code Surface

This module is fully implemented in the current tree.

## Primary Code Files

- [`cmd/server/shutdown.go`](../../cmd/server/shutdown.go)
- [`cmd/server/main.go`](../../cmd/server/main.go)
- [`internal/handlers/handlers.go`](../../internal/handlers/handlers.go)

## Proof Commands

```bash
go build ./11-flagship/01-opslane/cmd/server
go test ./11-flagship/01-opslane/...
go run ./11-flagship/01-opslane/scripts/progress.go
```

## What This Proves

- `setupGracefulShutdown` listens for `SIGINT` and `SIGTERM`.
- `/health` responds with HTTP 503 (`status: "draining"`) immediately upon receiving a shutdown signal to prevent load balancers from sending new requests.
- `server.Shutdown` drains in-flight HTTP requests.
- `bus.Close` prevents new jobs from entering the asynchronous queues.
- `pool.Stop` waits for background workers (orders and payments) to finish processing their buffered work.
- The `main` goroutine blocks until all the above is complete before dropping the database connection and exiting.

## What To Read Next

If you want the current learner map, go back to [MODULES.md](../../MODULES.md). This is the final module in the flagship path.
