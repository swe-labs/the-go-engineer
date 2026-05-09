# TI.12 Type Assertions

## Mission

Learn how to extract concrete types from interface values using type assertions and the comma-ok pattern.

## Why This Lesson Exists Now

You have learned that interfaces can hold any type that implements their methods. Sometimes you need to extract the underlying concrete type to access type-specific fields or methods. Type assertions let you do this safely.

## Prerequisites

- `TI.3` interfaces

## Mental Model

Think of a vending machine that returns items in generic boxes. You know it might contain a soda, a snack, or water, but you need to open the box (type assert) to find out which one you got and access its specific properties.

## Visual Model

```mermaid
graph TD
    A["data"] --> B["type definition"]
    B --> C["methods or interface behavior"]
```
```go
// Basic type assertion
var i interface{} = "hello"
s := i.(string)  // s is now "hello"

// Safe type assertion with comma-ok
s, ok := i.(string)  // ok is true if assertion succeeded
```

## Machine View

A type assertion checks if the underlying concrete type matches the expected type. If it does, the value is returned. If not, the panic can be avoided using the comma-ok form.

## Run Instructions

```bash
go run ./04-types-design/12-type-assertions
```

## Code Walkthrough

### Basic assertion

`value.(Type)` panics if the type doesn't match.

### Comma-ok pattern

`value, ok := value.(Type)` returns ok=false instead of panicking.

### Type switch

Use `value.(type)` inside a switch to handle multiple types.

## Try It

1. Create an interface{} holding different types and assert each.
2. Use comma-ok to safely handle failed assertions.
3. Combine type assertions with type switches.

## ⚠️ In Production
Type assertions are used when reading from generic containers, handling dynamic data, and working with interface types from external sources.

## 🤔 Thinking Questions

1. What problem is this lesson trying to solve?
2. What would change if you removed this idea from the program?
3. Where do you expect to see this pattern again in real Go code?
## Next Step

Continue to `TI.13` nil interfaces.