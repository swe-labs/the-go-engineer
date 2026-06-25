# Unit testing

## Learning objective

Write unit tests using Go's `testing` package, organize tests in `_test.go` files, and differentiate between `t.Error` and `t.Fatal` to report failures effectively.

## Why this matters

Unit tests are the most granular and most numerous tests in any codebase. They execute in milliseconds, run in isolation, and pinpoint failures to a single function. Every Go engineer writes unit tests daily. Mastering the `testing` package — its conventions, control flow, and reporting methods — is the foundation of all other testing work.

## Mental model

A unit test is a function that calls another function and checks the result. Think of it as a contract verifier: the production code promises "given input X, I return output Y," and the test checks that promise. If the contract is broken, the test reports it. Unlike a production function, a test function receives a `*testing.T` value that is the test's communication channel with the world — use it to log, fail, or skip.

```
Production:   func Add(a, b int) int
Test:         func TestAdd(t *testing.T)
Contract:     Add(2, 3) == 5
```

## Core idea

The `testing` package provides the framework for automated tests. Every test function must follow the pattern:

```go
func TestXxx(t *testing.T) {
    // setup
    // call the function under test
    // compare result to expected value
    // report failure with t.Error or t.Fatal
}
```

Key conventions:

| Convention | Rule |
|---|---|
| File name | `*_test.go` |
| Function name | `TestXxx` where `Xxx` starts with a capital letter |
| Parameter | `t *testing.T` |
| Import | `"testing"` |
| Reporting | `t.Error` / `t.Errorf` (non-fatal), `t.Fatal` / `t.Fatalf` (fatal) |

- `t.Error` reports a failure but continues the test. Use it when you want to check multiple conditions and report all failures.
- `t.Fatal` reports a failure and stops the test immediately. Use it when continuing would panic or produce misleading errors.

## Under the hood

`*testing.T` is a struct that tracks test state. When you call `t.Error` or `t.Fatal`, it marks the test as failed internally. The difference is that `t.Error` sets a flag and returns, while `t.Fatal` sets the flag and calls `runtime.Goexit()` to terminate the current goroutine. This means `t.Fatal` runs deferred functions; a bare `panic` does not.

The test binary is linked with the package under test in the same package. If the test file declares `package main`, it has access to all exported and unexported symbols in the `main` package. For non-main packages, tests typically use `package foo_test` (external test package) to enforce black-box testing, or `package foo` for white-box access to internals.

## How Go uses it

Every Go project uses the `testing` package the same way — there is no alternative test framework baked in. Third-party libraries like `testify` add assertions and mocking, but the runner is always `go test`. The convention is:

- One test file per source file: `foo.go` → `foo_test.go`.
- One test function per behavior, not per production function.
- Test functions should be independent and order-independent.
- Use `go test -v` for verbose output, `go test -run TestXxx` to run a specific test.

## Go example

```go
package main

import (
	"fmt"
	"strings"
)

// Greet returns a greeting for the given name.
func Greet(name string) string {
	return "Hello, " + strings.TrimSpace(name) + "!"
}

// Sum returns the sum of a slice of integers.
func Sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func main() {
	fmt.Println(Greet("Alice"))
	fmt.Println(Sum([]int{1, 2, 3, 4}))
}
```

```go
// main_test.go
package main

import "testing"

func TestGreet(t *testing.T) {
	got := Greet("Alice")
	want := "Hello, Alice!"
	if got != want {
		t.Errorf("Greet(\"Alice\") = %q; want %q", got, want)
	}
}

func TestGreetEmpty(t *testing.T) {
	got := Greet("")
	want := "Hello, !"
	if got != want {
		t.Errorf("Greet(\"\") = %q; want %q", got, want)
	}
}

func TestSum(t *testing.T) {
	got := Sum([]int{1, 2, 3, 4})
	want := 10
	if got != want {
		t.Errorf("Sum({1,2,3,4}) = %d; want %d", got, want)
	}
}

func TestSumEmpty(t *testing.T) {
	got := Sum([]int{})
	want := 0
	if got != want {
		t.Errorf("Sum({}) = %d; want %d", got, want)
	}
}
```

## Step-by-step execution

Running `go test -v`:

