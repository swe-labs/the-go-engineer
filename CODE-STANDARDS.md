# Code Quality & Style Standards

This document establishes code quality and style standards for The Go Engineer curriculum.

## File Header Template

Every lesson file must start with this header:

```go
// ============================================================================
// Section N: Section Name — Lesson Title
// Level: Beginner | Intermediate | Advanced
// ============================================================================
//
// WHAT YOU'LL LEARN:
//   - Concept 1
//   - Concept 2
//   - Concept 3
//
// ENGINEERING DEPTH:
//   Why this matters in production Go. What problems does it solve?
//   Include real-world context and examples.
//
// RUN: go run ./path/to/lesson
// ============================================================================

package main

import "fmt"

func main() {
    // ...
}
```

## Formatting Standards

### Code Formatting

**Required**: All code must be formatted with `gofmt`

```bash
gofmt -w .
make fmt
```

No custom formatting allowed. Go's standard formatter ensures consistency across all lessons.

### Import Organization

```go
import (
    // Standard library
    "fmt"
    "strings"

    // Third-party packages
    "github.com/stretchr/testify/assert"

    // Internal packages
    "github.com/rasel9t6/the-go-engineer/pkg/util"
)
```

Order:
1. Standard library (alphabetical)
2. Third-party packages (alphabetical)
3. Internal packages (alphabetical)

### Line Length

**Maximum**: 120 characters (soft limit, 100 is better)

Go prefers longer lines to avoid wrapping, but readability comes first.

### Naming Conventions

#### Package Names
- Lowercase, no underscores
- Short and descriptive: `user`, `config`, `storage`
- Avoid: `util`, `utils`, `helper`, `common`

```go
// ✅ Good
package user
package database
package http

// ❌ Avoid
package userutil
package helpers
package db_helpers
```

#### Function Names
- Exported functions: PascalCase
- Unexported functions: camelCase
- Verb + Noun: `GetUser`, `validateEmail`, `parseConfig`

```go
// ✅ Good
func GetUser(id int) (*User, error)
func validateEmail(email string) bool
func calculateTotal(items []Item) float64

// ❌ Avoid
func get_user(id int)
func ValidateEmail(email string)
func ComputeTot(items []Item)
```

#### Variable Names
- Short and meaningful: `u` for user, `err` for error, `count` for item count
- Loop variables: `i`, `j` for index; otherwise descriptive
- Avoid: `temp`, `tmp`, `x`, `data`

```go
// ✅ Good
for i, item := range items { }
for name, value := range config { }
for err := range errors { }

// ❌ Avoid
for idx, x := range items { }
for k, v := range config { }
for e := range errors { }
```

#### Constants
- UPPER_CASE with underscores for multi-word
- Or PascalCase if single word

```go
// ✅ Good
const MaxTimeout = 30 * time.Second
const DefaultPageSize = 50
const HTTP_STATUS_OK = 200

// ❌ Avoid
const max_timeout = 30
const DEFAULTPAGESIZE = 50
```

## Error Handling Patterns

### Pattern 1: Basic Error Check

```go
if err != nil {
    return err
}
```

This is the standard pattern - short and clear.

### Pattern 2: Wrapped Errors

```go
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}
```

Always use `%w` for error wrapping (Go 1.13+).

### Pattern 3: With Context

```go
if err != nil {
    return fmt.Errorf("open database (host=%s, port=%d): %w", host, port, err)
}
```

Include relevant context to help debugging.

### Pattern 4: Logging

```go
if err != nil {
    log.Printf("warning: failed to cleanup: %v", err)
    // Continue execution after logging
}
```

Use logging for non-critical errors that shouldn't stop execution.

### Pattern 5: Custom Error Types

```go
// Define custom error
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}

// Use it
if !isValid(email) {
    return ValidationError{Field: "email", Message: "invalid format"}
}
```

## Comments

### Comment Standards

**Every exported symbol must have a doc comment:**

```go
// User represents a user in the system.
type User struct {
    ID    int
    Email string
}

// NewUser creates a new user with the given email.
func NewUser(email string) *User {
    return &User{Email: email}
}
```

**Complex logic needs inline comments:**

```go
func process(items []Item) {
    // Sort by creation time, newest first
    sort.Slice(items, func(i, j int) bool {
        return items[i].CreatedAt.After(items[j].CreatedAt)
    })
    
    // Filter to only published items
    var published []Item
    for _, item := range items {
        if item.Published {
            published = append(published, item)
        }
    }
}
```

**DO NOT**: Comment obvious code

