# Maintainer Checklist

Use this checklist to keep repository workflow, branch maintenance, documentation, and releases consistent.

All module references follow the canonical v3 architecture from [ARCHITECTURE.md](./ARCHITECTURE.md).

## Daily Triage

- Enforce issue-first workflow.
- Confirm new contributor PRs target the correct branch (`release/v3` or `main`).
- Add labels early.
- Move issues into the correct milestone and the GitHub Project: `The Go Engineer v3`.
- Redirect release-line work to the correct release branch.
- Reject unapproved architectural drift.

## Project and milestone tracking

For every issue and PR:

- assign the owner when appropriate
- apply relevant labels
- attach the active milestone when one exists
- add the issue or PR to `The Go Engineer v3`
- update tracking metadata when scope changes

## Required Issue Title Format

```text
[TYPE] short description
```

Allowed types:

```text
[FEAT] [FIX] [DOCS] [TEST] [CHORE] [REFACTOR] [SECURITY] [RELEASE]
```

## Closed-Loop PR Maintenance

Before a PR is ready for final maintainer review:

- findings-first self-review must be complete
- P0/P1/P2 findings must be fixed or explicitly documented
- CI failures must be fixed or explained with evidence
- top-level PR comments must be answered
- inline review threads must be fixed, explained, or marked outdated when the diff no longer applies
- addressed threads should be resolved when permitted
- a final readiness comment should summarize changes, handled findings, validation, and remaining risks

The final squash-merge remains maintainer-only.

## PR Review Rules

- Reject PRs that bypass the README-first teaching contract.
- Reject PRs that introduce "magic" early in the curriculum.
- Reject public architecture drift unless explicitly approved.
- Require linked issues.
- Require validation evidence.
- Use **Squash and Merge**.
- Never develop directly on long-lived branches.
- Do not self-merge unless explicitly approved for that PR.

## Merge Commit Format

Final squash commits should use:

```text
[TYPE] short imperative description
```

Examples:

```text
[FEAT] add bounded worker pool lesson
[FIX] correct tenant-scoping example in Opslane
[DOCS] align release docs with v3.0.0 architecture
```

## Branch Roles

- `main`: active development and next-generation feature work
- `release/v3`: active v3 release preparation and maintenance line

## Backports

- For a v3 bug, fix on `release/v3` first and cherry-pick to `main`.
- If a fix belongs in multiple supported lines, use `git cherry-pick -x`.

```bash
git switch <target-branch>
git pull origin <target-branch>
git cherry-pick -x <merged-commit-sha>
git push origin <target-branch>
```

## Release Flow

For v3.x maintenance releases:

1. Branch from `release/v3`.
2. Update `CHANGELOG.md` and any required metadata.
3. Run local verification.
4. Open a release PR into `release/v3`.
5. Squash and merge after green CI and approval.

## Maintainer vs. Learner Isolation Invariant

Maintainers and authors must strictly preserve the boundary between internal tooling and learner-facing content:

1. **Maintainer-Only Scope**:
   - `internal/tools/curriculum/`: Automated graph engine, schema validation, contract validation, and verification CLI.
   - `scripts/`: Operational runners (`validate_curriculum.go`), maintenance scripts, catalog generators.
   - `curriculum/*.json` + `curriculum/modules/*/manifest.json`: Authoritative schema and graph manifests (thin `path.core.json`/`path.electives.json` indexes, per-module item/project/assessment manifests, `schema.v3.json`, etc.).
   - `MAINTAINER-CHECKLIST.md`, `CONTRIBUTING.md`, `ARCHITECTURE.md`: Contributor and maintainer standards.

2. **Learner-Facing Scope**:
   - `curriculum/modules/`: 18 core backend modules (instructional lessons, project capstones, assessment checkpoints).
   - `curriculum/companion-frontend/`: Server-driven full-stack hypermedia track.
   - `curriculum/electives/`: Advanced specialization tracks.
   - Learner tools: Terminal, text editor/IDE, standard Go toolchain (`go run`, `go test`, `go build`).

3. **Strict Isolation Rule**:
   - Learner-facing READMEs, lessons, projects, and assessments must **never** reference internal maintainer scripts (`scripts/validate_curriculum.go`), internal packages (`internal/tools/`), JSON manifest schemas (`schema.v3.json`), or automated quality gates.
   - Learners are never instructed to run maintainer validation scripts or tested on repo maintainer infrastructure.

## Local Verification

Run this sequence before committing or opening pull requests:

```bash
go build ./...
go vet ./...
unformatted=$(gofmt -l . | grep -v '_starter' || true); test -z "$unformatted" || (echo "$unformatted" && exit 1)
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go test -race ./curriculum/modules/11-lifecycle-context-concurrency/... ./curriculum/modules/18-flagship-opslane/...
go test -v ./internal/tools/curriculum/...
go run ./scripts/validate_curriculum.go
```

