# ADR-001: Tracing Strategy

- Status: accepted
- Date: 2026-05-07

## Context

Opslane needs distributed trace context and export behavior that is production-shaped while still readable for learners.

## Decision

Use an internal tracer boundary that:

- parses W3C `traceparent` when present
- generates cryptographically random trace/span IDs
- exports spans to OTLP JSON endpoint
- supports configurable sampling via environment

## Consequences

- better interoperability with external tracing stacks
- explicit educational surface for propagation and exporter failure behavior
