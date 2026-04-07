# Contributing to The Go Engineer

Thank you for your interest in contributing! This guide will help you get started.

## Quick Links

- **[Testing Standards](./TESTING-STANDARDS.md)** - How to write tests
- **[Code Quality Standards](./CODE-STANDARDS.md)** - Code style and best practices
- **[Curriculum Map](./docs/curriculum/README.md)** - Complete curriculum structure

## Getting Started

```bash
# Clone the repository
git clone https://github.com/rasel9t6/the-go-engineer.git
cd the-go-engineer

# Verify your environment
make build
make test
make lint
```

## Branch Strategy

The repository uses long-lived branches for supported major versions:

- `main`: active v2 development and the default target for new work
- `release/v1`: stable v1 maintenance for bug fixes and support updates
- `release/v2`: created from `main` when v2 enters beta and feature freeze

Create short-lived topic branches from the branch your change should ship to:

- `feat/...` from `main` for new sections, lessons, and v2 features
- `fix/...` from `main` for v2 bug fixes
- `fix/v1-...` from `release/v1` for v1-only bug fixes
- `docs/...` or `chore/...` from the appropriate target branch

### Pull Request Targeting

- Target `main` by default.
- Only target `release/v1` when the issue is explicitly v1-only or the fix must reach current stable users first.
- If a fix belongs in both lines, merge it into the correct source branch first and then `git cherry-pick -x` it to the other supported branch.
- Do not try to keep `main` and `release/v1` identical once v2 development begins. Divergence between supported major versions is expected.

### Merge Policy

- Maintainers use **Squash and Merge** for pull requests into `main`, `release/v1`, and future release branches.
- Do not work directly on long-lived branches.
- Auto-delete topic branches after merge when possible.

## Code Style

### Every Go File Must Follow This Template

```go
package main

import "fmt"

// ============================================================================
// Section N: Section Name — Lesson Title
// Level: Beginner | Intermediate | Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Point 1
//   - Point 2
//
// RUN: go run ./NN-section-name/N-lesson-name
// ============================================================================

func main() {
    // Your code with inline comments explaining every concept
    fmt.Println("Hello")

    // KEY TAKEAWAY:
    // - Summary point 1
    // - Summary point 2
    fmt.Println("\n---------------------------------------------------")
    fmt.Println("🚀 NEXT UP: SS.N next-lesson-name")
    fmt.Println("   Current: SS.N (current-lesson-name)")
    fmt.Println("---------------------------------------------------")
}
```

### Comment Guidelines

1. **Every non-obvious line gets a comment** — This is a teaching repo, not a production codebase
2. **Explain WHY, not WHAT** — `// Increment i` is useless. `// Move to the next element because...` teaches
3. **Use block comments for concepts** — Multi-line `//` comments above code blocks to explain the concept before showing the code
4. **Include KEY TAKEAWAY** — Every `main()` function should end with a KEY TAKEAWAY summary
5. **Cross-reference sections** — Use `(See Section 05: Interfaces)` when a concept is covered elsewhere

### Formatting Rules

- Run `gofmt` on every file (or use `make fmt`)
- Run `go vet ./...` before committing (or use `make vet`)
- Run `make lint` to verify all checks pass

## Section Numbering

- Sections are numbered `00-22` with descriptive names
- Lessons within a section are numbered `1-`, `2-`, etc.
- Exercises are the LAST numbered item in a section

## Curriculum Metadata During V2 Migration

During the v1 to v2 transition, the repo has two metadata surfaces with different jobs:

- `curriculum.json` remains the legacy lesson graph used for current learner-facing compatibility
- `curriculum.v2.json` is the additive metadata surface for migrated v2 sections and typed v2 items

Use these rules while both files exist:

1. Keep shared lesson truth aligned across both files when a lesson appears in both places.
   Shared truth means the stable lesson id, the live repo path, and the lesson-level prerequisite graph.
2. Put v2-only item types in `curriculum.v2.json`, not in `curriculum.json`.
   That includes exercises, checkpoints, mini-projects, and later capstones.
3. Treat `curriculum.v2.json` as optional during migration.
   The validator will keep supporting `curriculum.json` while v2 coverage grows section by section.
4. Use current live repo paths in `curriculum.v2.json` until the matching live migration issue changes them.
   Planning-only prototype ids and paths must not be added to `main` before the real section migration lands.

## Adding a New Lesson

1. Create a directory: `NN-section-name/N-lesson-name/`
2. Create `main.go` following the template above
3. Update the section's `README.md` with the new lesson in the `Learning Path` table
4. Update `curriculum.json` with the new lesson ID, concept, and prerequisites
5. Run `go run scripts/validate_curriculum.go` to verify the mapping
6. Update the root `README.md` if adding a new section
7. Open the pull request against `main`
8. Verify: `make build && make test && make lint`

## Adding an Exercise

Exercises include both complete solutions and starter stubs:

```text
NN-section-name/
└── N-exercise-name/
    ├── main.go              ← Complete solution with comments
    └── _starter/
        └── main.go          ← TODO stubs for self-challenge
```

The `_starter/main.go` should:

1. Have the same file header and REQUIREMENTS checklist
2. Contain function signatures with `// TODO: implement this` bodies
3. Compile successfully (return zero-values from stubs)
4. Print a message directing students to the requirements

## Commit Messages

Use clear, descriptive commit messages:

```text
Add: Section 17 Context deep-dive (4 lessons)
Fix: go vet warning in 07-strings formatting
Update: Backfill comments in Section 02 control flow
```

## Backports and Stable Fixes

When a bug affects both stable v1 and active v2:

1. Fix it on the branch that needs the release first.
2. Merge that pull request with squash.
3. Cherry-pick the resulting commit with `git cherry-pick -x` onto the other supported branch.

This keeps history intentional and avoids branch-to-branch sync merges that make maintenance harder over time.

## Quality Checklist

Before submitting, verify:

- [ ] `make build` passes (all packages compile)
- [ ] `make vet` passes (no suspicious code)
- [ ] `make fmt-check` passes (code is formatted)
- [ ] `make test` passes (all tests pass)
- [ ] `go run scripts/validate_curriculum.go` passes
- [ ] Every new Go file has the standard header template
- [ ] Every new Go file has the "NEXT UP" footer
- [ ] Every concept has inline teaching comments
- [ ] The section README is updated with the ID-based table
- [ ] `curriculum.json` and [docs/curriculum/dependency-graphs.html](./docs/curriculum/dependency-graphs.html) are synced
- [ ] CI pipeline passes on push (GitHub Actions)
