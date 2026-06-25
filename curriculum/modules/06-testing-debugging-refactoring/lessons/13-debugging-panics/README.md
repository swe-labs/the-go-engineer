# Debugging panics

## Learning objective

Analyse Go panic traces for nil pointer dereferences, index-out-of-range errors, and deadlocks, and use `runtime.Stack` and `recover` to inspect and handle panics in production.

## Why this matters

Panics are the most common runtime failures in Go. Nil pointer dereferences, slice index overflows, and closed-channel sends account for the majority of production crashes. Every Go engineer must recognise these patterns at a glance and know exactly which line to fix. Debugging panics is not about avoiding them entirely; it is about fixing them in minutes instead of hours.

## Mental model

A panic is an interrupt that unwinds the stack until a `recover()` is found — or until the program exits. Think of it as an uncaught exception that carries a value (the panic argument) and a snapshot of the stack at the crash site. The three most common panics each have a distinct signature:

| Panic message | Root cause |
|---|---|
| `invalid memory address or nil pointer dereference` | Accessing field/method on a nil pointer or nil interface value |
| `index out of range [N] with length L` | Slice/array access with index >= length |
| `send on closed channel` | Sending to a channel that has been closed |
| `fatal error: all goroutines are asleep - deadlock!` | Every goroutine is blocked waiting on a channel, mutex, or select |

## Core idea

Every panic trace tells you **exactly** what went wrong and **exactly** where. The skill is not guessing — it is reading the trace with precision.

**Nil pointer dereference**: Look at the line in the stack trace top frame. Find the pointer access — either `.` field access, `[i]` index, or `()` method call on a nil value. The nil variable is the one before the dot.

**Index out of range**: The trace includes the invalid index `[N]` and the length `L`. Check the slice/array length and the index calculation. The index is often off-by-one or negative.

**Deadlock**: The trace shows every goroutine's stack. Each goroutine is blocked in `chan send`, `chan receive`, `sync.Mutex.Lock`, or similar. Find the cyclic wait.

## Under the hood

