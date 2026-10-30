# Hello World

## Learning objective

By the end of this lesson, you will understand the structure of a minimal Go program — package declaration, import, function, and statement — and you will write and run your own Hello World program.

## Why this matters

Hello World is the smallest complete Go program. Every Go program you write will share the same structure: a `package` declaration, the `imports` it needs, and a `func main()` as the entry point. Understanding this skeleton means you can recognize the shape of any Go program from the start.

## Mental model

A Go source file is a text file with three sections: package declaration at the top, imports in the middle, and code at the bottom. The `main` function is where execution begins — like the front door of your program. The `fmt.Println` call sends text to your terminal.

## Visual model

```mermaid
flowchart TD
    subgraph SourceFile["main.go Source File"]
        PKG["package main (Declares executable binary entrypoint)"]
        IMP["import \"fmt\" (Imports formatted I/O standard package)"]
        FN["func main() (Entrypoint: execution begins here)"]
        CALL["fmt.Println(\"Hello, World!\")"]

        PKG --> IMP --> FN --> CALL
    end

    CALL -->|Emits bytes to file descriptor 1| STDOUT["OS Standard Output (Terminal Console)"]
```

## Core idea

Every Go program starts with `package main` and a `func main()`. The `import` statement brings in packages from the standard library (or third parties) that provide functionality like printing, math, or file access.

## Under the hood

When you run `go run main.go`, the Go compiler:

1. Reads the source file and parses it into an abstract syntax tree (AST).
2. Resolves the import `"fmt"` to the standard library package.
3. Compiles the code into a temporary binary.
4. Executes the binary, which calls `main()`, which calls `fmt.Println`, which writes to stdout.
5. Deletes the temporary binary.

The entire process takes milliseconds for a small program.

## How Go uses it

The `fmt` package is the most-used package in the Go standard library. `fmt.Println` prints values separated by spaces followed by a newline. `fmt.Printf` prints formatted output using verbs like `%s` (string), `%d` (integer), and `%v` (default format). Every Go developer uses `fmt` daily for logging, debugging, and user-facing output.

## Go example

```go
package main

import (
	"fmt"
	"io"
	"os"
)

// FormatGreeting constructs a standard Go welcome message.
func FormatGreeting(target string) string {
	if target == "" {
		target = "World"
	}
	return fmt.Sprintf("Hello, %s!", target)
}

// WriteGreeting sends formatted output to an io.Writer (e.g., os.Stdout).
func WriteGreeting(w io.Writer, target string) error {
	_, err := fmt.Fprintln(w, FormatGreeting(target))
	return err
}

func main() {
	// Demonstrating output to standard output stream (os.Stdout)
	if err := WriteGreeting(os.Stdout, "World"); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing greeting: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Program completed successfully via package main -> func main() entry point.")
}
```

This program demonstrates three core concepts:
1. **Entry Point Requirement**: `package main` coupled with `func main()` forms the mandatory execution entry point. Without both, the Go toolchain will not compile a standalone binary.
2. **Deterministic String Formatting**: `FormatGreeting` isolates business logic from I/O, allowing pure unit testing without terminal side effects.
3. **Standard Stream Writing**: `WriteGreeting` writes to `io.Writer`, accepting `os.Stdout` (File Descriptor 1) in production, but allowing `bytes.Buffer` in automated tests.

## Step-by-step execution

1. Open your terminal and navigate to this lesson directory.
2. Inspect `main.go` and observe the separation between string formatting and stream writing.
3. Execute the program:
   ```bash
   go run .
   ```
4. Run the automated unit tests:
   ```bash
   go test -v .
   ```

## Common mistakes

- **Missing `package main`**: Without it, `go run` produces `package is not a main package`. Only `package main` can produce an executable.
- **Wrong import syntax**: Imports use double quotes, not single quotes: `import "fmt"` not `import 'fmt'`.
- **Missing braces**: Every opening `{` must have a closing `}`. Go forces you to put the opening brace on the same line as the declaration — a newline before `{` is a syntax error.
- **Misspelling `Println`**: Go is case-sensitive. `println` (lowercase) is a built-in that works differently from `fmt.Println`. Use `fmt.Println` for portability.

## Debugging walkthrough

**Scenario**: You run `go run main.go` and get `package is not a main package`.

1. Check the first line of your file. It must be exactly `package main` with no extra spaces.
2. If the file is in a subdirectory, ensure you run `go run` from the directory containing the file, or specify the full path.
3. If you have multiple `.go` files in the directory, only one can have `func main()`. Remove or rename extra main functions.

**Scenario**: You run `go run main.go` and get `undefined: fmt`.

1. Check that you have `import "fmt"` in your file.
2. Ensure the import is after the package declaration and before the function definitions.
3. Run `go mod tidy` from the module root to ensure the module file is up to date.

## Production notes

- In real Go projects, `main.go` is often minimal — it parses flags, sets up dependencies, and calls into other packages. The core logic lives in separate packages.
- Many production projects use `cmd/` layout: `cmd/server/main.go`, `cmd/cli/main.go` for multiple binaries in one repo.
- Never put business logic in `main.go`. It should only wire things together.

## Performance implications

- A Hello World binary compiled with `go build` is about 1.5 MB. Most of that size is the Go runtime (scheduler, garbage collector), not your code.
- `go run` does not leave a binary on disk — it compiles to a temp directory and runs from there. Use `go build` to produce a permanent binary.

## Practice task

1. Run `go run .` from this lesson directory to see the greeting.
2. Modify the program to print your name instead of "World".
3. Add a third line that prints your favorite programming language.

## Tests / verification

```bash
go run .
```

Expected output:
```
Hello, World!
My Go installation works.
```

After your modification, the output should include your name and favorite language.

## Review questions

1. What is the purpose of `package main`?
2. What does `import "fmt"` do?
3. Why does `func main()` need the `{` on the same line?
4. What happens if you omit `package main` from a program you want to run?
5. What is the difference between `go run main.go` and `go build main.go`?

## NEXT UP

[go run](../03-go-run/README.md) — learn how `go run` compiles and executes Go programs.
