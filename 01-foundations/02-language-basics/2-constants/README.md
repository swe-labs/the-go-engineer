# LB.2 Constants

## Mission

Learn how to declare constants in Go and understand why constants matter for safety and performance.

This lesson separates immutable values from ordinary variables.

## Why This Lesson Exists Now

After variables, the next fundamental question is: "What if a value should never change?"

In programming, some values are fixed by design — port numbers, error messages, configuration keys. Using variables for these would allow accidental changes that cause bugs. Constants make the intent explicit and the compiler enforces it.

This lesson builds on `LB.1` by showing how to declare values that are locked at compile time.

## Production Relevance

In production Go code, constants matter because:

- **Safety**: A constant cannot be modified at runtime, no matter what code runs
- **Clarity**: Code that uses constants says "this value is fixed by design"
- **Performance**: The compiler inlines constants, making them faster than variables

Real services use constants for server ports, API endpoints, error codes, and feature flags.

## Prerequisites

- `LB.1` variables

## Mental Model

Constants are immutable values set at compile time.

Unlike variables, constants cannot change after declaration. The compiler inlines their values directly into the machine code.

## Visual Model

```text
const Host = "localhost"
const Port = 8080
```

```text
Variable:  lives in memory, can change at runtime
Constant:  inlined at compile time, never changes
```

## Machine View

When you declare a constant, the Go compiler replaces every use of that constant with its actual value. No memory is allocated for constants at runtime.

This makes constants:
- Impossible to modify accidentally
- Faster to access (no memory lookup)
- Clear documentation that a value is fixed

## Run Instructions

```bash
go run ./01-foundations/02-language-basics/2-constants
```

## Code Walkthrough

### `const Host = "127.0.0.1"`

This declares a string constant. The value is fixed and cannot change.

### `const pi float64 = 3.1415926`

This declares a typed constant. It can only be used where a float64 is expected.

### `const ( ... )`

Grouped constants use parentheses. This keeps related values organized.

### `var isRunning bool`

Variables are declared with `var`. They start with zero values and can be reassigned.

## Try It

1. Change the `Host` constant and rerun.
2. Add this line inside `main()` and observe the compile error:

```go
Host = "different-value"  // This will fail - constants cannot be reassigned
```

3. Add a new constant to the group.

## Common Questions

- Why use constants instead of variables?
  - Constants document that a value should never change
  - They catch accidental changes at compile time
  - They can be slightly faster (inlined at compile time)

## Next Step

Continue to `LB.3` enums with iota.