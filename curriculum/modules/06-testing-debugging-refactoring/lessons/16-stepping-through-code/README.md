# Stepping through code

## Learning objective

Use Delve's stepping commands — `step`, `next`, `stepout`, and `continue` — to navigate code execution line by line, into functions, and back out again.

## Why this matters

Breakpoints stop the program. Stepping moves it forward in controlled increments. Without stepping, you can only observe the state when execution happens to stop. With stepping, you can follow every intermediate state, every assignment, and every function call boundary. This is essential for understanding unfamiliar code, verifying logic, and finding off-by-one errors or subtle control-flow bugs.

## Mental model

Think of the program counter (PC) as a cursor on a line of source code. Stepping commands move that cursor:

- **`next` (n)**: Advance to the next line in the *same* function. If the current line calls another function, execute it entirely and stop after it returns.
- **`step` (s)**: If the current line calls a function, move the cursor into that function's first line. If not, behave like `next`.
- **`stepout` (so)**: Finish the current function and stop at the line after the call site in the caller.
- **`continue` (c)**: Resume full-speed execution until the next breakpoint or program exit.

## Core idea

| Command | Short | Where you land |
|---|---|---|
| `next` | `n` | Next line in the same function, stepping over calls |
| `step` | `s` | Into a called function, or next line if no call |
| `stepout` | `so` | Return to caller, at the line after the call |
| `continue` | `c` | Next breakpoint or program end |
| `restart` | `r` | Start the program over from the beginning |

The key distinction: `next` steps *over* a function call (running it invisibly), while `step` steps *into* it (so you can trace its internals).

## Under the hood

`next` works by setting a temporary breakpoint on the next source line in the current function, then continuing execution. When that breakpoint is hit, Delve removes it and stops. If the current line is a function call, `next` sets the temporary breakpoint on the line *after* the call in the current function — the called function runs to completion before stopping.

`step` checks if the current line contains a function call. If yes, it sets a temporary breakpoint on the first line of the called function (using its symbol address). If the line has no function call, `step` behaves like `next`.

`stepout` sets a temporary breakpoint at the return address of the current function (the line in the caller after the call site) and continues.

All three are implemented as temporary breakpoints — they fire once and are automatically cleared.

## How Go uses it

- **Trace data flow**: Step into a function to trace how input transforms into output line by line.
- **Verify conditionals**: Put a breakpoint before an `if` statement, then `next` through each branch to confirm the correct one executes.
- **Follow recursion**: Step into a recursive function and watch the call stack grow. Use `stack` to see the depth.
- **Skip library code**: If you accidentally `step` into `fmt.Println`, use `stepout` to return to your code immediately.
- **Navigate test failures**: In `dlv test`, break in a test function and `next` through assertions.

## Go example

```go
package main

import "fmt"

func main() {
	result := factorial(5)
	fmt.Println("Factorial:", result)
}

func factorial(n int) int {
	fmt.Println("factorial called with:", n)
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}
```

Debugging session:
```bash
dlv debug
(dlv) break main.factorial
(dlv) continue
> main.factorial() ./main.go:9
(dlv) next                 # step over fmt.Println
> main.go:10               # if n <= 1
(dlv) print n              # 5
(dlv) next                 # step over if, into return
> main.go:12               # return n * factorial(n-1)
(dlv) step                 # step into factorial(n-1)
> main.factorial() ./main.go:9  # new call, n=4
(dlv) stack                # see both frames
(dlv) stepout              # return from this call
> main.go:12               # back in the first call, after factorial(n-1)
(dlv) continue             # let it finish
```

## Step-by-step execution

Starting from `main` with a breakpoint at `factorial(5)`:

1. `continue` → hits breakpoint at `func factorial(n int) int {`.
2. `next` → executes `fmt.Println` call (prints "factorial called with: 5"), lands on `if n <= 1`.
3. `n` is 5, condition is false.
4. `next` → lands on `return n * factorial(n-1)`.
5. `step` → **enters** `factorial(4)`. A new frame appears on the stack.
6. `next` → executes `fmt.Println`, lands on `if n <= 1`.
7. `stepout` → this call (n=4) completes, returns to the caller's `return` line.
8. The original call's `factorial(n-1)` now has its result. `continue` to finish.

