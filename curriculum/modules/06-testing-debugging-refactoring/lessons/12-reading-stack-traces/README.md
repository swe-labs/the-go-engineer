# Reading stack traces

## Learning objective

Read and interpret Go stack traces to locate the source of a panic or error, and use `runtime.Caller` and `runtime.Stack` to inspect the call stack programmatically.

## Why this matters

When a Go program crashes in production, the stack trace is your first and best clue. It tells you exactly where the panic happened, which function called it, and the entire call chain back to `main` or a goroutine entry point. Engineers who cannot read stack traces waste hours grepping logs and guessing. Those who can read them identify the root cause in seconds.

## Mental model

A stack trace is a top-to-bottom timeline of function calls, written backwards. The top frame is the crash site — the line that panicked. Each frame below shows the function that called the one above, its file, and its line number. Frames are like footprints: follow them from the crash back to the starting point. Each goroutine gets its own stack trace section prefixed with `goroutine N [state]:`.

```
panic: runtime error: invalid memory address or nil pointer dereference
[signal 0xc0000005 code=0x0 addr=0x0 pc=0x...]

goroutine 1 [running]:
main.buggyFunc(...)
    /home/user/proj/main.go:42
main.callChain(...)
    /home/user/proj/main.go:38
main.main()
    /home/user/proj/main.go:10
```

Read from top: `buggyFunc` at line 42 panicked, called by `callChain` at line 38, called by `main` at line 10.

## Core idea

Every stack trace contains:

| Component | Meaning |
|---|---|
| `goroutine N [state]:` | Goroutine ID and current state (running, sleeping, IO wait) |
| `pkg.FuncName(...)` | Fully qualified function name, with argument ellipsis |
| `file.go:line` | Source file path and line number |
| `+0xNN` | Offset within the function's machine code (rarely needed) |

A **panic trace** includes the panic message and value at the top. An **error trace** (from `errors.New` or `fmt.Errorf`) includes no automatic stack — you must add one with `runtime.Caller` or use a library like `pkg/errors`.

## Under the hood

Go's runtime maintains a stack of _call frames_ for every goroutine. Each frame holds the program counter (PC) and stack pointer. When a panic occurs, the runtime walks the linked list of frames using the PC to look up function metadata in the symbol table (embedded in the binary unless stripped). The stack walker stops when it reaches `main.main()`, the goroutine entry point, or `goexit`.

`runtime.Caller(skip)` walks up the stack by `skip` frames (0 = the function calling Caller). `runtime.Stack(buf, all)` writes formatted stack traces into a byte slice. Both rely on the same internal frame-walking mechanism.

## How Go uses it

- **Panic recovery**: `recover()` inspects the current goroutine's stack for a pending panic.
- **Logging**: Production servers capture `runtime.Stack` on `SIGQUIT` or on unexpected panics and write it to a crash file.
- **Profiling**: `pprof` uses stack traces to attribute CPU and memory samples to call sites.
- **Testing**: `t.Log` and `testing/quick` include stack traces on failure.
- **Error wrapping**: `fmt.Errorf("context: %w", err)` preserves the original error but not the stack; use `runtime.Caller` to annotate.

## Go example

```go
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("=== Walk the stack with runtime.Caller ===")
	First()

	fmt.Println("\n=== Capture full trace with runtime.Stack ===")
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	fmt.Printf("%s\n", buf[:n])

	fmt.Println("\n=== Panic and recover ===")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered:", r)
			}
		}()
		trigger()
	}()
}

func First()  { Second() }
func Second() { Third() }
func Third() {
	for skip := 0; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		fmt.Printf("  skip=%d  %s (%s:%d)\n", skip, fn.Name(), file, line)
	}
}

func trigger() { deep() }
func deep()    { panic("boom") }
```

## Step-by-step execution

For the `Third()` call chain `main → First → Second → Third`:

1. `Third()` is called. The stack has four frames.
2. Loop starts at `skip=0`: `runtime.Caller(0)` returns `main.Third`.
3. `skip=1`: returns `main.Second` (the caller of Third).
4. `skip=2`: returns `main.First` (caller of Second).
5. `skip=3`: returns `main.main` (caller of First).
6. `skip=4`: `ok=false`, loop stops.

