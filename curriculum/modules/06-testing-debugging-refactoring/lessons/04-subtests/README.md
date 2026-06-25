# Subtests

## Learning objective

Use `t.Run` to create named subtests, run subtests in parallel with `t.Parallel`, perform setup and teardown per subtest, and selectively execute subsets with `-run`.

## Why this matters

As test suites grow, running every test for every change becomes slow. Subtests let you group related tests, run them independently, and parallelize expensive ones. The `-run` flag lets you target a single subtest during development without waiting for the entire suite. This is essential for any codebase with more than a handful of tests.

## Mental model

Think of `t.Run(name, fn)` as creating a mini test within a test. Each subtest has its own `*testing.T`, its own pass/fail state, and its own name in the output tree. The parent test is a container that waits for all subtests to finish. Subtests can be run selectively: `go test -run Parent/Child` runs only one branch of the tree.

```
TestAPI              ← parent test
├── TestAPI/GET      ← subtest
├── TestAPI/POST     ← subtest
└── TestAPI/DELETE   ← subtest
```

## Core idea

`t.Run` creates a subtest:

```go
func TestParent(t *testing.T) {
    t.Run("child", func(t *testing.T) {
        // subtest body
    })
}
```

Key behaviors:

| Feature | How |
|---|---|
| **Named subtests** | `t.Run("name", fn)` — name appears in output |
| **Parallel subtests** | `t.Parallel()` inside the subtest block |
| **Selective run** | `go test -run 'TestParent/child$'` |
| **Nested subtests** | `t.Run` can be called inside another subtest |
| **Shared setup** | Code before `t.Run` runs once for all subtests |
| **Per-subtest setup** | Call setup inside each subtest's block |

Parallel subtests are especially valuable. When a test calls `t.Parallel()`, it pauses until the parent test returns from `t.Run`, then all parallel subtests run concurrently:

```go
func TestParallel(t *testing.T) {
    for _, tc := range cases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()
            // tc is safe here because of shadow copy
        })
    }
}
```

## Under the hood

When `t.Run(name, fn)` is called:

1. Go creates a new `*testing.T` with a name of `Parent/Child` (using `/` as separator).
2. The subtest runs in a new goroutine.
3. The parent `t` waits for the subtest goroutine to complete before returning from `t.Run`.
4. If the subtest calls `t.Parallel()`, it blocks inside `t.Parallel()` until the parent finishes all top-level subtests, then resumes.
5. Failure in a subtest does not affect sibling subtests — only the parent test is marked as failed.

The `-run` flag accepts a regex pattern. Go matches it against the full test name (parent + `/` + child). Only matching tests execute.

## How Go uses it

Subtests are used throughout Go's standard library:

- `go test -run 'TestSplit/simple'` in `strings` tests.
- `go test -run 'TestFormat/'` to run all formatting subtests.
- `go test -run 'TestEncode/(JSON|XML)'` to run only JSON and XML encoding subtests.

The pattern `-run ^TestName$` runs only the top-level test with no subtests. `-run TestName/sub` runs all subtests matching `sub`.

## Go example

```go
package main

import (
	"fmt"
	"time"
)

// FetchData simulates a data fetch that takes variable time.
func FetchData(source string) string {
	time.Sleep(10 * time.Millisecond) // simulate work
	return "data from " + source
}

func main() {
	fmt.Println(FetchData("cache"))
	fmt.Println(FetchData("db"))
}
```

```go
// main_test.go
package main

import (
	"testing"
	"time"
)

func TestFetchData(t *testing.T) {
	sources := []struct {
		name string
		src  string
	}{
		{name: "cache", src: "cache"},
		{name: "database", src: "db"},
		{name: "api", src: "api"},
	}
	for _, s := range sources {
		s := s
		t.Run(s.name, func(t *testing.T) {
			t.Parallel()
			start := time.Now()
			got := FetchData(s.src)
			elapsed := time.Since(start)
			if elapsed > 100*time.Millisecond {
				t.Errorf("FetchData(%q) took %v; want <100ms", s.src, elapsed)
			}
			if got == "" {
				t.Errorf("FetchData(%q) returned empty", s.src)
			}
		})
	}
}
```

## Step-by-step execution

Running `go test -v -run TestFetchData`:

