# Testing Standards & Guidelines

This document establishes testing standards for The Go Engineer curriculum.

Section references follow the locked v2.1 architecture from `ARCHITECTURE.md`.

## Testing Promise

Tests should prove behavior, not merely execute code.

Every lesson or exercise should include the level of verification appropriate to its section, complexity, and learner goal.

## Testing Hierarchy

### Level 1: Unit Tests

Required for exercises and reusable logic.

```go
func TestCheckUsername(t *testing.T) {
	tests := []struct {
		name string
		input string
		want bool
	}{
		{name: "valid username", input: "validuser", want: true},
		{name: "reserved username", input: "admin", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckUsername(tt.input); got != tt.want {
				t.Errorf("CheckUsername(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
```

### Level 2: Table-Driven Tests

Required when testing more than two cases.

Use descriptive `name` fields and `t.Run`.

### Level 3: HTTP Handler Tests

Required for HTTP handlers in s06 onward.

Use `httptest` instead of launching a real server unless the lesson explicitly teaches server lifecycle behavior.

```go
req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
w := httptest.NewRecorder()

handler.ServeHTTP(w, req)

if w.Code != http.StatusOK {
	t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
}
```

### Level 4: Benchmarks

Required for performance-critical code and performance-focused lessons.

```go
func BenchmarkSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// benchmarked work
	}
}
```

Run with:

```bash
go test -bench=. -benchmem
```

### Level 5: Integration Tests

Required where multiple components must work together.

Use helpers with `t.Helper()` and cleanup with `t.Cleanup`.

### Level 6: Fuzz Tests

Required or strongly recommended for parsing, decoding, validation, and input-boundary code.

### Level 7: API and gRPC Tests

Required for API contracts in s06 and related sections.

## Error Handling in Tests

Use `t.Fatalf` for setup failures or when continuing would be misleading.

Use `t.Errorf` when the test can continue and report more failures.

Always include enough context in failure messages to diagnose the issue.

## Test Helpers

Helpers must call `t.Helper()`.

```go
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}
```

## Coverage Goals

| Area | Target |
| --- | --- |
| General minimum where no stricter section target exists | 75% |
| Foundation-level lesson items with testable behavior | 80% |
| s06 Backend, APIs & Databases | 85% |
| s08 Quality & Testing | 95% |
| s09 Architecture & Security | 85% |
| s11 Opslane Flagship | 90% |

Coverage is a signal, not a substitute for meaningful assertions.

## Required Local Verification

For full PR readiness, use the same CI-equivalent bundle as `CODE-STANDARDS.md`:

```bash
go build ./...
go vet ./...
unformatted=$(gofmt -l .); test -z "$unformatted" || (echo "$unformatted" && exit 1)
go mod tidy
git diff --exit-code -- go.mod go.sum
go test ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
go tool cover -func coverage.out
go run ./scripts/validate_curriculum.go
```

On PowerShell, quote the coverage flag if needed:

```powershell
go test "-coverprofile=coverage.out" ./...
```

Do not commit generated `coverage.out` or `coverage.html` artifacts.

For benchmark changes:

```bash
go test -bench=. -benchmem -count=1 ./08-quality-test/01-quality-and-performance/testing/benchmarks/
```

## CI Expectations

The current CI workflow verifies:

- build
- vet
- gofmt check
- `go mod tidy` no-diff check
- tests
- race tests
- coverage report generation
- benchmark command
- curriculum validator

Staticcheck is recommended locally unless CI is updated to enforce it.

## Lesson Author Checklist

- [ ] Unit tests cover exported behavior.
- [ ] Table-driven tests are used for multiple cases.
- [ ] Edge cases are covered.
- [ ] Error paths are tested.
- [ ] Helpers use `t.Helper()`.
- [ ] Cleanup uses `t.Cleanup()` where appropriate.
- [ ] Concurrency tests avoid arbitrary sleeps.
- [ ] Race-sensitive changes pass `go test -race ./...`.
- [ ] Tests verify behavior rather than implementation trivia.
- [ ] Exercise starter code compiles.