For `runtime.Stack(buf, false)`:
1. The runtime walks the calling goroutine's frame list.
2. Each frame is formatted as `file:line`, function name, offset.
3. The written bytes are printed. Only the current goroutine is shown (`all=false`).

For the panic in `deep()`:
1. Execution stops at `panic("boom")`.
2. Stack unwinding begins.
3. `recover()` in the deferred function catches the panic value.
4. The panic does not propagate to `main`. Program continues.

## Common mistakes

- **Reading bottom-to-top**: The crash is at the **top**, not the bottom. Beginners often start at `main()` and get confused.
- **Ignoring file:line numbers**: The function name tells you where, the line number tells you exactly. Always look at the line.
- **Confusing goroutine traces**: If you see multiple `goroutine N` sections, each is independent. Find the one with `[running]` — that's where the panic is.
- **Assuming all panics produce traces**: Some runtime panics (e.g. stack overflow) produce truncated traces. Silent nil panics in `defer` may swallow the original trace.
- **Not capturing stack on error**: `fmt.Errorf("wrapped: %w", err)` carries no stack. Use `runtime.Caller` at the error site if you need the call site.

## Debugging walkthrough

Given this broken code:

```go
package main

type User struct {
	Name string
}

func main() {
	var u *User
	fmt.Println(u.Name)
}
```

**Symptom**: Panic: `invalid memory address or nil pointer dereference`.

**Stack trace**:
```
goroutine 1 [running]:
main.main()
    /tmp/example/main.go:8 +0x35
```

**Investigation**: The trace points to line 8 in `main()`. Line 8 is `fmt.Println(u.Name)`. Variable `u` is a nil pointer — `var u *User` initialises it to `nil`.

**Fix**: Initialise `u` before use: `u = &User{Name: "Alice"}`. Or check for nil: `if u != nil { fmt.Println(u.Name) }`.

## Production notes

- **Always log the full stack trace** on unexpected panics. Use `runtime.Stack` in a deferred `recover()` and write to a dedicated error log.
- **Strip binary paths** in release builds. Use `-trimpath` in `go build` to remove absolute file paths from stack traces, keeping only relative module paths.
- **Error monitoring tools** (Sentry, Datadog, Rollbar) parse stack traces from `runtime.Stack`. Send them as structured fields, not as log text.
- **Goroutine dumps**: Send `SIGQUIT` (Unix) to a running Go process to print all goroutine stacks to stderr. On Windows, use Ctrl-Break or `runtime.Gosched` tricks.
- **`net/http/pprof`** serves `/debug/pprof/goroutine?debug=2` which returns full goroutine stack traces — invaluable for debugging deadlocks in production.

## Performance implications

- `runtime.Caller` and `runtime.Stack` are not free. Each call walks the frame list and resolves PCs to symbols. On a hot path, avoid calling them on every iteration.
- Capturing a stack trace takes ~1–5 microseconds depending on stack depth. In a high-throughput server, sample stack traces (e.g. 1 in 1000 requests) rather than capturing every one.
- `runtime.Stack(buf, true)` (all goroutines) acquires a global lock and blocks other goroutines briefly. Use sparingly in production — typical use is on `SIGQUIT` or on severe errors only.
- Stack trace formatting allocates memory for the formatted string. Pre-allocate a buffer if capturing traces in a tight loop.

## Practice task

Write a function `StackDepth() int` that returns the current number of frames on the goroutine's stack by using `runtime.Caller` in a loop. Then write a recursive function `recurse(n int) int` that calls itself until `n == 0`, and at depth 0 calls `StackDepth()`. Print the depth at each recursion level. Run the program and confirm the depth matches the recursion count + a constant offset for main and runtime frames.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/12-reading-stack-traces
go test ./curriculum/modules/06-testing-debugging-refactoring/lessons/12-reading-stack-traces
```

## Review questions

1. Which end of a stack trace is the crash site — the top or the bottom?
2. What does the `skip` parameter of `runtime.Caller(skip)` represent?
3. How would you capture a stack trace from all goroutines, not just the current one?
4. Why does `fmt.Errorf("wrapped: %w", err)` not include a stack trace? How would you add one?
5. In a panic trace, what does the `goroutine N [running]:` header tell you?

## NEXT UP

Debugging panics — how to dissect common panic types and fix them systematically.
