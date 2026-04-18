# Contributing to The Go Engineer

Thank you for your interest in contributing! This guide will help you get started.

## Quick Links

- **[Architecture Blueprint](./ARCHITECTURE.md)** — The 12-section v2.1 structure (source of truth)
- **[Curriculum Blueprint](./CURRICULUM-BLUEPRINT.md)** — Teaching and lesson contract standards
- **[Code Quality Standards](./CODE-STANDARDS.md)** — Code style and engineering best practices
- **[Testing Standards](./TESTING-STANDARDS.md)** — How to write tests

## Getting Started

```bash
# Clone the repository
git clone https://github.com/rasel9t6/the-go-engineer.git
cd the-go-engineer

# Verify your environment
go version
go test ./...
```

---

## The Strict GitHub Workflow

To maintain a high standard of quality, all contributors (including maintainers) MUST follow this workflow:

### 1. Create an Issue First

Never start coding without an approved issue.

- Ensure the issue maps to a specific gap or lesson in the 12-section architecture.
- Add labels, assign yourself, and add the issue to the **"The Go Engineer v2"** project.
- A maintainer must approve the issue before work begins.

### 2. Branch from `main`

The repository uses long-lived branches for supported major versions.

- `main`: active v2 development and the default target for new work.
- `release/v1`: stable v1 maintenance for bug fixes and support updates.

Create short-lived topic branches from `main` (for v2 work):

- `feat/...` for new sections or lessons
- `fix/...` for bug fixes
- `docs/...` or `chore/...` for documentation and tooling

### 3. Open a Draft PR

As soon as you push your branch, open a Draft Pull Request linking to the issue (`Closes #123`).

### 4. Logical Commits

Make logical, atomic commits. Do not lump unrelated formatting changes in with feature logic.

### 5. Final CI & Review

Ensure the CI pipeline passes. Once ready, mark the PR as ready for review.
Maintainers will use **Squash and Merge** to merge your PR into `main`.

---

## Code Style

### Every Go Lesson File Must Follow This Template

```go
// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section NN: Section Name — Lesson Title
// Level: Foundation | Core | Stretch
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Concept one
//   - Concept two
//
// WHY THIS MATTERS:
//   One sentence on why this exists in a real production Go codebase.
//
// RUN: go run ./NN-section-slug/N-lesson-slug
// ============================================================================

package main

import "fmt"

func main() {
    // Your code with inline comments explaining every concept
    fmt.Println("Hello")

    // KEY TAKEAWAY:
    // - Summary point 1
    // - Summary point 2
    fmt.Println()
    fmt.Println("---------------------------------------------------")
    fmt.Println("NEXT UP: XY.N next-lesson-slug")
    fmt.Println("Run    : go run ./NN-section-slug/N-next-lesson-slug")
    fmt.Println("Current: XY.N (current-lesson-slug)")
    fmt.Println("---------------------------------------------------")
}
```

**NEXT UP footer format:** Must use `NEXT UP:` (without emoji) followed by the exact item ID and slug from `curriculum.v2.json`. The validator checks this format with the pattern `NEXT UP:\s*([A-Z]{2,6}\.\d+)`.

### Comment Guidelines

1. **Every non-obvious line gets a comment** — This is a teaching repo, not a production codebase
2. **Explain WHY, not WHAT** — `// Increment i` is useless. `// Move to the next element because...` teaches
3. **Use block comments for concepts** — Multi-line `//` comments above code blocks to explain the concept before showing the code
4. **Include KEY TAKEAWAY** — Every `main()` function should end with a KEY TAKEAWAY summary
5. **Cross-reference lessons** — Use the lesson ID format: _(Context cancellation is covered in CT.3.)_

### Formatting Rules

- Run `go fmt ./...` on every file
- Run `go vet ./...` before committing

---

## Section Numbering

- Sections are numbered `s00–s11` according to `ARCHITECTURE.md`.
- Lessons within a section use their subsystem ID prefix (e.g., HC, GT, LB, CF, DS, FE, etc.).
- Exercises are the LAST item in their subsystem group.

## Adding a New Lesson

1. Ensure an Issue is open and approved.
2. Create the directory: `NN-section-slug/N-lesson-slug/`
3. Create the `README.md` following the Lesson README Contract (see `ARCHITECTURE.md`).
4. Create `main.go` following the template above.
5. Register the lesson in `curriculum.v2.json`.
6. Update the section's root `README.md` with the new lesson.
7. Verify locally: `go run ./scripts/validate_curriculum.go && go test ./... && go vet ./...`
8. Push and request review.

## Adding an Exercise

Exercises include both complete solutions and starter stubs:

```text
NN-section-slug/
└── N-exercise-name/
    ├── README.md            ← Exercise requirements and walkthrough
    ├── main.go              ← Complete solution with comments
    ├── main_test.go         ← Tests that verify the solution
    └── _starter/
        └── main.go          ← TODO stubs for self-challenge
```

The `_starter/main.go` should:

1. Have the same file header and requirements checklist.
2. Contain function signatures with `// TODO: implement this` bodies.
3. Compile successfully (return zero-values from stubs).
4. Print a message directing students to the requirements.

---

## Commit Messages

Use clear, descriptive commit messages:

```text
feat(s07): add SY.1 sync.Mutex lesson
fix(s06): correct SQL injection example in DB.3
docs: update ROADMAP.md to match ARCHITECTURE.md v2.1
chore: update CI workflow for staticcheck
```

---

## Quality Checklist

Before submitting your PR for final review, verify:

- [ ] `go build ./...` passes
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes
- [ ] `go fmt ./...` shows no diff
- [ ] `go run ./scripts/validate_curriculum.go` passes
- [ ] Every new Go file has the standard header template and NEXT UP footer
- [ ] Every concept has inline teaching comments
- [ ] The section README is updated
- [ ] The lesson is registered in `curriculum.v2.json`
- [ ] CI pipeline passes on GitHub Actions
