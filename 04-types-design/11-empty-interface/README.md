# TI.11 Empty Interface

## Mission

Learn about the empty interface (any) and when it is appropriate to use it.

## Why This Lesson Exists Now

You have learned about interfaces with specific methods. But sometimes you need a type that can hold anything. The empty interface (`any` or `interface{}`) is that type.

## Prerequisites

- `TI.3` interfaces

## Mental Model

Think of a moving box. A specialized box might only hold books, or only clothes. But a "any" box can hold anythingâ€”dishes, clothes, books, toys. The empty interface is like that flexible box.

## Visual Model

```mermaid
graph TD
    A["data"] --> B["type definition"]
    B --> C["methods or interface behavior"]
```
```text
// These are equivalent:
interface{}
any

// Can hold any type:
var x interface{} = 42
x = "hello"
x = []int{1, 2, 3}
x = map[string]int{"a": 1}
```

## Machine View

An empty interface has no methods, so every type satisfies it. This is powerful but loses type safety. Use it sparingly and try to use generic constraints instead when possible.

## Run Instructions

```bash
go run ./04-types-design/11-empty-interface
```

## Code Walkthrough

### Storing any type

The empty interface can hold any Go value.

### Type assertion needed

To use the value, you must assert its type.

### Use cases

JSON decoding, printing debug info, logging any value.

## Try It

1. Store different types in an empty interface variable.
2. Use type assertion to extract each type.
3. Use type switch to handle multiple types.

## Common Questions

- When to use empty interface vs generics?
  Generics preserve type safety. Use empty interface when you truly do not know the type at compile time.

- Is interface{} the same as any?
  Yes, in Go 1.18+, `any` is an alias for `interface{}`.

## ⚠️ In Production
Empty interface is used in JSON decoding (encoding/json), reflection, and generic logging utilities.

## 🤔 Thinking Questions

1. What problem is this lesson trying to solve?
2. What would change if you removed this idea from the program?
3. Where do you expect to see this pattern again in real Go code?
## Next Step

Continue to `TI.12` type assertions.
