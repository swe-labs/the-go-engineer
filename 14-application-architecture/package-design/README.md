# Section 18: Package Design

## Learning Objectives

Good package design is what separates beginner Go code from production-grade software. This section teaches you how to organize a real Go application into clean, testable, reusable packages.

## Beginner → Expert Mapping

| Topic | Level | Importance | Engineering Concept |
| ----- | ----- | ---------- | ------------------- |
| Package naming | Beginner | High | Conventions and anti-patterns |
| Export rules | Beginner | **Critical** | Uppercase = public, lowercase = private |
| Internal packages | Intermediate | High | `internal/` visibility restriction |
| Project layout | Intermediate | **Critical** | The standard Go project structure |
| Dependency Inversion | Advanced | **Critical** | Interfaces at boundaries |

## Engineering Depth

Go packages are not just folders — they are the unit of compilation, testing, documentation, and visibility. A well-designed package:

1. Has a **clear single responsibility** (e.g., `auth`, `storage`, `email`)
2. Exports a **small surface area** (few public types/functions)
3. Defines **interfaces at the consumer**, not the provider
4. Uses `internal/` to prevent external access to implementation details
5. Avoids circular dependencies (Go's compiler forbids them)

## Contents

| Directory | Topic | Level |
| --------- | ----- | ----- |
| `1-naming/` | Package naming conventions and anti-patterns | Beginner |
| `2-visibility/` | Export rules, internal packages | Intermediate |
| `3-project-layout/` | Standard Go project structure | Intermediate |

## How to Run

```bash
go run ./18-package-design/1-naming
go run ./18-package-design/2-visibility
go run ./18-package-design/3-project-layout
```

## References

- [Effective Go: Package Names](https://go.dev/doc/effective_go#package-names)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Blog: Package Names](https://go.dev/blog/package-names)
