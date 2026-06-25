# Refactoring safely

## Learning objective

Refactor Go code with confidence by using a test safety net, applying small-step transformations, and leveraging `go vet` to catch structural errors.

## Why this matters

Refactoring changes the structure of code without changing its behaviour. Without tests, every refactor is a blind gamble — you hope the output is the same, but you cannot prove it. With tests, you can refactor aggressively, run the tests after each small change, and immediately know if you broke something. This is the difference between refactoring as a disciplined practice and "rewriting and hoping."

## Mental model

Refactoring is a sequence of behaviour-preserving transformations. Each transformation is small (rename a variable, extract a function, change a signature) and reversible. The test suite acts as a canary: if all tests still pass, the behaviour is preserved. If a test fails, the last transformation introduced a behavioural change — undo it or fix it.

```
[working code] → (small refactor) → [run tests] → (pass?) → [next refactor]
                                          ↓ (fail?)
                                    [undo or fix]
```

## Core idea

Safe refactoring relies on three pillars:

1. **Test safety net**: Before any refactor, the test suite must pass. After each change, run the tests. If they fail, the change was not behaviour-preserving.
2. **Small steps**: Change one thing at a time. Rename a variable, then extract a function, then change a signature. Never combine steps.
3. **Tool support**: `go vet` catches structural errors (e.g. unreachable code, mistmatched printf args). `gofmt` enforces consistent formatting. IDEs provide automated refactoring (rename, extract) that can be safer than manual editing.

Common refactoring operations in Go:

| Refactoring | How |
|---|---|
| Extract function | Move a block of code into a new function with named parameters and return values |
| Rename | Change a variable, function, or type name (IDE or `gorename`) |
| Change signature | Add/remove/reorder parameters (update all call sites) |
| Inline function | Replace a function call with its body |
| Move declaration | Move a type or function to a different file or package |
| Introduce interface | Extract a method set into an interface, use it as a parameter type |

## Under the hood

`go vet` runs static analysis passes (checkers) over the source code. Each checker is a Go program that analyses the AST and reports diagnostics. The `printf` checker verifies format strings match argument types. The `copylocks` checker warns about copying `sync.Mutex` by value. The `nilfunc` checker detects comparisons of functions to nil. `go vet` is not a linter for style — it detects likely bugs.

Automated refactoring tools (like `gorename` or IDE refactorings) operate on the AST and symbol table. They rename a symbol everywhere it appears, including in other files in the same package. The Go compiler then re-resolves all references; if the refactoring introduced ambiguities or broke imports, the compiler reports errors.

## How Go uses it

- **`gofmt -w`**: Rewrites Go source to canonical formatting. Always run before committing.
- **`go vet ./...`**: Runs all checkers on the package. Essential in CI.
- **`gorename`**: A standalone tool (`golang.org/x/tools/cmd/gorename`) for renaming symbols across packages.
- **`go fix`**: Migrates code to use newer APIs (e.g. `golang.org/x/net/context` → `context`).
- **IDE integrated refactoring**: VS Code (gopls), GoLand, and vim-go provide rename, extract function, and change signature with preview.

## Go example

Before refactoring (a single function doing too much):

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	names := []string{"alice", "bob", "charlie"}
	result := process(names)
	fmt.Println(result)
}

func process(names []string) string {
	result := ""
	for _, n := range names {
		result += strings.ToUpper(n[:1]) + strings.ToLower(n[1:]) + ","
	}
	return strings.TrimSuffix(result, ",")
}
```

After extracting a helper function:

```go
func process(names []string) string {
	result := ""
	for _, n := range names {
		result += capitalize(n) + ","
	}
	return strings.TrimSuffix(result, ",")
}

