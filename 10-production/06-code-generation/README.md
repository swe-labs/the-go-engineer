# Track CG: Code Generation

## Mission

Master the "Automated Hands." Learn how to use code generation tools to eliminate repetitive "Boilerplate" code while maintaining type safety and performance. Understand the `go generate` workflow and how to use industry-standard tools like **Mockery** (for testing) and **sqlc** (for database access) to make your development faster and less error-prone.

## Stage Ownership

This track belongs to [10 Production Operations](../README.md).

## Track Map

| ID | Type | Surface | Mission | Requires |
| --- | --- | --- | --- | --- |
| `CG.1` | Lesson | [go generate Primer](./1-go-generate) | Master the built-in `go:generate` directive. | entry |
| `CG.2` | Lesson | [Mockery Workflow](./2-mockery) | Generate clean, type-safe mocks for testing. | `CG.1` |
| `CG.3` | Lesson | [sqlc Workflow](./3-sqlc) | Generate typed Go code from raw SQL queries. | `CG.2` |

## Why This Track Matters

In Go, we prefer "Explicit" over "Implicit." However, being explicit often requires writing a lot of repetitive code (e.g., implementing an interface for testing or mapping SQL rows to structs).

1. **Productivity**: Code generation tools can write hundreds of lines of code for you in seconds.
2. **Reliability**: Tools don't make typos. If your interface changes, your mocks can be updated instantly with one command.
3. **Type Safety**: Unlike reflection-heavy ORMs, code generation produces standard Go code that is checked by the compiler.
4. **No Magic**: The generated code is just Go code. You can read it, debug it, and commit it to your repository.

## Next Step

Congratulations! You have completed Section 10. You are now ready to apply everything you've learned to a real-world, high-scale application. Continue to [11 Flagship](../../11-flagship).
