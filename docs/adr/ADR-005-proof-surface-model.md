# ADR-005: Proof-Surface Curriculum Model

- Status: accepted
- Date: 2026-05-07

## Context

Coverage percentage alone is not enough to prove learning or engineering quality.

## Decision

Use a proof-surface model combining:

- behavior tests (unit/integration/race/fuzz)
- curriculum validator checks
- CI drift checks (formatting, module tidy, migration policy)
- learner-visible documentation and checkpoint links

## Consequences

- quality claims remain evidence-backed
- educational and production concerns stay aligned
