# LB.4 Application Logger

## Mission

Build a small application logger that combines variables, constants, iota, and custom types into one readable, runnable program.

This is the milestone exercise for the language basics section.

## Why This Lesson Exists Now

This is the first real integration exercise. So far, each lesson taught one concept in isolation. Now the learner combines everything they've learned:

- Variables for storing data
- Constants for fixed values
- iota for ordered enums
- Named types for type safety
- Methods for behavior

This mirrors real development: you don't use one concept in isolation, you compose them.

## Prerequisites

- `GT.4` development environment
- `LB.1` variables
- `LB.2` constants
- `LB.3` enums with iota

## What You Will Build

Implement a logger that:

1. Defines a `LogLevel` type based on `int`
2. Creates level constants with `iota`
3. Maps each level to a readable name
4. Implements a `String()` method with bounds checking
5. Prints readable log level output through a helper function

## Run Instructions

Run the completed solution:

```bash
go run ./01-foundations/02-language-basics/4-application-logger
```

Run the starter file:

```bash
go run ./01-foundations/02-language-basics/4-application-logger/_starter
```

## Solution Walkthrough

### Define LogLevel type

```go
type LogLevel int
```

This creates a named type based on int for type-safe enums.

### Create iota constants

```go
const (
    LevelTrace   LogLevel = iota
    LevelDebug
    LevelInfo
    LevelWarning
    LevelError
)
```

iota auto-increments: 0, 1, 2, 3, 4.

### Map levels to names

```go
var levelNames = []string{"Trace", "Debug", "Info", "Warning", "Error"}
```

Parallel array indexing matches constant values to readable names.

### Implement String() method

```go
func (l LogLevel) String() string {
    if l < LevelTrace || l > LevelError {
        return "Unknown"
    }
    return levelNames[l]
}
```

Bounds checking prevents crashes on invalid values.

## Try It

1. Add a new log level constant and its name.
2. Test with an invalid level like 99.
3. Change the output format in `printLogLevel`.

## Verification Surface

Run the solution:

```bash
go run ./01-foundations/02-language-basics/4-application-logger
```

Expected output includes all log levels and "Unknown" for invalid input.

## Success Criteria

Your finished solution should:

- Define a named `LogLevel` type
- Use `iota` for ordered log level constants
- Convert levels into readable text safely
- Avoid crashing on invalid levels
- Print a few example log levels in `main()`

## Next Step

Continue to `03-control-flow` after this milestone.