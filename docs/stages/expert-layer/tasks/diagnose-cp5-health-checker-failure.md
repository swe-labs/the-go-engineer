# CP.5 Health Checker Diagnosis

## Mission

Diagnose a failure or weakness in the `CP.5` health-checker milestone by working from evidence
instead of guesswork.

## Type

- diagnosis task

## Level

- stretch

## Prerequisites

- [07 Concurrency](../../07-concurrency.md)
- `CP.5`
- [rubric and checkpoint template](../../../templates/rubric-checkpoint-template.md)

## Task

1. Inspect the `CP.5` health-checker surface and identify one likely concurrency or coordination
   failure mode.
2. Explain what evidence you would inspect first to confirm it.
3. Propose one bounded fix and explain why it is safer than the most obvious shallow alternative.

## Evidence

- describe the suspected failure in concrete terms
- explain what signal, behavior, or code surface would support the diagnosis
- justify the proposed fix in terms of bounded concurrency, cancellation, or ownership

## Rubric

### 1. Correctness

The diagnosis should match a plausible concurrency failure mode.

### 2. Completeness

The learner should cover the failure, the evidence trail, and the proposed fix.

### 3. Boundary Handling

The learner should discuss cancellation, sibling work, ownership, or bounded concurrency directly.

### 4. Code Quality

The proposed fix should reduce ambiguity or coordination risk without overengineering the surface.

### 5. Verification Discipline

The diagnosis should rely on observable evidence instead of intuition alone.

### 6. Explanation Quality

The learner should explain why the chosen fix is better than a shallow alternative.

## Common Weak Answers

- guessing at a race or leak without naming the evidence trail
- proposing "just add more locks" without explaining the coordination problem
- describing a failure mode without suggesting a bounded fix

## Next Step

Use this task as a bridge into
[11 Flagship](../../11-flagship.md)
when you want the same judgment pressure applied to a larger system.
