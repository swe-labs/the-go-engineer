---
name: github-workflow
description: Strict GitHub operational workflow for issue creation, branching, review-before-commit, committing, PR management, conversation resolution, CI verification, and maintainer handoff in The Go Engineer v2.1-final repository.
---

# GitHub Workflow

Use this skill whenever making code, documentation, curriculum, CI, release, or tooling changes.

This workflow is closed-loop: implement, review, fix, commit, open/update PR, process feedback, resolve addressed threads, and hand off only when ready. The final squash-merge remains maintainer-only.

## Repository context

- Repository: `swe-labs/the-go-engineer`
- Default branch: `main`
- Architecture: v2.1-final, locked
- Public spine: 5 phases, 12 sections, `s00` through `s11`
- Source of truth: `ARCHITECTURE.md`
- Curriculum registry: `curriculum.v2.json`
- Validator: `go run ./scripts/validate_curriculum.go`

Do not propose new public root sections or migration restructuring unless explicitly requested by the maintainer.

## Required references

Read the relevant docs before implementation:

- `AGENTS.md`
- `ARCHITECTURE.md`
- `CURRICULUM-BLUEPRINT.md`
- `CODE-STANDARDS.md`
- `TESTING-STANDARDS.md`
- `CONTRIBUTING.md`
- `docs/ENGINEERING_ERROR_FRAMEWORK.md` for backend or error-handling work
- `docs/flagship/OPSLANE_SAAS_BACKEND.md` for s11 or flagship work

## Workflow

### 1. Issue first

Create or confirm an issue.

Title format:

```text
[TYPE] short description
```

Allowed types (match the branch prefix mapping below):

```text
[FEAT] [FIX] [DOCS] [TEST] [CHORE] [REFACTOR] [SECURITY] [RELEASE]
```

Type-to-branch mapping — brackets for human-facing surfaces, slash-prefixed for branch names:

| Type | Branch prefix |
| --- | --- |
| `[FEAT]` | `feat/` |
| `[FIX]` | `fix/` |
| `[DOCS]` | `docs/` |
| `[TEST]` | `test/` |
| `[CHORE]` | `chore/` |
| `[REFACTOR]` | `refactor/` |
| `[SECURITY]` | `security/` |
| `[RELEASE]` | `release/` |

Issue body must include:

- affected section: `s00` through `s11`, if applicable
- affected lesson IDs, if applicable
- expected files or folders
- validation plan
- risk notes


### 1.1. Project and milestone tracking

For every issue and PR, keep tracking metadata current:

- assign `rasel9t6` when appropriate
- apply relevant labels
- attach the active milestone when one exists
- add the issue or PR to the GitHub Project: `The Go Engineer v2`
- update labels, milestone, and project placement when scope changes

Command examples:

```bash
gh issue edit <issue-number> --add-assignee rasel9t6
gh issue edit <issue-number> --milestone "<milestone>"
gh issue edit <issue-number> --add-label "<label>"
gh project item-add "<project-number-or-url>" --owner rasel9t6 --url "<issue-or-pr-url>"
```

### 2. Branch from the correct line

For post-v2.1 work:

```bash
git switch main
git pull origin main
git switch -c <branch>
```

For stable v2.1.x fixes:

```bash
git switch release/v2
git pull origin release/v2
git switch -c <branch>
```

For stable v1 fixes:

```bash
git switch release/v1
git pull origin release/v1
git switch -c <branch>
```

Branch naming:

Branches use slash-prefixed types (no brackets — brackets break git URLs and shell tab-completion).

| Pattern | Use |
| --- | --- |
| `feat/sNN-lesson-slug` | lesson or section content |
| `fix/sNN-description` | bug fix |
| `docs/document-or-section` | documentation-only change |
| `test/sNN-description` | test-only improvement |
| `chore/tooling-description` | tooling, CI, dependencies |
| `refactor/description` | code improvement without behavior change |
| `security/description` | security fix |
| `release/vX.Y.Z-prep` | release metadata work |

### 2.1. PR title format

PR titles use the same bracketed format as issues and commits:

```text
[TYPE] short description
```

### 3. Open a draft PR early

Open a draft PR after the first push.

PR body must include:

