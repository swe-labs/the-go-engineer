# defer mechanics

## Learning objective

Use the `defer` statement to schedule function calls at return, predict LIFO execution order, understand argument capture versus closure capture, and exploit deferred functions with named return values.

## Why this matters

Resource management is the most common source of bugs in systems programming: files left open, mutexes left locked, HTTP connections leaked. `defer` eliminates entire categories of cleanup bugs by binding resource release to function exit. It is one of Go's signature features, and mastering its mechanics -- especially the subtle differences between argument evaluation and closure capture -- is essential for writing correct production code.

## Mental model

`defer` is a stack of to-do items attached to the current function call. When you say `defer f()`, you are adding `f()` to the top of the stack. The arguments to `f` are evaluated and frozen immediately (like taking a snapshot), but the function body does not execute until the surrounding function returns. When the function returns (normally, or via panic), the stack is popped in LIFO order -- last deferred, first executed.

## Core idea

Three key rules govern `defer`:

1. **Arguments are evaluated immediately**: `defer fmt.Println(x)` captures the value of `x` at the defer statement, not at return time.
2. **Execution is deferred until the surrounding function returns**: The deferred call runs after the function's normal return or panic, but before the function returns to its caller.
3. **LIFO order**: Multiple defers execute in last-in-first-out order -- like a stack.

Named return values interact with defer in a powerful way: a deferred closure can read and modify the named return value of the enclosing function.

## Under the hood

In Go's runtime, each goroutine maintains a linked list of `_defer` structs. When the compiler encounters a `defer` statement, it inserts a call to `runtime.deferproc` which:

1. Allocates a `_defer` struct on the heap.
2. Stores the function pointer and frozen arguments.
3. Prepends the struct to the goroutine's defer list.

When the function returns, the compiler inserts `runtime.deferreturn` in the epilogue. This function pops `_defer` structs from the list (LIFO) and executes them. If multiple returns exist, `deferreturn` is called at each one.

In Go 1.14+, the compiler optimized the common case: for functions with defer but no panic-recover, deferred calls are inlined at each return site instead of going through the runtime. This makes defer in hot paths nearly free.

## How Go uses it

```go
// File close
f, err := os.Open(path)
if err != nil { return err }
defer f.Close()

// Mutex unlock
mu.Lock()
defer mu.Unlock()

// Transaction rollback on error
tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }
defer tx.Rollback() // no-op if committed

// Timing
defer func() { fmt.Println("elapsed:", time.Since(start)) }()
```

## Go example

```go
package main

import "fmt"

func main() {
	fmt.Println("start")
	defer fmt.Println("first defer")
	defer fmt.Println("second defer")
	defer fmt.Println("third defer")
	fmt.Println("end")
}
```

Output:

```
start
end
third defer
second defer
first defer
```

The key observation: "end" prints before the deferred calls. The defers run after the function body completes. LIFO order means "third defer" runs first.

Named return example:

```go
func double(x int) (result int) {
	defer func() {
		result *= 2
	}()
	return x + 1
}
```

`double(5)` returns `12`, not `6`. The deferred closure modifies `result` after `x + 1` (6) is assigned to `result`, doubling it to 12.

## Step-by-step execution

For the LIFO example:

1. `fmt.Println("start")` executes immediately → prints "start".
2. `defer fmt.Println("first defer")` -- the argument `"first defer"` is evaluated and frozen. The `_defer` struct is prepended to the goroutine's defer list.
3. `defer fmt.Println("second defer")` -- same process. This struct is now at the front of the list.
4. `defer fmt.Println("third defer")` -- now at the front.
5. `fmt.Println("end")` executes immediately → prints "end".
6. Function is about to return. The runtime calls `deferreturn`.
7. Pop the front of the list: `fmt.Println("third defer")` executes → prints "third defer".
8. Pop next: `fmt.Println("second defer")` → prints "second defer".
9. Pop next: `fmt.Println("first defer")` → prints "first defer".
10. Function returns to caller.

For the named return example:

1. `double(5)` is called. `result` is initialized to zero value (0 for int).
2. `defer func() { result *= 2 }()` -- the closure captures `result` by reference (closure capture, not argument capture).
3. `return x + 1` evaluates to `6`. This value is assigned to `result` (the named return).
4. The deferred function runs: `result *= 2` → `result` becomes `12`.
5. Function returns `result` (12).

## Common mistakes

