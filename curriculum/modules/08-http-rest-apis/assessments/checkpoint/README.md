# Module 08 Assessment — HTTP and REST APIs

## Purpose

Verify that the learner can independently apply and explain HTTP handlers, routing, middleware, REST API design, JSON encoding/decoding, client requests, and gRPC basics.

## Targets

- `core-08-01`
- `core-08-02`
- `core-08-03`
- `core-08-04`
- `core-08-05`
- `core-08-06`
- `core-08-07`
- `core-08-08`
- `core-08-09`
- `core-08-10`
- `core-08-11`
- `core-08-12`
- `core-08-13`
- `core-08-14`
- `core-08-15`
- `core-08-16`
- `core-08-17`
- `core-08-18`
- `core-08-19`
- `core-08-20`
- `core-08-21`

## Scope

This checkpoint covers:

- HTTP protocol mechanics: request/response lifecycle, status codes, headers, and streaming bodies
- Standard library routing: `http.ServeMux` with Go 1.22+ method-path patterns and path values
- Middleware composition: onion-style wrapping, request interception, and panic recovery
- Context propagation: request-scoped deadlines, cancellation, and metadata injection
- JSON serialization: stream encoding/decoding with `json.NewDecoder`, envelope schemas, and payload validation
- Production resilience: connection pooling, client timeouts, and graceful server draining with `os/signal`
- RPC architectures: gRPC unary and streaming contracts, Protocol Buffers schema definition, and HTTP/2 multiplexing

## Evidence required

Submit:

1. **HTTP API implementation**: Multi-endpoint REST service featuring CRUD handlers, standardized JSON data/error envelopes, and input validation
2. **Middleware chain**: Composable pipeline implementing panic recovery, correlation ID header injection, and structured request duration logging
3. **HTTP client consumer**: Client program executing requests against the API with explicit timeouts and structured error handling
4. **Written explanations**: Answers to the five review questions below

## Time expectation

Plan for 2 to 3 hours.

## Rubric

| Criterion | Weight | Excellent (5) | Good (3) | Needs Work (0) |
|-----------|--------|--------------|----------|----------------|
| **Conceptual clarity** | 30% | Clearly explains HTTP request lifecycle, handler execution, middleware wrapping, context propagation, and gRPC/REST trade-offs without memorized magic | Explains basics of handlers, routing, and middleware with minor gaps | Cannot explain HTTP request dispatch, middleware wrapping, or server lifecycle |
| **Implementation** | 35% | Complete, clean implementation with proper Go idioms, error envelopes, correct HTTP status codes, and no resource leaks | Working implementation with minor status code mismatches or unhandled edge cases | Incomplete or non-functional API; panics crash the server |
| **Proof** | 25% | Strong test suite verifying routes, status codes, middleware recovery, and client timeout behavior | Basic tests covering happy path only | No test evidence or verification |
| **Communication** | 10% | Professional explanation linking HTTP and REST design to system reliability, scalability, and production operations | Basic explanation with some clarity | Unable to articulate design trade-offs or decisions |

**Passing score**: 80%

## Review questions

1. Walk me through what happens during an HTTP request lifecycle in Go's `net/http` from `ListenAndServe` to response write.
2. Explain middleware chaining. How does wrapping `http.Handler` work, and what are the correct patterns for ordering middleware?
3. What is a graceful shutdown? Show me the signal handling code and explain what happens to in-flight requests.
4. When should a backend engineer choose gRPC over REST/JSON, and what are the operational trade-offs?
5. Why is it dangerous to use `http.DefaultClient` or `http.DefaultServeMux` in production services?

## Passing standard

Score of 80% or higher across the four rubric criteria, with all required evidence submitted.

## Retake policy

Learner may retry after fixing rubric gaps and writing a short correction note.

## NEXT UP

[Module 09: SQL and PostgreSQL Persistence](../../09-sql-postgres-persistence/README.md)
