# CLI testability

## Learning objective

Write testable CLI applications by injecting `io.Writer` dependencies for output, replacing `os.Stdin` in tests, using golden files for output comparison, and testing exit codes.

## Why this matters

CLI tools are the backbone of DevOps, CI/CD, and developer tooling. But CLI apps are notoriously hard to test because they print to stdout, read from stdin, and call `os.Exit`. Untested CLI tools break silently in production pipelines. By applying dependency injection and testing patterns, you can achieve the same confidence in CLI apps as in library code.

## Mental model

Think of a CLI tool as a function: `func run(input io.Reader, output io.Writer, args []string) exitCode`. By pushing side effects (printing, reading, exiting) to explicit parameters, the core logic becomes testable. The `main()` function is just a thin wrapper that wires real `os.Stdin`, `os.Stdout`, and `os.Args` into `run()`.

## Core idea

Three key patterns make CLI apps testable:

1. **Inject `io.Writer` for output**: Instead of `fmt.Println(...)`, write to a configurable writer. In production it's `os.Stdout`; in tests it's `bytes.Buffer`.
2. **Inject `io.Reader` for input**: Read commands from a configurable reader. In production it's `os.Stdin`; in tests it's `strings.NewReader(...)`.
3. **Return exit codes instead of `os.Exit`**: Return an `int` from the core function. Only `main()` calls `os.Exit`.

**Golden file testing**: Write expected output to a `.golden` file, then compare actual output against it. If output changes, update the golden file (after careful review).

## Under the hood

`fmt.Fprintf(w, ...)` writes to any `io.Writer`. `fmt.Fscanln(r, ...)` reads from any `io.Reader`. These are the building blocks. By making your core function accept these interfaces:

```go
func Run(r io.Reader, w io.Writer, args []string) int {
    // core logic here
}
```

`main()` becomes trivial:

```go
func main() {
    os.Exit(Run(os.Stdin, os.Stdout, os.Args[1:]))
}
```

Testing becomes trivial:

```go
func TestRun(t *testing.T) {
    var buf bytes.Buffer
    exitCode := Run(strings.NewReader("input\n"), &buf, []string{})
    // assert buf.String() and exitCode
}
```

## How Go uses it

- **Testing CLI tools**: All major Go CLI frameworks (Cobra, urfave/cli) support writer injection.
- **Golden file tests**: Used extensively in the Go standard library (`testing/golden`) and projects like `gofmt`.
- **Pipe testing**: Simulate `cat file.txt | mytool` by providing file content as reader.
- **Error path testing**: Test what happens when stdout is closed (broken pipe) or stdin returns error.

## Go example

```go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Run is the core CLI logic. Returns exit code.
func Run(r io.Reader, w io.Writer, args []string) int {
	fs := flag.NewFlagSet("greet", flag.ContinueOnError)
	fs.SetOutput(w)
	upper := fs.Bool("upper", false, "convert to uppercase")
	name := fs.String("name", "World", "name to greet")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	greeting := fmt.Sprintf("Hello, %s!", *name)
	if *upper {
		greeting = strings.ToUpper(greeting)
	}
	fmt.Fprintln(w, greeting)

	// Read additional names from stdin (one per line)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		g := fmt.Sprintf("Hello, %s!", line)
		if *upper {
			g = strings.ToUpper(g)
		}
		fmt.Fprintln(w, g)
	}
	return 0
}

func main() {
	os.Exit(Run(os.Stdin, os.Stdout, os.Args[1:]))
}
```

## Step-by-step execution

1. `Run` receives `os.Stdin` (reader), `os.Stdout` (writer), and `os.Args[1:]` (slice).
2. `flag.NewFlagSet` creates a flag set that writes errors to `w` (not stderr by default).
3. `fs.Parse(args)` parses flags like `--upper --name=Alice`.
4. Greeting is constructed, optionally uppercased, written to `w`.
5. If stdin has data (piped or typed), scanner reads each line and prints a greeting.
6. Returns 0 on success.

In tests:
1. `Run(strings.NewReader("Bob\nCarol\n"), &buf, []string{"--upper"})` is called.
2. Same logic but reads from the string reader and writes to the buffer.
3. Assert `buf.String()` contains `"HELLO, WORLD!"`, `"HELLO, BOB!"`, `"HELLO, CAROL!"`.

## Common mistakes

- Mistake: Calling `os.Exit(1)` inside the core function — prevents clean-up and makes testing impossible.
  - Fix: Return an exit code and let `main()` call `os.Exit`.

- Mistake: Importing `flag` and using `flag.Parse()` in core function — the global `flag` package reads `os.Args` directly, which you can't replace in tests.
  - Fix: Use `flag.NewFlagSet` and pass `args` explicitly.

- Mistake: Not testing with piped input — the tool works interactively but fails in CI pipelines.
  - Fix: Always test with a reader that provides input.

- Mistake: Golden files that are never reviewed — stale golden files hide regressions.
  - Fix: Review golden file diffs in PRs. Update them deliberately, not automatically.

## Debugging walkthrough

Test fails with "unexpected exit":

```go
func TestSomething(t *testing.T) {
	os.Exit(Run(...)) // This kills the test process!
}
```

**Symptom**: Test runner reports the test exited unexpectedly or the test process terminates.

**Root cause**: Calling `os.Exit` in a test function terminates the entire test binary, not just the current test.

**Fix**: Only `main()` should call `os.Exit`. Core functions return exit codes:

```go
func TestSomething(t *testing.T) {
	code := Run(...)
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}
```

## Production notes

- **Always inject dependencies** in production CLI code too — it makes the code modular and testable from day one.
- **Golden files should be checked into version control** but reviewed carefully in PRs. Use `-update` flag to regenerate: `go test -update`.
- **Test edge cases**: empty input, very long input, unicode, binary input on stdin.
- **Exit code convention**: 0 = success, 1 = general error, 2 = misuse (wrong args), 3+ = tool-specific.
- **Consider `testing/iotest`** for simulating error readers and writers (timeouts, broken pipes).

## Performance implications

- Writer injection adds zero overhead — `fmt.Fprintf` is as fast as `fmt.Printf` (which internally writes to `os.Stdout`).
- Reading from a `strings.Reader` or `bytes.Buffer` in tests is faster than real I/O.
- Golden file comparison is a simple byte comparison — negligible cost.
- The overhead of abstraction (interfaces) is inlined by the compiler — no runtime cost.

## Practice task

Write a function `Run(r io.Reader, w io.Writer, args []string) int` for a `wc`-like tool that counts lines in input. If the `-l` flag is passed, count lines (default). If `-w`, count words. If `-c`, count characters. Read from stdin. Return 0 on success, 1 on error. Write a test that pipes "hello world\nfoo bar baz\n" and checks line count is 2 (with `-l`), word count is 5 (with `-w`), char count is 23 (with `-c`).

## Tests / verification

```bash
go run ./curriculum/modules/07-cli-files-json-config/lessons/21-cli-testability
go test ./curriculum/modules/07-cli-files-json-config/lessons/21-cli-testability
```

## Review questions

1. Why should you avoid `os.Exit` inside a CLI tool's core logic function?
2. How does injecting `io.Writer` make a CLI tool testable?
3. What is a golden file test and when would you use it?
4. Why use `flag.NewFlagSet` instead of the global `flag.Parse()` in a testable CLI?
5. What exit code convention is standard for Unix CLI tools?

## NEXT UP

Congratulations on completing Module 07! You now know how to build CLI tools, read/write files, handle JSON, and manage configuration in Go. Next up: Module 08 — HTTP, REST APIs, and Web Services.
