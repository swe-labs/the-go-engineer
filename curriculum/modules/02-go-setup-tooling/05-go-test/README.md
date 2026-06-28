# go test

## Learning objective

By the end of this lesson, you will understand what `go test` does, how to write and run a basic test, and how to interpret test output. You will run existing tests and observe pass/fail results.

## Why this matters

Testing is how you prove your code works. `go test` is the command that runs those proofs. Every professional Go project relies on `go test` to catch regressions, verify fixes, and ensure that changes do not break existing behavior.

## Mental model

A test is a function that calls your code and checks the result against an expected value. `go test` finds all test functions (functions starting with `Test` in files ending with `_test.go`), runs them, and reports which passed and which failed.

## Core idea

`go test` compiles the package and all `_test.go` files, runs every function matching `func TestXxx(t *testing.T)`, and reports PASS or FAIL for each one. A test fails when it calls `t.Errorf` or `t.Fatalf`.

## Under the hood

`go test` creates a separate test binary that links your package code with the test files. It runs this binary, which calls each `TestXxx` function. The `*testing.T` parameter provides methods like `Errorf` (report failure, continue) and `Fatalf` (report failure, stop immediately). `go test -v` prints the name and result of every test function.

## How Go uses it

The Go standard library has over 30,000 test functions. Every package in the standard library has tests, and the Go team requires all tests to pass before any change is merged. Common testing patterns:

- `go test .` — run tests in the current package
- `go test ./...` — run tests in the current package and all subpackages
- `go test -v .` — verbose output showing each test name
- `go test -run TestAdd .` — run only tests matching the pattern `TestAdd`

## Go example

```go
package main

import "fmt"

func Add(a, b int) int {
	return a + b
}

func main() {
	fmt.Println("2 + 3 =", Add(2, 3))
	fmt.Println("Testing with go test is next.")
}
```

```go
package main

import "testing"

func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5
	if got != want {
		t.Errorf("Add(2, 3) = %d; want %d", got, want)
	}
}
```

## Step-by-step execution

1. Run `go test .` from this lesson directory to run all tests.
2. Check the output — it should say `ok` or `PASS`.
3. Run `go test -v .` to see each test name and its result.
4. Introduce a deliberate bug in `Add` (change `+` to `-`) and run `go test .` again to see a failure.
5. Fix the bug and run `go test .` again to confirm the fix.

## Common mistakes

- **Putting tests in the wrong file**: Test files must end with `_test.go` (e.g., `main_test.go`, not `test_main.go`).
- **Forgetting the `*testing.T` parameter**: `func TestAdd()` without `(t *testing.T)` is not a test function — `go test` ignores it silently.
- **Using `t.Fatalf` when `t.Errorf` is better**: `Fatalf` stops the current test immediately, hiding any later failures. Use `Errorf` when you want to report all failures in one run.
- **Running tests from the wrong directory**: `go test .` runs tests for the package in the current directory. `go test ./...` runs all tests in the module.

## Debugging walkthrough

**Scenario**: `go test .` outputs `ok` but you know the code is wrong.

1. Check that your test file ends in `_test.go`. If it does not, `go test` ignores it.
2. Check that your test function is named `TestXxx` with a capital T. `testAdd` (lowercase t) is not a test — `go test` only runs exported functions (capital letter).
3. Run `go test -v .` to see which tests actually ran. If your test name does not appear, it was not found.

## Production notes

- CI pipelines always run `go test ./...` to verify the entire module.
- Use `go test -count=1 .` to bypass the test cache and force a fresh run.
- Use `go test -shuffle=on .` to randomize test order and detect hidden dependencies between tests.
- Integration tests use a build tag: `//go:build integration` and run with `go test -tags=integration .`.

## Performance implications

- `go test` caches results when the source and test files have not changed. Cached results print `(cached)` instead of rerunning.
- Use `-count=1` to bypass the cache when you need a fresh run (e.g., after switching branches).
- Tests run in parallel within a package by default. Use `t.Parallel()` to mark tests for parallel execution.

## Practice task

1. Run `go test -v .` from this lesson directory and observe the output.
2. Modify the `Add` function to return `a * b` instead of `a + b`, then run `go test .` and watch the tests fail.
3. Restore the correct implementation and add a new test: `TestAddLarge` that checks `Add(1000000, 2000000) == 3000000`.

## Tests / verification

```bash
go test -v .
```

Expected output (passing):
```
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
=== RUN   TestAddNegative
--- PASS: TestAddNegative (0.00s)
=== RUN   TestAddZero
--- PASS: TestAddZero (0.00s)
PASS
ok  	02-go-setup-tooling/05-go-test	0.xxx
```

Also run:
```bash
go run .
```
Expected output:
```
2 + 3 = 5
Testing with go test is next.
```

## Review questions

1. What must a test file name end with?
2. What is the difference between `t.Errorf` and `t.Fatalf`?
3. What does `go test -v` show that `go test` does not?
4. How do you run all tests in a module and its subpackages?
5. Why might `go test` print `(cached)` instead of rerunning tests?

## NEXT UP

[gofmt](../06-gofmt/README.md) — format your Go code automatically.
