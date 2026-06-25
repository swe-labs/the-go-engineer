# Delve basics

## Learning objective

Install Delve, start a debugging session with `dlv debug` and `dlv test`, and use basic CLI commands to inspect program execution.

## Why this matters

Print statement debugging (`fmt.Println`) works for simple programs, but it breaks down when you need to inspect intermediate state, follow complex control flow, or debug concurrent code. A debugger like Delve lets you pause execution at any line, examine variables, and step through code instruction by instruction. Every professional Go engineer should know how to use a debugger.

## Mental model

Think of Delve as a remote control for your program. The debugger starts your program in a paused state. You place markers (breakpoints) where you want to pause. When the program hits a marker, it stops, and you can inspect the current state, then tell it to continue (run to the next marker) or step forward one line at a time. Delve controls one goroutine at a time, so you can follow exactly what that goroutine is doing.

## Core idea

Delve is Go's official debugger. It understands Go's runtime, goroutines, and data structures natively. Key commands:

| Command | Short | Action |
|---|---|---|
| `break` | `b` | Set a breakpoint on a file:line or function |
| `continue` | `c` | Resume execution until next breakpoint |
| `next` | `n` | Step over to next line in current function |
| `step` | `s` | Step into a function call |
| `stepout` | `so` | Step out of current function to caller |
| `print` | `p` | Print a variable or expression |
| `locals` | — | Print local variables in current frame |
| `args` | — | Print function arguments in current frame |
| `stack` | `bt` | Print the current goroutine's stack trace |
| `goroutines` | `grs` | List all goroutines |
| `restart` | `r` | Restart the program from the beginning |

## Under the hood

Delve starts your compiled binary as a child process but pauses it immediately before `main`. It uses OS tracing facilities to control execution:
- **Linux**: `ptrace` syscall
- **macOS**: `mach` traps via `task_for_pid`
- **Windows**: `WaitForDebugEvent` and `ContinueDebugEvent`

Delve inserts `INT3` breakpoint instructions (0xCC on x86) at target addresses. When the CPU hits an `INT3`, it raises a signal that Delve intercepts. Delve then looks up the address in its breakpoint table, stops the process, and returns control to the user. On `continue`, Delve restores the original instruction, single-steps past it, re-inserts the breakpoint, and resumes.

## How Go uses it

- **`dlv debug`**: Compiles and debugs the current package. Equivalent to `go build` then `dlv exec`.
- **`dlv test`**: Works like `dlv debug` but for tests. You can set breakpoints in test functions.
- **`dlv exec`**: Attaches to an already-compiled binary.
- **`dlv attach`**: Attaches to a running process by PID.
- **`dlv connect`**: Connects to a headless Delve server (for remote debugging).
- **`dlv core`**: Examines a core dump without running the program.

## Go example

```go
package main

import "fmt"

func main() {
	msg := greeting("Alice")
	fmt.Println(msg)

	result := compute(10, 5)
	fmt.Println("Result:", result)
}

func greeting(name string) string {
	return "Hello, " + name + "!"
}

func compute(a, b int) int {
	sum := a + b
	diff := a - b
	product := a * b
	return sum + diff + product
}
```

To debug this program:

```bash
dlv debug
(dlv) break main.main
(dlv) continue
(dlv) break main.compute
(dlv) continue
(dlv) print a
(dlv) print b
(dlv) next
(dlv) print sum
(dlv) continue
```

## Step-by-step execution

Starting a `dlv debug` session:

1. Run `dlv debug` in the lesson directory. Delve compiles the package and pauses at the first instruction (in `runtime` setup, before `main`).
2. Set a breakpoint: `break main.main` — Delve resolves the symbol and sets a breakpoint.
3. `continue` — execution runs to `main.main` and stops.
4. `break main.compute` — set another breakpoint inside `compute`.
5. `continue` — execution runs to the `compute(10, 5)` call and stops at the first line of `compute`.
6. `print a` — shows `10`. `print b` — shows `5`.
7. `next` — executes `sum := a + b`, stops before `diff := a - b`.
8. `print sum` — shows `15`.
9. `continue` — runs to the end of the program and exits.

## Common mistakes

- **Forgetting to set breakpoints before `continue`**: If no breakpoint is set, the program runs to completion and exits.
- **Using `next` when you meant `step`**: `next` steps over function calls; `step` steps into them. If you want to trace into `compute`, use `step` at the call site.
- **Not using tab completion**: Delve supports tab completion for commands, file names, and symbols. Press Tab liberally.
- **Trying to attach to a stripped binary**: Without debug symbols, Delve cannot resolve function names and line numbers. Compile without `-ldflags="-s -w"` for debugging.
- **Debugging in optimized code**: The compiler may inline functions or reorder statements. Compile with `-gcflags="-N -l"` to disable optimisation when debugging (`dlv debug` does this automatically).

## Debugging walkthrough

Start debugging our example:

```bash
cd path/to/lesson
dlv debug
```

```
(dlv) break main:12
Breakpoint 1 set at 0x...
(dlv) continue
> main.main() ./main.go:12 (hits goroutine(1):1 total:1)
    12:     result := compute(10, 5)
(dlv) break main.compute
Breakpoint 2 set at 0x...
(dlv) continue
> main.compute() ./main.go:17 (hits goroutine(1):1 total:1)
    17: func compute(a, b int) int {
(dlv) print a
10
(dlv) print b
5
(dlv) next
> main.compute() ./main.go:18 (hits goroutine(1):1 total:1)
    18:     sum := a + b
(dlv) next
    19:     diff := a - b
(dlv) print sum
15
(dlv) continue
Hello, Alice!
Result: 30
```

## Production notes

- **Headless mode**: Run `dlv debug --headless --listen=:2345 --api-version=2` to start a remote debugging server. Connect with `dlv connect :2345` or from VS Code / GoLand.
- **Never debug in production**: Debugging requires pausing execution. Use Delve on a staging environment or a local replica. For production issues, rely on logging, metrics, and core dumps.
- **Core dump debugging**: Capture a core dump (`GOTRACEBACK=crash ./app`) and analyse offline with `dlv core <binary> <core>`. This explores the state at the crash without any execution.
- **CI integration**: `dlv test` with `--headless` can be used in CI to collect detailed failure information from failing tests.

## Performance implications

- Delve inserts a trap instruction (INT3) at each breakpoint, which causes a context switch to the OS debugger handler — ~1 microsecond overhead per breakpoint hit.
- With 0 breakpoints set, Delve adds zero runtime overhead. There is no performance cost to running under Delve unless you stop at breakpoints.
- Stepping through code (one line at a time) is slow — hundreds of microseconds per step. It is meant for detailed inspection, not for running through large loops.
- Conditional breakpoints (when you are ready for lesson 15) require Delve to evaluate the condition at every hit, which adds overhead proportional to the condition complexity.

## Practice task

Our `main.go` has a function `compute(a, b int) int`. Start `dlv debug`, set a breakpoint on `main.compute`, and step through each line. Use `print` to inspect `a`, `b`, `sum`, `diff`, and `product` after each `next`. Use `locals` to see all local variables at once. Exit with `quit`.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/14-delve-basics
go test ./curriculum/modules/06-testing-debugging-refactoring/lessons/14-delve-basics
```

## Review questions

1. What is the difference between `dlv debug` and `dlv exec`?
2. What command resumes execution after a breakpoint hit?
3. If you want to follow execution into a called function, do you use `next` or `step`?
4. What does `INT3` have to do with how Delve implements breakpoints?
5. Why should you not attach Delve to a production process?

## NEXT UP

Breakpoints — how to set, list, clear, and use conditional breakpoints in Delve.
