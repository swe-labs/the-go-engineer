# TI.16 Complex Generic Constraints

## Mission

Learn advanced constraint patterns including parameterized constraints, using interfaces as constraints, and creating reusable generic utilities.

## Why This Lesson Exists Now

You know basic generics with simple constraints like `int | float64`. But real-world code often needs more sophisticated constraints—constraints that require methods, parameterized types, or multiple interface requirements.

## Prerequisites

- `TI.9` generics

## Mental Model

Think of a vending machine that accepts only certain payment methods. The constraint is not just "some type"—it's "anything with Pay() method that returns error." Similarly, generic constraints can require methods, not just type identity.

## Visual Model

```go
// Constraint requiring methods
type Adder interface {
    Add(other int) int
}

// Constraint requiring multiple interfaces  
type Serializer interface {
    fmt.Stringer
    json.Marshaler
}
```

## Run Instructions

```bash
go run ./01-foundations/06-types-and-interfaces/16-complex-generics
```

## Code Walkthrough

### Interface as constraint

Interfaces can be constraints—anything implementing the interface works.

### Multiple interface constraints

Use embedded interfaces to require multiple behaviors.

### Comparable constraint

The built-in `comparable` constraint allows equality operators.

## Try It

1. Create a constraint that requires both String() and a custom method.
2. Use the comparable constraint to create a generic key-value pair.
3. Build a constraint for numeric types with multiple operations.

## Production Relevance

Complex constraints are used in real Go code for data structures, serialization, and anywhere you need type-safe generic utilities.

## Next Step

Continue to `TI.17` generic data structures.