1. `TestFetchData` starts. The loop creates 3 subtests: "cache", "database", "api".
2. Each subtest calls `t.Parallel()`, which blocks until all subtests in the parent are registered.
3. The parent finishes its loop and returns. Now the parallel barrier lifts.
4. All 3 subtests run concurrently, each calling `FetchData` with their source.
5. Each subtest verifies the result is non-empty and the duration is under 100ms.
6. All pass. Output shows `--- PASS: TestFetchData` with all subtests listed.

If "api" took 200ms, only that subtest fails. "cache" and "database" still pass.

## Common mistakes

- **Loop variable capture.** Without `s := s` (or Go 1.22+), parallel subtests share the loop variable and get the last iteration's value. Always shadow the variable when using `t.Parallel()`.

- **Mixing sequential and parallel subtests in the same parent.** Once a subtest calls `t.Parallel()`, all remaining subtests in that parent also wait for the parallel barrier. To run sequential then parallel, use separate test functions.

- **Misunderstanding `-run` patterns.** `-run TestFoo` matches any test containing "TestFoo" anywhere in the name. Use `^TestFoo$` for exact match.

- **Using `t.Fatal` when `t.Error` would do in a subtest.** `t.Fatal` stops only that subtest, so the impact is limited, but `t.Error` still gives more information.

## Debugging walkthrough

A parallel test is flaky:

```go
func TestCounter(t *testing.T) {
    counter := 0
    for i := 1; i <= 3; i++ {
        i := i
        t.Run(fmt.Sprintf("inc_%d", i), func(t *testing.T) {
            t.Parallel()
            counter += i
        })
    }
    t.Log("counter =", counter)
}
```

**Symptom**: `counter` is 0 or 3 instead of 6. The test is nondeterministic.

**Investigation**: `counter` is a shared variable accessed from multiple goroutines without synchronization. `t.Parallel()` causes subtests to run in separate goroutines, and the parent logs `counter` before the subtests increment it.

**Root cause**: Race condition. The parent returns from `t.Run` only after each subtest completes, but with `t.Parallel()`, the subtest does not run until all subtests are registered and the parent returns from the top-level function. So `counter` is read before writes complete.

**Fix**: Remove `t.Parallel()` or add proper synchronization. For parallel tests, never mutate shared state.

## Production notes

- **Use parallel subtests for slow tests.** If a test makes HTTP calls, waits on timers, or does heavy computation, parallel execution speeds up the suite.
- **Avoid `t.Parallel()` in tests that modify shared state.** Each subtest should operate on its own data.
- **`-run` is a development tool.** In CI, run the full suite. Use `-run` locally to focus on failing tests during debugging.
- **Test output order with parallelism.** Since subtests run concurrently, their output interleaves. Use `t.Log` to add context to failures.
- **Subtest name uniqueness.** Duplicate subtest names within a parent are allowed but make selective execution ambiguous. Keep names unique.

## Performance implications

- Parallel subtests can reduce wall-clock time for a suite from minutes to seconds, at the cost of higher CPU and memory usage.
- The concurrency is limited by `GOMAXPROCS`. For CPU-bound tests, more parallelism may not help.
- Each `t.Run` goroutine has stack overhead (~2KB min). For thousands of subtests, this adds up but is rarely a problem.
- IO-bound parallel subtests benefit significantly from concurrency — database tests, HTTP tests, file system tests.

## Practice task

Write a function `Process(id int) (int, error)` that returns `id * 2` if `id > 0`, or an error if `id <= 0`. Write a test with 5 subtests:
- "positive", "zero", "negative" — all sequential
- Add a parallel subtest group for "large_1", "large_2", "large_3" that sleep 10ms each

Run `go test -v -run TestProcess` and observe the timing difference with and without `t.Parallel()`.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/04-subtests
go test -v ./curriculum/modules/06-testing-debugging-refactoring/lessons/04-subtests
```

## Review questions

1. What does `t.Run("name", fn)` do and how does it affect test output?
2. How do you selectively run a single subtest named "add" inside "TestMath"?
3. Why is the shadow copy (`tt := tt`) necessary when using `t.Parallel()` inside a range loop on Go < 1.22?
4. What happens if you call `t.Parallel()` inside a subtest — does the parent wait?
5. Can subtests be nested? If so, how would you run only the innermost subtest with `-run`?

## NEXT UP

`t.Cleanup` — registering cleanup functions that run when a test completes.
