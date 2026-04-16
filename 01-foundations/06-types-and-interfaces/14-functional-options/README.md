# TI.14 Functional Options

## Mission

Learn the functional options pattern—a common Go pattern for building configurable APIs without requiring many constructor parameters.

## Why This Lesson Exists Now

When a type has many optional fields, passing all of them to a constructor becomes unwieldy. The functional options pattern lets callers customize only what they need using small, composable functions.

## Prerequisites

- `TI.2` methods

## Mental Model

Think of ordering a pizza. You could have a constructor with 20 parameters (crust, sauce, cheese, toppings, size, etc.). Or you could have `WithExtraCheese()`, `WithPepperoni()`, `LargeSize()` functions that you chain together. Much cleaner!

## Visual Model

```go
// Without options: too many parameters
NewServer("web", "us-east", 4, 16, true, false, "linux", "10.0.0.1", ...)

// With functional options
NewServer(
    WithName("web"),
    WithRegion("us-east"),
    WithCPUs(4),
)
```

## Run Instructions

```bash
go run ./01-foundations/06-types-and-interfaces/14-functional-options
```

## Code Walkthrough

### Option type

Define a function type that modifies a config struct.

### Option function

Each option function returns an Option that gets applied.

### WithDefault pattern

Use functional composition to build up configuration.

## Try It

1. Add a new option function for a missing field.
2. Create a server with multiple options chained together.
3. Make some options have default values.

## Production Relevance

Functional options are used throughout Go APIs—gRPC, Terraform provider, Cobra CLI, etc. Essential for building clean, extensible libraries.

## Next Step

Continue to `TI.15` method values.