# Exit codes

## Learning objective

Evaluate process exit codes (0 for success, non-zero for failure) to construct robust shell script error handling and CI pipelines.

## Why this matters

Programs need a simple, universal way to signal success or failure to their caller (the shell, a parent program, or an automated CI pipeline). Exit codes provide this with a single integer — no parsing of unpredictable console text required. Understanding exit code contracts is the foundation of automated testing, container health checks, and production deployment scripts.

## Mental model

When an operating system process terminates, it returns an unsigned 8-bit integer (0–255 on POSIX, 32-bit on Windows) to the OS kernel:
- **Exit Code 0:** The universal contract for **Success** ("completed cleanly without errors").
- **Non-Zero (1–255):** Indicates **Failure** or abnormal termination.
- **The Caller Contract:** Automated callers (Bash, PowerShell, GitHub Actions, Docker, Kubernetes) do not inspect the words printed on the screen. They inspect *only* the numeric exit code to decide whether a command succeeded or failed.

## Visual model

```mermaid
flowchart TD
    PROG["Go Application Process terminates"] --> EXIT{"Termination Path"}

    EXIT -->|main() returns normally| SUCCESS["os.Exit(0)
Kernel exit status: 0"]
    EXIT -->|os.Exit(1) or handled error| FAILURE["os.Exit(1)
Kernel exit status: 1"]
    EXIT -->|panic() occurs without recover| PANIC["runtime panic unwinds
Kernel exit status: 2"]
    EXIT -->|Killed by OS OOM / SIGKILL| KILLED["Killed by Signal 9
Kernel exit status: 137 (128+9)"]

    SUCCESS --> SHELL["Parent Shell / CI Runner inspects status"]
    FAILURE --> SHELL
    PANIC --> SHELL
    KILLED --> SHELL

    SHELL --> DECISION{"Status == 0 ?"}
    DECISION -->|Yes: 0| CI_PASS["CI Step: SUCCESS (Proceed to next step)"]
    DECISION -->|No: != 0| CI_FAIL["CI Step: FAILED (Halt pipeline immediately)"]
```

## Core idea

Without exit codes, scripts would need to parse arbitrary strings like `"error"`, `"FAILED"`, or `"done"` — a fragile approach that breaks whenever log formats change, translations are added, or casing changes. A single standardized integer creates an unambiguous, universal machine-to-machine contract.

## Standard Exit Code Reference

While applications can define custom codes between 1 and 125, the software engineering industry adheres to standard conventions:

| Exit Code | Meaning | Typical Trigger |
|:---:|---|---|
| **0** | **Success** | Program completed its task without errors (`main()` returned). |
| **1** | **General Catch-All Error** | Catch-all for application errors (file not found, invalid argument, DB connection failure). |
| **2** | **Misuse of Shell Builtin / Panic** | Missing required flags, shell syntax error, or unrecovered Go runtime panic. |
| **126** | **Command Cannot Execute** | Target file was found in `$PATH`, but lacks executable permissions (`chmod +x`). |
| **127** | **Command Not Found** | Executable does not exist in any directory listed in `$PATH`. |
| **128** | **Invalid Exit Argument** | Calling `exit` with an out-of-range argument on POSIX. |
| **128 + N** | **Fatal Error Signal `N`** | Process was forcibly terminated by OS kernel signal `N`. |
| **130** | **Terminated by `SIGINT` (128 + 2)** | User pressed `Ctrl+C` in the terminal to interrupt the process. |
| **137** | **Terminated by `SIGKILL` (128 + 9)** | Process was killed by `kill -9` or the OS **Out of Memory (OOM) Killer**. |
| **143** | **Terminated by `SIGTERM` (128 + 15)**| Graceful termination request from Docker, Kubernetes, or system shutdown. |

## Shell Evaluation and Automation

How the calling shell retrieves the exit code of the immediately preceding command:
- **Bash / Zsh:** Special variable `$?`
- **PowerShell:** Special variable `$LASTEXITCODE` (for native executables) and `$?` (boolean `True`/`False`)

### Script Chaining
- `cmd1 && cmd2`: Runs `cmd2` only if `cmd1` returns exit code `0`.
- `cmd1 || cmd2`: Runs `cmd2` only if `cmd1` returns a non-zero exit code.
- In Bash scripts, adding `set -e` at the top causes the entire script to abort immediately if any command returns a non-zero exit code.

## How Go uses it

- When `func main()` returns, the Go runtime automatically calls `os.Exit(0)`.
- To exit immediately with a specific error code, call `os.Exit(code)`. Note: `os.Exit()` does not run deferred functions (`defer`).
- An unrecovered `panic()` prints a stack trace to `stderr` and exits with code `2`.
- When running `go test`, the test runner exits with `0` if all tests pass, and `1` if any test fails.

## Go example

The example program accepts a filename argument, inspects its existence, demonstrates proper exit code signaling, and shows how unhandled panics produce standard exit codes:

