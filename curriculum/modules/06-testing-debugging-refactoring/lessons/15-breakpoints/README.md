# Breakpoints

## Learning objective

Set, list, clear, and use conditional breakpoints in Delve to pause execution at specific functions, lines, or when a condition is true.

## Why this matters

Without breakpoints, debugging is a firehose: the program runs from start to finish, and you must infer everything from output logs. Breakpoints let you freeze time at exactly the moment something interesting happens — when a specific function is called, when a variable reaches a certain value, or when a particular line is hit. This is the fundamental building block of all interactive debugging.

## Mental model

A breakpoint is a bookmark that says "stop here". When execution reaches a bookmarked line, the debugger pauses and hands you control. You can have multiple bookmarks. Each one can be:
- **Location-based**: stops at a file:line or function
- **Conditional**: stops only if a boolean expression is true
- **Temporary**: stops once and is automatically removed

Think of breakpoints as selective logging without having to edit, recompile, and redeploy code.

## Core idea

Delve breakpoint commands:

| Command | Example | Effect |
|---|---|---|
| `break <loc>` | `b main.compute` | Break on function entry |
| `break <file>:<line>` | `b main.go:18` | Break at a specific line |
| `break <loc> if <cond>` | `b main.compute if a > 100` | Conditional break |
| `break -<type> <loc>` | — | See `help break` for types |
| `breakpoints` | `bp` | List all breakpoints |
| `clear <n>` | `clear 1` | Remove breakpoint by ID |
| `clearall` | — | Remove all breakpoints |
| `on <n> <cmd>` | `on 1 p a` | Execute command each hit without stopping |

Location specifiers:
- `main.compute` — function by package-qualified name
- `main.go:18` — file name (relative to package) and line
- `compute` — Delve searches for a matching function in the current package
- `*<address>` — break at a specific memory address (advanced)

## Under the hood

When you set a breakpoint with `b main.go:18`, Delve:
1. Opens the source file to find line 18's byte offset.
2. Compiles to find the machine code address of that line.
3. Reads the byte at that address and saves it.
4. Overwrites it with `INT3` (0xCC on x86/amd64).
5. Adds the breakpoint to its internal table with the ID, address, and condition.

When the CPU executes `INT3`, the OS delivers a debug event to Delve. Delve looks up the address, evaluates the condition (if any), and either stops the process (returning control to you) or continues transparently.

## How Go uses it

- **Function breakpoints**: `b main.myHandler` — stop when a specific handler is called.
- **Library breakpoints**: `b (*net/http.Server).Serve` — stop in the standard library.
- **Test breakpoints**: `dlv test` then `b TestMyFunc` — stop at the start of a test.
- **Goroutine-aware breakpoints**: `b -goroutine 5 main.go:18` — only break in goroutine 5.
- **Conditional breakpoints**: `b main.go:42 if userID == 100` — stop when condition is true.

## Go example

```go
package main

import "fmt"

func main() {
	users := []string{"Alice", "Bob", "Charlie", "Diana"}

	for i, name := range users {
		processUser(i, name)
	}

	fmt.Println("All users processed.")
}

func processUser(id int, name string) {
	score := computeScore(id)
	fmt.Printf("User %d (%s): score=%d\n", id, name, score)
}

func computeScore(id int) int {
	base := id * 10
	bonus := id * id
	return base + bonus
}
```

Debugging session:
```bash
dlv debug
(dlv) break main.go:10                 # break at processUser call
(dlv) break main.processUser if id == 2 # conditional break
(dlv) break main.computeScore           # break on function entry
(dlv) bp                                # list all breakpoints
(dlv) clear 1                           # remove first breakpoint
(dlv) continue                          # run to next breakpoint
(dlv) print id
(dlv) print name
(dlv) continue                          # next iteration or next breakpoint
```

## Step-by-step execution

With `break main.computeScore if id == 2`:

1. Program starts. No breakpoint hit in `main`.
2. Loop iteration 0: `processUser(0, "Alice")`. `computeScore` is called. `id` is 0, not 2. Condition is false — breakpoint skipped.
3. Loop iteration 1: `processUser(1, "Bob")`. `id` is 1, not 2. Skipped.
4. Loop iteration 2: `processUser(2, "Charlie")`. `id` is 2. Condition is true. Delve stops at the first line of `computeScore`.
5. You can inspect `id`, `base`, `bonus` with `print`.
6. `continue` resumes. `computeScore` returns.
7. Loop iteration 3: `processUser(3, "Diana")`. `id` is 3, not 2. Skipped.
8. Program finishes.

