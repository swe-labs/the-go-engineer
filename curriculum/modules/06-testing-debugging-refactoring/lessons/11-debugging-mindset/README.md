# Debugging mindset

## Learning objective

Apply systematic debugging techniques including hypothesis-driven debugging, printf debugging, error message analysis, and rubber duck debugging to diagnose and fix Go programs.

## Why this matters

Bugs are inevitable. The skill that separates effective engineers from frustrated ones is not writing perfect code — it is debugging efficiently. A systematic approach turns debugging from a guessing game into a scientific process: form a hypothesis, test it with evidence, and narrow down until the root cause is found. This lesson gives you the mental framework to debug any Go program, from a misbehaving function to a production outage.

## Mental model

Debugging is a scientific investigation. Your code is a system that produces an observable output. The bug is a discrepancy between actual and expected behavior. Your job is to find the cause by:

1. **Observing** the failure (symptom).
2. **Forming a hypothesis** about the root cause.
3. **Testing the hypothesis** with an experiment (a modified print, a unit test, a breakpoint).
4. **Narrowing** the search space based on results.
5. **Iterating** until the root cause is confirmed.
6. **Fixing and verifying**.

```
Symptom → Hypothesis → Experiment → Result
                  ↑                    │
                  └── Narrow ──────────┘
```

## Core idea

Systematic debugging relies on a few core techniques:

| Technique | How it works |
|---|---|
| **Hypothesis-driven** | Form a falsifiable guess. Test it. If wrong, form a new one. |
| **Binary search** | Comment out half the code. If the bug disappears, it is in the commented half. Repeat. |
| **Printf debugging** | Add `fmt.Printf` at key points to trace state. |
| **Error message analysis** | Read the full stack trace. Identify the exact line and call chain. |
| **Rubber duck** | Explain the problem aloud to an inanimate object. The act of explaining often reveals the flaw. |
| **Minimal reproduction** | Create the smallest possible program that exhibits the bug. |

The key mindset shift: **do not assume anything**. Verify every assumption. "I know this function works" is a hypothesis, not a fact.

## Under the hood

When you run a Go program and it panics, the runtime prints a stack trace:

```
panic: runtime error: index out of range [5] with length 3

goroutine 1 [running]:
main.sumSlice(...)
    main.go:10 +0x45
main.main()
    main.go:6 +0x29
```

This trace tells you:
- The panic type: `index out of range [5] with length 3`.
- The exact file and line: `main.go:10`.
- The call chain: `main.main` called `main.sumSlice`.
- The goroutine: 1 (the main goroutine).

Reading stack traces is a superpower. It tells you *exactly* where the runtime detected a problem. For logic errors (wrong output, no panic), stack traces do not help — you need printf debugging or a debugger.

## How Go uses it

Go's tooling supports debugging at multiple levels:

- **`go vet`** catches common mistakes at compile time (unused variables, impossible conditions).
- **`go test -v`** shows test output and failures.
- **`fmt.Printf` / `log.Println`** for printf debugging.
- **`runtime/debug.PrintStack()`** to print the current goroutine's stack.
- **`panic` / `recover`** for crash analysis.
- **Delve** (`dlv`) for interactive debugging (breakpoints, variable inspection).

For production debugging, Go's `net/http/pprof` and `runtime/pprof` provide profiling, and `expvar` exposes runtime variables.

## Go example

```go
package main

import (
	"fmt"
)

// SumUntil returns the sum of numbers from 1 to n.
func SumUntil(n int) int {
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
	}
	return sum
}

// Average returns the average of a slice.
func Average(nums []int) float64 {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return float64(sum) / float64(len(nums))
}

func main() {
	fmt.Println("SumUntil(5) =", SumUntil(5))
	fmt.Println("Average([2,4,6]) =", Average([]int{2, 4, 6}))
	fmt.Println("Average([]) =", Average([]int{})) // bug!
}
```

This program will panic at the last line: division by zero because `len(nums)` is 0.

## Step-by-step execution

Running `go run main.go`:

1. Prints `SumUntil(5) = 15` (correct).
2. Prints `Average([2,4,6]) = 4` (correct: (2+4+6)/3 = 4).
3. Calls `Average([]int{})`. `sum` is 0, `len(nums)` is 0.
4. `0 / 0` triggers a runtime panic: `panic: runtime error: integer divide by zero`.
5. Stack trace is printed. The trace points to `main.go:22`: `return float64(sum) / float64(len(nums))`.