```markdown
Closes #<issue>

## Scope
- Section:
- Lesson IDs:
- Files:

## Type
- [ ] [FEAT]
- [ ] [FIX]
- [ ] [DOCS]
- [ ] [TEST]
- [ ] [CHORE]
- [ ] [REFACTOR]
- [ ] [SECURITY]
- [ ] [RELEASE]

## Validation
- [ ] go build ./...
- [ ] go vet ./...
- [ ] gofmt check
- [ ] go mod tidy no-diff check
- [ ] go test ./...
- [ ] go test -race ./...
- [ ] go test -coverprofile=coverage.out ./...
- [ ] go run ./scripts/validate_curriculum.go

## Review Loop
- [ ] findings-first self-review complete
- [ ] P0/P1/P2 findings fixed or documented
- [ ] review threads answered
- [ ] addressed threads resolved

## Tracking
- [ ] issue/PR assigned when appropriate
- [ ] labels applied
- [ ] milestone attached when available
- [ ] project item added to `The Go Engineer v2`

## Risk
-
```

Keep the PR as draft until the review loop is complete.

### 4. Implement the change

Follow the relevant skill:

- `lesson-authoring` for lessons, exercises, READMEs, section maps, and `curriculum.v2.json`
- `go-engineer-code-review` for review passes
- this skill for issue/branch/commit/PR operations

### 5. Review before committing

Before every commit, run a findings-first review on the pending diff.

Use `go-engineer-code-review` against:

- staged changes
- unstaged changes
- changed files against the base branch
- relevant source-of-truth docs

Fix P0/P1/P2 findings before committing unless the maintainer explicitly accepts the risk.

Default pre-commit review checks:

- architecture v2.1 remains locked
- no new public root section unless approved
- `curriculum.v2.json` matches files and lesson metadata
- README next step and source `NEXT UP:` match curriculum metadata
- source `Level` and `RUN:` headers match curriculum metadata
- lesson README follows the required contract
- Machine Role comments explain role, boundary, invariant, or failure mode when lessons or examples change
- README cross-references use `[!NOTE]` or `[!TIP]` alerts and link specific lesson references
- starter code compiles
- Go files are formatted and idiomatic
- errors, contexts, concurrency, and resources are handled correctly
- tests cover behavior changes
- CI-equivalent commands are likely to pass

### 6. Commit logically

Use bracketed commit type format:

```text
[TYPE] short imperative description
```

Examples:

```text
[FEAT] add SY.1 sync.Mutex lesson
[FIX] correct tenant scoping in DB.3
[DOCS] align release workflow with v2.1 maintenance
[TEST] add benchmark coverage for PR.5
[CHORE] update validator documentation
```

Keep commits focused. Do not mix unrelated formatting with feature logic.

### 7. Push and monitor checks

After pushing:

1. read CI status
2. inspect failing workflow jobs when present
3. fix actionable failures
4. push follow-up commits
5. repeat until green or document unrelated failure evidence

Required full readiness commands:

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

### 8. Process PR conversation and review threads

For every open PR iteration, inspect:

- top-level PR comments
- review submissions
- inline review threads
- unresolved threads
- requested changes
- failing checks

For each valid finding:

1. fix it on the PR branch
2. run targeted validation
3. reply to the thread with the fix summary
4. resolve the thread when permitted

Use concise resolution replies:

```text
Fixed in the latest commit: <what changed>. Verified with `<command>`.
```

If the comment is outdated because the diff changed:

```text
The underlying diff no longer exists after the latest update. Marking resolved.
```

If the comment is not actionable:

```text
Leaving this unchanged because <reason>. The remaining risk is <risk>, and validation is <evidence>.
```

Do not ignore unresolved conversations. Do not resolve a thread without either fixing it or explaining why it no longer applies.

### 9. Final self-review before ready-for-review

Before marking the PR ready:

1. rerun `go-engineer-code-review` on the full PR diff
2. verify no P0/P1/P2 findings remain
3. verify issue linkage and PR checklist
4. verify CI status
5. leave a final PR comment

Final PR comment format:

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

### 10. Maintainer handoff

When ready:

- mark the PR ready for review
- request maintainer review if appropriate
- do not squash-merge
- do not enable auto-merge unless the maintainer explicitly requests it
- the maintainer performs final approval and squash-merge

## Hard stops

Stop and report rather than continuing when:

- requested work contradicts locked Architecture v2.1
- a P0/P1 finding cannot be fixed safely
- CI fails for a reason that requires maintainer decision
- a review thread requests a product or architecture decision
- credentials, secrets, or private data would be exposed
- final merge is requested without maintainer approval