## Common mistakes

- **Stepping into `fmt.Println`**: It happens to every beginner. The call is in your code; `step` enters it. Use `next` to skip over library calls, or use `stepout` to escape.
- **Using `next` when you meant `step`**: If you `next` over a function call, you lose the ability to trace its internals. The program runs the function invisibly. Set a breakpoint inside the function first if you want to enter it later.
- **Losing your place with multiple `continue`s**: `continue` runs until the next breakpoint. If you have no breakpoints, the program finishes. Set a breakpoint further down first.
- **Thinking `stepout` returns to `main`**: `stepout` returns to the *caller* of the current function — which might be another intermediate function, not `main`.
- **Stepping into closures**: Stepping into a closure or anonymous function might land you inside the runtime's defer or goroutine machinery. Use `stepout` to get back.

## Debugging walkthrough

The factorial program above has a bug: there's no `fmt.Println` for the n=0 case because the base case is `n <= 1`. Let's trace it:

```
(dlv) break main.factorial
(dlv) continue
> main.factorial() ./main.go:9
(dlv) print n
5
(dlv) next                  # print
(dlv) next                  # if
(dlv) next                  # return n * factorial(n-1)
(dlv) step                  # into factorial(4)
> main.factorial() ./main.go:9
(dlv) print n
4
(dlv) next; next; next; step  # repeat until n=1
> main.factorial() ./main.go:9
(dlv) print n
1
(dlv) next                  # prints "factorial called with: 1"
(dlv) next                  # if n <= 1 -> true
(dlv) next                  # return 1
(dlv) stepout               # back to n=2's return
(dlv) continue              # finish
```

The step-by-step confirms that each recursive call gets the correct `n` and returns the right value.

## Production notes

- **Stepping is slow**: Each step requires a context switch to the debugger. Stepping through a loop with 10,000 iterations is impractical. Set a conditional breakpoint inside the loop instead.
- **Time-sensitive code**: Stepping pauses execution. If your code talks to external services or has timeouts, the connections may drop while you step. Be aware of this for network code.
- **Multiple goroutines**: Stepping controls one goroutine. Other goroutines continue running. Use `goroutine <id>` to switch, then `step` in that goroutine.
- **Reversed execution**: Delve does not support stepping backwards (reverse debugging). If you miss a state, you must restart and set a breakpoint earlier.

## Performance implications

- `next`, `step`, and `stepout` each trigger a breakpoint hit and a context switch (~1–10 µs). Human reaction time (~200ms per keystroke) dominates, so this matters only for automated tooling.
- Temporary breakpoints from stepping are tracked in Delve's internal table. They are cleaned up immediately after firing, so memory impact is negligible.
- When stepping through a function that the compiler inlined, Delve may show source lines out of order or skip lines. Use `-gcflags="-N -l"` to disable optimisations if exact line fidelity is required.

## Practice task

Start `dlv debug` on this lesson's program (factorial with a debug print). Set a breakpoint on `main.factorial`. Use the stepping commands to:
1. `next` over the `fmt.Println` call inside `factorial`.
2. `step` into the recursive `factorial(n-1)` call.
3. Use `stack` to see the call stack depth at recursion level 3.
4. Use `stepout` to return from a deeply nested call back to its caller.
5. Use `continue` to finish execution.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/16-stepping-through-code
go test ./curriculum/modules/06-testing-debugging-refactoring/lessons/16-stepping-through-code
```

## Review questions

1. What is the difference between `next` and `step`?
2. If you accidentally `step` into `fmt.Sprintf`, what command gets you back to your code?
3. What does `stepout` do, and where does execution resume?
4. How does `next` work internally in terms of breakpoints?
5. Why might stepping through a function show lines in a different order than the source file?

## NEXT UP

Inspecting variables — how to print, watch, and understand variable state in Delve.