**Debugging process**:

1. **Observe**: The panic says "integer divide by zero" at `Average`.
2. **Hypothesis**: The function does not check for an empty slice before dividing.
3. **Test**: Read the code — `return float64(sum) / float64(len(nums))` — there is no guard.
4. **Fix**: Add a guard: `if len(nums) == 0 { return 0 }`.
5. **Verify**: Re-run. No panic.

## Common mistakes

- **Fixing symptoms, not causes.** If you add a nil check where a nil value appears, but the real bug is upstream where the value should have been initialized, you are patching a symptom.

- **Changing code without understanding it.** Making random changes to see if the bug goes away wastes time. Each change should test a specific hypothesis.

- **Ignoring error messages.** "I see the panic but I did not read the line number." The stack trace tells you exactly where to look.

- **Not reproducing the bug.** If you cannot reproduce it, you cannot fix it. If the bug is intermittent, look for race conditions, timing, or random input.

- **Debugging in production without isolation.** Adding debug prints to production code can change timing and mask the bug. Reproduce locally first.

## Debugging walkthrough

A function returns the wrong result:

```go
// Discount applies a percentage discount to a price.
func Discount(price, percent int) int {
    return price * percent / 100
}
```

**Symptom**: `Discount(199, 10)` returns `19`, but the business expects `19.9` rounded properly.

**Hypothesis 1**: Integer division truncates toward zero. Test: `fmt.Println(199 * 10 / 100)` → `19`. Correct hypothesis.

**Hypothesis 2**: The fix should use floating-point or round properly. Test:
```go
fmt.Println(math.Round(float64(199) * 10.0 / 100.0)) // 20
```
But the business wants `19.9` (no rounding). So the return type should be `float64`.

**Fix**: Change return type to `float64` and use floating-point arithmetic:
```go
func Discount(price, percent float64) float64 {
    return price * percent / 100.0
}
```

**Verification**: `Discount(199, 10)` returns `19.9`.

## Production notes

- **Reproduce in a test.** Before fixing a bug, write a failing test that reproduces it. Then fix the code, watch the test pass, and keep the test as a regression guard.
- **Log strategically.** Production logging should include enough context to diagnose failures. `log.Printf("user %s: payment failed: %v", userID, err)` is better than `log.Printf("payment failed")`.
- **Use structured logging.** Libraries like `log/slog` (Go 1.21+) produce structured output that can be queried and aggregated.
- **Postmortems.** After fixing a production bug, write a brief postmortem: what happened, why, how it was found, and how to prevent it next time. This systematizes learning.
- **Debugging is a team sport.** Pair debug when stuck. A fresh perspective often spots what you missed.

## Performance implications

- **Printf debugging is slow.** Adding thousands of `fmt.Printf` calls in a hot loop changes timing and can mask race conditions. Remove debug prints after fixing.
- **Conditional logging** with a debug flag (`var debug = flag.Bool("debug", false, ...)`) keeps debug prints available without performance impact in production.
- **Stack traces are expensive.** `debug.PrintStack()` captures the full goroutine stack, which involves symbol lookup. Use sparingly.
- **The best debugging performance optimization is getting it right the first time.** Write tests before code (TDD) to reduce debugging time.

## Practice task

Take this buggy program and debug it:

```go
func main() {
    numbers := []int{1, 2, 3, 4, 5}
    fmt.Println("Sum of evens:", SumEvens(numbers))
}

func SumEvens(nums []int) int {
    sum := 0
    for _, n := range nums {
        if n%2 == 0 {
            continue
        }
        sum += n
    }
    return sum
}
```

The program prints `Sum of evens: 9` but it should print `Sum of evens: 6` (2 + 4). Use systematic debugging to find and fix the bug. Write a test that reproduces it before fixing.

## Tests / verification

```bash
go run ./curriculum/modules/06-testing-debugging-refactoring/lessons/11-debugging-mindset
go test -v ./curriculum/modules/06-testing-debugging-refactoring/lessons/11-debugging-mindset
```

## Review questions

1. What is the first step in a systematic debugging approach?
2. How does the binary search technique apply to debugging?
3. Why should you reproduce a bug in a test before fixing it?
4. What information does a Go panic stack trace provide?
5. What is rubber duck debugging and why does it work?

## NEXT UP

Reading stack traces — interpreting Go panic and error stack traces to pinpoint failures.