1. `go test` compiles `main.go` and `main_test.go` together.
2. Test discovery: finds `TestGreet`, `TestGreetEmpty`, `TestSum`, `TestSumEmpty`.
3. Runs `TestGreet`: calls `Greet("Alice")`, gets `"Hello, Alice!"`, compares to `want`. Match. Pass.
4. Runs `TestGreetEmpty`: calls `Greet("")`, gets `"Hello, !"`. Match. Pass.
5. Runs `TestSum`: calls `Sum({1,2,3,4})`, gets `10`. Match. Pass.
6. Runs `TestSumEmpty`: calls `Sum({})`, gets `0`. Match. Pass.
7. All pass. Output: `ok` and exit code 0.

If `Greet` had a bug and returned `"Hello Alice!"` (missing comma), the test would print: `Greet("Alice") = "Hello Alice!"; want "Hello, Alice!"`.

## Common mistakes

- **Using `t.Fatal` when `t.Error` would do.** `t.Fatal` aborts the test immediately. If the first assertion fails but later assertions would also fail, you lose information. Prefer `t.Error` unless continuing would cause a panic (e.g., dereferencing a nil pointer from the result).

- **Forgetting to run the test.** Writing a test file and never executing it is common. Always run `go test` after writing tests.

- **Using `t.Error` after a nil check.** If a function returns a value and an error, and the error is non-nil, using the value will panic. Check the error first with `t.Fatal`, then assert on the value.

- **Testing the wrong package.** `go test` runs tests for the current package. To test all packages: `go test ./...`.

- **Neglecting edge cases.** Test zero values, empty slices, negative numbers, and boundary conditions — not just the happy path.

## Debugging walkthrough

A unit test is failing:

```go
func TestGreet(t *testing.T) {
	got := Greet(" Alice ") // note: input has spaces
	want := "Hello, Alice!"
	if got != want {
		t.Errorf("Greet() = %q; want %q", got, want)
	}
}
```

**Symptom**: `Greet(" Alice ")` returns `"Hello,  Alice !"` with extra spaces because `TrimSpace` is not called.

**Investigation**: Run `go test -v` and observe the failure output. The `%q` formatting shows the exact characters including spaces.

**Root cause**: The `Greet` function wraps the name in the greeting before trimming.

**Fix**: Apply `strings.TrimSpace` before interpolation: `return "Hello, " + strings.TrimSpace(name) + "!"`. Or, if the function is correct but the test expectation is wrong, fix the test.

## Production notes

- **Internal vs external tests.** Use `package foo` (internal) for testing unexported helpers. Use `package foo_test` (external) to enforce that tests use only the public API.
- **Test helpers.** If multiple tests share setup logic, extract it into a helper function that takes `*testing.T` as a parameter.
- **`t.Helper()`.** Mark helper functions with `t.Helper()` so test failure lines point to the caller, not the helper internals.
- **`go test -count=1`.** Disables caching. Useful during development when the same test should run fresh each time.
- **`go test -shuffle=on`.** Randomizes test order to surface hidden dependencies between tests.

## Performance implications

- Unit tests should be fast. A slow unit test discourages frequent execution.
- Tests that write to disk, make network calls, or sleep are not unit tests — they are integration tests.
- `t.Error` and `t.Errorf` have minimal overhead. They format the message only on failure, so the happy path is cheap.
- The `-v` flag adds IO overhead. Use it only when debugging.

## Practice task

Write a function `SplitHostPort(addr string) (host string, port int, err error)` that parses `"host:port"` strings. The port must be a number between 1 and 65535. Write unit tests for:
- Valid input (`"localhost:8080"`)
- Missing port (`"localhost"`)
- Invalid port (`"localhost:abc"`)
- Port out of range (`"localhost:0"`, `"localhost:70000"`)

Run `go test -v` and verify all cases pass.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/02-unit-testing
go test -v ./curriculum/modules/06-testing-debugging-refactoring/lessons/02-unit-testing
```

## Review questions

1. What is the difference between `t.Error` and `t.Fatal`? When would you use each?
2. How does `go test` discover which functions to run?
3. What happens if a test function has the signature `func TestAdd(t *testing.T)` but the file is named `add_test.go` versus `add.go`?
4. Why might you choose `package main` (internal test) over `package main_test` (external test)?
5. What does `t.Helper()` do and why is it useful?

## NEXT UP

Table-driven tests — structuring multiple test cases cleanly with slices of test data.
