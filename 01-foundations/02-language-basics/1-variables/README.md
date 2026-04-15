# LB.1 Variables

## Mission

Learn how to declare and use variables in Go, and understand Go's zero value system.

This lesson teaches the three ways to create variables in Go.

## Why This Lesson Exists Now

Variables are the foundation of every program. Before a learner can store user input, calculate values, or make decisions, they need to understand how Go holds and updates data.

This lesson builds on `GT.4` where the learner first ran a program. Now they learn how the program remembers values.

## Production Relevance

In production Go code, variables and zero values matter because:

- **Predictability**: Every type has a known zero value, so you never wonder "what's in this variable before I assign it?"
- **Safety**: Go's compiler catches unused variables, preventing dead code from hiding bugs
- **Clarity**: The three declaration styles (`var`, `:=`, `var =`) communicate intent: explicit type, quick local, or inferred

Real services use variables for configuration, counters, state machines, and user data.

## Mental Model

In Go, every variable has a type, and every type has a predictable zero value.

The three declaration styles are:

1. `var name string` — explicit type, starts with zero value
2. `name := "value"` — short declaration, type inferred
3. `var name = "value"` — var with inference

## Visual Model

```text
var greeting string   →  ""      (empty string)
var count int         →  0
var isActive bool     →  false
var price float64     →  0.0
```

```text
greeting := "hello"   →  inferred as string
age := 25             →  inferred as int
```

## Machine View

When a variable is declared but not assigned a value, Go automatically initializes it to its type's zero value.

This is a safety feature. In some languages, uninitialized memory contains random garbage data. Go guarantees every variable starts in a known, safe state.

## Run Instructions

```bash
go run ./01-foundations/02-language-basics/1-variables
```

## Code Walkthrough

### `var greeting string`

This declares a variable with an explicit type. The value starts as an empty string `""`.

### `greeting = "Hello, world!"`

This assigns a value to the variable. Now the variable holds meaningful data.

### `var count int`

This declares an integer variable. Its zero value is `0`.

### `var isRunning bool`

This declares a boolean variable. Its zero value is `false`.

### `firstName, lastName := "John", "Doe"`

This uses short declaration to create two variables at once. Go infers both are strings.

### `var year = 2025`

This uses `var` with inference. The type is inferred from the value.

### Unused variable rule

Go refuses to compile programs with unused variables. This catches mistakes early.

## Try It

1. Change `count = 10` to a different number and rerun.
2. Add a new variable using short declaration.
3. Declare a variable but never use it — observe the compile error.

## Common Questions

- Why does Go have three ways to declare variables?
  - `var` for explicit control and zero values
  - `:=` for quick local declarations
  - `var` with inference when you want clarity with inference

- What happens if I don't assign a value?
  - The variable gets its type's zero value automatically.

## Next Step

Continue to `LB.2` constants.