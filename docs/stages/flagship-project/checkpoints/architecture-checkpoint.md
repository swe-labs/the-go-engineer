# Architecture Checkpoint

## Mission

Explain why the flagship seed is split the way it is and identify one architectural seam that could
become a risk if it grows carelessly.

## Type

- flagship checkpoint

## Level

- core

## Prerequisites

- [Foundation checkpoint](./foundation-checkpoint.md)
- [7 Architecture](../../07-architecture.md)
- [rubric and checkpoint template](../../../templates/rubric-checkpoint-template.md)

## Task

1. Describe the most important package or service boundaries in the enterprise capstone.
2. Identify one seam that is strong and one seam that could become risky as the system grows.
3. Propose one bounded architectural improvement and explain why it matters now more than a larger
   redesign.

## Evidence

- point to the concrete project boundaries you are describing
- justify why one seam is healthy and another is risky
- explain why the proposed change is bounded and well-timed

## Rubric

### 1. Correctness

The learner should describe the real system seams accurately.

### 2. Completeness

The checkpoint should cover current boundaries, one risk, and one bounded improvement.

### 3. Boundary Handling

The learner should discuss ownership, package seams, or interface boundaries directly.

### 4. Code Quality

The proposed improvement should improve structure without collapsing into architecture theater.

### 5. Verification Discipline

The explanation should point to evidence in the current flagship seed.

### 6. Explanation Quality

The learner should justify why the selected improvement matters more than a larger speculative
redesign.

## Next Step

Move to the
[Operations checkpoint](./operations-checkpoint.md).
