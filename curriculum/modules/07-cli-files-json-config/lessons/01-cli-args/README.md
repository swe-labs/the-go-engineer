# CLI args

## Learning objective

Read and validate command-line arguments using `os.Args`, distinguish positional args from the program name, and parse string arguments into typed values with proper error handling.

## Why this matters

Every CLI tool, dev script, and server binary accepts arguments. The `os.Args` slice is the raw gateway to everything your program receives from the shell. Whether you build a tiny utility or a production microservice, understanding how to safely extract and validate positional arguments is the foundation of every Go CLI.

## Mental model

Imagine your program as a function `main([]string)` — the operating system hands it a slice of strings: `os.Args`. The first element (`os.Args[0]`) is always the path used to invoke the program. Everything after it (`os.Args[1:]`) are the actual arguments the user typed. Your job is to validate and convert those strings.

```
shell:  go run main.go Alice 30
              |                |
              v                v
        os.Args[0]        os.Args[1:]
```

## Core idea

`os.Args` is a `[]string` (a slice of strings) declared in the `os` package. It is populated by the Go runtime before `main()` executes. The length `len(os.Args)` tells you how many tokens were on the command line.

- Positional arguments are unlabeled values like `Alice 30`.
- Flags are labeled values like `--name=Alice` (covered in the next lesson).
- There is no built-in parser for positional args — you must write the validation yourself.

## Under the hood

The Go runtime receives the argument vector from the operating system's `exec` system call. On Windows this comes from `GetCommandLineW` parsed into argc/argv style; on Unix it is the `argv` pointer passed to `main`. The runtime copies these C strings into Go-managed memory and presents them as `os.Args`. There is no limit enforced by Go, though the OS may impose one (typically 128KB–2MB of total argument data on Linux, ~32K on older Windows).

## How Go uses it

- `os.Args` is a package-level variable, set once at startup.
- It is read-only after initialization.
- The zero-element `os.Args[0]` is used in usage messages, error reports, and `log` output.
- You must handle `len(os.Args)` checks before indexing to avoid panics.
- Parsing is string-based; use `strconv`, `fmt.Sscanf`, or custom logic for typed conversion.

## Go example

```go
package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <name> [age]")
		return
	}
	name := os.Args[1]
	if len(os.Args) >= 3 {
		age, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Invalid age %q: must be a number\n", os.Args[2])
			return
		}
		fmt.Printf("Hello %s, age %d\n", name, age)
		return
	}
	fmt.Printf("Hello %s\n", name)
}
```

## Step-by-step execution

Run `go run main.go Alice 30`:

1. `len(os.Args)` → `3` (program, "Alice", "30").
2. Condition `len(os.Args) < 2` is false, skip usage block.
3. `name = os.Args[1]` → `"Alice"`.
4. `len(os.Args) >= 3` is true.
5. `strconv.Atoi("30")` → `30, nil`.
6. Print `"Hello Alice, age 30"`.

Run `go run main.go`:

1. `len(os.Args)` → `1`.
2. Condition `len(os.Args) < 2` is true.
3. Print usage message and return early.

## Common mistakes

- **Indexing `os.Args[1]` without checking length**: Panics with `index out of range [1] with length 1`.
- **Assuming `os.Args[0]` is the source file**: It is the compiled binary path (e.g., `/tmp/go-build123/main`), not the `.go` file.
- **Forgetting `strconv` errors**: `strconv.Atoi` returns an error on non-numeric input; ignoring it silently uses zero.
- **Treating `os.Args` as mutable**: While modifying the slice is possible, doing so is confusing and never necessary.

## Debugging walkthrough

Buggy code:

```go
func main() {
	fmt.Println("Hello", os.Args[1])
}
```

**Symptom**: Panics when run without arguments.

**Investigation**: Add `fmt.Println("args:", os.Args)` before the print. You'll see `args: [/tmp/go-build.../exe/main]` — length 1 with no user-provided arg.

**Fix**: Guard with `if len(os.Args) < 2 { ... }`.

## Production notes

- Always print a usage message when required args are missing — do not panic or silently exit 0.
- Use `os.Args[0]` in the usage line (e.g., `Usage: %s <name>`, `os.Args[0]`) so the message works regardless of how the binary is invoked.
- For production CLIs, prefer the `flag` package or a third-party library like `cobra` over raw `os.Args` — but understand `os.Args` as the foundation.
- Validate inputs early and fail with a clear message. A confusing error from a deep parse path is worse than an immediate rejection.

## Performance implications

- `os.Args` is allocated once at startup; reading it costs nothing.
- Parsing large numbers of arguments with `strconv` is fast — microseconds for hundreds of args.
- The real cost is in validation logic, not slice access.

## Practice task

Write a program that accepts three positional arguments: `name`, `count` (int), and `price` (float64). Print a formatted receipt line: `"Total for <name>: <count> x <price> = <total>"` where total = count x price. If any argument is missing or the wrong type, print a usage message. Save your code as `main.go`.

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/01-cli-args
go test ./curriculum/modules/07-cli-files-json-config/lessons/01-cli-args
```

## Review questions

1. What is the value of `os.Args[0]` when you run `go run main.go hello`?
2. What happens if you access `os.Args[2]` when only one argument is provided?
3. How would you safely read an optional third argument that may or may not be present?
4. Why should you call `strconv.Atoi` instead of using a type assertion on `os.Args[1]`?
5. What does `len(os.Args)` equal when the program is run with no additional arguments?

## NEXT UP

Flags — parsing named command-line options with the `flag` package.
