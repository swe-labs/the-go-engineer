# Opslane Production Semantics

This document tracks production-oriented semantics layered on top of the teaching implementation.

## Current Implemented Baseline

- idempotency key support on order creation flow
- distributed rate-limit persistence path
- readiness/liveness endpoints (`/health`, `/readyz`, `/livez`)
- OpenAPI baseline at [`./openapi.yaml`](./openapi.yaml)

## Next Semantics Track

- role-based authorization matrix for protected endpoints
- payment idempotency key standardization
- outbox pattern for event delivery guarantees
- durable queue mode for worker processing
- generated client examples from OpenAPI
- uniform structured error response contract
