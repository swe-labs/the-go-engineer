# Code Quality & Style Standards

This document establishes code quality and style standards for The Go Engineer curriculum.
All standards here follow the official Go conventions from [Effective Go](https://golang.org/doc/effective_go)
and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

---

## File Header Template

Every lesson `main.go` file must start with this header:

```go
// Copyright (c) 2026 Rasel Hossen
// Licensed under The Go Engineer License v1.0

// ============================================================================
// Section NN: Section Name — Lesson Title
// Level: Foundation | Core | Stretch
// Foundation = Phase 0–1, Core = Phase 2–3, Stretch = Phase 4
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

package main

import "fmt"

func main() {
    // KEY TAKEAWAY:
    // - Summary point 1
    // - Summary point 2
    fmt.Println("\n---------------------------------------------------")
    fmt.Println("NEXT UP: XY.N next-lesson-slug")
    fmt.Println("Run    : go run ./NN-section-slug/N-next-lesson-slug")
    fmt.Println("Current: XY.N (current-lesson-slug)")
    fmt.Println("---------------------------------------------------")
}
```

**NEXT UP footer format:** Must use `NEXT UP:` (without emoji) followed by the exact item ID and slug from `curriculum.v2.json`. The validator checks this format with the pattern `NEXT UP:\s*([A-Z]{2,6}\.\d+)`. Emoji in the footer will cause the validator to fail.

---

## Formatting Standards

### Code Formatting

All code must be formatted with `gofmt`. No exceptions.

```bash
go fmt ./...
```

Tabs for indentation (never spaces). Brace on same line as the statement. This is enforced by `gofmt` — do not configure your editor to override it.

### Import Organisation

```go
import (
    // Standard library — alphabetical
    "fmt"
    "strings"

    // Third-party packages — alphabetical
    "github.com/stretchr/testify/assert"

    // Internal packages — alphabetical
    "github.com/rasel9t6/the-go-engineer/internal/auth"
)
```

Blank line between each group. `goimports` or `gopls` handles this automatically on save.

### Line Length

**Soft limit: 100 characters.** Lines over 120 characters should be broken. Go prefers longer lines to avoid wrapping, but readability always wins.

---

## Naming Conventions

### Package Names

Lowercase, no underscores. Compound words are acceptable when conventional (e.g., `httputil`).

```go
// Correct
package user
package database
package httputil

// Wrong — avoid
package userutil       // compound name
package helpers        // vague
package utils          // vague
package db_helpers     // underscore
```

Short and descriptive. The package name is part of every call site (`user.Get`, not `userutil.GetUser`).

### Function Names

Exported: `PascalCase`. Unexported: `camelCase`. Verb + Noun pattern.

```go
// Correct
func GetUser(id int) (*User, error)
func validateEmail(email string) bool
func calculateTotal(items []Item) float64

// Wrong
func get_user(id int)
func Validate_Email(email string)     // underscore
func ComputeTot(items []Item)         // abbreviated, unclear
```

### Variable Names

Short and meaningful. Single-letter variables are acceptable in short scopes (`i` for index, `err` for error, `n` for count). Avoid generic names like `temp`, `tmp`, `x`, `data`.

```go
// Correct — descriptive
for i, item := range items { }
for name, value := range config { }

// Acceptable — conventional
for i, v := range values { }

// Wrong — too short in non-trivial scope
for idx, x := range items { }    // x tells us nothing
for k, v := range config { }     // k and v are barely better than nothing
```

### Constants

**Go constants use `PascalCase` for exported names and `camelCase` for unexported.** Go does NOT use `SCREAMING_SNAKE_CASE` for constants — that is a C/Java convention.

```go
// Correct Go convention
const MaxTimeout = 30 * time.Second
const DefaultPageSize = 50
const httpStatusOK = 200         // unexported: camelCase

// Wrong — not Go style
const MAX_TIMEOUT = 30           // C-style
const DEFAULT_PAGE_SIZE = 50     // C-style
const HTTP_STATUS_OK = 200       // C-style — staticcheck ST1003 flags this
```

The only exception is `iota`-based enum groups, where the convention follows the exported rule:

```go
type LogLevel int

const (
    LogDebug LogLevel = iota
    LogInfo
    LogWarn
    LogError
)
```

### Error Variables

Sentinel errors are `var`, not `const`, and use the `Err` prefix for exported errors:

```go
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidInput = errors.New("invalid input")
```

---

## Error Handling Patterns

See `docs/ENGINEERING_ERROR_FRAMEWORK.md` for the full three-tier framework (UserError / SystemError / FatalError).

### Pattern 1: Basic error check

```go
if err != nil {
    return err
}
```

### Pattern 2: Error wrapping with context

```go
if err != nil {
    return fmt.Errorf("open database (host=%s): %w", host, err)
}
```

Always use `%w` (not `%v`) when wrapping an error. `%w` preserves the error chain for `errors.Is()` and `errors.As()`.

### Pattern 3: Logging non-critical errors

```go
if err != nil {
    log.Printf("warning: failed to cleanup temporary file: %v", err)
    // Continue — this is non-fatal
}
```

### Pattern 4: The three-tier framework (Phase 2 onward)

```go
// Input validation → UserError
if !isValidEmail(email) {
    return nil, &UserError{Code: "invalid_email", Message: "Email format is incorrect"}
}

// Infrastructure failure → SystemError
if err := db.Query(...); err != nil {
    return nil, &SystemError{Code: "db_query_failed", Message: "Internal error", Cause: err}
}
```

---

## Comments

### Doc comments

Every exported symbol must have a doc comment starting with the symbol's name:

```go
// User represents an authenticated user in the system.
type User struct {
    ID    int
    Email string
}

// GetUser retrieves a user by their unique ID.
// Returns ErrUserNotFound if no user exists with the given ID.
// Returns an error wrapping the database error if the query fails.
func GetUser(id int) (*User, error) { ... }
```

### Teaching comments (curriculum-specific)

In lesson files, every non-trivial line gets a comment explaining WHY, not WHAT:

```go
// Pre-allocate the slice with the expected capacity.
// Without this, append() will copy the underlying array multiple times.
results := make([]Result, 0, len(items))
```

Not this:

```go
i++ // Increment i  ← useless: the reader can see that
```

### Cross-reference comments

When a lesson uses a concept taught elsewhere:

```go
// context.WithTimeout is covered in CT.3. Here we use it to prevent
// the database query from running indefinitely if the DB is slow.
ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
defer cancel()
```

### KEY TAKEAWAY comment

Every `main()` in a lesson should end with a summary comment:

```go
// KEY TAKEAWAY:
//   - Errors are values. The (value, error) return pattern is Go's idiom for explicit failure.
//   - Always check errors before using the returned value.
//   - Use fmt.Errorf("context: %w", err) to wrap errors with context.
fmt.Println("NEXT UP: FE.5 validation")
```

---

## Code Organisation

### Function order within a file

```go
package main

// 1. Types
type User struct { ... }

// 2. Constructors
func NewUser(name string) *User { ... }

// 3. Methods (receivers)
func (u *User) String() string { ... }

// 4. Unexported helpers
func validate(s string) bool { ... }

// 5. main() last
func main() { ... }
```

### Lesson file structure

```directory
lesson-name/
├── README.md          ← Written first. Learner reads this before opening main.go.
├── main.go            ← Primary lesson code
├── main_test.go       ← Tests (required for exercises)
└── _starter/
    └── main.go        ← TODO stubs (exercises only)
```

---

## Concurrency Standards

### Always use defer for WaitGroup.Done

```go
// Correct
go func() {
    defer wg.Done()
    doWork()
}()

// Wrong — Done() not called if doWork() panics
go func() {
    doWork()
    wg.Done()
}()
```

### Always pass WaitGroup by pointer

```go
func worker(wg *sync.WaitGroup) { // POINTER
    defer wg.Done()
}
go worker(&wg)
```

### Always cancel contexts

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel() // Always, even for timeouts — prevents goroutine leak
```

### Only the sender closes a channel

```go
// Correct — producer closes
go func() {
    defer close(ch)
    for _, item := range items {
        ch <- item
    }
}()

// Wrong — receiver closing causes panic if sender sends after
go func() {
    close(ch) // Never close a channel you don't own
}()
```

---

## Testing Standards

See `TESTING-STANDARDS.md` for full coverage requirements.

### Test function naming

```go
// Correct — describes the case being tested
func TestProcessValidInput(t *testing.T)
func TestProcessInvalidEmail(t *testing.T)

// Wrong — too vague
func TestProcess(t *testing.T)
func TestProcess1(t *testing.T)
```

### Table-driven test template

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"missing @", "userexample.com", true},
        {"empty string", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateEmail(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateEmail(%q) error = %v, wantErr = %v", tt.input, err, tt.wantErr)
            }
        })
    }
}
```

---

## Context Propagation

- `ctx context.Context` must ALWAYS be the first parameter of any function that performs I/O or may block.
- Never store `context.Context` inside a struct. It must flow through the call stack.
- Always pass a context to database queries (`db.QueryContext`), HTTP requests (`http.NewRequestWithContext`), and any operation that can timeout.

```go
// Correct
func GetUser(ctx context.Context, id int) (*User, error) {
    return db.QueryRowContext(ctx, "SELECT ...", id)
}

// Wrong — ignores cancellation
func GetUser(id int) (*User, error) {
    return db.QueryRow("SELECT ...", id)
}
```

---

## Generics Guidelines

- Use generics for data structures (`Stack[T]`, `Set[T]`) and utility functions that work identically regardless of type.
- Do NOT use generics when a standard interface (`io.Reader`, `fmt.Stringer`) solves the problem.
- Do NOT use generics when the code is clearer with concrete types.
- Taught in: s04 (types-design) — TI.9 through TI.17.

---

## Security Standards

- Never log raw passwords, tokens, session IDs, or PII. Pass them to `log.Printf` only after redaction.
- Never build SQL queries with string concatenation. Always use parameterised queries (`db.QueryContext(ctx, "... WHERE id = ?", id)`).
- Never write custom cryptographic functions. Use `golang.org/x/crypto/bcrypt` or `crypto/aes` from the standard library.
- Always validate and sanitise input at the HTTP boundary before any business logic.
- Taught in: s09 (architecture-security) — SEC.1 through SEC.10.

---

## Vet & Lint — Required Before Every Commit

```bash
go vet ./...         # Catch suspicious patterns
gofmt -d .           # Check formatting diff (no changes = correct)
go test -race ./...  # Check for race conditions
staticcheck ./...    # Install: go install honnef.co/go/tools/cmd/staticcheck@latest
```

The CI pipeline enforces all four. A PR that fails any of these will not be merged.

---

## Anti-Patterns to Avoid

### 1. Global mutable state

```go
// Wrong
var db *sql.DB
var config *Config

// Correct — use dependency injection through struct fields
type App struct {
    db     *sql.DB
    config *Config
}
```

### 2. Overusing `any` / `interface{}`

```go
// Wrong — type assertions everywhere
func Process(data any) any { ... }

// Correct — concrete types or proper interfaces
func Process(data []string) (map[string]int, error) { ... }
```

### 3. Named returns with bare `return`

```go
// Wrong — unclear what is returned
func process() (result string, err error) {
    result = "value"
    return
}

// Correct — explicit
func process() (result string, err error) {
    result = "value"
    return result, nil
}
```

### 4. Variable shadowing

```go
// Wrong — outer x and inner x are different variables
x := value
{
    x := otherValue  // Shadows outer x — confusing
}

// Correct
x := value
{
    y := otherValue  // Different name — clear intent
}
```

### 5. Ignoring errors

```go
// Wrong — error silently discarded
result, _ := doImportantWork()

// Correct — always handle
result, err := doImportantWork()
if err != nil {
    return fmt.Errorf("doing important work: %w", err)
}
```

---

## Code Review Checklist

- [ ] `go fmt ./...` — no formatting diff
- [ ] `go vet ./...` — no suspicious patterns
- [ ] `go test -race ./...` — no race conditions
- [ ] Imports organised in three groups (stdlib / third-party / internal)
- [ ] Naming follows Go convention (not C-style SCREAMING_SNAKE_CASE)
- [ ] Error handling is explicit — no ignored errors in meaningful paths
- [ ] Context is the first parameter in all I/O functions
- [ ] `defer` is used correctly — WaitGroup.Done, file.Close, cancel
- [ ] No global mutable state (exceptions must be documented)
- [ ] Exported symbols have doc comments
- [ ] Lesson files have header template and NEXT UP footer
- [ ] Tests exist for exercises (coverage > 75%)

---

## References

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Proverbs](https://go-proverbs.github.io/)
- [ENGINEERING_ERROR_FRAMEWORK.md](./docs/ENGINEERING_ERROR_FRAMEWORK.md)
- [TESTING-STANDARDS.md](./TESTING-STANDARDS.md)
