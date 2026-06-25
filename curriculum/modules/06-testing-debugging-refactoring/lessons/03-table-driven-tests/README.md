# Table-driven tests

## Learning objective

Write table-driven tests using anonymous struct slices, iterate over test cases, name them clearly, and use `t.Run` for subtest isolation.

## Why this matters

Repetitive test code is brittle and hard to maintain. When you copy-paste a test block for each input, you increase the chance of errors and make it painful to add new cases. Table-driven tests solve this by collecting all test cases in a single data structure and iterating over them. This is the idiomatic Go testing pattern, used throughout the standard library and every production Go codebase.

## Mental model

Think of a table-driven test as a spreadsheet. Each row is one test case with columns: name, input, expected output. The test iterates rows — one row, one assertion. Adding a new case is a new row, not a new function. The structure is always:

```
table → iterate rows → run test → report failure
```

The data structure describes *what* to test; the loop describes *how* to test it.

## Core idea

A table-driven test uses a slice of anonymous structs to define test cases:

```go
func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{name: "positive", a: 2, b: 3, want: 5},
		{name: "negative", a: -1, b: 1, want: 0},
		{name: "zero",    a: 0, b: 0, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
```

Benefits of this pattern:

| Benefit | Why |
|---|---|
| **Additive** | New case = new struct literal row |
| **Readable** | All inputs and expectations in one place |
| **Isolated** | `t.Run` gives each case its own subtest |
| **Self-documenting** | Each case has a descriptive name |
| **Parallelizable** | Subtests can run in parallel |

## Under the hood

When `t.Run(name, fn)` is called, Go creates a new `*testing.T` value derived from the parent. The subtest runs in a separate goroutine. The parent test does not complete until all subtests finish. If a subtest calls `t.Fatal`, only that subtest stops — the others continue.

The subtest name is appended to the parent name with a `/` separator. Running `go test -run TestAdd/positive` runs only the subtest named `positive` inside `TestAdd`.

## How Go uses it

The Go standard library uses table-driven tests extensively. For example, `fmt` tests use a slice of struct cases for every formatting verb. `strings` tests use tables for `Split`, `Replace`, `Trim`, and so on. The `go test` tool supports the `-run` flag with pattern matching that works with subtest names:

```bash
go test -run 'TestAdd/positive'    # runs only the "positive" subtest
go test -run 'TestAdd/(positive|negative)'  # regex pattern
```

## Go example

```go
package main

import "fmt"

// Max returns the larger of two integers.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	fmt.Println("Max(3, 7) =", Max(3, 7))
	fmt.Println("Max(10, 2) =", Max(10, 2))
	fmt.Println("Max(5, 5) =", Max(5, 5))
}
```

```go
// main_test.go
package main

import "testing"

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{name: "first larger",  a: 10, b: 2, want: 10},
		{name: "second larger", a: 3,  b: 7, want: 7},
		{name: "equal",         a: 5,  b: 5, want: 5},
		{name: "negative",      a: -3, b: -1, want: -1},
		{name: "zero",          a: 0,  b: 0, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Max(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Max(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
```

## Step-by-step execution

Running `go test -v`:

1. `TestMax` is discovered and called.
2. The test slice `tests` is created with 5 cases.
3. Loop iteration 1: `t.Run("first larger", ...)` creates a subtest. Inside, `Max(10, 2)` returns `10`. Match. Subtest passes.
4. Loop iteration 2: `t.Run("second larger", ...)`. `Max(3, 7)` returns `7`. Match.
5. Iteration 3: "equal", `Max(5, 5)` returns `5`. Match.
6. Iteration 4: "negative", `Max(-3, -1)` returns `-1`. Match.
7. Iteration 5: "zero", `Max(0, 0)` returns `0`. Match.
8. All 5 subtests pass. Output shows `--- PASS: TestMax` for each.

