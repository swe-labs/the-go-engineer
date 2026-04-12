# Operations Checkpoint

## Mission

Explain how the flagship seed is started, observed, and stopped, and identify one operational risk
that matters before feature expansion.

## Type

- flagship checkpoint

## Level

- stretch

## Prerequisites

- [Architecture checkpoint](./architecture-checkpoint.md)
- [8 Production Engineering](../../08-production-engineering.md)
- [rubric and checkpoint template](../../../templates/rubric-checkpoint-template.md)

## Task

1. Explain the current runtime and deployment-facing path for the enterprise capstone.
2. Identify one operational risk involving logs, shutdown, or deployment workflow.
3. Propose one bounded improvement that would make the system safer to operate.

## Evidence

- point to the startup, runtime, or deployment surfaces you are discussing
- explain the operational risk in concrete terms
- justify why the proposed change is worth doing before broader feature work

## Rubric

### 1. Correctness

The learner should describe the current operational path accurately.

### 2. Completeness

The checkpoint should cover runtime flow, one real risk, and one bounded improvement.

### 3. Boundary Handling

The learner should connect the operational risk to lifecycle, observability, or deployment seams.

### 4. Code Quality

The proposed improvement should make the system easier to operate or reason about.

### 5. Verification Discipline

The explanation should rely on visible system evidence instead of guesses about production.

### 6. Explanation Quality

The learner should explain why the chosen operational risk matters before more feature work.

## Next Step

Move to the
[Iteration checkpoint](./iteration-checkpoint.md).
