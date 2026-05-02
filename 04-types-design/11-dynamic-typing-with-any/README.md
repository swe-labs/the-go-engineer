# TI.11 Dynamic Typing with Any

## Mission

- Utilize the `any` keyword to handle data of unknown or dynamic types.
- Recover concrete types safely using type assertions and type switches.
- Understand and avoid the "Typed Nil" pitfall in interface comparisons.

## Prerequisites

- `TI.3` Interfaces

## Mental Model

The **Empty Interface** (`interface{}` or `any`) is an interface that specifies zero methods. Because every type satisfies at least zero methods, the empty interface can hold a value of any concrete type. This provides a mechanism for dynamic typing within Go's statically-typed system, which is useful at system boundaries where data shapes are not guaranteed.

## Visual Model

```mermaid
graph TD
    A[any / interface{}] --> B[Type Assertion]
    B --> C[Concrete Value]
```

## Machine View

The runtime representation of an interface value consists of a **Type Descriptor** and a **Data Pointer**. 

- **Nil Interface**: Both the type and data pointers are `nil`. (`v == nil` is true).
- **Typed Nil**: The type pointer is non-nil (pointing to a concrete type like `*emailNotifier`), but the data pointer is `nil`. (`v == nil` is false).

This is a common source of bugs. If a function returns a nil pointer of a concrete type into an interface return value, the caller's check against `nil` will fail because the interface now "contains" type information.

## Run Instructions

```bash
go run ./04-types-design/11-dynamic-typing-with-any
```

## Code Walkthrough

- **Type Assertions**: `v.(T)` extracts the concrete value. The comma-ok idiom (`v, ok := i.(T)`) should always be used to prevent runtime panics if the assertion fails.
- **Dynamic Switches**: Type switches provide a clean way to branch logic based on the dynamic type stored in `any`.
- **Constraint-Free**: Unlike specific interfaces, `any` provides no behavioral guarantees. You must always verify the underlying type before performing operations.

## Try It

1. In `main.go`, observe the output of the "Typed Nil" comparison.
2. Add a new function `process(v any)` that uses a type switch to handle `float64` and `[]string`.
3. Test the function with several different types and verify the `default` case handles unexpected inputs correctly.

## In Production

- **Logging**: Capturing arbitrary metadata fields of varying types.
- **Decoding**: Handling dynamic JSON payloads where fields may contain numbers, strings, or nested objects.
- **Messaging**: Implementing message brokers that pass arbitrary payloads between services.

## Thinking Questions

1. Why does Go allow an interface to be non-nil even if its underlying data is nil?
2. What are the performance and safety trade-offs of using `any` instead of specific interfaces?
3. In a modern Go codebase (Go 1.18+), when would you prefer `any` over a Generic type constraint?

## Next Step

Next: `TI.12` -> [`04-types-design/12-functional-options`](../12-functional-options/README.md)
