# Code Quality & Style Standards

This document establishes code quality and style standards for The Go Engineer curriculum.

All standards follow official Go conventions from Effective Go and Go Code Review Comments, with additional teaching requirements for this repository.

## File Header Template

Every lesson `main.go` file must start with this header:

```go
// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section NN: Section Name - Lesson Title
// Level: Foundation | Core | Stretch
// Foundation = Phase 0-1, Core = Phase 2-3, Stretch = Phase 4
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Concept one described precisely
//   - Concept two described precisely
//   - Concept three described precisely
//
// WHY THIS MATTERS:
//   One sentence on why this exists in a real production Go codebase.
//
// RUN: go run ./NN-section-slug/N-lesson-slug
// ============================================================================
```

Every lesson `main()` should end with a clear takeaway and terminal footer:

```go
// KEY TAKEAWAY:
// - Summary point 1
// - Summary point 2
fmt.Println()
fmt.Println("---------------------------------------------------")
fmt.Println("NEXT UP: XY.N next-lesson-slug")
fmt.Println("Run    : go run ./NN-section-slug/N-next-lesson-slug")
fmt.Println("Current: XY.N current-lesson-slug")
fmt.Println("---------------------------------------------------")
```

## NEXT UP Footer Format

The footer must use `NEXT UP:` exactly, without emoji.

The next item ID and slug must match `curriculum.v2.json`.

The validator checks this pattern:

```text
NEXT UP:\s*([A-Z]{2,6}\.\d+)
```

## Formatting Standards

All Go code must be formatted with `gofmt`.

```bash
gofmt -w .
```

CI checks formatting with:

```bash
unformatted=$(gofmt -l .)
test -z "$unformatted"
```

Rules:

- tabs for indentation
- brace on the same line as the statement
- soft line limit of 100 characters
- break lines above 120 characters unless readability clearly improves otherwise

## Import Organization

Use this grouping when multiple groups exist:

```go
import (
	// Standard library
	"context"
	"fmt"
	"strings"

	// Third-party
	"github.com/stretchr/testify/assert"

	// Internal
	"github.com/rasel9t6/the-go-engineer/internal/auth"
)
```

Use `goimports` or `gopls` when available.

## Naming Conventions

### Packages

Package names are lowercase, short, and descriptive.

Good:

```go
package user
package database
package httputil
```

Avoid:

```go
package userutil
package helpers
package utils
package db_helpers
```

### Functions

Exported functions use `PascalCase`. Unexported functions use `camelCase`.

Prefer verb + noun naming:

```go
func GetUser(id int) (*User, error)
func validateEmail(email string) bool
func calculateTotal(items []Item) float64
```

### Variables

Use short names in short scopes and descriptive names in wider scopes.

Avoid vague names like `data`, `tmp`, `x`, or `result` when the surrounding context does not make them obvious.

### Constants

Go constants do not use `SCREAMING_SNAKE_CASE`.

Good:

```go
const MaxTimeout = 30 * time.Second
const DefaultPageSize = 50
const httpStatusOK = 200
```

Avoid:

```go
const MAX_TIMEOUT = 30
const DEFAULT_PAGE_SIZE = 50
```

### Error Codes

For public or production-shaped error codes, use stable uppercase machine codes:

```go
INVALID_EMAIL
DB_QUERY_FAILED
INTERNAL_ERROR
```

For local lesson-only examples, keep the convention consistent inside that lesson.

## Error Handling

Use explicit error handling.

### Basic check

```go
if err != nil {
	return err
}
```

### Wrapping

Use `%w` when preserving the cause matters:

```go
if err != nil {
	return fmt.Errorf("open database: %w", err)
}
```

### Three-tier framework

For backend and production-shaped code, follow `docs/ENGINEERING_ERROR_FRAMEWORK.md`:

- UserError: validation and business rule failures
- SystemError: infrastructure and external failures
- FatalError: unrecoverable startup or invariant failures

## Comments

### Doc comments

Every exported symbol must have a doc comment that starts with the symbol name.

```go
// User represents an authenticated user in the system.
type User struct {
	ID    int
	Email string
}
```

### Teaching comments

Lesson files explain why behavior exists.

Good:

```go
// Pre-allocate the slice with the expected capacity so append avoids repeated reallocations.
results := make([]Result, 0, len(items))
```

Avoid:

```go
i++ // Increment i
```

### Cross-reference comments

When a lesson uses a concept taught elsewhere, reference the lesson ID:

```go
// context.WithTimeout is covered in CT.3. Here we use it to prevent the query
// from running indefinitely when the database is slow.
```

## Code Organization

Preferred order:

1. package declaration
2. imports
3. constants
4. variables
5. types
6. constructors
7. methods
8. unexported helpers
9. `main()`

## Lesson Directory Shape

```text
lesson-name/
|-- README.md
|-- main.go
|-- main_test.go
`-- _starter/
    `-- main.go
```

`main_test.go` and `_starter/` are required for exercises and strongly recommended where behavior should be proven.

## Concurrency Standards

- Use `defer wg.Done()`.
- Pass `sync.WaitGroup` by pointer.
- Always call context cancellation functions.
- Only the sending owner closes a channel.
- Avoid unbounded goroutine creation.
- Avoid sleeps for synchronization in tests.
- Run race-sensitive changes with `go test -race ./...`.

## Context Propagation

`ctx context.Context` must be the first parameter of functions that perform I/O, call external services, or may block.

Good:

```go
func GetUser(ctx context.Context, id int) (*User, error)
```

Avoid storing `context.Context` inside structs.

## Security Standards

- Never log raw passwords, tokens, session IDs, or secrets.
- Never build SQL queries through string concatenation.
- Validate and sanitize input at HTTP boundaries.
- Use standard cryptography libraries; do not write custom crypto.
- Do not leak internal error details through user-facing responses.
- Enforce tenant/user scoping where applicable.

## Required Local Checks

Before final review:

```bash
go build ./...
go vet ./...
unformatted=$(gofmt -l .); test -z "$unformatted" || (echo "$unformatted" && exit 1)
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go test -race ./...
go run ./scripts/validate_curriculum.go
```

Recommended when installed:

```bash
staticcheck ./...
```

Staticcheck is recommended locally. It should only be described as required if CI is updated to enforce it.

## Review Checklist

- [ ] Code is formatted.
- [ ] Names follow Go conventions.
- [ ] Errors are handled explicitly.
- [ ] Error wrapping preserves causes where needed.
- [ ] Context flows through I/O and blocking paths.
- [ ] Resources are closed.
- [ ] Concurrency avoids leaks and races.
- [ ] Exported symbols have doc comments.
- [ ] Lesson files have standard headers.
- [ ] `NEXT UP:` footers match curriculum metadata.
- [ ] Tests exist for exercises and behavior changes.
- [ ] No secrets or sensitive data are logged.
