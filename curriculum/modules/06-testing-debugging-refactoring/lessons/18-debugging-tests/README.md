# Debugging tests

## Learning objective

Use `dlv test` to debug Go test functions, set breakpoints inside specific tests, inspect test fixtures, and identify race conditions during test execution.

## Why this matters

A failing test gives you a failure message and a line number. But the failure message often tells you *what* went wrong, not *why*. The test's internal state — the mock return values, loop iterations, intermediate assertions — is invisible from the failure output alone. `dlv test` opens the test's runtime and lets you step through the exact failing path with full visibility, turning a 30-minute guessing game into a 2-minute inspection.

## Mental model

`dlv test` works like `dlv debug` but for the test binary. The test runner is a Go program that discovers and runs test functions. You set breakpoints inside specific test functions or helper functions they call. When a test function runs and hits a breakpoint, you inspect state, step through assertions, and watch the test's data flow — exactly like debugging a regular program.

## Core idea

| Command | Effect |
|---|---|
| `dlv test` | Build and debug tests in the current package |
| `dlv test -- -test.run TestName` | Debug only a specific test |
| `dlv test -- -test.v` | Run tests with verbose output under the debugger |
| `b TestSomething` | Break at the start of a specific test function |
| `b helperFunc` | Break inside a helper the test calls |

Key difference from `dlv debug`: the test binary may run multiple test functions. If you set a breakpoint on `TestFoo`, only `TestFoo` stops — other tests run unimpeded.

Race condition debugging requires the `-race` flag:
```bash
dlv test -- -race ./...
```
Note: debugging with the race detector enabled is slower but lets you catch data races in the debugger.

## Under the hood

`dlv test` compiles the test binary (with debug symbols), then starts it as a child process. The test binary is a regular Go program that uses the `testing` package's `M.Run()` to discover and execute test functions. Delve intercepts execution at the same level as any other program. Breakpoints in test functions work identically to breakpoints in production code.

For `dlv test -- -test.run TestFoo`, the `-test.run` flag is passed to the test binary, which uses `regexp.MatchString` against test function names. Only matching tests execute — the others are skipped.

## How Go uses it

- **Debugging assertion failures**: Step through a table-driven test to see exactly which input produces the wrong output.
- **Inspecting mocks and stubs**: Verify that mock return values are what you expect at the point of use.
- **Debugging test helpers**: `t.Helper()` marks a function as a test helper. Delve respects this and can skip helpers in stack traces.
- **Race reproduction**: Run `dlv test -- -race` with a breakpoint before a suspected race to inspect goroutine state.
- **Benchmark debugging**: `dlv test -- -bench=. -benchtime=1x` runs a benchmark once under the debugger.

## Go example

```go
package main

import (
	"errors"
	"testing"
)

func TestProcessOrder(t *testing.T) {
	tests := []struct {
		name     string
		items    []string
		wantErr  bool
		wantCost int
	}{
		{"valid order", []string{"item1", "item2"}, false, 200},
		{"empty order", []string{}, true, 0},
		{"unknown item", []string{"nope"}, true, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cost, err := processOrder(tc.items)
			if (err != nil) != tc.wantErr {
				t.Errorf("processOrder() error = %v, wantErr %v", err, tc.wantErr)
			}
			if cost != tc.wantCost {
				t.Errorf("processOrder() cost = %d, want %d", cost, tc.wantCost)
			}
		})
	}
}

type Item struct {
	Name  string
	Price int
}

var catalog = map[string]int{
	"item1": 100,
	"item2": 100,
}

func processOrder(items []string) (int, error) {
	if len(items) == 0 {
		return 0, errors.New("empty order")
	}
	total := 0
	for _, name := range items {
		price, ok := catalog[name]
		if !ok {
			return 0, errors.New("unknown item: " + name)
		}
		total += price
	}
	return total, nil
}
```

Debugging session:
```bash
dlv test
(dlv) break TestProcessOrder
(dlv) continue
(dlv) break processOrder
(dlv) continue
(dlv) print items
(dlv) next; print total; next
(dlv) continue
```

To run only the "empty order" test case:
```bash
dlv test -- -test.run TestProcessOrder/empty_order
```

## Step-by-step execution

