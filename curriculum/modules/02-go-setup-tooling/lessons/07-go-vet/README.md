# go vet

## Learning objective

By the end of this lesson, you will understand what `go vet` does, the kinds of bugs it catches, and how to use it to detect suspicious code patterns that compile but are likely wrong.

## Why this matters

Compilation success does not mean code correctness. `go vet` catches classes of bugs that the compiler does not check: mismatched printf verbs, unused variables, defer misuse, and suspicious assignments. Running `go vet` before every commit prevents subtle bugs from reaching production.

## Mental model

If `go build` checks whether your code is valid Go, `go vet` checks whether your code is *sensible* Go. It looks for patterns that are syntactically correct but logically suspicious — a `fmt.Printf` call with the wrong format verb, an assignment that is never read, a lock that is copied by value.

## Visual model

```mermaid
flowchart TD
    SRC["Go Source Code\n(Syntactically valid Go)"] --> TC["Compiler Frontend\n(Type-Checked AST + Types Info)"]

    subgraph Analyzers["go/analysis Framework Pass"]
        A_PRINTF["printf: fmt.Printf(\"%d\", \"hello\")"]
        A_LOCK["copylocks: sync.Mutex passed by value"]
        A_UNREACH["unreachable: dead code following return"]
        A_SHADOW["shadow: accidental inner-scope := shadowing"]
    end

    TC --> A_PRINTF
    TC --> A_LOCK
    TC --> A_UNREACH
    TC --> A_SHADOW

    A_PRINTF --> OUT["CLI Diagnostics Output:\n./main.go:12: Printf format %d has arg of type string\nexit status 1"]
    A_LOCK --> OUT
    A_UNREACH --> OUT
    A_SHADOW --> OUT
```

## Core idea

`go vet` runs static analysis on your Go source to detect suspicious constructs. It compiles the package, then runs built-in analyzers on the typed AST to find patterns that are valid Go but likely unintended.

## Under the hood

`go vet` uses the `go/analysis` framework. Each analyzer implements the `go/analysis.Analyzer` interface and operates on the compiled AST and type information. Built-in analyzers include:

- `printf` — checks format verb and argument type matching
- `shadow` — detects variables shadowed by inner declarations
- `copylocks` — detects passing a mutex or similar type by value
- `unreachable` — detects dead code after a return or panic
- `tests` — checks that test functions have the correct signature

## How Go uses it

The Go standard library is vetted before every release. `go vet` is integrated into the Go toolchain and runs as part of the standard CI process. Key commands:

- `go vet .` — vet the current package
- `go vet ./...` — vet all packages in the module
- `go vet -vettool=$(which staticcheck) .` — run `staticcheck` as a vet pass

## Go example

```go
package main

import (
	"fmt"
	"strings"
)

// VetDiagnostic represents an issue flagged by static analysis.
type VetDiagnostic struct {
	CheckName string
	Severity  string
	Message   string
	Example   string
}

// DiagnoseCommonVetIssues demonstrates the primary bug patterns go vet catches.
func DiagnoseCommonVetIssues() []VetDiagnostic {
	return []VetDiagnostic{
		{
			CheckName: "printf",
			Severity:  "warning",
			Message:   "Printf format %d has arg 'hello' of wrong type string",
			Example:   `fmt.Printf("%d\n", "hello")`,
		},
		{
			CheckName: "copylocks",
			Severity:  "error",
			Message:   "assignment copies lock value to pass-by-value receiver",
			Example:   `type SafeCounter struct { mu sync.Mutex; val int }; func (c SafeCounter) Inc() {}`,
		},
		{
			CheckName: "unreachable",
			Severity:  "warning",
			Message:   "unreachable code after return or panic",
			Example:   `func foo() int { return 1; fmt.Println("never reached") }`,
		},
		{
			CheckName: "nilfunc",
			Severity:  "warning",
			Message:   "comparison of function f == nil is always false",
			Example:   `if myFunc == nil { ... } (when myFunc is a declared func)`,
		},
	}
}

// CheckFormatArgs checks whether a printf format specifier matches argument counts.
func CheckFormatArgs(format string, argCount int) error {
	specifiers := 0
	for i := 0; i < len(format); i++ {
		if format[i] == '%' {
			if i+1 < len(format) && format[i+1] == '%' {
				i++ // escaped %%
				continue
			}
			specifiers++
		}
	}

	if specifiers != argCount {
		return fmt.Errorf("format %q requires %d args, but %d were provided", format, specifiers, argCount)
	}
	return nil
}

func main() {
	fmt.Println("=== Deconstructing 'go vet' Static Analysis Engine ===")

	diagnostics := DiagnoseCommonVetIssues()
	for i, d := range diagnostics {
		fmt.Printf("[%d] Checker: %-12s | Severity: %-7s\n", i+1, d.CheckName, strings.ToUpper(d.Severity))
		fmt.Printf("    Problem: %s\n", d.Message)
		fmt.Printf("    Bad Pattern: %s\n\n", d.Example)
	}

	// Demonstration of correct code verified by go vet
	user := "Senior Gopher"
	years := 5
	fmt.Printf("Verified Clean: User %s has %d years of Go experience.\n", user, years)
	fmt.Println("Status: go vet runs in CI to reject subtle bugs before compilation/deployment.")
}
```

