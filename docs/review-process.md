# Review Process

This document defines the review standard for Go Engineer changes.

The review process is findings-first. Reviewers should identify concrete issues, severity, and required fixes before approval.

## Review goals

A review protects:

- learner clarity
- technical correctness
- curriculum graph integrity
- validation strictness
- migration traceability
- repository maintainability
- release safety

## Severity levels

| Severity | Meaning | Approval allowed? |
|---|---|---|
| P0 | Broken graph, invalid JSON, missing required file, failing release gate, secret leak | No |
| P1 | Incorrect teaching, broken code/test, missing proof surface, invalid assessment target | No |
| P2 | Weak explanation, generic crossref, missing production context, incomplete rubric | Usually no |
| P3 | Style, wording, naming, minor polish | Yes with follow-up if needed |

## Metadata review checklist

- [ ] JSON parses.
- [ ] IDs are stable and unique.
- [ ] Module paths use `curriculum/modules/` or `curriculum/electives/`.
- [ ] Lesson paths use typed folders.
- [ ] Prerequisites resolve.
- [ ] Next item links resolve.
- [ ] Concepts have owners and reinforcement.
- [ ] Crossrefs have specific reasons.
- [ ] Projects have rubrics and assessment bindings.
- [ ] Assessments have targets and evidence.
- [ ] Migration coverage remains complete.
- [ ] No placeholder/scaffold statuses remain.

## Content review checklist

- [ ] README follows required section contract.
- [ ] Mental model is clear.
- [ ] Visual model is useful.
- [ ] Machine/runtime view is accurate.
- [ ] Code walkthrough explains why, not just what.
- [ ] Common mistakes are practical.
- [ ] Debugging signals are concrete.
- [ ] Production context is realistic.
- [ ] Security/performance notes are present where relevant.
- [ ] Proof of understanding is observable.
- [ ] Next step is correct.

## Code review checklist

- [ ] Code compiles.
- [ ] Code is gofmt-formatted.
- [ ] Errors are handled explicitly.
- [ ] Tests prove behavior.
- [ ] Tests cover edge cases.
- [ ] HTTP/database/concurrency examples avoid unsafe shortcuts.
- [ ] Starter and solution match.
- [ ] Comments explain non-obvious behavior.
- [ ] No secrets or credentials.
- [ ] No unnecessary dependencies.

## Project review checklist

- [ ] Project mission is clear.
- [ ] Requirements are testable.
- [ ] Constraints are realistic.
- [ ] Architecture is explained.
- [ ] Deliverables are explicit.
- [ ] Rubric weights sum correctly.
- [ ] Verification path is runnable.
- [ ] Portfolio guidance is useful.

## Assessment review checklist

- [ ] Assessment target IDs are valid.
- [ ] Questions/tasks match taught content.
- [ ] Evidence requirements are clear.
- [ ] Rubric is objective.
- [ ] Answer key exists when appropriate.
- [ ] Retake policy exists.
- [ ] Passing standard is explicit.

## Review comment format

Use:

```text
[P1] Missing assessment evidence requirement

The assessment asks learners to explain transactions but does not define what evidence counts as passing. Add a rubric row for isolation-level reasoning and a required example failure scenario.
```

A good finding includes:

- severity
- exact problem
- why it matters
- requested fix

## Approval rule

Approve only when:

- P0/P1 findings are fixed
- P2 findings are fixed or explicitly accepted with tracking
- validation passes
- remaining risks are documented

Do not approve because a change is large or urgent. Large changes need stricter review, not weaker review.

## Post-review

After fixes:

1. Re-run relevant validation.
2. Confirm changed files are still aligned.
3. Check that generated artifacts were not hand-edited.
4. Re-review previously failing areas.
5. Approve only after the evidence is clean.