When the runtime detects a nil pointer access (via the CPU's memory protection), it delivers a signal (SIGSEGV on Unix). The Go runtime's signal handler translates this into a panic with the message `invalid memory address or nil pointer dereference`. For index-out-of-range, the compiler inserts a bounds check before every slice/array access; if it fails, the runtime panics immediately. Deadlock detection is a periodic check in the scheduler: if all goroutines are in a "waiting" state and no goroutine can make progress, the runtime prints all stacks and exits.

## How Go uses it

- **Recovery middleware**: HTTP frameworks (gin, echo) use deferred `recover()` in middleware to catch panics per-request, log the stack, and return HTTP 500 without crashing the server.
- **Startup validation**: Panic in `init()` or `main()` if configuration is invalid — fail fast, show the trace.
- **Testing**: `t.Helper()` marks a function as a test helper so panic traces skip it, pointing directly to the failing test line.
- **`runtime/debug.Stack`**: Captures the current goroutine's stack as a `[]byte` without requiring a buffer.

## Go example

```go
package main

import (
	"fmt"
	"runtime/debug"
)

type Server struct {
	DB *Database
}

type Database struct {
	DSN string
}

func main() {
	s := &Server{DB: nil}
	// Call without nil check — panics.
	_ = s.GetConfig()
}

func (s *Server) GetConfig() string {
	stack := debug.Stack()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIC: %v\n%s\n", r, stack)
		}
	}()
	// s.DB is nil — panic on next line.
	return s.DB.DSN
}
```

## Step-by-step execution

For `s.DB.DSN` when `s.DB` is nil:

1. Evaluate `s` → pointer to `Server`.
2. Evaluate `s.DB` → dereference `s` to get `DB` field → value is `nil` (a nil `*Database`).
3. Evaluate `s.DB.DSN` → the runtime dereferences the nil pointer to access `DSN` → **SIGSEGV**.
4. Go's signal handler catches SIGSEGV and creates a panic with `invalid memory address or nil pointer dereference`.
5. The deferred function catches the panic via `recover()`, prints the stack (captured before the panic) and the message.
6. The panic is handled. The program continues.

For an index-out-of-range: `arr := [3]int{}; _ = arr[5]`:

1. Compiler emits bounds check: is `5 < 3`? No.
2. Runtime panics: `index out of range [5] with length 3`.
3. Stack trace points to the exact line.

## Common mistakes

- **Fixing the symptom, not the cause**: Seeing a nil pointer panic and adding a nil check at the crash site without finding why the pointer is nil upstream.
- **Ignoring the full trace in deadlocks**: Only looking at one goroutine. Deadlocks always involve at least two goroutines. Read all goroutine stacks.
- **Recovering and swallowing**: `recover()` without logging the stack hides the bug. Always log the full stack trace on recovery.
- **Nil pointer in an interface value**: An interface value that holds a nil concrete pointer is not itself nil. `var p *int; var i interface{} = p; i == nil` is `false`. This causes baffling nil panics.
- **Panic in a goroutine**: A panic in a goroutine crashes the whole program, not just the goroutine. `recover()` must be in the same goroutine.

## Debugging walkthrough

Given this code:

```go
package main

import "fmt"

func main() {
	grades := []int{85, 90, 78}
	average := computeAverage(grades)
	fmt.Println("Average:", average)
}

func computeAverage(nums []int) float64 {
	total := 0
	for i := 0; i <= len(nums); i++ {
		total += nums[i]
	}
	return float64(total) / float64(len(nums))
}
```

**Symptom**: Panic: `index out of range [3] with length 3`.

**Stack trace**:
```
goroutine 1 [running]:
main.computeAverage(...)
    /tmp/main.go:11
main.main(...)
    /tmp/main.go:5
```

**Investigation**: Line 11 is `total += nums[i]`. The trace says index `[3]` with length `3`. Valid indices are 0, 1, 2. The loop condition is `i <= len(nums)` which goes one past the end. It should be `i < len(nums)`.

**Fix**: Change `i <= len(nums)` to `i < len(nums)`.

## Production notes

- **Always recover in top-level goroutines**: In a server, every HTTP handler runs in its own goroutine. Install a recovery middleware that logs the stack and continues.
- **Crash-only software**: Some teams advocate crashing on unexpected panics so the process manager (Kubernetes, supervisor) restarts the process. In this model, `recover()` is only used at the top level to log the crash before exiting.
- **Use `net/http/pprof`**: It provides `/debug/pprof/goroutine?debug=2` to dump all goroutine stacks on demand — helpful for diagnosing deadlocks without restarting.
- **Alert on recovered panics**: If you recover, treat it as a critical alert. A recovered panic is still a bug that must be fixed.

## Performance implications

- **Bounds checks**: The compiler omits bounds checks when it can prove the index is safe (e.g. `for i := range slice` has no check). In hot loops where bounds checks are unavoidable (random index access, dynamic indexing), the CPU branch misprediction cost is ~1–2 cycles — negligible for most code.
- **Recover overhead**: A deferred `recover()` is fast if no panic occurs (~10ns). If a panic does occur, stack unwinding is O(depth) and includes printing the trace. Panics are not for control flow; they are for exceptional conditions.
- **Signal handling**: Nil pointer panics go through the OS signal handler (SIGSEGV), which is slower than a software-triggered panic (e.g. index out of bounds). Both are fatal unless recovered.

## Practice task

Write a function `SafeDivide(a, b []int) []int` that divides each element of `a` by the corresponding element of `b`. If `b[i]` is 0, recover the resulting panic (integer division by zero panics?) — actually, integer division by zero panics in Go. Recover from it and append `0` to the result. If `i` is out of range for `b`, treat it as division by 1 (no panic, just use 1). Write a test that exercises both error paths.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/13-debugging-panics
go test ./curriculum/modules/06-testing-debugging-refactoring/lessons/13-debugging-panics
```

## Review questions

1. What does `invalid memory address or nil pointer dereference` mean in terms of what the CPU experienced?
2. A slice has length 5. What index values are valid? What does the panic message look like for index 5?
3. How can a `recover()` call in one goroutine catch a panic from another goroutine?
4. What is the difference between a nil interface value and an interface holding a nil concrete pointer?
5. In a deadlock trace, what should you look for across all goroutine dumps?

## NEXT UP

Delve basics — using a debugger to step through Go code interactively.
