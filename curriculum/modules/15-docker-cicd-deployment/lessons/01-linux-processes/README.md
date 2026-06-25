# Linux processes

## Learning objective

Use Go's `os/exec` package to spawn child processes, capture their stdout/stderr, inspect exit codes, and manage process lifecycle without creating zombies or leaking resources.

## Why this matters

Every Docker container is a Linux process with an isolated view of the world. When you deploy a Go service, it runs as PID 1 inside the container. Understanding parent-child process relationships, exit codes, and zombie reaping is essential for writing reliable containerized services. Go's `os/exec` is the primary tool for running shell commands, database migrations, or sidecar processes from within a Go program.

## Mental model

A process is a running program with its own address space, file descriptors, and process ID (PID). When you call `exec.Command`, Go forks a child process. The parent (your Go program) and child (the command) run concurrently. The parent can wait for the child to finish, read its output, or send it signals. If the parent exits without waiting, the child becomes an orphan and is adopted by PID 1 (init). If the parent never calls `Wait`, the child's exit status is never collected and it becomes a zombie — a dead process occupying a slot in the process table.

## Core idea

The `os/exec` package provides a safe, cross-platform way to run external commands. It separates command construction (`exec.Command`) from command execution (`Run`, `Start`, `Output`, `CombinedOutput`). The key types are:

| Type | Purpose |
|---|---|
| `exec.Cmd` | Represents an external command ready to run |
| `exec.ExitError` | Returned when a command exits with a non-zero status |
| `Process` | OS-level process handle (PID, signal methods) |
| `ProcessState` | Exited process metadata (exit code, resource usage) |

## Under the hood

`exec.Command` populates a `Cmd` struct but does not fork. Calling `Start()` invokes `os.StartProcess`, which on Linux calls `syscall.ForkExec`. Fork clones the current process, then execve replaces the child's memory with the target binary. The child inherits the parent's file descriptors unless `Cmd.Stdout`, `Cmd.Stderr`, or `Cmd.Stdin` are set — in which case Go creates pipes. `Cmd.Wait()` calls `syscall.Wait4(pid, &wstatus, 0, &rusage)`, which blocks until the child changes state. The `wstatus` encodes the exit code and signal information.

## How Go uses it

Go's own toolchain uses process spawning everywhere:
- `go build` invokes the compiler and linker as subprocesses
- `go test` compiles and runs test binaries
- `go generate` runs user-specified commands
- `go vet` and `go fmt` are separate binaries called as subprocesses

Production Go services use `os/exec` for running backup scripts, invoking `git` in CI workflows, calling database migration binaries, and executing shell commands in admin APIs.

## Go example

```go
package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Command failed with exit code %d\n", exitErr.ExitCode())
			fmt.Printf("Stderr: %s\n", exitErr.Stderr)
		} else {
			fmt.Printf("Failed to run command: %v\n", err)
		}
		return
	}
	fmt.Printf("Output: %s", output)
	fmt.Printf("Exit code: %d\n", cmd.ProcessState.ExitCode())
}
```

## Step-by-step execution

1. `exec.Command("go", "version")` creates a `Cmd` struct. No process is spawned yet.
2. `cmd.Output()` calls `cmd.Start()` which forks the `go` binary as a child process.
3. The child process writes its stdout to a pipe created by Go.
4. `cmd.Output()` calls `cmd.Wait()` which blocks until the child exits.
5. `syscall.Wait4` returns, and Go reads the exit status from `wstatus`.
6. If the child exited with code 0, `Output()` returns the stdout bytes and nil error.
7. If the child exited non-zero, `Output()` returns an `*exec.ExitError` containing the exit code and stderr.

## Common mistakes