func capitalize(name string) string {
	return strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
}
```

Step-by-step refactor:
1. Identify the block: `strings.ToUpper(n[:1]) + strings.ToLower(n[1:])`.
2. Write test that captures current behaviour.
3. Extract into `capitalize(name string) string` — copy the block, create function, call it.
4. Run tests — confirm they pass.
5. Commit the extraction. Then consider further refactoring (e.g. using `strings.Join` for the main loop).

## Step-by-step execution

Refactoring `process` step by small step:

1. Write a test: `TestProcess` and `TestCapitalize` that fully cover all paths.
2. Verify tests pass: `go test ./...` → all green.
3. Extract `capitalize`:
   - Create `func capitalize(name string) string`.
   - Move the block into it.
   - Call it from `process`.
4. Run tests → all green.
5. Extract the join logic:
   - Create `func joinNames(names []string) string` that builds the comma-separated string.
   - Use `strings.Join` with a mapped slice.
6. Run tests → all green. (If red, fix: the refactor introduced a behavioural change.)

Each step is ~30 seconds. The entire refactor takes 2 minutes, with tests verifying every intermediate state.

## Common mistakes

- **Refactoring without tests**: You cannot prove behaviour preservation. Write tests first, even if they are for the code you are about to refactor.
- **Large refactors in one commit**: Combine 10 changes into one commit. A test failure cannot be attributed to a specific change. Keep commits small and atomic.
- **Refactoring and fixing bugs simultaneously**: If a test was failing before, the refactor might accidentally "fix" it, and you lose the regression coverage. Separate bug fixes from refactoring.
- **Not running `go vet`**: `go vet` catches errors that tests might miss (e.g. passing `int` to `%s`). Run it after every refactor.
- **Forgetting `gofmt`**: The compiler accepts non-idiomatic formatting. `gofmt` ensures consistency and prevents spurious diffs.

## Debugging walkthrough

You rename `process` to `formatNames` across the package. Tests pass. But another package in the module imports `process`:

```
$ go build ./...
# github.com/example/pkg
pkg/caller.go:10: undefined: process
```

The rename missed an external call site. Solution:
1. Use `gorename` from `golang.org/x/tools`: `gorename -from '"github.com/example/main".process' -to formatNames`. This handles cross-package renaming.
2. If using IDE refactoring, select "Find references" before renaming and verify all call sites are included.
3. After the rename, `go build ./...` reports all broken imports. Fix them one by one.

## Production notes

- **Refactoring in a monorepo**: Changes may affect many packages. Use `go build ./...` and `go test ./...` from the module root to catch all breakage.
- **CI pipeline**: Run `go vet ./...` and `go test ./...` on every PR. A red CI after refactoring means the change was not behaviour-preserving.
- **Code review**: Refactoring commits should be reviewed for correctness (does the new structure match the old behaviour?) and style. The reviewer can run the tests themselves to confirm.
- **Large-scale refactoring**: Use `gorename`, `eg` (example-based refactoring tool in `golang.org/x/tools`), or `gofmt -r` for automated rewrites across hundreds of files.

## Performance implications

- Refactoring for readability may introduce performance regressions:
  - Extracting a function adds a function call overhead (a few nanoseconds). In hot paths, this matters. The compiler may inline small functions, mitigating the cost.
  - Introducing an interface for testability replaces concrete method calls with dynamic dispatch (vtable lookup). This prevents inlining and adds ~1-5 ns per call.
  - `strings.Builder` is faster than `+=` for string concatenation. If your refactoring changes how strings are built, benchmark it.
- Always measure: run benchmarks before and after a refactoring when performance is critical.

## Practice task

The `main.go` contains a `process` function that capitalises names and joins them with commas. Write tests for `process` that cover empty input, single name, and multiple names. Then:
1. Extract a `capitalize` function.
2. Extract a `joinNames` function that uses `strings.Join`.
3. Rename `process` to `formatNames`.
4. Run `go vet ./...` and `go test ./...` after each step.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/19-refactoring-safely
go test ./curriculum/modules/06-testing-debugging-refactoring/lessons/19-refactoring-safely
go vet ./curriculum/modules/06-testing-debugging-refactoring/lessons/19-refactoring-safely
```

## Review questions

1. Why should you run the test suite before starting a refactoring?
2. What is the danger of refactoring in large batches rather than small steps?
3. What does `go vet` check that `go test` does not?
4. Why should you not fix bugs in the same commit as a refactoring?
5. How does extracting a function affect performance? When does the compiler mitigate this?

## NEXT UP

Fuzz testing — automatically generating random inputs to find edge cases and bugs.
