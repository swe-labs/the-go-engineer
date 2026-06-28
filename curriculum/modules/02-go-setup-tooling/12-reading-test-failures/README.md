# Reading Test Failures

## Learning objective

By the end of this lesson, you will be able to run tests with `go test -v`, read the output to identify which test failed and why, and fix both the test and the code under test.

## Why this matters

Test output is the most detailed feedback you will get about your code's correctness. A failing test tells you exactly which behavior is broken, what value was expected, and what value was produced. Reading test output correctly turns test failures into a precise debugging guide.

## Mental model

A test function calls your code, checks the result against an expected value, and reports a failure if they do not match. `go test -v` prints the name of every test function and whether it passed or failed, along with any failure messages.

## Core idea

Test output has three parts: the test name, the PASS or FAIL status, and the failure message (if any). The failure message includes the file and line number of the failing assertion, the expected value, and the actual value.

## Under the hood

`go test` compiles a test binary that links your package code with the test files. Each `TestXxx(t *testing.T)` function is called. When `t.Errorf` is called, it records the failure with the file and line using runtime.Caller. After all tests run, the test binary prints the summary. `go test -v` prints each test's progress in real time.

## How Go uses it

Test output format:

```
=== RUN   TestName
    file_test.go:10: expected message
--- PASS: TestName (0.00s)
=== RUN   TestFailing
    file_test.go:20: Add(2, 2) = 5; want 4
--- FAIL: TestFailing (0.00s)
FAIL
```

The message `Add(2, 2) = 5; want 4` tells you:
- The input: `Add(2, 2)`
- The actual value: `5`
- The expected value: `4`

## Go example

```go
package main

import "fmt"

func Add(a, b int) int {
	return a + b
}

func main() {
	fmt.Println(Add(2, 2))
}
```

```go
package main

import "testing"

func TestAdd(t *testing.T) {
	got := Add(2, 2)
	want := 4
	if got != want {
		t.Errorf("Add(2, 2) = %d; want %d", got, want)
	}
}
```

## Step-by-step execution

1. Run `go test -v .` from this lesson directory — tests should pass.
2. Modify `Add` to return `a * b` instead of `a + b` (introduce a bug).
3. Run `go test -v .` again and observe the failure message.
4. Read the failure message: note the test name, file:line, actual vs expected.
5. Fix the bug and confirm tests pass.

## Common mistakes

- **Running `go test` without `-v`**: Without `-v`, you only see `FAIL` or `ok`. You miss the individual test results and failure messages.
- **Only reading the expected value and ignoring the actual value**: Both are important. The actual value tells you what your code *did* produce; the expected tells you what it *should* produce.
- **Fixing the test instead of fixing the code**: If a test fails, the code is likely wrong. Only change the test if the expected value in the test is incorrect.
- **Not running tests after fixing**: Always run `go test .` again after making a change to verify the fix.

## Debugging walkthrough

**Scenario**: `go test -v` shows `TestAdd: got 5, want 4`.

1. The test calls `Add(2, 2)` and expects `4` but got `5`.
2. Check the `Add` function — it might be doing `a * b` instead of `a + b`, or returning `a + b + 1`.
3. Fix the `Add` function implementation.
4. Run `go test -v .` again to confirm the fix.

**Scenario**: `go test -v` shows `TestAddZero: got 0, want 0` — this is a PASS.

1. No failure here. The message shows the values for transparency.
2. If you see a PASS, the test assertions all passed.

## Production notes

- CI pipelines run `go test -v ./...` and archive the output as a build artifact.
- Focus on the first test failure — later failures may be cascading from the same root cause.
- Use `go test -run TestName .` to run a single failing test while debugging.

## Performance implications

- `go test -v` adds minimal overhead — the output is produced as tests run.
- Tests that call `t.Errorf` instead of `t.Fatalf` continue running after a failure, collecting more information in one run.

## Practice task

1. Run `go test -v .` from this lesson directory — confirm all tests pass.
2. Change `Add` to return `a - b` and run tests again — observe the failure.
3. Read the failure message and identify the test name, line number, actual vs expected.
4. Restore the correct implementation and add a new test case.
5. Run tests again to confirm all pass.

## Tests / verification

```bash
go test -v .
```

Expected output:
```
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
=== RUN   TestAddZero
--- PASS: TestAddZero (0.00s)
PASS
ok  	02-go-setup-tooling/12-reading-test-failures	0.xxx
```

```bash
go run .
```
Expected output: `4`

## Review questions

1. What information does `go test -v` show that `go test` does not?
2. In the failure message `Add(2, 2) = 5; want 4`, which value came from the code?
3. What is the difference between `t.Errorf` and `t.Fatalf` in test output?
4. Why should you fix the first test failure before addressing later ones?
5. What does `PASS` mean in test output?

## NEXT UP

[Go module root basics](../13-go-module-root-basics/README.md) — understand go.mod, module paths, and how Go organizes packages.