This program categorizes the most dangerous bugs that escape compiler syntax checks:
- **`printf`**: Detects format verbs (`%d`, `%s`, `%f`) mismatched with supplied runtime argument types.
- **`copylocks`**: Catches concurrency race bugs where structs containing `sync.Mutex` or `sync.WaitGroup` are passed by value instead of pointer.
- **`unreachable`**: Identifies dead code blocks positioned after `return` or `panic`.
- **`nilfunc`**: Flags comparisons between function signatures and `nil` that can never evaluate to true.

## Step-by-step execution

1. Run the static analyzer demonstration:
   ```bash
   go run .
   ```
2. Run the automated unit tests:
   ```bash
   go test -v .
   ```
3. Run the official Go static analysis tool on this lesson:
   ```bash
   go vet .
   ```

## Common mistakes

- **Confusing `go vet` with `gofmt`**: `gofmt` fixes formatting; `go vet` finds bugs. Both should run before every commit, but they catch different things.
- **Ignoring vet warnings because the code compiles**: A `go vet` warning means the code is probably wrong even though it compiles. Treat warnings as errors.
- **Running only `go vet` without fixing the issues**: `go vet` only reports issues — it does not fix them. Each warning must be manually corrected.

## Debugging walkthrough

**Scenario**: `go vet .` reports `Printf format %d has arg name of wrong type string`.

1. Look at the line number in the warning. Check the `fmt.Printf` call on that line.
2. The format verb `%d` expects an integer, but the argument is a string. Change `%d` to `%s`.
3. Run `go vet .` again to confirm the warning is gone.

**Scenario**: `go vet .` reports `declared and not used`.

1. You declared a variable but never read its value. Either use the variable (e.g., print it) or remove the declaration.
2. If you intentionally do not need the value, assign it to `_` (blank identifier): `_ = result`.

## Production notes

- CI pipelines should fail on `go vet ./...` warnings. Add `go vet ./...` as a step before running tests.
- Many teams use `golangci-lint` which includes `go vet` and dozens of additional linters in a single run.
- `go vet` is fast — it runs in parallel across packages and completes in seconds even for large modules.

## Performance implications

- `go vet` runs after compilation and is fast — typically under a second per package.
- The analysis is built into the Go toolchain, so there are no additional dependencies to install.
- Running `go vet` on every file in a module adds negligible time to the CI pipeline.

## Practice task

1. Run `go vet .` from this lesson directory — confirm zero warnings.
2. Create a file with a `fmt.Printf` format mismatch (e.g., `%d` with a string argument) and run `go vet .` to see the warning.
3. Add an unused variable and run `go vet .` again.
4. Fix both issues and confirm `go vet .` produces no output.

## Tests / verification

```bash
go vet .
```

Expected output: nothing (exit 0). Any output means there are issues to fix.

```bash
go run .
```

Expected output:
```
Hello, Go!
go vet found no issues in this file.
```

## Review questions

1. What is the difference between a compiler error and a `go vet` warning?
2. What kind of bug does the `printf` analyzer catch?
3. What does the `copylocks` analyzer detect?
4. How do you run `go vet` on all packages in a module?
5. Why should CI fail on `go vet` warnings even though the code compiles?

## NEXT UP

[go doc](../08-go-doc/README.md) — read and write Go documentation from the terminal.
