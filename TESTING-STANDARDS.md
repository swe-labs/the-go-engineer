# Testing Standards & Guidelines

> This document establishes testing standards for The Go Engineer curriculum.
> Section references use the v2.1 architecture from `ARCHITECTURE.md`.
> Testing is taught formally in s08 (quality-test): TE.1–TE.10 and PR.1–PR.5.

---

## Overview

Every lesson should include appropriate tests that demonstrate:

- ✅ Core functionality
- ✅ Edge cases and error handling
- ✅ Integration points (where applicable)
- ✅ Performance characteristics (for performance-critical code)

---

## Testing Hierarchy

### Level 1: Unit Tests (Required for all exercises)

**Purpose**: Test individual functions and methods in isolation

**File**: `*_test.go` in the same package

**Taught in**: s08 — TE.1: Unit Testing Basics

```go
package user

import "testing"

func TestCheckUsername(t *testing.T) {
    tests := []struct {
        input    string
        expected bool
    }{
        {"validuser", true},
        {"user1", false},
        {"admin", false},
    }
    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            if got := CheckUsername(tt.input); got != tt.expected {
                t.Errorf("got %v, want %v", got, tt.expected)
            }
        })
    }
}
```

### Level 2: Table-Driven Tests (Required for > 2 cases)

**Purpose**: Test multiple scenarios systematically

**Taught in**: s08 — TE.2: Table-Driven Tests

```go
func TestFunction(t *testing.T) {
    cases := []struct {
        name     string  // Description of test case
        input    string  // Input value
        expected string  // Expected output
        wantErr  bool    // Whether error expected
    }{
        {"valid case", "input", "output", false},
        {"invalid case", "bad", "", true},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got, err := Function(tc.input)
            if (err != nil) != tc.wantErr {
                t.Fatalf("unexpected error: %v", err)
            }
            if got != tc.expected {
                t.Errorf("got %q, want %q", got, tc.expected)
            }
        })
    }
}
```

### Level 3: HTTP Handler Tests (Required for s06 onward)

**Purpose**: Test HTTP handlers without running a server

**Taught in**: s08 — TE.3: HTTP Handler Testing

```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/api/users", nil)
    w := httptest.NewRecorder()

    handler(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("status code %d, want %d", w.Code, http.StatusOK)
    }
}
```

### Level 4: Benchmarks (Required for performance-critical code)

**Purpose**: Measure and track performance

**Taught in**: s08 — TE.4: Benchmarking; PR.1–PR.5

```go
func BenchmarkSort(b *testing.B) {
    data := generateData()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        sort.Slice(data, func(i, j int) bool {
            return data[i] < data[j]
        })
    }
}
```

**Run**: `go test -bench=. -benchmem`

### Level 5: Integration Tests (Required for multi-component lessons)

**Purpose**: Test integration between multiple packages/components

**Taught in**: s08 — TE.9: Integration Tests

```go
func TestDatabaseOperations(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    // Insert
    id, err := db.InsertUser("john")
    if err != nil {
        t.Fatalf("insert failed: %v", err)
    }

    // Query
    user, err := db.GetUser(id)
    if err != nil {
        t.Fatalf("query failed: %v", err)
    }

    if user.Name != "john" {
        t.Errorf("got %q, want %q", user.Name, "john")
    }
}
```

### Level 6: Fuzz Tests (Required for parsing/validation code)

**Purpose**: Test with randomly generated inputs to find edge cases

**Taught in**: s08 — TE.6: Fuzz Testing (Go 1.18+)

```go
func FuzzParseInt(f *testing.F) {
    // Add seed corpus
    f.Add("123")
    f.Add("-1")
    f.Add("0")

    f.Fuzz(func(t *testing.T, s string) {
        _, err := ParseInt(s)
        if err == nil {
            // Ensure no panic occurred and behavior is handled
        }
    })
}
```

### Level 7: API & gRPC Tests (Required for s06 endpoints)

**Purpose**: Test APIs using HTTP clients or gRPC test clients

**Taught in**: s06 — HS.10: REST API; API.9: gRPC Service

```go
func TestGRPCServer(t *testing.T) {
    client := setupTestGRPCClient(t)

    resp, err := client.GetUser(context.Background(), &pb.GetUserRequest{Id: "123"})
    assert.NoError(t, err)
    assert.NotNil(t, resp)
}
```

---

## Error Handling in Tests

### Consistent Error Patterns

Use these patterns throughout tests:

```go
// For setup failures (use Fatal — stop the test immediately)
if err := setupTest(); err != nil {
    t.Fatalf("setup failed: %v", err)
}

// For unexpected errors (use Errorf — continue to check other assertions)
if err != nil {
    t.Errorf("unexpected error: %v", err)
}

// For assertion failures
if got != want {
    t.Errorf("got %v, want %v", got, want)
}

// For checking both error and value
if err != nil {
    t.Fatalf("got error: %v", err)
}
if got != want {
    t.Errorf("got %v, want %v", got, want)
}
```

---

## Using testify/assert

For cleaner tests, use the testify/assert library (already in go.mod):

```go
import "github.com/stretchr/testify/assert"

func TestWithAssert(t *testing.T) {
    result, err := Function()
    assert.NoError(t, err)
    assert.Equal(t, "expected", result)
    assert.True(t, someCondition)
}
```

---

## Test Helpers

Create helper functions for common test setup:

```go
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper() // Marks this as a helper — errors point to the caller, not here
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open database: %v", err)
    }
    if err := initSchema(db); err != nil {
        t.Fatalf("failed to init schema: %v", err)
    }
    t.Cleanup(func() { db.Close() }) // Automatic cleanup when test ends
    return db
}
```

---

## Code Coverage

### Computing Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage percentage
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

### Coverage Goals by Section

| Section                                   | Minimum Coverage |
| ----------------------------------------- | ---------------- |
| Phase 0–1: Foundation (s00–s04)           | 80%              |
| s08: Quality & Testing                    | 95%              |
| s06: Backend, APIs & Databases            | 85%              |
| s09: Architecture & Security              | 85%              |
| s11: Opslane Flagship                     | 90%              |

---

## Checklist for Lesson Authors

When creating a new lesson, ensure:

- [ ] Unit tests written for all exported functions
- [ ] Table-driven tests for multiple cases
- [ ] Edge cases covered (empty inputs, nil, max values)
- [ ] Error cases tested separately
- [ ] All error paths covered
- [ ] Tests use consistent naming: `TestFunctionName`
- [ ] Tests use `t.Run` for sub-tests
- [ ] Tests use `t.Helper()` in helper functions
- [ ] Tests use proper error messages with context
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Code passes vet checks (`go vet ./...`)
- [ ] Coverage > 75% (aim for higher)
- [ ] Benchmarks added for performance-critical code
- [ ] Helper functions documented
- [ ] Test fixtures use descriptive names

---

## Running Tests Locally

```bash
# Run all tests
make test

# Run with race detector
make test-race

# Run specific test
go test -run TestName ./path/to/package

# Run specific sub-test
go test -run TestName/SubName ./path/to/package

# Verbose output
go test -v ./...

# With timeout
go test -timeout 30s ./...

# Generate coverage
make cover
```

---

## CI/CD Integration

The CI pipeline automatically:

- Runs all tests on every push
- Checks formatting with `gofmt`
- Runs race detector (`-race`)
- Generates coverage reports
- Runs `staticcheck`
- Validates curriculum structure (`go run ./scripts/validate_curriculum.go`)

See `.github/workflows/ci.yml` for details.

---

## Common Testing Mistakes to Avoid

❌ **DON'T**: Use global variables in tests

```go
var db *sql.DB  // ← Bad: shared state across tests
```

✅ **DO**: Setup fresh state for each test

```go
func TestFunction(t *testing.T) {
    db := setupTestDB(t)  // ← Good: isolated state
    // ...
}
```

---

❌ **DON'T**: Ignore errors without comment

```go
_ = err  // ← Bad: error silently ignored
```

✅ **DO**: Explicitly handle or check errors

```go
assert.NoError(t, err)  // ← Good: explicit check
```

---

❌ **DON'T**: Use sleeps for synchronisation

```go
time.Sleep(100 * time.Millisecond)  // ← Bad: flaky test
```

✅ **DO**: Use channels or proper synchronisation

```go
<-done  // ← Good: synchronisation primitive
```

---

## References

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [testify/assert](https://github.com/stretchr/testify#assert)
- [httptest Package](https://golang.org/pkg/net/http/httptest/)
- [Go Benchmark Timers](https://golang.org/pkg/testing/#B)
- [ARCHITECTURE.md](./ARCHITECTURE.md) — s08 lesson plan
- [CODE-STANDARDS.md](./CODE-STANDARDS.md) — code review checklist
