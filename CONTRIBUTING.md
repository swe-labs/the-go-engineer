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

## Adding a New Lesson

1. Create a directory: `NN-section-name/N-lesson-name/`
2. Create `main.go` following the template above
3. Update the section's `README.md` with the new lesson in the `Learning Path` table
4. Update `curriculum.json` with the new lesson ID, concept, and prerequisites
5. Run `go run scripts/validate_curriculum.go` to verify the mapping
6. Update the root `README.md` if adding a new section
7. Verify: `make build && make test && make lint`

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
