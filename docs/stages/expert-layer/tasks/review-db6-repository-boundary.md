# DB.6 Repository Boundary Review

## Mission

Review the `DB.6` repository-pattern milestone and explain whether its boundaries are clear enough
for a small backend service.

## Type

- review task

## Level

- core

## Prerequisites

- [06 Backend & DB](../../06-backend-db.md)
- `DB.6`
- [rubric and checkpoint template](../../../templates/rubric-checkpoint-template.md)

## Task

1. Review the `DB.6` repository pattern surface and identify the three most important engineering
   risks or weaknesses.
2. Explain which one matters most and why.
3. Propose one boundary improvement that would make the surface easier to test, reason about, or
   evolve.

## Evidence

- point to the specific boundary or code shape you are criticizing
- explain why it is a real engineering problem instead of a style preference
- justify the priority order of your top three findings

## Rubric

### 1. Correctness

The review should describe real issues, not imagined ones.

### 2. Completeness

The learner should cover three findings and rank them clearly.

### 3. Boundary Handling

The review should discuss repository seams, SQL ownership, or service-boundary clarity directly.

### 4. Code Quality

The proposed improvement should make the surface clearer, safer, or easier to change.

### 5. Verification Discipline

The review should point to specific evidence from the milestone surface.

### 6. Explanation Quality

The learner should justify why the top issue matters most instead of only listing problems.

## Common Weak Answers

- naming three style nits instead of meaningful engineering risks
- proposing a large refactor without explaining the problem it solves
- criticizing boundaries without pointing to concrete evidence

## Next Step

Move to the diagnosis task at
[CP.5 health checker diagnosis](./diagnose-cp5-health-checker-failure.md)
or continue to the
[11 Flagship stage](../../11-flagship.md)
if you want longer-form pressure next.
