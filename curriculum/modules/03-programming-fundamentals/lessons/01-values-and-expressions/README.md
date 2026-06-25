# Values and expressions

## Learning objective

Write and evaluate Go expressions using literals, variables, and arithmetic/boolean operators, and predict the result of any expression by applying Go's operator precedence and type rules.

## Why this matters

Every line of Go code is or contains an expression. Conditionals, arithmetic, slice indices, function arguments, and struct field access all evaluate to values. If you cannot reason about expressions precisely, you cannot reason about program behavior at all. Professional Go engineers read and debug expressions daily, so mastering them is the first step to writing correct code.

## Mental model

An expression is a recipe that produces a value. You combine ingredients (operands) with tools (operators) in a specific order. The expression `3 + 4 * 2` is a recipe: start with `4 * 2 = 8`, then `3 + 8 = 11`. Operators have a priority (precedence), and when two operators have the same priority, evaluation proceeds left to right (associativity). Parentheses override the default order, just like in algebra.

This model breaks only for short-circuit boolean operators (`&&`, `||`) which stop evaluating early, and for function calls whose arguments are evaluated before the call but the function itself may have side effects.

## Core idea

A _value_ is data -- a number, a boolean, a string, or a compound type. An _expression_ is source code that the Go compiler and runtime evaluate to produce a value. Every expression has a _type_ (known at compile time) and a _value_ (known or computed at runtime).

Kinds of expressions in Go:

| Expression kind | Example | Evaluates to |
|---|---|---|
| Literal | `42`, `"hello"`, `true` | The value itself |
| Variable | `x`, `name` | The current value stored in the variable |
| Binary | `a + b`, `x == y` | The result of the operator |
| Unary | `-x`, `!ok` | The transformed operand |
| Function call | `math.Sqrt(9)` | The return value |
| Index | `arr[0]` | The element at the position |
| Selector | `p.Name` | The field or method value |

## Under the hood

Go parses source code into an Abstract Syntax Tree (AST). The compiler type-checks every expression node: it verifies that operator operand types match (e.g., `int + int` is valid, `int + string` is not) and resolves any type conversions. At runtime, the generated machine code evaluates the expression according to Go's evaluation order:

1. Operands evaluate first, left to right.
2. Operator applies once all operands are ready.
3. Result becomes an operand for the enclosing expression.

Integer arithmetic uses two's complement wraparound on overflow -- no panic occurs for signed or unsigned overflow. Floating-point follows IEEE 754. There is no implicit type coercion; every conversion must be explicit.

## How Go uses it

Go is an expression-oriented language at the statement level. Common constructs that rely on expressions:

- **Assignments**: `x = y + z` -- the right side is an expression.
- **If conditions**: `if x > 0` -- the condition is a boolean expression.
- **For loops**: `for i < 10` -- the condition is an expression.
- **Switch**: `switch x` and `case 1:` -- both are expressions.
- **Return**: `return a + b` -- the return value is an expression.
- **Defer, go, send/receive**: all take expressions.

Go's standard library provides `go/ast` and `go/types` packages that let you inspect expressions programmatically -- the same tools the compiler uses internally.

## Go example

```go
package main

import "fmt"

func main() {
	a := 10
	b := 3

	fmt.Println("a + b =", a+b)
	fmt.Println("a - b =", a-b)
	fmt.Println("a * b =", a*b)
	fmt.Println("a / b =", a/b)
	fmt.Println("a % b =", a%b)

	fmt.Println("(a + b) * 2 =", (a+b)*2)
	fmt.Println("a + b * 2 =", a+b*2)
}
```

Run with `go run .` to see the output. Notice `a / b` produces `3` (integer division discards the remainder). Also compare `(a+b)*2` vs `a+b*2` -- parentheses change the result from `26` to `16` because `*` has higher precedence than `+`.

## Step-by-step execution

For the expression `(a + b) * 2` with `a=10`, `b=3`:

1. Evaluate `a` → `10`.
2. Evaluate `b` → `3`.
3. Apply `+` to `10` and `3` → `13`.
4. Evaluate literal `2` → `2`.
5. Apply `*` to `13` and `2` → `26`.
6. Result `26` is assigned or printed.

For `a + b * 2`:

1. Evaluate `b` → `3`.
2. Evaluate literal `2` → `2`.
3. Apply `*` to `3` and `2` → `6`.
4. Evaluate `a` → `10`.
5. Apply `+` to `10` and `6` → `16`.
6. Result `16` is assigned or printed.

The key difference: `*` binds more tightly than `+`, so `b * 2` happens first in the second expression.

## Common mistakes

