# Code Quality & Style Standards

This document establishes code quality and style standards for The Go Engineer curriculum.

All standards follow official Go conventions from Effective Go and Go Code Review Comments, with additional teaching requirements for this repository.

Architecture and curriculum metadata remain governed by `ARCHITECTURE.md`, `curriculum.v2.json`, and `CURRICULUM-BLUEPRINT.md`.

This document is the public code and authoring contract. If a lesson, README, validator, or authoring helper changes the expected quality bar, update this file first or in the same change.

## Standard Layers

The standard has two layers:

- **Machine-enforced rules**: formatting, curriculum metadata, lesson source headers, `RUN:` commands, `NEXT UP:` footers, Machine Role comment presence, README headings, local links, and validator-backed architecture contracts.
- **Review-enforced rules**: naming judgment, Machine Role comment quality, teaching clarity, production realism, security posture, and whether an example explains the engineering tradeoff instead of only demonstrating syntax.

High-quality lesson code should let a learner answer four questions without guessing:

- What role does this construct play in the program?
- Why does this choice exist instead of a simpler-looking alternative?
- What invariant, failure mode, or boundary does it protect?
- Where does this idea connect to the curriculum path?

## Curriculum Registry Standard

`curriculum.v2.json` is a source artifact, not a loose data dump.

Keep it canonical:

- top-level fields stay ordered as `schema_version`, `sections`, then `items`
- section objects and item objects keep the field order defined by the validator structs
- sections stay ordered `s00` through `s11`
- items stay grouped by section and ordered by path, then curriculum ID
- arrays use `[]` when empty; do not use `null`
- references do not duplicate values or point back to the same item
- paths stay under the owning section `path_prefix`
- slugs use lowercase kebab-case

The validator enforces these rules so the registry remains readable, reviewable, and stable across releases.

## File Header Template

Every learner-facing completed lesson `main.go` file must start with this header:

```go
// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section NN: Section Name - Lesson Title
// Level: Foundation | Core | Production | Stretch
// Use the exact item level from curriculum.v2.json.
// Do not infer level from section alone; section phase is separate metadata.
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

Starter files may use a smaller exercise-focused header when the learner is expected to fill in the implementation. Any completed runnable lesson surface should use the full header unless the lesson intentionally teaches a minimal file shape.

The validator reads completed lesson source headers. `Level` must match the item `level` in `curriculum.v2.json`, and `RUN:` must match that item's `run_command` exactly.

Both `RUN:` forms are accepted. Use the single-line form when it fits, and the two-line form when it keeps the header readable:

```go
// RUN: go run ./NN-section-slug/N-lesson-slug
```

```go
// RUN:
//   go run ./NN-section-slug/N-lesson-slug
```

`Level` is item-level difficulty and proof depth. It is different from section `phase`, which is defined in `curriculum.v2.json` as `foundations`, `engineering-core`, or `systems`.

Level labels mean:

- `Foundation`: first-principles material or a concept's first formal appearance.
- `Core`: day-to-day Go engineering skill that learners should internalize.
- `Production`: production-shaped engineering with reliability, deployment, service, persistence, observability, or operating concerns.
- `Stretch`: advanced depth that is useful but not required for the main path.

Section bands are a review aid, not the source of truth. Foundation and core lessons can appear inside later sections when a new domain starts. Production lessons are concentrated in backend, concurrency, architecture, production, and flagship sections. Stretch lessons are used only when the curriculum explicitly marks optional advanced depth. The registry value in `curriculum.v2.json` always wins.

Every lesson `main()` should end with a clear takeaway and terminal footer:

```go
// KEY TAKEAWAY:
// - Summary point 1
// - Summary point 2
fmt.Println()
fmt.Println("---------------------------------------------------")
fmt.Println("NEXT UP: XY.N -> NN-section-slug/N-next-lesson-slug")
fmt.Println("Run    : go run ./NN-section-slug/N-next-lesson-slug")
fmt.Println("Current: XY.N current-lesson-slug")
fmt.Println("---------------------------------------------------")
```

## NEXT UP Footer Format

The footer must use `NEXT UP:` exactly, without emoji.

The next item ID and path must match `curriculum.v2.json`.

The validator checks this pattern:

```text
NEXT UP:\s*([A-Z]{2,6}\.\d+)\s*->\s*([A-Za-z0-9._/\-]+)
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
	"github.com/swe-labs/the-go-engineer/internal/auth"
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
const maxTimeout = 30 * time.Second
const defaultPageSize = 50
const httpStatusOK = 200
```

Avoid:

```go
const MAX_TIMEOUT = 30
const DEFAULT_PAGE_SIZE = 50
```

### Error Codes

For public or production-shaped error codes, use stable uppercase machine-code values, not `SCREAMING_SNAKE_CASE` Go identifiers:

```go
const (
	errorCodeInvalidEmail   = "INVALID_EMAIL"
	errorCodeDBQueryFailed  = "DB_QUERY_FAILED"
	errorCodeInternalError  = "INTERNAL_ERROR"
)
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

Machine Role comments can satisfy this requirement for exported symbols when they start with the symbol name and explain the role clearly.

### Teaching comments

Lesson files explain why behavior exists.

Prefer comments that explain one of these:

- role: what job this construct performs for the program
- reason: why this approach is used here
- invariant: what must remain true
- boundary: where input, output, cancellation, ownership, or trust changes
- failure mode: what would break if the code changed carelessly

Good:

```go
// Pre-allocate the slice with the expected capacity so append avoids repeated reallocations.
results := make([]Result, 0, len(items))
```

