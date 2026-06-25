# Return values

## Learning objective

Write functions that return values using both explicit and named return syntax, understand naked returns and their pitfalls, and use early returns to simplify control flow.

## Why this matters

Return values are how functions communicate results back to callers. Using named returns and early returns effectively makes code more readable and less error-prone. Misunderstanding zero-value returns or accidentally using naked returns in complex functions leads to subtle bugs.

## Mental model

A return value is the output of a function. Think of a function as a vending machine: you insert arguments (coins), press a button, and a product (the return value) comes out. The `return` statement is the slot where the result appears.

Named return values are like pre-labelled bins. Before the function body runs, Go creates variables with those names and initialises them to their zero values. A bare `return` returns whatever those variables currently hold.

## Core idea

A function declares its return type after the parameter list:

```go
func add(a, b int) int {
    return a + b
}
```

**Named return values** give the return value a name and type:

```go
func divide(a, b int) (result int) {
    result = a / b
    return // naked return — returns current value of result
}
```

A **naked return** (`return` without arguments) returns the current values of all named return variables. Use it only in short functions where it improves readability.

**Zero value returns**: If a function declares a return type but you `return` without a value (or reach the end without any return), Go returns the zero value for the type — but only if the function has named returns.

**Early return**: Using `return` before the end of the function to exit early, typically when an error or special case is detected.

## Under the hood

When Go compiles a function with a return type, it allocates space on the stack for the return value. For named returns, the compiler creates local variables and inserts an implicit `return` statement at the end that returns their current values.

A naked `return` is syntactic sugar: the compiler rewrites it to `return result` (using the named variable). Naked returns are disallowed in functions without named return variables.

## How Go uses it

- **Named returns** appear in interface implementations and long functions where the return variable's purpose is documented by its name.
- **Naked returns** are used in short functions where the named return variable's purpose is obvious.
- **Early returns** are the idiomatic Go pattern for error handling — the "happy path" continues at the bottom of the function.
- **Zero returns** appear in recursive or loop-based functions where the base case returns zero.

## Go example

```go
package main

import "fmt"

// Single explicit return.
func square(n int) int {
	return n * n
}

// Named return with naked return.
func divide(a, b int) (result int) {
	if b == 0 {
		return 0 // explicit early return, not naked
	}
	result = a / b
	return // naked — returns result
}

// Named returns with multiple variables.
func rectangleProps(w, h int) (area, perimeter int) {
	area = w * h
	perimeter = 2 * (w + h)
	return
}

// Early return pattern.
func safeDivide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

func main() {
	fmt.Println("square(5) =", square(5))
	fmt.Println("divide(10, 3) =", divide(10, 3))
	fmt.Println("divide(5, 0) =", divide(5, 0))

	a, p := rectangleProps(4, 5)
	fmt.Printf("rectangleProps(4, 5): area=%d, perimeter=%d\n", a, p)

	fmt.Println("safeDivide(9, 3) =", safeDivide(9, 3))
}
```

## Step-by-step execution

For `divide(10, 3)`:

1. Stack frame for `divide` is created. Named return variable `result` is initialised to `0` (zero value for `int`).
2. `b == 0` is false, so early return is skipped.
3. `result = 10 / 3` assigns `3` to `result`.
4. Naked `return` returns the current value of `result`, which is `3`.

For `divide(5, 0)`:

1. `result` initialised to `0`.
2. `b == 0` is true, explicit `return 0` returns `0` immediately. The naked `return` at the end is never reached.

## Common mistakes

- **Using naked returns in long functions**: Naked returns hide what is being returned. In functions longer than ~10 lines, prefer explicit `return result`.
- **Forgetting to return a value**: A function with a return type must have a `return` on every path. Go flags this as a compile error.
- **Shadowing named returns**: `result := a / b` inside the function creates a new local variable that shadows the named return. Use `=` not `:=`.
- **Naked return without named returns**: A bare `return` in a function without named return values is a compile error.

## Debugging walkthrough

This code compiles but produces unexpected output:

```go
package main

import "fmt"

func getValue(flag bool) (val int) {
	if flag {
		val := 42 // bug: shadows the named return
		return val
	}
	return
}

func main() {
	fmt.Println(getValue(true))  // expects 42
}
```

**Symptom**: Prints `0` instead of `42`.

**Root cause**: Inside the `if` block, `val := 42` uses short variable declaration (`:=`), which creates a *new* local variable `val` that shadows the named return variable. The `return val` returns the shadowed variable's value (42), but the naked `return` at the end returns the original `val` (0) when `flag` is `true` — wait, no. In this specific case, `if flag` is true, we hit `return val` which returns 42. But if `flag` is false, we hit `return` which returns the named return `val` (0). So this actually works for `true`. Let's re-examine.

Actually the real bug is subtler. Let me reconsider. The code would print 42 when flag is true. But it's bad practice. A more accurate bug:

```go
func getValue(flag bool) (val int) {
	if flag {
		val = 42
	}
	// forgot return — but named returns have implicit return?
}
```

Actually Go does not insert implicit return for named returns. You must still have `return` or `return val`. From the spec: "Named return values act like local variables... A return statement without arguments returns the named return values." But you still need a `return` statement. If you omit `return` entirely, it's a compile error.

Let me provide a correct walkthrough instead.

Consider this code that does not compile:

```go
package main

func half(n int) int {
	if n%2 == 0 {
		return n / 2
	}
	// missing return when n is odd
}

func main() {}
```

**Symptom**: `missing return at end of function`.

**Root cause**: When `n` is odd, the `if` body is skipped and no `return` follows.

**Fix**: Add a return after the if block:

```go
func half(n int) int {
	if n%2 == 0 {
		return n / 2
	}
	return n / 2 // or return 0
}
```

## Production notes

- **Prefer explicit returns** over naked returns in all but the shortest functions (2-3 lines).
- **Named returns document intent**: `func parse(input string) (result int, err error)` makes the return values self-documenting.
- **Zero-value returns** are relied upon in `switch`/`select` with default cases, but always include an explicit `return` for clarity.
- **Early return** is the idiomatic Go pattern: handle errors first, return early, keep the happy path unindented.

## Performance implications

- Named return values do not introduce overhead — they are identical to local variables.
- Early returns do not affect performance (no "single exit point" dogma in Go).
- Returning large structs copies the value. The compiler may optimise this with a hidden pointer, but for large types consider returning a pointer.

## Practice task

Write a function `factorial(n int) (result int)` that:
- Returns `1` when `n <= 1`
- Returns `n * factorial(n-1)` otherwise
- Uses a named return with a naked return at the end

Then in `main()`, call `factorial(5)` and `factorial(0)` and print results.

## Tests / verification

```bash
go run ./curriculum/modules/04-functions-errors-data-semantics/lessons/03-return-values
go test ./curriculum/modules/04-functions-errors-data-semantics/lessons/03-return-values
```

## Review questions

1. What syntax declares a named return value?
2. What is a naked return, and when is it valid?
3. What happens if you use `:=` to assign to a named return variable inside the function?
4. Write a function signature that returns an `int` named `sum`.
5. Is `return` by itself valid in a function with signature `func() int`? Why or why not?

## NEXT UP

Multiple return values — the `(T, error)` idiom and the blank identifier.
