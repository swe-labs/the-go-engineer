# Testing Standards & Guidelines

This document establishes testing standards for The Go Engineer curriculum to ensure consistent quality and comprehensive coverage.

## Overview

Every lesson should include appropriate tests that demonstrate:
- ✅ Core functionality
- ✅ Edge cases and error handling
- ✅ Integration points (where applicable)
- ✅ Performance characteristics (for performance-critical code)

## Testing Hierarchy

### Level 1: Unit Tests (Required for all lessons)

**Purpose**: Test individual functions and methods in isolation

**File**: `*_test.go` in the same package

**Example**:
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

**File**: `*_test.go`

**Pattern**:
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

### Level 3: HTTP Handler Tests (Required for web code)

**Purpose**: Test HTTP handlers without running a server

**File**: `*_test.go`

**Example**:
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

**File**: `*_test.go`

**Example**:
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

**File**: `integration_test.go` (in a `*_test` build tag file or separate test file)

**Example**:
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

## Error Handling in Tests

### Consistent Error Patterns

Use these patterns throughout tests:

```go
// For setup failures (use Fatal)
if err := setupTest(); err != nil {
    t.Fatalf("setup failed: %v", err)
}

// For unexpected errors
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

## Test Helpers

Create helper functions for common test setup:

```go
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open database: %v", err)
    }
    // Initialize schema
    if err := initSchema(db); err != nil {
        t.Fatalf("failed to init schema: %v", err)
    }
    return db
}
```

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

| Section | Minimum Coverage |
|---------|-----------------|
| Core Foundations | 80% |
| Functions & Errors | 85% |
| Data Structures | 80% |
| Web & Database | 75% |
| Testing (§13) | 95% |

## Checklist for Lesson Authors

When creating a new lesson, ensure:

- [ ] Unit tests written for all exported functions
- [ ] Table-driven tests for multiple cases
- [ ] Edge cases covered (empty inputs, nil, max values)
- [ ] Error cases tested separately
- [ ] All error paths covered
- [ ] Tests use consistent naming: `TestFunctionName`
- [ ] Tests use t.Run for sub-tests
- [ ] Tests use proper error messages
- [ ] Code is formatted (`gofmt`)
- [ ] Code passes vet checks (`go vet`)
- [ ] Coverage > 75% (aim for higher)
- [ ] Benchmarks added for performance-critical code
- [ ] Helper functions documented
- [ ] Test fixtures use descriptive names

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

## CI/CD Integration

The CI pipeline automatically:
- Runs all tests on every push
- Checks formatting with gofmt
- Runs race detector
- Generates coverage reports
- Runs benchmarks
- Validates curriculum structure

See `.github/workflows/ci.yml` for details.

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

❌ **DON'T**: Use sleeps for synchronization
```go
time.Sleep(100 * time.Millisecond)  // ← Bad: flaky test
```

✅ **DO**: Use channels or proper synchronization
```go
<-done  // ← Good: synchronization primitive
```

---

## References

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [testify/assert](https://github.com/stretchr/testify#assert)
- [httptest Package](https://golang.org/pkg/net/http/httptest/)
- [Go Benchmark Timers](https://golang.org/pkg/testing/#B)