### The Authoritative Curriculum Quality Gates

All validation gates executed by `go run ./scripts/validate_curriculum.go` must pass. The authoritative command list lives in the `validators()` registry in `internal/tools/curriculum/main.go`; this table is a reference:

| Gate | Description |
|:---|:---|
| `validate-schema` | Validates all 27 JSON manifests (8 top-level + 19 per-module) against `schema.v3.json`. |
| `validate-graph` | Detects cycles and verifies DAG reachability. |
| `validate-lessons` | Enforces zero machine-role comments and valid `NEXT UP` footers. |
| `validate-prerequisite-closure` | Verifies prerequisite closures across modules. |
| `validate-project-binding` | Ensures all 29 projects bind to valid anchors, deliverables, and rubrics. |
| `validate-no-placeholder-zm` | Scans for boilerplate or placeholder phrases in `zero_magic` blocks. |
| `validate-assessment-completeness` | Verifies all assessments have non-generic criteria and passing scores. |
| `validate-unique-tags` | Enforces tag taxonomy consistency. |
| `validate-crossrefs` | Validates bidirectional integrity across cross-reference edges. |
| `validate-concept-ownership` | Ensures all 69 canonical concepts map to unambiguous owners. |
| `validate-cognitive-load` | Verifies cognitive density ratings and pacing. |
| `validate-operational-failures` | Verifies 100% mathematical symmetry across all 13 failure categories. |
| `validate-failure-engineering` | Ensures production modules include concrete failure examples. |
| `validate-zero-magic` | Validates adherence of all 356 golden lessons to zero-magic authoring spec. |
| `validate-physical-artifacts` | Verifies existence of all 447 instructional READMEs on disk. |
| `validate-assessment-artifacts` | Verifies 4-file assessment suites across all checkpoints and milestones. |
| `validate-starter-parity` | Compiles and validates Go AST correctness across starters and solutions. |
| `validate-markdown-links` | Crawls Markdown files and verifies all relative links exist on disk. |
| `validate-inline-paths` | Crawls Markdown files and verifies all relative links exist on disk. |
| `validate-anti-slop` | Crawls Markdown files and verifies all relative links exist on disk. |
| `validate-sese-navigation` | Enforces SESE entry (`## Start here`) and exit gate integrity. |
| `validate-project-scripts` | Verifies cross-platform script parity (`verify.sh` and `verify.ps1`). |
| `validate-no-emojis` | Enforces strict zero-emoji invariant across all Markdown prose. |
| `validate-anchor-quotes` | Enforces pioneer anchor quotes at the top of locked lessons. |
| `validate-diagram-palette` | Enforces Dark Engineering Palette classes and colors on Mermaid flowcharts. |
| `validate-readme-contracts` | Resolves each item's document contract (family + subtype) from `readme.contracts.json`; verifies required headings, code policy, and file layout. Locked-module content-shape defects block; other items follow `severity_policy.by_status` (draft/scaffolded skip, review/reviewed warn, ready-and-later block); `--strict` promotes the remaining advisories. |

## Documentation Alignment Check

Before release or workflow changes, verify these documents agree with `ARCHITECTURE.md` and each other:

- `README.md`
- `ARCHITECTURE.md`
- `ROADMAP.md`
- `LEARNING-PATH.md`
- `CODE-STANDARDS.md`
- `TESTING-STANDARDS.md`
- `COMMON-MISTAKES.md`
- `docs/PROGRESSION.md`
- `docs/ENGINEERING_ERROR_FRAMEWORK.md`
- `docs/flagship/OPSLANE_SAAS_BACKEND.md`
- `CONTRIBUTING.md`
- `RELEASE.md`
- `MAINTAINER-CHECKLIST.md`
- `.github/PULL_REQUEST_TEMPLATE.md`
- `.github/ISSUE_TEMPLATE/*`

## Architecture Checklist

- [ ] Core module count remains 18 (`00` through `16`, plus `18-flagship-opslane`).
- [ ] Module directories remain canonically organized under `curriculum/modules/`.
- [ ] Item IDs belong to the declared module progression.
- [ ] Metadata manifests in `curriculum/` match filesystem paths.
- [ ] All automated curriculum quality gates pass.

## Lesson Checklist

- [ ] README follows the required pedagogical section order.
- [ ] `main.go` follows the standard header and key takeaway footer.
- [ ] `NEXT UP:` matches the item's per-module manifest entry.
- [ ] Module README is updated.
- [ ] Starter code compiles and establishes red verification state.
- [ ] Tests prove behavior.
- [ ] Run command works exactly as printed.
