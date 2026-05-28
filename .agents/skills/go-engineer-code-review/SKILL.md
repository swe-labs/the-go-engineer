---
name: go-engineer-code-review
description: Review commits, PRs, diffs, staged changes, review threads, or the whole The Go Engineer repository for v2.1 architecture compliance, curriculum correctness, Go correctness, tests, CI risk, security, reliability, and Opslane risks.
---

# The Go Engineer Code Review

Use this skill for review requests such as:

- review this commit
- review this PR
- review staged changes
- review this branch against main
- review open PR conversations
- review unresolved PR threads
- review the whole project
- review a lesson, section, curriculum entry, or Opslane module

## Review stance

Act like a senior Go engineer and curriculum maintainer.

Do not trust this skill as a complete contract when repository files disagree. Inspect the current source-of-truth documents and let them override this abbreviated checklist.

Use findings-first review:

1. report actionable findings first
2. prioritize P0/P1/P2
3. fix confirmed findings when operating in implementation mode
4. summarize only after findings

Do not produce only a checklist.

## Operating modes

### Report-only mode

Use when the user asks only for a review report.

Return findings, verification, residual risks, and summary.

### Fix mode

Use when the user asks to implement, prepare a PR, resolve feedback, or complete workflow.

In fix mode:

1. inspect the diff and PR context
2. identify findings
3. fix P0/P1/P2 findings before commit
4. run relevant validation
5. commit with `[TYPE]` format when requested by the workflow
6. update the PR
7. answer and resolve addressed review threads
8. leave a final readiness comment

Final squash-merge is not part of this skill.

## Required context

Inspect relevant source-of-truth files:

- `AGENTS.md` for internal workflow only when present or provided by the environment
- `ARCHITECTURE.md`
- `CURRICULUM-BLUEPRINT.md`
- `CODE-STANDARDS.md`
- `TESTING-STANDARDS.md`
- `CONTRIBUTING.md`
- `curriculum.v2.json`
- `.github/workflows/ci.yml`
- `docs/ENGINEERING_ERROR_FRAMEWORK.md` for backend/error changes
- `docs/flagship/OPSLANE_SAAS_BACKEND.md` for s11 or Opslane changes

## Architecture lock

Flag as P1 unless explicitly requested by the maintainer:

- new public root-level section outside `s00-s11`
- renamed canonical section folders
- section ownership changes contrary to `ARCHITECTURE.md`
- revived legacy folder structures
- lesson prefix or section mapping changes without explicit architecture approval

## Review-before-commit checklist

Before a commit, review the pending diff for:

- Architecture v2.1 drift
- changed lesson IDs, prefixes, folders, or paths
- `curriculum.v2.json` mismatches
- `curriculum.v2.json` non-canonical formatting, field order drift, duplicate keys, `null` arrays, out-of-order items, duplicate references, or self-references
- README `Next Step` mismatch
- source `NEXT UP:` mismatch
- source `Level` or `RUN:` header mismatch
- missing section README update
- missing or stale run/test commands
- missing or low-quality Machine Role comments
- stale or detached cross-reference blocks
- Go formatting or compile risk
- test gaps
- CI failure risk
- stale Makefile/script path usage
- security or data leakage risk
- Opslane tenant/auth/payment/worker/shutdown risk when relevant

Fix P0/P1/P2 findings before committing unless explicitly accepted.

## PR conversation review

When reviewing an open PR, inspect:

- top-level conversation comments
- review submissions
- inline review threads
- unresolved threads
- check status
- workflow logs for failed jobs

For each thread:

- classify as fixed, valid-and-unfixed, outdated, not-actionable, or maintainer-decision
- make code/doc changes for valid findings
- reply with the fix summary and validation command
- resolve the thread when the issue is addressed or no longer applies
- do not resolve maintainer-decision threads without explicit direction

Suggested reply formats:

```text
Fixed in the latest commit: <short fix>. Verified with `<command>`.
```

```text
The changed lines no longer exist after the latest update, so this thread is outdated. Marking resolved.
```

```text
Keeping this unchanged because <reason>. Remaining risk: <risk>. Validation: `<command>`.
```

## Severity rubric

- **P0 Blocker**: data loss, severe security issue, CI cannot pass, repository cannot build, curriculum validator broken globally.
- **P1 High**: architecture contract violation, learner path break, wrong lesson registration, broken run command, auth/tenant/security risk, race/deadlock, Opslane correctness issue.
- **P2 Medium**: real bug, missing required test, README/source mismatch, incomplete starter, Go standard violation with maintainability impact.
- **P3 Low**: useful improvement but not merge-blocking.
- **Nit**: style only. Avoid unless requested.

## The Go Engineer-specific checks

### Architecture and curriculum

Check:

- `ARCHITECTURE.md` remains the source of truth.
- All sections remain `s00-s11`.
- Folder path matches `curriculum.v2.json`.
- The registry remains canonical JSON with sections ordered `s00-s11` and items grouped by section, then path/ID order.
- Lesson ID prefix matches the owning section.
- `next_items` aligns with README `Next Step` and terminal `NEXT UP:`.
- New lessons are registered exactly once.
- Removed or renamed lessons do not leave dead curriculum entries.
- Section README is updated when lesson order or content changes.
- `go run ./scripts/validate_curriculum.go` should pass.