If the "negative" case had a bug (e.g., `want: -3`), the output would show:
```
--- FAIL: TestMax (0.00s)
    --- FAIL: TestMax/negative (0.00s)
        main_test.go:18: Max(-3, -1) = -1; want -3
```

## Common mistakes

- **Not naming test cases.** Without a name field, failures show as `TestMax/#01`, which is not descriptive. Always include a `name` field.

- **Capturing loop variables incorrectly.** In Go versions before 1.22, the loop variable `tt` was reused across iterations. Use `tt := tt` inside the loop or upgrade to Go 1.22+ which fixed this. With `t.Run`, the subtest closure captures `tt` by reference — if running parallel tests, you must shadow the variable.

- **Mixing table-driven and ad-hoc patterns.** Do not put some cases in a table and write separate test functions for others. Use tables consistently.

- **Testing too many dimensions in one table.** If different cases require different setup or different assertions, split into multiple test functions or use separate tables for distinct behaviors.

- **Forgetting edge cases.** Table-driven tests make it easy to add many cases. Take advantage: include boundary values, nil inputs, empty slices, and overflow conditions.

## Debugging walkthrough

A table-driven test is failing:

```go
func TestDivide(t *testing.T) {
	tests := []struct {
		name   string
		a, b   int
		want   int
		wantOK bool
	}{
		{name: "simple", a: 10, b: 2, want: 5, wantOK: true},
		{name: "by zero", a: 5, b: 0, want: 0, wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Divide(tt.a, tt.b)
			if ok != tt.wantOK {
				t.Errorf("Divide(%d,%d) ok=%v; want %v", tt.a, tt.b, ok, tt.wantOK)
			}
			if ok && got != tt.want {
				t.Errorf("Divide(%d,%d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
```

**Symptom**: Test "by zero" fails because `Divide(5, 0)` panics instead of returning an error.

**Investigation**: The panic appears in the test output as a stack trace. The `wantOK: false` case expects a clean error return, but the production code panics at `a / b` when `b == 0`.

**Root cause**: The `Divide` function does not check for zero divisor.

**Fix**: Add a zero check: `if b == 0 { return 0, false }`.

## Production notes

- **Table-driven tests are the default.** When reviewing Go code, if a test has more than 2 similar-looking assertions, a table-driven approach is expected.
- **Named return values in test structs.** Use clear field names like `want`, `wantErr`, `wantCount`. Avoid cryptic abbreviations.
- **Large tables.** If a table has 20+ cases, consider splitting by category or moving the test data to a separate file (e.g., `testdata/cases.json`).
- **Subtest naming convention.** Use lowercase, hyphen-separated names: `"empty-input"`, `"nil-slice"`, `"max-boundary"`.

## Performance implications

- Table-driven tests are not inherently faster than individual test functions. The benefit is maintainability, not speed.
- Each `t.Run` creates a new `*testing.T` and goroutine. Overhead is negligible for tens of cases. For thousands of cases, consider whether the loop logic itself is efficient.
- Subtests marked with `t.Parallel()` can run concurrently, potentially reducing wall-clock time for slow tests (e.g., those that wait on IO).

## Practice task

Write a function `ParseBool(s string) (bool, error)` that parses `"true"`, `"false"`, `"1"`, `"0"` (case-insensitive) and returns an error for other inputs. Write a table-driven test with at least 6 cases covering valid inputs, invalid inputs, and edge cases (empty string, mixed case).

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/03-table-driven-tests
go test -v ./curriculum/modules/06-testing-debugging-refactoring/lessons/03-table-driven-tests
```

## Review questions

1. What are the advantages of table-driven tests over writing individual test functions for each case?
2. How does `t.Run` affect test output and selective execution with `-run`?
3. Before Go 1.22, why was `tt := tt` necessary inside a range loop when using goroutines?
4. What fields should every table-driven test case struct include?
5. Given `go test -run TestParseBool/valid`, which subtests will execute?

## NEXT UP

Subtests — deeper patterns for `t.Run`, parallelism, setup/teardown, and selective execution.
