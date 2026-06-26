# Lesson 04: How to run code in this repository

## Learning objective

By the end of this lesson, you will understand the difference between `go run`, `go test`, `go build`, and `go vet`. You will know when to use each command and what each one tells you about your code. You will demonstrate this by running a simple calculator program using each command.

## Why this matters

Every software engineer runs code hundreds of times per day. The command you choose determines what you learn: `go run` tells you "does it work?", `go test` tells you "does it still work?", `go build` tells you "can it ship?", and `go vet` tells you "is it suspicious?". Using the wrong command wastes time and hides bugs. Mastering these four commands is the foundation of every other skill in this curriculum.

## Mental model

Think of each command as a different inspection tool for a house.

- `go run` is like walking through the house turning on lights. You see if anything is immediately broken.
- `go test` is like checking every outlet with a tester. You verify that every individual circuit works correctly.
- `go build` is like taking a photograph of the fully assembled house. It confirms that all the pieces fit together into a single structure that can be shipped.
- `go vet` is like a building inspector looking for code violations — wires that are unsafe, missing smoke detectors, or structural issues that are technically legal but dangerous.

## Core idea

| Command | What it does | When to use it | Output |
|---------|-------------|----------------|--------|
| `go run .` | Compiles and runs the main package immediately | While developing, to see quick results | Program output to stdout/stderr |
| `go test .` | Compiles test files and runs all test functions | Before and after making changes, to verify correctness | PASS/FAIL with test details |
| `go build .` | Compiles the package into a binary (or checks that it compiles) | Before committing, to ensure the code compiles standalone | Binary file or "nothing" if -o is not set |
| `go vet` | Analyzes the package for suspicious constructs | Before committing, to catch common mistakes | Warnings for suspicious code, or nothing if clean |

The chart order is intentional: run first to explore, test to verify, build to ship, vet to audit.

## Under the hood

When you run `go run .`, the Go toolchain compiles the package to a temporary directory, executes the resulting binary, and then deletes the temporary file. You see the program's output immediately.

`go test .` works differently. It compiles both the `_test.go` files and the package files into a separate test binary. This binary calls each function that matches the pattern `func TestXxx(t *testing.T)`. Each test function receives a `*testing.T` value that provides methods like `Errorf`, `Fatalf`, and `Log`. When a test calls `t.Fatalf`, the test stops immediately. When it calls `t.Errorf`, the test reports the failure but continues running other checks.

`go build .` compiles the package and writes the binary to the current directory (for a `package main`) or checks that a library package compiles (for other packages). The binary name matches the directory name.

`go vet` runs analyzers that look for common mistakes: passing the wrong number of arguments to `Printf`, unreachable code, or variable shadowing. Each analyzer is a separate Go program that reads the package's abstract syntax tree.

## How Go uses it

The Go project itself uses all four commands continuously. Every pull request to the Go standard library must pass `go test ./...` and `go vet ./...`. The `go build` command is used in release pipelines to produce binaries for all supported platforms. The Go team treats a failing vet check as a blocker, not a suggestion.

In this curriculum, you will use these same commands in every module. Each lesson's README tells you the exact command to run. By the end of Module 00, you will have run each command dozens of times.

## Go example

```go
package main

import (
	"errors"
	"fmt"
	"os"
)

func Add(a, b int) int {
	return a + b
}

func Subtract(a, b int) int {
	return a - b
}

func Multiply(a, b int) int {
	return a * b
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run . <op> <a> <b>")
		fmt.Println("Ops: add, sub, mul, div")
		return
	}
	op := os.Args[1]
	a := atoi(os.Args[2])
	b := atoi(os.Args[3])

	switch op {
	case "add":
		fmt.Println(Add(a, b))
	case "sub":
		fmt.Println(Subtract(a, b))
	case "mul":
		fmt.Println(Multiply(a, b))
	case "div":
		result, err := Divide(a, b)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Println(result)
	default:
		fmt.Fprintln(os.Stderr, "Unknown operation:", op)
		os.Exit(1)
	}
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}
```

## Step-by-step execution

1. If you copy the inline code into a `.go` file, `go run . add 3 4` compiles the package, runs `main()` with `os.Args = ["...", "add", "3", "4"]`. The switch matches `"add"`, calls `Add(3, 4)`, and prints `7`.
2. `go test .` compiles the test files alongside the package files. It calls each `TestXxx` function. `TestAdd` runs four cases: 2+3=5, -1+1=0, 0+0=0, 100+200=300. If any case fails, `t.Errorf` reports it and `go test` exits with a non-zero status.
3. `go build .` writes a binary to the current directory. You can run it directly.
4. `go vet .` checks for suspicious patterns. If it finds one, it prints a warning to stderr. If the code is clean, it prints nothing.

