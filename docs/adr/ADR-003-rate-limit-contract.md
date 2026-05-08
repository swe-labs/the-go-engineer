# ADR-003: Rate Limit API Contract

- Status: accepted
- Date: 2026-05-07

## Context

Clients need deterministic, machine-readable rate-limit behavior to implement retries and backoff.

## Decision

Standardize response behavior:

- dynamic `X-RateLimit-Limit`
- dynamic `X-RateLimit-Remaining`
- `X-RateLimit-Reset` and `Retry-After` on blocked requests
- JSON error envelope with `error.code` and `error.message`

## Consequences

- consistent client behavior across endpoints
- improved observability and integration testability
