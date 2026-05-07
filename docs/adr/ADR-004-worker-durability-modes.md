# ADR-004: Worker Durability Modes

- Status: accepted
- Date: 2026-05-07

## Context

In-memory channels are excellent for teaching concurrency but not durable during process failure.

## Decision

Define explicit durability modes for asynchronous workflows:

- `in_memory` (default educational mode)
- `durable` (future mode using outbox + external queue integration)

## Consequences

- learners understand the tradeoff explicitly
- implementation can evolve without rewriting workflow boundaries