## Common mistakes

- **Setting a breakpoint on a function that never gets called**: The program finishes without stopping. Verify the function is actually in the call path.
- **Breakpoint at wrong line**: Line numbers shift when you edit the file. Re-set breakpoints after editing. Delve detects stale breakpoints and warns.
- **Using `clearall` accidentally**: Removes all breakpoints. Use `clear <id>` to remove specific ones.
- **Conditional breakpoint with side effects**: The condition expression is evaluated each time. `if x++ > 5` would modify `x`. Avoid mutating conditions.
- **Breakpoints in inlined functions**: The compiler may inline small functions. Compile with `-gcflags="-N -l"` to disable inlining when debugging.

## Debugging walkthrough

Given the Go example above, suppose the score seems wrong for user 2.

```
(dlv) break main.computeScore
(dlv) continue
> main.computeScore() ./main.go:21 (hits goroutine(1):1 total:1)
    21: func computeScore(id int) int {
(dlv) print id
0
(dlv) continue
> main.computeScore() ./main.go:21 (hits goroutine(1):2 total:2)
    21: func computeScore(id int) int {
(dlv) print id
1
(dlv) continue
> main.computeScore() ./main.go:21 (hits goroutine(1):3 total:3)
    21: func computeScore(id int) int {
(dlv) print id
2
(dlv) next
    22:     base := id * 10
(dlv) next
    23:     bonus := id * id
(dlv) print base
20
(dlv) print bonus
4
(dlv) next
    24:     return base + bonus
(dlv) print base + bonus
24
```

You find `computeScore(2)` returns `20 + 4 = 24`. If you expected a different value, the bug is in `computeScore`'s formula — not in the caller. The breakpoint proved the inputs and intermediate values are correct.

## Production notes

- **Never leave breakpoints in production**: Compiling with debug symbols and leaving breakpoints unattended is a security risk. Debug builds should never be deployed.
- **Remote debugging with breakpoints**: `dlv debug --headless --listen=:2345` starts a headless server. Connect from an IDE or another Delve instance. Breakpoints work the same way.
- **Logpoints** (breakpoints that log instead of stopping): Use `on <id> <cmd>` — e.g. `on 1 p name` prints `name` each hit without stopping. Great for tracing production-like environments without modifying code.
- **Breakpoint explosion on hot paths**: Avoid setting breakpoints in frequently called functions (e.g. hot loop, middleware). Each hit triggers a debug event. Use conditional breakpoints or logpoints instead.

## Performance implications

- Each breakpoint hit causes a context switch from the target process to the debugger to evaluate the breakpoint — ~1–10 microseconds.
- Conditional breakpoints add expression evaluation overhead. Keep conditions simple (integer comparison, not complex function calls).
- In headless mode, the network round-trip to the remote client adds latency on each breakpoint hit.
- Breakpoints with `on` (logpoints) still cause a breakpoint hit, but Delve automatically continues — the overhead is still there, just invisible to you.

## Practice task

Start `dlv debug` on this lesson's program. Set the following breakpoints:
1. A function breakpoint on `main.processUser`
2. A line breakpoint on the line where `computeScore` is called (inside `processUser`)
3. A conditional breakpoint on `main.computeScore` that stops only when `id == 3`
4. Use `bp` to list all breakpoints
5. Use `clear` to remove the first breakpoint
6. Use `continue` to run through the program and observe the conditional breakpoint trigger

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/15-breakpoints
go test ./curriculum/modules/06-testing-debugging-refactoring/lessons/15-breakpoints
```

## Review questions

1. How do you set a breakpoint on a function named `doSomething` in package `main`?
2. What is the difference between `b main.go:10` and `b main.compute`?
3. How do you set a conditional breakpoint that triggers only when `x` is greater than 5?
4. What does `bp` do in a Delve session?
5. Why might a breakpoint on a small function never get hit, even though the function is called?

## NEXT UP

Stepping through code — controlling execution line by line with next, step, stepout, and continue.
