# LB.1 Variables

## Mission

Learn how Go declares variables and why every type has a predictable zero value.

## Prerequisites

- `GT.4` development environment

## Mental Model

A variable is a named slot that holds a value while the program runs.
Go gives you three common declaration shapes:

1.  `var name string`: Explicit type, starts as a zero value.
2.  `var name = "value"`: Type inferred from the right-hand side.
3.  `name := "value"`: Short declaration (local only), type inferred.

Every declared variable also starts in a known **Zero State**.

## Visual Model

```mermaid
graph LR
    A["Variable Name"] --- B["Type Constraint"]
    B --- C["Allocated Memory Slot"]
    C --> D["Zero Value (initial)"]
    D --> E["Assigned Value (updated)"]
```

## Machine View

When a variable is declared, Go reserves space in memory (usually on the stack for simple locals) and initializes it to the type's zero value. This guarantee prevents "uninitialized memory" bugs where a variable could contain random leftover data from a previous execution.

> [!NOTE]
> For a deeper look at how Go manages memory for variables, revisit [HC.3 Memory Basics](../../00-how-computers-work/3-memory-basics/README.md).

## Run Instructions

```bash
go run ./02-language-basics/1-variables
```

## Code Walkthrough

-   **`var greeting string`**: Declares a string. The zero value is `""` (empty string).
-   **`greeting = "Hello"`**: Assignment updates the existing slot.
-   **`var count int`**: Zero value is `0`.
-   **`var isRunning bool`**: Zero value is `false`.
-   **`firstName, lastName := "John", "Doe"`**: Short declaration creates and initializes multiple variables at once.

## Common Mistakes

### Unused Variables
Unlike Python or JavaScript, Go will **not** let you keep unused local variables. This is a common source of frustration for beginners, but it's designed to keep production code clean.
- **Error:** `name declared and not used`
- **Solution:** Either use the variable or delete it.

### Short Declaration vs. Package Level
The `:=` syntax is for **local** variables (inside functions) only. If you try to use it at the package level (outside of any function), the compiler will complain.
- **Error:** `non-declaration statement outside function body`
- **Solution:** Use `var` for package-level variables.

## Try It
> The compiler strictly enforces code quality. If you declare a local variable but never use it, the compiler will refuse to build the program. This was introduced in [GT.6 Reading Compiler Errors](../../01-getting-started/6-reading-compiler-errors/README.md).

## Try It

1.  Change one assigned value in `main.go` and rerun.
2.  Add a new variable using `:=`.
3.  Declare a variable and leave it unused to see the compiler error.

## In Production

Predictable initialization is critical for reliability. In production systems, you want to know exactly what state your system starts in. Go's zero-value guarantee and unused-variable checks reduce "dead code" and hidden bugs that often plague larger codebases.

## Thinking Questions

1.  Why is a guaranteed zero value safer than leaving memory uninitialized?
2.  When might explicit `var name string` be clearer than `name := ""`?
3.  How does the "unused variable" rule help you during refactoring?

## Next Step

Next: `LB.2` -> [`02-language-basics/2-constants`](../2-constants/README.md)