- **Assuming deferred calls run at block scope**: `defer` is function-scoped, not block-scoped. If you `defer f()` inside a `for` loop, the deferred calls do not run at the end of each iteration; they all run when the function returns.

```go
func readFiles(files []string) error {
    for _, f := range files {
        file, err := os.Open(f)
        if err != nil { return err }
        defer file.Close() // defers all add up; only closes at function return!
        // ... read from file
    }
    return nil
}
```

- **Argument capture vs closure capture**: `defer fmt.Println(x)` captures `x`'s value at defer time. `defer func() { fmt.Println(x) }()` captures the variable `x` by reference (closure), seeing the value at defer-execution time. These are different.

- **Defer in a loop leaks resources**: If a function processes many items in a loop and defers cleanup for each, cleanup only runs when the function returns, not after each item. This causes O(n) resource retention.

- **Assuming deferred code runs before `os.Exit`**: `os.Exit` terminates the program immediately; deferred functions do not run.

## Debugging walkthrough

```go
package main

import "fmt"

func main() {
	x := 1
	defer fmt.Println("defer 1:", x)
	x = 2
	defer fmt.Println("defer 2:", x)
	x = 3
	fmt.Println("end:", x)
}
```

**Prediction attempt**: Many beginners expect both defers to see `x=3`, printing "defer 1: 3" and "defer 2: 3".

**Actual output**:
```
end: 3
defer 2: 2
defer 1: 1
```

**Explanation**: Arguments to deferred functions are evaluated at the defer statement, not at return. At `defer fmt.Println("defer 1:", x)`, `x` is `1`. At `defer fmt.Println("defer 2:", x)`, `x` is `2`.

**LIFO order**: "defer 2" (deferred second) runs before "defer 1" (deferred first).

**Fix if you want the current value**: Use a closure:

```go
x := 1
defer func(v int) { fmt.Println("defer 1:", v) }(x)
x = 2
defer func(v int) { fmt.Println("defer 2:", v) }(x)
```

Or capture the variable reference:

```go
x := 1
defer func() { fmt.Println("defer 1:", x) }()
x = 2
defer func() { fmt.Println("defer 2:", x) }()
// Now both see x=3 at return time
```

## Production notes

- **Pair every acquire with a defer**: The idiomatic Go pattern is `resource, err := acquire(); if err != nil { return }; defer resource.Release()`. This prevents leaks even as the function grows.

- **Defer is not free in hot loops**: Go 1.14+ optimizes defer in functions without panic/recover, but each defer still has a small cost. In extreme hot loops, inline cleanup instead of defer.

- **Named returns with defer**: Use `defer` to modify named returns for logging, metrics, or recovery. This is a powerful pattern but can make control flow harder to follow. Document it when you use it.

- **Defers run on panic**: This is critical for cleanup reliability. Even if your function panics, all deferred cleanup runs during stack unwinding.

## Performance implications

- Before Go 1.14, each `defer` caused a heap allocation for the `_defer` struct. This made defer expensive in tight loops.
- Go 1.14+ optimized most defer calls by inlining them at return sites when the compiler can prove no panic/recover is needed. In this case, defer costs roughly the same as a direct call.
- Defer in a function with `recover` still goes through the runtime path, which is slower.
- For most production code, defer cost is negligible. Only optimize if profiling shows defer in a hot path.

## Practice task

Write a function `countdown()` that:

1. Prints `"start"`.
2. Defers `fmt.Println("done")`.
3. Uses a loop to defer three calls: `defer fmt.Println(i)` for `i = 1, 2, 3`.
4. Prints `"launch"`.

Call `countdown()` and predict the output before running.

Then write a function `addLogger(fn func(int) int) func(int) (int, error)` that wraps a function with a deferred log statement printing `"called with %d, result %d\n"`. Use named returns.

## Tests / verification

```bash
go run ./curriculum/modules/04-functions-errors-data-semantics/lessons/15-defer-mechanics
go test ./curriculum/modules/04-functions-errors-data-semantics/lessons/15-defer-mechanics
```

## Review questions

1. In what order do multiple deferred calls execute?
2. When are the arguments to a deferred function evaluated?
3. How can a deferred function modify the return value of the enclosing function?
4. If you defer `file.Close()` in a loop over 10 files, when do the close operations actually run?
5. Does `os.Exit(0)` trigger deferred functions?

## NEXT UP

defer for cleanup