```go
package main

import (
	"fmt"
	"os"
)

// CheckFileStatus validates if a file exists and returns a standardized exit code.
func CheckFileStatus(path string) int {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: file '%s' does not exist.\n", path)
			return 1 // General error: file missing
		}
		fmt.Fprintf(os.Stderr, "Error accessing file: %v\n", err)
		return 1
	}

	if info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: '%s' is a directory, not a file.\n", path)
		return 1
	}

	fmt.Printf("Success: file '%s' found (%d bytes).\n", path, info.Size())
	return 0 // Success
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: exit-codes <filename>")
		fmt.Println("Exiting with status code 2 (Invalid Arguments)")
		os.Exit(2)
	}

	target := os.Args[1]

	// Handle simulated panic demo
	if target == "--simulate-panic" {
		fmt.Println("Triggering runtime panic...")
		panic("critical unexpected failure") // Will exit with status 2
	}

	exitCode := CheckFileStatus(target)
	fmt.Printf("Process terminating with exit code: %d\n", exitCode)
	os.Exit(exitCode)
}
```

## Step-by-step execution

1. The user invokes `./exit-codes myfile.txt` from their terminal shell.
2. The operating system spawns the process and passes command-line arguments in `os.Args`.
3. If no arguments are provided, the program prints usage text to `stdout` and calls `os.Exit(2)`.
4. If a filename is supplied, `os.Stat(target)` tests whether the file exists.
5. If the file exists, the program prints confirmation and exits with `0`.
6. If the file is missing, the program prints an error message to `stderr` and exits with `1`.
7. The OS kernel captures the return code in the process's termination record.
8. The parent shell reads the code via the `wait4()` syscall and populates `$?` (or `$LASTEXITCODE`).

## Common mistakes

| Mistake | Why it happens | Fix |
|---|---|---|
| Exiting with 0 when an error occurred | The program printed an error message, but main returned normally | If an error occurs, explicitly exit with a non-zero code (`os.Exit(1)`) |
| Assuming non-zero means specifically exit code 1 | Many errors produce specific codes (2, 127, 137) | In scripts, check `$? -ne 0` rather than strictly `$? -eq 1` |
| Calling `os.Exit()` and expecting `defer` to run | `os.Exit` terminates the process immediately via kernel syscall | Return errors up to `main()` and let `main()` return, or log and exit at the top level |
| Using negative exit codes (e.g., `-1`) | POSIX converts integers modulo 256; `-1` becomes `255` | Always use clear positive integers in the range `0` to `125` for application codes |

## Debugging walkthrough

**Scenario: A CI pipeline reports a build step passed, even though the test command printed error logs.**
- **Cause:** The script executed a test command that failed, but the very next command was `echo "Tests completed"`, which succeeded with exit code 0. The shell returns the exit code of the *last* command run in the step.
- **Diagnosis:** Run `false; echo "Done"; echo $?`. The output is `0`.
- **Resolution:** In Bash CI scripts, add `set -e` (or `set -o pipefail`) at the beginning of the script so any failing command halts the script immediately.

**Scenario: A Docker container suddenly disappears with exit code 137.**
- **Cause:** Exit code `137` = `128 + 9` (`SIGKILL`). The container consumed more memory than its limit, and the Linux kernel's Out of Memory (OOM) Killer terminated it violently.
- **Diagnosis:** Run `docker inspect <container_id> --format='{{.State.OOMKilled}}'`. If `true`, the process was killed by OOM.
- **Resolution:** Increase the container memory limits or optimize your Go application's memory allocations.

## Production notes

In Kubernetes, Pod restart policies and health probes (Liveness and Readiness probes) rely strictly on container exit codes. When a container exits with non-zero, Kubernetes records the exit code in `CrashLoopBackOff` events. If the exit code is 137, SREs immediately investigate memory limits; if it is 143, they know the container was gracefully drained during a cluster autoscaling event.

## Performance implications

- Capturing an exit code is a zero-cost operation; the OS kernel already tracks process termination status in the process table.
- Exit codes are restricted to 8 bits on POSIX systems (0–255); any integer greater than 255 is truncated via modulo 256 (`256 % 256 = 0`), which can accidentally turn a failure into a success!

## Tests / verification

1. Run a command that succeeds and check its exit code:
   - Bash: `echo "hello" && echo $?` (Prints `0`)
   - PowerShell: `Write-Host "hello"; $LASTEXITCODE` (Or `$?`)
2. Run a command that fails and check its exit code:
   - Bash: `ls /nonexistent_directory 2> /dev/null; echo $?` (Prints non-zero)
   - PowerShell: `Get-ChildItem C:\nonexistent_directory -ErrorAction SilentlyContinue; $?` (Prints `False`)
3. Test conditional chaining:
   - Bash: `true && echo "Runs on success"` vs `false && echo "Will not run"`
   - Bash: `false || echo "Runs on failure"`

## Review questions

1. What is the universal meaning of exit code `0` across all operating systems?
2. What does exit code `137` indicate, and how is that number calculated from kernel signals?
3. What special shell variable stores the exit code of the previous command in Bash vs PowerShell?
4. Why is relying on exit codes superior to parsing console output text in CI/CD automation pipelines?
5. What happens to Go `defer` statements when `os.Exit(1)` is invoked directly?

## NEXT UP

[Lesson 07: How the OS manages processes](../07-how-the-os-manages-processes/README.md)
