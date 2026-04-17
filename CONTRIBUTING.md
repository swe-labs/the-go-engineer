# Contributing to The Go Engineer

Thank you for your interest in contributing! This guide will help you get started.

## Quick Links

- **[Testing Standards](./TESTING-STANDARDS.md)** - How to write tests
- **[Code Quality Standards](./CODE-STANDARDS.md)** - Code style and best practices
- **[Curriculum Blueprint](./CURRICULUM-BLUEPRINT.md)** - How we teach and deliver lessons
- **[Architecture Blueprint](./ARCHITECTURE.md)** - The 21-section v2 structure

## Getting Started

```bash
# Clone the repository
git clone https://github.com/rasel9t6/the-go-engineer.git
cd the-go-engineer

# Verify your environment
go version
go test ./...
```

## The Strict GitHub Workflow

To maintain a high standard of quality, all contributors (including maintainers) MUST follow this workflow:

### 1. Create an Issue First
Never start coding without an approved issue.
- Ensure the issue maps to a specific gap or lesson in the 21-section blueprint.
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

## Code Style

### Every Go File Must Follow This Template

```go
package main

import "fmt"

// ============================================================================
// Section N: Section Name — Lesson Title
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
5. **Cross-reference sections** — Use `(See Section 08: Interface Contracts)` when a concept is covered elsewhere

### Formatting Rules

- Run `gofmt -s -w .` on every file
- Run `go vet ./...` before committing

## Section Numbering

- Sections are numbered `01-21` according to `ARCHITECTURE.md`.
- Lessons within a section are numbered `1-`, `2-`, etc.
- Exercises are the LAST numbered item in a section.

## Adding a New Lesson

1. Ensure an Issue is open and approved.
2. Create the directory: `NN-section-name/N-lesson-name/`
3. Create the `README.md` explaining the mental model (Doc first, then code).
4. Create `main.go` following the template above.
5. Update the section's root `README.md` with the new lesson in the path table.
6. Verify locally: `go test ./...` and `go vet ./...`
7. Push and request review.

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

1. Have the same file header and requirements checklist.
2. Contain function signatures with `// TODO: implement this` bodies.
3. Compile successfully (return zero-values from stubs).
4. Print a message directing students to the requirements.

## Commit Messages

Use clear, descriptive commit messages:

```text
feat: add Section 17 Context deep-dive (4 lessons)
fix: go vet warning in 07-strings formatting
docs: update comments in Section 04 control flow
```

## Quality Checklist

Before submitting your PR for final review, verify:

- [ ] `go build ./...` passes
- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes
- [ ] Code is formatted with `gofmt`
- [ ] Every new Go file has the standard header template and "NEXT UP" footer
- [ ] Every concept has inline teaching comments
- [ ] The section README is updated
- [ ] CI pipeline passes on GitHub Actions