```go
// ❌ Bad: Obvious to anyone reading the code
i++ // Increment i

// ✅ Good: Explains why, not what
i++ // Skip first occurrence due to header row
```

### Doc Comment Format

```go
// Function describes what the function does.
// It can span multiple lines.
//
// Parameters are documented in the description if complex.
//
// Returns are documented if not obvious.
//
// Example:
//  result := Function(input1, input2)
//
// See also: RelatedFunction, AnotherFunction
func Function(input string) error {
```

## Code Organization

### Function Order Within a File

```go
package main

import "fmt"

// Types first
type User struct {
    Name string
}

// Constructors
func NewUser(name string) *User {
    return &User{Name: name}
}

// Receivers (methods)
func (u *User) String() string {
    return u.Name
}

// Helper functions
func validate(s string) bool {
    return s != ""
}

// Main function (if applicable)
func main() {
    // ...
}
```

### File Organization

```
lesson/
├── main.go              # Primary implementation
├── types.go             # Custom types (if > 100 lines)
├── handlers.go          # HTTP handlers (if relevant)
├── repository.go        # Database logic (if relevant)
├── main_test.go         # Tests
└── _starter/            # Optional: starting point
    └── main.go
```

## Concurrency

### Goroutine Patterns

**Always defer channel close:**

```go
// ✅ Good
ch := make(chan string)
go func() {
    defer close(ch)
    ch <- "hello"
}()

// ❌ Avoid: Causes panic or deadlock
go func() {
    ch <- "hello"
    close(ch)  // Not using defer
}()
```

**Always use WaitGroup pointer:**

```go
// ✅ Good
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // ...
}()
wg.Wait()

// ❌ Avoid: wg copied, not shared
wg := sync.WaitGroup{}
go func() {
    defer wg.Done()  // Doesn't affect original wg!
}()
wg.Wait()
```

**Context cancellation:**

```go
// ✅ Good
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
    select {
    case <-ctx.Done():
        return
    // ... do work ...
    }
}()
```

## Testing Patterns

### Test Function Naming

```go
// ✅ Good
func TestProcessValidInput(t *testing.T)
func TestProcessInvalidEmail(t *testing.T)
func TestProcessEmptyString(t *testing.T)

// ❌ Avoid
func TestProcess(t *testing.T)  // Too vague
func TestProcess1(t *testing.T)  // Use names, not numbers
func TestFunctionDoesStuff(t *testing.T)  // Too generic
```

### Table-Driven Test Template

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name      string   // Describe the case
        input     string   // Input value
        expected  int      // Expected output
        wantErr   bool     // Expect an error?
    }{
        {"valid case", "input", 42, false},
        {"invalid case", "bad", 0, true},
        {"edge case", "", 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Fatalf("unexpected error: %v", err)
            }
            if err == nil && got != tt.expected {
                t.Errorf("got %d, want %d", got, tt.expected)
            }
        })
    }
}
```

## Vet & Lint

**All code must pass:**

```bash
go vet ./...         # Check for suspicious patterns
gofmt -d .          # Check formatting
go test -race ./... # Check for race conditions (optional, use in CI)
```

**Before committing:**

```bash
make lint
```

## Checklist for Code Review

- [ ] Code formatted with `gofmt`
- [ ] All imports organized correctly
- [ ] Naming conventions followed
- [ ] Error handling is consistent
- [ ] Comments on exported symbols
- [ ] No global state (except in specific cases)
- [ ] Concurrency safe (no races)
- [ ] Tests written and passing
- [ ] Coverage > 75%
- [ ] No unused variables or imports
- [ ] Consistent with other lessons

## Anti-Patterns to Avoid

### 1. Global Variables

```go
// ❌ Bad
var db *sql.DB
var config *Config

// ✅ Good
type App struct {
    db     *sql.DB
    config *Config
}
```

### 2. Interface{} Overuse

```go
// ❌ Bad
func Process(data interface{}) interface{} {
    // Type assertions everywhere
}

// ✅ Good
func Process(data []string) (map[string]int, error) {
    // Clear types
}
```

### 3. Bare `return`

```go
// ❌ Avoid in functions with named returns
func Process() (result string, err error) {
    result = "value"
    return  // Unclear what's returned
}

// ✅ Good
func Process() (result string, err error) {
    result = "value"
    return result, nil  // Clear
}
```

### 4. Shadowing Variables

```go
// ❌ Bad
x := value
{
    x := otherValue  // Shadows outer x
    // ...
}
// Is this the original x or the inner one?

// ✅ Good
x := value
{
    y := otherValue  // New variable
    // ...
}
```

## Resources

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Proverbs](https://go-proverbs.github.io/)