- **Not checking exit codes.** `exec.Command("false").Run()` returns an `ExitError` because `false` exits 1. Many developers ignore the error return and assume success. Always check the error from `Run`, `Wait`, or `Output`.
- **Creating zombie processes.** Calling `cmd.Start()` and forgetting `cmd.Wait()` leaves the child as a zombie. The OS keeps the process table entry until `Wait` collects the exit status.
- **Using `exec.Command("sh", "-c", ...)` instead of the direct binary.** This wraps your command in a shell, which changes signal propagation and adds quoting complexity. Prefer calling the binary directly with arguments.
- **Assuming `FindProcess` verifies the process exists.** `os.FindProcess` on Unix returns a `Process` even if the PID does not exist. It only checks that PID > 0. Use `process.Signal(syscall.Signal(0))` to test if a process is alive.
- **Not setting a timeout.** A child process that hangs forever will leak goroutines. Use `cmd.Context` or `cmd.Cancel` to enforce a deadline.

## Debugging walkthrough

Consider code that runs a build script and always exits with 0 even when the build fails:

```go
cmd := exec.Command("./build.sh")
cmd.Run()
fmt.Println("done")
```

**Symptom:** The program always prints "done" even when `build.sh` fails.

**Investigation:** Add error checking and stderr capture:

```go
cmd := exec.Command("./build.sh")
output, err := cmd.CombinedOutput()
if err != nil {
    fmt.Printf("Error: %v\nOutput: %s\n", err, output)
}
```

**Root cause:** `cmd.Run()` returns an `*exec.ExitError` when the command exits non-zero, but the code ignores the return value. The program continues as if nothing went wrong.

**Fix:** Always check the error from `Run`. Use `cmd.Output()` or `cmd.CombinedOutput()` when you need the command's output, and inspect `ExitCode()` from the error.

## Production notes

- Always set a timeout or context on external commands to prevent hung processes from leaking goroutines.
- Use `exec.CommandContext` with a `context.WithTimeout` to enforce deadlines.
- In Kubernetes, PID 1 inside the container should handle reaping orphaned zombie processes. Go's runtime does this automatically for child processes created via `os/exec`.
- Never pipe user input directly into `exec.Command` without sanitization — command injection is a critical vulnerability.
- Log the command being run (with sanitized arguments) for debugging, but never log secret values like passwords or tokens.

## Performance implications

- Process creation is expensive: fork+exec takes hundreds of microseconds. Minimize subprocess calls in hot paths.
- Each pipe created by `Cmd.StdoutPipe` is a kernel buffer (typically 64 KB on Linux). Reading output incrementally with `cmd.StdoutPipe` instead of `cmd.Output()` avoids blocking if the child writes more than the pipe buffer.
- Go reuses `Clone` goroutines internally but each `Start` call creates a real OS process — not a lightweight goroutine.
- The `CombinedOutput` method reads both stdout and stderr into memory; for large outputs, use separate pipes or write to files.

## Practice task

Write a function `RunWithTimeout(name string, args []string, timeout time.Duration) (string, error)` that:
- Runs the command identified by `name` with `args`.
- Returns the combined stdout+stderr as a string.
- Returns an error if the command exceeds `timeout`.
- Returns an error if the command exits non-zero (include the exit code in the error message).

Then write a `main()` that calls this with `go version` (2s timeout) and `sleep 10` (100ms timeout to trigger timeout). Run it and verify the output.

## Tests / verification

```bash
go run ./curriculum/modules/15-docker-cicd-deployment/lessons/01-linux-processes
go test -v ./curriculum/modules/15-docker-cicd-deployment/lessons/01-linux-processes
```

## Review questions

1. What is the difference between `cmd.Run()` and `cmd.Start()` followed by `cmd.Wait()`?
2. Why does `os.FindProcess(pid)` return a non-nil `Process` even when the PID does not exist?
3. What happens to a child process if the parent exits without calling `Wait()`?
4. How does `exec.CommandContext` prevent resource leaks from hung subprocesses?
5. When would you use `cmd.Output()` vs `cmd.CombinedOutput()` vs reading from `cmd.StdoutPipe`?

## NEXT UP

Signals — how the OS interrupts processes and how Go handles them with `os/signal`.
