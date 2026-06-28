# go run

## Learning objective

By the end of this lesson, you will understand what `go run` does, when to use it, and how it differs from `go build`. You will use `go run` to execute Go programs without producing a permanent binary.

## Why this matters

`go run` is the fastest way to execute a Go program during development. You write code, run it, see output, and iterate. No separate compile step, no binary to clean up. Every Go developer uses `go run` dozens of times per day during the edit-run-debug cycle.

## Mental model

`go run` is like a REPL for compiled languages — it takes your source file, compiles it in a temp directory, runs the result, and cleans up. The source file remains; the binary disappears. Think of it as "compile, execute, forget."

## Core idea

`go run` compiles a Go package and immediately executes the resulting binary, then discards it. It is the fastest feedback loop for Go development.

## Under the hood

When you run `go run main.go`, Go:

1. Creates a temporary directory.
2. Compiles the source files into a binary in that temp directory.
3. Executes the binary with stdout and stderr connected to your terminal.
4. Deletes the temp directory and binary after execution completes.

The binary name is a random hash — you never see it. The entire cycle happens in milliseconds for small programs.

## How Go uses it

The Go standard library uses `go run` extensively in its own development. The `go tool` subcommands (like `go test`, `go build`) are themselves Go programs that are run during development. The `go run` command accepts file paths, package paths, and module patterns:

- `go run main.go` — run a single file
- `go run .` — run the package in the current directory
- `go run ./cmd/server` — run a package in a subdirectory

## Go example

```go
package main

import "fmt"

func main() {
	fmt.Println("Running with go run")
	fmt.Println("This program was compiled and executed in one step.")
}
```

## Step-by-step execution

1. Run `go run .` from this lesson directory.
2. Observe the output in your terminal.
3. Modify the message string and run `go run .` again — see how fast the edit-run cycle is.
4. Run `go build .` afterward to see the difference — notice the binary that remains on disk.

## Common mistakes

- **Using `go run` for production**: `go run` produces no permanent binary. For production, use `go build` and deploy the resulting executable.
- **Running `go run main.go` when there are multiple files**: If the directory has multiple `.go` files, `go run main.go` only compiles that one file. If `main.go` references functions in other files, compilation fails. Use `go run .` instead.
- **Forgetting that `go run` recompiles every time**: Unlike `go test` which caches results, `go run` always recompiles. For large programs, this adds startup time. Use `go build` + run the binary for repeated executions.

## Debugging walkthrough

**Scenario**: `go run .` produces `package main is not a main package`.

1. Check that one file in the directory has `func main()`. You can only have one `main` function per package.
2. Check that the file with `func main()` is in `package main` (not `package foo`).
3. Run `go run main.go` explicitly if you know which file contains `main`.

## Production notes

- Never use `go run` in CI/CD pipelines. Always use `go build` and run the binary explicitly.
- `go run` is ideal for small CLI tools, one-off scripts, and exploratory coding.
- The `-race` flag works with `go run`: `go run -race .` compiles and runs with the race detector enabled.

## Performance implications

- `go run` recompiles every time — no cache reuse. For a tiny program this is instant; for a large module it adds 1-3 seconds.
- The temporary binary is compiled with debug symbols and no optimizations (like `go build` without flags).
- `go run -count=1 .` is redundant — `go run` never caches.

## Practice task

1. Run `go run .` from this lesson directory.
2. Change the message to something personal and run `go run .` again.
3. Time how long it takes from pressing Enter to seeing output.

## Tests / verification

```bash
go run .
```

Expected output:
```
Running with go run
This program was compiled and executed in one step.
```

After your modification, the output should reflect your new message.

## Review questions

1. What happens to the binary after `go run` finishes?
2. Why should you use `go build` instead of `go run` for production?
3. What is the difference between `go run main.go` and `go run .`?
4. Does `go run` cache compiled results like `go test`?
5. What flag would you add to `go run` to enable the race detector?

## NEXT UP

[go build](../04-go-build/README.md) — compile Go programs into standalone binaries.
