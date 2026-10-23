# Book Catalog REST API

## Learning objective

Build a production-grade, zero-external-dependency REST API in Go for managing a book catalog. Master `net/http` routing with Go 1.22+ pattern matching, custom middleware chaining, JSON request/response envelope handling, input validation, and graceful server shutdown.

## Why this matters

The vast majority of real-world Go backend services are HTTP web APIs. Writing high-throughput, maintainable HTTP services requires mastering context propagation, structured error responses, defense against handler panics, and connection draining during deployments.

## Prerequisites

Complete all Module 08 lessons (01–20) on HTTP fundamentals, standard library routing, middleware patterns, JSON streaming, and server lifecycle.

## Available concepts

This project draws on `net/http`, `http.ServeMux` pattern matching, `http.Handler` interface, custom middleware wrappers, `log/slog` structured logging, `context.Context` cancellation and value injection, JSON serialization with `encoding/json`, input sanitization, and graceful shutdown with `os/signal`.

## Project overview

The `_starter/` directory contains scaffolding with data models, error definitions, and handler stubs. Your job is to implement the storage methods, middleware functions, and HTTP handlers so that:

- The service compiles and runs cleanly on port `8080` (or `PORT` environment variable).
- All unit and integration tests pass.
- Panics within HTTP handlers are caught and returned as structured 500 errors without crashing the process.
- Requests carry correlation IDs via `X-Correlation-ID` response headers and context.
- Graceful shutdown handles `SIGINT` / `SIGTERM` within a 5-second deadline.

The `_solution/` directory provides a verified reference implementation.

## API Specification

| Method | Path | Status | Description |
|---|---|---|---|
| `GET` | `/healthz` | `200 OK` | Liveness health probe returning `{"data": {"status": "ok"}}` |
| `POST` | `/books` | `201 Created` | Create a book with validation (`Title`, `Author`, `ISBN`, `PriceCents` > 0) |
| `GET` | `/books` | `200 OK` | List books with author filtering and pagination (`limit`, `offset`) |
| `GET` | `/books/{id}` | `200 OK` / `404 Not Found` | Retrieve a single book by ID |
| `PUT` | `/books/{id}` | `200 OK` / `404 Not Found` | Update book title, author, and price |
| `DELETE` | `/books/{id}` | `204 No Content` / `404 Not Found` | Delete book from catalog |

### Error Envelope Format

All non-2xx responses must return a standardized JSON error envelope:

```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "title is required"
  }
}
```

## Tasks

### Task 1 — Implement thread-safe storage & validation

In `_starter/main.go`, implement:
- `ValidateISBN(isbn string) bool` — validates 10- or 13-digit numeric ISBNs.
- `BookStore.Create`, `BookStore.Get`, `BookStore.List`, `BookStore.Update`, `BookStore.Delete` using `sync.RWMutex` to guarantee safe concurrent reads and exclusive writes.

### Task 2 — Implement HTTP handlers & JSON envelope

Implement the HTTP handlers:
- `handleCreateBook`: decode JSON, validate fields, return `201 Created` with `{"data": book}`.
- `handleListBooks`: parse `author`, `limit`, and `offset` query parameters. Return `{"data": books, "meta": {"total": N, "limit": L, "offset": O}}`.
- `handleGetBook`, `handleUpdateBook`, `handleDeleteBook`: parse `r.PathValue("id")` and handle `ErrNotFound`.

### Task 3 — Implement production middleware pipeline

Implement the middleware wrappers:
1. `RecoveryMiddleware`: recover from panics, log error with `slog`, write structured 500 JSON.
2. `CorrelationIDMiddleware`: preserve incoming `X-Correlation-ID` or generate new hex ID, set on response header and inject into `r.Context()`.
3. `LoggingMiddleware`: log request method, path, duration, and correlation ID.

### Task 4 — Implement graceful shutdown

In `main()`, set explicit server timeouts (`ReadTimeout: 5s`, `WriteTimeout: 10s`, `IdleTimeout: 60s`). Listen for `SIGINT` / `SIGTERM` signals and invoke `server.Shutdown(ctx)` with a 5-second deadline context.

## Verification

Run all automated tests from `_starter/`:

```bash
go test -v ./curriculum/modules/08-http-rest-apis/projects/book-catalog-rest-api/_starter
```

Expected output:
```
=== RUN   TestHealthz
--- PASS: TestHealthz (0.00s)
=== RUN   TestCreateBook_Success
--- PASS: TestCreateBook_Success (0.00s)
=== RUN   TestCreateBook_ValidationErrors
--- PASS: TestCreateBook_ValidationErrors (0.00s)
=== RUN   TestGetBook_FoundAndNotFound
--- PASS: TestGetBook_FoundAndNotFound (0.00s)
=== RUN   TestListBooks_PaginationAndFilter
--- PASS: TestListBooks_PaginationAndFilter (0.00s)
=== RUN   TestUpdateAndDeleteBook
--- PASS: TestUpdateAndDeleteBook (0.00s)
=== RUN   TestCorrelationIDMiddleware
--- PASS: TestCorrelationIDMiddleware (0.00s)
=== RUN   TestPanicRecoveryMiddleware
--- PASS: TestPanicRecoveryMiddleware (0.00s)
PASS
```

Run the server locally:
```bash
go run ./curriculum/modules/08-http-rest-apis/projects/book-catalog-rest-api/_solution
```

## Rubric

| Criterion | Excellent (5) | Good (3) | Needs Work (0) |
|---|---|---|---|
| **Correctness** | All CRUD endpoints, health check, and error codes match specification; tests pass | Most endpoints work; minor status code mismatch | Missing endpoints or non-functional routing |
| **Middleware & Reliability** | Panic recovery, correlation ID propagation, and timeouts work seamlessly | Logging or correlation ID missing | Unhandled panics crash server |
| **Concurrency & Thread-Safety** | In-memory store uses `sync.RWMutex` with proper lock acquisition and defer unlock | Mutex used but potential race conditions | Unsynchronized map access causing runtime panics |
| **API Quality** | Consistent JSON envelope for both data and errors; graceful shutdown implemented | JSON returned but inconsistent envelope | Non-JSON responses or abrupt server termination |

## NEXT UP

[Task API with Postgres](../09-sql-postgres-persistence/projects/task-api-with-postgres/README.md)
