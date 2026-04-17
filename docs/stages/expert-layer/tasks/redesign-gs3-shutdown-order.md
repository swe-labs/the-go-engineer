# GS.3 Shutdown-Order Redesign

## Mission

Redesign part of the `GS.3` shutdown capstone so service shutdown order becomes clearer and safer
under operational pressure.

## Type

- redesign task

## Level

- production

## Prerequisites

- [10 Production](../../10-production.md)
- `GS.3`
- [rubric and checkpoint template](../../../templates/rubric-checkpoint-template.md)

## Task

1. Inspect the `GS.3` shutdown capstone and identify one shutdown-order or lifecycle boundary that
   could become risky as the system grows.
2. Propose a bounded redesign that improves the shutdown model without rebuilding the whole
   application.
3. Explain why your redesign is better than a simpler but riskier alternative.

## Evidence

- point to the lifecycle seam you want to redesign
- explain the operational risk in concrete terms
- justify why the redesign is bounded, realistic, and worth doing now

## Rubric

### 1. Correctness

The redesign should target a plausible shutdown or lifecycle risk.

### 2. Completeness

The learner should describe the current risk, the redesign, and the trade-off versus an
alternative.

### 3. Boundary Handling

The redesign should discuss shutdown order, readiness, worker drain, or resource ownership
explicitly.

### 4. Code Quality

The proposed redesign should clarify the system instead of adding accidental complexity.

### 5. Verification Discipline

The redesign should be defended with concrete lifecycle reasoning, not generic “best practice”
language.

### 6. Explanation Quality

The learner should explain why the chosen redesign is better than a shallow alternative.

## Common Weak Answers

- suggesting “just shut everything down faster” without discussing order
- proposing a broad refactor with no bounded seam
- naming lifecycle risks without tying them to actual service behavior

## Next Step

Use this task as a bridge into the
[11 Flagship stage](../../11-flagship.md)
when you want the same redesign pressure applied to a larger system.