Avoid:

```go
i++ // Increment i
```

### Machine Role Comments

For every major type, complex data structure, or core function used in a lesson, add a comment that explicitly defines its **Machine Role** or technical purpose. This links Go's syntax to its functional behavior.

Format: `// SymbolName (Tool Type): [direct technical role]`

Use the real symbol name without square brackets. For methods, use either `MethodName` or `Receiver.MethodName` when the receiver makes the role clearer.

Place the comment immediately above the declaration or local construct it explains.

Use this for:

- major structs, interfaces, functions, methods, and error types
- non-obvious slices, maps, channels, mutexes, contexts, goroutines, or pipelines
- test helpers that hide setup or cleanup complexity

Preferred tool type labels:

- `Struct`
- `Interface`
- `Function`
- `Method`
- `Constructor`
- `Error`
- `Slice`
- `Map`
- `Channel`
- `Mutex`
- `Context`
- `Goroutine`
- `Pipeline`
- `Test Helper`
- `Boundary`
- `Adapter`

Do not use Machine Role comments for every temporary variable or obvious operation. The comment should explain the runtime or design role, not restate the identifier.

Good:

```go
// ServerConfig (Struct): aggregates configuration state into one validated runtime boundary.
type ServerConfig struct { ... }

// UserStore (Interface): defines the persistence behavior the service needs without binding it to SQL.
type UserStore interface { ... }

// parseConfig (Function): transforms raw file text into validated key-value settings.
func parseConfig(content string) (map[string]string, error) { ... }

// ServeHTTP (Method): turns one HTTP request into validation, service work, and a response.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { ... }

// entries (Slice): provides a sortable view of map data so output order stays deterministic.
entries := make([]configEntry, 0, len(config))
```

Avoid:

```go
// ServerConfig is a server config.
type ServerConfig struct { ... }

// Increment i.
i++
```

This pattern is required for all new and refactored lessons to ensure a "Zero-Magic" learning experience.

### Cross-reference comments

When source code uses a concept taught elsewhere, reference the lesson ID and explain why the borrowed idea appears here:

```go
// context.WithTimeout is covered in CT.3. Here we use it to prevent the query
// from running indefinitely when the database is slow.
```

Source cross-references should be local and short. Do not create large comment blocks that retell another lesson. If a concept needs more than two source-comment sentences, move the explanation into the README and keep the code comment focused on why the concept is used at that line.

In `README.md` files, use GitHub-style alerts for cross-references:

- use `[!NOTE]` for prerequisite context, backward references, and gentle forward references
- use `[!TIP]` for actionable navigation, rerun suggestions, or learner practice advice
- keep the alert inside the relevant README section instead of creating a detached heading
- include the lesson ID and a clickable local `README.md` link when referencing a specific lesson
- do not put `[!NOTE]` or `[!TIP]` syntax inside Go source comments
- do not use legacy `Forward Reference` or `Backward Reference` labels

```markdown
> [!NOTE]
> This concept is covered in depth in [HC.1 What is a Program?](./00-how-computers-work/1-what-is-a-program/README.md).
```

```markdown
> [!TIP]
> If the terminal output feels surprising, rerun [GT.2 Hello World](./01-getting-started/2-hello-world/README.md) before continuing.
```

Avoid detached, standalone "Forward/Backward Reference" headlines.

## Lesson Proof Surface

Each lesson needs one coherent proof surface:

- `curriculum.v2.json` names the path, level, verification mode, run command, test command, starter path, prerequisites, and next items.
- The source header repeats the level and exact run command.
- The source footer repeats the next item as `NEXT UP:`.
- The README explains how to run the lesson and links to the next lesson's `README.md`.
- Tests, benchmarks, starter code, or rubric text prove the behavior appropriate to the lesson type.

If any one of those surfaces changes, update the others in the same change.

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

## Production-Shaped Code

As lessons move from foundation to production depth, examples should look like the kind of Go a learner can grow into:

- make boundaries explicit: input parsing, validation, persistence, network calls, cancellation, and output formatting should have visible ownership
- prefer deterministic output where tests or learners compare behavior
- clean up resources at the same level that acquired them
- keep timeouts, retries, backoff, idempotency, and shutdown behavior explicit when they are part of the lesson goal
- avoid adding frameworks or abstractions only to make code look enterprise-like
- choose the smallest realistic example that still teaches the production tradeoff

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
go test -coverprofile=coverage.out ./...
go run ./scripts/validate_curriculum.go
```

On PowerShell, quote the coverage flag if needed:

```powershell
go test "-coverprofile=coverage.out" ./...
```

Do not commit generated `coverage.out` or `coverage.html` artifacts.

For benchmark-related changes:

```bash
go test -bench=. -benchmem -count=1 ./08-quality-test/01-quality-and-performance/testing/benchmarks/
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
- [ ] `Level` and `RUN:` headers match `curriculum.v2.json`.
- [ ] Machine Role comments explain role, boundary, invariant, or failure mode without restating syntax.
- [ ] Source cross-references use lesson IDs and explain local relevance.
- [ ] `NEXT UP:` footers match curriculum metadata.
- [ ] README cross-references use `[!NOTE]` or `[!TIP]` alerts.
- [ ] README `Next Step` entries use clickable links to the next `README.md`.
- [ ] Curriculum metadata, source header/footer, README run instructions, and tests describe the same proof surface.
- [ ] Tests exist for exercises and behavior changes.
- [ ] Full PR-readiness checks pass, including coverage generation.
- [ ] No secrets or sensitive data are logged.
