# Project Assessment — Book Catalog REST API

## Purpose

Evaluate your ability to build a production-grade, zero-external-dependency REST API in Go. This assessment verifies that you can implement standard library routing with Go 1.22+ pattern matching, custom middleware chaining (panic recovery, correlation IDs, structured logging), thread-safe in-memory storage (`sync.RWMutex`), standardized JSON envelopes for data and errors, input validation, and graceful server shutdown (`os/signal`).

## Targets

- `project-book-catalog-rest-api`

## Scope

This assessment covers:

- `net/http` routing with Go 1.22+ method-path patterns (`GET /books`, `POST /books`, `GET /books/{id}`)
- Path parameter extraction with `r.PathValue("id")`
- Query parameter parsing for pagination (`limit`, `offset`) and author filtering
- Thread-safe storage with `sync.RWMutex` (read lock for lookups/lists, write lock for mutations)
- Custom middleware pipeline: `RecoveryMiddleware`, `CorrelationIDMiddleware`, and `LoggingMiddleware`
- Standardized response envelopes for successful responses (`{"data": ...}`) and errors (`{"error": {"code": "...", "message": "..."}}`)
- Input validation: 10- and 13-digit ISBN verification, non-empty fields, positive price values
- Server timeout configuration (`ReadTimeout`, `WriteTimeout`, `IdleTimeout`)
- Graceful server shutdown on `SIGINT` / `SIGTERM` using `server.Shutdown(ctx)` with a deadline

## Evidence required

Submit:

1. Your completed `_starter/main.go` with storage, middleware, handlers, and graceful shutdown implemented
2. Output of `go test -v .` showing all unit and integration tests passing
3. A brief explanation answering: How does the panic recovery middleware prevent a single faulty request from crashing the entire server, and why must the status code be written before the response body?

## Time expectation

Plan for 90 to 120 minutes.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Correctness** | 35% | All CRUD endpoints, health check, and error codes match specification; input validation catches bad ISBNs/prices; all automated tests pass | Most endpoints work; minor status code mismatch or partial validation | Missing endpoints or non-functional routing |
| **Middleware & Reliability** | 25% | Panic recovery catches handler panics and returns structured 500 without crashing; correlation IDs propagated via headers and context; structured logging active | Recovery works but logging or correlation ID missing | Unhandled panics crash server |
| **Concurrency & Thread-Safety** | 25% | In-memory store uses `sync.RWMutex` with proper lock acquisition and defer unlock; race detector passes under concurrent requests | Mutex used but potential race condition or unnecessary write locking during reads | Unsynchronized map access causing runtime panics |
| **API Quality & Lifecycle** | 15% | Consistent JSON envelope for data and errors; timeouts configured; graceful shutdown drains in-flight connections cleanly | JSON returned but inconsistent envelope, or abrupt server termination | Non-JSON responses or no signal handling |

**Passing score**: 80%

## Review questions

1. How do you distinguish between a client error (400 Bad Request, 404 Not Found, 422 Unprocessable) versus an internal server error (500) in your HTTP handlers?
2. Walk me through your JSON serialization design: why use `json.NewDecoder(r.Body)` instead of reading all bytes into memory with `io.ReadAll`?
3. How does your middleware chain handle panics, request logging, and correlation ID propagation across HTTP request contexts?
4. What happens to active, in-flight HTTP requests when `server.Shutdown(ctx)` is called upon receiving `SIGINT`?
5. Why is `sync.RWMutex` preferable over a standard `sync.Mutex` for an in-memory catalog with frequent reads and infrequent writes?

## Passing standard

You pass when all test suites pass, CRUD endpoints adhere to the REST specification, middleware prevents server panics, and you can explain the HTTP request lifecycle and context propagation.

## Retake policy

Learner may revise code and resubmit after addressing all rubric gaps.

## NEXT UP

[Module 08 Checkpoint](../checkpoint/README.md)
