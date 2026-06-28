# Exit codes

## Learning objective

Understand and apply Exit codes in the context of professional Go software engineering.

## Why this matters

Programs need a simple, universal way to signal success or failure to their caller (shell, CI system, or another program). Exit codes provide this with a single integer — no parsing of output text required.

## Mental model

Compiler errors tell you exactly what is wrong and where. Read them from bottom to top: the last line is the root cause. Learning to read errors without panic is the highest-leverage skill for a new Go developer.

## Core idea

Without exit codes, scripts and shells would need to parse program output to determine success — fragile, language-dependent, and slow. A single integer provides reliable, universal success/failure signaling.

## Under the hood

On Linux, the exit syscall (nr 60 on x86_64) takes a single integer argument. The kernel sets the process's exit_code field in task_struct, wakes up the parent waiting on waitid(), and stores the code until the parent collects it with wait4() (zombie state). On Windows, ExitProcess sets the exit code in the EPROCESS block.

## How Go uses it

Go programs use os.Exit(code) for explicit exit codes and panic() for unrecoverable errors (which exits with code 2). The 'go vet' tool reports when os.Exit is called in unexpected places. Go tests use t.FailNow() and t.Fatal() to signal failure with appropriate exit codes.

## Go example

```go
package main

import (
	"fmt"
	"os"
)

// ExitCodeByFileStatus returns 0 for success (file exists) or 1 for failure.
func ExitCodeByFileStatus(fileExists bool) int {
	if fileExists {
		return 0
	}
	return 1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: exit-codes <filename>")
		os.Exit(2)
	}
	_, err := os.Stat(os.Args[1])
	exists := err == nil
	code := ExitCodeByFileStatus(exists)
	fmt.Printf("exit code: %d\n", code)
	os.Exit(code)
}
```

## Step-by-step execution

1. A Go program completes normally — the runtime calls os.Exit(0) automatically when main returns.
2. If the program encounters an error, it calls os.Exit(1) (or another code) to signal failure.
3. The OS kernel stores the exit code in the process's exit status field in the process table.
4. The parent shell retrieves the exit code using waitpid() on Unix or GetExitCodeProcess on Windows.
5. The shell stores the exit code in $? (bash) or $LASTEXITCODE (PowerShell) for use in conditionals.

## Common mistakes

| Mistake | Why it happens | Fix |
|---------|---------------|-----|
| Ignoring exit codes entirely and assuming a program succeeded because no error text was printed | Developers often only check printed output, not the exit code | Always check the exit code after running a command: echo $? or $LASTEXITCODE |
| Using exit code 0 for errors and non-zero for success | The convention is 0 = success, non-zero = failure; reversing it breaks all tooling | Always use 0 for success and non-zero for failure |
| Forgetting that a program exiting with code 0 may still have produced incorrect output | Exit code only indicates the program ran without internal errors, not that the output is correct | Validate output separately using tests or assertions |

## Debugging walkthrough

### Scenario: A CI pipeline step fails but the script continues running, masking the failure.

**Cause:** The script did not check the exit code of the failed command and did not use set -e (bash) or $? checks.

**Fix:** Add error checking after each command: `command || exit 1` in bash, or use `set -e` to exit on any failure.

### Scenario: A program always exits with code 0 even when it encounters an unrecoverable error.

**Cause:** The program does not call os.Exit(1) or return a non-zero code on error paths.

**Fix:** Always return a non-zero exit code on error: os.Exit(1) for generic errors, or os.Exit(64) for usage errors.

## Production notes

CI/CD pipelines, shell scripts, Makefiles, and container orchestrators all rely on exit codes to determine whether a step succeeded. 'docker run' exits with the container's exit code. 'kubectl apply' uses exit codes for success/failure signaling.

## Performance implications

- Reading an exit code is a zero-cost operation — the kernel already captures it during process termination.
- Exit codes beyond 255 are truncated (modulo 256) on Unix — always use 0-255 range.
- Exit code 128+N on Unix conventionally indicates termination by signal N (e.g., 130 = SIGINT, 137 = SIGKILL).

## Tests / verification

1. Run `true` in your terminal, then check the exit code with `echo $?` (bash) or `$LASTEXITCODE` (PowerShell). It should be 0. Repeat with `false` — it should be 1.
2. Chain commands: `true && echo "ran"` prints "ran". `false && echo "ran"` does not. Try the same with `||`.
3. Run a command that fails (e.g., `ls nonexistent`), then check the exit code. Note how different failures produce different non-zero codes.

Apply this in the [Shell + Git Lab](../projects/shell-git-lab/README.md): the verify script checks exit codes from every Git command. A failed merge or missing file produces a non-zero exit, and the script reports it as FAIL. Understanding exit codes helps you read and debug verification output.

## Review questions

1. What exit code does a successful command return? What range indicates failure?
2. Explain the difference between a panic (non-zero exit + stack trace) and os.Exit (clean exit with code).
3. Chain two commands with `&&` and `||` and predict which runs based on exit codes.
4. What exit code does a program return by default when it completes normally? Why is 0 used for success?
5. Why does `git merge` produce a non-zero exit code when there is a merge conflict, and how should a script handle this?

## NEXT UP

Lesson 07: [How the OS manages processes](../07-how-the-os-manages-processes/README.md)
