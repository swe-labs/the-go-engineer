# Variables

## Learning objective

Declare and initialize Go variables using `var` and `:=`, understand zero-value initialization, shadowing, and scope rules, and write safe code that avoids common variable pitfalls.

## Why this matters

Every nontrivial program stores and manipulates data through variables. A misplaced shadow, an unintended zero value, or a scoping bug can introduce subtle production failures that are hard to reproduce. Professional Go codebases rely on clear variable discipline: knowing when to use `var` vs `:=`, understanding that uninitialized variables are never truly "empty," and reasoning about where a name is visible. These skills prevent entire categories of bugs before they happen.

## Mental model

A variable is a named box that holds a value. You create the box (declare), optionally put something inside (initialize), and later look inside or replace the contents (assign). The box has a fixed type: once you declare `var age int`, that box will only ever hold integers. The short declaration `:=` creates the box and fills it in one step. Scope is like a set of nested rooms — a name declared in an inner room temporarily hides (shadows) the same name in an outer room.

## Core idea

Go provides two declaration forms:

```
var name type = value     // full form
var name type             // zero-value initialization
var name = value          // type inferred from value
name := value             // short declaration (inside functions only)
```

The `var` form works at both package and function level. The `:=` form is a shorthand available only inside functions; it declares and initializes in one step.

Multiple variables can be declared together:

```go
var x, y int = 1, 2
a, b := "hello", true
```

Multiple assignment lets you swap values without a temporary variable:

```go
a, b = b, a
```

Shadowing occurs when a declaration in an inner scope uses the same name as a declaration in an outer scope. The inner name takes precedence until the inner scope ends.

## Under the hood

At the compiler level, each variable declaration allocates storage — either on the stack (for local variables whose lifetime is contained within the function) or on the heap (for variables that escape). The compiler tracks each variable's scope by maintaining a symbol table. When a name is referenced, the compiler walks outward through nested scopes until it finds a matching declaration. Shadowing works because the inner declaration is a separate entry in the symbol table that takes priority.

Zero-value initialization is a Go safety guarantee: every variable is set to a well-defined zero before first use. For integers this is `0`, for booleans `false`, for strings `""`, and for pointers/slices/maps/channels/ interfaces/functions `nil`. This means no "undefined variable" bugs at runtime — a variable always holds a legal value of its type.

## How Go uses it

- **Package-level `var`** for configuration, globals, and shared state (rare in idiomatic Go, but present).
- **Function-level `:=`** is the most common form — roughly 90% of local variable declarations in Go use `:=`.
- **Multiple return values** are idiomatic: `val, err := doSomething()`. This is Go's primary error-handling pattern.
- **Blank identifier `_`** discards values you don't need: `count, _ = fmt.Println(...)`.
- **Short variable redeclaration** is allowed as long as at least one variable on the left is new: `a, b := 1, 2; a, c := 3, 4` — `a` is reused, `c` is new.

## Go example

```go
package main

import "fmt"

var packageLevel = "visible everywhere in the package"

func main() {
	// var declaration with zero value
	var count int
	fmt.Println("zero value:", count)

	// var with initializer
	var msg string = "hello"
	fmt.Println("var with init:", msg)

	// short declaration
	price := 42
	fmt.Println("short decl:", price)

	// multiple assignment / swap
	x, y := 10, 20
	fmt.Println("before swap:", x, y)
	x, y = y, x
	fmt.Println("after swap:", x, y)

	// shadowing
	shadow := "outer"
	{
		shadow := "inner"
		fmt.Println("inside block:", shadow)
	}
	fmt.Println("outside block:", shadow)
}
```

## Step-by-step execution

For the swap `x, y := 10, 20; x, y = y, x`:

1. Declare `x` with value `10` and `y` with value `20`.
2. Evaluate the right side of the assignment: read `y` → `20`, read `x` → `10`.
3. Assign `20` to `x` and `10` to `y`.
4. Result: `x = 20`, `y = 10`.

For the shadowing example:

1. Declare `shadow` in `main` with value `"outer"`.
2. Enter the inner block (new scope).
3. Declare a new `shadow` in the inner scope with value `"inner"`.
4. Print `shadow` — the compiler picks the innermost declaration → `"inner"`.
5. Exit the inner block — the inner `shadow` is gone.
6. Print `shadow` — only the outer declaration is visible → `"outer"`.

## Common mistakes