1. `dlv test` compiles and starts the test binary. It pauses at the test binary's entry point.
2. `break TestProcessOrder` — sets a breakpoint at the start of the test function.
3. `continue` — runs to the breakpoint. The test function's local `tests` slice is initialised.
4. `break processOrder` — sets a breakpoint inside the production code under test.
5. `continue` — runs to the first call of `processOrder` (for "valid order").
6. `print items` — shows `[]string{"item1", "item2"}` confirming the first test case.
7. `next` — steps through the loop, showing `total` accumulating.
8. `continue` — finishes this call and proceeds to the next test case.
9. `print items` — shows `[]string{}` (empty order case).
10. The test returns early with an error.

## Common mistakes

- **Forgetting `-test.run`**: Without it, *all* tests in the package run, hitting breakpoints multiple times. Use `-test.run TestName` to focus on one test.
- **Setting breakpoints in skipped tests**: If a test is skipped by `-test.run`, its breakpoints never fire. Confirm the test name matches the regex.
- **Breaking on `t.Run` subtests**: A breakpoint on `TestProcessOrder` stops at the enclosing function, not inside subtests. To break in a subtest, set a breakpoint on the inner test function passed to `t.Run` — but because it's a closure, use a line breakpoint inside it.
- **Using `dlv debug` instead of `dlv test`**: `dlv debug` runs the `main` package. Tests are in `*_test.go` files and only run with `dlv test` or `go test`.
- **Race detector masking bugs**: The race detector (`-race`) changes memory ordering and may make some bugs non-deterministic. It adds significant overhead. Debug without `-race` first, then add it to confirm.

## Debugging walkthrough

Suppose the "empty order" test fails but only intermittently. We suspect a race condition in `processOrder` perhaps through a shared data structure.

```bash
dlv test -- -race -test.run TestProcessOrder/empty_order
(dlv) break processOrder
(dlv) continue
(dlv) print items
[]
(dlv) next                  # len(items) == 0
(dlv) next                  # return 0, error
(dlv) continue              # test assertion runs
```

No race — but the test fails because `wantErr` is `true` and the function returns an error correctly... wait, that should pass. The intermittent failure must be elsewhere. We add breakpoints on the other test cases and watch. The race is in `catalog` being read in `processOrder` while another goroutine writes to it. We add a breakpoint in the writer goroutine to confirm.

## Production notes

- **CI with `dlv test`**: In CI, use `dlv test --headless --listen=:2345 --api-version=2 -- -test.run TestName` to start a headless debug server for remote debugging.
- **Race flag in staging**: Run tests with `-race` in CI but not in production binaries. The race detector adds 5–10x CPU and memory overhead.
- **Parallel tests**: `t.Parallel()` runs subtests concurrently. The debugger follows one goroutine at a time. Use `goroutines` to list all test goroutines and switch between them.
- **Test caching**: `go test` caches results. `dlv test` always recompiles. Run `go clean -testcache` if `go test` is acting stale.

## Performance implications

- `dlv test` adds breakpoint overhead to each test that hits a breakpoint. Tests without breakpoints run at full speed.
- The `-race` flag under `dlv test` is slower than without it due to race detector instrumentation (5–10x slower).
- Debugging individual tests with `-test.run` is faster than debugging the whole suite because fewer functions run.

## Practice task

Write a test `TestProcessOrder` that calls `processOrder` with three different inputs. Start `dlv test` with `-test.run TestProcessOrder`. Set a breakpoint on `processOrder`. Step through each call and inspect the `items` slice and `total` accumulator. Use `continue` between test cases to see how each case exercises different branches.

## Tests / verification

```bash
go test ./curriculum/modules/06-testing-debugging-refactoring/lessons/18-debugging-tests
```

(Note: `dlv test` is run interactively; the above command runs tests normally.)

## Review questions

1. What is the difference between `dlv debug` and `dlv test`?
2. How do you debug a single test function instead of all tests in a package?
3. Where does execution stop when you set a breakpoint on a `t.Run` subtest?
4. Why might you pass `-race` to `dlv test`?
5. If you set a breakpoint on a helper function used by multiple tests, how do you know which test triggered it?

## NEXT UP

Refactoring safely — using tests as a safety net for code changes.