## Common mistakes

| Mistake | Why it happens | Fix |
|---------|---------------|-----|
| Running `go run .` from the wrong directory | The `.` means "current directory." If you are in the repository root, `go run .` tries to compile the root module as a binary, which fails if the root has no main package. | Always use the full relative path: `go run ./curriculum/modules/00-orientation/lessons/04-how-to-run-code-in-this-repository/` |
| Forgetting to pass arguments to `go run` | `go run .` runs `main()` with no arguments. The calculator's `main()` prints the Usage message when `len(os.Args) < 4`. | Run `go run . add 2 3` to see output, not just the usage text. |
| Confusing `go test .` with `go run .` | Both commands start with `go`, but they do very different things. | Remember: test is for verification, run is for exploration. |
| Ignoring vet warnings | `go vet` warnings are not errors, so the command succeeds (exit code 0). But the warnings point to real bugs. | Treat every vet warning as a bug report. Fix it before committing. |

## Debugging walkthrough

Scenario: You run `go test .` and get a compilation error: "undefined: Add".

Step 1: Check that the test file and the file defining `Add` are in the same directory. Go test compiles all `*_test.go` files with the package files. If the test file is in a different directory, the function `Add` will not be visible.

Step 2: Check the package declaration at the top of both files. Both must say `package main` (the same package). If one says `package main_test`, the test is in an external test package and cannot access unexported symbols.

Step 3: If the packages match but `Add` is still undefined, check whether `Add` is exported (capital A). Go visibility rules: a name starting with a capital letter is exported and visible to external test packages. A lowercase name is visible only within the same package.

Step 4: Run `go vet .` — vet can sometimes detect name mismatches or missing imports that cause confusing test failures.

Scenario: You run `go vet .` and get "Printf: missing arguments for %s".

Step 1: Look at the line number in the warning.
Step 2: Find the `fmt.Sprintf` or `fmt.Printf` call on that line.
Step 3: Count the format verbs (like `%s`, `%d`) and the arguments. If you have two verbs but only one argument, add the missing argument.
Step 4: Run `go vet .` again to confirm the warning is gone.

## Production notes

- In CI/CD pipelines, `go test ./...` is always run. `go vet ./...` should also be run and failures should fail the build. Many teams use `golangci-lint` which includes vet and dozens of other linters.
- `go build` is typically run with flags for the target OS and architecture (e.g., `GOOS=linux GOARCH=amd64 go build .`) to produce deployment artifacts.
- Always run `go vet` before `go build` in CI. Catching a vet warning early is cheaper than debugging a production issue.
- The calculator's `atoi` function is deliberately naive. Production code should use `strconv.Atoi` or `fmt.Sscanf` for robust parsing.

## Performance implications

- `go run .` compiles the package every time. For small programs like this calculator, compilation takes < 1 second. For large projects (millions of lines), it can take minutes. Use `go build` once and run the binary directly for repeated testing.
- `go test .` caches results. If the source files have not changed, `go test` prints `(cached)` and returns immediately. Use `go test -count=1 .` to bypass the cache and force a fresh run.
- The test binary is deleted after `go test` completes. Use `go test -c` to keep the test binary for later use.
- `go vet` is fast because it only analyzes the AST; it does not run the code. It typically completes in milliseconds for small packages.

## Practice task

Add a new operation `mod` (modulus) to the calculator. Create a `Mod(a, b int) (int, error)` function that returns `a % b` and an error if `b == 0`. Add a case for `"mod"` in the switch statement. Add a table-driven test for `Mod` that includes a zero-divisor case. Run `go test .` and `go vet .` to verify your changes.

## Tests / verification

The inline code example is for reading and understanding. To verify your understanding, complete the practice task above and check your answers against the description. You can also copy the inline code into a local `.go` file and run `go run .` and `go test .` in that directory to experiment with the output.

## Review questions

1. What is the difference between `go run .` and `go build .`?
2. Why does `go test .` succeed even when some tests fail (as long as no test panics)?
3. What does `go vet` check that `go build` does not?
4. In the calculator program, why does the `main()` function check `len(os.Args) < 4`?
5. What would happen if you ran `go run .` from a subdirectory that has no `.go` files?

## NEXT UP

[Lesson 05: How starter folders work](../05-how-starter-folders-work/README.md) — where you learn the `_starter/` and `_solution/` workflow used throughout this curriculum.
