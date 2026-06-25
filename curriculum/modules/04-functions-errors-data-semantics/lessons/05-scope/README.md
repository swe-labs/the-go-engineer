# Scope

## Learning objective

Identify the scope of variables in Go — package scope, function scope, and block scope — and explain how lexical scoping and shadowing affect variable visibility and lifetime.

## Why this matters

Scope determines where a variable is visible and usable. Bugs caused by shadowing, unintended variable reuse, or accessing variables outside their scope are common in every codebase. Understanding scope is essential for writing correct code, especially as programs grow beyond a single file.

## Mental model

Scope is the region of source code where a name (variable, function, type) is visible. Think of it as concentric circles: the innermost block can see outwards, but outer blocks cannot see inwards.

- **Package scope**: visible everywhere in the package (all files).
- **Function scope**: visible from declaration to the end of the function.
- **Block scope**: visible only within `{ }` braces.

Lexical scoping means the scope is determined by the source code structure at compile time, not by runtime call paths.

## Core idea

**Package scope**: Declarations at the top level (outside any function) are visible in all files of the same package. This includes variables, constants, types, and functions.

```go
package main

var packageVar = "visible everywhere in main"

func hello() {
    fmt.Println(packageVar) // ok
}
```

**Function scope**: Variables declared in a function are visible from their declaration to the end of the function.

**Block scope**: Variables declared inside `{ }` (if, for, switch, etc.) are visible only within that block. Each `case` in a `switch` or `select` is also a block.

**Shadowing**: A variable declared in an inner scope with the same name as one in an outer scope "shadows" (hides) the outer one. The outer variable is inaccessible within that inner block.

**Scope vs lifetime**: Scope is compile-time visibility. Lifetime is how long a variable exists at runtime. A variable can be out of scope but still alive (e.g., closure capturing a local variable).

## Under the hood

The Go compiler builds a symbol table that maps names to declarations. When the compiler encounters a name, it searches from the innermost scope outward. The first match wins — that's why shadowing occurs.

At runtime, local variables typically live on the stack (function scope) or in registers. Variables captured by closures escape to the heap, extending their lifetime beyond their scope. The compiler performs escape analysis to decide this.

## How Go uses it

- **Package-level variables** for configuration, shared state, or caches (use sparingly).
- **Local variables** for computation within a function.
- **Block-scoped variables** in `if`, `for`, `switch` to limit visibility.
- **Shadowing** is common (and often accidental) in `if err := ...` — you think you are assigning to an outer `err`, but instead you create a new one.

## Go example

```go
package main

import "fmt"

// Package-level variable — visible in all functions in this package.
var level = "package"

func main() {
	fmt.Println("outer level:", level) // "package"

	// Function scope variable shadows package-level.
	level := "function"
	fmt.Println("main level:", level) // "function"

	if true {
		// Block scope variable shadows function scope.
		level := "block"
		fmt.Println("block level:", level) // "block"

		// Inner block — no new shadow.
		{
			fmt.Println("inner block level:", level) // "block" (inherits from outer block)
		}
	}

	// Back to function scope — "block" is gone.
	fmt.Println("after block level:", level) // "function"

	// Demonstrating shadowing with if-scoped err.
	val := 10
	if val > 5 {
		result := val * 2 // block-scoped
		fmt.Println("if result:", result)
	}
	// fmt.Println(result) // compile error: result undefined here

	_ = val
}

func helper() {
	fmt.Println("helper level:", level) // "package" — the original package-level var
}
```

## Step-by-step execution

For the shadowing example at `level := "function"` in `main`:

1. Go encounters the identifier `level` in the assignment `level := "function"`.
2. It searches scopes: innermost (function body) finds the new declaration — no outer search needed.
3. This `level` shadows the package-level `level` for the rest of `main`.
4. Inside the `if` block, another `level := "block"` shadows the function-level `level`.
5. After the `if` block, the block-scoped `level` is out of scope. `level` refers to the function-scoped one again.
6. In `helper()`, there is no local `level`, so Go finds the package-level `level`.

## Common mistakes

- **Accidental shadowing with `:=`**: Inside an `if` or `for`, `err := ...` creates a new `err`, shadowing an outer one. Use `=` (with `err` already declared) to reuse the outer variable.
- **Using loop variable after loop**: Go 1.22 fixed this, but in older versions, `for i := 0; i < 3; i++` declared a single `i` reused across iterations. Capturing `&i` in a goroutine caused bugs.
- **Assuming if/for variable is accessible outside**: `if x := compute(); x > 0 { ... }` — `x` is only visible inside the `if` block and its `else`.
- **Shadowing needed packages**: `fmt := "hello"` shadows the `"fmt"` package import, preventing any use of `fmt.Println` in that scope.

## Debugging walkthrough

This code has a subtle shadowing bug:

```go
package main

import (
	"errors"
	"fmt"
)

func main() {
	err := errors.New("outer error")
	fmt.Println("before:", err)

	if true {
		err := errors.New("inner error")
		fmt.Println("inside:", err)
	}

	fmt.Println("after:", err) // still prints "outer error"!
}
```

**Symptom**: The `after` print shows `"outer error"` even though "inner error" was assigned.

**Root cause**: Inside the `if` block, `err := ...` creates a new variable `err` that shadows the outer `err`. The outer `err` is never modified.

**Fix**: Use `=` instead of `:=` to reuse the existing outer variable:

```go
if true {
    err = errors.New("inner error") // assignment, not declaration
}
```

## Production notes

- **Minimise package-level variables**: They increase coupling and make testing harder. Prefer passing state explicitly through function parameters.
- **Limit variable scope**: Declare variables as close to their use as possible. This reduces the cognitive load of tracking where a variable is modified.
- **Shadow detection**: The `go vet` tool can detect some shadowing cases. Use `go vet ./...` regularly.
- **Named returns**: A named return variable has function scope. If you redeclare it with `:=` inside the function, you shadow it — the named return is no longer accessible.

## Performance implications

- Package-level variables live for the entire program lifetime — they cannot be garbage collected.
- Block-scoped variables are typically stack-allocated and freed when the block exits.
- Variables captured by closures escape to the heap, increasing GC pressure.
- Scoping has zero runtime cost — it is purely a compile-time concept.

## Practice task

Write a function `demoScope() string` that:
- Declares a package-level variable `name = "global"`
- Declares a function-level `name` with value `"local"`
- Inside an `if true` block, declares a block-level `name` with value `"block"`
- Returns the concatenation of all three values in order (block, local, global) — but you cannot access them directly once shadowed. Hint: assign to different variables before shadowing.

## Tests / verification

```bash
go run ./curriculum/modules/04-functions-errors-data-semantics/lessons/05-scope
go test ./curriculum/modules/04-functions-errors-data-semantics/lessons/05-scope
```

## Review questions

1. What is the difference between package scope and block scope?
2. What is variable shadowing, and how does `:=` cause accidental shadowing?
3. Can a block-scoped variable be accessed after the block ends?
4. If a function declares `x := 1` and an inner `if` declares `x := 2`, does the outer `x` change? Explain.
5. What does `go vet` report about shadowing?

## NEXT UP

Call stack mental model — stack frames, push/pop, and LIFO behaviour.
