# ADR-002: Migration Governance

- Status: accepted
- Date: 2026-05-07

## Context

Schema drift and out-of-order migration execution are common production failure modes.

## Decision

Opslane migration runner enforces:

- strict sequential migration ordering
- checksum capture and verification for applied migrations
- dirty-state detection gate before applying new migrations

## Consequences

- safer releases with earlier drift detection
- clearer operator expectations for rollback and recovery