- Mistake: `x := 1` inside a loop creates a new `x` on every iteration.
  - Why it happens: `:=` always declares. If you intend to reuse a variable, use `=` instead.
  - Fix: Declare `x` before the loop and assign with `=` inside.

- Mistake: Expecting `:=` to work at package level.
  - Why it happens: `:=` is only valid inside function bodies.
  - Fix: Use `var` at package level.

- Mistake: Assuming shadowing modifies the outer variable.
  - Why it happens: Inner `:=` creates a new variable that hides but does not change the outer one.
  - Fix: Use direct assignment `=` if you intend to modify the outer variable; avoid shadowing by choosing distinct names.

- Mistake: `var x int = 0` is redundant.
  - Why it happens: Zero-value initialization already sets `x` to `0`.
  - Fix: Use `var x int` to rely on the zero value, or use `:=` with an explicit value.

## Debugging walkthrough

Consider this broken code:

```go
package main

import "fmt"

func main() {
	var discount float64
	price := 100
	discount := 0.1
	total := float64(price) * (1 - discount)
	fmt.Println(total)
}
```

**Symptom**: Compile error: `no new variables on left side of :=`.

**Investigation**: The compiler message points to `discount := 0.1`. Looking at the code, `discount` was already declared with `var discount float64` on the line above. The `:=` operator requires at least one new variable on the left side.

**Root cause**: `discount` appears only on the left of `:=` and it was already declared. The short declaration operator cannot redeclare a variable without introducing at least one new variable.

**Fix**: Change `discount := 0.1` to `discount = 0.1` — a simple assignment that overwrites the zero value:

```go
var discount float64
price := 100
discount = 0.1
total := float64(price) * (1 - discount)
```

Alternatively, combine declaration and initialization on one line and remove the `var` line:

```go
price := 100
discount := 0.1
total := float64(price) * (1 - discount)
```

## Production notes

- **Prefer `:=` inside functions** — it is concise and type inference reduces repetition.
- **Use `var` at package level** — package-level `:=` is not allowed, and `var` makes the zero-value behavior explicit.
- **Beware of shadowing `err`** — inside an `if` block, `if err := doSomething(); err != nil` shadows an outer `err`, making it inaccessible. This is often intentional and correct (the `err` scoped to the `if`), but can confuse newcomers.
- **Avoid package-level variables** for mutable state. They create implicit coupling between functions and make testing difficult. Prefer explicit function parameters.
- **Group related declarations** with `var (...)`. This improves readability and documents that the variables belong together.

## Performance implications

- Stack-allocated local variables are essentially free — the compiler allocates space by adjusting the stack pointer at function entry.
- Variables that escape to the heap (e.g., returned pointers, closures capturing locals) incur allocation and garbage-collection overhead.
- `:=` and `var` produce identical machine code. There is no performance difference — choose based on readability.
- Shadowing has zero runtime cost; the compiler resolves names at compile time.
- Multiple assignment compiles to the same instructions as separate assignments — the swap via `a, b = b, a` is compiled into register-level moves, often more efficient than using a temporary.

## Practice task

Write a function `calculate(r, h float64) (area, volume float64)` that computes the surface area and volume of a right circular cylinder. Use `r` for radius and `h` for height. Formulas: area = 2πr(r + h), volume = πr²h. Use `math.Pi`.

Then in `main()`:
- Declare three sets of radius and height values (e.g., `(3, 5)`, `(1.5, 2)`, `(10, 1)`).
- For each set, call `calculate` and print the result.
- Use short declaration, multiple assignment, and at least one `var` declaration with zero value (then assign before use).

## Tests / verification

```bash
go run ./curriculum/modules/03-programming-fundamentals/lessons/02-variables
go test ./curriculum/modules/03-programming-fundamentals/lessons/02-variables
```

## Review questions

1. Explain the difference between `var x int` and `x := 0`. Are the resulting variables different in any way?
2. Given `a, b := 1, 2; a, b = b, a`, what are the values of `a` and `b` after execution? Explain step by step.
3. Debug: `var x = 10; if true { x := 20 }; fmt.Println(x)` — what does it print and why? How would you make it print `20`?
4. Tradeoff: When would you choose `var x int` over `x := 0` in a production codebase? Consider readability, intent, and code review concerns.
5. What happens if you write `x, y := 1, 2; x, y := 3, 4` in Go? Does it compile? Explain.

## NEXT UP

Types — the set of basic types in Go (int, float64, string, bool, byte, rune), type inference with `:=`, and the distinction between typed and untyped constants.
