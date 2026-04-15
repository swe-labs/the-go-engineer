# LB.3 Enums with Iota

## Mission

Learn how Go creates enum-like values using iota, and understand type-safe enums.

Go does not have an `enum` keyword. Instead, it uses constants with iota for enum-like behavior.

## Why This Lesson Exists Now

After variables and constants, the next question is: "How do I represent a fixed set of related values?"

Real programs need to model things like log levels (debug, info, error), HTTP status categories, or user roles. In other languages, you might use enums. Go uses `iota` with named types.

This lesson builds on `LB.2` by showing how to create ordered, related constants that are type-safe.

## Production Relevance

In production Go code, iota and type-safe enums matter because:

- **Type safety**: A named type like `LogLevel` prevents passing random integers where a log level is expected
- **Organization**: iota auto-increments, so adding new values is easy and error-free
- **Readability**: String() methods make debug output human-readable

Real services use enums for log levels, HTTP methods, status codes, and configuration modes.

## Prerequisites

- `LB.2` constants

## Mental Model

`iota` is a special constant generator that produces sequential integers.

Inside a `const` block:
- First constant gets iota = 0
- Each subsequent constant gets iota + 1
- iota resets to 0 at each new const block

## Visual Model

```text
const (
    Red   = iota  // 0
    Blue           // 1
    Green          // 2
)
```

```text
type LogLevel int

const (
    LogError = iota  // 0
    LogWarn          // 1
    LogInfo          // 2
)
```

## Machine View

iota generates numeric constants automatically. For production code, wrap them in a named type for type safety.

When you create a custom type like `type LogLevel int`, the compiler prevents passing random integers where a LogLevel is expected.

## Run Instructions

```bash
go run ./01-foundations/02-language-basics/3-enums
```

## Code Walkthrough

### `const ( Sunday = iota + 1 ... )`

This starts iota at 1 instead of 0. Useful when 0 should mean "unset".

### `type LogLevel int`

This creates a named type. It provides type safety over raw integers.

### `const ( LogError LogLevel = iota ... )`

This attaches iota constants to the named type. Now you have type-safe enums.

### `func (l LogLevel) String() string`

This implements the stringer interface so fmt.Println prints readable names like "ERROR" instead of "0".

## Try It

1. Add a new log level to the enum and see the value auto-increment.
2. Create a custom type for a different enum.
3. Print an invalid enum value and see the default case.

## Common Questions

- Why doesn't Go have enums?
  - Go prefers simplicity. iota with named types achieves the same goal with less syntax.

- When should I use iota+1 instead of iota?
  - When 0 should represent "unset" or "invalid" rather than a valid enum value.

## Next Step

Continue to `LB.4` application logger exercise.