- Mistake: `7 / 2` produces `3` instead of `3.5`.
  - Why it happens: Both operands are `int`, so Go performs integer division, truncating toward zero.
  - Fix: Convert to `float64` first: `float64(7) / float64(2)` → `3.5`.

- Mistake: `"hello" + 42` expecting `"hello42"`.
  - Why it happens: Go does not implicitly convert between types. String + int is a compile error.
  - Fix: Use `fmt.Sprintf("%s%d", "hello", 42)` or `strconv.Itoa(42)`.

- Mistake: `x = y = 1` does not compile.
  - Why it happens: Go does not support chained assignment expressions. Assignment is a statement, not an expression.
  - Fix: Write two statements: `y = 1; x = y`.

- Mistake: Using `&&` and `||` thinking both sides always evaluate.
  - Why it happens: `&&` short-circuits: if the left operand is `false`, the right is never evaluated. `||` short-circuits: if the left is `true`, the right is skipped.
  - Fix: Never place expressions with side effects (like function calls modifying state) on the right of short-circuit operators unless you intend the short-circuit behavior.

## Debugging walkthrough

Consider this broken code:

```go
package main

import "fmt"

func main() {
	price := 19
	quantity := 3
	total := price * quantity / 100 * 5
	fmt.Println("Discount:", total)
}
```

**Symptom**: Prints `Discount: 0` instead of the expected `2` (5% of 57 = 2.85, truncated to int).

**Investigation**: Add `fmt.Printf` to trace each subexpression:

```go
fmt.Println("price * quantity =", price*quantity)           // 57
fmt.Println("price * quantity / 100 =", price*quantity/100)  // 0 (integer division!)
```

**Root cause**: `*` and `/` have the same precedence and evaluate left to right. `price * quantity / 100 * 5` evaluates as `((price * quantity) / 100) * 5` → `(57 / 100) * 5` → `0 * 5` → `0`. The integer division `57 / 100` truncates to `0`.

**Fix**: Use parentheses to control order, and convert to float for fractional percentages:

```go
total := price * quantity * 5 / 100  // (57 * 5) / 100 = 285 / 100 = 2
```

Or for precision:

```go
discount := float64(price*quantity) * 0.05
fmt.Printf("Discount: %.2f\n", discount)
```

## Production notes

In real Go codebases, expressions appear in hot paths, configuration parsing, and business logic. Key considerations:

- **Integer overflow** is silent. `math.MaxInt64 + 1` wraps to `math.MinInt64`. Use `math.Add64` from `math/bits` for overflow-checked arithmetic in critical paths like financial calculations.
- **Operator precedence bugs** are common in code reviews. When in doubt, add parentheses. They cost nothing at runtime and make intent explicit.
- **Short-circuit evaluation** is relied upon for safe nil checks: `if p != nil && p.Name != ""`. Reversing the order causes a nil pointer dereference panic.
- **Avoid magic numbers**: replace bare literals like `86400` with named constants `const secondsPerDay = 86400`. This makes expressions self-documenting.

## Performance implications

Arithmetic on built-in types (`int`, `float64`, `bool`) compiles to single CPU instructions. There is no allocation overhead. However:

- Integer division is 10-30x slower than addition or multiplication on modern CPUs.
- Bounds-checked array/slice access has a negligible branch, predictable by the branch predictor.
- String concatenation with `+` inside a loop allocates on every iteration. Use `strings.Builder` for repeated concatenation.

## Practice task

Write a function `compute(a, b int, op string) (int, error)` that:
- Returns `a + b` when `op` is `"add"`.
- Returns `a - b` when `op` is `"sub"`.
- Returns `a * b` when `op` is `"mul"`.
- Returns `a / b` when `op` is `"div"`, or an error if `b == 0`.
- Returns `a % b` when `op` is `"mod"`, or an error if `b == 0`.
- Returns `0, fmt.Errorf("unknown operator: %s", op)` for any other `op`.

Then write a `main()` that calls this function with at least five different inputs and prints results. Run it and verify the output.

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/01-values-and-expressions
go test ./curriculum/modules/03-programming-fundamentals/lessons/01-values-and-expressions
```

The existing tests verify `evaluate` with basic arithmetic. After completing the practice task, add tests for your `compute` function covering all operators and error cases.

## Review questions

1. What is the result of `10 + 3 * 2` in Go? Why?
2. Write an expression that computes 15% of 200 using only integer arithmetic.
3. Given `x := 10; y := 3; fmt.Println(x / y * y)`, does the output equal `x`? Explain.
4. What happens when you write `"count: " + 5` in Go? How would you fix it?
5. Why might `a/b*c` produce a different result than `a*c/b` with integer operands? Give an example.

## NEXT UP

Variables and how they bind names to values with type inference and zero-value initialization.
