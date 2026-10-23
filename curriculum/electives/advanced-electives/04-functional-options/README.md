# Functional options

## When you need this

Go constructors with many parameters create functions with long parameter lists that are hard to read, easy to misorder, and impossible to extend without breaking existing callers. Functional options solve this by using variadic option functions that set fields on a config struct after construction — each option is independent, optional, and discoverable.

## When you do not need this

When simpler concrete types or standard core library primitives suffice without adding abstraction overhead.

## Prerequisites

- elective-03

## Mental model

A functional option is a function that modifies a config struct. The constructor accepts an empty config struct, applies each option function to it, and uses the resulting config. Each option is a standalone function that sets exactly one field. Callers specify only the options they care about — unset fields use their zero values or defaults. This is equivalent to the builder pattern but uses functions instead of a builder struct.

## How it works

The functional options pattern uses a type Option func(*Config). The constructor NewServer(opts ...Option) creates a default Config, then calls each option in sequence: for _, opt := range opts { opt(&cfg) }. Each option is a closure: func WithTimeout(d time.Duration) Option { return func(c *Config) { c.Timeout = d } }. The closure captures the parameter (d) and returns a function that applies it. Options are applied in order, so later options override earlier ones for the same field.

## Go example

```go
package main

import "fmt"

func main() {
    fmt.Println("Elective: Functional options")
}
```

## Tradeoffs

- Each option function call is an indirect function call — overhead is negligible for typical configs (<10 options), measurable at scale (1000+ instantiations/s)
- The functional options pattern allocates no heap memory per option — option functions are usually closures that capture no variables, so they are allocated in the binary's read-only data section
- Variadic ...Option allocates a slice — for hot-path constructors, consider a fixed-size array or builder pattern to avoid allocation
- Compared to a config struct literal, functional options add ~5ns per option call — irrelevant for startup code, measurable in high-frequency object creation

## Common mistakes

- Using positional parameters for optional config — adding a new parameter breaks all callers and makes the function signature unreadable.
- Mutating the option closure's captured variable after creating the option — the option captures a reference, and changes to the variable affect the option's behavior.
- Not providing defaults for unset options — callers must specify every option or accept zero values that may not be valid.
- Applying options that conflict — two options set the same field to different values, and the last one wins (silently), which may surprise callers who expect an error.
- Using functional options for required parameters — required parameters should be explicit constructor arguments, not optional options.

## Practice task

The learner must implement a constructor with functional options: NewServer(addr string, opts ...Option) that accepts at least 3 options (WithTimeout, WithMaxConns, WithTLS), provide sensible defaults, and demonstrate that options are applied in order and that unset options use defaults.

## Production notes

Functional options are used in every Go library with configurable constructors: gRPC dial options, HTTP client options, database connection options, and logger configuration. Google Cloud Go client libraries use functional options for all service clients. Production Go services use functional options for server configuration, middleware setup, and client initialization.

## Related core lessons

- Module 16.03 — Generic data structures
- Module 16.15 — Demo preparation
