# Standard input and standard output

## Learning objective

Read from `os.Stdin` line by line, write to `os.Stdout` and `os.Stderr`, detect whether input is piped or interactive, and handle streaming input correctly.

## Why this matters

The three standard streams (stdin, stdout, stderr) are the UNIX pipe model that makes small tools composable. A Go program that reads stdin and writes stdout can be chained with any other program: `cat log.txt | go run filter.go | go run count.go`. This is the heart of the Unix philosophy and essential for every CLI engineer.

## Mental model

Imagine your program as a filter function: `func(stdin io.Reader) -> (stdout, stderr io.Writer)`. Data flows in through stdin, results go to stdout, errors and diagnostics go to stderr. Pipes connect one program's stdout to another's stdin.

```
echo "hello" | myprogram | go run transform.go
     stdin          stdout       stdin
```

## Core idea

- `os.Stdin` is an `*os.File` that implements `io.Reader`. Read from it with `bufio.Scanner` for line-by-line processing.
- `os.Stdout` is an `*os.File` that implements `io.Writer`. Write to it with `fmt.Println`, `fmt.Printf`, or direct `Write` calls.
- `os.Stderr` is the same but for error/diagnostic output — separated from stdout so pipes don't mix data with logs.
- Pipe detection: use `os.Stdin.Stat()` and check `ModeCharDevice` — if the mode has no `ModeCharDevice` bit, stdin is a pipe.

## Under the hood

Each process starts with three file descriptors: 0 (stdin), 1 (stdout), 2 (stderr). The shell sets these up before `exec`. When you pipe, the shell connects fd 1 of the left process to fd 0 of the right process via a kernel pipe buffer. `os.Stdin` is a Go wrapper around fd 0.

`bufio.Scanner` reads from the underlying reader in chunks (default 64KB buffer) and returns full lines. It handles `\n` and `\r\n` line endings transparently. Scanner stops when it hits `io.EOF` (pipe closed) or an error.

## How Go uses it

- `bufio.NewScanner(os.Stdin)` creates a scanner that reads from stdin.
- `scanner.Scan()` advances to the next line, returns `true` if data is available.
- `scanner.Text()` returns the current line as a string (without the trailing newline).
- `scanner.Err()` returns any I/O error encountered.
- `fmt.Println(...)` writes to stdout with a trailing newline.
- `fmt.Fprintln(os.Stderr, ...)` writes to stderr.
- `os.Stdin.Stat()` returns `FileInfo`; check `info.Mode() & os.ModeCharDevice` to detect pipes.

## Go example

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(strings.ToUpper(line))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
```

## Step-by-step execution

Run `echo "hello world" | go run main.go`:

1. `bufio.NewScanner(os.Stdin)` wraps stdin in a scanner.
2. First `scanner.Scan()` reads "hello world\n" from the pipe, returns `true`.
3. `scanner.Text()` returns `"hello world"` (no newline).
4. `strings.ToUpper("hello world")` → `"HELLO WORLD"`.
5. `fmt.Println("HELLO WORLD")` writes to stdout.
6. Second `scanner.Scan()` reads `io.EOF` (pipe closed), returns `false`.
7. Loop exits, `scanner.Err()` is nil, program exits cleanly.

## Common mistakes

- **Using `fmt.Scan` instead of `bufio.Scanner`**: `fmt.Scan` reads word-by-word, not line-by-line, and has confusing whitespace handling.
- **Ignoring `scanner.Err()`**: A broken pipe or encoding error is silently lost if you don't check `Err()` after the loop.
- **Writing diagnostics to stdout**: `fmt.Println("error:", err)` mixes error text into piped data. Always use `fmt.Fprintln(os.Stderr, ...)`.
- **Blocking on stdin in non-interactive mode**: If stdin is a terminal, the program hangs waiting for input. Use pipe detection to show a prompt only in interactive mode.

## Debugging walkthrough

Buggy code:

```go
func main() {
	data, _ := os.ReadFile("/dev/stdin") // Not portable!
	fmt.Println(string(data))
}
```

**Symptom**: Works on Linux but fails on Windows with "The system cannot find the file specified".

**Investigation**: On Windows, there is no `/dev/stdin`. The portable approach uses `os.Stdin` directly.

**Fix**: Use `io.ReadAll(os.Stdin)` or `bufio.NewScanner(os.Stdin)`.

## Production notes

- Always close pipes: a program that reads stdin in a loop should handle EOF gracefully.
- Use `os.Stderr` for all logging and diagnostics in CLI tools. This lets users run `./tool --verbose 2>debug.log`.
- For production services, stdin/stdout are often not used — logs go to files or structured logging systems. But the stream model still applies.
- `bufio.Scanner` has a max line length of 64KB by default. For longer lines, use `scanner.Buffer(buf, max)`.

## Performance implications

- `bufio.Scanner` is efficient for line-oriented I/O — it buffers reads internally, reducing system calls.
- Reading byte-by-byte from `os.Stdin` without buffering is extremely slow (one syscall per byte).
- Writing to stdout with `fmt.Fprintf` is buffered when stdout is attached to a terminal; for pipes it may be unbuffered. Use `bufio.NewWriter` for high-throughput output.

## Practice task

Write a program that reads lines from stdin, counts the number of words per line (using `strings.Fields`), and prints `"<line_number>: <word_count> words — <original_line>"`. Handle errors from scanner. Test by piping `echo` input. Save as `main.go`.

## Tests / verification

```bash
echo -e "hello world\nfoo bar baz" | go run ./curriculum/modules/07-cli-files-json-config/lessons/04-standard-input-and-standard-output
go test ./curriculum/modules/07-cli-files-json-config/lessons/04-standard-input-and-standard-output
```

## Review questions

1. Why is `bufio.Scanner` preferred over reading `os.Stdin` one byte at a time?
2. How do you check whether a program's stdin is a pipe or a terminal?
3. What is the difference between `fmt.Println` and `fmt.Fprintln(os.Stderr, ...)`?
4. What does `scanner.Scan()` return when the pipe is closed?
5. How would you read the entire stdin content as a single string?

## NEXT UP

`io.Reader` — the fundamental interface for reading data in Go.