### README-first teaching contract

For every lesson README, verify required sections and order:

1. `## Mission`
2. `## Prerequisites`
3. `## Mental Model`
4. `## Visual Model`
5. `## Machine View`
6. `## Run Instructions`
7. `## Code Walkthrough`
8. `## Try It`
9. `## In Production`
10. `## Thinking Questions`
11. `## Next Step`

For exercises, accept `## Solution Walkthrough` in place of `## Code Walkthrough` and require `## Verification Surface`.

Verify that cross-references in `README.md` use standardized GitHub-style alerts (`[!NOTE]` or `[!TIP]`) and avoid standalone "Forward/Backward Reference" headlines.

Check that specific lesson references include the lesson ID and a clickable local `README.md` link when practical.

### Go source contract

Check:

- required header from `CODE-STANDARDS.md`
- `Level` and `RUN:` header values match `curriculum.v2.json`
- package and naming conventions
- Go constants are not `SCREAMING_SNAKE_CASE`
- exported symbols have doc comments when applicable
- Machine Role comments for major constructs are present and explain role, boundary, invariant, or failure mode rather than restating syntax
- source cross-references use lesson IDs and explain local relevance
- no ignored meaningful errors
- `%w` is used for wrapping where preserving cause matters
- no panic for control flow
- no global mutable state unless justified
- `ctx context.Context` is first parameter for I/O or blocking work
- contexts are canceled
- files and resources are closed
- concurrency uses `defer wg.Done()`, pointer WaitGroups, owned channel close, and no goroutine leaks
- terminal footer uses `NEXT UP:` without emoji and matches curriculum metadata

### Tests

Check:

- exercises include meaningful tests
- table-driven tests are used for more than two cases
- HTTP handlers use `httptest`
- performance-critical lessons include benchmarks where appropriate
- parsing or validation code has fuzz tests when appropriate
- concurrency changes are covered by deterministic tests and `go test -race ./...`
- coverage expectations are respected for the relevant section

### CI and tooling

The repository CI-equivalent validation set is:

```bash
go build ./...
go vet ./...
unformatted=$(gofmt -l .); test -z "$unformatted" || (echo "$unformatted" && exit 1)
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
go run ./scripts/validate_curriculum.go
```

Benchmark-related changes:

```bash
go test -bench=. -benchmem -count=1 ./08-quality-test/01-quality-and-performance/02-testing/04-benchmarks/
```

### Backend and security

For s06, s09, s10, and s11, check:

- input validation at HTTP boundaries
- parameterized SQL only
- no raw secrets, tokens, passwords, session IDs, or PII in logs
- error responses do not leak internals
- error categories follow UserError/SystemError/FatalError where applicable
- DB operations use context-aware methods
- transactions wrap multi-step state changes
- idempotency for retries and payment/order workflows
- tenant/user scoping is enforced in middleware, service, and repository layers as appropriate

### Opslane

For s11 or Opslane changes, additionally check:

- modular monolith boundaries are preserved
- auth and tenant identity are never optional state
- order and payment state transitions are explicit and safe
- retries use timeout/backoff and duplicate protection
- workers are bounded and support backpressure/drain behavior
- observability includes correlation IDs, structured logs, metrics/traces where relevant
- graceful shutdown drains in-flight work and closes shared resources
- configuration validates at startup and never commits real secrets

## Evidence standard

Every finding must include:

- exact file and line/range when available
- failing scenario
- why this matters for learners, CI, architecture, security, or production behavior
- concrete fix or test

If line numbers are unavailable, cite file and nearby symbol or section.

## Output format

```markdown
## Findings

### P1: <short actionable title>
- **Location:** `path/to/file.go:line-line`
- **Problem:** <what breaks and in which scenario>
- **Why it matters:** <impact>
- **Suggested fix:** <specific fix or test>

## Open questions / assumptions

- <only focused questions that affect review confidence>

## Verification performed

- <commands run, files inspected, tests run/not run>

## Residual risks

- <what remains uncertain>

## Summary

<1-3 sentences after findings>
```

If no issues are found:

```markdown
## Findings

No blocking or high-confidence issues found.

## Verification performed

- <commands/files/tests>

## Residual risks

- <testing gaps, unreviewed areas, assumptions>

## Summary

<brief summary>
```

## Final PR readiness comment

After fixing feedback and before maintainer handoff, leave:

```markdown
## Final readiness check

### Changes
- 

### Review findings handled
- 

### Validation
- 

### Remaining risks
- 
```

## Anti-patterns

Do not:

- produce only a checklist
- report style-only nits by default
- claim commands passed unless actually run
- approve architecture changes that contradict v2.1-final
- hide serious issues in the summary
- ignore unresolved PR threads
- resolve a review thread without fixing or explaining it
- comment on generated/vendor files unless the generator/config changed or output reveals a real issue
- merge the PR
