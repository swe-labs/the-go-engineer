# Why testing exists

## Learning objective

Explain the purpose of software testing, articulate how tests prevent regressions and increase confidence, and apply Go's testing philosophy by writing a simple test.

## Why this matters

Without tests, every change to a codebase is a leap of faith. You modify one function and hope nothing else breaks. In a codebase of any size, that hope is misplaced. Tests convert hope into evidence. They are the mechanism by which a team of engineers can ship changes rapidly without fear. Professional Go engineers treat tests as a first-class deliverable, not an afterthought.

## Mental model

Imagine a safety net below a tightrope walker. The net does not help the walker cross, but it makes the fall safe. Tests are the same: they don't write the code, but they make mistakes safe. When you refactor, the test suite catches what you accidentally broke. When a new teammate changes unfamiliar code, the tests tell them instantly whether their change is correct. The test suite is an executable specification — it documents what the code is supposed to do in a form the machine can verify.

## Core idea

Testing exists to answer one question: *does my code work as intended?* This breaks down into several concrete motivations:

| Motivation | Why it matters |
|---|---|
| **Confidence** | Deploy with proof, not hope. |
| **Regression prevention** | Catch bugs reintroduced by new code. |
| **Documentation** | Tests show how code is *supposed* to behave. |
| **Design feedback** | Hard-to-test code is usually poorly designed. |
| **Refactoring safety** | Change internal structure without changing behavior. |

The **test pyramid** describes the ideal distribution of test types:

```
    /\
   /  \        E2E (few)
  /    \
 / unit \      Unit tests (many)
/________\
```

Unit tests form the base — they are fast, isolated, and numerous. Integration tests sit above. End-to-end tests sit at the top — they are slow and brittle, so you write fewer of them.

## Under the hood

When `go test` runs, the Go toolchain:

1. Scans each package for files matching `*_test.go`.
2. Compiles a temporary binary that links the package under test with the testing framework.
3. Discovers every function with the signature `func TestXxx(t *testing.T)`.
4. Calls each `TestXxx` function in sequence, passing a `*testing.T` value.
5. Records pass/fail for each test and reports the summary.

If any test calls `t.Error` or `t.Fatal`, that test fails. If any test in a package fails, the package overall fails. The exit code is non-zero when tests fail, which is how CI systems detect breakage.

## How Go uses it

Go's testing philosophy is deliberately minimal:

- **No assertions library in the stdlib** — you write `if got != want { t.Errorf(...) }` directly.
- **Test files live alongside source** — `foo.go` is tested by `foo_test.go` in the same package.
- **`go test` is the single command** — no test runner configuration needed.
- **`go test -cover`** measures code coverage.
- **`go test -bench`** runs benchmarks using the same `testing` package.
- **`go vet`** and **`go test -race`** catch deeper issues.

This simplicity means tests are just Go code. There is no special DSL, no magic setup — just functions that call your code and check results.

## Go example

```go
package main

import "fmt"

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

// IsEven returns true if n is even.
func IsEven(n int) bool {
	return n%2 == 0
}

func main() {
	fmt.Println("Add(2, 3) =", Add(2, 3))
	fmt.Println("IsEven(4) =", IsEven(4))
	fmt.Println("IsEven(5) =", IsEven(5))
}
```

## Step-by-step execution

When you run `go test` on this package:

1. The compiler builds a test binary containing `main_test.go` (not yet written — we will) and `main.go` together.
2. The test runner finds `TestAdd` in the test file and calls it.
3. Inside `TestAdd`, `Add(2, 3)` executes, returns `5`.
4. The test compares `got` (5) to `want` (5). They match, so no failure is reported.
5. The test runner moves to `TestIsEven`.
6. `IsEven(4)` returns `true`, matching the want value. Pass.
7. All tests pass. `go test` prints `ok` and exits with code 0.

If `IsEven(5)` had been expected to return `true`, the test would fail and print a descriptive message.

## Common mistakes

- **Not testing failure cases.** A function might work for normal inputs but panic on edge cases. Test empty inputs, nil slices, zero values, and boundary conditions.

- **Testing implementation, not behavior.** Tests that check internal state rather than observable results break when the implementation changes, even if the behavior is correct.

- **Writing tests that are too coupled to production code structure.** Refactoring should not require rewriting tests. Test the public API, not private details.

- **Ignoring test output.** A passing test suite is not proof of correctness — it is proof that the code behaves as *the tests* specify. If the tests are wrong or incomplete, the suite is misleading.

- **Manual testing instead of automated.** Clicking through a UI or running ad-hoc commands is not reproducible. Write a test once, run it forever.

## Debugging walkthrough

Consider a function that is supposed to compute a discount:

```go
func Discount(price, rate int) int {
	return price * rate / 100
}
```

**Symptom**: `Discount(200, 10)` returns `20` (correct), but `Discount(199, 10)` returns `19` — seems correct, but the decimal was truncated. The caller expected cent-level precision.

**Investigation**: The function returns `int`, so fractional cents are lost. Add a test:

```go
func TestDiscount(t *testing.T) {
	got := Discount(199, 10)
	want := 19.9 // but want is int — compile error
}
```

**Root cause**: The return type `int` cannot represent fractional values. The design trades precision for speed.

**Fix**: Change the signature to return `float64` or use a `cents` int64 pattern (store 1990 cents, compute 199 cents discount).

## Production notes

- **CI gates.** Every PR must pass the full test suite. Green tests are a prerequisite for merge. Red builds block the pipeline.
- **Coverage targets.** Many teams enforce a minimum coverage threshold (e.g., 80%). Use `go test -coverprofile=coverage.out` to measure.
- **Flaky tests.** Tests that pass sometimes and fail unpredictably destroy trust. Track them down immediately or disable them.
- **Test readability.** Tests are read more often than they are written. Use descriptive names, clear assertions, and consistent formatting.
- **Test organization.** Group tests by behavior, not by function. Use `t.Run` subtests for related scenarios.

## Performance implications

- Test execution speed matters. Slow tests discourage running them. Aim for unit tests that complete in milliseconds.
- `go test` compiles each package afresh. Use `go test -count=1` to disable result caching during development.
- The `testing` package is designed to have zero allocation overhead in its hot path. `t.Error` and `t.Fatal` are efficient — they do not panic or allocate unless a failure occurs.

## Practice task

Write a test for the `Add` and `IsEven` functions in `main.go`. Then add a function `Divide(a, b int) (int, error)` that returns an error when `b == 0`. Write tests for the success case and the error case. Run `go test -v` to see your tests execute.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/01-why-testing-exists
go test -v ./curriculum/modules/06-testing-debugging-refactoring/lessons/01-why-testing-exists
```

## Review questions

1. What is the test pyramid and why are unit tests at the base?
2. When you run `go test` and all tests pass, what does that guarantee about the code?
3. What is the difference between `t.Error` and `t.Fatal`?
4. Why does the Go standard library not include assertion functions like `assert.Equal`?
5. A developer says "I tested it manually and it works." What arguments would you give for writing an automated test anyway?

## NEXT UP

Unit testing — writing `TestXxx` functions with `testing.T` and organizing test files.
