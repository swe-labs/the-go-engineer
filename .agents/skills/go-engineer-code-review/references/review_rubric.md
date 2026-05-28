# The Go Engineer Review Rubric

## Highest-risk areas

1. Architecture v2.1 drift.
2. Curriculum registry, README, source, and footer mismatch.
3. Broken run or test commands.
4. Broken CI gates.
5. Security and tenant isolation in s06, s09, and s11.
6. Concurrency leaks or races in s07 and s11.
7. Starter code that does not compile.
8. Missing tests for exercises or behavior changes.
9. Unresolved PR feedback left unanswered.
10. Review threads resolved without a fix or explanation.

## Finding quality bar

A finding should be reported when it is:

- introduced or exposed by the change
- actionable
- supported by code or documentation evidence
- important enough to fix before merge or track intentionally

## Closed-loop expectation

When operating in fix mode:

1. find the issue
2. fix the issue
3. run targeted validation
4. commit with `[TYPE]` format when appropriate
5. reply to the related PR thread
6. resolve the thread when permitted
7. document remaining risk in the final readiness comment

Final squash-merge is outside the review loop and remains maintainer-only.
